package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Filters represents the query parameters from the frontend
type Filters struct {
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
	MinSalary  int    `json:"minSalary"`
	Education  string `json:"education"`
	Experience string `json:"experience"`
}

// CalculationResult represents the response data
type CalculationResult struct {
	Percentage   float64    `json:"percentage"`
	MatchingJobs int        `json:"matchingJobs"`
	TotalJobs    int        `json:"totalJobs"`
	Location     string     `json:"location"`
	MinSalaryMet bool       `json:"minSalaryMet"`
	SalaryInfo   SalaryInfo `json:"salaryInfo"`
}

// SalaryInfo provides detailed salary information
type SalaryInfo struct {
	MedianSalary int `json:"medianSalary"`
	Pct10Salary  int `json:"pct10Salary"`
	Pct25Salary  int `json:"pct25Salary"`
	Pct75Salary  int `json:"pct75Salary"`
	Pct90Salary  int `json:"pct90Salary"`
}

// Handlers struct holds the database connection
type Handlers struct {
	db *sql.DB
}

// NewHandlers creates a new Handlers instance
func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{db: db}
}

// CalculateHandler handles the /api/calculate endpoint
func (h *Handlers) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := Filters{
		Location:   r.URL.Query().Get("location"),
		Occupation: r.URL.Query().Get("occupation"),
		MinSalary:  parseMinSalary(r.URL.Query().Get("minSalary")),
		Education:  r.URL.Query().Get("education"),
		Experience: r.URL.Query().Get("experience"),
	}

	// Validate required fields
	if filters.Location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	// Calculate results based on filters
	result, err := h.calculateJobOpportunities(filters)
	if err != nil {
		log.Printf("Error calculating job opportunities: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HealthHandler provides a simple health check endpoint
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// calculateJobOpportunities performs the main calculation logic
func (h *Handlers) calculateJobOpportunities(filters Filters) (*CalculationResult, error) {
	// Build the SQL query based on filters
	query, args := buildQuery(filters)

	// Execute the query to get matching jobs count and salary info
	var matchingJobs int
	var medianSalary, pct10Salary, pct25Salary, pct75Salary, pct90Salary sql.NullInt64
	var totalEmp sql.NullInt64

	err := h.db.QueryRow(query, args...).Scan(
		&matchingJobs, &medianSalary, &pct10Salary, &pct25Salary, &pct75Salary, &pct90Salary, &totalEmp,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying matching jobs: %v", err)
	}

	// Get total jobs count across all locations
	var totalJobs int
	err = h.db.QueryRow("SELECT SUM(TOT_EMP) FROM jobs").Scan(&totalJobs)
	if err != nil {
		return nil, fmt.Errorf("error querying total jobs: %v", err)
	}

	// Calculate percentage
	var percentage float64
	if totalJobs > 0 {
		percentage = (float64(matchingJobs) / float64(totalJobs)) * 100
	}

	// Check if minimum salary requirement is met
	minSalaryMet := false
	if medianSalary.Valid && filters.MinSalary > 0 {
		minSalaryMet = medianSalary.Int64 >= int64(filters.MinSalary)
	}

	// Build salary info
	salaryInfo := SalaryInfo{
		MedianSalary: int(medianSalary.Int64),
		Pct10Salary:  int(pct10Salary.Int64),
		Pct25Salary:  int(pct25Salary.Int64),
		Pct75Salary:  int(pct75Salary.Int64),
		Pct90Salary:  int(pct90Salary.Int64),
	}

	return &CalculationResult{
		Percentage:   percentage,
		MatchingJobs: matchingJobs,
		TotalJobs:    totalJobs,
		Location:     filters.Location,
		MinSalaryMet: minSalaryMet,
		SalaryInfo:   salaryInfo,
	}, nil
}

// buildQuery constructs the SQL query and arguments based on filters
func buildQuery(filters Filters) (string, []interface{}) {
	baseQuery := `
		SELECT 
			SUM(TOT_EMP) as matching_jobs,
			AVG(A_MEDIAN) as median_salary,
			AVG(A_PCT10) as pct10_salary,
			AVG(A_PCT25) as pct25_salary,
			AVG(A_PCT75) as pct75_salary,
			AVG(A_PCT90) as pct90_salary,
			SUM(TOT_EMP) as total_emp
		FROM jobs 
		WHERE 1=1`

	var args []interface{}
	argCount := 1

	// Add location filter
	if filters.Location != "" {
		baseQuery += fmt.Sprintf(" AND AREA_TITLE ILIKE $%d", argCount)
		args = append(args, "%"+filters.Location+"%")
		argCount++
	}

	// Add occupation filter
	if filters.Occupation != "" {
		baseQuery += fmt.Sprintf(" AND OCC_TITLE ILIKE $%d", argCount)
		args = append(args, "%"+filters.Occupation+"%")
		argCount++
	}

	// Add education filter
	if filters.Education != "" {
		baseQuery += fmt.Sprintf(" AND Education = $%d", argCount)
		args = append(args, filters.Education)
		argCount++
	}

	// Add experience filter
	if filters.Experience != "" {
		baseQuery += fmt.Sprintf(" AND Experience = $%d", argCount)
		args = append(args, filters.Experience)
		argCount++
	}

	// Add salary filter - check if median salary meets minimum requirement
	if filters.MinSalary > 0 {
		baseQuery += fmt.Sprintf(" AND A_MEDIAN >= $%d", argCount)
		args = append(args, filters.MinSalary)
		argCount++
	}

	return baseQuery, args
}

// parseMinSalary converts the minSalary string to an integer
func parseMinSalary(salaryStr string) int {
	if salaryStr == "" {
		return 0
	}

	salary, err := strconv.Atoi(salaryStr)
	if err != nil {
		return 0
	}

	return salary
}
