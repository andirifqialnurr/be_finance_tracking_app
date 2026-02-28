package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduledFund type constants
const (
	ScheduledFundTypeTopUp    = "TOP_UP"
	ScheduledFundTypeTransfer = "TRANSFER"
)

// NotificationSetting type constants
const (
	NotificationTypeAllocationReminder = "ALLOCATION_REMINDER"
	NotificationTypeBudgetAlert        = "BUDGET_ALERT"
	NotificationTypeSavingsGoal        = "SAVINGS_GOAL"
	NotificationTypeScheduledFund      = "SCHEDULED_FUND"
)

// ScheduledFund represents a recurring scheduled top-up or transfer
type ScheduledFund struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	AccountID      uuid.UUID  `gorm:"type:uuid;not null" json:"account_id"`
	FromAccountID  *uuid.UUID `gorm:"type:uuid" json:"from_account_id,omitempty"` // For TRANSFER type
	ScheduleType   string     `gorm:"size:20;not null" json:"schedule_type"`      // TOP_UP | TRANSFER
	Amount         float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	DayOfMonth     int        `gorm:"not null" json:"day_of_month"` // 1-31; 0 = last day of month
	Description    string     `gorm:"type:text" json:"description"`
	IsActive       bool       `gorm:"not null;default:true" json:"is_active"`
	LastExecutedAt *time.Time `json:"last_executed_at,omitempty"`
	NextExecuteAt  time.Time  `gorm:"not null" json:"next_execute_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations
	Account     Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	FromAccount *Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
}

// BeforeCreate hook to generate UUID
func (sf *ScheduledFund) BeforeCreate(tx *gorm.DB) error {
	if sf.ID == uuid.Nil {
		sf.ID = uuid.New()
	}
	return nil
}

// NotificationSetting represents a push notification configuration via OneSignal
type NotificationSetting struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Type              string     `gorm:"size:50;not null" json:"type"` // ALLOCATION_REMINDER | BUDGET_ALERT | SAVINGS_GOAL | SCHEDULED_FUND
	Title             string     `gorm:"size:200;not null" json:"title"`
	Body              string     `gorm:"type:text;not null" json:"body"`
	DayOfMonth        int        `gorm:"not null" json:"day_of_month"`       // 1-31; 0 = last day of month
	TimeOfDay         string     `gorm:"size:5;not null" json:"time_of_day"` // "HH:MM" format e.g. "09:00"
	IsEnabled         bool       `gorm:"not null;default:true" json:"is_enabled"`
	OneSignalPlayerID string     `gorm:"size:200;not null" json:"onesignal_player_id"`
	LastSentAt        *time.Time `json:"last_sent_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// BeforeCreate hook to generate UUID
func (ns *NotificationSetting) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == uuid.Nil {
		ns.ID = uuid.New()
	}
	return nil
}

// OneSignalNotification is used for sending push notifications
type OneSignalNotification struct {
	AppID            string            `json:"app_id"`
	IncludePlayerIDs []string          `json:"include_player_ids"`
	Headings         map[string]string `json:"headings"`
	Contents         map[string]string `json:"contents"`
}

// DeviceRegistration is used for registering a device's OneSignal player ID
type DeviceRegistration struct {
	OneSignalPlayerID string `json:"onesignal_player_id" binding:"required"`
}
