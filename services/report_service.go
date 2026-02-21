package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"fmt"
)

type ReportService interface {
	GetMonthlyReport(month, year int) (*MonthlyReport, error)
	GetYearlyReport(year int) (*YearlyReport, error)
}

type reportService struct {
	incomeRepo     repositories.IncomeRepository
	expenseRepo    repositories.ExpenseRepository
	budgetRepo     repositories.CategoryBudgetRepository
	categoryRepo   repositories.ExpenseCategoryRepository
	allocationRepo repositories.BudgetAllocationRepository
}

func NewReportService(
	incomeRepo repositories.IncomeRepository,
	expenseRepo repositories.ExpenseRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
	allocationRepo repositories.BudgetAllocationRepository,
) ReportService {
	return &reportService{
		incomeRepo:     incomeRepo,
		expenseRepo:    expenseRepo,
		budgetRepo:     budgetRepo,
		categoryRepo:   categoryRepo,
		allocationRepo: allocationRepo,
	}
}

// Monthly Report Structures
type IncomeSource struct {
	Source string  `json:"source"`
	Amount float64 `json:"amount"`
}

type IncomeSummary struct {
	Total   float64        `json:"total"`
	Sources []IncomeSource `json:"sources"`
}

type ExpenseByCategory struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

type ExpenseSummary struct {
	Total      float64             `json:"total"`
	ByCategory []ExpenseByCategory `json:"by_category"`
}

type CategorySummary struct {
	Category  string  `json:"category"`
	Allocated float64 `json:"allocated"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
}

type ReportBudgetSummary struct {
	TotalAllocated    float64           `json:"total_allocated"`
	TotalSpent        float64           `json:"total_spent"`
	TotalRemaining    float64           `json:"total_remaining"`
	CategoriesSummary []CategorySummary `json:"categories_summary"`
}

type MonthlyReport struct {
	Month             int                 `json:"month"`
	Year              int                 `json:"year"`
	Income            IncomeSummary       `json:"income"`
	Expenses          ExpenseSummary      `json:"expenses"`
	Budget            ReportBudgetSummary `json:"budget"`
	Savings           float64             `json:"savings"`
	Unallocated       float64             `json:"unallocated"`
	TopCategory       string              `json:"top_category"`
	TopCategoryAmount float64             `json:"top_category_amount"`
	Insights          []string            `json:"insights"`
}

// Yearly Report Structures
type MonthBreakdown struct {
	Month    int     `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Savings  float64 `json:"savings"`
}

type CategoryYearlySummary struct {
	Category       string  `json:"category"`
	TotalSpent     float64 `json:"total_spent"`
	AverageMonthly float64 `json:"average_monthly"`
	HighestMonth   int     `json:"highest_month"`
	HighestAmount  float64 `json:"highest_amount"`
}

type YearlyTrends struct {
	AverageMonthlyIncome  float64 `json:"average_monthly_income"`
	AverageMonthlyExpense float64 `json:"average_monthly_expense"`
	HighestIncomeMonth    int     `json:"highest_income_month"`
	HighestExpenseMonth   int     `json:"highest_expense_month"`
	TotalSavingsRate      float64 `json:"total_savings_rate"`
}

type YearlyReport struct {
	Year                  int                     `json:"year"`
	TotalIncome           float64                 `json:"total_income"`
	TotalExpenses         float64                 `json:"total_expenses"`
	TotalSavings          float64                 `json:"total_savings"`
	MonthlyBreakdown      []MonthBreakdown        `json:"monthly_breakdown"`
	CategoryYearlySummary []CategoryYearlySummary `json:"category_yearly_summary"`
	Trends                YearlyTrends            `json:"trends"`
}

func (s *reportService) GetMonthlyReport(month, year int) (*MonthlyReport, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	// Get incomes for the month
	incomes, err := s.incomeRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	// Build income summary
	var totalIncome float64
	sources := make([]IncomeSource, 0)
	sourceMap := make(map[string]float64)

	for _, income := range incomes {
		totalIncome += income.Amount
		sourceMap[income.Source] += income.Amount
	}

	for source, amount := range sourceMap {
		sources = append(sources, IncomeSource{
			Source: source,
			Amount: amount,
		})
	}

	// Get expenses for the month
	expenses, err := s.expenseRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	// Build expense summary
	var totalExpense float64
	categoryExpenseMap := make(map[string]struct {
		amount float64
		count  int
	})

	for _, expense := range expenses {
		totalExpense += expense.Amount
		entry := categoryExpenseMap[expense.Category.Name]
		entry.amount += expense.Amount
		entry.count++
		categoryExpenseMap[expense.Category.Name] = entry
	}

	expenseByCategory := make([]ExpenseByCategory, 0)
	for category, data := range categoryExpenseMap {
		expenseByCategory = append(expenseByCategory, ExpenseByCategory{
			Category: category,
			Amount:   data.amount,
			Count:    data.count,
		})
	}

	// Get budget summary
	budgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		budgets = []models.CategoryBudget{}
	}

	var totalAllocated, totalSpent, totalRemaining float64
	categoriesSummary := make([]CategorySummary, 0)

	for _, budget := range budgets {
		totalAllocated += budget.AllocatedAmount
		totalSpent += budget.SpentAmount
		totalRemaining += budget.RemainingAmount

		categoriesSummary = append(categoriesSummary, CategorySummary{
			Category:  budget.Category.Name,
			Allocated: budget.AllocatedAmount,
			Spent:     budget.SpentAmount,
			Remaining: budget.RemainingAmount,
		})
	}

	// Calculate metrics
	savings := totalIncome - totalExpense
	unallocated := totalIncome - totalAllocated

	// Find top category
	topCategory := ""
	topAmount := 0.0
	for category, data := range categoryExpenseMap {
		if data.amount > topAmount {
			topAmount = data.amount
			topCategory = category
		}
	}

	// Generate insights
	insights := make([]string, 0)

	if totalAllocated > 0 {
		usagePercentage := (totalSpent / totalAllocated) * 100
		insights = append(insights, fmt.Sprintf("Total budget usage: %.1f%%", usagePercentage))
	}

	if totalIncome > 0 {
		savingsRate := (savings / totalIncome) * 100
		insights = append(insights, fmt.Sprintf("Savings rate: %.1f%%", savingsRate))
	}

	if topCategory != "" {
		insights = append(insights, fmt.Sprintf("Top spending category: %s (Rp %.0f)", topCategory, topAmount))
	}

	return &MonthlyReport{
		Month: month,
		Year:  year,
		Income: IncomeSummary{
			Total:   totalIncome,
			Sources: sources,
		},
		Expenses: ExpenseSummary{
			Total:      totalExpense,
			ByCategory: expenseByCategory,
		},
		Budget: ReportBudgetSummary{
			TotalAllocated:    totalAllocated,
			TotalSpent:        totalSpent,
			TotalRemaining:    totalRemaining,
			CategoriesSummary: categoriesSummary,
		},
		Savings:           savings,
		Unallocated:       unallocated,
		TopCategory:       topCategory,
		TopCategoryAmount: topAmount,
		Insights:          insights,
	}, nil
}

func (s *reportService) GetYearlyReport(year int) (*YearlyReport, error) {
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	var totalIncome, totalExpenses float64
	monthlyBreakdown := make([]MonthBreakdown, 0, 12)

	// Track category spending across the year
	categoryYearlyMap := make(map[string]struct {
		total         float64
		highestMonth  int
		highestAmount float64
		monthCount    int
	})

	var highestIncomeMonth, highestExpenseMonth int
	var highestIncome, highestExpense float64

	// Iterate through all months
	for month := 1; month <= 12; month++ {
		// Get incomes
		incomes, err := s.incomeRepo.FindByMonthYear(month, year)
		if err != nil {
			incomes = []models.Income{}
		}

		var monthIncome float64
		for _, income := range incomes {
			monthIncome += income.Amount
		}

		// Get expenses
		expenses, err := s.expenseRepo.FindByMonthYear(month, year)
		if err != nil {
			expenses = []models.Expense{}
		}

		var monthExpense float64
		categoryMonthMap := make(map[string]float64)

		for _, expense := range expenses {
			monthExpense += expense.Amount
			categoryMonthMap[expense.Category.Name] += expense.Amount
		}

		// Track highest income/expense months
		if monthIncome > highestIncome {
			highestIncome = monthIncome
			highestIncomeMonth = month
		}
		if monthExpense > highestExpense {
			highestExpense = monthExpense
			highestExpenseMonth = month
		}

		// Update category yearly stats
		for category, amount := range categoryMonthMap {
			entry := categoryYearlyMap[category]
			entry.total += amount
			entry.monthCount++
			if amount > entry.highestAmount {
				entry.highestAmount = amount
				entry.highestMonth = month
			}
			categoryYearlyMap[category] = entry
		}

		totalIncome += monthIncome
		totalExpenses += monthExpense

		monthlyBreakdown = append(monthlyBreakdown, MonthBreakdown{
			Month:    month,
			Income:   monthIncome,
			Expenses: monthExpense,
			Savings:  monthIncome - monthExpense,
		})
	}

	// Build category yearly summary
	categoryYearlySummary := make([]CategoryYearlySummary, 0)
	for category, data := range categoryYearlyMap {
		avgMonthly := data.total / float64(data.monthCount)
		categoryYearlySummary = append(categoryYearlySummary, CategoryYearlySummary{
			Category:       category,
			TotalSpent:     data.total,
			AverageMonthly: avgMonthly,
			HighestMonth:   data.highestMonth,
			HighestAmount:  data.highestAmount,
		})
	}

	// Calculate trends
	avgMonthlyIncome := totalIncome / 12
	avgMonthlyExpense := totalExpenses / 12
	totalSavings := totalIncome - totalExpenses
	savingsRate := 0.0
	if totalIncome > 0 {
		savingsRate = (totalSavings / totalIncome) * 100
	}

	return &YearlyReport{
		Year:                  year,
		TotalIncome:           totalIncome,
		TotalExpenses:         totalExpenses,
		TotalSavings:          totalSavings,
		MonthlyBreakdown:      monthlyBreakdown,
		CategoryYearlySummary: categoryYearlySummary,
		Trends: YearlyTrends{
			AverageMonthlyIncome:  avgMonthlyIncome,
			AverageMonthlyExpense: avgMonthlyExpense,
			HighestIncomeMonth:    highestIncomeMonth,
			HighestExpenseMonth:   highestExpenseMonth,
			TotalSavingsRate:      savingsRate,
		},
	}, nil
}
