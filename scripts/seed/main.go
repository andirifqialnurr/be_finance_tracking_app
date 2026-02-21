package main

import (
	"encoding/json"
	"finance-tracking-app/config"
	"finance-tracking-app/database"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"fmt"
	"log"
	"strings"
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
	incomeRepo := repositories.NewIncomeRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)
	alertRepo := repositories.NewBudgetAlertRepository(db)

	log.Println("🌱 Starting comprehensive database seeder...")
	log.Println("⚠️  This will populate the database with sample data for testing")

	// Check if data already exists
	var existingCount int64
	db.Model(&models.ExpenseCategory{}).Count(&existingCount)
	if existingCount > 0 {
		log.Printf("⚠️  Database already has %d categories. Please run cleanup.go first!", existingCount)
		log.Println("❌ SEEDING ABORTED - Database is not empty")
		return
	}

	// Get current month and year
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// ==================== STEP 1: SEED EXPENSE CATEGORIES ====================
	log.Println("\n📂 Step 1: Creating Expense Categories...")

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

	for _, category := range categories {
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		if err := categoryRepo.Create(&category); err != nil {
			log.Printf("❌ Failed to create category %s: %v", category.Name, err)
		} else {
			log.Printf("✅ Created category: %s (Budget: Rp %.0f)", category.Name, category.MonthlyBudget)
		}
	}

	// ==================== STEP 2: SEED INCOMES ====================
	log.Println("\n💰 Step 2: Creating Income Entries...")

	incomes := []models.Income{
		{
			ID:          uuid.New(),
			Source:      "Gaji Februari 2026",
			Amount:      8000000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 1, 9, 0, 0, 0, time.Local),
			Description: "Gaji bulanan Transfer Bank",
		},
		{
			ID:          uuid.New(),
			Source:      "Freelance Project - Website",
			Amount:      3500000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 10, 14, 30, 0, 0, time.Local),
			Description: "Payment untuk website development client",
		},
		{
			ID:          uuid.New(),
			Source:      "Bonus Kinerja Q1",
			Amount:      2000000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 15, 10, 0, 0, 0, time.Local),
			Description: "Bonus performa kuartal pertama",
		},
	}

	var createdIncomes []models.Income
	for _, income := range incomes {
		income.CreatedAt = time.Now()
		income.UpdatedAt = time.Now()

		if err := incomeRepo.Create(&income); err != nil {
			log.Printf("❌ Failed to create income %s: %v", income.Source, err)
		} else {
			log.Printf("✅ Created income: %s (Rp %.0f)", income.Source, income.Amount)
			createdIncomes = append(createdIncomes, income)
		}
	}

	// ==================== STEP 3: CREATE BUDGET ALLOCATIONS ====================
	log.Println("\n📊 Step 3: Creating Budget Allocations from Incomes...")

	totalMonthlyBudget := 0.0
	for _, cat := range categories {
		totalMonthlyBudget += cat.MonthlyBudget
	}

	for _, income := range createdIncomes {
		log.Printf("  Allocating income: %s (Rp %.0f)", income.Source, income.Amount)

		for _, category := range categories {
			// Calculate allocation proportionally
			allocationRatio := category.MonthlyBudget / totalMonthlyBudget
			allocatedAmount := income.Amount * allocationRatio

			allocation := models.BudgetAllocation{
				ID:              uuid.New(),
				IncomeID:        income.ID,
				CategoryID:      category.ID,
				AllocatedAmount: allocatedAmount,
				Month:           currentMonth,
				Year:            currentYear,
				CreatedAt:       time.Now(),
			}

			if err := allocationRepo.Create(&allocation); err != nil {
				log.Printf("    ❌ Failed to allocate to %s: %v", category.Name, err)
			} else {
				log.Printf("    ✅ Allocated Rp %.0f to %s", allocatedAmount, category.Name)
			}

			// Update or create category budget
			budget, err := budgetRepo.FindByCategoryAndMonth(category.ID, currentMonth, currentYear)
			if err != nil || budget == nil {
				// Create new budget
				budget = &models.CategoryBudget{
					ID:              uuid.New(),
					CategoryID:      category.ID,
					Month:           currentMonth,
					Year:            currentYear,
					AllocatedAmount: allocatedAmount,
					SpentAmount:     0,
					RemainingAmount: allocatedAmount,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				if err := budgetRepo.Create(budget); err != nil {
					log.Printf("    ❌ Failed to create budget for %s: %v", category.Name, err)
				}
			} else {
				// Update existing budget
				budget.AllocatedAmount += allocatedAmount
				budget.RemainingAmount += allocatedAmount
				budget.UpdatedAt = time.Now()
				if err := budgetRepo.Update(budget); err != nil {
					log.Printf("    ❌ Failed to update budget for %s: %v", category.Name, err)
				}
			}
		}
	}

	// ==================== STEP 4: SEED EXPENSES ====================
	log.Println("\n💸 Step 4: Creating Expense Entries...")

	// Helper to find category ID by name
	findCategoryID := func(name string) uuid.UUID {
		for _, cat := range categories {
			if cat.Name == name {
				return cat.ID
			}
		}
		return uuid.Nil
	}

	expenses := []models.Expense{
		// Makan expenses
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Makan"),
			Amount:      45000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 2, 12, 30, 0, 0, time.Local),
			Description: "Makan siang Warteg + Minum",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Makan"),
			Amount:      38000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 3, 13, 0, 0, 0, time.Local),
			Description: "Nasi Padang + Es Teh",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Makan"),
			Amount:      52000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 4, 19, 30, 0, 0, time.Local),
			Description: "Ayam Geprek + Minuman",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Makan"),
			Amount:      75000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 5, 20, 0, 0, 0, time.Local),
			Description: "Makan berdua di Restoran",
		},

		// Bensin expenses
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Bensin"),
			Amount:      50000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 3, 8, 0, 0, 0, time.Local),
			Description: "Isi bensin Pertamax Shell",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Bensin"),
			Amount:      50000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 8, 7, 30, 0, 0, time.Local),
			Description: "Isi bensin Pertalite",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Bensin"),
			Amount:      50000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 14, 18, 0, 0, 0, time.Local),
			Description: "Isi bensin full tank",
		},

		// Subscriptions
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("GitHub Copilot"),
			Amount:      180000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 1, 0, 5, 0, 0, time.Local),
			Description: "GitHub Copilot Monthly Subscription",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Netflix"),
			Amount:      120000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 5, 0, 10, 0, 0, time.Local),
			Description: "Netflix Premium Plan",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Spotify"),
			Amount:      60000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 7, 0, 15, 0, 0, time.Local),
			Description: "Spotify Individual Plan",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Kuota Internet"),
			Amount:      100000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 2, 10, 0, 0, 0, time.Local),
			Description: "Paket Internet 50GB Telkomsel",
		},
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Fore Coffee"),
			Amount:      24000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 1, 0, 20, 0, 0, time.Local),
			Description: "Fore+ Membership Monthly",
		},

		// Listrik
		{
			ID:          uuid.New(),
			CategoryID:  findCategoryID("Listrik"),
			Amount:      185000,
			Date:        time.Date(currentYear, time.Month(currentMonth), 12, 15, 30, 0, 0, time.Local),
			Description: "Bayar tagihan listrik token PLN",
		},
	}

	for _, expense := range expenses {
		expense.CreatedAt = time.Now()
		expense.UpdatedAt = time.Now()

		if err := expenseRepo.Create(&expense); err != nil {
			log.Printf("❌ Failed to create expense: %v", err)
		} else {
			// Update budget
			budget, _ := budgetRepo.FindByCategoryAndMonth(expense.CategoryID, currentMonth, currentYear)
			if budget != nil {
				budget.SpentAmount += expense.Amount
				budget.RemainingAmount -= expense.Amount
				budget.UpdatedAt = time.Now()
				budgetRepo.Update(budget)
			}

			// Get category name for logging
			catName := ""
			for _, cat := range categories {
				if cat.ID == expense.CategoryID {
					catName = cat.Name
					break
				}
			}
			log.Printf("✅ Created expense: %s - Rp %.0f (%s)", catName, expense.Amount, expense.Description)
		}
	}

	// ==================== STEP 5: CREATE BUDGET ALERTS ====================
	log.Println("\n🔔 Step 5: Creating Budget Alerts...")

	alerts := []models.BudgetAlert{
		{
			ID:                  uuid.New(),
			CategoryID:          findCategoryID("Makan"),
			ThresholdPercentage: 80,
			IsEnabled:           true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  uuid.New(),
			CategoryID:          findCategoryID("Bensin"),
			ThresholdPercentage: 75,
			IsEnabled:           true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  uuid.New(),
			CategoryID:          findCategoryID("Listrik"),
			ThresholdPercentage: 90,
			IsEnabled:           true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}

	for _, alert := range alerts {
		if err := alertRepo.Create(&alert); err != nil {
			log.Printf("❌ Failed to create alert: %v", err)
		} else {
			catName := ""
			for _, cat := range categories {
				if cat.ID == alert.CategoryID {
					catName = cat.Name
					break
				}
			}
			log.Printf("✅ Created alert for %s (threshold: %d%%)", catName, alert.ThresholdPercentage)
		}
	}

	// ==================== SUMMARY ====================
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("🎉 SEEDING COMPLETED SUCCESSFULLY!")
	log.Println(strings.Repeat("=", 60))

	log.Printf("\n📊 Summary for %s %d:", time.Month(currentMonth).String(), currentYear)
	log.Printf("  ✅ Categories Created: %d", len(categories))
	log.Printf("  ✅ Incomes Created: %d", len(createdIncomes))
	log.Printf("  ✅ Expenses Created: %d", len(expenses))
	log.Printf("  ✅ Budget Alerts Created: %d", len(alerts))

	totalIncome := 0.0
	for _, income := range createdIncomes {
		totalIncome += income.Amount
	}

	totalExpense := 0.0
	for _, expense := range expenses {
		totalExpense += expense.Amount
	}

	log.Println("\n💰 Financial Summary:")
	log.Printf("  Total Income: Rp %s", formatCurrency(totalIncome))
	log.Printf("  Total Allocated: Rp %s", formatCurrency(totalMonthlyBudget))
	log.Printf("  Total Spent: Rp %s", formatCurrency(totalExpense))
	log.Printf("  Remaining Budget: Rp %s", formatCurrency(totalIncome-totalExpense))
	log.Printf("  Savings: Rp %s", formatCurrency(totalIncome-totalExpense))

	log.Println("\n📂 Categories Budget Breakdown:")
	for _, cat := range categories {
		budget, _ := budgetRepo.FindByCategoryAndMonth(cat.ID, currentMonth, currentYear)
		if budget != nil {
			usagePercent := (budget.SpentAmount / budget.AllocatedAmount) * 100
			log.Printf("  • %s:", cat.Name)
			log.Printf("      Allocated: Rp %s", formatCurrency(budget.AllocatedAmount))
			log.Printf("      Spent: Rp %s (%.1f%%)", formatCurrency(budget.SpentAmount), usagePercent)
			log.Printf("      Remaining: Rp %s", formatCurrency(budget.RemainingAmount))
		}
	}

	log.Println("\n✨ Database is now ready for testing!")
	log.Println("🌐 You can now test all endpoints via Swagger UI")
	log.Println("📍 http://localhost:8081/swagger/index.html")
	log.Println("")
}

// mustJSON converts a map to datatypes.JSON, panics on error
func mustJSON(data interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return datatypes.JSON(jsonData)
}

// formatCurrency formats a number as Indonesian Rupiah
func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
}
