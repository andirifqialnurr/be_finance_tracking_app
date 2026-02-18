package handlers

import (
	"encoding/json"
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create new expense category
// @Description Create a new expense category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body CreateCategoryRequest true "Category data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	metadata, err := json.Marshal(req.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid metadata format",
		})
		return
	}

	category := &models.ExpenseCategory{
		Name:               req.Name,
		Type:               req.Type,
		MonthlyBudget:      req.MonthlyBudget,
		AllocationPriority: req.AllocationPriority,
		Metadata:           datatypes.JSON(metadata),
	}

	createdCategory, err := h.categoryService.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Category created successfully",
		Data:    createdCategory,
	})
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all expense categories
// @Tags categories
// @Produce json
// @Param active query boolean false "Filter active only"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	var categories []models.ExpenseCategory
	var err error

	if activeOnly {
		categories, err = h.categoryService.GetActiveCategories()
	} else {
		categories, err = h.categoryService.GetAllCategories()
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
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "category not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update an existing category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body UpdateCategoryRequest true "Category data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	metadata, err := json.Marshal(req.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid metadata format",
		})
		return
	}

	category := &models.ExpenseCategory{
		Name:               req.Name,
		Type:               req.Type,
		MonthlyBudget:      req.MonthlyBudget,
		AllocationPriority: req.AllocationPriority,
		Metadata:           datatypes.JSON(metadata),
	}

	updatedCategory, err := h.categoryService.UpdateCategory(id, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Category updated successfully",
		Data:    updatedCategory,
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category by ID (hard delete)
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Category deleted successfully",
	})
}

// Request DTOs
type CreateCategoryRequest struct {
	Name               string                 `json:"name" binding:"required"`
	Type               string                 `json:"type" binding:"required"`
	MonthlyBudget      float64                `json:"monthly_budget" binding:"required"`
	AllocationPriority int                    `json:"allocation_priority"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type UpdateCategoryRequest struct {
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	MonthlyBudget      float64                `json:"monthly_budget"`
	AllocationPriority int                    `json:"allocation_priority"`
	Metadata           map[string]interface{} `json:"metadata"`
}
