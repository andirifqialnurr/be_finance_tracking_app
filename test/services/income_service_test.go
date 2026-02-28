package services

import (
"finance-tracking-app/models"
"finance-tracking-app/repositories"
"finance-tracking-app/services"
"finance-tracking-app/test/helpers"
"testing"
"time"

"github.com/stretchr/testify/assert"
)

// createSalaryCardAccount creates a CARD account with SALARY income type and the given balance.
func createSalaryCardAccount(t *testing.T, accountRepo repositories.AccountRepository, balance float64) *models.Account {
t.Helper()
incomeType := models.CardIncomeTypeSalary
acc := &models.Account{
Name:       "Salary Card",
Type:       models.AccountTypeCard,
IncomeType: &incomeType,
Balance:    balance,
IsActive:   true,
}
if err := accountRepo.Create(acc); err != nil {
t.Fatalf("failed to create salary card account: %v", err)
}
return acc
}

func TestIncomeService_CreateIncome_WithAutoAllocation(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

incomeRepo := repositories.NewIncomeRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
allocationRepo := repositories.NewBudgetAllocationRepository(db)
accountRepo := repositories.NewAccountRepository(db)

incomeService := services.NewIncomeService(incomeRepo, accountRepo, categoryRepo, budgetRepo, allocationRepo, db)

// SALARY CARD account — auto-allocation will be triggered
account := createSalaryCardAccount(t, accountRepo, 0)

categories := []models.ExpenseCategory{
{Name: "Category 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true},
{Name: "Category 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true},
}
for i := range categories {
categoryRepo.Create(&categories[i])
}

income := &models.Income{
AccountID:   account.ID,
Source:      "Test Income",
Amount:      300000,
Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
Description: "Test description",
}

result, err := incomeService.CreateIncome(income)

assert.NoError(t, err)
assert.NotNil(t, result)
assert.Len(t, result.Allocations, 2)
assert.False(t, result.AlreadyAllocatedWarning)

totalAllocated := 0.0
for _, allocation := range result.Allocations {
totalAllocated += allocation.AllocatedAmount
}
assert.InDelta(t, income.Amount, totalAllocated, 0.01)

budgets, err := budgetRepo.FindByMonthYear(2, 2026)
assert.NoError(t, err)
assert.Len(t, budgets, 2)

for _, budget := range budgets {
assert.Greater(t, budget.AllocatedAmount, 0.0)
assert.Equal(t, budget.AllocatedAmount, budget.RemainingAmount)
assert.Equal(t, 0.0, budget.SpentAmount)
}
}

func TestIncomeService_CreateIncome_NoActiveCategories(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

incomeRepo := repositories.NewIncomeRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
allocationRepo := repositories.NewBudgetAllocationRepository(db)
accountRepo := repositories.NewAccountRepository(db)

incomeService := services.NewIncomeService(incomeRepo, accountRepo, categoryRepo, budgetRepo, allocationRepo, db)

// SALARY CARD — auto-allocation attempted but no categories exist
account := createSalaryCardAccount(t, accountRepo, 0)

income := &models.Income{
AccountID:   account.ID,
Source:      "Test Income",
Amount:      300000,
Date:        time.Now(),
Description: "Test description",
}

result, err := incomeService.CreateIncome(income)

assert.NoError(t, err)
assert.NotNil(t, result)
assert.Empty(t, result.Allocations)
assert.False(t, result.AlreadyAllocatedWarning)
}

func TestIncomeService_GetIncomesByMonthYear(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

incomeRepo := repositories.NewIncomeRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
allocationRepo := repositories.NewBudgetAllocationRepository(db)
accountRepo := repositories.NewAccountRepository(db)

incomeService := services.NewIncomeService(incomeRepo, accountRepo, categoryRepo, budgetRepo, allocationRepo, db)

incomes := []models.Income{
{Source: "Feb Income 1", Amount: 1000000, Date: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)},
{Source: "Feb Income 2", Amount: 2000000, Date: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)},
{Source: "Mar Income", Amount: 3000000, Date: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
}
for i := range incomes {
incomeRepo.Create(&incomes[i])
}

found, err := incomeService.GetIncomesByMonthYear(2, 2026)

assert.NoError(t, err)
assert.Len(t, found, 2)
for _, income := range found {
assert.Equal(t, 2, int(income.Date.Month()))
assert.Equal(t, 2026, income.Date.Year())
}
}

func TestIncomeService_DeleteIncome(t *testing.T) {
db := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, db)

incomeRepo := repositories.NewIncomeRepository(db)
categoryRepo := repositories.NewExpenseCategoryRepository(db)
budgetRepo := repositories.NewCategoryBudgetRepository(db)
allocationRepo := repositories.NewBudgetAllocationRepository(db)
accountRepo := repositories.NewAccountRepository(db)

incomeService := services.NewIncomeService(incomeRepo, accountRepo, categoryRepo, budgetRepo, allocationRepo, db)

// Need a real account so DeleteIncome can reverse the balance
account := createSalaryCardAccount(t, accountRepo, 1000000)

income := &models.Income{
AccountID: account.ID,
Source:    "Test Income",
Amount:    1000000,
Date:      time.Now(),
}
incomeRepo.Create(income)

err := incomeService.DeleteIncome(income.ID)
assert.NoError(t, err)

found, err := incomeRepo.FindByID(income.ID)
assert.Error(t, err)
assert.Nil(t, found)
}
