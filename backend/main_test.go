package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	// Test default value when environment variable is not set
	result := getEnv("NONEXISTENT_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestParseMinSalary(t *testing.T) {
	// Test valid salary string
	result := parseMinSalary("100000")
	if result != 100000 {
		t.Errorf("Expected 100000, got %d", result)
	}

	// Test empty string
	result = parseMinSalary("")
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}

	// Test invalid string
	result = parseMinSalary("invalid")
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(3, 200*time.Millisecond) // small window for test
	hits := 5
	blocked := 0
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	for i := 0; i < hits; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/health", nil)
		handler.ServeHTTP(rr, req)
		if rr.Code == http.StatusTooManyRequests {
			blocked++
		}
	}
	if blocked == 0 {
		t.Fatalf("expected some requests to be rate limited")
	}

	// Wait for window reset and ensure we can make requests again
	time.Sleep(250 * time.Millisecond)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/health", nil)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected OK after window reset, got %d", rr.Code)
	}
}
