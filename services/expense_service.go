package services

import (
"errors"
"finance-tracking-app/models"
"finance-tracking-app/repositories"

"github.com/google/uuid"
"gorm.io/gorm"
)

// ExpenseCreateResult wraps expense with optional budget warning
type ExpenseCreateResult struct {
Expense       *models.Expense `json:"expense"`
BudgetWarning string          `json:"budget_warning,omitempty"`
}

type ExpenseService interface {
CreateExpense(expense *models.Expense) (*ExpenseCreateResult, error)
GetExpenseByID(id uuid.UUID) (*models.Expense, error)
GetAllExpenses(limit, offset int) ([]models.Expense, error)
GetExpensesByCategory(categoryID uuid.UUID, limit, offset int) ([]models.Expense, error)
GetExpensesByMonthYear(month, year int) ([]models.Expense, error)
UpdateExpense(id uuid.UUID, expense *models.Expense) (*ExpenseCreateResult, error)
DeleteExpense(id uuid.UUID) error
}

type expenseService struct {
expenseRepo  repositories.ExpenseRepository
categoryRepo repositories.ExpenseCategoryRepository
budgetRepo   repositories.CategoryBudgetRepository
accountRepo  repositories.AccountRepository
db           *gorm.DB
}

func NewExpenseService(
expenseRepo repositories.ExpenseRepository,
categoryRepo repositories.ExpenseCategoryRepository,
budgetRepo repositories.CategoryBudgetRepository,
accountRepo repositories.AccountRepository,
db *gorm.DB,
) ExpenseService {
return &expenseService{
expenseRepo:  expenseRepo,
categoryRepo: categoryRepo,
budgetRepo:   budgetRepo,
accountRepo:  accountRepo,
db:           db,
}
}

func (s *expenseService) CreateExpense(expense *models.Expense) (*ExpenseCreateResult, error) {
if expense.Amount <= 0 {
return nil, errors.New("expense amount must be greater than 0")
}

// Validate account (HARD: must exist and have sufficient balance)
account, err := s.accountRepo.FindByID(expense.AccountID)
if err != nil {
return nil, errors.New("account not found")
}
if !account.IsActive {
return nil, errors.New("account is archived")
}
if account.Balance < expense.Amount {
return nil, errors.New("insufficient account balance")
}

// Validate category
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

month := int(expense.Date.Month())
year := expense.Date.Year()
result := &ExpenseCreateResult{}

err = s.db.Transaction(func(tx *gorm.DB) error {
if err := s.expenseRepo.Create(expense); err != nil {
return err
}

// HARD: deduct account balance
if err := s.accountRepo.UpdateBalance(expense.AccountID, -expense.Amount); err != nil {
return err
}

// SOFT: update budget if it exists; warn if insufficient or missing
budget, err := s.budgetRepo.FindByCategoryAndMonth(expense.CategoryID, month, year)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
result.BudgetWarning = "no budget allocated for this category this month"
return nil
}
return err
}

if budget.RemainingAmount < expense.Amount {
result.BudgetWarning = "expense exceeds remaining budget for this category"
}

budget.SpentAmount += expense.Amount
budget.RemainingAmount -= expense.Amount
return s.budgetRepo.Update(budget)
})

if err != nil {
return nil, err
}

expense.Category = *category
result.Expense = expense
return result, nil
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

func (s *expenseService) UpdateExpense(id uuid.UUID, expense *models.Expense) (*ExpenseCreateResult, error) {
if expense.Amount <= 0 {
return nil, errors.New("expense amount must be greater than 0")
}

result := &ExpenseCreateResult{}

err := s.db.Transaction(func(tx *gorm.DB) error {
existing, err := s.expenseRepo.FindByID(id)
if err != nil {
return errors.New("expense not found")
}

oldAccountID := existing.AccountID
oldAmount := existing.Amount
oldCategoryID := existing.CategoryID
oldMonth := int(existing.Date.Month())
oldYear := existing.Date.Year()
newMonth := int(expense.Date.Month())
newYear := expense.Date.Year()

// Validate new account (HARD)
newAccount, err := s.accountRepo.FindByID(expense.AccountID)
if err != nil {
return errors.New("account not found")
}
if !newAccount.IsActive {
return errors.New("account is archived")
}

// Validate new category
category, err := s.categoryRepo.FindByID(expense.CategoryID)
if err != nil {
return errors.New("category not found")
}
if !category.IsActive {
return errors.New("category is not active")
}

// Restore old account balance
if err := s.accountRepo.UpdateBalance(oldAccountID, oldAmount); err != nil {
return err
}

// HARD: check new account balance for new amount (after restore)
refreshedAccount, _ := s.accountRepo.FindByID(expense.AccountID)
if refreshedAccount != nil && refreshedAccount.Balance < expense.Amount {
return errors.New("insufficient account balance")
}

// Deduct new account
if err := s.accountRepo.UpdateBalance(expense.AccountID, -expense.Amount); err != nil {
return err
}

// Rollback old budget (SOFT)
oldBudget, err := s.budgetRepo.FindByCategoryAndMonth(oldCategoryID, oldMonth, oldYear)
if err == nil {
oldBudget.SpentAmount -= oldAmount
oldBudget.RemainingAmount += oldAmount
_ = s.budgetRepo.Update(oldBudget)
}

// Update expense
expense.ID = id
if err := s.expenseRepo.Update(expense); err != nil {
return err
}

// Update new budget (SOFT)
newBudget, err := s.budgetRepo.FindByCategoryAndMonth(expense.CategoryID, newMonth, newYear)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
result.BudgetWarning = "no budget allocated for this category this month"
return nil
}
return err
}

if newBudget.RemainingAmount < expense.Amount {
result.BudgetWarning = "expense exceeds remaining budget for this category"
}
newBudget.SpentAmount += expense.Amount
newBudget.RemainingAmount -= expense.Amount
return s.budgetRepo.Update(newBudget)
})

if err != nil {
return nil, err
}

category, _ := s.categoryRepo.FindByID(expense.CategoryID)
if category != nil {
expense.Category = *category
}
result.Expense = expense
return result, nil
}

func (s *expenseService) DeleteExpense(id uuid.UUID) error {
return s.db.Transaction(func(tx *gorm.DB) error {
expense, err := s.expenseRepo.FindByID(id)
if err != nil {
return errors.New("expense not found")
}

month := int(expense.Date.Month())
year := expense.Date.Year()

// Restore account balance
if err := s.accountRepo.UpdateBalance(expense.AccountID, expense.Amount); err != nil {
return err
}

// Restore budget (SOFT)
budget, err := s.budgetRepo.FindByCategoryAndMonth(expense.CategoryID, month, year)
if err == nil {
budget.SpentAmount -= expense.Amount
budget.RemainingAmount += expense.Amount
_ = s.budgetRepo.Update(budget)
}

return s.expenseRepo.Delete(id)
})
}
