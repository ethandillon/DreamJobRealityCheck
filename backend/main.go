package main

import (
	"context"
	"encoding/json"
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

	// API key middleware (no-op if API_KEYS not set)
	apiKeys := getAPIKeys()
	r.Use(apiKeyMiddleware(apiKeys))

	// Initialize handlers with database connection
	handlers := NewHandlers(db)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/calculate", handlers.CalculateHandler).Methods("GET")
	api.HandleFunc("/occupations", handlers.OccupationsHandler).Methods("GET")
	api.HandleFunc("/locations", handlers.LocationsHandler).Methods("GET")
	api.HandleFunc("/states", handlers.StatesHandler).Methods("GET")
	api.HandleFunc("/areas-by-state", handlers.AreasByStateHandler).Methods("GET")
	api.HandleFunc("/health", handlers.HealthHandler).Methods("GET") // remains open even when API key auth enabled

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

// getAPIKeys returns a slice of allowed API keys from env var API_KEYS (comma separated).
// If none provided, authentication is disabled.
func getAPIKeys() []string {
	raw := strings.TrimSpace(os.Getenv("API_KEYS"))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	keys := make([]string, 0, len(parts))
	for _, k := range parts {
		k = strings.TrimSpace(k)
		if k != "" {
			keys = append(keys, k)
		}
	}
	return keys
}

// apiKeyMiddleware enforces presence of X-API-Key header matching configured keys.
// Skips enforcement when no keys configured or for /api/health.
func apiKeyMiddleware(keys []string) mux.MiddlewareFunc {
	keySet := map[string]struct{}{}
	for _, k := range keys {
		keySet[k] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		// Disabled if no keys
		if len(keySet) == 0 {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow unauthenticated health checks
			if r.URL.Path == "/api/health" {
				next.ServeHTTP(w, r)
				return
			}
			provided := r.Header.Get("X-API-Key")
			if provided == "" {
				unauthorizedJSON(w, "missing API key")
				return
			}
			if _, ok := keySet[provided]; !ok {
				unauthorizedJSON(w, "invalid API key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func unauthorizedJSON(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
