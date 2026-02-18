package repositories

import (
	"finance-tracking-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	FindByID(id uuid.UUID) (*models.Expense, error)
	FindAll(limit, offset int) ([]models.Expense, error)
	FindByCategory(categoryID uuid.UUID, limit, offset int) ([]models.Expense, error)
	FindByDateRange(start, end time.Time) ([]models.Expense, error)
	FindByMonthYear(month, year int) ([]models.Expense, error)
	FindByCategoryAndMonth(categoryID uuid.UUID, month, year int) ([]models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) FindByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("Category").First(&expense, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) FindAll(limit, offset int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) FindByCategory(categoryID uuid.UUID, limit, offset int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("category_id = ?", categoryID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) FindByDateRange(start, end time.Time) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("date BETWEEN ? AND ?", start, end).
		Order("date DESC").
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) FindByMonthYear(month, year int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?", month, year).
		Order("date DESC").
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) FindByCategoryAndMonth(categoryID uuid.UUID, month, year int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Preload("Category").
		Where("category_id = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			categoryID, month, year).
		Order("date DESC").
		Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Expense{}, "id = ?", id).Error
}
