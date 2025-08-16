// src/App.jsx
import { useState } from 'react';
import Filters from './components/Filters';
import Results from './components/Results';
import { calculate } from './api/client';

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
      const data = await calculate(filters);
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