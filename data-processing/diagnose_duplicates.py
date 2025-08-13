import pandas as pd

# Load the original, raw data file
df = pd.read_excel('all_data_M_2023.xlsx')

# Apply the first filter we already discovered
df_filtered = df[df['I_GROUP'] == 'cross-industry']

# Now, let's isolate the exact key that is causing the error
problem_rows = df_filtered[
    (df_filtered['AREA_TITLE'] == 'U.S.') &
    (df_filtered['OCC_CODE'] == '13-1020')
]

# Print out the problematic rows and look for a column that is different
# We print it transposed (.T) to make it easier to compare all columns
print("Found the following duplicate rows for the key (U.S., 13-1020):")
print(problem_rows.T)