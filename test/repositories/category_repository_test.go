package repositories

import (
	"encoding/json"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"finance-tracking-app/test/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestCategoryRepository_Create(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewExpenseCategoryRepository(db)

	// Test data
	metadata, _ := json.Marshal(map[string]interface{}{"test": "data"})
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
		Metadata:           datatypes.JSON(metadata),
	}

	// Execute
	err := repo.Create(category)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, category.ID)
}

func TestCategoryRepository_FindByID(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewExpenseCategoryRepository(db)

	// Create test category
	metadata, _ := json.Marshal(map[string]interface{}{"test": "data"})
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
		Metadata:           datatypes.JSON(metadata),
	}
	err := repo.Create(category)
	assert.NoError(t, err)

	// Execute
	found, err := repo.FindByID(category.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, category.ID, found.ID)
	assert.Equal(t, category.Name, found.Name)
	assert.Equal(t, category.MonthlyBudget, found.MonthlyBudget)
}

func TestCategoryRepository_FindAllActive(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewExpenseCategoryRepository(db)

	// Create active categories
	activeCategories := []*models.ExpenseCategory{
		{Name: "Active 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true},
		{Name: "Active 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true},
	}

	for _, cat := range activeCategories {
		err := repo.Create(cat)
		assert.NoError(t, err)
	}

	// Create inactive category (need to update after creation due to GORM default:true)
	inactiveCategory := &models.ExpenseCategory{
		Name:               "Inactive",
		Type:               models.CategoryTypeOneTime,
		MonthlyBudget:      200000,
		AllocationPriority: 3,
		IsActive:           true, // Create as true first
	}
	err := repo.Create(inactiveCategory)
	assert.NoError(t, err)

	// Now update to inactive
	inactiveCategory.IsActive = false
	err = repo.Update(inactiveCategory)
	assert.NoError(t, err)

	// Execute
	found, err := repo.FindAllActive()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	for _, cat := range found {
		assert.True(t, cat.IsActive)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewExpenseCategoryRepository(db)

	// Create test category
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	err := repo.Create(category)
	assert.NoError(t, err)

	// Update
	category.Name = "Updated Name"
	category.MonthlyBudget = 200000

	// Execute
	err = repo.Update(category)

	// Assert
	assert.NoError(t, err)

	// Verify update
	found, err := repo.FindByID(category.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", found.Name)
	assert.Equal(t, float64(200000), found.MonthlyBudget)
}

func TestCategoryRepository_Delete(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewExpenseCategoryRepository(db)

	// Create test category
	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	err := repo.Create(category)
	assert.NoError(t, err)

	// Execute
	err = repo.Delete(category.ID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	found, err := repo.FindByID(category.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
