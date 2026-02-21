package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetService interface {
	GetBudgetsByMonthYear(month, year int) ([]BudgetSummary, error)
	GetBudgetSummary(month, year int) (*OverallSummary, error)
	ReallocateBudget(reallocation *BudgetReallocationRequest) (*BudgetReallocationResponse, error)
	GetReallocations(month, year int) ([]models.BudgetReallocation, error)
	CancelReallocation(id uuid.UUID) error
}

type BudgetSummary struct {
	Category        models.ExpenseCategory `json:"category"`
	AllocatedAmount float64                `json:"allocated_amount"`
	SpentAmount     float64                `json:"spent_amount"`
	RemainingAmount float64                `json:"remaining_amount"`
	PercentageUsed  float64                `json:"percentage_used"`
}

type OverallSummary struct {
	TotalIncome      float64 `json:"total_income"`
	TotalAllocated   float64 `json:"total_allocated"`
	TotalSpent       float64 `json:"total_spent"`
	TotalRemaining   float64 `json:"total_remaining"`
	UnallocatedFunds float64 `json:"unallocated_funds"`
}

type budgetService struct {
	budgetRepo       repositories.CategoryBudgetRepository
	incomeRepo       repositories.IncomeRepository
	allocationRepo   repositories.BudgetAllocationRepository
	reallocationRepo repositories.BudgetReallocationRepository
	categoryRepo     repositories.ExpenseCategoryRepository
	db               *gorm.DB
}

func NewBudgetService(
	budgetRepo repositories.CategoryBudgetRepository,
	incomeRepo repositories.IncomeRepository,
	allocationRepo repositories.BudgetAllocationRepository,
	reallocationRepo repositories.BudgetReallocationRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
	db *gorm.DB,
) BudgetService {
	return &budgetService{
		budgetRepo:       budgetRepo,
		incomeRepo:       incomeRepo,
		allocationRepo:   allocationRepo,
		reallocationRepo: reallocationRepo,
		categoryRepo:     categoryRepo,
		db:               db,
	}
}

func (s *budgetService) GetBudgetsByMonthYear(month, year int) ([]BudgetSummary, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	budgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	summaries := make([]BudgetSummary, 0, len(budgets))
	for _, budget := range budgets {
		percentageUsed := 0.0
		if budget.AllocatedAmount > 0 {
			percentageUsed = (budget.SpentAmount / budget.AllocatedAmount) * 100
		}

		summaries = append(summaries, BudgetSummary{
			Category:        budget.Category,
			AllocatedAmount: budget.AllocatedAmount,
			SpentAmount:     budget.SpentAmount,
			RemainingAmount: budget.RemainingAmount,
			PercentageUsed:  percentageUsed,
		})
	}

	return summaries, nil
}

func (s *budgetService) GetBudgetSummary(month, year int) (*OverallSummary, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	// Get all incomes for the month
	incomes, err := s.incomeRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalIncome float64
	for _, income := range incomes {
		totalIncome += income.Amount
	}

	// Get all allocations for the month
	allocations, err := s.allocationRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalAllocated float64
	for _, allocation := range allocations {
		totalAllocated += allocation.AllocatedAmount
	}

	// Get all budgets for the month
	budgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalSpent, totalRemaining float64
	for _, budget := range budgets {
		totalSpent += budget.SpentAmount
		totalRemaining += budget.RemainingAmount
	}

	summary := &OverallSummary{
		TotalIncome:      totalIncome,
		TotalAllocated:   totalAllocated,
		TotalSpent:       totalSpent,
		TotalRemaining:   totalRemaining,
		UnallocatedFunds: totalIncome - totalAllocated,
	}

	return summary, nil
}

type BudgetReallocationRequest struct {
	FromCategoryID uuid.UUID
	ToCategoryID   uuid.UUID
	Amount         float64
	Reason         string
	Month          int
	Year           int
}

type BudgetReallocationResponse struct {
	Reallocation       *models.BudgetReallocation `json:"reallocation"`
	FromCategoryBudget *CategoryBudgetInfo        `json:"from_category_budget"`
	ToCategoryBudget   *CategoryBudgetInfo        `json:"to_category_budget"`
}

type CategoryBudgetInfo struct {
	Remaining float64 `json:"remaining"`
}

func (s *budgetService) ReallocateBudget(req *BudgetReallocationRequest) (*BudgetReallocationResponse, error) {
	if req.Month < 1 || req.Month > 12 {
		return nil, errors.New("invalid month")
	}
	if req.Year < 2000 {
		return nil, errors.New("invalid year")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if req.FromCategoryID == req.ToCategoryID {
		return nil, errors.New("from and to category cannot be the same")
	}

	var reallocation *models.BudgetReallocation
	var fromBudget, toBudget *models.CategoryBudget

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check if both categories exist and are active
		fromCategory, err := s.categoryRepo.FindByID(req.FromCategoryID)
		if err != nil {
			return errors.New("from category not found")
		}
		if !fromCategory.IsActive {
			return errors.New("from category is not active")
		}

		toCategory, err := s.categoryRepo.FindByID(req.ToCategoryID)
		if err != nil {
			return errors.New("to category not found")
		}
		if !toCategory.IsActive {
			return errors.New("to category is not active")
		}

		// 2. Get budgets for both categories
		fromBudget, err = s.budgetRepo.FindByCategoryAndMonth(req.FromCategoryID, req.Month, req.Year)
		if err != nil {
			return errors.New("no budget found for from category in this month")
		}

		// Check if from budget has enough
		if fromBudget.RemainingAmount < req.Amount {
			return errors.New("insufficient budget in from category")
		}

		toBudget, err = s.budgetRepo.FindByCategoryAndMonth(req.ToCategoryID, req.Month, req.Year)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new budget for to category
				toBudget = &models.CategoryBudget{
					CategoryID:      req.ToCategoryID,
					Month:           req.Month,
					Year:            req.Year,
					AllocatedAmount: req.Amount,
					SpentAmount:     0,
					RemainingAmount: req.Amount,
				}
				if err := s.budgetRepo.Create(toBudget); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// Update existing to budget
			toBudget.AllocatedAmount += req.Amount
			toBudget.RemainingAmount += req.Amount
			if err := s.budgetRepo.Update(toBudget); err != nil {
				return err
			}
		}

		// 3. Update from budget
		fromBudget.RemainingAmount -= req.Amount
		fromBudget.AllocatedAmount -= req.Amount
		if err := s.budgetRepo.Update(fromBudget); err != nil {
			return err
		}

		// 4. Create reallocation record
		reallocation = &models.BudgetReallocation{
			FromCategoryID: req.FromCategoryID,
			ToCategoryID:   req.ToCategoryID,
			Amount:         req.Amount,
			Reason:         req.Reason,
			Month:          req.Month,
			Year:           req.Year,
		}

		if err := s.reallocationRepo.Create(reallocation); err != nil {
			return err
		}

		// Load relations
		reallocation.FromCategory = *fromCategory
		reallocation.ToCategory = *toCategory

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BudgetReallocationResponse{
		Reallocation: reallocation,
		FromCategoryBudget: &CategoryBudgetInfo{
			Remaining: fromBudget.RemainingAmount,
		},
		ToCategoryBudget: &CategoryBudgetInfo{
			Remaining: toBudget.RemainingAmount,
		},
	}, nil
}

func (s *budgetService) GetReallocations(month, year int) ([]models.BudgetReallocation, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	return s.reallocationRepo.FindByMonthYear(month, year)
}

func (s *budgetService) CancelReallocation(id uuid.UUID) error {
	// 1. Find the reallocation record
	reallocation, err := s.reallocationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reallocation not found")
		}
		return err
	}

	// 2. Reverse the budget changes in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get both budgets
		fromBudget, err := s.budgetRepo.FindByCategoryAndMonth(
			reallocation.FromCategoryID,
			reallocation.Month,
			reallocation.Year,
		)
		if err != nil {
			return errors.New("from category budget not found")
		}

		toBudget, err := s.budgetRepo.FindByCategoryAndMonth(
			reallocation.ToCategoryID,
			reallocation.Month,
			reallocation.Year,
		)
		if err != nil {
			return errors.New("to category budget not found")
		}

		// Restore from budget (add back the amount)
		fromBudget.AllocatedAmount += reallocation.Amount
		fromBudget.RemainingAmount += reallocation.Amount
		if err := s.budgetRepo.Update(fromBudget); err != nil {
			return err
		}

		// Restore to budget (remove the amount)
		toBudget.AllocatedAmount -= reallocation.Amount
		toBudget.RemainingAmount -= reallocation.Amount
		if err := s.budgetRepo.Update(toBudget); err != nil {
			return err
		}

		// Delete the reallocation record
		if err := s.reallocationRepo.Delete(id); err != nil {
			return err
		}

		return nil
	})
}
