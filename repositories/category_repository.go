package repositories

import (
	"finance-tracking-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseCategoryRepository interface {
	Create(category *models.ExpenseCategory) error
	FindByID(id uuid.UUID) (*models.ExpenseCategory, error)
	FindAll() ([]models.ExpenseCategory, error)
	FindAllActive() ([]models.ExpenseCategory, error)
	Update(category *models.ExpenseCategory) error
	Delete(id uuid.UUID) error
	SoftDelete(id uuid.UUID) error
}

type expenseCategoryRepository struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) ExpenseCategoryRepository {
	return &expenseCategoryRepository{db: db}
}

func (r *expenseCategoryRepository) Create(category *models.ExpenseCategory) error {
	return r.db.Create(category).Error
}

func (r *expenseCategoryRepository) FindByID(id uuid.UUID) (*models.ExpenseCategory, error) {
	var category models.ExpenseCategory
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *expenseCategoryRepository) FindAll() ([]models.ExpenseCategory, error) {
	var categories []models.ExpenseCategory
	err := r.db.Order("allocation_priority ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *expenseCategoryRepository) FindAllActive() ([]models.ExpenseCategory, error) {
	var categories []models.ExpenseCategory
	err := r.db.Where("is_active = ?", true).
		Order("allocation_priority ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *expenseCategoryRepository) Update(category *models.ExpenseCategory) error {
	return r.db.Save(category).Error
}

func (r *expenseCategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExpenseCategory{}, "id = ?", id).Error
}

func (r *expenseCategoryRepository) SoftDelete(id uuid.UUID) error {
	return r.db.Model(&models.ExpenseCategory{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
