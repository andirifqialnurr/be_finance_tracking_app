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

func TestExpenseService_CreateExpense_ValidBudget(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)

	// Create a category
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	// Create budget for this category
	budget := &models.CategoryBudget{
		CategoryID:      category.ID,
		Month:           2,
		Year:            2026,
		AllocatedAmount: 100000,
		SpentAmount:     0,
		RemainingAmount: 100000,
	}
	budgetRepo.Create(budget)

	// Create expense
	expense := &models.Expense{
		CategoryID:  category.ID,
		Amount:      50000,
		Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		Description: "Test expense",
	}

	// Execute
	createdExpense, err := expenseService.CreateExpense(expense)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdExpense)
	assert.NotEqual(t, uuid.Nil, createdExpense.ID)

	// Verify budget was deducted
	updatedBudget, err := budgetRepo.FindByCategoryAndMonth(category.ID, 2, 2026)
	assert.NoError(t, err)
	assert.Equal(t, 50000.0, updatedBudget.SpentAmount)
	assert.Equal(t, 50000.0, updatedBudget.RemainingAmount)
}

func TestExpenseService_CreateExpense_ExceedsBudget(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)

	// Create a category
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	// Create budget with limited amount
	budget := &models.CategoryBudget{
		CategoryID:      category.ID,
		Month:           2,
		Year:            2026,
		AllocatedAmount: 50000,
		SpentAmount:     0,
		RemainingAmount: 50000,
	}
	budgetRepo.Create(budget)

	// Try to create expense that exceeds budget
	expense := &models.Expense{
		CategoryID:  category.ID,
		Amount:      75000, // Exceeds remaining budget of 50000
		Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		Description: "Expensive item",
	}

	// Execute
	createdExpense, err := expenseService.CreateExpense(expense)

	// Assert - should return error
	assert.Error(t, err)
	assert.Nil(t, createdExpense)
	assert.Contains(t, err.Error(), "insufficient budget")
}

func TestExpenseService_CreateExpense_NoBudget(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)

	// Create a category
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	// No budget created - expense should fail

	expense := &models.Expense{
		CategoryID:  category.ID,
		Amount:      50000,
		Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		Description: "Test expense",
	}

	// Execute
	createdExpense, err := expenseService.CreateExpense(expense)

	// Assert - should return error
	assert.Error(t, err)
	assert.Nil(t, createdExpense)
}

func TestExpenseService_GetExpensesByCategory(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)

	// Create categories
	category1 := &models.ExpenseCategory{Name: "Category 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true}
	category2 := &models.ExpenseCategory{Name: "Category 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true}

	categoryRepo.Create(category1)
	categoryRepo.Create(category2)

	// Create budgets
	budget1 := &models.CategoryBudget{CategoryID: category1.ID, Month: 2, Year: 2026, AllocatedAmount: 100000, RemainingAmount: 100000}
	budget2 := &models.CategoryBudget{CategoryID: category2.ID, Month: 2, Year: 2026, AllocatedAmount: 50000, RemainingAmount: 50000}

	budgetRepo.Create(budget1)
	budgetRepo.Create(budget2)

	// Create expenses
	expenses := []models.Expense{
		{CategoryID: category1.ID, Amount: 10000, Date: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC), Description: "Expense 1"},
		{CategoryID: category1.ID, Amount: 20000, Date: time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), Description: "Expense 2"},
		{CategoryID: category2.ID, Amount: 15000, Date: time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC), Description: "Expense 3"},
	}

	for _, exp := range expenses {
		expenseService.CreateExpense(&exp)
	}

	// Execute - get expenses for category1
	found, err := expenseService.GetExpensesByCategory(category1.ID, 10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, found, 2)

	for _, exp := range found {
		assert.Equal(t, category1.ID, exp.CategoryID)
	}
}

func TestExpenseService_DeleteExpense(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)

	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)

	// Create category and budget
	category := &models.ExpenseCategory{Name: "Test", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true}
	categoryRepo.Create(category)

	budget := &models.CategoryBudget{CategoryID: category.ID, Month: 2, Year: 2026, AllocatedAmount: 100000, RemainingAmount: 100000}
	budgetRepo.Create(budget)

	// Create expense
	expense := &models.Expense{
		CategoryID:  category.ID,
		Amount:      30000,
		Date:        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC),
		Description: "To be deleted",
	}
	createdExpense, _ := expenseService.CreateExpense(expense)

	// Execute
	err := expenseService.DeleteExpense(createdExpense.ID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	found, err := expenseRepo.FindByID(createdExpense.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
