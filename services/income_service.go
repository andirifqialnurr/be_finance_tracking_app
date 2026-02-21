package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeService interface {
	CreateIncome(income *models.Income) (*models.Income, []models.BudgetAllocation, error)
	GetIncomeByID(id uuid.UUID) (*models.Income, error)
	GetAllIncomes(limit, offset int) ([]models.Income, error)
	GetIncomesByMonthYear(month, year int) ([]models.Income, error)
	UpdateIncome(id uuid.UUID, income *models.Income) (*models.Income, []models.BudgetAllocation, error)
	DeleteIncome(id uuid.UUID) error
}

type incomeService struct {
	incomeRepo     repositories.IncomeRepository
	categoryRepo   repositories.ExpenseCategoryRepository
	budgetRepo     repositories.CategoryBudgetRepository
	allocationRepo repositories.BudgetAllocationRepository
	db             *gorm.DB
}

func NewIncomeService(
	incomeRepo repositories.IncomeRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	allocationRepo repositories.BudgetAllocationRepository,
	db *gorm.DB,
) IncomeService {
	return &incomeService{
		incomeRepo:     incomeRepo,
		categoryRepo:   categoryRepo,
		budgetRepo:     budgetRepo,
		allocationRepo: allocationRepo,
		db:             db,
	}
}

// CreateIncome creates a new income and automatically allocates budget to active categories
func (s *incomeService) CreateIncome(income *models.Income) (*models.Income, []models.BudgetAllocation, error) {
	var allocations []models.BudgetAllocation

	// Start transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Save income
		if err := s.incomeRepo.Create(income); err != nil {
			return err
		}

		// 2. Get all active categories
		categories, err := s.categoryRepo.FindAllActive()
		if err != nil {
			return err
		}

		if len(categories) == 0 {
			// No categories to allocate, just return income
			return nil
		}

		// 3. Calculate total monthly budget from all active categories
		var totalMonthlyBudget float64
		for _, cat := range categories {
			totalMonthlyBudget += cat.MonthlyBudget
		}

		if totalMonthlyBudget == 0 {
			return errors.New("total monthly budget is zero, cannot allocate")
		}

		// 4. Get month and year from income date
		month := int(income.Date.Month())
		year := income.Date.Year()

		// 5. Allocate budget to each category proportionally
		for _, category := range categories {
			// Calculate allocation based on category's budget proportion
			allocationAmount := (category.MonthlyBudget / totalMonthlyBudget) * income.Amount

			if allocationAmount <= 0 {
				continue
			}

			// Find or create category budget for this month/year
			budget, err := s.budgetRepo.FindByCategoryAndMonth(category.ID, month, year)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Create new budget
					budget = &models.CategoryBudget{
						CategoryID:      category.ID,
						Month:           month,
						Year:            year,
						AllocatedAmount: allocationAmount,
						SpentAmount:     0,
						RemainingAmount: allocationAmount,
					}
					if err := s.budgetRepo.Create(budget); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// Update existing budget
				budget.AllocatedAmount += allocationAmount
				budget.RemainingAmount += allocationAmount
				if err := s.budgetRepo.Update(budget); err != nil {
					return err
				}
			}

			// Create allocation record
			allocation := models.BudgetAllocation{
				IncomeID:        income.ID,
				CategoryID:      category.ID,
				AllocatedAmount: allocationAmount,
				Month:           month,
				Year:            year,
			}

			if err := s.allocationRepo.Create(&allocation); err != nil {
				return err
			}

			// Load category info for response
			allocation.Category = category
			allocations = append(allocations, allocation)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return income, allocations, nil
}

func (s *incomeService) GetIncomeByID(id uuid.UUID) (*models.Income, error) {
	return s.incomeRepo.FindByID(id)
}

func (s *incomeService) GetAllIncomes(limit, offset int) ([]models.Income, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.incomeRepo.FindAll(limit, offset)
}

func (s *incomeService) GetIncomesByMonthYear(month, year int) ([]models.Income, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}
	return s.incomeRepo.FindByMonthYear(month, year)
}

func (s *incomeService) UpdateIncome(id uuid.UUID, income *models.Income) (*models.Income, []models.BudgetAllocation, error) {
	var allocations []models.BudgetAllocation

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get existing income
		existingIncome, err := s.incomeRepo.FindByID(id)
		if err != nil {
			return err
		}

		oldMonth := int(existingIncome.Date.Month())
		oldYear := existingIncome.Date.Year()

		newMonth := int(income.Date.Month())
		newYear := income.Date.Year()
		newAmount := income.Amount

		// 2. Rollback old allocations
		oldAllocations, err := s.allocationRepo.FindByIncome(id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		for _, oldAlloc := range oldAllocations {
			budget, err := s.budgetRepo.FindByCategoryAndMonth(oldAlloc.CategoryID, oldMonth, oldYear)
			if err == nil {
				budget.AllocatedAmount -= oldAlloc.AllocatedAmount
				budget.RemainingAmount -= oldAlloc.AllocatedAmount
				if err := s.budgetRepo.Update(budget); err != nil {
					return err
				}
			}
			// Delete old allocation record
			if err := s.allocationRepo.Delete(oldAlloc.ID); err != nil {
				return err
			}
		}

		// 3. Update income
		income.ID = id
		if err := s.incomeRepo.Update(income); err != nil {
			return err
		}

		// 4. Create new allocations with new amount
		categories, err := s.categoryRepo.FindAllActive()
		if err != nil {
			return err
		}

		if len(categories) == 0 {
			return nil
		}

		var totalMonthlyBudget float64
		for _, cat := range categories {
			totalMonthlyBudget += cat.MonthlyBudget
		}

		if totalMonthlyBudget == 0 {
			return errors.New("total monthly budget is zero, cannot allocate")
		}

		for _, category := range categories {
			allocationAmount := (category.MonthlyBudget / totalMonthlyBudget) * newAmount

			if allocationAmount <= 0 {
				continue
			}

			budget, err := s.budgetRepo.FindByCategoryAndMonth(category.ID, newMonth, newYear)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					budget = &models.CategoryBudget{
						CategoryID:      category.ID,
						Month:           newMonth,
						Year:            newYear,
						AllocatedAmount: allocationAmount,
						SpentAmount:     0,
						RemainingAmount: allocationAmount,
					}
					if err := s.budgetRepo.Create(budget); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				budget.AllocatedAmount += allocationAmount
				budget.RemainingAmount += allocationAmount
				if err := s.budgetRepo.Update(budget); err != nil {
					return err
				}
			}

			allocation := models.BudgetAllocation{
				IncomeID:        id,
				CategoryID:      category.ID,
				AllocatedAmount: allocationAmount,
				Month:           newMonth,
				Year:            newYear,
			}

			if err := s.allocationRepo.Create(&allocation); err != nil {
				return err
			}

			allocation.Category = category
			allocations = append(allocations, allocation)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return income, allocations, nil
}

func (s *incomeService) DeleteIncome(id uuid.UUID) error {
	// TODO: Consider implementing logic to rollback allocations when deleting income
	return s.incomeRepo.Delete(id)
}
