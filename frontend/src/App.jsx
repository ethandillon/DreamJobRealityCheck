// src/App.jsx
import { useState } from 'react';
import Filters from './components/Filters';
import Results from './components/Results';

function App() {
  // State to hold the final result from the (simulated) API
  const [resultData, setResultData] = useState(null);
  // State to show a loading message while we "fetch"
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // This function fetches data from your Go backend
  const handleCalculate = async (filters) => {
    setIsLoading(true);
    setResultData(null);
    setError(null);

    try {
      // Build query parameters
      const queryParams = new URLSearchParams();
      if (filters.location) queryParams.append('location', filters.location);
      if (filters.occupation) queryParams.append('occupation', filters.occupation);
      if (filters.minSalary) queryParams.append('minSalary', filters.minSalary);
      if (filters.education) queryParams.append('education', filters.education);
      if (filters.experience) queryParams.append('experience', filters.experience);

      // Fetch from your Go API
      const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${apiBase}/api/calculate?${queryParams}`);
      
      if (!response.ok) {
        throw new Error(`API request failed: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      setResultData(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex flex-col md:flex-row min-h-screen bg-gray-900 font-secondary text-white">
      {/* Left Panel: Filters */}
      <div className="w-full md:max-w-md bg-white p-8 shadow-2xl">
        <Filters onCalculate={handleCalculate} />
      </div>

      {/* Right Panel: Results */}
      <div className="flex w-full md:flex-grow items-center justify-center p-8">
        <Results isLoading={isLoading} data={resultData} error={error} />
      </div>
    </main>
  );
}

export default App;