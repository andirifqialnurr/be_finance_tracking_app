package repositories

import (
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"finance-tracking-app/test/helpers"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIncomeRepository_Create(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Test data
	income := &models.Income{
		Source:      "Test Income",
		Amount:      1000000,
		Date:        time.Now(),
		Description: "Test description",
	}

	// Execute
	err := repo.Create(income)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, income.ID)
	assert.NotZero(t, income.CreatedAt)
}

func TestIncomeRepository_FindByID(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Create test income
	income := &models.Income{
		Source:      "Test Income",
		Amount:      1000000,
		Date:        time.Now(),
		Description: "Test description",
	}
	err := repo.Create(income)
	assert.NoError(t, err)

	// Execute
	found, err := repo.FindByID(income.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, income.ID, found.ID)
	assert.Equal(t, income.Source, found.Source)
	assert.Equal(t, income.Amount, found.Amount)
}

func TestIncomeRepository_FindByID_NotFound(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Execute
	found, err := repo.FindByID(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestIncomeRepository_FindAll(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Create test incomes
	incomes := []*models.Income{
		{Source: "Income 1", Amount: 1000000, Date: time.Now()},
		{Source: "Income 2", Amount: 2000000, Date: time.Now()},
		{Source: "Income 3", Amount: 3000000, Date: time.Now()},
	}

	for _, income := range incomes {
		err := repo.Create(income)
		assert.NoError(t, err)
	}

	// Execute
	found, err := repo.FindAll(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, found, 3)
}

func TestIncomeRepository_FindByMonthYear(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Create test incomes with different dates
	targetDate := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	otherDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)

	incomes := []*models.Income{
		{Source: "Feb Income 1", Amount: 1000000, Date: targetDate},
		{Source: "Feb Income 2", Amount: 2000000, Date: targetDate},
		{Source: "Mar Income", Amount: 3000000, Date: otherDate},
	}

	for _, income := range incomes {
		err := repo.Create(income)
		assert.NoError(t, err)
	}

	// Execute
	found, err := repo.FindByMonthYear(2, 2026)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	for _, income := range found {
		assert.Equal(t, 2, int(income.Date.Month()))
		assert.Equal(t, 2026, income.Date.Year())
	}
}

func TestIncomeRepository_Delete(t *testing.T) {
	// Setup
	db := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, db)
	repo := repositories.NewIncomeRepository(db)

	// Create test income
	income := &models.Income{
		Source: "Test Income",
		Amount: 1000000,
		Date:   time.Now(),
	}
	err := repo.Create(income)
	assert.NoError(t, err)

	// Execute
	err = repo.Delete(income.ID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	found, err := repo.FindByID(income.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
