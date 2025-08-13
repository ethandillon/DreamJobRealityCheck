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

const ResultDisplay = ({ data }) => (
  // 2. Replace the old div with our AnimatedGradientBorder component
  <AnimatedGradientBorder
    // 3. Move all layout and styling classes here
    className="text-center text-white p-10 rounded-xl shadow-lg backdrop-blur-sm max-w-2xl"
  >
    {/* The inner content remains exactly the same */}
    <p className="font-secondary text-gray-400 text-lg">
      Based on your criteria for jobs in {data.location}:
    </p>
    <div className="font-primary text-7xl font-bold my-4 text-indigo-400">
      {data.percentage.toFixed(2)}%
    </div>
    <p className="font-secondary text-xl">
      An estimated{' '}
      <span className="font-bold">{data.matchingJobs.toLocaleString()}</span> out of{' '}
      <span className="font-bold">{data.totalJobs.toLocaleString()}</span> jobs meet your
      standards.
    </p>
    
    {/* Salary Information */}
    {data.salaryInfo && (
      <div className="mt-6 p-4 bg-gray-800/50 rounded-lg">
        <h3 className="font-bold text-lg mb-3 text-indigo-300">Salary Information</h3>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-400">Median:</span>
            <span className="ml-2 font-bold text-green-400">${data.salaryInfo.medianSalary?.toLocaleString()}</span>
          </div>
          <div>
            <span className="text-gray-400">10th Percentile:</span>
            <span className="ml-2 font-bold text-yellow-400">${data.salaryInfo.pct10Salary?.toLocaleString()}</span>
          </div>
          <div>
            <span className="text-gray-400">25th Percentile:</span>
            <span className="ml-2 font-bold text-blue-400">${data.salaryInfo.pct25Salary?.toLocaleString()}</span>
          </div>
          <div>
            <span className="text-gray-400">75th Percentile:</span>
            <span className="ml-2 font-bold text-purple-400">${data.salaryInfo.pct75Salary?.toLocaleString()}</span>
          </div>
          <div className="col-span-2">
            <span className="text-gray-400">90th Percentile:</span>
            <span className="ml-2 font-bold text-pink-400">${data.salaryInfo.pct90Salary?.toLocaleString()}</span>
          </div>
        </div>
        
        {/* Salary Requirement Status */}
        {data.minSalaryMet !== undefined && (
          <div className="mt-3 pt-3 border-t border-gray-600">
            <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
              data.minSalaryMet 
                ? 'bg-green-900 text-green-300' 
                : 'bg-red-900 text-red-300'
            }`}>
              {data.minSalaryMet ? '✓' : '✗'} 
              {data.minSalaryMet 
                ? ' Salary requirement met' 
                : ' Salary requirement not met'
              }
            </span>
          </div>
        )}
      </div>
    )}
  </AnimatedGradientBorder>
);

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