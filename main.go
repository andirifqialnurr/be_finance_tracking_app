package main

import (
	"finance-tracking-app/database"
	"finance-tracking-app/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Create router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Transaction routes
	api.HandleFunc("/transactions", handlers.GetAllTransactions).Methods("GET")
	api.HandleFunc("/transactions/{id}", handlers.GetTransactionByID).Methods("GET")
	api.HandleFunc("/transactions", handlers.CreateTransaction).Methods("POST")
	api.HandleFunc("/transactions/{id}", handlers.UpdateTransaction).Methods("PUT")
	api.HandleFunc("/transactions/{id}", handlers.DeleteTransaction).Methods("DELETE")

	// Statistics route
	api.HandleFunc("/statistics", handlers.GetStatistics).Methods("GET")

	// Start server
	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
