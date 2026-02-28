package database

import (
	"finance-tracking-app/config"
	"finance-tracking-app/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection and runs migrations
func InitDB() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	runMigrations()
}

// addNotNullColumn safely adds a NOT NULL uuid column to an existing table by:
// 1. Adding the column as nullable first
// 2. Deleting rows that have no value (they are invalid without a FK anyway)
// 3. Setting the column to NOT NULL so AutoMigrate skips the ALTER
func addNotNullColumn(table, column string) {
	DB.Exec(`ALTER TABLE "` + table + `" ADD COLUMN IF NOT EXISTS "` + column + `" uuid`)
	DB.Exec(`DELETE FROM "` + table + `" WHERE "` + column + `" IS NULL`)
	DB.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "` + column + `" SET NOT NULL`)
}

// preMigrate handles column additions that would fail due to NOT NULL on existing rows
func preMigrate() {
	type tableColumn struct {
		table  string
		column string
		model  interface{}
	}

	checks := []tableColumn{
		{"incomes", "account_id", &models.Income{}},
		{"expenses", "account_id", &models.Expense{}},
	}

	for _, c := range checks {
		if DB.Migrator().HasTable(c.table) && !DB.Migrator().HasColumn(c.model, c.column) {
			addNotNullColumn(c.table, c.column)
		}
	}
}

// runMigrations automatically migrates database schema
func runMigrations() {
	preMigrate()

	err := DB.AutoMigrate(
		&models.Account{},
		&models.Income{},
		&models.ExpenseCategory{},
		&models.CategoryBudget{},
		&models.Expense{},
		&models.BudgetAllocation{},
		&models.BudgetReallocation{},
		&models.BudgetAlert{},
		&models.AccountTransfer{},
		&models.ScheduledFund{},
		&models.NotificationSetting{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed successfully")

	// Add unique constraint for CategoryBudget
	DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_category_budget_unique 
		ON category_budgets(category_id, month, year)
	`)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Println("Error closing database:", err)
		return
	}

	log.Println("Database connection closed")
}
