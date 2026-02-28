package services

import (
	"finance-tracking-app/models"
	"time"

	"gorm.io/gorm"
)

type StatisticsService interface {
	GetMonthlyStats(year int) (*MonthlyStatsResponse, error)
	GetOverview(month, year int) (*OverviewResponse, error)
}

type statisticsService struct {
	db *gorm.DB
}

func NewStatisticsService(db *gorm.DB) StatisticsService {
	return &statisticsService{db: db}
}

// MonthlyData holds aggregated income/expense for one month
type MonthlyData struct {
	Month        int     `json:"month"`
	MonthName    string  `json:"month_name"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetSavings   float64 `json:"net_savings"`
}

// CategoryBreakdown holds expense breakdown per category for a month
type CategoryBreakdown struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSpent   float64 `json:"total_spent"`
	Percentage   float64 `json:"percentage"`
}

// MonthlyStatsResponse is the full yearly stats response
type MonthlyStatsResponse struct {
	Year       int           `json:"year"`
	Monthly    []MonthlyData `json:"monthly"`
	YearTotal  float64       `json:"year_total_income"`
	YearSpent  float64       `json:"year_total_expense"`
	YearSaving float64       `json:"year_net_savings"`
}

// OverviewResponse is the current snapshot
type OverviewResponse struct {
	Month             int                 `json:"month"`
	Year              int                 `json:"year"`
	TotalBalance      float64             `json:"total_balance"`
	TotalIncome       float64             `json:"total_income"`
	TotalExpense      float64             `json:"total_expense"`
	TotalTransfers    float64             `json:"total_transfers_out"`
	NetSavings        float64             `json:"net_savings"`
	Accounts          []AccountSummaryRow `json:"accounts"`
	CategoryBreakdown []CategoryBreakdown `json:"category_breakdown"`
}

// AccountSummaryRow is a single account row in overview
type AccountSummaryRow struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

func (s *statisticsService) GetMonthlyStats(year int) (*MonthlyStatsResponse, error) {
	result := &MonthlyStatsResponse{Year: year}

	for m := 1; m <= 12; m++ {
		start := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)

		var income float64
		s.db.Model(&models.Income{}).
			Where("date >= ? AND date < ?", start, end).
			Select("COALESCE(SUM(amount), 0)").Scan(&income)

		var expense float64
		s.db.Model(&models.Expense{}).
			Where("date >= ? AND date < ?", start, end).
			Select("COALESCE(SUM(amount), 0)").Scan(&expense)

		md := MonthlyData{
			Month:        m,
			MonthName:    time.Month(m).String(),
			TotalIncome:  income,
			TotalExpense: expense,
			NetSavings:   income - expense,
		}
		result.Monthly = append(result.Monthly, md)
		result.YearTotal += income
		result.YearSpent += expense
	}
	result.YearSaving = result.YearTotal - result.YearSpent
	return result, nil
}

func (s *statisticsService) GetOverview(month, year int) (*OverviewResponse, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	resp := &OverviewResponse{Month: month, Year: year}

	// Total balance across all active accounts
	s.db.Model(&models.Account{}).
		Where("is_active = ?", true).
		Select("COALESCE(SUM(balance), 0)").Scan(&resp.TotalBalance)

	// Income this month
	s.db.Model(&models.Income{}).
		Where("date >= ? AND date < ?", start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&resp.TotalIncome)

	// Expense this month
	s.db.Model(&models.Expense{}).
		Where("date >= ? AND date < ?", start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&resp.TotalExpense)

	// Transfers out this month
	s.db.Model(&models.AccountTransfer{}).
		Where("transfer_date >= ? AND transfer_date < ?", start, end).
		Select("COALESCE(SUM(amount), 0)").Scan(&resp.TotalTransfers)

	resp.NetSavings = resp.TotalIncome - resp.TotalExpense

	// List all active accounts
	var accounts []models.Account
	s.db.Where("is_active = ?", true).Find(&accounts)
	for _, a := range accounts {
		resp.Accounts = append(resp.Accounts, AccountSummaryRow{
			ID:      a.ID.String(),
			Name:    a.Name,
			Type:    a.Type,
			Balance: a.Balance,
		})
	}

	// Category breakdown
	type catRow struct {
		CategoryID   string
		CategoryName string
		TotalSpent   float64
	}
	var rows []catRow
	s.db.Model(&models.Expense{}).
		Select("expenses.category_id, expense_categories.name as category_name, SUM(expenses.amount) as total_spent").
		Joins("JOIN expense_categories ON expense_categories.id = expenses.category_id").
		Where("expenses.date >= ? AND expenses.date < ?", start, end).
		Group("expenses.category_id, expense_categories.name").
		Scan(&rows)

	for _, row := range rows {
		pct := 0.0
		if resp.TotalExpense > 0 {
			pct = (row.TotalSpent / resp.TotalExpense) * 100
		}
		resp.CategoryBreakdown = append(resp.CategoryBreakdown, CategoryBreakdown{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			TotalSpent:   row.TotalSpent,
			Percentage:   pct,
		})
	}

	return resp, nil
}
