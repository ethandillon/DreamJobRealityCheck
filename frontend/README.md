# Dream Job Reality Check – Frontend

React + Vite single‑page UI for exploring how many U.S. jobs meet a user‑defined set of criteria (occupation, geography, salary, education, experience). It consumes a read‑only Go backend and presents both regional and national perspectives.

---
## Table of Contents
1. [Tech Stack](#tech-stack)
2. [Core Concept](#core-concept)
3. [Environment Variables](#environment-variables)
4. [Data Flow](#data-flow)
5. [Components](#components)
6. [API Usage](#api-usage)
7. [Request / Response Example](#request--response-example)
8. [UX / Logic Notes](#ux--logic-notes)
9. [Development](#development)
10. [Build & Deploy](#build--deploy)
11. [Future Enhancements](#future-enhancements)

---
## Tech Stack
- React 19 + Vite (fast dev server & HMR)
- Tailwind CSS (utility styling)
- Headless UI for accessible primitives
- Fetch API (no external state/query library for simplicity)

## Core Concept
User defines filters → frontend gathers them → backend calculates proportion of jobs (regional & national) that meet those criteria → UI displays percentages and formatted counts with a toggle between regional and national perspectives.

## Environment Variables
| Variable | Purpose | Default |
|----------|---------|---------|
| `VITE_API_BASE_URL` | Base URL of backend (no trailing slash) | `http://localhost:8080` |

Example `.env.local`:
```
VITE_API_BASE_URL=http://localhost:8080
```

## Data Flow
1. `Filters` component loads selectable data (`occupations`, `states`, then `areas`) using API client helpers.
2. User chooses Occupation → State → Area + adjusts salary slider / education / experience.
3. Clicking the action button triggers `calculate(filters)` → `GET /api/calculate?location=...&occupation=...&minSalary=...&education=...&experience=...`.
4. Result stored in `App` state and passed to `Results`.
5. `Results` shows percentage & counts; user can toggle between national/regional views (no new API call; just different fields in same payload).

## Components
| Component | Purpose |
|-----------|---------|
| `App` | Orchestrates layout and ties together `Filters` + `Results`. Holds result/error/loading state. |
| `Filters` | Manages form state, fetches supporting lists, validates required selections, emits filter payload. |
| `Results` | Renders loading, error, placeholder, or computed result with toggle for national vs regional. |
| `AnimatedGradientBorder` | Decorative container with continuously animating gradient border around results panel. |
| `CustomSelect` | Simple dropdown for fixed option arrays (education, experience). |
| `SearchableDropdown` | Filterable list for large option sets (occupations, states, areas). |
| `DataInfoModal` | Modal explaining underlying dataset / attribution. |
| `api/client.js` | Centralized fetch wrapper (query building & error handling). |

## API Usage
Endpoints consumed:
| Path | Used In | Purpose |
|------|---------|---------|
| `/api/occupations` | `Filters` (on mount) | Populate occupation dropdown |
| `/api/states` | `Filters` (on mount) | Populate state dropdown |
| `/api/areas-by-state?state=...` | `Filters` (state change) | Populate area list for selected state |
| `/api/calculate` | `Filters` submit | Fetch calculation result |
| `/api/health` | (optional) | Basic availability check (not wired to UI) |

Each lookup endpoint returns a JSON object with a plural key (e.g. `{ "occupations": [...] }`). The calculate endpoint returns a structured result (see below).

## Request / Response Example
Request:
```
GET /api/calculate?location=Detroit-Warren-Dearborn%2C+MI&occupation=Software+Developers&minSalary=70000&education=Bachelor's%20degree&experience=None
```
Response (abbreviated):
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

Fields used directly in UI today: `percentage`, `percentageRegion`, `matchingJobs`, `totalJobs`, `totalJobsRegion`, `location` (salaryInfo currently hidden but preserved for future expansion).

## UX / Logic Notes
- Form cannot submit until Occupation + State + Area are selected (and occupation not blank).
- Salary slider (30,000 → 250,000 step 5,000) updates value label live.
- Education / Experience default to "Any" which omits those filters.
- Percentages: UI formats trailing zeros intelligently (expanding decimals only until a significant digit appears; pure zero collapses to `0`).
- National vs Regional toggle is client-side (no re-fetch) to minimize latency and API load.
- Errors bubble up from fetch helper; message displayed within styled box.

## Development
Install deps and start dev server:
```
npm install
npm run dev
```
Frontend expects the backend reachable at `VITE_API_BASE_URL` (default localhost:8080). Ensure CORS on backend allows the Vite dev origin (5173 by default).

Lint:
```
npm run lint
```

Build production bundle:
```
npm run build
```
Preview production build locally:
```
npm run preview
```

## Build & Deploy
Produces a static build (`dist/`) that can be served by any static host (e.g., Vercel, Netlify, Fly static, S3 + CloudFront). Ensure runtime environment injects the correct `VITE_API_BASE_URL` at build time (environment variables must be present during `npm run build`). If multiple environments are needed, create per‑env builds or introduce a lightweight runtime config script.

## Future Enhancements
- Persist recent user selections in localStorage for session continuity.
- Add a small sparkline or distribution visualization using `salaryInfo` fields.
- Introduce skeleton loading states for dropdowns.
- Replace manual fetches with a query caching layer (React Query / TanStack Query) if complexity increases.
- Accessibility audit: ensure focus trapping in `DataInfoModal` and ARIA roles in custom dropdowns.
- Internationalization (metric currency formatting, localization of labels).

---
This README focuses on how the current frontend works and integrates with the backend rather than a generic Vite template description.
