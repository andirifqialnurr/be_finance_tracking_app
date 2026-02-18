package services

import (
	"errors"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(category *models.ExpenseCategory) (*models.ExpenseCategory, error)
	GetCategoryByID(id uuid.UUID) (*models.ExpenseCategory, error)
	GetAllCategories() ([]models.ExpenseCategory, error)
	GetActiveCategories() ([]models.ExpenseCategory, error)
	UpdateCategory(id uuid.UUID, category *models.ExpenseCategory) (*models.ExpenseCategory, error)
	DeleteCategory(id uuid.UUID) error
	DeactivateCategory(id uuid.UUID) error
}

type categoryService struct {
	categoryRepo repositories.ExpenseCategoryRepository
}

func NewCategoryService(categoryRepo repositories.ExpenseCategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(category *models.ExpenseCategory) (*models.ExpenseCategory, error) {
	// Validate category type
	validTypes := map[string]bool{
		models.CategoryTypeSubscription:    true,
		models.CategoryTypeDailyContinuous: true,
		models.CategoryTypeUsageBased:      true,
		models.CategoryTypeOneTime:         true,
	}

	if !validTypes[category.Type] {
		return nil, errors.New("invalid category type")
	}

	// Validate monthly budget
	if category.MonthlyBudget <= 0 {
		return nil, errors.New("monthly budget must be greater than 0")
	}

	// Set default values
	if category.AllocationPriority == 0 {
		category.AllocationPriority = 1
	}
	category.IsActive = true

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategoryByID(id uuid.UUID) (*models.ExpenseCategory, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) GetAllCategories() ([]models.ExpenseCategory, error) {
	return s.categoryRepo.FindAll()
}

func (s *categoryService) GetActiveCategories() ([]models.ExpenseCategory, error) {
	return s.categoryRepo.FindAllActive()
}

func (s *categoryService) UpdateCategory(id uuid.UUID, updatedCategory *models.ExpenseCategory) (*models.ExpenseCategory, error) {
	// Check if category exists
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updatedCategory.Name != "" {
		category.Name = updatedCategory.Name
	}
	if updatedCategory.Type != "" {
		category.Type = updatedCategory.Type
	}
	if updatedCategory.MonthlyBudget > 0 {
		category.MonthlyBudget = updatedCategory.MonthlyBudget
	}
	if updatedCategory.AllocationPriority > 0 {
		category.AllocationPriority = updatedCategory.AllocationPriority
	}
	if updatedCategory.Metadata != nil {
		category.Metadata = updatedCategory.Metadata
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) DeactivateCategory(id uuid.UUID) error {
	return s.categoryRepo.SoftDelete(id)
}
