package services

import (
"finance-tracking-app/models"
"finance-tracking-app/repositories"
"finance-tracking-app/services"
"finance-tracking-app/test/helpers"
"testing"
"time"

"github.com/google/uuid"
"github.com/stretchr/testify/assert"
)

// createTestAccount inserts an account with the given balance into the DB and returns it.
func createTestAccount(t *testing.T, accountRepo repositories.AccountRepository, balance float64) *models.Account {
t.Helper()
acc := &models.Account{
Name:     "Test Account",
Type:     models.AccountTypeCash,
Balance:  balance,
IsActive: true,
}
if err := accountRepo.Create(acc); err != nil {
t.Fatalf("failed to create test account: %v", err)
}
return acc
}

func TestExpenseService_CreateExpense_ValidBudget(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

account := createTestAccount(t, accountRepo, 200000)

category := &models.ExpenseCategory{
Name:               "Test Category",
Type:               models.CategoryTypeSubscription,
MonthlyBudget:      100000,
AllocationPriority: 1,
IsActive:           true,
}
categoryRepo.Create(category)

budget := &models.CategoryBudget{
CategoryID:      category.ID,
Month:           2,
Year:            2026,
AllocatedAmount: 100000,
SpentAmount:     0,
RemainingAmount: 100000,
}
budgetRepo.Create(budget)

expense := &models.Expense{
AccountID:   account.ID,
CategoryID:  category.ID,
Amount:      50000,
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "Test expense",
}

result, err := expenseService.CreateExpense(expense)

assert.NoError(t, err)
assert.NotNil(t, result)
assert.NotEqual(t, uuid.Nil, result.Expense.ID)
assert.Empty(t, result.BudgetWarning)

updatedBudget, err := budgetRepo.FindByCategoryAndMonth(category.ID, 2, 2026)
assert.NoError(t, err)
assert.Equal(t, 50000.0, updatedBudget.SpentAmount)
assert.Equal(t, 50000.0, updatedBudget.RemainingAmount)
}

func TestExpenseService_CreateExpense_ExceedsBudget(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

// Account has enough balance — budget check is SOFT (warning only)
account := createTestAccount(t, accountRepo, 200000)

category := &models.ExpenseCategory{
Name:               "Test Category",
Type:               models.CategoryTypeSubscription,
MonthlyBudget:      100000,
AllocationPriority: 1,
IsActive:           true,
}
categoryRepo.Create(category)

budget := &models.CategoryBudget{
CategoryID:      category.ID,
Month:           2,
Year:            2026,
AllocatedAmount: 50000,
SpentAmount:     0,
RemainingAmount: 50000,
}
budgetRepo.Create(budget)

expense := &models.Expense{
AccountID:   account.ID,
CategoryID:  category.ID,
Amount:      75000, // Exceeds budget remaining but account balance is fine
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "Expensive item",
}

// Budget exceeded is now a SOFT warning — expense succeeds with a warning
result, err := expenseService.CreateExpense(expense)

assert.NoError(t, err)
assert.NotNil(t, result)
assert.NotEmpty(t, result.BudgetWarning, "should carry a budget warning")
assert.Contains(t, result.BudgetWarning, "exceeds remaining budget")
}

func TestExpenseService_CreateExpense_InsufficientAccountBalance(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

// Account balance is too low — HARD check must block the expense
account := createTestAccount(t, accountRepo, 10000)

category := &models.ExpenseCategory{
Name:               "Test Category",
Type:               models.CategoryTypeSubscription,
MonthlyBudget:      100000,
AllocationPriority: 1,
IsActive:           true,
}
categoryRepo.Create(category)

budget := &models.CategoryBudget{
CategoryID:      category.ID,
Month:           2,
Year:            2026,
AllocatedAmount: 100000,
SpentAmount:     0,
RemainingAmount: 100000,
}
budgetRepo.Create(budget)

expense := &models.Expense{
AccountID:   account.ID,
CategoryID:  category.ID,
Amount:      50000,
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "Too expensive",
}

result, err := expenseService.CreateExpense(expense)

assert.Error(t, err)
assert.Nil(t, result)
assert.Contains(t, err.Error(), "insufficient account balance")
}

func TestExpenseService_CreateExpense_NoBudget(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

account := createTestAccount(t, accountRepo, 200000)

category := &models.ExpenseCategory{
Name:               "Test Category",
Type:               models.CategoryTypeSubscription,
MonthlyBudget:      100000,
AllocationPriority: 1,
IsActive:           true,
}
categoryRepo.Create(category)

// No budget created — expense succeeds with a "no budget" warning (SOFT)
expense := &models.Expense{
AccountID:   account.ID,
CategoryID:  category.ID,
Amount:      50000,
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "Test expense",
}

result, err := expenseService.CreateExpense(expense)

assert.NoError(t, err)
assert.NotNil(t, result)
assert.NotEmpty(t, result.BudgetWarning, "should carry a no-budget warning")
assert.Contains(t, result.BudgetWarning, "no budget allocated")
}

func TestExpenseService_GetExpensesByCategory(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

account := createTestAccount(t, accountRepo, 500000)

category1 := &models.ExpenseCategory{Name: "Category 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true}
category2 := &models.ExpenseCategory{Name: "Category 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true}
categoryRepo.Create(category1)
categoryRepo.Create(category2)

budget1 := &models.CategoryBudget{CategoryID: category1.ID, Month: 2, Year: 2026, AllocatedAmount: 100000, RemainingAmount: 100000}
budget2 := &models.CategoryBudget{CategoryID: category2.ID, Month: 2, Year: 2026, AllocatedAmount: 50000, RemainingAmount: 50000}
budgetRepo.Create(budget1)
budgetRepo.Create(budget2)

expenses := []models.Expense{
{AccountID: account.ID, CategoryID: category1.ID, Amount: 10000, Date: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC), Description: "Expense 1"},
{AccountID: account.ID, CategoryID: category1.ID, Amount: 20000, Date: time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), Description: "Expense 2"},
{AccountID: account.ID, CategoryID: category2.ID, Amount: 15000, Date: time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC), Description: "Expense 3"},
}

for i := range expenses {
expenseService.CreateExpense(&expenses[i])
}

found, err := expenseService.GetExpensesByCategory(category1.ID, 10, 0)

assert.NoError(t, err)
assert.Len(t, found, 2)
for _, exp := range found {
assert.Equal(t, category1.ID, exp.CategoryID)
}
}

func TestExpenseService_DeleteExpense(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

expenseRepo := repositories.NewExpenseRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
accountRepo := repositories.NewAccountRepository(db)

expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)

account := createTestAccount(t, accountRepo, 200000)

category := &models.ExpenseCategory{Name: "Test", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true}
categoryRepo.Create(category)

budget := &models.CategoryBudget{CategoryID: category.ID, Month: 2, Year: 2026, AllocatedAmount: 100000, RemainingAmount: 100000}
budgetRepo.Create(budget)

expense := &models.Expense{
AccountID:   account.ID,
CategoryID:  category.ID,
Amount:      30000,
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "To be deleted",
}
createdResult, err := expenseService.CreateExpense(expense)
assert.NoError(t, err)

err = expenseService.DeleteExpense(createdResult.Expense.ID)
assert.NoError(t, err)

found, err := expenseRepo.FindByID(createdResult.Expense.ID)
assert.Error(t, err)
assert.Nil(t, found)
}
