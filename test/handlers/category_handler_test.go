package handlers

import (
	"finance-tracking-app/handlers"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"finance-tracking-app/services"
	"finance-tracking-app/test/helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupCategoryHandler(t *testing.T) (*gin.Engine, *handlers.CategoryHandler) {
	helpers.SetupGinTestMode()
	db := helpers.SetupTestDB(t)

	categoryRepo := repositories.NewExpenseCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	router := gin.New()
	return router, categoryHandler
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	router, handler := setupCategoryHandler(t)
	router.POST("/categories", handler.CreateCategory)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "Valid category creation",
			requestBody: map[string]interface{}{
				"name":                "Test Category",
				"type":                "SUBSCRIPTION",
				"monthly_budget":      100000,
				"allocation_priority": 1,
				"metadata": map[string]interface{}{
					"test": "data",
				},
			},
			expectedStatus: http.StatusCreated,
			expectSuccess:  true,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"name": "Test Category",
			},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name:           "Empty request body",
			requestBody:    nil,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := helpers.MakeJSONRequest(t, "POST", "/categories", tt.requestBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.requestBody != nil {
				var response models.Response
				helpers.ParseJSONResponse(t, w, &response)
				assert.Equal(t, tt.expectSuccess, response.Success)
			}
		})
	}
}

func TestCategoryHandler_GetCategories(t *testing.T) {
	router, handler := setupCategoryHandler(t)

	// Seed some data
	db := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	categories := []*models.ExpenseCategory{
		{Name: "Cat 1", Type: models.CategoryTypeSubscription, MonthlyBudget: 100000, AllocationPriority: 1, IsActive: true},
		{Name: "Cat 2", Type: models.CategoryTypeUsageBased, MonthlyBudget: 50000, AllocationPriority: 2, IsActive: true},
	}

	for _, cat := range categories {
		categoryRepo.Create(cat)
	}

	router.GET("/categories", handler.GetCategories)

	req := helpers.MakeJSONRequest(t, "GET", "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	helpers.ParseJSONResponse(t, w, &response)
	assert.True(t, response.Success)
}

func TestCategoryHandler_GetCategoryByID(t *testing.T) {
	router, handler := setupCategoryHandler(t)

	// Seed a category
	db := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	category := &models.ExpenseCategory{
		Name:               "Test Category",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	router.GET("/categories/:id", handler.GetCategoryByID)

	// Test valid ID
	t.Run("Valid ID", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/categories/"+category.ID.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Response
		helpers.ParseJSONResponse(t, w, &response)
		assert.True(t, response.Success)
	})

	// Test invalid ID
	t.Run("Invalid ID", func(t *testing.T) {
		req := helpers.MakeJSONRequest(t, "GET", "/categories/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	router, handler := setupCategoryHandler(t)

	// Seed a category
	db := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	category := &models.ExpenseCategory{
		Name:               "Original Name",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	router.PUT("/categories/:id", handler.UpdateCategory)

	updateData := map[string]interface{}{
		"name":           "Updated Name",
		"monthly_budget": 200000,
	}

	req := helpers.MakeJSONRequest(t, "PUT", "/categories/"+category.ID.String(), updateData)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	helpers.ParseJSONResponse(t, w, &response)
	assert.True(t, response.Success)
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	router, handler := setupCategoryHandler(t)

	// Seed a category
	db := helpers.SetupTestDB(t)
	categoryRepo := repositories.NewExpenseCategoryRepository(db)

	category := &models.ExpenseCategory{
		Name:               "To Delete",
		Type:               models.CategoryTypeSubscription,
		MonthlyBudget:      100000,
		AllocationPriority: 1,
		IsActive:           true,
	}
	categoryRepo.Create(category)

	router.DELETE("/categories/:id", handler.DeleteCategory)

	req := helpers.MakeJSONRequest(t, "DELETE", "/categories/"+category.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	helpers.ParseJSONResponse(t, w, &response)
	assert.True(t, response.Success)

	// Verify it was deleted
	found, err := categoryRepo.FindByID(category.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
