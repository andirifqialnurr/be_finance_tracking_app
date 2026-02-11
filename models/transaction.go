package models

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // "income" or "expense"
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Statistics represents financial statistics
type Statistics struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	Transactions int     `json:"transactions"`
}

// Response represents API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
