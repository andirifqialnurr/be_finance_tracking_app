package integration

import (
	"finance-tracking-app/handlers"
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

// setupReportIntegrationTest sets up environment for report integration testing
func setupReportIntegrationTest(t *testing.T) (*gin.Engine, func()) {
	helpers.SetupGinTestMode()

	// Setup test database
	db := helpers.SetupTestDB(t)

	// Initialize repositories
	incomeRepo := repositories.NewIncomeRepository(db)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	budgetRepo := repositories.NewCategoryBudgetRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	allocationRepo := repositories.NewBudgetAllocationRepository(db)

	// Initialize services
	reportService := services.NewReportService(incomeRepo, expenseRepo, budgetRepo, categoryRepo, allocationRepo)

	// Initialize handler
	reportHandler := handlers.NewReportHandler(reportService)

	// Setup router
	router := gin.New()
	api := router.Group("/api/v1")
	{
		reports := api.Group("/reports")
		{
			reports.GET("/monthly", reportHandler.GetMonthlyReport)
			reports.GET("/monthly/export", reportHandler.ExportMonthlyReport)
			reports.GET("/yearly", reportHandler.GetYearlyReport)
			reports.GET("/yearly/export", reportHandler.ExportYearlyReport)
		}
	}

	cleanup := func() {
		helpers.CleanupTestDB(t, db)
	}

	return router, cleanup
}

// TestMonthlyReportExport_PDF tests PDF export functionality for monthly report
func TestMonthlyReportExport_PDF(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Test PDF export
	req, _ := http.NewRequest("GET", "/api/v1/reports/monthly/export?month=2&year=2026&format=pdf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".pdf")
	assert.Greater(t, len(w.Body.Bytes()), 0, "PDF file should have content")
}

// TestMonthlyReportExport_Excel tests Excel export functionality for monthly report
func TestMonthlyReportExport_Excel(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Test Excel export
	req, _ := http.NewRequest("GET", "/api/v1/reports/monthly/export?month=2&year=2026&format=excel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".xlsx")
	assert.Greater(t, len(w.Body.Bytes()), 0, "Excel file should have content")
}

// TestYearlyReportExport_PDF tests PDF export functionality for yearly report
func TestYearlyReportExport_PDF(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Test PDF export
	req, _ := http.NewRequest("GET", "/api/v1/reports/yearly/export?year=2026&format=pdf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".pdf")
}

// TestYearlyReportExport_Excel tests Excel export functionality for yearly report
func TestYearlyReportExport_Excel(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Test Excel export
	req, _ := http.NewRequest("GET", "/api/v1/reports/yearly/export?year=2026&format=excel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".xlsx")
}

// TestReportExport_InvalidFormat tests validation for invalid export format
func TestReportExport_InvalidFormat(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Test invalid format
	req, _ := http.NewRequest("GET", "/api/v1/reports/monthly/export?month=2&year=2026&format=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestReportExport_InvalidParameters tests validation for invalid parameters
func TestReportExport_InvalidParameters(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Invalid month",
			url:  "/api/v1/reports/monthly/export?month=13&year=2026&format=pdf",
		},
		{
			name: "Invalid year",
			url:  "/api/v1/reports/monthly/export?month=2&year=1999&format=pdf",
		},
		{
			name: "Invalid year for yearly report",
			url:  "/api/v1/reports/yearly/export?year=1999&format=pdf",
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

// TestReportExport_WithRealisticData tests export with realistic data
func TestReportExport_WithRealisticData(t *testing.T) {
	router, cleanup := setupReportIntegrationTest(t)
	defer cleanup()

	// Use current month and year for testing
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	// Test PDF export with realistic timing
	url := fmt.Sprintf("/api/v1/reports/monthly/export?month=%d&year=%d&format=pdf", month, year)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
}
