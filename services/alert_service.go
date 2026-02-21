package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"time"

	"github.com/google/uuid"
)

type AlertService interface {
	CreateAlert(alert *models.BudgetAlert) (*models.BudgetAlert, error)
	GetAlertByID(id uuid.UUID) (*models.BudgetAlert, error)
	GetAllAlerts(status string) ([]AlertWithStatus, error)
	UpdateAlert(id uuid.UUID, alert *models.BudgetAlert) (*models.BudgetAlert, error)
	DeleteAlert(id uuid.UUID) error
	CheckBudgetAlert(categoryID uuid.UUID, month, year int) (*AlertStatus, error)
}

type alertService struct {
	alertRepo    repositories.BudgetAlertRepository
	budgetRepo   repositories.CategoryBudgetRepository
	categoryRepo repositories.ExpenseCategoryRepository
}

func NewAlertService(
	alertRepo repositories.BudgetAlertRepository,
	budgetRepo repositories.CategoryBudgetRepository,
	categoryRepo repositories.ExpenseCategoryRepository,
) AlertService {
	return &alertService{
		alertRepo:    alertRepo,
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

type AlertWithStatus struct {
	ID           uuid.UUID               `json:"id"`
	Category     *models.ExpenseCategory `json:"category"`
	Threshold    int                     `json:"threshold"`
	CurrentUsage float64                 `json:"current_usage"`
	Status       string                  `json:"status"` // "triggered", "warning", "safe"
	Message      string                  `json:"message"`
	Level        string                  `json:"level"` // "critical", "warning", "info"
	TriggeredAt  *time.Time              `json:"triggered_at,omitempty"`
}

type AlertStatus struct {
	Triggered  bool    `json:"triggered"`
	Level      string  `json:"level"` // "warning", "critical"
	Message    string  `json:"message"`
	Percentage float64 `json:"percentage_used"`
}

func (s *alertService) CreateAlert(alert *models.BudgetAlert) (*models.BudgetAlert, error) {
	// Validate
	if alert.ThresholdPercentage < 0 || alert.ThresholdPercentage > 100 {
		return nil, errors.New("threshold percentage must be between 0 and 100")
	}

	// Check if category exists
	category, err := s.categoryRepo.FindByID(alert.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if !category.IsActive {
		return nil, errors.New("category is not active")
	}

	// Check if alert already exists for this category
	existingAlert, err := s.alertRepo.FindByCategoryID(alert.CategoryID)
	if err == nil && existingAlert != nil {
		return nil, errors.New("alert already exists for this category")
	}

	if err := s.alertRepo.Create(alert); err != nil {
		return nil, err
	}

	// Load category
	alert.Category = *category

	return alert, nil
}

func (s *alertService) GetAlertByID(id uuid.UUID) (*models.BudgetAlert, error) {
	return s.alertRepo.FindByID(id)
}

func (s *alertService) GetAllAlerts(status string) ([]AlertWithStatus, error) {
	var alerts []models.BudgetAlert
	var err error

	if status == "active" {
		alerts, err = s.alertRepo.FindAllActive()
	} else {
		alerts, err = s.alertRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	alertsWithStatus := make([]AlertWithStatus, 0, len(alerts))

	for _, alert := range alerts {
		// Get current budget for this category
		budget, err := s.budgetRepo.FindByCategoryAndMonth(alert.CategoryID, month, year)

		var currentUsage float64
		var alertStatus, message, level string

		if err == nil && budget.AllocatedAmount > 0 {
			currentUsage = (budget.SpentAmount / budget.AllocatedAmount) * 100

			if currentUsage >= 100 {
				alertStatus = "triggered"
				level = "critical"
				message = "Budget fully depleted"
			} else if currentUsage >= float64(alert.ThresholdPercentage) {
				alertStatus = "triggered"
				level = "warning"
				message = "Budget threshold exceeded"
			} else {
				alertStatus = "safe"
				level = "info"
				message = "Budget within limits"
			}
		} else {
			alertStatus = "safe"
			level = "info"
			message = "No budget data available"
			currentUsage = 0
		}

		alertsWithStatus = append(alertsWithStatus, AlertWithStatus{
			ID:           alert.ID,
			Category:     &alert.Category,
			Threshold:    alert.ThresholdPercentage,
			CurrentUsage: currentUsage,
			Status:       alertStatus,
			Message:      message,
			Level:        level,
			TriggeredAt:  alert.LastTriggered,
		})
	}

	return alertsWithStatus, nil
}

func (s *alertService) UpdateAlert(id uuid.UUID, alert *models.BudgetAlert) (*models.BudgetAlert, error) {
	// Validate
	if alert.ThresholdPercentage < 0 || alert.ThresholdPercentage > 100 {
		return nil, errors.New("threshold percentage must be between 0 and 100")
	}

	// Check if alert exists
	existingAlert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("alert not found")
	}

	// Update fields
	existingAlert.ThresholdPercentage = alert.ThresholdPercentage
	existingAlert.IsEnabled = alert.IsEnabled

	if err := s.alertRepo.Update(existingAlert); err != nil {
		return nil, err
	}

	return existingAlert, nil
}

func (s *alertService) DeleteAlert(id uuid.UUID) error {
	return s.alertRepo.Delete(id)
}

func (s *alertService) CheckBudgetAlert(categoryID uuid.UUID, month, year int) (*AlertStatus, error) {
	// Get alert for this category
	alert, err := s.alertRepo.FindByCategoryID(categoryID)
	if err != nil {
		// No alert configured, return safe status
		return &AlertStatus{
			Triggered:  false,
			Level:      "info",
			Message:    "No alert configured",
			Percentage: 0,
		}, nil
	}

	if !alert.IsEnabled {
		return &AlertStatus{
			Triggered:  false,
			Level:      "info",
			Message:    "Alert is disabled",
			Percentage: 0,
		}, nil
	}

	// Get budget
	budget, err := s.budgetRepo.FindByCategoryAndMonth(categoryID, month, year)
	if err != nil {
		return &AlertStatus{
			Triggered:  false,
			Level:      "info",
			Message:    "No budget data available",
			Percentage: 0,
		}, nil
	}

	if budget.AllocatedAmount == 0 {
		return &AlertStatus{
			Triggered:  false,
			Level:      "info",
			Message:    "Budget not allocated",
			Percentage: 0,
		}, nil
	}

	usagePercentage := (budget.SpentAmount / budget.AllocatedAmount) * 100

	if usagePercentage >= 100 {
		// Update last triggered
		now := time.Now()
		alert.LastTriggered = &now
		s.alertRepo.Update(alert)

		return &AlertStatus{
			Triggered:  true,
			Level:      "critical",
			Message:    "Budget fully depleted",
			Percentage: usagePercentage,
		}, nil
	}

	if usagePercentage >= float64(alert.ThresholdPercentage) {
		// Update last triggered
		now := time.Now()
		alert.LastTriggered = &now
		s.alertRepo.Update(alert)

		return &AlertStatus{
			Triggered:  true,
			Level:      "warning",
			Message:    "Budget threshold exceeded",
			Percentage: usagePercentage,
		}, nil
	}

	return &AlertStatus{
		Triggered:  false,
		Level:      "info",
		Message:    "Budget within limits",
		Percentage: usagePercentage,
	}, nil
}
