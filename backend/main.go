package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize router
	r := mux.NewRouter()

	// Initialize handlers with database connection
	handlers := NewHandlers(db)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/calculate", handlers.CalculateHandler).Methods("GET")
	api.HandleFunc("/occupations", handlers.OccupationsHandler).Methods("GET")
	api.HandleFunc("/locations", handlers.LocationsHandler).Methods("GET")
	api.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: getAllowedOrigins(),
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Get port from environment or use default
	port := getEnv("SERVER_PORT", "8080")

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getAllowedOrigins parses CORS_ORIGIN which may be a comma-separated list
// Defaults to allowing Vite ports 5173 and 5174 for local dev
func getAllowedOrigins() []string {
	raw := getEnv("CORS_ORIGIN", "http://localhost:5173,http://localhost:5174")
	parts := strings.Split(raw, ",")
	allowed := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			allowed = append(allowed, v)
		}
	}
	return allowed
}
