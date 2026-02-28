package handlers

import (
	"finance-tracking-app/models"
	"finance-tracking-app/services"
	"finance-tracking-app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduledFundHandler struct {
	sfService services.ScheduledFundService
}

func NewScheduledFundHandler(sfService services.ScheduledFundService) *ScheduledFundHandler {
	return &ScheduledFundHandler{sfService: sfService}
}

// CreateScheduledFund godoc
// @Summary Create a recurring scheduled fund
// @Tags scheduled-funds
// @Accept json
// @Produce json
// @Param request body CreateScheduledFundRequest true "Scheduled fund data"
// @Success 201 {object} models.Response
// @Router /scheduled-funds [post]
func (h *ScheduledFundHandler) CreateScheduledFund(c *gin.Context) {
	var req CreateScheduledFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid account_id"})
		return
	}

	sf := &models.ScheduledFund{
		AccountID:    accountID,
		ScheduleType: req.ScheduleType,
		Amount:       req.Amount,
		DayOfMonth:   req.DayOfMonth,
		Description:  req.Description,
	}

	if req.FromAccountID != "" {
		fromID, err := uuid.Parse(req.FromAccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid from_account_id"})
			return
		}
		sf.FromAccountID = &fromID
	}

	created, err := h.sfService.CreateScheduledFund(sf)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "Scheduled fund created successfully", Data: created})
}

// GetScheduledFunds godoc
// @Summary Get all scheduled funds
// @Tags scheduled-funds
// @Produce json
// @Param active_only query bool false "Filter active only"
// @Success 200 {object} models.Response
// @Router /scheduled-funds [get]
func (h *ScheduledFundHandler) GetScheduledFunds(c *gin.Context) {
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "false"))

	funds, err := h.sfService.GetAll(activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Scheduled funds retrieved", Data: funds})
}

// GetScheduledFundByID godoc
// @Summary Get a scheduled fund by ID
// @Tags scheduled-funds
// @Produce json
// @Param id path string true "Scheduled Fund ID"
// @Success 200 {object} models.Response
// @Router /scheduled-funds/{id} [get]
func (h *ScheduledFundHandler) GetScheduledFundByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	sf, err := h.sfService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Success: false, Error: "scheduled fund not found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Scheduled fund retrieved", Data: sf})
}

// UpdateScheduledFund godoc
// @Summary Update a scheduled fund
// @Tags scheduled-funds
// @Accept json
// @Produce json
// @Param id path string true "Scheduled Fund ID"
// @Param request body UpdateScheduledFundRequest true "Update data"
// @Success 200 {object} models.Response
// @Router /scheduled-funds/{id} [patch]
func (h *ScheduledFundHandler) UpdateScheduledFund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	var req UpdateScheduledFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	updates := &models.ScheduledFund{
		Amount:      req.Amount,
		DayOfMonth:  req.DayOfMonth,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	updated, err := h.sfService.UpdateScheduledFund(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Scheduled fund updated", Data: updated})
}

// DeleteScheduledFund godoc
// @Summary Delete a scheduled fund
// @Tags scheduled-funds
// @Produce json
// @Param id path string true "Scheduled Fund ID"
// @Success 200 {object} models.Response
// @Router /scheduled-funds/{id} [delete]
func (h *ScheduledFundHandler) DeleteScheduledFund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	if err := h.sfService.DeleteScheduledFund(id); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Scheduled fund deleted"})
}

// ---- DTOs ----

type CreateScheduledFundRequest struct {
	AccountID     string  `json:"account_id" binding:"required"`    // destination account
	FromAccountID string  `json:"from_account_id"`                  // source for TRANSFER
	ScheduleType  string  `json:"schedule_type" binding:"required"` // TOP_UP | TRANSFER
	Amount        float64 `json:"amount" binding:"required"`
	DayOfMonth    int     `json:"day_of_month"` // 1-31; 0 = last day
	Description   string  `json:"description"`
}

type UpdateScheduledFundRequest struct {
	Amount      float64 `json:"amount"`
	DayOfMonth  int     `json:"day_of_month"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
}

// ensure time import is used
var _ = time.Now
var _ = utils.FormatCurrency
