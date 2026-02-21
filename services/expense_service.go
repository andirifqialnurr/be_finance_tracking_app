package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseService interface {
	CreateExpense(expense *models.Expense) (*models.Expense, error)
	GetExpenseByID(id uuid.UUID) (*models.Expense, error)
	GetAllExpenses(limit, offset int) ([]models.Expense, error)
	GetExpensesByCategory(categoryID uuid.UUID, limit, offset int) ([]models.Expense, error)
	GetExpensesByMonthYear(month, year int) ([]models.Expense, error)
	UpdateExpense(id uuid.UUID, expense *models.Expense) (*models.Expense, error)
	DeleteExpense(id uuid.UUID) error
}

type expenseService struct {
	expenseRepo  repositories.ExpenseRepository
	categoryRepo repositories.ExpenseCategoryRepository
	budgetRepo   repositories.CategoryBudgetRepository
	db           *gorm.DB
}

func NewExpenseService(
	expenseRepo repositories.ExpenseRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	db *gorm.DB,
) ExpenseService {
	return &expenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
		budgetRepo:   budgetRepo,
		db:           db,
	}
}

// CreateExpense creates a new expense and updates category budget
func (s *expenseService) CreateExpense(expense *models.Expense) (*models.Expense, error) {
	// Validate amount
	if expense.Amount <= 0 {
		return nil, errors.New("expense amount must be greater than 0")
	}

	// Check if category exists and is active
	category, err := s.categoryRepo.FindByID(expense.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	if !category.IsActive {
		return nil, errors.New("category is not active")
	}

	// Get month and year from expense date
	month := int(expense.Date.Month())
	year := expense.Date.Year()

	// Start transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Check budget availability
		budget, err := s.budgetRepo.FindByCategoryAndMonth(expense.CategoryID, month, year)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("no budget allocated for this category in this month")
			}
			return err
		}

		// Check if budget is sufficient
		if budget.RemainingAmount < expense.Amount {
			return errors.New("insufficient budget for this expense")
		}

		// Create expense
		if err := s.expenseRepo.Create(expense); err != nil {
			return err
		}

		// Update budget
		budget.SpentAmount += expense.Amount
		budget.RemainingAmount -= expense.Amount
		if err := s.budgetRepo.Update(budget); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load category info
	expense.Category = *category

	return expense, nil
}

func (s *expenseService) GetExpenseByID(id uuid.UUID) (*models.Expense, error) {
	return s.expenseRepo.FindByID(id)
}

func (s *expenseService) GetAllExpenses(limit, offset int) ([]models.Expense, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.expenseRepo.FindAll(limit, offset)
}

func (s *expenseService) GetExpensesByCategory(categoryID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.expenseRepo.FindByCategory(categoryID, limit, offset)
}

func (s *expenseService) GetExpensesByMonthYear(month, year int) ([]models.Expense, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}
	return s.expenseRepo.FindByMonthYear(month, year)
}

func (s *expenseService) UpdateExpense(id uuid.UUID, expense *models.Expense) (*models.Expense, error) {
	// Validate amount
	if expense.Amount <= 0 {
		return nil, errors.New("expense amount must be greater than 0")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get existing expense
		existingExpense, err := s.expenseRepo.FindByID(id)
		if err != nil {
			return err
		}

		oldCategoryID := existingExpense.CategoryID
		oldAmount := existingExpense.Amount
		oldMonth := int(existingExpense.Date.Month())
		oldYear := existingExpense.Date.Year()

		newCategoryID := expense.CategoryID
		newAmount := expense.Amount
		newMonth := int(expense.Date.Month())
		newYear := expense.Date.Year()

		// 2. Check if new category exists and is active
		category, err := s.categoryRepo.FindByID(newCategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("category not found")
			}
			return err
		}

		if !category.IsActive {
			return errors.New("category is not active")
		}

		// 3. Rollback old budget
		oldBudget, err := s.budgetRepo.FindByCategoryAndMonth(oldCategoryID, oldMonth, oldYear)
		if err == nil {
			oldBudget.SpentAmount -= oldAmount
			oldBudget.RemainingAmount += oldAmount
			if err := s.budgetRepo.Update(oldBudget); err != nil {
				return err
			}
		}

		// 4. Check new budget availability
		newBudget, err := s.budgetRepo.FindByCategoryAndMonth(newCategoryID, newMonth, newYear)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("no budget allocated for this category in this month")
			}
			return err
		}

		// Check if budget is sufficient
		if newBudget.RemainingAmount < newAmount {
			return errors.New("insufficient budget for this expense")
		}

		// 5. Update expense
		expense.ID = id
		if err := s.expenseRepo.Update(expense); err != nil {
			return err
		}

		// 6. Update new budget
		newBudget.SpentAmount += newAmount
		newBudget.RemainingAmount -= newAmount
		if err := s.budgetRepo.Update(newBudget); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load category info
	category, _ := s.categoryRepo.FindByID(expense.CategoryID)
	if category != nil {
		expense.Category = *category
	}

	return expense, nil
}

func (s *expenseService) DeleteExpense(id uuid.UUID) error {
	// TODO: Consider implementing logic to rollback budget when deleting expense
	return s.expenseRepo.Delete(id)
}
