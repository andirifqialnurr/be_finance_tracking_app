package main

import (
	"finance-tracking-app/config"
	"finance-tracking-app/database"
	"finance-tracking-app/models"
	"log"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Get database instance
	db := database.GetDB()

	log.Println("🧹 Starting database cleanup...")
	log.Println("⚠️  WARNING: This will delete ALL data from the database!")

	// Delete in correct order to respect foreign key constraints
	tables := []struct {
		label string
		table string
	}{
		{"Notification Settings", "notification_settings"},
		{"Scheduled Funds", "scheduled_funds"},
		{"Budget Alerts", "budget_alerts"},
		{"Budget Reallocations", "budget_reallocations"},
		{"Budget Allocations", "budget_allocations"},
		{"Expenses", "expenses"},
		{"Category Budgets", "category_budgets"},
		{"Incomes", "incomes"},
		{"Account Transfers", "account_transfers"},
		{"Accounts", "accounts"},
		{"Expense Categories", "expense_categories"},
	}

	for i, t := range tables {
		log.Printf("\n%d. Deleting %s...", i+1, t.label)
		if err := db.Exec("DELETE FROM " + t.table).Error; err != nil {
			log.Printf("   Error: %v", err)
		} else {
			log.Printf("   OK: %s deleted", t.label)
		}
	}

	// Verify cleanup
	var catCount, incCount, expCount, accCount int64
	db.Model(&models.ExpenseCategory{}).Count(&catCount)
	db.Model(&models.Income{}).Count(&incCount)
	db.Model(&models.Expense{}).Count(&expCount)
	db.Model(&models.Account{}).Count(&accCount)

	log.Println("\n" + "============================================================")
	log.Println("DATABASE CLEANUP COMPLETED!")
	log.Println("============================================================")
	log.Printf("  Categories : %d", catCount)
	log.Printf("  Incomes    : %d", incCount)
	log.Printf("  Expenses   : %d", expCount)
	log.Printf("  Accounts   : %d", accCount)

	if catCount == 0 && incCount == 0 && expCount == 0 && accCount == 0 {
		log.Println("\nDatabase is now clean. Run: go run scripts/seed/main.go")
	} else {
		log.Println("\nWarning: Some records still exist in database")
	}
}
