package repositories

import (
	"finance-tracking-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetAlertRepository interface {
	Create(alert *models.BudgetAlert) error
	FindByID(id uuid.UUID) (*models.BudgetAlert, error)
	FindByCategoryID(categoryID uuid.UUID) (*models.BudgetAlert, error)
	FindAll() ([]models.BudgetAlert, error)
	FindAllActive() ([]models.BudgetAlert, error)
	Update(alert *models.BudgetAlert) error
	Delete(id uuid.UUID) error
}

type budgetAlertRepository struct {
	db *gorm.DB
}

func NewBudgetAlertRepository(db *gorm.DB) BudgetAlertRepository {
	return &budgetAlertRepository{db: db}
}

func (r *budgetAlertRepository) Create(alert *models.BudgetAlert) error {
	return r.db.Create(alert).Error
}

func (r *budgetAlertRepository) FindByID(id uuid.UUID) (*models.BudgetAlert, error) {
	var alert models.BudgetAlert
	err := r.db.Preload("Category").First(&alert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *budgetAlertRepository) FindByCategoryID(categoryID uuid.UUID) (*models.BudgetAlert, error) {
	var alert models.BudgetAlert
	err := r.db.Preload("Category").First(&alert, "category_id = ?", categoryID).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *budgetAlertRepository) FindAll() ([]models.BudgetAlert, error) {
	var alerts []models.BudgetAlert
	err := r.db.Preload("Category").Find(&alerts).Error
	return alerts, err
}

func (r *budgetAlertRepository) FindAllActive() ([]models.BudgetAlert, error) {
	var alerts []models.BudgetAlert
	err := r.db.Preload("Category").Where("is_enabled = ?", true).Find(&alerts).Error
	return alerts, err
}

func (r *budgetAlertRepository) Update(alert *models.BudgetAlert) error {
	return r.db.Save(alert).Error
}

func (r *budgetAlertRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BudgetAlert{}, "id = ?", id).Error
}
