package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BudgetHandler struct {
	budgetService services.BudgetService
}

func NewBudgetHandler(budgetService services.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// GetBudgets godoc
// @Summary Get budgets by month and year
// @Description Get all budget summaries for a specific month and year
// @Tags budgets
// @Produce json
// @Param month query int false "Month (1-12)" default(current month)
// @Param year query int false "Year" default(current year)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /budgets [get]
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	// Default to current month/year if not provided
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	budgets, err := h.budgetService.GetBudgetsByMonthYear(month, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Budgets retrieved successfully",
		Data:    budgets,
	})
}

// GetBudgetSummary godoc
// @Summary Get overall budget summary
// @Description Get overall budget summary for a specific month and year
// @Tags budgets
// @Produce json
// @Param month query int false "Month (1-12)" default(current month)
// @Param year query int false "Year" default(current year)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /budgets/summary [get]
func (h *BudgetHandler) GetBudgetSummary(c *gin.Context) {
	// Default to current month/year if not provided
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	summary, err := h.budgetService.GetBudgetSummary(month, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Budget summary retrieved successfully",
		Data:    summary,
	})
}

// ReallocateBudget godoc
// @Summary Reallocate budget between categories
// @Description Manually reallocate budget from one category to another
// @Tags budgets
// @Accept json
// @Produce json
// @Param reallocation body ReallocateBudgetRequest true "Reallocation data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /budgets/reallocate [post]
func (h *BudgetHandler) ReallocateBudget(c *gin.Context) {
	var req ReallocateBudgetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Parse category IDs
	fromCategoryID, err := uuid.Parse(req.FromCategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid from_category_id format",
		})
		return
	}

	toCategoryID, err := uuid.Parse(req.ToCategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid to_category_id format",
		})
		return
	}

	// Default to current month/year if not provided
	now := time.Now()
	month := req.Month
	year := req.Year
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	reallocationReq := &services.BudgetReallocationRequest{
		FromCategoryID: fromCategoryID,
		ToCategoryID:   toCategoryID,
		Amount:         req.Amount,
		Reason:         req.Reason,
		Month:          month,
		Year:           year,
	}

	response, err := h.budgetService.ReallocateBudget(reallocationReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Budget reallocated successfully",
		Data:    response,
	})
}

// GetReallocations godoc
// @Summary Get budget reallocation history
// @Description Get history of budget reallocations for a specific month/year
// @Tags budgets
// @Produce json
// @Param month query int false "Month (1-12)" default(current month)
// @Param year query int false "Year" default(current year)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /budgets/reallocations [get]
func (h *BudgetHandler) GetReallocations(c *gin.Context) {
	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	reallocations, err := h.budgetService.GetReallocations(month, year)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Reallocations retrieved successfully",
		Data:    reallocations,
	})
}

// CancelReallocation godoc
// @Summary Cancel a budget reallocation
// @Description Cancel/delete a budget reallocation and restore budgets to previous state
// @Tags budgets
// @Produce json
// @Param id path string true "Reallocation ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /budgets/reallocate/{id} [delete]
func (h *BudgetHandler) CancelReallocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid reallocation id format",
		})
		return
	}

	if err := h.budgetService.CancelReallocation(id); err != nil {
		if err.Error() == "reallocation not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Reallocation cancelled successfully and budgets restored",
		Data:    nil,
	})
}

// Request DTOs
type ReallocateBudgetRequest struct {
	FromCategoryID string  `json:"from_category_id" binding:"required"`
	ToCategoryID   string  `json:"to_category_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	Reason         string  `json:"reason"`
	Month          int     `json:"month"`
	Year           int     `json:"year"`
}
