package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	// Prefer full connection URL if provided (e.g., Supabase)
	if raw := getEnv("DATABASE_URL", ""); raw != "" {
		conn := ensureSSLModeInURL(raw)
		db, err := sql.Open("postgres", conn)
		if err != nil {
			return nil, fmt.Errorf("error opening database using DATABASE_URL")
		}
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("error connecting to database using DATABASE_URL")
		}
		log.Println("Successfully connected to database via DATABASE_URL")
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(0)
		return db, nil
	}
	if raw := getEnv("DB_URL", ""); raw != "" {
		conn := ensureSSLModeInURL(raw)
		db, err := sql.Open("postgres", conn)
		if err != nil {
			return nil, fmt.Errorf("error opening database using DB_URL")
		}
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("error connecting to database using DB_URL")
		}
		log.Println("Successfully connected to database via DB_URL")
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(0)
		return db, nil
	}

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
	// Basic pooling defaults; override via env in future if needed
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)
	return db, nil
}

// getEnv function is defined in main.go

// ensureSSLModeInURL appends sslmode=require to Postgres URLs that do not already specify sslmode
// If the provided string does not look like a URL (no "://"), it is returned unchanged
func ensureSSLModeInURL(conn string) string {
	if !strings.Contains(conn, "://") {
		return conn
	}
	parsed, err := url.Parse(conn)
	if err != nil {
		return conn
	}
	// Only adjust for postgres schemes
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return conn
	}
	q := parsed.Query()
	if q.Get("sslmode") == "" {
		q.Set("sslmode", "require")
		parsed.RawQuery = q.Encode()
		return parsed.String()
	}
	return conn
}
