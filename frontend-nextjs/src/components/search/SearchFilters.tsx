'use client';

import { useState } from 'react';

export interface SearchFiltersType {
  category: string;
  dietary: string[];
  difficulty: string;
  timeRange: { min: number; max: number };
  calories: { min: number; max: number };
  isHalal: boolean;
}

interface SearchFiltersProps {
  filters: SearchFiltersType;
  onFiltersChange: (filters: SearchFiltersType) => void;
  className?: string;
}

export function SearchFilters({ filters, onFiltersChange, className = '' }: SearchFiltersProps) {
  const [isExpanded, setIsExpanded] = useState(false);

  const updateFilters = (updates: Partial<SearchFiltersType>) => {
    onFiltersChange({ ...filters, ...updates });
  };

  const toggleDietary = (diet: string) => {
    const current = [...filters.dietary];
    const index = current.indexOf(diet);
    
    if (index > -1) {
      current.splice(index, 1);
    } else {
      current.push(diet);
    }
    
    updateFilters({ dietary: current });
  };

  const categories = [
    { value: '', label: 'All Categories' },
    { value: 'breakfast', label: 'Breakfast' },
    { value: 'lunch', label: 'Lunch' },
    { value: 'dinner', label: 'Dinner' },
    { value: 'snack', label: 'Snack' },
    { value: 'dessert', label: 'Dessert' }
  ];

  const difficulties = [
    { value: '', label: 'All Levels' },
    { value: 'easy', label: 'Easy' },
    { value: 'medium', label: 'Medium' },
    { value: 'hard', label: 'Hard' }
  ];

  const dietaryOptions = [
    { value: 'vegetarian', label: 'Vegetarian' },
    { value: 'vegan', label: 'Vegan' },
    { value: 'gluten-free', label: 'Gluten-Free' },
    { value: 'dairy-free', label: 'Dairy-Free' },
    { value: 'keto', label: 'Keto' },
    { value: 'paleo', label: 'Paleo' },
    { value: 'low-carb', label: 'Low Carb' },
    { value: 'low-sodium', label: 'Low Sodium' }
  ];

  return (
    <div className={`bg-white rounded-lg border border-gray-200 ${className}`}>
      {/* Filters Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-medium text-gray-900">Search Filters</h3>
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="text-gray-500 hover:text-gray-700 focus:outline-none focus:text-gray-700"
          >
            <svg
              className={`h-5 w-5 transform transition-transform ${isExpanded ? 'rotate-180' : ''}`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </button>
        </div>

        {/* Active Filters Summary */}
        <div className="mt-2 flex flex-wrap gap-2">
          {filters.category && (
            <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800">
              {categories.find(c => c.value === filters.category)?.label}
              <button
                onClick={() => updateFilters({ category: '' })}
                className="ml-1 text-blue-600 hover:text-blue-800"
              >
                ×
              </button>
            </span>
          )}
          {filters.difficulty && (
            <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-green-100 text-green-800">
              {filters.difficulty}
              <button
                onClick={() => updateFilters({ difficulty: '' })}
                className="ml-1 text-green-600 hover:text-green-800"
              >
                ×
              </button>
            </span>
          )}
          {filters.dietary.map(diet => (
            <span
              key={diet}
              className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-purple-100 text-purple-800"
            >
              {diet}
              <button
                onClick={() => toggleDietary(diet)}
                className="ml-1 text-purple-600 hover:text-purple-800"
              >
                ×
              </button>
            </span>
          ))}
          {filters.isHalal && (
            <span className="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
              Halal
              <button
                onClick={() => updateFilters({ isHalal: false })}
                className="ml-1 text-yellow-600 hover:text-yellow-800"
              >
                ×
              </button>
            </span>
          )}
        </div>
      </div>

      {/* Collapsible Filters Content */}
      {isExpanded && (
        <div className="p-4 space-y-6">
          {/* Category Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Category
            </label>
            <select
              value={filters.category}
              onChange={(e) => updateFilters({ category: e.target.value })}
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              {categories.map(category => (
                <option key={category.value} value={category.value}>
                  {category.label}
                </option>
              ))}
            </select>
          </div>

          {/* Difficulty Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Difficulty
            </label>
            <select
              value={filters.difficulty}
              onChange={(e) => updateFilters({ difficulty: e.target.value })}
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              {difficulties.map(difficulty => (
                <option key={difficulty.value} value={difficulty.value}>
                  {difficulty.label}
                </option>
              ))}
            </select>
          </div>

          {/* Dietary Restrictions */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Dietary Restrictions
            </label>
            <div className="space-y-2">
              {dietaryOptions.map(option => (
                <label key={option.value} className="flex items-center">
                  <input
                    type="checkbox"
                    checked={filters.dietary.includes(option.value)}
                    onChange={() => toggleDietary(option.value)}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <span className="ml-2 text-sm text-gray-700">{option.label}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Time Range Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Cooking Time (minutes)
            </label>
            <div className="flex items-center space-x-4">
              <div className="flex-1">
                <input
                  type="number"
                  min="0"
                  max="120"
                  value={filters.timeRange.min}
                  onChange={(e) => updateFilters({ 
                    timeRange: { ...filters.timeRange, min: parseInt(e.target.value) || 0 }
                  })}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Min"
                />
              </div>
              <span className="text-gray-500">to</span>
              <div className="flex-1">
                <input
                  type="number"
                  min="0"
                  max="120"
                  value={filters.timeRange.max}
                  onChange={(e) => updateFilters({ 
                    timeRange: { ...filters.timeRange, max: parseInt(e.target.value) || 120 }
                  })}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Max"
                />
              </div>
            </div>
          </div>

          {/* Calorie Range Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Calories
            </label>
            <div className="flex items-center space-x-4">
              <div className="flex-1">
                <input
                  type="number"
                  min="0"
                  max="2000"
                  value={filters.calories.min}
                  onChange={(e) => updateFilters({ 
                    calories: { ...filters.calories, min: parseInt(e.target.value) || 0 }
                  })}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Min"
                />
              </div>
              <span className="text-gray-500">to</span>
              <div className="flex-1">
                <input
                  type="number"
                  min="0"
                  max="2000"
                  value={filters.calories.max}
                  onChange={(e) => updateFilters({ 
                    calories: { ...filters.calories, max: parseInt(e.target.value) || 2000 }
                  })}
                  className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Max"
                />
              </div>
            </div>
          </div>

          {/* Halal Filter */}
          <div>
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={filters.isHalal}
                onChange={(e) => updateFilters({ isHalal: e.target.checked })}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <span className="ml-2 text-sm font-medium text-gray-700">
                Halal Only
              </span>
            </label>
          </div>

          {/* Clear Filters Button */}
          <div className="pt-4 border-t border-gray-200">
            <button
              onClick={() => updateFilters({
                category: '',
                dietary: [],
                difficulty: '',
                timeRange: { min: 0, max: 120 },
                calories: { min: 0, max: 2000 },
                isHalal: false
              })}
              className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              Clear All Filters
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default SearchFilters;
