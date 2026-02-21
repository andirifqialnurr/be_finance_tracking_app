package repositories

import (
	"finance-tracking-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetReallocationRepository interface {
	Create(reallocation *models.BudgetReallocation) error
	FindByID(id uuid.UUID) (*models.BudgetReallocation, error)
	FindByMonthYear(month, year int) ([]models.BudgetReallocation, error)
	FindAll(limit, offset int) ([]models.BudgetReallocation, error)
	Delete(id uuid.UUID) error
}

type budgetReallocationRepository struct {
	db *gorm.DB
}

func NewBudgetReallocationRepository(db *gorm.DB) BudgetReallocationRepository {
	return &budgetReallocationRepository{db: db}
}

func (r *budgetReallocationRepository) Create(reallocation *models.BudgetReallocation) error {
	return r.db.Create(reallocation).Error
}

func (r *budgetReallocationRepository) FindByID(id uuid.UUID) (*models.BudgetReallocation, error) {
	var reallocation models.BudgetReallocation
	err := r.db.Preload("FromCategory").Preload("ToCategory").
		First(&reallocation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reallocation, nil
}

func (r *budgetReallocationRepository) FindByMonthYear(month, year int) ([]models.BudgetReallocation, error) {
	var reallocations []models.BudgetReallocation
	err := r.db.Preload("FromCategory").Preload("ToCategory").
		Where("month = ? AND year = ?", month, year).
		Order("created_at DESC").
		Find(&reallocations).Error
	return reallocations, err
}

func (r *budgetReallocationRepository) FindAll(limit, offset int) ([]models.BudgetReallocation, error) {
	var reallocations []models.BudgetReallocation
	err := r.db.Preload("FromCategory").Preload("ToCategory").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reallocations).Error
	return reallocations, err
}

func (r *budgetReallocationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BudgetReallocation{}, "id = ?", id).Error
}
