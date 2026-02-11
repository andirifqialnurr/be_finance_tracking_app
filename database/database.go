package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./finance_tracking.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Create tables
	createTables()
}

// createTables creates necessary tables if they don't exist
func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		amount REAL NOT NULL,
		type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
		category TEXT NOT NULL,
		description TEXT,
		date TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	log.Println("Tables created successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
