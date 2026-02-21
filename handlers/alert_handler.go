package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlertHandler struct {
	alertService services.AlertService
}

func NewAlertHandler(alertService services.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// GetAlerts godoc
// @Summary Get all budget alerts
// @Description Get all budget alerts with current status
// @Tags alerts
// @Produce json
// @Param status query string false "Filter by status: active or all (default: all)"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /alerts [get]
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	status := c.DefaultQuery("status", "all")

	alerts, err := h.alertService.GetAllAlerts(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Alerts retrieved successfully",
		Data:    alerts,
	})
}

// CreateAlert godoc
// @Summary Create a new budget alert
// @Description Create a new budget alert for a category
// @Tags alerts
// @Accept json
// @Produce json
// @Param alert body CreateAlertRequest true "Alert data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req CreateAlertRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid category_id format",
		})
		return
	}

	alert := &models.BudgetAlert{
		CategoryID:          categoryID,
		ThresholdPercentage: req.ThresholdPercentage,
		IsEnabled:           req.IsEnabled,
	}

	createdAlert, err := h.alertService.CreateAlert(alert)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Alert created successfully",
		Data:    createdAlert,
	})
}

// UpdateAlert godoc
// @Summary Update budget alert
// @Description Update budget alert configuration
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param alert body UpdateAlertRequest true "Alert data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alerts/{id} [patch]
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	var req UpdateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	alert := &models.BudgetAlert{
		ThresholdPercentage: req.ThresholdPercentage,
		IsEnabled:           req.IsEnabled,
	}

	updatedAlert, err := h.alertService.UpdateAlert(id, alert)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Alert updated successfully",
		Data:    updatedAlert,
	})
}

// DeleteAlert godoc
// @Summary Delete budget alert
// @Description Delete budget alert by ID
// @Tags alerts
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alerts/{id} [delete]
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid ID format",
		})
		return
	}

	if err := h.alertService.DeleteAlert(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Alert deleted successfully",
	})
}

// Request DTOs
type CreateAlertRequest struct {
	CategoryID          string `json:"category_id" binding:"required"`
	ThresholdPercentage int    `json:"threshold_percentage" binding:"required,min=0,max=100"`
	IsEnabled           bool   `json:"is_enabled"`
}

type UpdateAlertRequest struct {
	ThresholdPercentage int  `json:"threshold_percentage" binding:"required,min=0,max=100"`
	IsEnabled           bool `json:"is_enabled"`
}
