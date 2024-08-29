package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() (*sql.DB, error) {
	// Open a connection to the SQLite database file (or create it if it doesn't exist)
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create a sample table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		startDate TEXT,
		startTime TEXT,
		endDate TEXT,
		endTime TEXT,
		description TEXT,
		created_at NOT NULL DEFAULT current_timestamp,
		updated_at NOT NULL DEFAULT current_timestamp
	);
	`
	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

func main() {
	mux := http.NewServeMux()
	db, err := initDB()
	if err != nil {
		log.Fatalln(err)
	}
	h := &Handler{
		DB: db,
	}
	mux.HandleFunc("GET /events", h.GetEvents())
	log.Println("listening on port 8000")
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
