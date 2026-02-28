package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"finance-tracking-app/config"
	"finance-tracking-app/models"
	"finance-tracking-app/repositories"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateSetting(ns *models.NotificationSetting) (*models.NotificationSetting, error)
	GetAllSettings() ([]models.NotificationSetting, error)
	GetSettingByID(id uuid.UUID) (*models.NotificationSetting, error)
	UpdateSetting(id uuid.UUID, ns *models.NotificationSetting) (*models.NotificationSetting, error)
	DeleteSetting(id uuid.UUID) error
	SendNotification(playerID, title, body string) error
	TestNotification(playerID, title, body string) error
	CheckAndSendDue() error
}

type notificationService struct {
	repo   repositories.NotificationRepository
	client *http.Client
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &notificationService{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *notificationService) CreateSetting(ns *models.NotificationSetting) (*models.NotificationSetting, error) {
	validTypes := map[string]bool{
		models.NotificationTypeAllocationReminder: true,
		models.NotificationTypeBudgetAlert:        true,
		models.NotificationTypeSavingsGoal:        true,
		models.NotificationTypeScheduledFund:      true,
	}
	if !validTypes[ns.Type] {
		return nil, errors.New("invalid notification type")
	}
	if ns.OneSignalPlayerID == "" {
		return nil, errors.New("onesignal_player_id is required")
	}
	if err := validateTimeOfDay(ns.TimeOfDay); err != nil {
		return nil, err
	}
	ns.IsEnabled = true
	return ns, s.repo.Create(ns)
}

func (s *notificationService) GetAllSettings() ([]models.NotificationSetting, error) {
	return s.repo.FindAll()
}

func (s *notificationService) GetSettingByID(id uuid.UUID) (*models.NotificationSetting, error) {
	return s.repo.FindByID(id)
}

func (s *notificationService) UpdateSetting(id uuid.UUID, updates *models.NotificationSetting) (*models.NotificationSetting, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Body != "" {
		existing.Body = updates.Body
	}
	if updates.TimeOfDay != "" {
		if err := validateTimeOfDay(updates.TimeOfDay); err != nil {
			return nil, err
		}
		existing.TimeOfDay = updates.TimeOfDay
	}
	if updates.DayOfMonth >= 0 {
		existing.DayOfMonth = updates.DayOfMonth
	}
	if updates.OneSignalPlayerID != "" {
		existing.OneSignalPlayerID = updates.OneSignalPlayerID
	}
	existing.IsEnabled = updates.IsEnabled
	return existing, s.repo.Update(existing)
}

func (s *notificationService) DeleteSetting(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *notificationService) SendNotification(playerID, title, body string) error {
	cfg := config.AppConfig
	if cfg.OneSignalAppID == "" || cfg.OneSignalAPIKey == "" {
		return errors.New("OneSignal is not configured")
	}

	payload := models.OneSignalNotification{
		AppID:            cfg.OneSignalAppID,
		IncludePlayerIDs: []string{playerID},
		Headings:         map[string]string{"en": title},
		Contents:         map[string]string{"en": body},
	}

	return s.sendToOneSignal(payload)
}

func (s *notificationService) TestNotification(playerID, title, body string) error {
	return s.SendNotification(playerID, title, body)
}

func (s *notificationService) CheckAndSendDue() error {
	now := time.Now()
	currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	dayOfMonth := now.Day()

	settings, err := s.repo.FindDueNow(dayOfMonth, currentTime)
	if err != nil {
		return err
	}

	for _, setting := range settings {
		err := s.SendNotification(setting.OneSignalPlayerID, setting.Title, setting.Body)
		if err == nil {
			_ = s.repo.UpdateLastSent(setting.ID, now)
		}
	}
	return nil
}

func (s *notificationService) sendToOneSignal(payload models.OneSignalNotification) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	cfg := config.AppConfig
	url := cfg.OneSignalAPIURL + "/notifications"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+cfg.OneSignalAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("OneSignal request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OneSignal returned status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func validateTimeOfDay(t string) error {
	if len(t) != 5 || t[2] != ':' {
		return errors.New("time_of_day must be in HH:MM format")
	}
	return nil
}
