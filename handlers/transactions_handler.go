package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetTransactions godoc
// @Summary Get combined transaction history
// @Description Get combined income and expense transactions with filtering and pagination
// @Tags transactions
// @Produce json
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param type query string false "Transaction type: income, expense, or all (default: all)"
// @Param category_id query string false "Filter by category ID (for expenses)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20)"
// @Param sort query string false "Sort by: date_asc, date_desc, amount_asc, amount_desc (default: date_desc)"
// @Success 200 {object} services.TransactionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	filters := services.TransactionFilters{
		Type: c.DefaultQuery("type", "all"),
		Sort: c.DefaultQuery("sort", "date_desc"),
	}

	// Parse month and year
	if monthStr := c.Query("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid month parameter",
			})
			return
		}
		filters.Month = month
	}

	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 2000 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid year parameter",
			})
			return
		}
		filters.Year = year
	}

	// Parse date range
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid start_date format, use YYYY-MM-DD",
			})
			return
		}
		filters.StartDate = startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid end_date format, use YYYY-MM-DD",
			})
			return
		}
		filters.EndDate = endDate
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filters.Page = page
	filters.Limit = limit

	// Category filter
	filters.CategoryID = c.Query("category_id")

	// Validate type
	if filters.Type != "income" && filters.Type != "expense" && filters.Type != "all" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid type parameter, must be: income, expense, or all",
		})
		return
	}

	// Validate sort
	validSorts := []string{"date_asc", "date_desc", "amount_asc", "amount_desc"}
	isValidSort := false
	for _, validSort := range validSorts {
		if filters.Sort == validSort {
			isValidSort = true
			break
		}
	}
	if !isValidSort {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid sort parameter, must be: date_asc, date_desc, amount_asc, or amount_desc",
		})
		return
	}

	response, err := h.transactionService.GetTransactions(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    response,
	})
}
