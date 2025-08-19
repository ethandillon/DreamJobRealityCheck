# Data Processing Pipeline

This folder contains the one‑off / reproducible scripts used to transform raw Bureau of Labor Statistics (OEWS + Employment Projections) data into the single `career_data` table consumed by the backend API.

## Source Inputs
| File | Origin | Purpose |
|------|--------|---------|
| `all_data_M_2023.xlsx` | OEWS (Occupational Employment & Wage Statistics) 2023 data extract | Provides employment counts & wage distribution percentiles by area + occupation |
| `education.xlsx` | BLS Employment Projections (Table 5.4) | Maps occupation codes to typical Education & Experience requirements |

## Output Artifacts
| File | Description |
|------|------------|
| `cleaned_oes_data.csv` | Filtered OEWS data (cross‑industry, detailed occupations + national total) with duplicates removed |
| `combined_career_data.csv` | Final merged dataset (employment + wages + education + experience) used to populate `career_data` |
| `data_schema_example.csv` | Example illustrating target schema layout (reference) |

## High-Level Steps
1. Load raw OEWS workbook, treating symbols (`*`, `**`, `#`) as missing values.
2. Keep only `I_GROUP == 'cross-industry'` rows (cross-industry totals across industries for each occupation).
3. Retain either detailed occupations (`O_GROUP == 'detailed'`) OR the national aggregate row (`OCC_CODE == '00-0000'`).
4. Deduplicate on `(AREA_TITLE, OCC_CODE)` as a safety net (first occurrence wins).
5. Select a minimal column subset & coerce numeric types (nullable Int64 for salary & employment fields).
6. Drop rows missing both employment and median salary.
7. Merge cleaned OEWS data with Employment Projections sheet `Table 5.4` (education & experience) on `OCC_CODE`.
8. Reorder & rename columns, export as `combined_career_data.csv`.
9. Upload / import into Supabase Postgres as `career_data` (schema matches backend expectations).

## Scripts
| Script | Role | Notes |
|--------|------|-------|
| `processData.py` | Cleans and filters raw OEWS workbook → `cleaned_oes_data.csv` | Ensures only detailed + national summary rows; removes duplicates |
| `tableCombinationGenerator.py` | Merges cleaned OEWS data with education & experience from EP `Table 5.4` → `combined_career_data.csv` | Validates expected sheet & columns |
| `diagnose_duplicates.py` | Diagnostic helper for investigating duplicate `(AREA_TITLE, OCC_CODE)` cases | Printed transposed for easy visual diff |

## Column Mapping Into Final Table
| Final Column | Source Column (OEWS / EP) | Meaning |
|--------------|---------------------------|---------|
| AREA_TITLE | `AREA_TITLE` | Geographic area (State / MSA / non-metro / U.S.) |
| OCC_CODE | `OCC_CODE` | Occupation code |
| OCC_TITLE | `OCC_TITLE` | Occupation title |
| Education | `Typical education needed for entry` (EP) | Ladder label for education requirement |
| Experience | `Work experience in a related occupation` (EP) | Ladder label for prior experience |
| TOT_EMP | `TOT_EMP` | Employment count |
| A_MEDIAN | `A_MEDIAN` | Median annual wage |
| A_PCT10 | `A_PCT10` | 10th percentile annual wage |
| A_PCT25 | `A_PCT25` | 25th percentile annual wage |
| A_PCT75 | `A_PCT75` | 75th percentile annual wage |
| A_PCT90 | `A_PCT90` | 90th percentile annual wage |

## Data Quality Choices
| Issue | Decision | Rationale |
|-------|----------|-----------|
| Multiple national rows for same occupation | Keep the one with largest `TOT_EMP` (handled downstream) | Ensures stable national denominator |
| Duplicate `(AREA_TITLE, OCC_CODE)` pairs | Drop subsequent duplicates | Preserve first consistent record after initial filters |
| Missing percentile wages | Retain if employment + median present | Still useful for denominator & threshold logic |
| Non-detailed occupations | Excluded except total `00-0000` | Avoid double-counting aggregated groups |

## Importing Into Supabase
1. Create table `career_data` with columns mirroring final CSV (see root / backend docs).
2. Enforce UNIQUE `(area_title, occ_code)` constraint.
3. Use Supabase web UI or `psql \copy` to load `combined_career_data.csv`.
4. Optionally create indexes:
	- `(occ_title)` for occupation lookups
	- `(area_title)` for area filtering
	- `(education)`, `(experience)` if query plans show benefit

## Re-running End-to-End
```
python processData.py
python tableCombinationGenerator.py
```
Confirm resulting `combined_career_data.csv` row count matches expectations before import.

## Potential Enhancements
- Parameterize year (currently hard-coded to 2023 filenames)
- Add unit tests validating no duplicate `(AREA_TITLE, OCC_CODE)` after processing
- Script to push directly to Supabase via API / `pg_copy` library
- Add hash-based change detection to skip rebuild when inputs unchanged

---
This folder documents how the dataset feeding the application was produced; it is not part of the runtime API deployment.
