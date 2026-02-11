package handlers

import (
	"database/sql"
	"encoding/json"
	"finance-tracking-app/database"
	"finance-tracking-app/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllTransactions returns all transactions
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, title, amount, type, category, description, date, created_at, updated_at 
	          FROM transactions ORDER BY date DESC, created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.Title, &t.Amount, &t.Type, &t.Category, &t.Description, &t.Date, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse transactions")
			return
		}
		transactions = append(transactions, t)
	}

	respondWithSuccess(w, http.StatusOK, "Transactions fetched successfully", transactions)
}

// GetTransactionByID returns a single transaction by ID
func GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	query := `SELECT id, title, amount, type, category, description, date, created_at, updated_at 
	          FROM transactions WHERE id = ?`

	var t models.Transaction
	err = database.DB.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Amount, &t.Type, &t.Category, &t.Description, &t.Date, &t.CreatedAt, &t.UpdatedAt)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Transaction not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch transaction")
		return
	}

	respondWithSuccess(w, http.StatusOK, "Transaction fetched successfully", t)
}

// CreateTransaction creates a new transaction
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t models.Transaction
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if t.Title == "" || t.Amount == 0 || t.Type == "" || t.Category == "" || t.Date == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Validate type
	if t.Type != "income" && t.Type != "expense" {
		respondWithError(w, http.StatusBadRequest, "Type must be 'income' or 'expense'")
		return
	}

	query := `INSERT INTO transactions (title, amount, type, category, description, date) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query, t.Title, t.Amount, t.Type, t.Category, t.Description, t.Date)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	id, _ := result.LastInsertId()
	t.ID = int(id)

	respondWithSuccess(w, http.StatusCreated, "Transaction created successfully", t)
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var t models.Transaction
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if t.Title == "" || t.Amount == 0 || t.Type == "" || t.Category == "" || t.Date == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Validate type
	if t.Type != "income" && t.Type != "expense" {
		respondWithError(w, http.StatusBadRequest, "Type must be 'income' or 'expense'")
		return
	}

	query := `UPDATE transactions 
	          SET title = ?, amount = ?, type = ?, category = ?, description = ?, date = ?, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = ?`

	result, err := database.DB.Exec(query, t.Title, t.Amount, t.Type, t.Category, t.Description, t.Date, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update transaction")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	t.ID = id
	respondWithSuccess(w, http.StatusOK, "Transaction updated successfully", t)
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	query := `DELETE FROM transactions WHERE id = ?`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete transaction")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Transaction not found")
		return
	}

	respondWithSuccess(w, http.StatusOK, "Transaction deleted successfully", nil)
}

// GetStatistics returns financial statistics
func GetStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stats models.Statistics

	// Get total income
	incomeQuery := `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'income'`
	database.DB.QueryRow(incomeQuery).Scan(&stats.TotalIncome)

	// Get total expense
	expenseQuery := `SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'expense'`
	database.DB.QueryRow(expenseQuery).Scan(&stats.TotalExpense)

	// Calculate balance
	stats.Balance = stats.TotalIncome - stats.TotalExpense

	// Get total transactions count
	countQuery := `SELECT COUNT(*) FROM transactions`
	database.DB.QueryRow(countQuery).Scan(&stats.Transactions)

	respondWithSuccess(w, http.StatusOK, "Statistics fetched successfully", stats)
}

// Helper function to send error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	response := models.Response{
		Success: false,
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// Helper function to send success response
func respondWithSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	response := models.Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
