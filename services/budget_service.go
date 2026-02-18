package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
)

type BudgetService interface {
	GetBudgetsByMonthYear(month, year int) ([]BudgetSummary, error)
	GetBudgetSummary(month, year int) (*OverallSummary, error)
}

type BudgetSummary struct {
	Category        models.ExpenseCategory `json:"category"`
	AllocatedAmount float64                `json:"allocated_amount"`
	SpentAmount     float64                `json:"spent_amount"`
	RemainingAmount float64                `json:"remaining_amount"`
	PercentageUsed  float64                `json:"percentage_used"`
}

type OverallSummary struct {
	TotalIncome      float64 `json:"total_income"`
	TotalAllocated   float64 `json:"total_allocated"`
	TotalSpent       float64 `json:"total_spent"`
	TotalRemaining   float64 `json:"total_remaining"`
	UnallocatedFunds float64 `json:"unallocated_funds"`
}

type budgetService struct {
	budgetRepo     repositories.CategoryBudgetRepository
	incomeRepo     repositories.IncomeRepository
	allocationRepo repositories.BudgetAllocationRepository
}

func NewBudgetService(
	budgetRepo repositories.CategoryBudgetRepository,
	incomeRepo repositories.IncomeRepository,
	allocationRepo repositories.BudgetAllocationRepository,
) BudgetService {
	return &budgetService{
		budgetRepo:     budgetRepo,
		incomeRepo:     incomeRepo,
		allocationRepo: allocationRepo,
	}
}

func (s *budgetService) GetBudgetsByMonthYear(month, year int) ([]BudgetSummary, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	budgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	summaries := make([]BudgetSummary, 0, len(budgets))
	for _, budget := range budgets {
		percentageUsed := 0.0
		if budget.AllocatedAmount > 0 {
			percentageUsed = (budget.SpentAmount / budget.AllocatedAmount) * 100
		}

		summaries = append(summaries, BudgetSummary{
			Category:        budget.Category,
			AllocatedAmount: budget.AllocatedAmount,
			SpentAmount:     budget.SpentAmount,
			RemainingAmount: budget.RemainingAmount,
			PercentageUsed:  percentageUsed,
		})
	}

	return summaries, nil
}

func (s *budgetService) GetBudgetSummary(month, year int) (*OverallSummary, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2000 {
		return nil, errors.New("invalid year")
	}

	// Get all incomes for the month
	incomes, err := s.incomeRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalIncome float64
	for _, income := range incomes {
		totalIncome += income.Amount
	}

	// Get all allocations for the month
	allocations, err := s.allocationRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalAllocated float64
	for _, allocation := range allocations {
		totalAllocated += allocation.AllocatedAmount
	}

	// Get all budgets for the month
	budgets, err := s.budgetRepo.FindByMonthYear(month, year)
	if err != nil {
		return nil, err
	}

	var totalSpent, totalRemaining float64
	for _, budget := range budgets {
		totalSpent += budget.SpentAmount
		totalRemaining += budget.RemainingAmount
	}

	summary := &OverallSummary{
		TotalIncome:      totalIncome,
		TotalAllocated:   totalAllocated,
		TotalSpent:       totalSpent,
		TotalRemaining:   totalRemaining,
		UnallocatedFunds: totalIncome - totalAllocated,
	}

	return summary, nil
}
