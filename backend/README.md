# Dream Job Calculator Backend

A Go-based REST API backend for the Dream Job Calculator application that connects to PostgreSQL to calculate job opportunities based on user criteria.

## Features

- RESTful API endpoints for job calculations
- PostgreSQL database integration
- CORS support for frontend integration
- Environment-based configuration
- Health check endpoint

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Setup

### 1. Install Dependencies

```bash
go mod tidy
```

**Note**: Make sure your existing PostgreSQL database has a `jobs` table with the required structure. The backend expects this table to exist and be properly configured.

### 2. Database Setup

Since you already have a PostgreSQL database set up, ensure your database has a `jobs` table with the following structure (based on your data schema):

- `id`: Primary key (SERIAL)
- `AREA_TITLE`: Location/area (VARCHAR)
- `OCC_CODE`: Occupation code (VARCHAR)
- `OCC_TITLE`: Occupation title (VARCHAR)
- `Education`: Required education level (VARCHAR)
- `Experience`: Required experience level (VARCHAR)
- `TOT_EMP`: Total employment count (INTEGER)
- `A_MEDIAN`: Median annual salary (INTEGER)
- `A_PCT10`: 10th percentile salary (INTEGER)
- `A_PCT25`: 25th percentile salary (INTEGER)
- `A_PCT75`: 75th percentile salary (INTEGER)
- `A_PCT90`: 90th percentile salary (INTEGER)

### 3. Environment Configuration

Copy the `.env.example` file to `.env` and update the values:

```bash
cp .env.example .env
```

Update the `.env` file with your database credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=dream_job_calculator
DB_SSLMODE=disable

SERVER_PORT=8080
CORS_ORIGIN=http://localhost:5173
```

### 4. Run the Application

```bash
go run .
```

The server will start on port 8080 (or the port specified in your .env file).

## API Endpoints

### GET /api/calculate

Calculates job opportunities based on provided filters.

**Query Parameters:**
- `location` (required): City and state (e.g., "San Francisco, CA")
- `occupation` (optional): Job title or field
- `minSalary` (optional): Minimum annual salary
- `education` (optional): Required education level
- `experience` (optional): Required work experience

**Example Request:**
```
GET /api/calculate?location=San Francisco, CA&occupation=Software Developer&minSalary=100000&education=Bachelor's degree&experience=2-4 years
```

**Response:**
```json
{
  "percentage": 12.5,
  "matchingJobs": 1250,
  "totalJobs": 10000,
  "location": "San Francisco, CA",
  "minSalaryMet": true,
  "salaryInfo": {
    "medianSalary": 120000,
    "pct10Salary": 80000,
    "pct25Salary": 95000,
    "pct75Salary": 150000,
    "pct90Salary": 180000
  }
}
```

### GET /api/health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Database Schema

The main table is `jobs` with the following structure:

- `id`: Primary key
- `title`: Job title
- `occupation`: Job category
- `location`: City and state
- `annual_salary`: Annual salary in USD
- `required_education`: Minimum education requirement
- `required_experience`: Minimum experience requirement
- `company`: Company name
- `industry`: Industry sector
- `created_at`: Record creation timestamp

## Development

### Project Structure

```
backend/
├── main.go          # Main application entry point
├── database.go      # Database connection and initialization
├── handlers.go      # HTTP request handlers
├── go.mod           # Go module file
├── go.sum           # Go dependencies checksum
└── .env             # Environment configuration
```

### Adding New Endpoints

1. Add the route in `main.go`
2. Create the handler function in `handlers.go`
3. Update the `Handlers` struct if needed

### Database Queries

All database queries are built dynamically based on the provided filters. The `buildQuery` function in `handlers.go` constructs parameterized SQL queries to prevent SQL injection.

## Testing

To test the API endpoints:

```bash
# Health check
curl http://localhost:8080/api/health

# Calculate job opportunities
curl "http://localhost:8080/api/calculate?location=San Francisco, CA&occupation=Software Developer&minSalary=100000"
```

## Deployment

1. Build the binary:
```bash
go build -o dream-job-calculator .
```

2. Set environment variables in production
3. Run the binary with proper database credentials
4. Ensure CORS origins are configured for your production domain

## Troubleshooting

### Common Issues

1. **Database Connection Failed**: Check your `.env` file and ensure PostgreSQL is running
2. **Port Already in Use**: Change the `SERVER_PORT` in your `.env` file
3. **CORS Errors**: Verify the `CORS_ORIGIN` setting matches your frontend URL

### Logs

The application logs connection status and any errors to stdout. Check the console output for debugging information.
