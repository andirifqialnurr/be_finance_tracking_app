package main

import (
	"finance-tracking-app/config"
	"finance-tracking-app/database"
	"finance-tracking-app/handlers"
	"finance-tracking-app/repositories"
	"finance-tracking-app/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()

	// Get database instance
	db := database.GetDB()

	// Initialize repositories
	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	// Initialize services
	incomeService := services.NewIncomeService(incomeRepo, categoryRepo, budgetRepo, allocationRepo, db)
	categoryService := services.NewCategoryService(categoryRepo)
	expenseService := services.NewExpenseService(expenseRepo, categoryRepo, budgetRepo, db)
	budgetService := services.NewBudgetService(budgetRepo, incomeRepo, allocationRepo)

	// Initialize handlers
	incomeHandler := handlers.NewIncomeHandler(incomeService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware (for development)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "finance-tracking-api",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Income routes
		incomes := v1.Group("/incomes")
		{
			incomes.POST("", incomeHandler.CreateIncome)
			incomes.GET("", incomeHandler.GetIncomes)
			incomes.GET("/:id", incomeHandler.GetIncomeByID)
			incomes.DELETE("/:id", incomeHandler.DeleteIncome)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Expense routes
		expenses := v1.Group("/expenses")
		{
			expenses.POST("", expenseHandler.CreateExpense)
			expenses.GET("", expenseHandler.GetExpenses)
			expenses.GET("/:id", expenseHandler.GetExpenseByID)
			expenses.DELETE("/:id", expenseHandler.DeleteExpense)
		}

		// Budget routes
		budgets := v1.Group("/budgets")
		{
			budgets.GET("", budgetHandler.GetBudgets)
			budgets.GET("/summary", budgetHandler.GetBudgetSummary)
		}
	}

	// Start server
	port := config.AppConfig.ServerPort
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("🚀 Server running on http://localhost%s\n", addr)
	fmt.Printf("📊 API Documentation: http://localhost%s/api/v1\n", addr)
	fmt.Printf("💚 Health Check: http://localhost%s/health\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
