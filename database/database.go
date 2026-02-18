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

// runMigrations automatically migrates database schema
func runMigrations() {
	err := DB.AutoMigrate(
		&models.Income{},
		&models.ExpenseCategory{},
		&models.CategoryBudget{},
		&models.Expense{},
		&models.BudgetAllocation{},
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
