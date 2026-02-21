package integration

import (
	"encoding/json"
	"finance-tracking-app/handlers"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"finance-tracking-app/services"
	"finance-tracking-app/test/helpers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupAnalyticsIntegrationTest sets up environment for analytics integration testing
func setupAnalyticsIntegrationTest(t *testing.T) (*gin.Engine, func()) {
	helpers.SetupGinTestMode()

	// Setup test database
	db := helpers.SetupTestDB(t)

	// Initialize repositories
	expenseRepo := repositories.NewExpenseRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	// Initialize services
	analyticsService := services.NewAnalyticsService(expenseRepo, budgetRepo, categoryRepo)

	// Initialize handler
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Setup router
	router := gin.New()
	api := router.Group("/api/v1")
	{
		analytics := api.Group("/analytics")
		{
			analytics.GET("/spending-pattern", analyticsHandler.GetSpendingPattern)
			analytics.GET("/category-comparison", analyticsHandler.GetCategoryComparison)
			analytics.GET("/top-spending", analyticsHandler.GetTopSpending)
			analytics.GET("/budget-performance", analyticsHandler.GetBudgetPerformance)
		}
	}

	cleanup := func() {
		helpers.CleanupTestDB(t, db)
	}

	return router, cleanup
}

// TestAnalytics_SpendingPattern tests spending pattern analytics endpoint
func TestAnalytics_SpendingPattern(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	// Create test category first (would need helper to create category)
	req, _ := http.NewRequest("GET", "/api/v1/analytics/spending-pattern?months=6", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestAnalytics_CategoryComparison tests category comparison analytics
func TestAnalytics_CategoryComparison(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	url := fmt.Sprintf("/api/v1/analytics/category-comparison?month=%d&year=%d", month, year)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestAnalytics_TopSpending tests top spending categories analytics
func TestAnalytics_TopSpending(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	url := fmt.Sprintf("/api/v1/analytics/top-spending?month=%d&year=%d&limit=5", month, year)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestAnalytics_BudgetPerformance tests budget performance analytics
func TestAnalytics_BudgetPerformance(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	now := time.Now()
	year := now.Year()

	url := fmt.Sprintf("/api/v1/analytics/budget-performance?year=%d", year)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestAnalytics_InvalidParameters tests validation for invalid parameters
func TestAnalytics_InvalidParameters(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Invalid month in category comparison",
			url:  "/api/v1/analytics/category-comparison?month=13&year=2026",
		},
		{
			name: "Invalid year in category comparison",
			url:  "/api/v1/analytics/category-comparison?month=2&year=1999",
		},
		{
			name: "Invalid year in budget performance",
			url:  "/api/v1/analytics/budget-performance?year=1999",
		},
		{
			name: "Invalid limit in top spending",
			url:  "/api/v1/analytics/top-spending?month=2&year=2026&limit=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

// TestAnalytics_EndToEnd tests complete analytics workflow
func TestAnalytics_EndToEnd(t *testing.T) {
	router, cleanup := setupAnalyticsIntegrationTest(t)
	defer cleanup()

	// Get spending pattern
	req1, _ := http.NewRequest("GET", "/api/v1/analytics/spending-pattern?months=3", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Get category comparison
	req2, _ := http.NewRequest("GET", "/api/v1/analytics/category-comparison?month=2&year=2026", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Get top spending
	req3, _ := http.NewRequest("GET", "/api/v1/analytics/top-spending?month=2&year=2026&limit=10", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// Get budget performance
	req4, _ := http.NewRequest("GET", "/api/v1/analytics/budget-performance?year=2026", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}
