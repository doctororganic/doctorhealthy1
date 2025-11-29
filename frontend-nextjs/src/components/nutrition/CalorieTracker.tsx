'use client';

import { useState, useEffect } from 'react';
import { LoadingSkeleton } from '../ui/LoadingSkeleton';
import { ErrorDisplay } from '../ui/ErrorDisplay';
import { EmptyState } from '../ui/EmptyState';
import { Pagination } from '../ui/Pagination';
import { calculateBMI, getBMICategory, calculateTDEE } from '../../utils/nutritionCalculations';

interface CalorieEntry {
  id: string;
  date: string;
  meal: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  food: string;
  calories: number;
  protein?: number;
  carbs?: number;
  fat?: number;
  fiber?: number;
  notes?: string;
}

interface DailySummary {
  date: string;
  totalCalories: number;
  totalProtein: number;
  totalCarbs: number;
  totalFat: number;
  totalFiber: number;
  meals: CalorieEntry[];
  goalCalories?: number;
  remainingCalories?: number;
}

interface CalorieTrackerProps {
  targetCalories?: number;
  weight?: number;
  height?: number;
  age?: number;
  gender?: 'male' | 'female';
  activityLevel?: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extra_active';
}

export function CalorieTracker({
  targetCalories = 2000,
  weight = 70,
  height = 175,
  age = 30,
  gender = 'male',
  activityLevel = 'moderately_active'
}: CalorieTrackerProps) {
  const [entries, setEntries] = useState<CalorieEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingEntry, setEditingEntry] = useState<CalorieEntry | null>(null);

  const [newEntry, setNewEntry] = useState<Partial<CalorieEntry>>({
    meal: 'breakfast',
    food: '',
    calories: 0,
    protein: 0,
    carbs: 0,
    fat: 0,
    fiber: 0,
    notes: ''
  });

  const entriesPerPage = 10;

  // Calculate daily summaries
  const dailySummaries = entries.reduce((acc: Record<string, DailySummary>, entry) => {
    const date = entry.date;
    if (!acc[date]) {
      acc[date] = {
        date,
        totalCalories: 0,
        totalProtein: 0,
        totalCarbs: 0,
        totalFat: 0,
        totalFiber: 0,
        meals: [],
        goalCalories: targetCalories
      };
    }
    
    acc[date].meals.push(entry);
    acc[date].totalCalories += entry.calories;
    acc[date].totalProtein += entry.protein || 0;
    acc[date].totalCarbs += entry.carbs || 0;
    acc[date].totalFat += entry.fat || 0;
    acc[date].totalFiber += entry.fiber || 0;
    acc[date].remainingCalories = targetCalories - acc[date].totalCalories;
    
    return acc;
  }, {});

  const sortedDates = Object.keys(dailySummaries).sort((a, b) => new Date(b).getTime() - new Date(a).getTime());
  const paginatedDates = sortedDates.slice(
    (currentPage - 1) * entriesPerPage,
    currentPage * entriesPerPage
  );

  // Load entries from localStorage
  useEffect(() => {
    try {
      const savedEntries = localStorage.getItem('calorieEntries');
      if (savedEntries) {
        setEntries(JSON.parse(savedEntries));
      }
    } catch (err) {
      console.error('Failed to load entries:', err);
    }
  }, []);

  // Save entries to localStorage
  const saveEntries = (newEntries: CalorieEntry[]) => {
    try {
      localStorage.setItem('calorieEntries', JSON.stringify(newEntries));
    } catch (err) {
      console.error('Failed to save entries:', err);
      setError('Failed to save entries to storage');
    }
  };

  // Add or update entry
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (!newEntry.food || !newEntry.calories || newEntry.calories <= 0) {
        throw new Error('Food name and calories are required');
      }

      const entry: CalorieEntry = {
        id: editingEntry?.id || Date.now().toString(),
        date: editingEntry?.date || new Date().toISOString().split('T')[0],
        meal: newEntry.meal || 'breakfast',
        food: newEntry.food || '',
        calories: newEntry.calories || 0,
        protein: newEntry.protein || 0,
        carbs: newEntry.carbs || 0,
        fat: newEntry.fat || 0,
        fiber: newEntry.fiber || 0,
        notes: newEntry.notes || ''
      };

      let updatedEntries: CalorieEntry[];
      if (editingEntry) {
        updatedEntries = entries.map(e => e.id === editingEntry.id ? entry : e);
      } else {
        updatedEntries = [entry, ...entries];
      }

      setEntries(updatedEntries);
      saveEntries(updatedEntries);
      
      // Reset form
      setNewEntry({
        meal: 'breakfast',
        food: '',
        calories: 0,
        protein: 0,
        carbs: 0,
        fat: 0,
        fiber: 0,
        notes: ''
      });
      setEditingEntry(null);
      setShowAddForm(false);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save entry');
    } finally {
      setLoading(false);
    }
  };

  // Delete entry
  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this entry?')) {
      const updatedEntries = entries.filter(e => e.id !== id);
      setEntries(updatedEntries);
      saveEntries(updatedEntries);
    }
  };

  // Edit entry
  const handleEdit = (entry: CalorieEntry) => {
    setEditingEntry(entry);
    setNewEntry(entry);
    setShowAddForm(true);
  };

  // Cancel editing
  const handleCancel = () => {
    setEditingEntry(null);
    setNewEntry({
      meal: 'breakfast',
      food: '',
      calories: 0,
      protein: 0,
      carbs: 0,
      fat: 0,
      fiber: 0,
      notes: ''
    });
    setShowAddForm(false);
  };

  // Calculate statistics
  const todayEntries = entries.filter(e => e.date === new Date().toISOString().split('T')[0]);
  const todayCalories = todayEntries.reduce((sum, e) => sum + e.calories, 0);
  const weekEntries = entries.filter(e => {
    const entryDate = new Date(e.date);
    const weekAgo = new Date();
    weekAgo.setDate(weekAgo.getDate() - 7);
    return entryDate >= weekAgo;
  });
  const weekAverage = weekEntries.length > 0 ? weekEntries.reduce((sum, e) => sum + e.calories, 0) / weekEntries.length : 0;

  const currentBMI = calculateBMI(weight, height);
  const bmiCategory = getBMICategory(currentBMI);
  const estimatedTDEE = calculateTDEE(
    10 * weight + 6.25 * height - 5 * age + (gender === 'male' ? 5 : -161),
    activityLevel
  );

  if (loading && entries.length === 0) {
    return <LoadingSkeleton count={3} height="h-32" />;
  }

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Calorie Tracker</h1>
        <p className="text-gray-600">Track your daily nutrition and monitor your goals</p>
      </div>

      {error && (
        <ErrorDisplay
          error={error}
          title="Tracking Error"
          onRetry={() => setError(null)}
        />
      )}

      {/* Statistics Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-blue-600">Today's Calories</p>
              <p className="text-2xl font-bold text-blue-900">{Math.round(todayCalories)}</p>
              <p className="text-sm text-blue-700">of {targetCalories} goal</p>
            </div>
            <div className="text-blue-200">
              <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                <path fillRule="evenodd" d="M4 5a2 2 0 012-2 1 1 0 000 2H6a2 2 0 100 4h2a2 2 0 100 4h2a1 1 0 100 2H6a2 2 0 01-2-2V5z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
          <div className="mt-2">
            <div className="w-full bg-blue-200 rounded-full h-2">
              <div 
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${Math.min(100, (todayCalories / targetCalories) * 100)}%` }}
              />
            </div>
          </div>
        </div>

        <div className="bg-green-50 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-green-600">Week Average</p>
              <p className="text-2xl font-bold text-green-900">{Math.round(weekAverage)}</p>
              <p className="text-sm text-green-700">calories/day</p>
            </div>
            <div className="text-green-200">
              <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M5 2a1 1 0 011 1v1h1a1 1 0 010 2H6v1a1 1 0 01-2 0V6H3a1 1 0 010-2h1V3a1 1 0 011-1zm0 10a1 1 0 011 1v1h1a1 1 0 110 2H6v1a1 1 0 11-2 0v-1H3a1 1 0 110-2h1v-1a1 1 0 011-1zM12 2a1 1 0 01.967.744L14.146 7.2 17.5 9.134a1 1 0 010 1.732l-3.354 1.935-1.18 4.455a1 1 0 01-1.933 0L9.854 12.8 6.5 10.866a1 1 0 010-1.732l3.354-1.935 1.18-4.455A1 1 0 0112 2z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-purple-50 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-purple-600">Current BMI</p>
              <p className="text-2xl font-bold text-purple-900">{currentBMI}</p>
              <p className="text-sm text-purple-700">{bmiCategory}</p>
            </div>
            <div className="text-purple-200">
              <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-1 1v2a1 1 0 002 0V8a1 1 0 012-2V3a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-orange-50 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-orange-600">Est. TDEE</p>
              <p className="text-2xl font-bold text-orange-900">{Math.round(estimatedTDEE)}</p>
              <p className="text-sm text-orange-700">calories/day</p>
            </div>
            <div className="text-orange-200">
              <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M12.395 2.553a1 1 0 00-1.45-.385c-.325.196-.594.467-.636.617L6 12.79l-4.195 2.34c-.229.126-.477.198-.759.466-1.068.721-1.123-.363-.432-.372-.654-.132-.623.21-.324.444-.317.894.138 1.395.628 1.39 2.553zM5 7H4a2 2 0 00-2 2v8a2 2 0 002 2h1v-2h4v2h1a2 2 0 002-2v-8a2 2 0 00-2-2z" clipRule="evenodd" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Add Entry Button */}
      <div className="flex justify-center">
        <button
          onClick={() => setShowAddForm(true)}
          className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Add New Entry
        </button>
      </div>

      {/* Add/Edit Entry Form */}
      {showAddForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl p-6 w-full max-w-md mx-4">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              {editingEntry ? 'Edit Entry' : 'Add Food Entry'}
            </h3>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Date</label>
                <input
                  type="date"
                  value={newEntry.date}
                  onChange={(e) => setNewEntry(prev => ({ ...prev, date: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Meal</label>
                <select
                  value={newEntry.meal}
                  onChange={(e) => setNewEntry(prev => ({ ...prev, meal: e.target.value as any }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="breakfast">Breakfast</option>
                  <option value="lunch">Lunch</option>
                  <option value="dinner">Dinner</option>
                  <option value="snack">Snack</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Food Item</label>
                <input
                  type="text"
                  value={newEntry.food}
                  onChange={(e) => setNewEntry(prev => ({ ...prev, food: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter food name"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Calories</label>
                  <input
                    type="number"
                    min="0"
                    value={newEntry.calories}
                    onChange={(e) => setNewEntry(prev => ({ ...prev, calories: parseInt(e.target.value) || 0 }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Protein (g)</label>
                  <input
                    type="number"
                    min="0"
                    step="0.1"
                    value={newEntry.protein}
                    onChange={(e) => setNewEntry(prev => ({ ...prev, protein: parseFloat(e.target.value) || 0 }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Carbs (g)</label>
                  <input
                    type="number"
                    min="0"
                    step="0.1"
                    value={newEntry.carbs}
                    onChange={(e) => setNewEntry(prev => ({ ...prev, carbs: parseFloat(e.target.value) || 0 }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Fat (g)</label>
                  <input
                    type="number"
                    min="0"
                    step="0.1"
                    value={newEntry.fat}
                    onChange={(e) => setNewEntry(prev => ({ ...prev, fat: parseFloat(e.target.value) || 0 }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Fiber (g)</label>
                <input
                  type="number"
                  min="0"
                  step="0.1"
                  value={newEntry.fiber}
                  onChange={(e) => setNewEntry(prev => ({ ...prev, fiber: parseFloat(e.target.value) || 0 }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Notes (optional)</label>
                <textarea
                  value={newEntry.notes}
                  onChange={(e) => setNewEntry(prev => ({ ...prev, notes: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows={3}
                  placeholder="Add any notes about this meal..."
                />
              </div>

              <div className="flex justify-end space-x-3">
                <button
                  type="button"
                  onClick={handleCancel}
                  className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {loading ? 'Saving...' : (editingEntry ? 'Update' : 'Add Entry')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Daily Entries */}
      {paginatedDates.length > 0 ? (
        <div className="space-y-6">
          {paginatedDates.map(date => {
            const summary = dailySummaries[date];
            return (
              <div key={date} className="bg-white rounded-lg shadow-lg p-6">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-lg font-semibold text-gray-900">
                    {new Date(date).toLocaleDateString('en-US', { 
                      weekday: 'long', 
                      year: 'numeric', 
                      month: 'long', 
                      day: 'numeric' 
                    })}
                  </h3>
                  <div className="flex items-center space-x-4">
                    <span className={`text-sm font-medium ${
                      summary.totalCalories <= targetCalories ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {summary.totalCalories} / {targetCalories} cal
                    </span>
                    <span className={`text-sm font-medium ${
                      summary.remainingCalories && summary.remainingCalories >= 0 ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {summary.remainingCalories && summary.remainingCalories >= 0 ? '+' : ''}{summary.remainingCalories || 0} cal left
                    </span>
                  </div>
                </div>

                <div className="space-y-3">
                  {summary.meals.map(meal => (
                    <div key={meal.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                      <div className="flex-1">
                        <div className="flex items-center space-x-3">
                          <span className="text-sm font-medium text-gray-900 capitalize">{meal.meal}</span>
                          <span className="text-sm font-semibold text-gray-900">{meal.food}</span>
                        </div>
                        <div className="flex items-center space-x-4 text-xs text-gray-600">
                          <span>{meal.calories} cal</span>
                          {meal.protein && <span>P: {meal.protein}g</span>}
                          {meal.carbs && <span>C: {meal.carbs}g</span>}
                          {meal.fat && <span>F: {meal.fat}g</span>}
                          {meal.fiber && <span>Fiber: {meal.fiber}g</span>}
                        </div>
                        {meal.notes && (
                          <p className="text-sm text-gray-600 mt-1">{meal.notes}</p>
                        )}
                      </div>
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => handleEdit(meal)}
                          className="text-blue-600 hover:text-blue-800 text-sm"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(meal.id)}
                          className="text-red-600 hover:text-red-800 text-sm"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Daily Totals */}
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div className="text-center">
                      <p className="text-gray-600">Total Calories</p>
                      <p className="font-bold text-lg">{summary.totalCalories}</p>
                    </div>
                    <div className="text-center">
                      <p className="text-gray-600">Protein</p>
                      <p className="font-bold text-lg">{summary.totalProtein}g</p>
                    </div>
                    <div className="text-center">
                      <p className="text-gray-600">Carbs</p>
                      <p className="font-bold text-lg">{summary.totalCarbs}g</p>
                    </div>
                    <div className="text-center">
                      <p className="text-gray-600">Fat</p>
                      <p className="font-bold text-lg">{summary.totalFat}g</p>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <EmptyState
          title="No Entries Yet"
          message="Start tracking your nutrition by adding your first food entry."
          actionLabel="Add Your First Entry"
          onAction={() => setShowAddForm(true)}
        />
      )}

      {/* Pagination */}
      {sortedDates.length > entriesPerPage && (
        <div className="mt-6">
          <Pagination
            currentPage={currentPage}
            totalPages={Math.ceil(sortedDates.length / entriesPerPage)}
            onPageChange={setCurrentPage}
          />
        </div>
      )}
    </div>
  );
}

export default CalorieTracker;
