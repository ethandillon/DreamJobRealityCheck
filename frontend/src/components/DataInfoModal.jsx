import { useEffect } from 'react';

// Accessible modal / popover for explaining data source & methodology
function DataInfoModal({ open, onClose }) {
  // Close on ESC
  useEffect(() => {
    if (!open) return;
    const handler = (e) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="data-info-overlay" role="dialog" aria-modal="true" aria-labelledby="data-info-title">
      <div className="data-info-modal">
        <div className="flex justify-between items-start mb-4">
          <h2 id="data-info-title" className="text-xl font-semibold">About the Data & Calculation</h2>
          <button aria-label="Close" onClick={onClose} className="text-gray-500 hover:text-gray-800 transition">✕</button>
        </div>
        <div className="text-sm leading-relaxed max-h-[60vh] overflow-y-auto pr-1">
          <p>
            We use the latest official data from the U.S. Bureau of Labor Statistics (May 2023). We calculate the percentage of jobs that match your standards for location, title, and the typical education or experience needed for that career. Since pay can vary, we consider a job a match if its median or higher-end salary meets your minimum. The percentage is always specific to your search, comparing your results against all the jobs in the area you select. A crucial point to remember is that this data doesn't represent every single job in America. The government's survey focuses specifically on jobs on company payrolls, so it doesn't include the millions of people who are self-employed, work in the gig economy, or own small businesses.
          </p>
        </div>
        <div className="mt-6 text-right">
          <button onClick={onClose} className="px-4 py-2 bg-gray-900 text-white rounded-md text-sm hover:bg-gray-700 transition">Close</button>
        </div>
      </div>
    </div>
  );
}

export default DataInfoModal;
