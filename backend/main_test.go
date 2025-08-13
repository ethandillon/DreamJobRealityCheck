package main

import (
	"testing"
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
