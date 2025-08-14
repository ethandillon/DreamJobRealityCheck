# Deployment Guide

## Backend (Go)
- Build container:
  - docker build -t dream-job-api:latest .
- Required env vars:
  - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
  - DB_SSLMODE=require (prod)
  - SERVER_PORT=8080
  - CORS_ORIGIN=https://your-frontend-domain
- Run:
  - docker run -p 8080:8080 --env-file /path/to/prod.env dream-job-api:latest

## Database
- Ensure indexes:

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS gin_area_trgm ON career_data USING gin (area_title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS gin_occ_trgm  ON career_data USING gin (occ_title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_edu      ON career_data(education);
CREATE INDEX IF NOT EXISTS idx_exp      ON career_data(experience);

## Frontend (Vite)
- Set API base:
  - frontend/.env.production -> VITE_API_BASE_URL=https://your-api-domain
- Build:
  - npm ci && npm run build
- Deploy frontend/dist/ to static hosting (Vercel/Netlify/S3+CloudFront).

## Hardening
- Graceful shutdown, server timeouts (already enabled).
- Restrict CORS to production domain(s).
- Use SSL to DB.

