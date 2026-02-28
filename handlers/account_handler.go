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

type AccountHandler struct {
	accountService  services.AccountService
	transferService services.TransferService
}

func NewAccountHandler(accountService services.AccountService, transferService services.TransferService) *AccountHandler {
	return &AccountHandler{accountService: accountService, transferService: transferService}
}

// CreateAccount godoc
// @Summary Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body CreateAccountRequest true "Account data"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	account := &models.Account{
		Name:        req.Name,
		Type:        req.Type,
		IncomeType:  req.IncomeType,
		Balance:     req.InitialBalance,
		Color:       req.Color,
		Description: req.Description,
		GoalAmount:  req.GoalAmount,
		GoalLabel:   req.GoalLabel,
	}

	created, err := h.accountService.CreateAccount(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: "Account created successfully", Data: created})
}

// GetAccounts godoc
// @Summary Get all accounts
// @Tags accounts
// @Produce json
// @Param type query string false "Account type filter (CARD|CASH|SAVINGS)"
// @Param active_only query bool false "Only active accounts"
// @Success 200 {object} models.Response
// @Router /accounts [get]

// GetAllAccountsSummary godoc
// @Summary Get all accounts overview with total balance
// @Tags accounts
// @Produce json
// @Success 200 {object} models.Response
// @Router /accounts/summary [get]
func (h *AccountHandler) GetAllAccountsSummary(c *gin.Context) {
	summary, err := h.accountService.GetAccountsSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Accounts summary retrieved", Data: summary})
}
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accountType := c.Query("type")
	activeOnly, _ := strconv.ParseBool(c.DefaultQuery("active_only", "true"))

	accounts, err := h.accountService.GetAllAccounts(accountType, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Accounts retrieved successfully", Data: accounts})
}

// GetAccountByID godoc
// @Summary Get account by ID
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} models.Response
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	account, err := h.accountService.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Success: false, Error: "account not found"})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Account retrieved successfully", Data: account})
}

// GetAccountSummary godoc
// @Summary Get account balance summary for a specific month
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} models.Response
// @Router /accounts/{id}/summary [get]
func (h *AccountHandler) GetAccountSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	summary, err := h.accountService.GetBalanceSummary(id, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Account summary retrieved", Data: summary})
}

// UpdateAccount godoc
// @Summary Update account details
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body UpdateAccountRequest true "Account update data"
// @Success 200 {object} models.Response
// @Router /accounts/{id} [patch]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	updates := &models.Account{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		IncomeType:  req.IncomeType,
		GoalAmount:  req.GoalAmount,
		GoalLabel:   req.GoalLabel,
	}

	updated, err := h.accountService.UpdateAccount(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Account updated successfully", Data: updated})
}

// ArchiveAccount godoc
// @Summary Archive an account
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} models.Response
// @Router /accounts/{id}/archive [post]
func (h *AccountHandler) ArchiveAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	if err := h.accountService.ArchiveAccount(id); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Message: "Account archived successfully"})
}

// TopUp godoc
// @Summary Add funds (income) to an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body TopUpRequest true "Top-up data"
// @Success 201 {object} models.Response
// @Router /accounts/{id}/topup [post]
func (h *AccountHandler) TopUp(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid date format, use RFC3339"})
		return
	}

	income := &models.Income{
		AccountID:   id,
		Source:      req.Source,
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
	}

	result, err := h.accountService.TopUp(id, income)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	msg := "Funds added successfully"
	if result.AlreadyAllocatedWarning {
		msg = "Funds added. Budget already allocated for this month."
	} else if len(result.Allocations) > 0 {
		msg = "Funds added and budget auto-allocated (SALARY)"
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: msg, Data: result})
}

// Spent godoc
// @Summary Record an expense from an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body SpentRequest true "Expense data"
// @Success 201 {object} models.Response
// @Router /accounts/{id}/spent [post]
func (h *AccountHandler) Spent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid ID format"})
		return
	}

	var req SpentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid category_id format"})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid date format, use RFC3339"})
		return
	}

	expense := &models.Expense{
		AccountID:   id,
		CategoryID:  categoryID,
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
	}

	result, err := h.accountService.Spent(id, expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	msg := "Expense recorded successfully"
	if result.BudgetWarning != "" {
		msg = "Expense recorded. Warning: " + result.BudgetWarning
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Message: msg, Data: result})
}

// ---- Request/Response DTOs ----

type CreateAccountRequest struct {
	Name           string   `json:"name" binding:"required"`
	Type           string   `json:"type" binding:"required"` // CARD|CASH|SAVINGS
	IncomeType     *string  `json:"income_type"`             // required for CARD
	InitialBalance float64  `json:"initial_balance"`
	Color          string   `json:"color"`
	Description    string   `json:"description"`
	GoalAmount     *float64 `json:"goal_amount"` // SAVINGS
	GoalLabel      *string  `json:"goal_label"`  // SAVINGS
}

type UpdateAccountRequest struct {
	Name        string   `json:"name"`
	Color       string   `json:"color"`
	Description string   `json:"description"`
	IncomeType  *string  `json:"income_type"` // CARD only
	GoalAmount  *float64 `json:"goal_amount"` // SAVINGS
	GoalLabel   *string  `json:"goal_label"`  // SAVINGS
}

type TopUpRequest struct {
	Source      string  `json:"source" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description"`
}

type SpentRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	Description string  `json:"description"`
}

type AccountTransferRequest struct {
	ToAccountID string  `json:"to_account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Note        string  `json:"note"`
	Date        string  `json:"date" binding:"required"`
}

// TransferFromAccount godoc
// @Summary Transfer funds from this account to another
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Source Account ID"
// @Param request body AccountTransferRequest true "Transfer data"
// @Success 201 {object} models.Response
// @Router /accounts/{id}/transfer [post]
func (h *AccountHandler) TransferFromAccount(c *gin.Context) {
	fromID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid account ID"})
		return
	}

	var req AccountTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: err.Error()})
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

// GetAccountTransfers godoc
// @Summary Get transfer history for a specific account
// @Tags accounts
// @Produce json
// @Param id path string true "Account ID"
// @Param month query int false "Month (1-12)"
// @Param year query int false "Year"
// @Success 200 {object} models.Response
// @Router /accounts/{id}/transfers [get]
func (h *AccountHandler) GetAccountTransfers(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "invalid account ID"})
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
