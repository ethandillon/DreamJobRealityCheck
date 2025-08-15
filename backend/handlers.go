package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	Percentage       float64    `json:"percentage"`
	PercentageRegion float64    `json:"percentageRegion"`
	MatchingJobs     int        `json:"matchingJobs"`
	TotalJobs        int        `json:"totalJobs"`
	TotalJobsRegion  int        `json:"totalJobsRegion"`
	Location         string     `json:"location"`
	MinSalaryMet     bool       `json:"minSalaryMet"`
	SalaryInfo       SalaryInfo `json:"salaryInfo"`
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

	// Get total jobs count across all locations using the NATIONAL aggregated row for occ_code '00-0000'.
	// Some datasets include many '00-0000' rows (one per area). We want the SINGLE national total, which should have the
	// largest tot_emp for that occ_code. Ordering by tot_emp DESC ensures we pick the correct national aggregate even if
	// area_title filters (e.g., 'U.S.') vary or were transformed during preprocessing.
	var totalJobs int
	err = h.db.QueryRow("SELECT tot_emp FROM career_data WHERE occ_code = '00-0000' ORDER BY tot_emp DESC LIMIT 1").Scan(&totalJobs)
	if err != nil {
		return nil, fmt.Errorf("error querying total jobs: %v", err)
	}

	// Get total jobs count for the selected region/location only (denominator for regional view)
	var totalJobsRegion int
	err = h.db.QueryRow("SELECT SUM(tot_emp) FROM career_data WHERE area_title = $1", filters.Location).Scan(&totalJobsRegion)
	if err != nil {
		return nil, fmt.Errorf("error querying regional total jobs: %v", err)
	}

	// Calculate percentage
	var percentage float64
	if totalJobs > 0 && matchingJobs.Valid {
		percentage = (matchingJobs.Float64 / float64(totalJobs)) * 100
	}

	// Regional percentage
	var percentageRegion float64
	if totalJobsRegion > 0 && matchingJobs.Valid {
		percentageRegion = (matchingJobs.Float64 / float64(totalJobsRegion)) * 100
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
		Percentage:       percentage,
		PercentageRegion: percentageRegion,
		MatchingJobs:     int(matchingJobs.Float64),
		TotalJobs:        totalJobs,
		TotalJobsRegion:  totalJobsRegion,
		Location:         filters.Location,
		MinSalaryMet:     minSalaryMet,
		SalaryInfo:       salaryInfo,
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

	// Add education filter with ladder semantics
	if filters.Education != "" && filters.Education != "Any" {
		allowedEdu := getAllowedEducationValues(filters.Education)
		if len(allowedEdu) == 1 && allowedEdu[0] == "__EXACT__POSTSECONDARY_NONDEGREE__" {
			// Exact match for non-ladder value
			baseQuery += fmt.Sprintf(" AND education = $%d", argCount)
			args = append(args, "Postsecondary nondegree award")
			argCount++
		} else if len(allowedEdu) > 0 {
			placeholders := make([]string, 0, len(allowedEdu))
			for _, v := range allowedEdu {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argCount))
				args = append(args, v)
				argCount++
			}
			baseQuery += fmt.Sprintf(" AND education IN (%s)", strings.Join(placeholders, ", "))
		}
	}

	// Add experience filter with ladder semantics
	if filters.Experience != "" && filters.Experience != "Any" {
		allowedExp := getAllowedExperienceValues(filters.Experience)
		if len(allowedExp) > 0 {
			// Detect if the selection implies including rows with NULL experience
			includesNone := false
			for _, v := range allowedExp {
				if strings.EqualFold(v, "None") {
					includesNone = true
					break
				}
			}

			placeholders := make([]string, 0, len(allowedExp))
			for _, v := range allowedExp {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argCount))
				args = append(args, v)
				argCount++
			}

			if includesNone {
				// Match both NULL experience and explicit 'None' values
				baseQuery += fmt.Sprintf(" AND (experience IS NULL OR experience IN (%s))", strings.Join(placeholders, ", "))
			} else {
				baseQuery += fmt.Sprintf(" AND experience IN (%s)", strings.Join(placeholders, ", "))
			}
		}
	}

	// Add salary filter - inclusive across distribution percentiles
	if filters.MinSalary > 0 {
		baseQuery += fmt.Sprintf(" AND (a_median >= $%d OR a_pct75 >= $%d OR a_pct90 >= $%d)", argCount, argCount, argCount)
		args = append(args, filters.MinSalary)
		argCount++
	}

	return baseQuery, args
}

// getAllowedEducationValues maps a UI-selected minimum education to DB values (ladder semantics)
// Returns special marker ["__EXACT__POSTSECONDARY_NONDEGREE__"] when the selection is the non-ladder value
func getAllowedEducationValues(uiValue string) []string {
	if uiValue == "Postsecondary nondegree award" {
		return []string{"__EXACT__POSTSECONDARY_NONDEGREE__"}
	}
	// Map UI labels to DB labels for the ladder
	ladderDB := []string{
		"No formal educational credential",
		"High school diploma or equivalent",
		"Associate degree",
		"Bachelor's degree",
		"Master's degree",
		"Doctoral or professional degree",
	}
	// Normalize a couple common UI variants
	uiNorm := strings.ToLower(uiValue)
	switch uiNorm {
	case "no formal education":
		uiValue = "No formal educational credential"
	case "high school diploma":
		uiValue = "High school diploma or equivalent"
	}
	// Find index in ladder
	idx := -1
	for i, v := range ladderDB {
		if v == uiValue {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil
	}
	return ladderDB[:idx+1]
}

// getAllowedExperienceValues maps a UI-selected experience to DB values (ladder semantics)
func getAllowedExperienceValues(uiValue string) []string {
	switch uiValue {
	case "None":
		return []string{"None"}
	case "Less than 5 years":
		return []string{"None", "Less than 5 years"}
	case "5 years or more":
		return []string{"None", "Less than 5 years", "5 years or more"}
	default:
		return nil
	}
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
