# Building a Go API for Current Toronto Time with MySQL Database Logging


## 1.Set Up MySQL Database:
Initialize MySQL:
.\mysqld --initialize --console
.\mysqld

•	SHOW DATABASES;

This command lists all the databases available in the MySQL server.

•	CREATE DATABASE toronto_time;

This command creates a new database named toronto_time in MySQL.

•	USE toronto_time;

This command switches the active database to toronto_time, so that any subsequent operations (like creating tables, querying data, etc.) are done within this database.

•	CREATE TABLE time_log ( id INT AUTO_INCREMENT PRIMARY KEY, timestamp DATETIME NOT NULL);

This command creates a table named time_log inside the toronto_time database.
 id INT AUTO_INCREMENT PRIMARY KEY: This creates a column named id with an integer data type. AUTO_INCREMENT ensures that each new entry in the table gets a unique, automatically generated value starting from 1. The PRIMARY KEY constraint ensures that this column uniquely identifies each row in the table.
timestamp DATETIME NOT NULL: This creates a column named timestamp of type DATETIME, which will store date and time values. The NOT NULL constraint ensures that every row must have a value for this column. 


•	INSERT INTO time_log (timestamp) VALUES (NOW());

This command inserts a new row into the time_log table. The value for the timestamp column is set to the current date and time (NOW() is a MySQL function that returns the current date and time).

•	SELECT * FROM time_log;

This command retrieves all rows from the time_log table. The * symbol means "all columns," so this will fetch all records from the table, including the id and timestamp.

| id  | timestamp           |
|-----|---------------------|
|  1  | 2024-11-28 12:34:56 |


![image](https://github.com/user-attachments/assets/1c0b65ca-5fe6-4c42-88d1-56ae7f503774)

 


## 2.API Development:

## Go API for Current Toronto Time with MySQL Database Logging

This project provides a simple API that returns the current time in Toronto and logs it to a MySQL database.

### Code

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

const (
    dsn               = "root:1234@tcp(127.0.0.1:3306)/toronto_time" // Replace with the correct password if needed
    torontoTimeZone    = "America/Toronto"
)

// Database connection
var db *sql.DB

// Response structure for JSON output
type Response struct {
    CurrentTime string `json:"current_time"`
    Timezone    string `json:"timezone"`
}

func init() {
    // Open a log file for writing
    logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }

    // Set log output to the log file
    log.SetOutput(logFile)
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

    // Register the /current-time and /logs endpoints
    http.HandleFunc("/current-time", currentTimeHandler)
    http.HandleFunc("/logs", logsHandler)

    // Start the HTTP server
    log.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

// Handler for /current-time endpoint
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

// Log the current time into the database
func logTimeToDatabase(timestamp time.Time) error {
    // Query to insert the current time into the time_log table
    query := "INSERT INTO time_log (timestamp) VALUES (?)"
    _, err := db.Exec(query, timestamp)
    if err != nil {
        return fmt.Errorf("failed to insert time into database: %w", err)
    }
    return nil
}

// Handler for /logs endpoint
func logsHandler(w http.ResponseWriter, r *http.Request) {
    // Query to fetch all logs
    rows, err := db.Query("SELECT id, timestamp FROM time_log")
    if err != nil {
        http.Error(w, "Failed to retrieve logs: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Create a slice to store the logs
    var logs []struct {
        ID        int    `json:"id"`
        Timestamp string `json:"timestamp"`
    }

    // Iterate through the rows and add them to the logs slice
    for rows.Next() {
        var log struct {
            ID        int    `json:"id"`
            Timestamp string `json:"timestamp"`
        }
        if err := rows.Scan(&log.ID, &log.Timestamp); err != nil {
            http.Error(w, "Failed to scan log: "+err.Error(), http.StatusInternalServerError)
            return
        }
        logs = append(logs, log)
    }

    // Check for errors while iterating
    if err := rows.Err(); err != nil {
        http.Error(w, "Error processing rows: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the logs as a JSON response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(logs); err != nil {
        http.Error(w, "Failed to encode logs: "+err.Error(), http.StatusInternalServerError)
    }
}


```

## 3.Time Zone Conversion in JSON:

![image](https://github.com/user-attachments/assets/bd820c6d-e2d6-4b84-b332-96301b22511a)


## 4.Database Connection:
![image](https://github.com/user-attachments/assets/2c1c8c82-6bab-49da-936e-080d9e880b76)



 


 
 

