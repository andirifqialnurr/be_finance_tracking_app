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

type TransferHandler struct {
	transferService services.TransferService
}

func NewTransferHandler(transferService services.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// CreateTransfer godoc
// @Summary Transfer funds between accounts
// @Tags transfers
// @Accept json
// @Produce json
// @Param request body CreateTransferRequest true "Transfer data"
// @Success 201 {object} models.Response
// @Router /transfers [post]
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	fromID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid from_account_id"})
		return
	}

	toID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid to_account_id"})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid date format, use RFC3339"})
		return
	}

	transfer, err := h.transferService.CreateTransfer(fromID, toID, req.Amount, req.Note, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "Transfer completed successfully", Data: transfer})
}

// GetTransfers godoc
// @Summary Get all transfers
// @Tags transfers
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} models.Response
// @Router /transfers [get]
func (h *TransferHandler) GetTransfers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transfers, err := h.transferService.GetAllTransfers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Transfers retrieved successfully", Data: transfers})
}

// GetTransferByID godoc
// @Summary Get transfer by ID
// @Tags transfers
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} models.Response
// @Router /transfers/{id} [get]
func (h *TransferHandler) GetTransferByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	transfer, err := h.transferService.GetTransferByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Success: false, Error: "transfer not found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Transfer retrieved successfully", Data: transfer})
}

// GetTransfersByAccount godoc
// @Summary Get transfers for a specific account
// @Tags transfers
// @Produce json
// @Param account_id query string true "Account ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} models.Response
// @Router /transfers/account [get]
func (h *TransferHandler) GetTransfersByAccount(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	if accountIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "account_id is required"})
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid account_id"})
		return
	}

	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	transfers, err := h.transferService.GetTransfersByAccount(accountID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Account transfers retrieved", Data: transfers})
}

// CancelTransfer godoc
// @Summary Cancel a transfer (within 24 hours)
// @Tags transfers
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} models.Response
// @Router /transfers/{id} [delete]
func (h *TransferHandler) CancelTransfer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	if err := h.transferService.CancelTransfer(id); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Transfer cancelled and balances restored"})
}

// ---- DTOs ----

type CreateTransferRequest struct {
	FromAccountID string  `json:"from_account_id" binding:"required"`
	ToAccountID   string  `json:"to_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Note          string  `json:"note"`
	Date          string  `json:"date" binding:"required"`
}
