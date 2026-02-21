package services

import (
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"sort"
	"time"
)

type TransactionService interface {
	GetTransactions(filters TransactionFilters) (*TransactionResponse, error)
}

type transactionService struct {
	incomeRepo  repositories.IncomeRepository
	expenseRepo repositories.ExpenseRepository
}

func NewTransactionService(
	incomeRepo repositories.IncomeRepository,
	expenseRepo repositories.ExpenseRepository,
) TransactionService {
	return &transactionService{
		incomeRepo:  incomeRepo,
		expenseRepo: expenseRepo,
	}
}

type TransactionFilters struct {
	Month      int
	Year       int
	StartDate  time.Time
	EndDate    time.Time
	Type       string // "income", "expense", "all"
	CategoryID string
	Page       int
	Limit      int
	Sort       string // "date_asc", "date_desc", "amount_asc", "amount_desc"
}

type Transaction struct {
	ID          string                  `json:"id"`
	Type        string                  `json:"type"` // "income" or "expense"
	Source      string                  `json:"source,omitempty"`
	Category    *models.ExpenseCategory `json:"category,omitempty"`
	Amount      float64                 `json:"amount"`
	Date        time.Time               `json:"date"`
	Description string                  `json:"description"`
	CreatedAt   time.Time               `json:"created_at"`
}

type TransactionSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type TransactionResponse struct {
	Data       []Transaction      `json:"data"`
	Pagination Pagination         `json:"pagination"`
	Summary    TransactionSummary `json:"summary"`
}

func (s *transactionService) GetTransactions(filters TransactionFilters) (*TransactionResponse, error) {
	var transactions []Transaction
	var totalIncome, totalExpense float64

	// Set defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Type == "" {
		filters.Type = "all"
	}

	// Fetch incomes
	if filters.Type == "all" || filters.Type == "income" {
		var incomes []models.Income
		var err error

		if filters.Month > 0 && filters.Year > 0 {
			incomes, err = s.incomeRepo.FindByMonthYear(filters.Month, filters.Year)
		} else if !filters.StartDate.IsZero() && !filters.EndDate.IsZero() {
			incomes, err = s.incomeRepo.FindByDateRange(filters.StartDate, filters.EndDate)
		} else {
			incomes, err = s.incomeRepo.FindAll(10000, 0) // Large limit to get all
		}

		if err != nil {
			return nil, err
		}

		for _, income := range incomes {
			totalIncome += income.Amount
			transactions = append(transactions, Transaction{
				ID:          income.ID.String(),
				Type:        "income",
				Source:      income.Source,
				Amount:      income.Amount,
				Date:        income.Date,
				Description: income.Description,
				CreatedAt:   income.CreatedAt,
			})
		}
	}

	// Fetch expenses
	if filters.Type == "all" || filters.Type == "expense" {
		var expenses []models.Expense
		var err error

		if filters.Month > 0 && filters.Year > 0 {
			expenses, err = s.expenseRepo.FindByMonthYear(filters.Month, filters.Year)
		} else if !filters.StartDate.IsZero() && !filters.EndDate.IsZero() {
			expenses, err = s.expenseRepo.FindByDateRange(filters.StartDate, filters.EndDate)
		} else {
			expenses, err = s.expenseRepo.FindAll(10000, 0) // Large limit to get all
		}

		if err != nil {
			return nil, err
		}

		for _, expense := range expenses {
			totalExpense += expense.Amount
			transactions = append(transactions, Transaction{
				ID:          expense.ID.String(),
				Type:        "expense",
				Category:    &expense.Category,
				Amount:      expense.Amount,
				Date:        expense.Date,
				Description: expense.Description,
				CreatedAt:   expense.CreatedAt,
			})
		}
	}

	// Sort transactions
	s.sortTransactions(&transactions, filters.Sort)

	// Calculate pagination
	total := len(transactions)
	totalPages := (total + filters.Limit - 1) / filters.Limit
	start := (filters.Page - 1) * filters.Limit
	end := start + filters.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedTransactions := transactions[start:end]

	return &TransactionResponse{
		Data: paginatedTransactions,
		Pagination: Pagination{
			Page:       filters.Page,
			Limit:      filters.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Summary: TransactionSummary{
			TotalIncome:  totalIncome,
			TotalExpense: totalExpense,
			NetBalance:   totalIncome - totalExpense,
		},
	}, nil
}

func (s *transactionService) sortTransactions(transactions *[]Transaction, sortBy string) {
	switch sortBy {
	case "date_asc":
		sort.Slice(*transactions, func(i, j int) bool {
			return (*transactions)[i].Date.Before((*transactions)[j].Date)
		})
	case "date_desc":
		sort.Slice(*transactions, func(i, j int) bool {
			return (*transactions)[i].Date.After((*transactions)[j].Date)
		})
	case "amount_asc":
		sort.Slice(*transactions, func(i, j int) bool {
			return (*transactions)[i].Amount < (*transactions)[j].Amount
		})
	case "amount_desc":
		sort.Slice(*transactions, func(i, j int) bool {
			return (*transactions)[i].Amount > (*transactions)[j].Amount
		})
	default:
		// Default: sort by date descending
		sort.Slice(*transactions, func(i, j int) bool {
			return (*transactions)[i].Date.After((*transactions)[j].Date)
		})
	}
}
