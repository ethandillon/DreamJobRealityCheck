// src/App.jsx
import { useEffect, useMemo, useState } from 'react';
import Filters from './components/Filters';
import Results from './components/Results';
import { calculate } from './api/client';

function App() {
  // State to hold the final result from the (simulated) API
  const [resultData, setResultData] = useState(null);
  // State to show a loading message while we "fetch"
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [lastFilters, setLastFilters] = useState(null);

  // Parse initial filters from URL on mount
  const initialFilters = useMemo(() => {
    const params = new URLSearchParams(window.location.search);
    const f = {
      location: params.get('location') || '',
      occupation: params.get('occupation') || '',
      minSalary: params.get('minSalary') ? Number(params.get('minSalary')) : 80000,
      education: params.get('education') || 'Any',
      experience: params.get('experience') || 'Any',
    };
    return f;
  }, []);

  // Helper to build shareable query string
  const filtersToQuery = (filters) => {
    const p = new URLSearchParams();
    Object.entries(filters || {}).forEach(([k, v]) => {
      if (v === undefined || v === null || v === '') return;
      p.set(k, String(v));
    });
    return p.toString();
  };

  // This function fetches data from your Go backend
  const handleCalculate = async (filters) => {
    setIsLoading(true);
    setResultData(null);
    setError(null);
    setLastFilters(filters);
    // Sync URL for shareability
    const qs = filtersToQuery(filters);
    const nextUrl = `${window.location.pathname}${qs ? `?${qs}` : ''}`;
    window.history.replaceState(null, '', nextUrl);

    try {
      const data = await calculate(filters);
      setResultData(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setIsLoading(false);
    }
  };

  // Auto-run calculation on load if URL has enough info (at least location and occupation)
  useEffect(() => {
    if (initialFilters.location && initialFilters.occupation) {
      handleCalculate(initialFilters);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <main className="flex flex-col md:flex-row min-h-screen bg-gray-900 font-secondary text-white">
      {/* Left Panel: Filters */}
      <div className="w-full md:max-w-md bg-white p-8 shadow-2xl">
        <Filters onCalculate={handleCalculate} initialValues={initialFilters} />
      </div>

      {/* Right Panel: Results */}
      <div className="flex w-full md:flex-grow items-center justify-center p-8">
        <Results isLoading={isLoading} data={resultData} error={error} shareFilters={lastFilters} />
      </div>
    </main>
  );
}

export default App;