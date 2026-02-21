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

type IncomeHandler struct {
	incomeService services.IncomeService
}

func NewIncomeHandler(incomeService services.IncomeService) *IncomeHandler {
	return &IncomeHandler{
		incomeService: incomeService,
	}
}

// CreateIncome godoc
// @Summary Create new income
// @Description Create a new income and automatically allocate budget to categories
// @Tags income
// @Accept json
// @Produce json
// @Param income body CreateIncomeRequest true "Income data"
// @Success 201 {object} CreateIncomeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /incomes [post]
func (h *IncomeHandler) CreateIncome(c *gin.Context) {
	var req CreateIncomeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate
	if req.Source == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "source is required",
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "amount must be greater than 0",
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

	income := &models.Income{
		Source:      req.Source,
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
	}

	createdIncome, allocations, err := h.incomeService.CreateIncome(income)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Build allocation summary
	allocationSummaries := make([]AllocationSummary, 0, len(allocations))
	for _, alloc := range allocations {
		allocationSummaries = append(allocationSummaries, AllocationSummary{
			Category:  alloc.Category.Name,
			Allocated: alloc.AllocatedAmount,
		})
	}

	response := CreateIncomeResponse{
		Income:      createdIncome,
		Allocations: allocationSummaries,
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Income created and budget allocated successfully",
		Data:    response,
	})
}

// GetIncomes godoc
// @Summary Get all incomes
// @Description Get all incomes with pagination
// @Tags income
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /incomes [get]
func (h *IncomeHandler) GetIncomes(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	var incomes []models.Income
	var err error

	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)
		incomes, err = h.incomeService.GetIncomesByMonthYear(month, year)
	} else {
		incomes, err = h.incomeService.GetAllIncomes(limit, offset)
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
		Message: "Incomes retrieved successfully",
		Data:    incomes,
	})
}

// GetIncomeByID godoc
// @Summary Get income by ID
// @Description Get a specific income by ID
// @Tags income
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /incomes/{id} [get]
func (h *IncomeHandler) GetIncomeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	income, err := h.incomeService.GetIncomeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "income not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Income retrieved successfully",
		Data:    income,
	})
}

// UpdateIncome godoc
// @Summary Update income
// @Description Update income by ID and re-calculate allocations
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Income ID"
// @Param income body CreateIncomeRequest true "Updated income data"
// @Success 200 {object} CreateIncomeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /incomes/{id} [patch]
func (h *IncomeHandler) UpdateIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	var req CreateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate
	if req.Source == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "source is required",
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "amount must be greater than 0",
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

	income := &models.Income{
		Source:      req.Source,
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
	}

	updatedIncome, allocations, err := h.incomeService.UpdateIncome(id, income)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Build allocation summary
	allocationSummaries := make([]AllocationSummary, 0, len(allocations))
	for _, alloc := range allocations {
		allocationSummaries = append(allocationSummaries, AllocationSummary{
			Category:  alloc.Category.Name,
			Allocated: alloc.AllocatedAmount,
		})
	}

	response := CreateIncomeResponse{
		Income:      updatedIncome,
		Allocations: allocationSummaries,
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Income updated and allocations adjusted successfully",
		Data:    response,
	})
}

// DeleteIncome godoc
// @Summary Delete income
// @Description Delete income by ID
// @Tags income
// @Produce json
// @Param id path string true "Income ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /incomes/{id} [delete]
func (h *IncomeHandler) DeleteIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	if err := h.incomeService.DeleteIncome(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Income deleted successfully",
	})
}

// Request/Response DTOs
type CreateIncomeRequest struct {
	Source      string  `json:"source" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description"`
}

type CreateIncomeResponse struct {
	Income      *models.Income      `json:"income"`
	Allocations []AllocationSummary `json:"allocations"`
}

type AllocationSummary struct {
	Category  string  `json:"category"`
	Allocated float64 `json:"allocated"`
}
