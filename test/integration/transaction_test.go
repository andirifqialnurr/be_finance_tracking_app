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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTransactionIntegrationTest sets up environment for transaction integration testing
func setupTransactionIntegrationTest(t *testing.T) (*gin.Engine, func()) {
	helpers.SetupGinTestMode()

	// Setup test database
	db := helpers.SetupTestDB(t)

	// Initialize repositories
	incomeRepo := repositories.NewIncomeRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)

	// Initialize services
	transactionService := services.NewTransactionService(incomeRepo, expenseRepo)

	// Initialize handler
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup router
	router := gin.New()
	api := router.Group("/api/v1")
	{
		transactions := api.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
		}
	}

	cleanup := func() {
		helpers.CleanupTestDB(t, db)
	}

	return router, cleanup
}

// TestTransactions_GetAllTypes tests getting all transaction types
func TestTransactions_GetAllTypes(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/api/v1/transactions?month=2&year=2026", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestTransactions_FilterByType tests filtering transactions by type
func TestTransactions_FilterByType(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name         string
		transType    string
		expectedCode int
	}{
		{
			name:         "Filter income only",
			transType:    "income",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Filter expense only",
			transType:    "expense",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get all transactions",
			transType:    "all",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?type="+tt.transType+"&month=2&year=2026", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response models.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.True(t, response.Success)
		})
	}
}

// TestTransactions_Pagination tests pagination functionality
func TestTransactions_Pagination(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name  string
		page  string
		limit string
	}{
		{
			name:  "First page with 10 items",
			page:  "1",
			limit: "10",
		},
		{
			name:  "Second page with 20 items",
			page:  "2",
			limit: "20",
		},
		{
			name:  "Third page with 5 items",
			page:  "3",
			limit: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestTransactions_Sorting tests sorting functionality
func TestTransactions_Sorting(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name string
		sort string
	}{
		{
			name: "Sort by date ascending",
			sort: "date_asc",
		},
		{
			name: "Sort by date descending",
			sort: "date_desc",
		},
		{
			name: "Sort by amount ascending",
			sort: "amount_asc",
		},
		{
			name: "Sort by amount descending",
			sort: "amount_desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/transactions?sort="+tt.sort, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestTransactions_DateRangeFilter tests flexible date range filtering
func TestTransactions_DateRangeFilter(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/api/v1/transactions?start_date=2026-01-01&end_date=2026-12-31", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestTransactions_InvalidParameters tests validation for invalid parameters
func TestTransactions_InvalidParameters(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Invalid month",
			url:  "/api/v1/transactions?month=13&year=2026",
		},
		{
			name: "Invalid year",
			url:  "/api/v1/transactions?month=2&year=1999",
		},
		{
			name: "Invalid type",
			url:  "/api/v1/transactions?type=invalid",
		},
		{
			name: "Invalid sort option",
			url:  "/api/v1/transactions?sort=invalid_sort",
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

// TestTransactions_CompleteWorkflow tests combined filter, sort, and pagination
func TestTransactions_CompleteWorkflow(t *testing.T) {
	router, cleanup := setupTransactionIntegrationTest(t)
	defer cleanup()

	// Complex query with multiple parameters
	url := "/api/v1/transactions?type=all&month=2&year=2026&page=1&limit=20&sort=date_desc"
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Verify response has pagination info and summary
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "data")
	assert.Contains(t, data, "pagination")
	assert.Contains(t, data, "summary")
}
