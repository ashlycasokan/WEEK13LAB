package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn             = "root:1234@tcp(127.0.0.1:3306)/toronto_time" // Replace with the correct password if needed
	torontoTimeZone = "America/Toronto"
)

// Database connection
var db *sql.DB

// Response structure for JSON output
type Response struct {
	CurrentTime string `json:"current_time"`
	Timezone    string `json:"timezone"`
}

func main() {
	var err error

	// Establish connection to MySQL database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}()

	// Verify the database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	log.Println("Connected to MySQL database successfully!")

	// Register the /current-time endpoint
	http.HandleFunc("/current-time", currentTimeHandler)

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Load Toronto timezone
	location, err := time.LoadLocation(torontoTimeZone)
	if err != nil {
		http.Error(w, "Could not load timezone: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the current time in Toronto
	currentTime := time.Now().In(location)

	// Log the current time into the database
	if err := logTimeToDatabase(currentTime); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the JSON response
	response := Response{
		CurrentTime: currentTime.Format("2006-01-02 15:04:05"),
		Timezone:    torontoTimeZone,
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func logTimeToDatabase(timestamp time.Time) error {
	// Query to insert the current time into the time_log table
	query := "INSERT INTO time_log (timestamp) VALUES (?)"
	_, err := db.Exec(query, timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert time into database: %w", err)
	}
	return nil
}
