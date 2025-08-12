// src/components/Filters.jsx

import { useState } from 'react';
import CustomSelect from './CustomSelect'; // The custom dropdown component
import SearchableDropdown from './SearchableDropdown'; // Import the new component
import { mockOccupations } from '../mockOccupations'; // Import our mock data

// Define the options for our dropdowns as constant arrays
const educationOptions = [
  "No formal education",
  "High school diploma",
  "Associate degree",
  "Bachelor's degree",
  "Master's degree",
  "Doctoral or professional degree",
];

const experienceOptions = [
  "None",
  "Less than 2 years",
  "2-4 years",
  "5+ years",
];

// A reusable component for our form rows to keep the code clean
const FilterRow = ({ label, children }) => (
  <div className="mb-6">
    <label className="block text-sm font-medium text-gray-700 mb-2">{label}</label>
    {children}
  </div>
);

function Filters({ onCalculate }) {
  // State hooks to manage the form inputs
  const [location, setLocation] = useState('San Francisco, CA');
  const [occupation, setOccupation] = useState(''); // Add state for occupation
  const [minSalary, setMinSalary] = useState(80000);
  const [education, setEducation] = useState(educationOptions[3]); // Default to "Bachelor's degree"
  const [experience, setExperience] = useState(experienceOptions[2]); // Default to "2-4 years"

  const handleSubmit = (event) => {
    event.preventDefault(); // Prevent full page reload on form submission
    onCalculate({ location, occupation, minSalary, education, experience });
  };

  return (
    <div className="flex flex-col h-full text-gray-800">
      <h1 className="font-primary text-4xl font-bold mb-2">
        Career Opportunity Calculator
      </h1>
      <p className="text-gray-500 mb-8">
        What percentage of jobs in the United States meet your standards?
      </p>

      <form onSubmit={handleSubmit} className="flex-grow">
        <FilterRow label="Occupation / Field">
          {/* Use our new searchable dropdown component */}
          <SearchableDropdown
            options={mockOccupations}
            value={occupation}
            onChange={setOccupation}
            placeholder="e.g., Software Developer"
          />
        </FilterRow>

        <FilterRow label="Location (City, State)">
          <input
            type="text"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            className="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
        </FilterRow>

        <FilterRow label={`Minimum Annual Salary: $${Number(minSalary).toLocaleString()}`}>
          <input
            type="range"
            min="30000"
            max="250000"
            step="5000"
            value={minSalary}
            onChange={(e) => setMinSalary(e.target.value)}
            className="w-full cursor-pointer custom-range" // Apply our custom slider style
          />
        </FilterRow>

        <FilterRow label="Minimum Education Level">
          <CustomSelect
            value={education}
            onChange={setEducation}
            options={educationOptions}
          />
        </FilterRow>

        <FilterRow label="Required Work Experience">
          <CustomSelect
            value={experience}
            onChange={setExperience}
            options={experienceOptions}
          />
        </FilterRow>
      </form>
      
      {/* This div pushes the button to the bottom of the container */}
      <div className="mt-auto pt-4">
        <button
          onClick={handleSubmit}
          className="w-full bg-gray-900 text-white font-bold py-3 px-4 rounded-lg hover:bg-gray-700 transition-colors duration-300 flex items-center justify-between"
        >
          <span>Let's Find Out</span>
          <span>&rarr;</span>
        </button>
        <p className="text-xs text-gray-400 mt-2 text-center">
          Calculated using U.S. Bureau of Labor Statistics Data
        </p>
      </div>
    </div>
  );
}

export default Filters;