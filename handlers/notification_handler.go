package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// RegisterDevice godoc
// @Summary Register a device for push notifications (store player ID)
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body RegisterDeviceRequest true "Device registration data"
// @Success 200 {object} models.Response
// @Router /notifications/register-device [post]
func (h *NotificationHandler) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}
	// Simply validate and acknowledge; player ID is persisted per notification setting
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Device registered successfully",
		Data:    gin.H{"onesignal_player_id": req.OneSignalPlayerID},
	})
}

// CreateNotificationSetting godoc
// @Summary Create a push notification setting
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body CreateNotificationSettingRequest true "Notification setting data"
// @Success 201 {object} models.Response
// @Router /notifications/settings [post]
func (h *NotificationHandler) CreateNotificationSetting(c *gin.Context) {
	var req CreateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	ns := &models.NotificationSetting{
		Type:              req.Type,
		Title:             req.Title,
		Body:              req.Body,
		DayOfMonth:        req.DayOfMonth,
		TimeOfDay:         req.TimeOfDay,
		OneSignalPlayerID: req.OneSignalPlayerID,
	}

	created, err := h.notificationService.CreateSetting(ns)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "Notification setting created", Data: created})
}

// GetNotificationSettings godoc
// @Summary Get all notification settings
// @Tags notifications
// @Produce json
// @Success 200 {object} models.Response
// @Router /notifications/settings [get]
func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
	settings, err := h.notificationService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Notification settings retrieved", Data: settings})
}

// GetNotificationSettingByID godoc
// @Summary Get notification setting by ID
// @Tags notifications
// @Produce json
// @Param id path string true "Setting ID"
// @Success 200 {object} models.Response
// @Router /notifications/settings/{id} [get]
func (h *NotificationHandler) GetNotificationSettingByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	setting, err := h.notificationService.GetSettingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Success: false, Error: "notification setting not found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Notification setting retrieved", Data: setting})
}

// UpdateNotificationSetting godoc
// @Summary Update a notification setting
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Setting ID"
// @Param request body UpdateNotificationSettingRequest true "Update data"
// @Success 200 {object} models.Response
// @Router /notifications/settings/{id} [patch]
func (h *NotificationHandler) UpdateNotificationSetting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	var req UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	updates := &models.NotificationSetting{
		Title:             req.Title,
		Body:              req.Body,
		DayOfMonth:        req.DayOfMonth,
		TimeOfDay:         req.TimeOfDay,
		IsEnabled:         req.IsEnabled,
		OneSignalPlayerID: req.OneSignalPlayerID,
	}

	updated, err := h.notificationService.UpdateSetting(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Notification setting updated", Data: updated})
}

// DeleteNotificationSetting godoc
// @Summary Delete a notification setting
// @Tags notifications
// @Produce json
// @Param id path string true "Setting ID"
// @Success 200 {object} models.Response
// @Router /notifications/settings/{id} [delete]
func (h *NotificationHandler) DeleteNotificationSetting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	if err := h.notificationService.DeleteSetting(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Notification setting deleted"})
}

// TestNotification godoc
// @Summary Send a test push notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body TestNotificationRequest true "Test notification data"
// @Success 200 {object} models.Response
// @Router /notifications/test [post]
func (h *NotificationHandler) TestNotification(c *gin.Context) {
	var req TestNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	if err := h.notificationService.TestNotification(req.OneSignalPlayerID, req.Title, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Test notification sent successfully"})
}

// ---- DTOs ----

type RegisterDeviceRequest struct {
	OneSignalPlayerID string `json:"onesignal_player_id" binding:"required"`
}

type CreateNotificationSettingRequest struct {
	Type              string `json:"type" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Body              string `json:"body" binding:"required"`
	DayOfMonth        int    `json:"day_of_month"`                   // 0 = every day
	TimeOfDay         string `json:"time_of_day" binding:"required"` // HH:MM
	OneSignalPlayerID string `json:"onesignal_player_id" binding:"required"`
}

type UpdateNotificationSettingRequest struct {
	Title             string `json:"title"`
	Body              string `json:"body"`
	DayOfMonth        int    `json:"day_of_month"`
	TimeOfDay         string `json:"time_of_day"`
	IsEnabled         bool   `json:"is_enabled"`
	OneSignalPlayerID string `json:"onesignal_player_id"`
}

type TestNotificationRequest struct {
	OneSignalPlayerID string `json:"onesignal_player_id" binding:"required"`
	Title             string `json:"title" binding:"required"`
	Body              string `json:"body" binding:"required"`
}
