'use client';

import { useState, useEffect, useCallback } from 'react';
import { useSearchParams, useRouter, usePathname } from 'next/navigation';
import { SearchFilters, SearchFiltersType } from './SearchFilters';
import { LoadingSkeleton } from '../ui/LoadingSkeleton';
import { ErrorDisplay } from '../ui/ErrorDisplay';
import { EmptyState } from '../ui/EmptyState';
import type { SearchSuggestion } from '../../hooks/useSearch';

interface AdvancedSearchProps {
  onSearch: (query: string, filters: SearchFiltersType) => void;
  loading?: boolean;
  error?: string | null;
  suggestions?: SearchSuggestion[];
  placeholder?: string;
}

export function AdvancedSearch({ 
  onSearch, 
  loading = false, 
  error = null, 
  suggestions = [],
  placeholder = "Search recipes, workouts, ingredients..."
}: AdvancedSearchProps) {
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();
  const [query, setQuery] = useState(searchParams?.get('q') || '');
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedSuggestionIndex, setSelectedSuggestionIndex] = useState(-1);
  const [searchHistory, setSearchHistory] = useState<string[]>([]);

  const [filters, setFilters] = useState({
    category: '',
    dietary: [],
    difficulty: '',
    timeRange: { min: 0, max: 120 },
    calories: { min: 0, max: 2000 },
    isHalal: false
  });

  // Load search history from localStorage
  useEffect(() => {
    const history = localStorage.getItem('searchHistory');
    if (history) {
      setSearchHistory(JSON.parse(history));
    }
  }, []);

  // Debounced search function
  const debouncedSearch = useCallback(
    (() => {
      let timeout: NodeJS.Timeout;
      return (searchQuery: string) => {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
          const params = new URLSearchParams();
          if (searchQuery) params.set('q', searchQuery);
          
          // Add filters to URL
          if (filters.category) params.set('category', filters.category);
          if (filters.dietary.length > 0) params.set('dietary', filters.dietary.join(','));
          if (filters.difficulty) params.set('difficulty', filters.difficulty);
          if (filters.isHalal) params.set('halal', 'true');
          
          router.push(`${pathname}?${params.toString()}`);
          onSearch(searchQuery, filters);
        }, 300);
      };
    })(),
    [onSearch, filters, router, pathname]
  );

  // Update search when query or filters change
  useEffect(() => {
    debouncedSearch(query);
  }, [query, filters, debouncedSearch]);

  // Handle search input
  const handleInputChange = (value: string) => {
    setQuery(value);
    setShowSuggestions(value.length > 0);
    setSelectedSuggestionIndex(-1);
  };

  // Handle search submission
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    performSearch();
  };

  const performSearch = () => {
    if (query.trim()) {
      // Add to search history
      const newHistory = [query, ...searchHistory.filter(h => h !== query)].slice(0, 10);
      setSearchHistory(newHistory);
      localStorage.setItem('searchHistory', JSON.stringify(newHistory));
    }
    setShowSuggestions(false);
    onSearch(query, filters);
  };

  // Handle suggestion selection
  const handleSuggestionSelect = (suggestion: SearchSuggestion) => {
    setQuery(suggestion.text);
    setShowSuggestions(false);
    performSearch();
  };

  // Keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedSuggestionIndex(prev => 
          prev < suggestions.length - 1 ? prev + 1 : 0
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedSuggestionIndex(prev => 
          prev > 0 ? prev - 1 : suggestions.length - 1
        );
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedSuggestionIndex >= 0) {
          handleSuggestionSelect(suggestions[selectedSuggestionIndex]);
        } else {
          performSearch();
        }
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedSuggestionIndex(-1);
        break;
    }
  };

  const clearSearch = () => {
    setQuery('');
    setShowSuggestions(false);
    router.push(pathname);
  };

  const applyFilters = (newFilters: SearchFiltersType) => {
    setFilters(newFilters as typeof filters);
  };

  return (
    <div className="w-full max-w-4xl mx-auto p-4">
      {/* Search Input Section */}
      <div className="relative mb-6">
        <form onSubmit={handleSubmit} className="relative">
          <div className="relative flex items-center">
            <div className="relative flex-1">
              {/* Search Icon */}
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>

              {/* Search Input */}
              <input
                type="text"
                value={query}
                onChange={(e) => handleInputChange(e.target.value)}
                onKeyDown={handleKeyDown}
                onFocus={() => setShowSuggestions(query.length > 0)}
                onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
                placeholder={placeholder}
                className="block w-full pl-10 pr-12 py-3 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 text-gray-900"
              />

              {/* Clear Button */}
              {query && (
                <button
                  type="button"
                  onClick={clearSearch}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <svg className="h-5 w-5 text-gray-400 hover:text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>

            {/* Search Button */}
            <button
              type="submit"
              disabled={loading}
              className="ml-3 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Searching...
                </div>
              ) : (
                'Search'
              )}
            </button>
          </div>

          {/* Suggestions Dropdown */}
          {showSuggestions && suggestions.length > 0 && (
            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-auto">
              {suggestions.map((suggestion, index) => (
                <div
                  key={suggestion.id}
                  className={`px-4 py-3 cursor-pointer hover:bg-gray-50 flex items-center ${
                    index === selectedSuggestionIndex ? 'bg-blue-50' : ''
                  }`}
                  onClick={() => handleSuggestionSelect(suggestion)}
                >
                  <div className="flex-1">
                    <div className="text-sm font-medium text-gray-900">{suggestion.text}</div>
                    <div className="text-xs text-gray-500 capitalize">{suggestion.type}</div>
                  </div>
                  <div className="ml-2">
                    <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium ${
                      suggestion.type === 'recipe' ? 'bg-green-100 text-green-800' :
                      suggestion.type === 'workout' ? 'bg-blue-100 text-blue-800' :
                      suggestion.type === 'ingredient' ? 'bg-yellow-100 text-yellow-800' :
                      'bg-purple-100 text-purple-800'
                    }`}>
                      {suggestion.type}
                    </span>
                  </div>
                </div>
              ))}

              {/* Search History */}
              {searchHistory.length > 0 && (
                <>
                  <div className="border-t border-gray-200 px-4 py-2 text-xs font-medium text-gray-500">
                    Recent Searches
                  </div>
                  {searchHistory.slice(0, 3).map((historyItem, index) => (
                    <div
                      key={`history-${index}`}
                      className="px-4 py-2 cursor-pointer hover:bg-gray-50 text-sm text-gray-600 flex items-center"
                      onClick={() => handleSuggestionSelect({
                        id: `history-${index}`,
                        text: historyItem,
                        type: 'recipe'
                      })}
                    >
                      <svg className="h-4 w-4 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      {historyItem}
                    </div>
                  ))}
                </>
              )}
            </div>
          )}
        </form>
      </div>

      {/* Error Display */}
      {error && (
        <ErrorDisplay
          error={error}
          title="Search Error"
          onRetry={() => performSearch()}
          className="mb-4"
        />
      )}

      {/* Search Filters */}
      <SearchFilters
        filters={filters}
        onFiltersChange={applyFilters}
        className="mb-6"
      />

      {/* Loading State */}
      {loading && <LoadingSkeleton count={3} height="h-32" />}

      {/* Quick Search Categories */}
      <div className="mb-6">
        <h3 className="text-lg font-medium text-gray-900 mb-3">Quick Search</h3>
        <div className="flex flex-wrap gap-2">
          {['Chicken', 'Vegetarian', 'Quick', 'Healthy', 'Low Carb', 'High Protein'].map((tag) => (
            <button
              key={tag}
              onClick={() => handleInputChange(tag)}
              className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {tag}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default AdvancedSearch;
