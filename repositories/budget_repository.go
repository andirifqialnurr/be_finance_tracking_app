package repositories

import (
	"finance-tracking-app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryBudgetRepository interface {
	Create(budget *models.CategoryBudget) error
	FindByID(id uuid.UUID) (*models.CategoryBudget, error)
	FindByCategoryAndMonth(categoryID uuid.UUID, month, year int) (*models.CategoryBudget, error)
	FindByMonthYear(month, year int) ([]models.CategoryBudget, error)
	Update(budget *models.CategoryBudget) error
	UpdateAllocatedAmount(id uuid.UUID, amount float64) error
	UpdateSpentAmount(id uuid.UUID, amount float64) error
	Delete(id uuid.UUID) error
}

type categoryBudgetRepository struct {
	db *gorm.DB
}

func NewCategoryBudgetRepository(db *gorm.DB) CategoryBudgetRepository {
	return &categoryBudgetRepository{db: db}
}

func (r *categoryBudgetRepository) Create(budget *models.CategoryBudget) error {
	return r.db.Create(budget).Error
}

func (r *categoryBudgetRepository) FindByID(id uuid.UUID) (*models.CategoryBudget, error) {
	var budget models.CategoryBudget
	err := r.db.Preload("Category").First(&budget, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *categoryBudgetRepository) FindByCategoryAndMonth(categoryID uuid.UUID, month, year int) (*models.CategoryBudget, error) {
	var budget models.CategoryBudget
	err := r.db.Preload("Category").
		Where("category_id = ? AND month = ? AND year = ?", categoryID, month, year).
		First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *categoryBudgetRepository) FindByMonthYear(month, year int) ([]models.CategoryBudget, error) {
	var budgets []models.CategoryBudget
	err := r.db.Preload("Category").
		Where("month = ? AND year = ?", month, year).
		Find(&budgets).Error
	return budgets, err
}

func (r *categoryBudgetRepository) Update(budget *models.CategoryBudget) error {
	return r.db.Save(budget).Error
}

func (r *categoryBudgetRepository) UpdateAllocatedAmount(id uuid.UUID, amount float64) error {
	return r.db.Model(&models.CategoryBudget{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"allocated_amount": gorm.Expr("allocated_amount + ?", amount),
			"remaining_amount": gorm.Expr("remaining_amount + ?", amount),
		}).Error
}

func (r *categoryBudgetRepository) UpdateSpentAmount(id uuid.UUID, amount float64) error {
	return r.db.Model(&models.CategoryBudget{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"spent_amount":     gorm.Expr("spent_amount + ?", amount),
			"remaining_amount": gorm.Expr("remaining_amount - ?", amount),
		}).Error
}

func (r *categoryBudgetRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CategoryBudget{}, "id = ?", id).Error
}
