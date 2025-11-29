'use client';

import { useState, useCallback, useEffect } from 'react';

// Enhanced meal plan types for 10K users
export interface MealPlanRequest {
  userId?: number;
  goal: 'weight_loss' | 'muscle_gain' | 'maintenance' | 'bulking' | 'cutting';
  calories: number;
  dietType: 'halal' | 'vegetarian' | 'vegan' | 'keto' | 'mediterranean' | 'balanced';
  days: number;
  allergies?: string[];
  healthConditions?: string[];
  preferences?: {
    cookingTime?: 'quick' | 'medium' | 'elaborate';
    mealFrequency?: 3 | 4 | 5 | 6;
    cuisineTypes?: string[];
    avoidIngredients?: string[];
  };
  macroTargets?: {
    protein: number;
    carbs: number;
    fat: number;
  };
}

export interface Meal {
  id: string;
  name: string | { en: string; ar: string };
  type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  calories: number;
  macros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber: number;
  };
  ingredients: Array<{
    name: string | { en: string; ar: string };
    amount: number;
    unit: string;
    calories: number;
  }>;
  instructions: string[] | { en: string[]; ar: string[] };
  prepTime: number;
  cookTime: number;
  servings: number;
  isHalal: boolean;
  allergens: string[];
  tags: string[];
}

export interface DayMealPlan {
  day: number;
  date: string;
  totalCalories: number;
  totalMacros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber: number;
  };
  meals: {
    breakfast: Meal;
    lunch: Meal;
    dinner: Meal;
    snacks: Meal[];
  };
  waterIntake: number; // liters
  supplements?: Array<{
    name: string;
    dosage: string;
    timing: string;
  }>;
}

export interface MealPlanResponse {
  id: string;
  userId: number;
  goal: string;
  dietType: string;
  totalDays: number;
  dailyCalorieTarget: number;
  dailyMacroTargets: {
    protein: number;
    carbs: number;
    fat: number;
  };
  days: DayMealPlan[];
  shoppingList: Array<{
    ingredient: string;
    totalAmount: number;
    unit: string;
    category: string;
    estimatedCost?: number;
  }>;
  nutritionalAnalysis: {
    averageDailyNutrition: {
      calories: number;
      protein: number;
      carbs: number;
      fat: number;
      fiber: number;
      vitamins: Record<string, number>;
      minerals: Record<string, number>;
    };
    weeklyTotals: {
      calories: number;
      cost?: number;
    };
  };
  createdAt: string;
  validUntil: string;
  lastModified: string;
}

export interface MealPlanHookState {
  mealPlan: MealPlanResponse | null;
  isGenerating: boolean;
  error: string | null;
  lastGenerated: string | null;
}

// Custom hook for meal plan management
export function useMealPlan() {
  const [state, setState] = useState<MealPlanHookState>({
    mealPlan: null,
    isGenerating: false,
    error: null,
    lastGenerated: null,
  });

  // Generate meal plan with smart error handling
  const generateMealPlan = useCallback(async (request: MealPlanRequest): Promise<MealPlanResponse | null> => {
    setState(prev => ({ ...prev, isGenerating: true, error: null }));

    try {
      // Validate request
      const validationError = validateMealPlanRequest(request);
      if (validationError) {
        throw new Error(validationError);
      }

      // Make API request
      const response = await fetch('/api/v1/actions/generate-meal-plan', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      
      if (data.status !== 'success') {
        throw new Error(data.error || 'Failed to generate meal plan');
      }

      const mealPlan = data.data as MealPlanResponse;
      
      setState(prev => ({
        ...prev,
        mealPlan,
        isGenerating: false,
        error: null,
        lastGenerated: new Date().toISOString(),
      }));

      return mealPlan;

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      setState(prev => ({
        ...prev,
        mealPlan: null,
        isGenerating: false,
        error: errorMessage,
        lastGenerated: null,
      }));

      return null;
    }
  }, []);

  // Save meal plan to local storage for offline access
  const saveMealPlan = useCallback((mealPlan: MealPlanResponse) => {
    try {
      localStorage.setItem(`meal-plan-${mealPlan.id}`, JSON.stringify(mealPlan));
      localStorage.setItem('last-meal-plan-id', mealPlan.id);
    } catch (error) {
      console.warn('Failed to save meal plan to local storage:', error);
    }
  }, []);

  // Load meal plan from local storage
  const loadSavedMealPlan = useCallback(async (planId?: string): Promise<MealPlanResponse | null> => {
    try {
      const id = planId || localStorage.getItem('last-meal-plan-id');
      if (!id) return null;

      const savedData = localStorage.getItem(`meal-plan-${id}`);
      if (!savedData) return null;

      const mealPlan = JSON.parse(savedData) as MealPlanResponse;
      
      // Check if meal plan is still valid
      if (new Date() > new Date(mealPlan.validUntil)) {
        localStorage.removeItem(`meal-plan-${id}`);
        return null;
      }

      setState(prev => ({
        ...prev,
        mealPlan,
        error: null,
        lastGenerated: mealPlan.createdAt,
      }));

      return mealPlan;

    } catch (error) {
      console.warn('Failed to load meal plan from local storage:', error);
      return null;
    }
  }, []);

  // Update specific meal in plan
  const updateMeal = useCallback(async (dayIndex: number, mealType: keyof DayMealPlan['meals'], newMeal: Meal) => {
    if (!state.mealPlan) return false;

    setState(prev => ({
      ...prev,
      isGenerating: true,
      error: null,
    }));

    try {
      // Create updated meal plan
      const updatedPlan = { ...state.mealPlan };
      updatedPlan.days[dayIndex].meals[mealType] = newMeal as any;
      
      // Recalculate day totals
      updatedPlan.days[dayIndex] = recalculateDayTotals(updatedPlan.days[dayIndex]);
      updatedPlan.lastModified = new Date().toISOString();

      // Update state
      setState(prev => ({
        ...prev,
        mealPlan: updatedPlan,
        isGenerating: false,
        error: null,
      }));

      // Save updated plan
      saveMealPlan(updatedPlan);

      return true;

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update meal';
      setState(prev => ({
        ...prev,
        isGenerating: false,
        error: errorMessage,
      }));

      return false;
    }
  }, [state.mealPlan, saveMealPlan]);

  // Clear current meal plan
  const clearMealPlan = useCallback(() => {
    setState({
      mealPlan: null,
      isGenerating: false,
      error: null,
      lastGenerated: null,
    });
  }, []);

  // Export meal plan to different formats
  const exportMealPlan = useCallback((format: 'json' | 'csv' | 'pdf' = 'json') => {
    if (!state.mealPlan) return null;

    switch (format) {
      case 'json':
        return JSON.stringify(state.mealPlan, null, 2);
      
      case 'csv':
        return convertMealPlanToCSV(state.mealPlan);
      
      case 'pdf':
        // Would integrate with a PDF generation library
        console.log('PDF export would be implemented with jsPDF or similar');
        return null;
      
      default:
        return null;
    }
  }, [state.mealPlan]);

  return {
    ...state,
    generateMealPlan,
    saveMealPlan,
    loadSavedMealPlan,
    updateMeal,
    clearMealPlan,
    exportMealPlan,
  };
}

// Validation helper
function validateMealPlanRequest(request: MealPlanRequest): string | null {
  if (!request.goal) {
    return 'Goal is required';
  }

  if (!request.calories || request.calories < 800 || request.calories > 5000) {
    return 'Calories must be between 800 and 5000';
  }

  if (!request.dietType) {
    return 'Diet type is required';
  }

  if (!request.days || request.days < 1 || request.days > 30) {
    return 'Days must be between 1 and 30';
  }

  if (request.allergies && request.allergies.length > 10) {
    return 'Maximum 10 allergies allowed';
  }

  return null;
}

// Helper function to recalculate day totals
function recalculateDayTotals(day: DayMealPlan): DayMealPlan {
  const meals = [day.meals.breakfast, day.meals.lunch, day.meals.dinner, ...day.meals.snacks];
  
  const totalCalories = meals.reduce((sum, meal) => sum + meal.calories, 0);
  const totalMacros = meals.reduce((totals, meal) => ({
    protein: totals.protein + meal.macros.protein,
    carbs: totals.carbs + meal.macros.carbs,
    fat: totals.fat + meal.macros.fat,
    fiber: totals.fiber + meal.macros.fiber,
  }), { protein: 0, carbs: 0, fat: 0, fiber: 0 });

  return {
    ...day,
    totalCalories,
    totalMacros,
  };
}

// Helper function to convert meal plan to CSV
function convertMealPlanToCSV(mealPlan: MealPlanResponse): string {
  const headers = ['Day', 'Meal Type', 'Meal Name', 'Calories', 'Protein', 'Carbs', 'Fat', 'Prep Time'];
  const rows = [headers.join(',')];

  mealPlan.days.forEach(day => {
    const addMealRow = (meal: Meal, type: string) => {
      const name = typeof meal.name === 'string' ? meal.name : meal.name.en;
      rows.push([
        day.day,
        type,
        name,
        meal.calories,
        meal.macros.protein,
        meal.macros.carbs,
        meal.macros.fat,
        meal.prepTime,
      ].join(','));
    };

    addMealRow(day.meals.breakfast, 'Breakfast');
    addMealRow(day.meals.lunch, 'Lunch');
    addMealRow(day.meals.dinner, 'Dinner');
    day.meals.snacks.forEach((snack, index) => {
      addMealRow(snack, `Snack ${index + 1}`);
    });
  });

  return rows.join('\n');
}

// Hook for meal plan statistics and analytics
export function useMealPlanAnalytics(mealPlan: MealPlanResponse | null) {
  const [analytics, setAnalytics] = useState<{
    adherenceScore: number;
    nutritionalBalance: Record<string, number>;
    costEffectiveness: number;
    varietyScore: number;
  } | null>(null);

  // Calculate analytics when meal plan changes
  useEffect(() => {
    if (!mealPlan) {
      setAnalytics(null);
      return;
    }

    // Calculate various metrics
    const adherenceScore = calculateAdherenceScore(mealPlan);
    const nutritionalBalance = calculateNutritionalBalance(mealPlan);
    const costEffectiveness = calculateCostEffectiveness(mealPlan);
    const varietyScore = calculateVarietyScore(mealPlan);

    setAnalytics({
      adherenceScore,
      nutritionalBalance,
      costEffectiveness,
      varietyScore,
    });

  }, [mealPlan]);

  return analytics;
}

// Helper analytics functions
function calculateAdherenceScore(mealPlan: MealPlanResponse): number {
  // Implementation for adherence score calculation
  return 85; // Placeholder
}

function calculateNutritionalBalance(mealPlan: MealPlanResponse): Record<string, number> {
  // Implementation for nutritional balance calculation
  return { protein: 90, carbs: 85, fat: 80, vitamins: 75 }; // Placeholder
}

function calculateCostEffectiveness(mealPlan: MealPlanResponse): number {
  // Implementation for cost effectiveness calculation
  return 78; // Placeholder
}

function calculateVarietyScore(mealPlan: MealPlanResponse): number {
  // Implementation for variety score calculation
  return 82; // Placeholder
}