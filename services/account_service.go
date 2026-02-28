package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountSummaryItem is one row in the all-accounts summary response
type AccountSummaryItem struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	IncomeType         *string  `json:"income_type,omitempty"`
	Balance            float64  `json:"balance"`
	Color              string   `json:"color"`
	IsActive           bool     `json:"is_active"`
	GoalAmount         *float64 `json:"goal_amount,omitempty"`
	GoalLabel          *string  `json:"goal_label,omitempty"`
	ProgressPercentage *float64 `json:"progress_percentage,omitempty"` // only SAVINGS with goal_amount
}

// AccountsSummaryResponse is returned by GET /accounts/summary
type AccountsSummaryResponse struct {
	Accounts     []AccountSummaryItem `json:"accounts"`
	TotalBalance float64              `json:"total_balance"`
}

// TopUpResult holds the result of a TopUp operation
type TopUpResult struct {
	Income                  *models.Income            `json:"income"`
	Account                 *models.Account           `json:"account"`
	Allocations             []models.BudgetAllocation `json:"allocations,omitempty"`
	AlreadyAllocatedWarning bool                      `json:"already_allocated_warning,omitempty"`
}

// SpentResult holds the result of a Spent operation
type SpentResult struct {
	Expense       *models.Expense `json:"expense"`
	Account       *models.Account `json:"account"`
	BudgetWarning string          `json:"budget_warning,omitempty"`
}

type AccountService interface {
	CreateAccount(account *models.Account) (*models.Account, error)
	GetAllAccounts(accountType string, activeOnly bool) ([]models.Account, error)
	GetAccountsSummary() (*AccountsSummaryResponse, error)
	GetAccountByID(id uuid.UUID) (*models.Account, error)
	UpdateAccount(id uuid.UUID, updates *models.Account) (*models.Account, error)
	ArchiveAccount(id uuid.UUID) error
	TopUp(accountID uuid.UUID, income *models.Income) (*TopUpResult, error)
	Spent(accountID uuid.UUID, expense *models.Expense) (*SpentResult, error)
	GetBalanceSummary(id uuid.UUID, month, year int) (*models.AccountBalanceSummary, error)
}

type accountService struct {
	accountRepo    repositories.AccountRepository
	incomeRepo     repositories.IncomeRepository
	expenseRepo    repositories.ExpenseRepository
	categoryRepo   repositories.ExpenseCategoryRepository
	budgetRepo     repositories.CategoryBudgetRepository
	allocationRepo repositories.BudgetAllocationRepository
	alertRepo      repositories.BudgetAlertRepository
	db             *gorm.DB
}

func NewAccountService(
	accountRepo repositories.AccountRepository,
	incomeRepo repositories.IncomeRepository,
	expenseRepo repositories.ExpenseRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	allocationRepo repositories.BudgetAllocationRepository,
	alertRepo repositories.BudgetAlertRepository,
	db *gorm.DB,
) AccountService {
	return &accountService{
		accountRepo:    accountRepo,
		incomeRepo:     incomeRepo,
		expenseRepo:    expenseRepo,
		categoryRepo:   categoryRepo,
		budgetRepo:     budgetRepo,
		allocationRepo: allocationRepo,
		alertRepo:      alertRepo,
		db:             db,
	}
}

func (s *accountService) CreateAccount(account *models.Account) (*models.Account, error) {
	validTypes := map[string]bool{
		models.AccountTypeCard:    true,
		models.AccountTypeCash:    true,
		models.AccountTypeSavings: true,
	}
	if !validTypes[account.Type] {
		return nil, errors.New("invalid account type, must be CARD, CASH or SAVINGS")
	}

	if account.Type == models.AccountTypeCard {
		if account.IncomeType == nil || *account.IncomeType == "" {
			return nil, errors.New("income_type is required for CARD account")
		}
		validIncomeTypes := map[string]bool{
			models.CardIncomeTypeSalary:    true,
			models.CardIncomeTypeProject:   true,
			models.CardIncomeTypeFreelance: true,
			models.CardIncomeTypeBusiness:  true,
			models.CardIncomeTypeOther:     true,
		}
		if !validIncomeTypes[*account.IncomeType] {
			return nil, errors.New("invalid income_type for CARD, must be SALARY|PROJECT|FREELANCE|BUSINESS|OTHER")
		}
	}

	if account.Name == "" {
		return nil, errors.New("account name is required")
	}
	if account.Balance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	account.IsActive = true

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAllAccounts(accountType string, activeOnly bool) ([]models.Account, error) {
	return s.accountRepo.FindAll(accountType, activeOnly)
}

func (s *accountService) GetAccountsSummary() (*AccountsSummaryResponse, error) {
	accounts, err := s.accountRepo.FindAll("", true) // active only
	if err != nil {
		return nil, err
	}

	total := 0.0
	items := make([]AccountSummaryItem, 0, len(accounts))
	for _, a := range accounts {
		total += a.Balance
		item := AccountSummaryItem{
			ID:         a.ID.String(),
			Name:       a.Name,
			Type:       a.Type,
			IncomeType: a.IncomeType,
			Balance:    a.Balance,
			Color:      a.Color,
			IsActive:   a.IsActive,
			GoalAmount: a.GoalAmount,
			GoalLabel:  a.GoalLabel,
		}
		if a.Type == models.AccountTypeSavings && a.GoalAmount != nil && *a.GoalAmount > 0 {
			pct := (a.Balance / *a.GoalAmount) * 100
			item.ProgressPercentage = &pct
		}
		items = append(items, item)
	}

	return &AccountsSummaryResponse{Accounts: items, TotalBalance: total}, nil
}

func (s *accountService) GetAccountByID(id uuid.UUID) (*models.Account, error) {
	return s.accountRepo.FindByID(id)
}

func (s *accountService) UpdateAccount(id uuid.UUID, updates *models.Account) (*models.Account, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if !account.IsActive {
		return nil, errors.New("cannot update archived account")
	}

	if updates.Name != "" {
		account.Name = updates.Name
	}
	if updates.Color != "" {
		account.Color = updates.Color
	}
	if updates.Description != "" {
		account.Description = updates.Description
	}
	if account.Type == models.AccountTypeCard && updates.IncomeType != nil {
		account.IncomeType = updates.IncomeType
	}
	if account.Type == models.AccountTypeSavings {
		if updates.GoalAmount != nil {
			account.GoalAmount = updates.GoalAmount
		}
		if updates.GoalLabel != nil {
			account.GoalLabel = updates.GoalLabel
		}
	}

	return account, s.accountRepo.Update(account)
}

func (s *accountService) ArchiveAccount(id uuid.UUID) error {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return errors.New("account not found")
	}
	if !account.IsActive {
		return errors.New("account is already archived")
	}
	return s.accountRepo.Archive(id)
}

// TopUp adds income to an account and optionally auto-allocates budget (for SALARY card)
func (s *accountService) TopUp(accountID uuid.UUID, income *models.Income) (*TopUpResult, error) {
	if income.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if !account.IsActive {
		return nil, errors.New("account is archived and cannot receive funds")
	}

	result := &TopUpResult{}
	var allocations []models.BudgetAllocation

	err = s.db.Transaction(func(tx *gorm.DB) error {
		income.AccountID = accountID

		// Create income record
		if err := s.incomeRepo.Create(income); err != nil {
			return err
		}

		// Update account balance
		if err := s.accountRepo.UpdateBalance(accountID, income.Amount); err != nil {
			return err
		}

		// Salary auto-allocation
		if account.Type == models.AccountTypeCard &&
			account.IncomeType != nil &&
			*account.IncomeType == models.CardIncomeTypeSalary {

			month := int(income.Date.Month())
			year := income.Date.Year()

			// Check if already allocated for this month
			existingAllocs, _ := s.allocationRepo.FindByMonthYear(month, year)
			if len(existingAllocs) > 0 {
				result.AlreadyAllocatedWarning = true
			} else {
				var allocErr error
				allocations, allocErr = s.runAutoAllocation(income, month, year)
				if allocErr != nil {
					return allocErr
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload account with updated balance
	updatedAccount, _ := s.accountRepo.FindByID(accountID)
	result.Income = income
	result.Account = updatedAccount
	result.Allocations = allocations
	return result, nil
}

// Spent records an expense against an account; HARD balance check, SOFT budget check
func (s *accountService) Spent(accountID uuid.UUID, expense *models.Expense) (*SpentResult, error) {
	if expense.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return nil, errors.New("account not found")
	}
	if !account.IsActive {
		return nil, errors.New("account is archived")
	}
	// HARD balance check
	if account.Balance < expense.Amount {
		return nil, errors.New("insufficient account balance")
	}

	// Validate category
	category, err := s.categoryRepo.FindByID(expense.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if !category.IsActive {
		return nil, errors.New("category is not active")
	}

	result := &SpentResult{}
	month := int(expense.Date.Month())
	year := expense.Date.Year()

	err = s.db.Transaction(func(tx *gorm.DB) error {
		expense.AccountID = accountID

		// Create expense record
		if err := s.expenseRepo.Create(expense); err != nil {
			return err
		}

		// Deduct account balance
		if err := s.accountRepo.UpdateBalance(accountID, -expense.Amount); err != nil {
			return err
		}

		// SOFT budget update
		budget, err := s.budgetRepo.FindByCategoryAndMonth(expense.CategoryID, month, year)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.BudgetWarning = "no budget allocated for this category this month"
				return nil
			}
			return err
		}

		if budget.RemainingAmount < expense.Amount {
			result.BudgetWarning = "expense exceeds remaining budget for this category"
		}

		budget.SpentAmount += expense.Amount
		budget.RemainingAmount -= expense.Amount
		return s.budgetRepo.Update(budget)
	})

	if err != nil {
		return nil, err
	}

	expense.Category = *category
	updatedAccount, _ := s.accountRepo.FindByID(accountID)
	result.Expense = expense
	result.Account = updatedAccount
	return result, nil
}

func (s *accountService) GetBalanceSummary(id uuid.UUID, month, year int) (*models.AccountBalanceSummary, error) {
	return s.accountRepo.GetBalanceSummary(id, month, year)
}

// runAutoAllocation proportionally distributes income across active categories
func (s *accountService) runAutoAllocation(income *models.Income, month, year int) ([]models.BudgetAllocation, error) {
	categories, err := s.categoryRepo.FindAllActive()
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, nil
	}

	// Compute effective monthly budget per category
	totalBudget := 0.0
	effectiveBudgets := make(map[uuid.UUID]float64)
	for _, cat := range categories {
		var eb float64
		if cat.Type == models.CategoryTypeDailyContinuous && cat.DailyAmount != nil {
			days := daysInMonth(month, year)
			eb = *cat.DailyAmount * float64(days)
		} else {
			eb = cat.MonthlyBudget
		}
		effectiveBudgets[cat.ID] = eb
		totalBudget += eb
	}

	if totalBudget == 0 {
		return nil, nil
	}

	var allocations []models.BudgetAllocation
	for _, cat := range categories {
		eb := effectiveBudgets[cat.ID]
		allocationAmount := (eb / totalBudget) * income.Amount
		if allocationAmount <= 0 {
			continue
		}

		budget, err := s.budgetRepo.FindByCategoryAndMonth(cat.ID, month, year)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				days := daysInMonth(month, year)
				var effectiveDaily *float64
				var daysPtr *int
				if cat.Type == models.CategoryTypeDailyContinuous && cat.DailyAmount != nil {
					effectiveDaily = cat.DailyAmount
					daysPtr = &days
				}
				budget = &models.CategoryBudget{
					CategoryID:           cat.ID,
					Month:                month,
					Year:                 year,
					AllocatedAmount:      allocationAmount,
					SpentAmount:          0,
					RemainingAmount:      allocationAmount,
					EffectiveDailyAmount: effectiveDaily,
					DaysInMonth:          daysPtr,
				}
				if err := s.budgetRepo.Create(budget); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			budget.AllocatedAmount += allocationAmount
			budget.RemainingAmount += allocationAmount
			if err := s.budgetRepo.Update(budget); err != nil {
				return nil, err
			}
		}

		alloc := models.BudgetAllocation{
			IncomeID:        income.ID,
			CategoryID:      cat.ID,
			AllocatedAmount: allocationAmount,
			Month:           month,
			Year:            year,
		}
		if err := s.allocationRepo.Create(&alloc); err != nil {
			return nil, err
		}
		alloc.Category = cat
		allocations = append(allocations, alloc)
	}
	return allocations, nil
}

// daysInMonth returns the number of days in a given month/year
func daysInMonth(month, year int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}
