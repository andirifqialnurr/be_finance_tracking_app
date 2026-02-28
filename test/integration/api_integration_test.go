package integration

import (
	"encoding/json"
	"finance-tracking-app/handlers"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"finance-tracking-app/services"
	"finance-tracking-app/test/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupIntegrationTest sets up a complete application environment for integration testing
func setupIntegrationTest(t *testing.T) *gin.Engine {
	helpers.SetupGinTestMode()

	// Setup test database
	db := helpers.SetupTestDB(t)

	// Initialize repositories
	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)
	reallocationRepo := repositories.NewBudgetReallocationRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	// alertRepo := repositories.NewBudgetAlertRepository(db) // Not used in this test

	// Initialize services
	incomeService := services.NewIncomeService(incomeRepo, accountRepo, categoryRepo, budgetRepo, allocationRepo, db)
	categoryService := services.NewCategoryService(categoryRepo)
	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, accountRepo, db)
	budgetService := services.NewBudgetService(budgetRepo, incomeRepo, allocationRepo, reallocationRepo, categoryRepo, db)

	// Initialize handlers
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)

	// Setup router
	router := gin.New()

	api := router.Group("/api/v1")
	{
		// Income routes
		api.POST("/incomes", incomeHandler.CreateIncome)
		api.GET("/incomes", incomeHandler.GetIncomes)
		api.GET("/incomes/:id", incomeHandler.GetIncomeByID)

		// Category routes
		api.POST("/categories", categoryHandler.CreateCategory)
		api.GET("/categories", categoryHandler.GetCategories)
		api.GET("/categories/:id", categoryHandler.GetCategoryByID)
		api.PUT("/categories/:id", categoryHandler.UpdateCategory)
		api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Expense routes
		api.POST("/expenses", expenseHandler.CreateExpense)
		api.GET("/expenses", expenseHandler.GetExpenses)

		// Budget routes
		api.GET("/budgets", budgetHandler.GetBudgets)
		api.GET("/budgets/summary", budgetHandler.GetBudgetSummary)
	}

	return router
}

// TestEndToEnd_CompleteWorkflow tests the complete workflow from category creation to expense tracking
func TestEndToEnd_CompleteWorkflow(t *testing.T) {
	router := setupIntegrationTest(t)

	// Step 1: Create categories
	t.Run("Create Categories", func(t *testing.T) {
		categories := []map[string]interface{}{
			{
				"name":                "Makan",
				"type":                "DAILY_CONTINUOUS",
				"monthly_budget":      1000000,
				"allocation_priority": 1,
			},
			{
				"name":                "Transport",
				"type":                "USAGE_BASED",
				"monthly_budget":      500000,
				"allocation_priority": 2,
			},
		}

		for _, cat := range categories {
			req := helpers.MakeJSONRequest(t, "POST", "/api/v1/categories", cat)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
		}
	})

	// Step 2: Verify categories were created
	var categoryIDs []string
	t.Run("Get All Categories", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/api/v1/categories", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)

		// Extract category IDs
		categoriesData, _ := json.Marshal(response.Data)
		var categories []map[string]interface{}
		json.Unmarshal(categoriesData, &categories)

		assert.Len(t, categories, 2)

		for _, cat := range categories {
			categoryIDs = append(categoryIDs, cat["id"].(string))
		}
	})

	// Step 3: Create income (should auto-allocate to categories)
	t.Run("Create Income with Auto-Allocation", func(t *testing.T) {
		incomeData := map[string]interface{}{
			"source":      "Salary",
			"amount":      3000000,
			"date":        time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
			"description": "Monthly salary",
		}

		req := helpers.MakeJSONRequest(t, "POST", "/api/v1/incomes", incomeData)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)

		// Verify auto-allocation happened
		responseData, _ := json.Marshal(response.Data)
		var incomeResponse map[string]interface{}
		json.Unmarshal(responseData, &incomeResponse)

		allocations := incomeResponse["allocations"].([]interface{})
		assert.Len(t, allocations, 2, "Should allocate to both categories")
	})

	// Step 4: Check budget summary
	t.Run("Verify Budget Summary After Income", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/api/v1/budgets/summary?month=2&year=2026", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)

		// Verify budget data
		summaryData, _ := json.Marshal(response.Data)
		var summary map[string]interface{}
		json.Unmarshal(summaryData, &summary)

		assert.InDelta(t, 3000000, summary["total_income"].(float64), 0.01)
		assert.InDelta(t, 3000000, summary["total_allocated"].(float64), 1.0)
		assert.Equal(t, 0.0, summary["total_spent"].(float64))
	})

	// Step 5: Create expense (should deduct from budget)
	t.Run("Create Expense", func(t *testing.T) {
		if len(categoryIDs) == 0 {
			t.Skip("No category IDs available")
		}

		expenseData := map[string]interface{}{
			"category_id": categoryIDs[0],
			"amount":      50000,
			"date":        time.Date(2026, 2, 17, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
			"description": "Lunch",
		}

		req := helpers.MakeJSONRequest(t, "POST", "/api/v1/expenses", expenseData)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)
	})

	// Step 6: Verify budget was deducted
	t.Run("Verify Budget After Expense", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/api/v1/budgets/summary?month=2&year=2026", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)

		summaryData, _ := json.Marshal(response.Data)
		var summary map[string]interface{}
		json.Unmarshal(summaryData, &summary)

		assert.Equal(t, 50000.0, summary["total_spent"].(float64))
		assert.InDelta(t, 2950000, summary["total_remaining"].(float64), 1.0)
	})

	// Step 7: Verify category budgets detail
	t.Run("Get Category Budgets Detail", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/api/v1/budgets?month=2&year=2026", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)

		budgetsData, _ := json.Marshal(response.Data)
		var budgets []map[string]interface{}
		json.Unmarshal(budgetsData, &budgets)

		assert.Len(t, budgets, 2)

		// Verify at least one category has spent amount
		hasSpentCategory := false
		for _, budget := range budgets {
			if budget["spent_amount"].(float64) > 0 {
				hasSpentCategory = true
				assert.Less(t, budget["remaining_amount"].(float64), budget["allocated_amount"].(float64))
			}
		}
		assert.True(t, hasSpentCategory, "At least one category should have spent amount")
	})
}

// TestConcurrentRequests tests if the API handles concurrent requests properly
func TestConcurrentRequests(t *testing.T) {
	router := setupIntegrationTest(t)

	// Create a category first
	categoryData := map[string]interface{}{
		"name":                "Test Category",
		"type":                "SUBSCRIPTION",
		"monthly_budget":      100000,
		"allocation_priority": 1,
	}

	req := helpers.MakeJSONRequest(t, "POST", "/api/v1/categories", categoryData)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Make concurrent GET requests
	t.Run("Concurrent GET requests", func(t *testing.T) {
		done := make(chan bool, 5)

		for i := 0; i < 5; i++ {
			go func() {
				req := helpers.MakeJSONRequest(t, "GET", "/api/v1/categories", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
				done <- true
			}()
		}

		// Wait for all requests to complete
		for i := 0; i < 5; i++ {
			<-done
		}
	})
}
