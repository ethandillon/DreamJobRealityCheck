import pandas as pd
import os

# --- Configuration ---
# The name of the file you downloaded from the BLS website
input_filename = 'all_data_M_2023.xlsx' 
# The name of the clean file we will create
output_filename = 'cleaned_oes_data.csv' 

# --- Main Script ---
print(f"Looking for input file: {input_filename}")

if not os.path.exists(input_filename):
    print(f"Error: The file '{input_filename}' was not found in this directory.")
    print("Please download it from the BLS OEWS website and place it here.")
else:
    print("File found. Starting data processing... (This may take a minute)")

    # Define the special characters the BLS uses for missing/invalid data
    # We will tell pandas to treat all of these as "Not a Number" (NaN)
    missing_value_formats = ['*', '**', '#']

    # Read the excel file, explicitly telling it what to consider as missing data
    # This is the most important step!
    df = pd.read_excel(input_filename, na_values=missing_value_formats)

    print("Initial data loaded. Now cleaning and selecting columns...")

    # Select only the columns we need for the project to save memory and make it easier to work with
    columns_to_keep = [
        'AREA_TITLE',  # Location name (e.g., "California")
        'OCC_CODE',    # Occupation code (e.g., "15-1252")
        'OCC_TITLE',   # Occupation name (e.g., "Software Developers")
        'TOT_EMP',     # Total employment in that job/area
        'A_PCT10',     # 10th percentile annual salary
        'A_PCT25',     # 25th percentile annual salary
        'A_MEDIAN',    # Median annual salary
        'A_PCT75',     # 75th percentile annual salary
        'A_PCT90'      # 90th percentile annual salary
    ]
    
    df_cleaned = df[columns_to_keep]

    # The columns with numbers were loaded as 'object' type because of the special characters.
    # Now that we've handled the special characters, we can convert them to a numeric type.
    # `errors='coerce'` will turn any remaining non-numeric values into NaN.
    numeric_cols = ['TOT_EMP', 'A_PCT10', 'A_PCT25', 'A_MEDIAN', 'A_PCT75', 'A_PCT90']
    for col in numeric_cols:
        df_cleaned[col] = pd.to_numeric(df_cleaned[col], errors='coerce')

    # For simplicity, we'll remove rows where key data (like employment or median salary) is missing.
    df_cleaned = df_cleaned.dropna(subset=['TOT_EMP', 'A_MEDIAN'])

    print("Data cleaned. Saving to a new, faster file...")

    # Save the cleaned data to a CSV file. CSV files are much faster to load next time.
    df_cleaned.to_csv(output_filename, index=False)

    print(f"Success! A clean version of the data has been saved as '{output_filename}'")
    print("\n--- First 5 rows of the clean data: ---")
    print(df_cleaned.head())
    print("\n--- Data types of the clean data: ---")
    print(df_cleaned.info())