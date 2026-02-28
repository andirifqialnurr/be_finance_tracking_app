package repositories

import (
	"finance-tracking-app/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduledFundRepository interface {
	Create(sf *models.ScheduledFund) error
	FindByID(id uuid.UUID) (*models.ScheduledFund, error)
	FindAll(activeOnly bool) ([]models.ScheduledFund, error)
	FindDueToday(today time.Time) ([]models.ScheduledFund, error)
	Update(sf *models.ScheduledFund) error
	UpdateExecution(id uuid.UUID, lastExecuted time.Time, nextExecute time.Time) error
	Delete(id uuid.UUID) error
}

type scheduledFundRepository struct {
	db *gorm.DB
}

func NewScheduledFundRepository(db *gorm.DB) ScheduledFundRepository {
	return &scheduledFundRepository{db: db}
}

func (r *scheduledFundRepository) Create(sf *models.ScheduledFund) error {
	return r.db.Create(sf).Error
}

func (r *scheduledFundRepository) FindByID(id uuid.UUID) (*models.ScheduledFund, error) {
	var sf models.ScheduledFund
	err := r.db.Preload("Account").Preload("FromAccount").
		First(&sf, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sf, nil
}

func (r *scheduledFundRepository) FindAll(activeOnly bool) ([]models.ScheduledFund, error) {
	var sfs []models.ScheduledFund
	q := r.db.Preload("Account").Preload("FromAccount")
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	err := q.Order("day_of_month ASC").Find(&sfs).Error
	return sfs, err
}

// FindDueToday finds scheduled funds whose next_execute_at <= today midnight
func (r *scheduledFundRepository) FindDueToday(today time.Time) ([]models.ScheduledFund, error) {
	var sfs []models.ScheduledFund
	todayMidnight := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, time.UTC)
	err := r.db.Preload("Account").Preload("FromAccount").
		Where("is_active = ? AND next_execute_at < ?", true, todayMidnight).
		Find(&sfs).Error
	return sfs, err
}

func (r *scheduledFundRepository) Update(sf *models.ScheduledFund) error {
	return r.db.Save(sf).Error
}

func (r *scheduledFundRepository) UpdateExecution(id uuid.UUID, lastExecuted time.Time, nextExecute time.Time) error {
	return r.db.Model(&models.ScheduledFund{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_executed_at": lastExecuted,
			"next_execute_at":  nextExecute,
		}).Error
}

func (r *scheduledFundRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ScheduledFund{}, "id = ?", id).Error
}
