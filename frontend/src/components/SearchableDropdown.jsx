// src/components/SearchableDropdown.jsx
import { useState, useRef, useEffect } from 'react';

function SearchableDropdown({ options, value, onChange, placeholder, disabled = false }) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [highlightedIndex, setHighlightedIndex] = useState(-1); // For keyboard navigation
  const dropdownRef = useRef(null);
  const listRef = useRef(null); // list container for scroll management

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
    setHighlightedIndex(-1);
  };

  const filteredOptions = options.filter(option =>
    option.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Reset highlight when search term changes
  useEffect(() => {
    if (searchTerm === "" && !isOpen) return; // don't force highlight when closed
    if (filteredOptions.length > 0) {
      setHighlightedIndex(0);
    } else {
      setHighlightedIndex(-1);
    }
  }, [searchTerm, filteredOptions.length, isOpen]);

  const moveHighlight = (direction) => {
    if (filteredOptions.length === 0) return;
    setHighlightedIndex((prev) => {
      if (prev === -1) return 0;
      const next = (prev + direction + filteredOptions.length) % filteredOptions.length;
      return next;
    });
  };

  // Ensure highlighted option is visible within scroll viewport
  useEffect(() => {
    if (!isOpen) return;
    if (highlightedIndex < 0) return;
    const listEl = listRef.current;
    if (!listEl) return;
    const optionEl = listEl.querySelector(`[data-index='${highlightedIndex}']`);
    if (!optionEl) return;
    const optionTop = optionEl.offsetTop;
    const optionBottom = optionTop + optionEl.offsetHeight;
    const viewTop = listEl.scrollTop;
    const viewBottom = viewTop + listEl.clientHeight;
    if (optionTop < viewTop) {
      listEl.scrollTop = optionTop; // scroll up
    } else if (optionBottom > viewBottom) {
      listEl.scrollTop = optionBottom - listEl.clientHeight; // scroll down
    }
  }, [highlightedIndex, isOpen]);

  const handleKeyDown = (e) => {
    if (disabled) return;
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        if (!isOpen) setIsOpen(true);
        moveHighlight(1);
        break;
      case 'ArrowUp':
        e.preventDefault();
        if (!isOpen) setIsOpen(true);
        moveHighlight(-1);
        break;
      case 'Tab':
        if (isOpen && filteredOptions.length > 0) {
          e.preventDefault();
          moveHighlight(e.shiftKey ? -1 : 1);
        }
        break;
      case 'Enter':
        if (isOpen && highlightedIndex >= 0) {
          e.preventDefault();
            handleSelect(filteredOptions[highlightedIndex]);
        }
        break;
      case 'Escape':
        if (isOpen) {
          e.preventDefault();
          setIsOpen(false);
        }
        break;
      default:
        break;
    }
  };

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
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        aria-autocomplete="list"
        aria-expanded={isOpen}
        aria-haspopup="listbox"
        className={`w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 ${
          disabled ? 'bg-gray-100 cursor-not-allowed' : ''
        }`}
      />
      {isOpen && (
  <ul
          role="listbox"
    className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto"
    ref={listRef}
        >
          {filteredOptions.length > 0 ? (
            filteredOptions.map((option, idx) => {
              const highlighted = idx === highlightedIndex;
              return (
                <li
                  key={option}
                  role="option"
      data-index={idx}
                  aria-selected={highlighted}
                  onMouseEnter={() => setHighlightedIndex(idx)}
                  onMouseDown={(e) => e.preventDefault()} // Prevent input blur before click
                  onClick={() => handleSelect(option)}
                  className={`px-4 py-2 cursor-pointer hover:bg-indigo-100 ${highlighted ? 'bg-indigo-600 text-white hover:bg-indigo-600' : ''}`}
                >
                  {option}
                </li>
              );
            })
          ) : (
            <li className="px-4 py-2 text-gray-500">No results found</li>
          )}
        </ul>
      )}
    </div>
  );
}

export default SearchableDropdown;