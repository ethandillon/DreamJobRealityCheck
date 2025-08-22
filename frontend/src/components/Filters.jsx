// src/components/Filters.jsx

import { useState, useEffect } from 'react';
import { getOccupations, getStates, getAreasByState } from '../api/client';
import CustomSelect from './CustomSelect'; // The custom dropdown component
import SearchableDropdown from './SearchableDropdown'; // Import the new component
import DataInfoModal from './DataInfoModal';

// Define the options for our dropdowns as constant arrays (added "Any")
const educationOptions = [
  "Any",
  "No formal education",
  "High school diploma",
  "Postsecondary nondegree award",
  "Associate degree",
  "Bachelor's degree",
  "Master's degree",
  "Doctoral or professional degree",
];

// Match DB strings; include "Any" which removes experience filter
const experienceOptions = [
  "Any",
  "None",
  "Less than 5 years",
  "5 years or more",
];

// Minimal USPS abbreviation to state name map for deriving state from shared URLs
const ABR_TO_STATE = {
  AL: 'Alabama', AK: 'Alaska', AZ: 'Arizona', AR: 'Arkansas', CA: 'California', CO: 'Colorado', CT: 'Connecticut',
  DE: 'Delaware', DC: 'District of Columbia', FL: 'Florida', GA: 'Georgia', HI: 'Hawaii', ID: 'Idaho', IL: 'Illinois',
  IN: 'Indiana', IA: 'Iowa', KS: 'Kansas', KY: 'Kentucky', LA: 'Louisiana', ME: 'Maine', MD: 'Maryland', MA: 'Massachusetts',
  MI: 'Michigan', MN: 'Minnesota', MS: 'Mississippi', MO: 'Missouri', MT: 'Montana', NE: 'Nebraska', NV: 'Nevada',
  NH: 'New Hampshire', NJ: 'New Jersey', NM: 'New Mexico', NY: 'New York', NC: 'North Carolina', ND: 'North Dakota',
  OH: 'Ohio', OK: 'Oklahoma', OR: 'Oregon', PA: 'Pennsylvania', RI: 'Rhode Island', SC: 'South Carolina', SD: 'South Dakota',
  TN: 'Tennessee', TX: 'Texas', UT: 'Utah', VT: 'Vermont', VA: 'Virginia', WA: 'Washington', WV: 'West Virginia',
  WI: 'Wisconsin', WY: 'Wyoming', PR: 'Puerto Rico', GU: 'Guam', VI: 'Virgin Islands'
};

// A reusable component for our form rows to keep the code clean
const FilterRow = ({ label, children }) => (
  <div className="mb-6">
    <label className="block text-sm font-medium text-gray-700 mb-2">{label}</label>
    {children}
  </div>
);

function Filters({ onCalculate, initialValues }) {
  // State hooks to manage the form inputs
  const [location, setLocation] = useState(initialValues?.location || '');
  const [selectedState, setSelectedState] = useState('');
  const [occupation, setOccupation] = useState(initialValues?.occupation || ''); // Add state for occupation
  const [minSalary, setMinSalary] = useState(initialValues?.minSalary ?? 80000);
  const [education, setEducation] = useState(initialValues?.education || educationOptions[0]); // Default to "Any"
  const [experience, setExperience] = useState(initialValues?.experience || experienceOptions[0]); // Default to "Any"
  const [occupations, setOccupations] = useState([]); // State for real occupation data
  const [isLoadingOccupations, setIsLoadingOccupations] = useState(true); // Loading state
  const [states, setStates] = useState([]); // State list for first dropdown
  const [isLoadingStates, setIsLoadingStates] = useState(true);
  const [areas, setAreas] = useState([]); // Areas for the selected state
  const [isLoadingAreas, setIsLoadingAreas] = useState(false);
  const [showDataInfo, setShowDataInfo] = useState(false);

  // Fetch occupations from the backend API
  useEffect(() => {
    const fetchOccupations = async () => {
      try {
        const data = await getOccupations();
        setOccupations(data.occupations || []);
      } catch (error) {
        console.error('Error fetching occupations:', error);
        setOccupations([]);
      } finally {
        setIsLoadingOccupations(false);
      }
    };
    fetchOccupations();
  }, []);

  // Fetch states on mount
  useEffect(() => {
    const fetchStates = async () => {
      try {
        const data = await getStates();
        setStates(data.states || []);
      } catch (error) {
        console.error('Error fetching states:', error);
        setStates([]);
      } finally {
        setIsLoadingStates(false);
      }
    };
    fetchStates();
  }, []);

  // Fetch areas when a state is selected
  useEffect(() => {
    if (!selectedState) {
      setAreas([]);
      setLocation('');
      return;
    }
    const controller = new AbortController();
    const fetchAreas = async () => {
      setIsLoadingAreas(true);
      try {
        const data = await getAreasByState(selectedState, { signal: controller.signal });
        setAreas(data.areas || []);
      } catch (error) {
        if (error.name !== 'AbortError') {
          console.error('Error fetching areas:', error);
        }
        setAreas([]);
      } finally {
        setIsLoadingAreas(false);
      }
    };
    fetchAreas();
    return () => controller.abort();
  }, [selectedState]);

  // Try to derive selectedState from location (on mount / when initialValues change)
  useEffect(() => {
    if (!initialValues?.location || selectedState) return;
    const loc = initialValues.location;
    // Heuristic: if location contains ", XY" use that; else if it equals a state name, set it directly
    const stateFromComma = /,\s*([A-Z]{2})(?:\b|-)/.exec(loc)?.[1];
    if (stateFromComma && ABR_TO_STATE[stateFromComma]) {
      setSelectedState(ABR_TO_STATE[stateFromComma]);
      return;
    }
    // If the location looks like a full state name (no comma), set it
    if (!loc.includes(',')) {
      setSelectedState(loc);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialValues]);

  const handleSubmit = (event) => {
    event.preventDefault(); // Prevent full page reload on form submission
    const payload = { location, occupation, minSalary };
    if (education && education !== 'Any') payload.education = education;
    // Only include experience if not "Any"
    if (experience && experience !== 'Any') payload.experience = experience;
    if (!isFormValid) return;
    onCalculate(payload);
  };

  const isFormValid = Boolean(selectedState && location && occupation);

  return (
    <div className="flex flex-col h-full text-gray-800">
      <h1 className="font-primary text-5xl font-bold mb-2">
        Dream Job Reality Check
      </h1>
      <p className="text-gray-500 mb-8">
        What percentage of jobs in the United States meet your standards?
      </p>

      <form onSubmit={handleSubmit} className="flex-grow">
        <FilterRow label="Occupation / Field">
          {/* Use our new searchable dropdown component */}
          <SearchableDropdown
            options={occupations}
            value={occupation}
            onChange={setOccupation}
            placeholder={isLoadingOccupations ? "Loading occupations..." : "e.g., Software Developer"}
            disabled={isLoadingOccupations}
          />
        </FilterRow>

        <FilterRow label="State">
          <SearchableDropdown
            options={states}
            value={selectedState}
            onChange={(val) => { setSelectedState(val); setLocation(''); }}
            placeholder={isLoadingStates ? 'Loading states...' : 'Select a state'}
            disabled={isLoadingStates}
          />
        </FilterRow>

        <FilterRow label="Area within State">
          <SearchableDropdown
            options={areas}
            value={location}
            onChange={setLocation}
            placeholder={!selectedState ? 'Select a state first' : (isLoadingAreas ? 'Loading areas...' : 'Choose statewide, MSA, or non-metro area')}
            disabled={!selectedState || isLoadingAreas}
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
          disabled={!isFormValid}
          aria-disabled={!isFormValid}
          className={`w-full bg-gray-900 text-white font-bold py-3 px-4 rounded-lg transition-colors duration-300 flex items-center justify-between ${
            isFormValid ? 'hover:bg-gray-700 cursor-pointer' : 'opacity-50 cursor-not-allowed'
          }`}
        >
          <span>Let's Find Out</span>
          <span>&rarr;</span>
        </button>
        <p className="text-xs text-gray-400 mt-2 text-center">
          <button type="button" onClick={() => setShowDataInfo(true)} className="underline hover:text-gray-300 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 rounded">
            Calculated using U.S. Bureau of Labor Statistics Data
          </button>
        </p>
      </div>
      <DataInfoModal open={showDataInfo} onClose={() => setShowDataInfo(false)} />
    </div>
  );
}

export default Filters;