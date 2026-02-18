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

type ExpenseHandler struct {
	expenseService services.ExpenseService
}

func NewExpenseHandler(expenseService services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense godoc
// @Summary Create new expense
// @Description Create a new expense and update category budget
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body CreateExpenseRequest true "Expense data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req CreateExpenseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Parse category ID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid category ID format",
		})
		return
	}

	// Parse date
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid date format, use RFC3339",
		})
		return
	}

	expense := &models.Expense{
		CategoryID:  categoryID,
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
	}

	createdExpense, err := h.expenseService.CreateExpense(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Expense created successfully",
		Data:    createdExpense,
	})
}

// GetExpenses godoc
// @Summary Get all expenses
// @Description Get all expenses with pagination and filtering
// @Tags expenses
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param category_id query string false "Filter by category ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/expenses [get]
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	categoryIDStr := c.Query("category_id")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	var expenses []models.Expense
	var err error

	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)
		expenses, err = h.expenseService.GetExpensesByMonthYear(month, year)
	} else if categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "invalid category ID format",
			})
			return
		}
		expenses, err = h.expenseService.GetExpensesByCategory(categoryID, limit, offset)
	} else {
		expenses, err = h.expenseService.GetAllExpenses(limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Expenses retrieved successfully",
		Data:    expenses,
	})
}

// GetExpenseByID godoc
// @Summary Get expense by ID
// @Description Get a specific expense by ID
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/expenses/{id} [get]
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	expense, err := h.expenseService.GetExpenseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "expense not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Expense retrieved successfully",
		Data:    expense,
	})
}

// DeleteExpense godoc
// @Summary Delete expense
// @Description Delete expense by ID
// @Tags expenses
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	if err := h.expenseService.DeleteExpense(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Expense deleted successfully",
	})
}

// Request DTOs
type CreateExpenseRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description"`
}
