package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	// Get database connection parameters from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "dream_job_calculator")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// getEnv function is defined in main.go
