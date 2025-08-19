# Dream Job Reality Check – Backend

Go (1.21) service powering the Dream Job vs Reality calculator. Provides a read‑only REST API over a processed labor statistics dataset stored in PostgreSQL (`career_data` table). For any set of user criteria it returns:
1. Regional perspective – how many jobs in the selected area meet the criteria.
2. National perspective – how many jobs across the U.S. meet the same criteria.

---

## Table of Contents
1. [Features](#features)
2. [Architecture Overview](#architecture-overview)
3. [User Inputs & Query Mapping](#user-inputs--query-mapping)
4. [Data Model](#data-model)
5. [Business Logic Conventions](#business-logic-conventions)
6. [API Endpoints](#api-endpoints)
7. [Rate Limiting](#rate-limiting)
8. [CORS](#cors)
9. [Environment Variables](#environment-variables)
10. [Request / Response Example](#request--response-example)
11. [Deployment](#deployment)
12. [Implementation Notes](#implementation-notes)
13. [Future Improvements](#future-improvements)
14. [Key Files](#key-files)

---

## Features
- Read‑only JSON API (idempotent GET endpoints only)
- Dynamic filtering with parameterized SQL
- Inclusive salary distribution filtering
- Education & experience “ladder” semantics
- Dual regional + national metrics
- Lightweight in‑memory per‑IP rate limiting (100 req/min)
- Graceful shutdown & conservative HTTP timeouts

## Architecture Overview
Frontend gathers filters → calls `GET /api/calculate` → backend composes WHERE clauses only for provided filters → database returns aggregated employment & salary distribution → backend augments with national denominator (largest `tot_emp` for `occ_code='00-0000'`) → response includes counts, percentages, and salary band snapshot.

## User Inputs & Query Mapping
Inputs (UI → query params on `/api/calculate`):
| UI Field | Param | Mapping Notes |
|----------|-------|--------------|
| Occupation | `occupation` | Matches `occ_title` (case-insensitive) |
| State | `location` or used in `/api/states` | Distinct state-level `area_title` |
| Area within State | `location` | Full `area_title` string (metro / non-metro) |
| Minimum annual salary | `minSalary` | Compared across percentile fields (see logic) |
| Minimum education | `education` | Ladder mapping helper expands allowed values |
| Required work experience | `experience` | Ladder mapping helper expands allowed values |


## Data Model
Single fact table `career_data` (one row per `(area_title, occ_code)` pair):

| Column      | Type           | Null | Description |
|-------------|----------------|------|-------------|
| id          | SERIAL (int)   | NO   | Surrogate primary key |
| area_title  | VARCHAR(255)   | NO   | Geographic area / metro / state / non‑metro label |
| occ_code    | VARCHAR(15)    | NO   | Standard occupation code (e.g. `15-1252`) |
| occ_title   | VARCHAR(255)   | YES  | Human readable occupation title |
| education   | VARCHAR(255)   | YES  | Minimum education requirement label (as published / normalized) |
| experience  | VARCHAR(255)   | YES  | Required prior work experience label |
| tot_emp     | INTEGER        | YES  | Total employment count for the occupation in the area |
| a_median    | INTEGER        | YES  | Median annual wage (USD) |
| a_pct10     | INTEGER        | YES  | 10th percentile annual wage |
| a_pct25     | INTEGER        | YES  | 25th percentile annual wage |
| a_pct75     | INTEGER        | YES  | 75th percentile annual wage |
| a_pct90     | INTEGER        | YES  | 90th percentile annual wage |

Composite uniqueness: `(area_title, occ_code)` ensures no duplicate occupation entries per area.

## Business Logic Conventions
- National denominator: select the row with largest `tot_emp` where `occ_code='00-0000'`.
- Salary filter: if ANY of `a_median, a_pct10, a_pct25, a_pct75, a_pct90` ≥ `minSalary`, the record qualifies (broad/inclusive to surface potential career paths even when central tendency is lower).
- Education & experience: ladder semantics include higher levels automatically except explicit non-ladder cases handled in helpers.
- Percentages: `percentageRegion = matchingJobs / totalJobsRegion`, `percentage = matchingJobs / totalJobs` (national).

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/calculate` | Returns employment match metrics & salary info |
| GET | `/api/occupations` | Distinct `occ_title` values |
| GET | `/api/locations` | Distinct non-national `area_title` values |
| GET | `/api/states` | State-level area titles (no commas) |
| GET | `/api/areas-by-state?state=STATE_NAME` | All granular areas for the state |
| GET | `/api/health` | Liveness/health check |

All endpoints return JSON and are safe to cache (dataset is static for end users).


## Rate Limiting
In-memory per-IP (100 req/min/IP). Keys on first valid client IP from headers: `X-Forwarded-For`, `X-Real-IP`, `CF-Connecting-IP`; falls back to `RemoteAddr`. Suitable for single instance or low scale. For horizontal scale, replace with shared store (Redis) or external gateway.

## CORS
Configured via `CORS_ORIGIN` (comma-separated). Local default: `http://localhost:5173,http://localhost:5174`. Parsed into allowed origins slice in `main.go`.

## Environment Variables
| Variable | Purpose | Example |
|----------|---------|---------|
| DB_HOST | PostgreSQL host | `localhost` |
| DB_PORT | PostgreSQL port | `5432` |
| DB_USER | DB user | `appuser` |
| DB_PASSWORD | DB password | `***` |
| DB_NAME | Database name | `dream_job` |
| DB_SSLMODE | TLS mode (`disable` local / `require` prod) | `require` |
| SERVER_PORT | HTTP listen port | `8080` |
| CORS_ORIGIN | Allowed origins (comma list) | `https://dream-job-reality-check.vercel.app` |

## Request / Response Example
Request:
```
GET /api/calculate?location=Detroit-Warren-Dearborn%2C+MI&occupation=Software+Developers&minSalary=70000
```
Response:
```json
{
  "percentage": 0.0158902766192261,
  "percentageRegion": 0.638723083235174,
  "matchingJobs": 24130,
  "totalJobs": 151853870,
  "totalJobsRegion": 3777850,
  "location": "Detroit-Warren-Dearborn, MI",
  "minSalaryMet": true,
  "salaryInfo": {
    "medianSalary": 107490,
    "pct10Salary": 76740,
    "pct25Salary": 90060,
    "pct75Salary": 133160,
    "pct90Salary": 162950
  }
}
```

## Deployment
Containerized and deployed on Fly.io. Secrets configured with `fly secrets` (DB credentials + CORS origins). Health check uses `/api/health`. Stateless binary; scaling is linear (add instances) since all persistence is in Postgres.

## Implementation Notes
- HTTP server timeouts: read 15s, write 30s, idle 60s, header 10s.
- `MaxHeaderBytes` capped to mitigate oversized header attacks.
- Graceful shutdown on SIGINT/SIGTERM with 10s context.
- Parameterized SQL only (no string concatenation of user input).

## Future Improvements
- Replace in-memory rate limiter with distributed backend (Redis) for multi-instance deployments.
- Cache static lookup endpoints (`/api/occupations`, `/api/states`).
- Precompute national denominator once at startup.
- Externalize education/experience ladders into reference table.
- Add observability: structured logging & basic metrics (latency, error counts).

## Key Files
| File | Purpose |
|------|---------|
| `main.go` | Server bootstrap, routing, middleware, shutdown |
| `handlers.go` | Request parsing, query building, response formatting |
| `rate_limiter.go` | In-memory per-IP rate limiting middleware |
| `database.go` | PostgreSQL connection initialization |

---
