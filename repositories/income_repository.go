package repositories

import (
	"finance-tracking-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncomeRepository interface {
	Create(income *models.Income) error
	FindByID(id uuid.UUID) (*models.Income, error)
	FindAll(limit, offset int) ([]models.Income, error)
	FindByDateRange(start, end time.Time) ([]models.Income, error)
	FindByMonthYear(month, year int) ([]models.Income, error)
	Update(income *models.Income) error
	Delete(id uuid.UUID) error
}

type incomeRepository struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) IncomeRepository {
	return &incomeRepository{db: db}
}

func (r *incomeRepository) Create(income *models.Income) error {
	return r.db.Create(income).Error
}

func (r *incomeRepository) FindByID(id uuid.UUID) (*models.Income, error) {
	var income models.Income
	err := r.db.Preload("BudgetAllocations").First(&income, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &income, nil
}

func (r *incomeRepository) FindAll(limit, offset int) ([]models.Income, error) {
	var incomes []models.Income
	err := r.db.Order("date DESC").Limit(limit).Offset(offset).Find(&incomes).Error
	return incomes, err
}

func (r *incomeRepository) FindByDateRange(start, end time.Time) ([]models.Income, error) {
	var incomes []models.Income
	err := r.db.Where("date BETWEEN ? AND ?", start, end).
		Order("date DESC").
		Find(&incomes).Error
	return incomes, err
}

func (r *incomeRepository) FindByMonthYear(month, year int) ([]models.Income, error) {
	var incomes []models.Income
	err := r.db.Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Order("date DESC").
		Find(&incomes).Error
	return incomes, err
}

func (r *incomeRepository) Update(income *models.Income) error {
	return r.db.Save(income).Error
}

func (r *incomeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Income{}, "id = ?", id).Error
}
