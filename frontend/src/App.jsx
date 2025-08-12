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

  // This function simulates fetching data from your Go backend
  const handleCalculate = async (filters) => {
    setIsLoading(true);
    setResultData(null);
    setError(null);

    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 1500));

    // In a real app, you would fetch from your Go API here:
    // const queryParams = new URLSearchParams(filters).toString();
    // const response = await fetch(`http://localhost:8080/calculate?${queryParams}`);
    // const data = await response.json();
    
    // For now, we'll return some mock data
    try {
      if (filters.location.toLowerCase() === 'error') {
        throw new Error("Invalid location provided.");
      }
      const mockData = {
        percentage: Math.random() * 25, // Random percentage between 0 and 25
        matchingJobs: Math.floor(Math.random() * 5000) + 1000,
        totalJobs: 250000,
        location: filters.location || "the selected area",
      };
      setResultData(mockData);
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