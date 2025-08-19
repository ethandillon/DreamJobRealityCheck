package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	api.HandleFunc("/states", handlers.StatesHandler).Methods("GET")
	api.HandleFunc("/areas-by-state", handlers.AreasByStateHandler).Methods("GET")
	api.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// Attach rate limiter (100 req/min/IP)
	limiter := NewRateLimiter(100, time.Minute)
	api.Use(limiter.Middleware)

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: getAllowedOrigins(),
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"*"},
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Get port from environment or use default
	port := getEnv("SERVER_PORT", "8080")

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
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
