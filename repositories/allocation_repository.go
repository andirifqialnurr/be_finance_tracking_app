package repositories

import (
	"finance-tracking-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetAllocationRepository interface {
	Create(allocation *models.BudgetAllocation) error
	FindByID(id uuid.UUID) (*models.BudgetAllocation, error)
	FindByIncome(incomeID uuid.UUID) ([]models.BudgetAllocation, error)
	FindByCategory(categoryID uuid.UUID) ([]models.BudgetAllocation, error)
	FindByMonthYear(month, year int) ([]models.BudgetAllocation, error)
	Delete(id uuid.UUID) error
}

type budgetAllocationRepository struct {
	db *gorm.DB
}

func NewBudgetAllocationRepository(db *gorm.DB) BudgetAllocationRepository {
	return &budgetAllocationRepository{db: db}
}

func (r *budgetAllocationRepository) Create(allocation *models.BudgetAllocation) error {
	return r.db.Create(allocation).Error
}

func (r *budgetAllocationRepository) FindByID(id uuid.UUID) (*models.BudgetAllocation, error) {
	var allocation models.BudgetAllocation
	err := r.db.Preload("Income").Preload("Category").
		First(&allocation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

func (r *budgetAllocationRepository) FindByIncome(incomeID uuid.UUID) ([]models.BudgetAllocation, error) {
	var allocations []models.BudgetAllocation
	err := r.db.Preload("Category").
		Where("income_id = ?", incomeID).
		Find(&allocations).Error
	return allocations, err
}

func (r *budgetAllocationRepository) FindByCategory(categoryID uuid.UUID) ([]models.BudgetAllocation, error) {
	var allocations []models.BudgetAllocation
	err := r.db.Preload("Income").
		Where("category_id = ?", categoryID).
		Find(&allocations).Error
	return allocations, err
}

func (r *budgetAllocationRepository) FindByMonthYear(month, year int) ([]models.BudgetAllocation, error) {
	var allocations []models.BudgetAllocation
	err := r.db.Preload("Income").Preload("Category").
		Where("month = ? AND year = ?", month, year).
		Find(&allocations).Error
	return allocations, err
}

func (r *budgetAllocationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BudgetAllocation{}, "id = ?", id).Error
}
