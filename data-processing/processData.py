import pandas as pd
import os

# --- Configuration ---
input_filename = 'all_data_M_2023.xlsx' 
output_filename = 'cleaned_oes_data.csv' 

# --- Main Script ---
print(f"Looking for input file: {input_filename}")

if not os.path.exists(input_filename):
    print(f"Error: The file '{input_filename}' was not found in this directory.")
else:
    print("File found. Starting data processing...")
    missing_value_formats = ['*', '**', '#']
    df = pd.read_excel(input_filename, na_values=missing_value_formats)
    print(f"Initial row count: {len(df)}")

    # Level 1 Filter: Get cross-industry totals
    df = df[df['I_GROUP'] == 'cross-industry']
    print(f"Row count after 'cross-industry' filter: {len(df)}")
    
    # --- THIS IS THE CRITICAL CHANGE ---
    # We now keep EITHER the detailed occupations OR the single 'Total' summary row.
    # The '|' symbol means OR in pandas.
    df = df[(df['O_GROUP'] == 'detailed') | (df['OCC_CODE'] == '00-0000')]
    print(f"Row count after 'detailed' OR 'Total' filter: {len(df)}")
    # --- END OF CHANGE ---

    # Safety Net: Programmatically remove any remaining duplicates
    initial_rows = len(df)
    df = df.drop_duplicates(subset=['AREA_TITLE', 'OCC_CODE'], keep='first')
    final_rows = len(df)
    print(f"Dropped {initial_rows - final_rows} additional duplicate row(s) with the final safety net.")

    columns_to_keep = [
        'AREA_TITLE',
        'OCC_CODE',
        'OCC_TITLE',
        'TOT_EMP',
        'A_PCT10',
        'A_PCT25',
        'A_MEDIAN',
        'A_PCT75',
        'A_PCT90'
    ]
    df_cleaned = df[columns_to_keep].copy()

    print("Now cleaning and converting data types...")
    numeric_cols = ['TOT_EMP', 'A_PCT10', 'A_PCT25', 'A_MEDIAN', 'A_PCT75', 'A_PCT90']
    for col in numeric_cols:
        df_cleaned[col] = pd.to_numeric(df_cleaned[col], errors='coerce')
        df_cleaned[col] = df_cleaned[col].astype('Int64')

    # We need to adjust the dropna logic slightly, because the 'Total' row has no salary percentiles.
    # We will only drop rows if they are missing BOTH employment and median salary.
    df_cleaned = df_cleaned.dropna(subset=['TOT_EMP', 'A_MEDIAN'])

    print("Data cleaned. Saving to a new, faster file...")
    df_cleaned.to_csv(output_filename, index=False)

    print(f"Success! A unique, aggregate version of the data has been saved as '{output_filename}'")