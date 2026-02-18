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

func TestIncomeService_CreateIncome_WithAutoAllocation(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	// Setup repositories
	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	// Create service
	incomeService := services.NewIncomeService(incomeRepo, categoryRepo, budgetRepo, allocationRepo, db)

	// Seed categories
	categories := []models.ExpenseCategory{
		{Name: "Category 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true},
		{Name: "Category 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true},
	}

	for _, cat := range categories {
		categoryRepo.Create(&cat)
	}

	// Test income creation
	income := &models.Income{
		Source:      "Test Income",
		Amount:      300000,
		Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		Description: "Test description",
	}

	// Execute
	createdIncome, allocations, err := incomeService.CreateIncome(income)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdIncome)
	assert.Len(t, allocations, 2) // Should allocate to 2 categories

	// Verify allocations proportions
	totalAllocated := 0.0
	for _, allocation := range allocations {
		totalAllocated += allocation.AllocatedAmount
	}

	// Total allocated should be close to income amount (accounting for float precision)
	assert.InDelta(t, income.Amount, totalAllocated, 0.01)

	// Verify category budgets were created
	month := 2
	year := 2026
	budgets, err := budgetRepo.FindByMonthYear(month, year)
	assert.NoError(t, err)
	assert.Len(t, budgets, 2)

	// Verify proportional allocation
	for _, budget := range budgets {
		assert.Greater(t, budget.AllocatedAmount, 0.0)
		assert.Equal(t, budget.AllocatedAmount, budget.RemainingAmount)
		assert.Equal(t, 0.0, budget.SpentAmount)
	}
}

func TestIncomeService_CreateIncome_NoActiveCategories(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	incomeService := services.NewIncomeService(incomeRepo, categoryRepo, budgetRepo, allocationRepo, db)

	// Test income creation without active categories
	income := &models.Income{
		Source:      "Test Income",
		Amount:      300000,
		Date:        time.Now(),
		Description: "Test description",
	}

	// Execute
	createdIncome, allocations, err := incomeService.CreateIncome(income)

	// Assert - should create income but no allocations
	assert.NoError(t, err)
	assert.NotNil(t, createdIncome)
	assert.Len(t, allocations, 0)
}

func TestIncomeService_GetIncomesByMonthYear(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	incomeService := services.NewIncomeService(incomeRepo, categoryRepo, budgetRepo, allocationRepo, db)

	// Create incomes for different months
	incomes := []models.Income{
		{Source: "Feb Income 1", Amount: 1000000, Date: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)},
		{Source: "Feb Income 2", Amount: 2000000, Date: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC)},
		{Source: "Mar Income", Amount: 3000000, Date: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
	}

	for _, income := range incomes {
		incomeRepo.Create(&income)
	}

	// Execute
	found, err := incomeService.GetIncomesByMonthYear(2, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, found, 2)

	for _, income := range found {
		assert.Equal(t, 2, int(income.Date.Month()))
		assert.Equal(t, 2026, income.Date.Year())
	}
}

func TestIncomeService_DeleteIncome(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	incomeService := services.NewIncomeService(incomeRepo, categoryRepo, budgetRepo, allocationRepo, db)

	// Create income
	income := &models.Income{
		Source: "Test Income",
		Amount: 1000000,
		Date:   time.Now(),
	}
	incomeRepo.Create(income)

	// Execute
	err := incomeService.DeleteIncome(income.ID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	found, err := incomeRepo.FindByID(income.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
