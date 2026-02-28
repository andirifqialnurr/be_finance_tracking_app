package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Income represents money received from various sources
type Income struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null" json:"account_id"` // Source account
	Source      string    `gorm:"size:100;not null" json:"source"`
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Date        time.Time `gorm:"not null" json:"date"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Account           Account            `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	BudgetAllocations []BudgetAllocation `gorm:"foreignKey:IncomeID" json:"allocations,omitempty"`
}

// BeforeCreate hook to generate UUID
func (i *Income) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// ExpenseCategory represents a category of expenses (user-defined)
type ExpenseCategory struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Name               string         `gorm:"size:100;not null" json:"name"`
	Type               string         `gorm:"size:50;not null" json:"type"` // SUBSCRIPTION, DAILY_CONTINUOUS, USAGE_BASED, ONE_TIME
	MonthlyBudget      float64        `gorm:"type:decimal(15,2);not null;default:0" json:"monthly_budget"`
	DailyAmount        *float64       `gorm:"type:decimal(15,2)" json:"daily_amount,omitempty"` // Only for DAILY_CONTINUOUS
	AllocationPriority int            `gorm:"not null;default:1" json:"allocation_priority"`
	IsActive           bool           `gorm:"not null;default:true" json:"is_active"`
	Metadata           datatypes.JSON `gorm:"type:jsonb" json:"metadata" swaggertype:"object"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	// Relations
	CategoryBudgets   []CategoryBudget   `gorm:"foreignKey:CategoryID" json:"budgets,omitempty"`
	Expenses          []Expense          `gorm:"foreignKey:CategoryID" json:"expenses,omitempty"`
	BudgetAllocations []BudgetAllocation `gorm:"foreignKey:CategoryID" json:"allocations,omitempty"`
}

// BeforeCreate hook to generate UUID
func (ec *ExpenseCategory) BeforeCreate(tx *gorm.DB) error {
	if ec.ID == uuid.Nil {
		ec.ID = uuid.New()
	}
	return nil
}

// CategoryBudget represents real-time budget per category (reset per month)
type CategoryBudget struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CategoryID           uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Month                int       `gorm:"not null" json:"month"` // 1-12
	Year                 int       `gorm:"not null" json:"year"`
	AllocatedAmount      float64   `gorm:"type:decimal(15,2);not null;default:0" json:"allocated_amount"`
	SpentAmount          float64   `gorm:"type:decimal(15,2);not null;default:0" json:"spent_amount"`
	RemainingAmount      float64   `gorm:"type:decimal(15,2);not null;default:0" json:"remaining_amount"`
	EffectiveDailyAmount *float64  `gorm:"type:decimal(15,2)" json:"effective_daily_amount,omitempty"` // Snapshot for DAILY_CONTINUOUS
	DaysInMonth          *int      `json:"days_in_month,omitempty"`                                    // Snapshot of days in month
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relations
	Category ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hook to generate UUID
func (cb *CategoryBudget) BeforeCreate(tx *gorm.DB) error {
	if cb.ID == uuid.Nil {
		cb.ID = uuid.New()
	}
	return nil
}

// Expense represents individual expense transactions
type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null" json:"account_id"` // Account used for payment
	CategoryID  uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Date        time.Time `gorm:"not null" json:"date"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Account  Account         `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Category ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hook to generate UUID
func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// BudgetAllocation represents history of budget allocation from income to categories
type BudgetAllocation struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	IncomeID        uuid.UUID `gorm:"type:uuid;not null" json:"income_id"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	AllocatedAmount float64   `gorm:"type:decimal(15,2);not null" json:"allocated_amount"`
	Month           int       `gorm:"not null" json:"month"`
	Year            int       `gorm:"not null" json:"year"`
	CreatedAt       time.Time `json:"created_at"`

	// Relations
	Income   Income          `gorm:"foreignKey:IncomeID" json:"income,omitempty"`
	Category ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hook to generate UUID
func (ba *BudgetAllocation) BeforeCreate(tx *gorm.DB) error {
	if ba.ID == uuid.Nil {
		ba.ID = uuid.New()
	}
	return nil
}

// BudgetReallocation represents manual budget reallocation between categories
type BudgetReallocation struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	FromCategoryID uuid.UUID `gorm:"type:uuid;not null" json:"from_category_id"`
	ToCategoryID   uuid.UUID `gorm:"type:uuid;not null" json:"to_category_id"`
	Amount         float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Reason         string    `gorm:"type:text" json:"reason"`
	Month          int       `gorm:"not null" json:"month"`
	Year           int       `gorm:"not null" json:"year"`
	CreatedAt      time.Time `json:"created_at"`

	// Relations
	FromCategory ExpenseCategory `gorm:"foreignKey:FromCategoryID" json:"from_category,omitempty"`
	ToCategory   ExpenseCategory `gorm:"foreignKey:ToCategoryID" json:"to_category,omitempty"`
}

// BeforeCreate hook to generate UUID
func (br *BudgetReallocation) BeforeCreate(tx *gorm.DB) error {
	if br.ID == uuid.Nil {
		br.ID = uuid.New()
	}
	return nil
}

// BudgetAlert represents alert configuration per category
type BudgetAlert struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	CategoryID          uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"category_id"`
	ThresholdPercentage int        `gorm:"not null;default:80" json:"threshold_percentage"`
	IsEnabled           bool       `gorm:"not null;default:true" json:"is_enabled"`
	LastTriggered       *time.Time `json:"last_triggered,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Relations
	Category ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate hook to generate UUID
func (ba *BudgetAlert) BeforeCreate(tx *gorm.DB) error {
	if ba.ID == uuid.Nil {
		ba.ID = uuid.New()
	}
	return nil
}

// Response represents standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// CategoryType constants
const (
	CategoryTypeSubscription    = "SUBSCRIPTION"
	CategoryTypeDailyContinuous = "DAILY_CONTINUOUS"
	CategoryTypeUsageBased      = "USAGE_BASED"
	CategoryTypeOneTime         = "ONE_TIME"
)
