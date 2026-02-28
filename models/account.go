package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account type constants
const (
	AccountTypeCard    = "CARD"
	AccountTypeCash    = "CASH"
	AccountTypeSavings = "SAVINGS"
)

// Card income type constants
const (
	CardIncomeTypeSalary    = "SALARY"
	CardIncomeTypeProject   = "PROJECT"
	CardIncomeTypeFreelance = "FREELANCE"
	CardIncomeTypeBusiness  = "BUSINESS"
	CardIncomeTypeOther     = "OTHER"
)

// Account represents a financial account (Card, Cash, or Savings)
type Account struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Type        string    `gorm:"size:20;not null" json:"type"`         // CARD | CASH | SAVINGS
	IncomeType  *string   `gorm:"size:30" json:"income_type,omitempty"` // Only for CARD: SALARY|PROJECT|FREELANCE|BUSINESS|OTHER
	Balance     float64   `gorm:"type:decimal(15,2);not null;default:0" json:"balance"`
	Color       string    `gorm:"size:10;default:'#4CAF50'" json:"color"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`

	// SAVINGS only
	GoalAmount *float64 `gorm:"type:decimal(15,2)" json:"goal_amount,omitempty"`
	GoalLabel  *string  `gorm:"size:200" json:"goal_label,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Incomes   []Income          `gorm:"foreignKey:AccountID" json:"incomes,omitempty"`
	Expenses  []Expense         `gorm:"foreignKey:AccountID" json:"expenses,omitempty"`
	Transfers []AccountTransfer `gorm:"foreignKey:FromAccountID" json:"transfers,omitempty"`
}

// BeforeCreate hook to generate UUID
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AccountTransfer represents a fund transfer between two accounts
type AccountTransfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FromAccountID uuid.UUID `gorm:"type:uuid;not null" json:"from_account_id"`
	ToAccountID   uuid.UUID `gorm:"type:uuid;not null" json:"to_account_id"`
	Amount        float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Note          string    `gorm:"type:text" json:"note"`
	TransferDate  time.Time `gorm:"not null" json:"transfer_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	FromAccount Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
}

// BeforeCreate hook to generate UUID
func (at *AccountTransfer) BeforeCreate(tx *gorm.DB) error {
	if at.ID == uuid.Nil {
		at.ID = uuid.New()
	}
	return nil
}

// AccountBalanceSummary is a non-persisted struct for API responses
type AccountBalanceSummary struct {
	CurrentBalance       float64 `json:"current_balance"`
	TotalIncomeThisMonth float64 `json:"total_income_this_month"`
	TotalSpentThisMonth  float64 `json:"total_spent_this_month"`
	TotalTransferredOut  float64 `json:"total_transferred_out"`
	TotalTransferredIn   float64 `json:"total_transferred_in"`
}

// SavingsProgress is a non-persisted struct for API responses
type SavingsProgress struct {
	GoalAmount         float64 `json:"goal_amount"`
	GoalLabel          string  `json:"goal_label"`
	ProgressPercentage float64 `json:"progress_percentage"`
	RemainingToGoal    float64 `json:"remaining_to_goal"`
}
