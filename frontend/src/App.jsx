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
      const response = await fetch(`http://localhost:8080/api/calculate?${queryParams}`);
      
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
    <main className="flex min-h-screen bg-gray-900 font-secondary text-white">
      {/* Left Panel: Filters */}
      <div className="w-full max-w-md bg-white p-8 shadow-2xl">
        <Filters onCalculate={handleCalculate} />
      </div>

      {/* Right Panel: Results */}
      <div className="flex-grow flex items-center justify-center p-8">
        <Results isLoading={isLoading} data={resultData} error={error} />
      </div>
    </main>
  );
}

export default App;