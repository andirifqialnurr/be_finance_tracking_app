package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"time"

	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetSpendingPattern(categoryID uuid.UUID, months int) (*SpendingPatternResponse, error)
	GetCategoryComparison(month, year int) (*CategoryComparisonResponse, error)
	GetTopSpending(month, year, limit int) ([]TopSpendingItem, error)
	GetBudgetPerformance(year int) (*BudgetPerformanceResponse, error)
}

type analyticsService struct {
	expenseRepo  repositories.ExpenseRepository
	budgetRepo   repositories.CategoryBudgetRepository
	categoryRepo repositories.ExpenseCategoryRepository
}

func NewAnalyticsService(
	expenseRepo repositories.ExpenseRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
) AnalyticsService {
	return &analyticsService{
		expenseRepo:  expenseRepo,
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

// Spending Pattern Response
type MonthlyPattern struct {
	Month      string  `json:"month"`
	Spent      float64 `json:"spent"`
	Budget     float64 `json:"budget"`
	Percentage float64 `json:"percentage"`
}

type SpendingPatternResponse struct {
	Category               *models.ExpenseCategory `json:"category"`
	Pattern                []MonthlyPattern        `json:"pattern"`
	AverageMonthlySpending float64                 `json:"average_monthly_spending"`
	Trend                  string                  `json:"trend"` // "increasing", "decreasing", "stable"
}

// Category Comparison Response
type CategorySpending struct {
	Name              string  `json:"name"`
	Spent             float64 `json:"spent"`
	PercentageOfTotal float64 `json:"percentage_of_total"`
}

type MonthSummary struct {
	Month      string             `json:"month"`
	Categories []CategorySpending `json:"categories"`
	TotalSpent float64            `json:"total_spent"`
}

type CategoryChange struct {
	Category         string  `json:"category"`
	ChangeAmount     float64 `json:"change_amount"`
	ChangePercentage float64 `json:"change_percentage"`
}

type CategoryComparisonResponse struct {
	CurrentMonth  MonthSummary     `json:"current_month"`
	PreviousMonth MonthSummary     `json:"previous_month"`
	Changes       []CategoryChange `json:"changes"`
}

// Top Spending Response
type TopSpendingItem struct {
	Category           string  `json:"category"`
	Amount             float64 `json:"amount"`
	Count              int     `json:"count"`
	AverageTransaction float64 `json:"avg_per_transaction"`
}

// Budget Performance Response
type MonthlyPerformance struct {
	Month                  int     `json:"month"`
	CategoriesOverBudget   int     `json:"categories_over_budget"`
	CategoriesUnderBudget  int     `json:"categories_under_budget"`
	AverageUsagePercentage float64 `json:"average_usage_percentage"`
}

type CategoryPerformance struct {
	Category         string  `json:"category"`
	TimesOverBudget  int     `json:"times_over_budget"`
	TimesUnderBudget int     `json:"times_under_budget"`
	AverageUsage     float64 `json:"average_usage"`
}

type BudgetPerformanceResponse struct {
	Year                int                   `json:"year"`
	MonthlyPerformance  []MonthlyPerformance  `json:"monthly_performance"`
	CategoryPerformance []CategoryPerformance `json:"category_performance"`
}

func (s *analyticsService) GetSpendingPattern(categoryID uuid.UUID, months int) (*SpendingPatternResponse, error) {
	var category *models.ExpenseCategory

	// If categoryID is provided (not nil), get specific category
	if categoryID != uuid.Nil {
		cat, err := s.categoryRepo.FindByID(categoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		category = cat
	}

	if months <= 0 {
		months = 6
	}

	now := time.Now()
	pattern := make([]MonthlyPattern, 0, months)
	var totalSpent float64

	for i := months - 1; i >= 0; i-- {
		targetDate := now.AddDate(0, -i, 0)
		month := int(targetDate.Month())
		year := targetDate.Year()

		var spent, budgetAmount float64

		if categoryID != uuid.Nil {
			// Get budget for specific category
			budget, err := s.budgetRepo.FindByCategoryAndMonth(categoryID, month, year)
			if err == nil {
				spent = budget.SpentAmount
				budgetAmount = budget.AllocatedAmount
			}
		} else {
			// Get all budgets for this month (aggregate)
			budgets, err := s.budgetRepo.FindByMonthYear(month, year)
			if err == nil {
				for _, budget := range budgets {
					spent += budget.SpentAmount
					budgetAmount += budget.AllocatedAmount
				}
			}
		}

		var percentage float64
		if budgetAmount > 0 {
			percentage = (spent / budgetAmount) * 100
		}

		pattern = append(pattern, MonthlyPattern{
			Month:      targetDate.Format("2006-01"),
			Spent:      spent,
			Budget:     budgetAmount,
			Percentage: percentage,
		})

		totalSpent += spent
	}

	avgSpending := totalSpent / float64(months)

	// Determine trend
	trend := "stable"
	if len(pattern) >= 2 {
		firstHalf := pattern[:len(pattern)/2]
		secondHalf := pattern[len(pattern)/2:]

		var firstHalfTotal, secondHalfTotal float64
		for _, p := range firstHalf {
			firstHalfTotal += p.Spent
		}
		for _, p := range secondHalf {
			secondHalfTotal += p.Spent
		}

		if secondHalfTotal > firstHalfTotal*1.1 {
			trend = "increasing"
		} else if secondHalfTotal < firstHalfTotal*0.9 {
			trend = "decreasing"
		}
	}

	return &SpendingPatternResponse{
		Category:               category,
		Pattern:                pattern,
		AverageMonthlySpending: avgSpending,
		Trend:                  trend,
	}, nil
}

func (s *analyticsService) GetCategoryComparison(month, year int) (*CategoryComparisonResponse, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	// Get previous month/year
	prevMonth := month - 1
	prevYear := year
	if prevMonth < 1 {
		prevMonth = 12
		prevYear--
	}

	// Get current month data
	currentBudgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	// Get previous month data
	previousBudgets, err := s.budgetRepo.FindByMonthYear(prevMonth, prevYear)
	if err != nil {
		previousBudgets = []models.CategoryBudget{}
	}

	// Build current month summary
	currentCategories := make([]CategorySpending, 0)
	var currentTotal float64
	for _, budget := range currentBudgets {
		currentTotal += budget.SpentAmount
		currentCategories = append(currentCategories, CategorySpending{
			Name:  budget.Category.Name,
			Spent: budget.SpentAmount,
		})
	}

	// Calculate percentages for current month
	for i := range currentCategories {
		if currentTotal > 0 {
			currentCategories[i].PercentageOfTotal = (currentCategories[i].Spent / currentTotal) * 100
		}
	}

	// Build previous month summary
	previousCategories := make([]CategorySpending, 0)
	var previousTotal float64
	for _, budget := range previousBudgets {
		previousTotal += budget.SpentAmount
		previousCategories = append(previousCategories, CategorySpending{
			Name:  budget.Category.Name,
			Spent: budget.SpentAmount,
		})
	}

	// Calculate percentages for previous month
	for i := range previousCategories {
		if previousTotal > 0 {
			previousCategories[i].PercentageOfTotal = (previousCategories[i].Spent / previousTotal) * 100
		}
	}

	// Calculate changes
	changes := make([]CategoryChange, 0)
	prevMap := make(map[string]float64)
	for _, cat := range previousCategories {
		prevMap[cat.Name] = cat.Spent
	}

	for _, curr := range currentCategories {
		prev := prevMap[curr.Name]
		changeAmount := curr.Spent - prev
		var changePercentage float64
		if prev > 0 {
			changePercentage = (changeAmount / prev) * 100
		}

		changes = append(changes, CategoryChange{
			Category:         curr.Name,
			ChangeAmount:     changeAmount,
			ChangePercentage: changePercentage,
		})
	}

	return &CategoryComparisonResponse{
		CurrentMonth: MonthSummary{
			Month:      time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01"),
			Categories: currentCategories,
			TotalSpent: currentTotal,
		},
		PreviousMonth: MonthSummary{
			Month:      time.Date(prevYear, time.Month(prevMonth), 1, 0, 0, 0, 0, time.UTC).Format("2006-01"),
			Categories: previousCategories,
			TotalSpent: previousTotal,
		},
		Changes: changes,
	}, nil
}

func (s *analyticsService) GetTopSpending(month, year, limit int) ([]TopSpendingItem, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}
	if limit <= 0 {
		limit = 5
	}

	// Get all budgets for this month
	_, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	// Get expenses for this month
	expenses, err := s.expenseRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	// Group expenses by category
	categoryExpenses := make(map[uuid.UUID]struct {
		name  string
		total float64
		count int
	})

	for _, expense := range expenses {
		entry := categoryExpenses[expense.CategoryID]
		entry.name = expense.Category.Name
		entry.total += expense.Amount
		entry.count++
		categoryExpenses[expense.CategoryID] = entry
	}

	// Build top spending list
	topSpending := make([]TopSpendingItem, 0)
	for _, data := range categoryExpenses {
		avg := data.total / float64(data.count)
		topSpending = append(topSpending, TopSpendingItem{
			Category:           data.name,
			Amount:             data.total,
			Count:              data.count,
			AverageTransaction: avg,
		})
	}

	// Sort by amount (descending)
	for i := 0; i < len(topSpending)-1; i++ {
		for j := i + 1; j < len(topSpending); j++ {
			if topSpending[j].Amount > topSpending[i].Amount {
				topSpending[i], topSpending[j] = topSpending[j], topSpending[i]
			}
		}
	}

	// Limit results
	if len(topSpending) > limit {
		topSpending = topSpending[:limit]
	}

	return topSpending, nil
}

func (s *analyticsService) GetBudgetPerformance(year int) (*BudgetPerformanceResponse, error) {
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	monthlyPerformance := make([]MonthlyPerformance, 0)
	categoryStats := make(map[string]*CategoryPerformance)

	// Iterate through all months
	for month := 1; month <= 12; month++ {
		budgets, err := s.budgetRepo.FindByMonthYear(month, year)
		if err != nil {
			continue
		}

		var overBudget, underBudget int
		var totalUsagePercentage float64
		validCategories := 0

		for _, budget := range budgets {
			if budget.AllocatedAmount == 0 {
				continue
			}

			usagePercentage := (budget.SpentAmount / budget.AllocatedAmount) * 100
			totalUsagePercentage += usagePercentage
			validCategories++

			// Track over/under budget
			if budget.SpentAmount > budget.AllocatedAmount {
				overBudget++
			} else {
				underBudget++
			}

			// Track category-level performance
			categoryName := budget.Category.Name
			if _, exists := categoryStats[categoryName]; !exists {
				categoryStats[categoryName] = &CategoryPerformance{
					Category: categoryName,
				}
			}

			if budget.SpentAmount > budget.AllocatedAmount {
				categoryStats[categoryName].TimesOverBudget++
			} else {
				categoryStats[categoryName].TimesUnderBudget++
			}
		}

		avgUsage := float64(0)
		if validCategories > 0 {
			avgUsage = totalUsagePercentage / float64(validCategories)
		}

		monthlyPerformance = append(monthlyPerformance, MonthlyPerformance{
			Month:                  month,
			CategoriesOverBudget:   overBudget,
			CategoriesUnderBudget:  underBudget,
			AverageUsagePercentage: avgUsage,
		})
	}

	// Calculate average usage for each category
	for categoryName, stats := range categoryStats {
		totalMonths := stats.TimesOverBudget + stats.TimesUnderBudget
		if totalMonths > 0 {
			// Recalculate average based on actual budgets
			var totalUsage float64
			for month := 1; month <= 12; month++ {
				budgets, _ := s.budgetRepo.FindByMonthYear(month, year)
				for _, budget := range budgets {
					if budget.Category.Name == categoryName && budget.AllocatedAmount > 0 {
						usagePercentage := (budget.SpentAmount / budget.AllocatedAmount) * 100
						totalUsage += usagePercentage
					}
				}
			}
			stats.AverageUsage = totalUsage / float64(totalMonths)
		}
	}

	// Convert map to slice
	categoryPerformance := make([]CategoryPerformance, 0, len(categoryStats))
	for _, perf := range categoryStats {
		categoryPerformance = append(categoryPerformance, *perf)
	}

	return &BudgetPerformanceResponse{
		Year:                year,
		MonthlyPerformance:  monthlyPerformance,
		CategoryPerformance: categoryPerformance,
	}, nil
}
