package repositories

import (
	"finance-tracking-app/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ns *models.NotificationSetting) error
	FindByID(id uuid.UUID) (*models.NotificationSetting, error)
	FindAll() ([]models.NotificationSetting, error)
	FindDueNow(dayOfMonth int, currentTime string) ([]models.NotificationSetting, error)
	Update(ns *models.NotificationSetting) error
	UpdateLastSent(id uuid.UUID, sentAt time.Time) error
	Delete(id uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ns *models.NotificationSetting) error {
	return r.db.Create(ns).Error
}

func (r *notificationRepository) FindByID(id uuid.UUID) (*models.NotificationSetting, error) {
	var ns models.NotificationSetting
	err := r.db.First(&ns, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ns, nil
}

func (r *notificationRepository) FindAll() ([]models.NotificationSetting, error) {
	var settings []models.NotificationSetting
	err := r.db.Order("created_at ASC").Find(&settings).Error
	return settings, err
}

// FindDueNow finds enabled notifications matching today's day and current time window (HH:MM)
func (r *notificationRepository) FindDueNow(dayOfMonth int, currentTime string) ([]models.NotificationSetting, error) {
	var settings []models.NotificationSetting
	// day_of_month: 0 = every day, else specific day
	query := r.db.Where("is_enabled = ?", true).
		Where("time_of_day = ?", currentTime).
		Where("(day_of_month = 0 OR day_of_month = ?)", dayOfMonth)

	// Avoid sending twice in same minute: last_sent_at should not be in current minute
	now := time.Now()
	minuteStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	query = query.Where("last_sent_at IS NULL OR last_sent_at < ?", minuteStart)

	err := query.Find(&settings).Error
	return settings, err
}

func (r *notificationRepository) Update(ns *models.NotificationSetting) error {
	return r.db.Save(ns).Error
}

func (r *notificationRepository) UpdateLastSent(id uuid.UUID, sentAt time.Time) error {
	return r.db.Model(&models.NotificationSetting{}).
		Where("id = ?", id).
		Update("last_sent_at", sentAt).Error
}

func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.NotificationSetting{}, "id = ?", id).Error
}

// ensure fmt is used
var _ = fmt.Sprintf
