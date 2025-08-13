// src/components/SearchableDropdown.jsx
import { useState, useRef, useEffect } from 'react';

function SearchableDropdown({ options, value, onChange, placeholder, disabled = false }) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const dropdownRef = useRef(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelect = (option) => {
    onChange(option); // Update parent state
    setSearchTerm(""); // Clear search term
    setIsOpen(false);   // Close dropdown
  };

  const filteredOptions = options.filter(option =>
    option.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="relative" ref={dropdownRef}>
      <input
        type="text"
        value={searchTerm || value}
        disabled={disabled}
        onFocus={() => {
            if (!disabled) {
                setIsOpen(true);
                setSearchTerm(""); // Clear search term when focusing to search again
            }
        }}
        onChange={(e) => {
            if (!disabled) {
                setIsOpen(true);
                setSearchTerm(e.target.value);
                if (e.target.value === "") {
                    onChange(""); // Clear selection in parent if input is cleared
                }
            }
        }}
        placeholder={placeholder}
        className={`w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 ${
          disabled ? 'bg-gray-100 cursor-not-allowed' : ''
        }`}
      />
      {isOpen && (
        <ul className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto">
          {filteredOptions.length > 0 ? (
            filteredOptions.map((option) => (
              <li
                key={option}
                onClick={() => handleSelect(option)}
                className="px-4 py-2 cursor-pointer hover:bg-indigo-100"
              >
                {option}
              </li>
            ))
          ) : (
            <li className="px-4 py-2 text-gray-500">No results found</li>
          )}
        </ul>
      )}
    </div>
  );
}

export default SearchableDropdown;