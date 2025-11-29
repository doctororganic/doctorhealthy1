'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRecipes, useWorkouts } from './useNutritionData';

interface SearchSuggestion {
  id: string;
  text: string;
  type: 'recipe' | 'workout' | 'ingredient' | 'diet';
}

interface SearchFilters {
  category: string;
  dietary: string[];
  difficulty: string;
  timeRange: { min: number; max: number };
  calories: { min: number; max: number };
  isHalal: boolean;
}

interface UseSearchOptions {
  debounceMs?: number;
  minQueryLength?: number;
  maxSuggestions?: number;
}

export function useSearch(options: UseSearchOptions = {}) {
  const {
    debounceMs = 300,
    minQueryLength = 2,
    maxSuggestions = 8
  } = options;

  const [query, setQuery] = useState('');
  const [filters, setFilters] = useState<SearchFilters>({
    category: '',
    dietary: [],
    difficulty: '',
    timeRange: { min: 0, max: 120 },
    calories: { min: 0, max: 2000 },
    isHalal: false
  });

  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchHistory, setSearchHistory] = useState<string[]>([]);

  // Nutrition data hooks for getting suggestions
  const { data: recipesData, loading: recipesLoading } = useRecipes({
    page: 1,
    limit: 20
  });

  const { data: workoutsData, loading: workoutsLoading } = useWorkouts({
    page: 1,
    limit: 20
  });

  // Load search history from localStorage
  useEffect(() => {
    try {
      const history = localStorage.getItem('searchHistory');
      if (history) {
        setSearchHistory(JSON.parse(history));
      }
    } catch (error) {
      console.error('Failed to load search history:', error);
    }
  }, []);

  // Save search history to localStorage
  const saveToHistory = useCallback((newQuery: string) => {
    if (!newQuery.trim()) return;

    const newHistory = [newQuery, ...searchHistory.filter(h => h !== newQuery)].slice(0, 10);
    setSearchHistory(newHistory);
    
    try {
      localStorage.setItem('searchHistory', JSON.stringify(newHistory));
    } catch (error) {
      console.error('Failed to save search history:', error);
    }
  }, [searchHistory]);

  // Generate suggestions based on query
  const generateSuggestions = useCallback((searchQuery: string): SearchSuggestion[] => {
    if (searchQuery.length < minQueryLength) return [];

    const suggestions: SearchSuggestion[] = [];
    const queryLower = searchQuery.toLowerCase();

    // Recipe suggestions
    if (recipesData?.items) {
      recipesData.items.slice(0, 5).forEach((recipe: any) => {
        if (recipe.name?.toLowerCase().includes(queryLower)) {
          suggestions.push({
            id: `recipe-${recipe.id}`,
            text: recipe.name,
            type: 'recipe'
          });
        }
      });
    }

    // Workout suggestions
    if (workoutsData?.items) {
      workoutsData.items.slice(0, 5).forEach((workout: any) => {
        const title = workout.title?.en || workout.title?.ar || '';
        if (title.toLowerCase().includes(queryLower)) {
          suggestions.push({
            id: `workout-${workout.id}`,
            text: title,
            type: 'workout'
          });
        }
      });
    }

    // History suggestions (limited to show first)
    const historySuggestions = searchHistory
      .filter(h => h.toLowerCase().includes(queryLower))
      .slice(0, 3)
      .map((historyItem, index) => ({
        id: `history-${index}`,
        text: historyItem,
        type: 'recipe' as const // Default to recipe type for history
      }));

    return [...suggestions, ...historySuggestions].slice(0, maxSuggestions);
  }, [recipesData?.items, workoutsData?.items, searchHistory, minQueryLength, maxSuggestions]);

  // Debounced suggestion generation
  const debouncedGenerateSuggestions = useCallback(
    debounce((searchQuery: string) => {
      const newSuggestions = generateSuggestions(searchQuery);
      setSuggestions(newSuggestions);
      setIsSearching(false);
    }, debounceMs),
    [generateSuggestions, debounceMs]
  );

  // Handle query change
  const handleQueryChange = useCallback((newQuery: string) => {
    setQuery(newQuery);
    
    if (newQuery.length >= minQueryLength) {
      setIsSearching(true);
      debouncedGenerateSuggestions(newQuery);
    } else {
      setSuggestions([]);
      setIsSearching(false);
    }
  }, [minQueryLength, debouncedGenerateSuggestions]);

  // Handle search execution
  const handleSearch = useCallback((searchQuery?: string, searchFilters?: SearchFilters) => {
    const finalQuery = searchQuery || query;
    const finalFilters = searchFilters || filters;

    if (finalQuery.trim()) {
      saveToHistory(finalQuery.trim());
    }

    // Return search parameters for parent component to handle
    return {
      query: finalQuery,
      filters: finalFilters
    };
  }, [query, filters, saveToHistory]);

  // Clear search
  const clearSearch = useCallback(() => {
    setQuery('');
    setSuggestions([]);
    setFilters({
      category: '',
      dietary: [],
      difficulty: '',
      timeRange: { min: 0, max: 120 },
      calories: { min: 0, max: 2000 },
      isHalal: false
    });
  }, []);

  // Update filters
  const updateFilters = useCallback((newFilters: Partial<SearchFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Add quick search term
  const addQuickSearch = useCallback((term: string) => {
    setQuery(term);
    handleSearch(term, filters);
  }, [handleSearch, filters]);

  // Check if loading (either suggestions or data loading)
  const isLoading = isSearching || recipesLoading || workoutsLoading;

  return {
    // State
    query,
    filters,
    suggestions,
    isSearching: isLoading,
    searchHistory,
    
    // Actions
    setQuery: handleQueryChange,
    setFilters: updateFilters,
    onSearch: handleSearch,
    clearSearch,
    addQuickSearch,
    
    // Data
    recipes: recipesData?.items || [],
    workouts: workoutsData?.items || [],
    
    // Utilities
    hasQuery: query.length >= minQueryLength,
    hasFilters: !!(filters.category || filters.dietary.length > 0 || filters.difficulty || filters.isHalal),
    
    // Debounced functions
    debouncedGenerateSuggestions
  };
}

// Simple debounce utility (inline to avoid lodash dependency)
function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

export default useSearch;
export type { SearchSuggestion, SearchFilters, UseSearchOptions };
