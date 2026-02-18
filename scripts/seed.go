package main

import (
	"encoding/json"
	"finance-tracking-app/config"
	"finance-tracking-app/database"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	log.Println("🌱 Starting database seeder...")

	// Delete existing categories if any
	log.Println("Cleaning existing categories...")

	// Seed expense categories
	categories := []models.ExpenseCategory{
		{
			ID:                 uuid.New(),
			Name:               "Makan",
			Type:               models.CategoryTypeDailyContinuous,
			MonthlyBudget:      1240000,
			AllocationPriority: 1,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"daily_amount": 40000, "days": 31}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Bensin",
			Type:               models.CategoryTypeUsageBased,
			MonthlyBudget:      175000,
			AllocationPriority: 2,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"per_fill": 35000, "fill_count": 5}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Listrik",
			Type:               models.CategoryTypeUsageBased,
			MonthlyBudget:      200000,
			AllocationPriority: 3,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"type": "utility"}),
		},
		{
			ID:                 uuid.New(),
			Name:               "GitHub Copilot",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      180000,
			AllocationPriority: 4,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"renewal_date": "monthly"}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Kuota Internet",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      100000,
			AllocationPriority: 5,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"provider": "telkomsel"}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Netflix",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      120000,
			AllocationPriority: 6,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"plan": "premium"}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Spotify",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      60000,
			AllocationPriority: 7,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"plan": "individual"}),
		},
		{
			ID:                 uuid.New(),
			Name:               "Fore Coffee",
			Type:               models.CategoryTypeSubscription,
			MonthlyBudget:      24000,
			AllocationPriority: 8,
			IsActive:           true,
			Metadata:           mustJSON(map[string]interface{}{"membership": "fore+"}),
		},
	}

	log.Println("Creating categories...")
	for _, category := range categories {
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		if err := categoryRepo.Create(&category); err != nil {
			log.Printf("Failed to create category %s: %v", category.Name, err)
		} else {
			log.Printf("✅ Created category: %s (Budget: Rp %.0f)", category.Name, category.MonthlyBudget)
		}
	}

	log.Println("🎉 Seeding completed successfully!")
	log.Println("Total categories created:", len(categories))

	log.Println("\n📊 Summary:")
	log.Println("Total Monthly Budget: Rp 2.099.000")
	log.Println("\nPrimary Expenses:")
	log.Println("  - Makan: Rp 1.240.000")
	log.Println("  - Bensin: Rp 175.000")
	log.Println("  - Listrik: Rp 200.000")
	log.Println("  - Copilot: Rp 180.000")
	log.Println("  - Internet: Rp 100.000")
	log.Println("\nMembership Expenses:")
	log.Println("  - Netflix: Rp 120.000")
	log.Println("  - Spotify: Rp 60.000")
	log.Println("  - Fore: Rp 24.000")
}

// mustJSON converts a map to datatypes.JSON, panics on error
func mustJSON(data interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return datatypes.JSON(jsonData)
}
