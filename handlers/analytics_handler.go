package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetSpendingPattern godoc
// @Summary Get spending pattern for a category
// @Description Get spending pattern for a specific category over multiple months
// @Tags analytics
// @Produce json
// @Param category_id query string false "Category ID (optional)"
// @Param months query int false "Number of months (default: 6)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analytics/spending-pattern [get]
func (h *AnalyticsHandler) GetSpendingPattern(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	var categoryID uuid.UUID
	var err error

	if categoryIDStr != "" {
		categoryID, err = uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid category_id format",
			})
			return
		}
	}

	months, _ := strconv.Atoi(c.DefaultQuery("months", "6"))
	if months < 1 {
		months = 6
	}

	// If no category_id provided, use nil UUID to get all categories pattern
	if categoryIDStr == "" {
		categoryID = uuid.Nil
	}

	pattern, err := h.analyticsService.GetSpendingPattern(categoryID, months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Spending pattern retrieved successfully",
		Data:    pattern,
	})
}

// GetCategoryComparison godoc
// @Summary Compare category spending between months
// @Description Compare spending across categories for current vs previous month
// @Tags analytics
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analytics/category-comparison [get]
func (h *AnalyticsHandler) GetCategoryComparison(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "month and year are required",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid month parameter",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	comparison, err := h.analyticsService.GetCategoryComparison(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Category comparison retrieved successfully",
		Data:    comparison,
	})
}

// GetTopSpending godoc
// @Summary Get top spending categories
// @Description Get categories with highest spending for a specific month
// @Tags analytics
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year"
// @Param limit query int false "Number of categories to return (default: 5)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analytics/top-spending [get]
func (h *AnalyticsHandler) GetTopSpending(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "month and year are required",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid month parameter",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid limit parameter, must be greater than 0",
		})
		return
	}

	topSpending, err := h.analyticsService.GetTopSpending(month, year, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Top spending retrieved successfully",
		Data:    topSpending,
	})
}

// GetBudgetPerformance godoc
// @Summary Get budget performance analysis
// @Description Get budget performance metrics for the entire year
// @Tags analytics
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /analytics/budget-performance [get]
func (h *AnalyticsHandler) GetBudgetPerformance(c *gin.Context) {
	yearStr := c.Query("year")

	if yearStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "year is required",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid year parameter",
		})
		return
	}

	performance, err := h.analyticsService.GetBudgetPerformance(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Budget performance retrieved successfully",
		Data:    performance,
	})
}
