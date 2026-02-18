package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
// @Router /api/v1/budgets [get]
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
// @Router /api/v1/budgets/summary [get]
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
