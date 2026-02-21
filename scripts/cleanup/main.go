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

	// Delete in correct order to avoid foreign key constraints
	log.Println("\n1. Deleting Budget Alerts...")
	if err := db.Exec("DELETE FROM budget_alerts").Error; err != nil {
		log.Printf("Error deleting budget_alerts: %v", err)
	} else {
		log.Println("✅ Budget alerts deleted")
	}

	log.Println("\n2. Deleting Budget Reallocations...")
	if err := db.Exec("DELETE FROM budget_reallocations").Error; err != nil {
		log.Printf("Error deleting budget_reallocations: %v", err)
	} else {
		log.Println("✅ Budget reallocations deleted")
	}

	log.Println("\n3. Deleting Budget Allocations...")
	if err := db.Exec("DELETE FROM budget_allocations").Error; err != nil {
		log.Printf("Error deleting budget_allocations: %v", err)
	} else {
		log.Println("✅ Budget allocations deleted")
	}

	log.Println("\n4. Deleting Expenses...")
	if err := db.Exec("DELETE FROM expenses").Error; err != nil {
		log.Printf("Error deleting expenses: %v", err)
	} else {
		log.Println("✅ Expenses deleted")
	}

	log.Println("\n5. Deleting Category Budgets...")
	if err := db.Exec("DELETE FROM category_budgets").Error; err != nil {
		log.Printf("Error deleting category_budgets: %v", err)
	} else {
		log.Println("✅ Category budgets deleted")
	}

	log.Println("\n6. Deleting Incomes...")
	if err := db.Exec("DELETE FROM incomes").Error; err != nil {
		log.Printf("Error deleting incomes: %v", err)
	} else {
		log.Println("✅ Incomes deleted")
	}

	log.Println("\n7. Deleting Expense Categories...")
	if err := db.Exec("DELETE FROM expense_categories").Error; err != nil {
		log.Printf("Error deleting expense_categories: %v", err)
	} else {
		log.Println("✅ Expense categories deleted")
	}

	// Verify cleanup
	var counts struct {
		Categories    int64
		Incomes       int64
		Expenses      int64
		Budgets       int64
		Allocations   int64
		Reallocations int64
		Alerts        int64
	}

	db.Model(&models.ExpenseCategory{}).Count(&counts.Categories)
	db.Model(&models.Income{}).Count(&counts.Incomes)
	db.Model(&models.Expense{}).Count(&counts.Expenses)
	db.Model(&models.CategoryBudget{}).Count(&counts.Budgets)
	db.Model(&models.BudgetAllocation{}).Count(&counts.Allocations)
	db.Model(&models.BudgetReallocation{}).Count(&counts.Reallocations)
	db.Model(&models.BudgetAlert{}).Count(&counts.Alerts)

	log.Println("\n" + "============================================================")
	log.Println("🎉 DATABASE CLEANUP COMPLETED!")
	log.Println("============================================================")
	log.Printf("\n📊 Remaining Records:")
	log.Printf("  Categories: %d", counts.Categories)
	log.Printf("  Incomes: %d", counts.Incomes)
	log.Printf("  Expenses: %d", counts.Expenses)
	log.Printf("  Budgets: %d", counts.Budgets)
	log.Printf("  Allocations: %d", counts.Allocations)
	log.Printf("  Reallocations: %d", counts.Reallocations)
	log.Printf("  Alerts: %d", counts.Alerts)

	if counts.Categories == 0 && counts.Incomes == 0 && counts.Expenses == 0 {
		log.Println("\n✨ Database is now clean and ready for fresh seeding!")
		log.Println("💡 Run: go run scripts/seed/main.go")
	} else {
		log.Println("\n⚠️  Warning: Some records still exist in database")
	}
}
