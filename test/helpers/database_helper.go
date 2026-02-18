package helpers

import (
	"finance-tracking-app/models"
	"fmt"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates a PostgreSQL test database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	// PostgreSQL connection for testing
	dsn := "host=localhost user=postgres password=Aran2706$ dbname=finance_tracking_test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&models.Income{},
		&models.ExpenseCategory{},
		&models.CategoryBudget{},
		&models.Expense{},
		&models.BudgetAllocation{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Clean all data before test
	TruncateAllTables(t, db)

	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	// Truncate all tables after test
	TruncateAllTables(t, db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("Failed to get database instance: %v", err)
		return
	}
	sqlDB.Close()
}

// TruncateAllTables removes all data from all tables
func TruncateAllTables(t *testing.T, db *gorm.DB) {
	tables := []string{
		"budget_allocations",
		"expenses",
		"category_budgets",
		"incomes",
		"expense_categories",
	}

	// Disable foreign key checks temporarily
	db.Exec("SET session_replication_role = 'replica'")

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}

	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = 'origin'")
}

// SeedTestData seeds the database with test data
func SeedTestData(t *testing.T, db *gorm.DB) {
	categories := []models.ExpenseCategory{
		{
			Name:               "Test Category 1",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      100000,
			AllocationPriority: 1,
			IsActive:           true,
		},
		{
			Name:               "Test Category 2",
			Type:               models.CategoryTypeUsageBased,
			MonthlyBudget:      50000,
			AllocationPriority: 2,
			IsActive:           true,
		},
	}

	for _, cat := range categories {
		if err := db.Create(&cat).Error; err != nil {
			t.Fatalf("Failed to seed category: %v", err)
		}
	}
}
