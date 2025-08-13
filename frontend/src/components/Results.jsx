// src/components/Results.jsx
import { AnimatedGradientBorder } from './AnimatedGradientBorder'; // 1. Import our new component

// ... (Placeholder and LoadingSpinner components remain the same) ...
const Placeholder = () => (
    <div className="text-center text-gray-500">
      <div className="border border-gray-700 rounded-lg p-4 max-w-sm mx-auto">
        <span role="img" aria-label="chart" className="text-2xl mb-2 block">📊</span>
        Start by specifying your preferences.
      </div>
    </div>
  );
  
const LoadingSpinner = () => (
    <div className="text-center text-gray-400">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto"></div>
      <p className="mt-4">Calculating...</p>
    </div>
  );

import { useState } from 'react';

const ResultDisplay = ({ data }) => {
  const [view, setView] = useState('national'); // 'national' | 'regional'
  const percent = view === 'national' ? data.percentage : (data.percentageRegion ?? data.percentage);
  const percentStr = (() => {
    const num = Number(percent || 0);
    if (!isFinite(num)) return '0.00';
    // Start with 2 decimals; extend until a non-zero digit appears after the decimal
    let decimals = 2;
    let out = num.toFixed(decimals);
    const hasNonZeroAfterDecimal = (s) => {
      const i = s.indexOf('.')
      if (i === -1) return false;
      return /[1-9]/.test(s.slice(i + 1));
    };
    while (decimals < 10 && !hasNonZeroAfterDecimal(out)) {
      decimals += 1;
      out = num.toFixed(decimals);
    }
    // If still zero after 10 places, collapse to a clean 0
    if (!hasNonZeroAfterDecimal(out)) {
      return '0';
    }
    return out;
  })();
  const denom = view === 'national' ? data.totalJobs : (data.totalJobsRegion ?? data.totalJobs);
  const toggle = () => setView(v => (v === 'national' ? 'regional' : 'national'));

  return (
  // 2. Replace the old div with our AnimatedGradientBorder component
  <AnimatedGradientBorder
    // 3. Move all layout and styling classes here
    className="relative text-center text-white p-10 rounded-xl shadow-lg backdrop-blur-sm max-w-2xl"
  >
    {/* Subtle toggle button in the top-right */}
    <button
      onClick={toggle}
      className="absolute top-3 right-3 p-0 m-0 bg-transparent cursor-pointer"
      aria-label="Toggle national/regional view"
      title={view === 'national' ? 'Switch to regional view' : 'Switch to national view'}
    >
      <img src="/tab%20icon.png" alt="toggle view" className="h-4 w-4 opacity-80 hover:opacity-100" />
    </button>
    <div className="font-primary text-7xl font-bold my-4 text-indigo-400">
      {percentStr}%
    </div>
    <p className="font-secondary text-xl">
      An estimated{' '}
      <span className="font-bold">{data.matchingJobs.toLocaleString()}</span> out of{' '}
      <span className="font-bold">{Number(denom || 0).toLocaleString()}</span> jobs in {view === 'national' ? 'the U.S.' : data.location} meet your standards.
    </p>
    
    {/* Salary Information removed by request */}
  </AnimatedGradientBorder>
  );
};

const ErrorDisplay = ({ message }) => (
    <div className="text-center text-red-400 border border-red-400/50 p-6 rounded-lg">
      <h3 className="font-bold text-lg mb-2">An Error Occurred</h3>
      <p>{message}</p>
    </div>
  );

function Results({ isLoading, data, error }) {
  if (isLoading) {
    return <LoadingSpinner />;
  }
  if (error) {
    return <ErrorDisplay message={error} />;
  }
  if (data) {
    return <ResultDisplay data={data} />;
  }
  return <Placeholder />;
}

export default Results;