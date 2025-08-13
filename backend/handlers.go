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

// OccupationsHandler provides a list of unique occupation titles
func (h *Handlers) OccupationsHandler(w http.ResponseWriter, r *http.Request) {
	// Query to get unique occupation titles
	query := "SELECT DISTINCT occ_title FROM career_data WHERE occ_title IS NOT NULL AND occ_title != '' ORDER BY occ_title"

	rows, err := h.db.Query(query)
	if err != nil {
		log.Printf("Error querying occupations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var occupations []string
	for rows.Next() {
		var occTitle string
		if err := rows.Scan(&occTitle); err != nil {
			log.Printf("Error scanning occupation: %v", err)
			continue
		}
		occupations = append(occupations, occTitle)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating occupations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"occupations": occupations,
		"count":       len(occupations),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// LocationsHandler provides a list of unique area titles (locations)
// Excludes generic U.S.-wide labels
func (h *Handlers) LocationsHandler(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT DISTINCT area_title
        FROM career_data
        WHERE area_title IS NOT NULL
          AND area_title <> ''
          AND area_title NOT IN ('U.S.', 'United States', 'USA', 'US')
        ORDER BY area_title`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Printf("Error querying locations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []string
	for rows.Next() {
		var area string
		if err := rows.Scan(&area); err != nil {
			log.Printf("Error scanning location: %v", err)
			continue
		}
		locations = append(locations, area)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating locations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"locations": locations,
		"count":     len(locations),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// StatesHandler returns distinct state-level area titles
func (h *Handlers) StatesHandler(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT DISTINCT area_title
        FROM career_data
        WHERE area_title IS NOT NULL
          AND area_title <> ''
          AND area_title NOT ILIKE '%,%'
          AND area_title NOT ILIKE '%nonmetropolitan area%'
          AND area_title NOT IN ('U.S.', 'United States', 'USA', 'US')
        ORDER BY area_title`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Printf("Error querying states: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var states []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			log.Printf("Error scanning state: %v", err)
			continue
		}
		states = append(states, s)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating states: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"states": states,
		"count":  len(states),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// AreasByStateHandler returns all area titles relevant to a given state
func (h *Handlers) AreasByStateHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return
	}
	abbr := stateNameToAbbr(state)
	// Build patterns:
	// 1) exact state name
	// 2) MSAs that have ", {ABBR}" or ", {ABBR}-" after the comma
	// 3) nonmetropolitan areas containing the state name
	query := `
        SELECT DISTINCT area_title
        FROM career_data
        WHERE area_title = $1
           OR area_title ILIKE '%' || $2 || '%'
           OR area_title ILIKE '%' || $3 || '%'
        ORDER BY area_title`
	commaPattern := ", " + abbr // matches ", GA" including cross-state like ", GA-SC"
	nonMetroPattern := state + " nonmetropolitan area"

	rows, err := h.db.Query(query, state, commaPattern, nonMetroPattern)
	if err != nil {
		log.Printf("Error querying areas by state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var areas []string
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			log.Printf("Error scanning area: %v", err)
			continue
		}
		areas = append(areas, a)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating areas: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"areas": areas,
		"count": len(areas),
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// stateNameToAbbr maps state full names to USPS abbreviations
func stateNameToAbbr(state string) string {
	m := map[string]string{
		"Alabama": "AL", "Alaska": "AK", "Arizona": "AZ", "Arkansas": "AR", "California": "CA",
		"Colorado": "CO", "Connecticut": "CT", "Delaware": "DE", "District of Columbia": "DC",
		"Florida": "FL", "Georgia": "GA", "Hawaii": "HI", "Idaho": "ID", "Illinois": "IL",
		"Indiana": "IN", "Iowa": "IA", "Kansas": "KS", "Kentucky": "KY", "Louisiana": "LA",
		"Maine": "ME", "Maryland": "MD", "Massachusetts": "MA", "Michigan": "MI", "Minnesota": "MN",
		"Mississippi": "MS", "Missouri": "MO", "Montana": "MT", "Nebraska": "NE", "Nevada": "NV",
		"New Hampshire": "NH", "New Jersey": "NJ", "New Mexico": "NM", "New York": "NY",
		"North Carolina": "NC", "North Dakota": "ND", "Ohio": "OH", "Oklahoma": "OK", "Oregon": "OR",
		"Pennsylvania": "PA", "Rhode Island": "RI", "South Carolina": "SC", "South Dakota": "SD",
		"Tennessee": "TN", "Texas": "TX", "Utah": "UT", "Vermont": "VT", "Virginia": "VA",
		"Washington": "WA", "West Virginia": "WV", "Wisconsin": "WI", "Wyoming": "WY",
		"Puerto Rico": "PR", "Guam": "GU", "Virgin Islands": "VI",
	}
	if v, ok := m[state]; ok {
		return v
	}
	return state // fallback
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
	var matchingJobs sql.NullFloat64
	var medianSalary, pct10Salary, pct25Salary, pct75Salary, pct90Salary sql.NullFloat64
	var totalEmp sql.NullFloat64

	err := h.db.QueryRow(query, args...).Scan(
		&matchingJobs, &medianSalary, &pct10Salary, &pct25Salary, &pct75Salary, &pct90Salary, &totalEmp,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying matching jobs: %v", err)
	}

	// Get total jobs count across all locations
	var totalJobs int
	err = h.db.QueryRow("SELECT SUM(tot_emp) FROM career_data").Scan(&totalJobs)
	if err != nil {
		return nil, fmt.Errorf("error querying total jobs: %v", err)
	}

	// Calculate percentage
	var percentage float64
	if totalJobs > 0 && matchingJobs.Valid {
		percentage = (matchingJobs.Float64 / float64(totalJobs)) * 100
	}

	// Check if minimum salary requirement is met
	minSalaryMet := false
	if medianSalary.Valid && filters.MinSalary > 0 {
		minSalaryMet = medianSalary.Float64 >= float64(filters.MinSalary)
	}

	// Build salary info
	salaryInfo := SalaryInfo{
		MedianSalary: int(medianSalary.Float64),
		Pct10Salary:  int(pct10Salary.Float64),
		Pct25Salary:  int(pct25Salary.Float64),
		Pct75Salary:  int(pct75Salary.Float64),
		Pct90Salary:  int(pct90Salary.Float64),
	}

	return &CalculationResult{
		Percentage:   percentage,
		MatchingJobs: int(matchingJobs.Float64),
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
			SUM(tot_emp) as matching_jobs,
			AVG(a_median) as median_salary,
			AVG(a_pct10) as pct10_salary,
			AVG(a_pct25) as pct25_salary,
			AVG(a_pct75) as pct75_salary,
			AVG(a_pct90) as pct90_salary,
			SUM(tot_emp) as total_emp
		FROM career_data 
		WHERE 1=1`

	var args []interface{}
	argCount := 1

	// Add location filter
	if filters.Location != "" {
		baseQuery += fmt.Sprintf(" AND area_title ILIKE $%d", argCount)
		args = append(args, "%"+filters.Location+"%")
		argCount++
	}

	// Add occupation filter
	if filters.Occupation != "" {
		baseQuery += fmt.Sprintf(" AND occ_title ILIKE $%d", argCount)
		args = append(args, "%"+filters.Occupation+"%")
		argCount++
	}

	// Add education filter
	if filters.Education != "" {
		baseQuery += fmt.Sprintf(" AND education = $%d", argCount)
		args = append(args, filters.Education)
		argCount++
	}

	// Add experience filter
	if filters.Experience != "" {
		baseQuery += fmt.Sprintf(" AND experience = $%d", argCount)
		args = append(args, filters.Experience)
		argCount++
	}

	// Add salary filter - check if median salary meets minimum requirement
	if filters.MinSalary > 0 {
		baseQuery += fmt.Sprintf(" AND a_median >= $%d", argCount)
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
