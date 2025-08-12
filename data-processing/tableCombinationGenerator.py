import pandas as pd
import os

# --- Configuration ---
# Input files (these must be in the same directory as the script)
oes_data_file = 'cleaned_oes_data.csv'
education_data_file = 'education.xlsx'

# Output file (the final, combined dataset)
output_file = 'combined_career_data.csv'

# --- Main Script ---

# 1. Check if the required files exist before we start
if not os.path.exists(oes_data_file):
    print(f"Error: The file '{oes_data_file}' was not found.")
    print("Please run the previous script ('process_data.py') to generate this file first.")
    exit()

if not os.path.exists(education_data_file):
    print(f"Error: The file '{education_data_file}' was not found.")
    print("Please download it from the BLS Employment Projections page and place it here.")
    exit()

print("All required files found. Starting the merge process...")

# 2. Load the cleaned OEWS job market data
print(f"Loading job market data from '{oes_data_file}'...")
df_oes = pd.read_csv(oes_data_file)
# Ensure OCC_CODE is treated as a string to prevent any merge issues
df_oes['OCC_CODE'] = df_oes['OCC_CODE'].astype(str)

# 3. Load the EP education and experience data from the CORRECT SHEET
print(f"Loading education data from '{education_data_file}', targeting sheet 'Table 5.4'...")

# --- THIS IS THE CRITICAL CHANGE ---
# We specify the sheet_name and the correct number of rows to skip.
try:
    df_education = pd.read_excel(
        education_data_file,
        sheet_name='Table 5.4', # Specify the exact sheet name
        skiprows=1             # Skip the top two title rows to get to the header
    )
except ValueError as e:
    print(f"Error: Could not find sheet 'Table 5.4' in '{education_data_file}'.")
    print(f"Please ensure the file is correct and the sheet is named exactly 'Table 5.4'. Original error: {e}")
    exit()
# --- END OF CHANGE ---

print("Successfully loaded data from 'Table 5.4'.")

# 4. Clean up and select columns from the education dataframe
# The column names from the BLS file for 2023 need to be exact.
# These names are taken from the headers on row 3 of the specified sheet.
ep_columns_to_keep = {
    '2023 National Employment Matrix code': 'OCC_CODE',
    'Typical education needed for entry': 'Education',
    'Work experience in a related occupation': 'Experience'
}

# Check if required columns exist in the loaded dataframe
if not all(col in df_education.columns for col in ep_columns_to_keep.keys()):
    print("Error: The column names in the Excel file do not match the expected names.")
    print("Expected columns:", list(ep_columns_to_keep.keys()))
    print("Found columns:", df_education.columns.tolist())
    exit()

df_education_cleaned = df_education[list(ep_columns_to_keep.keys())]
df_education_cleaned = df_education_cleaned.rename(columns=ep_columns_to_keep)

# Ensure the merge key (OCC_CODE) is a string here as well
df_education_cleaned['OCC_CODE'] = df_education_cleaned['OCC_CODE'].astype(str)

# Drop any potential duplicates from the education file, keeping the first instance
df_education_cleaned = df_education_cleaned.drop_duplicates(subset=['OCC_CODE'])

print("Education data cleaned and prepared for merging.")

# 5. Perform the merge
print(f"Merging {len(df_oes)} job records with {len(df_education_cleaned)} education records...")
df_combined = pd.merge(
    left=df_oes,
    right=df_education_cleaned,
    on='OCC_CODE',
    how='left'
)

# 6. Post-merge verification and saving
print("Merge complete. Verifying the results...")

# Reorder columns for better readability
final_columns_order = [
    'AREA_TITLE',
    'OCC_CODE',
    'OCC_TITLE',
    'Education',
    'Experience',
    'TOT_EMP',
    'A_MEDIAN',
    'A_PCT10',
    'A_PCT25',
    'A_PCT75',
    'A_PCT90'
]
df_combined = df_combined[final_columns_order]

# Save the final, combined dataframe to a new CSV
df_combined.to_csv(output_file, index=False)

print(f"\nSuccess! The combined data has been saved to '{output_file}'")
print("\n--- First 5 rows of the final combined data: ---")
print(df_combined.head())