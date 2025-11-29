'use client';

import { useState, useEffect } from 'react';
import { LoadingSkeleton } from '../ui/LoadingSkeleton';
import { ErrorDisplay } from '../ui/ErrorDisplay';
import { calculateMacros } from '../../utils/nutritionCalculations';

interface UserProfile {
  age: number;
  gender: 'male' | 'female';
  height: number; // cm
  weight: number; // kg
  activityLevel: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extra_active';
  goal: 'lose_weight' | 'maintain' | 'gain_weight' | 'gain_muscle';
}

interface CalculationResults {
  bmr: number;
  tdee: number;
  targetCalories: number;
  macros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber: number;
  };
  water: number; // liters
}

export function NutritionCalculator() {
  const [profile, setProfile] = useState<UserProfile>({
    age: 30,
    gender: 'male',
    height: 175,
    weight: 70,
    activityLevel: 'moderately_active',
    goal: 'maintain'
  });

  const [results, setResults] = useState<CalculationResults | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const activityMultipliers = {
    sedentary: 1.2,
    lightly_active: 1.375,
    moderately_active: 1.55,
    very_active: 1.725,
    extra_active: 1.9
  };

  const goalAdjustments = {
    lose_weight: -500, // 500 calorie deficit
    maintain: 0,
    gain_weight: 500, // 500 calorie surplus
    gain_muscle: 300 // 300 calorie surplus
  };

  const calculateNutrition = () => {
    setLoading(true);
    setError(null);

    try {
      // Validate inputs
      if (profile.age < 15 || profile.age > 100) {
        throw new Error('Age must be between 15 and 100');
      }
      if (profile.height < 100 || profile.height > 250) {
        throw new Error('Height must be between 100 and 250 cm');
      }
      if (profile.weight < 30 || profile.weight > 300) {
        throw new Error('Weight must be between 30 and 300 kg');
      }

      // Calculate BMR using Mifflin-St Jeor equation
      let bmr: number;
      if (profile.gender === 'male') {
        bmr = 10 * profile.weight + 6.25 * profile.height - 5 * profile.age + 5;
      } else {
        bmr = 10 * profile.weight + 6.25 * profile.height - 5 * profile.age - 161;
      }

      // Calculate TDEE
      const tdee = bmr * activityMultipliers[profile.activityLevel];

      // Apply goal adjustment
      const targetCalories = tdee + goalAdjustments[profile.goal];

      // Calculate macros (40% protein, 40% carbs, 20% fat for muscle gain)
      let proteinRatio = 0.3, carbsRatio = 0.4, fatRatio = 0.3;
      
      if (profile.goal === 'gain_muscle') {
        proteinRatio = 0.4; carbsRatio = 0.3; fatRatio = 0.3;
      } else if (profile.goal === 'lose_weight') {
        proteinRatio = 0.4; carbsRatio = 0.3; fatRatio = 0.3;
      }

      const macros = calculateMacros(targetCalories, proteinRatio, carbsRatio, fatRatio);

      // Water recommendation (ml per kg body weight)
      const water = Math.round((profile.weight * 35) / 1000 * 10) / 10; // in liters

      setResults({
        bmr: Math.round(bmr),
        tdee: Math.round(tdee),
        targetCalories: Math.round(targetCalories),
        macros,
        water
      });

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Calculation failed');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    calculateNutrition();
  }, [profile]);

  const updateProfile = (field: keyof UserProfile, value: any) => {
    setProfile(prev => ({ ...prev, [field]: value }));
  };

  const getGoalDescription = (goal: string) => {
    switch (goal) {
      case 'lose_weight':
        return 'Lose Weight (500 calorie deficit)';
      case 'gain_weight':
        return 'Gain Weight (500 calorie surplus)';
      case 'gain_muscle':
        return 'Gain Muscle (300 calorie surplus)';
      default:
        return 'Maintain Weight';
    }
  };

  const getActivityDescription = (level: string) => {
    switch (level) {
      case 'sedentary':
        return 'Little or no exercise';
      case 'lightly_active':
        return 'Light exercise 1-3 days/week';
      case 'moderately_active':
        return 'Moderate exercise 3-5 days/week';
      case 'very_active':
        return 'Hard exercise 6-7 days/week';
      case 'extra_active':
        return 'Very hard exercise, physical job';
      default:
        return '';
    }
  };

  if (loading) {
    return <LoadingSkeleton count={3} height="h-20" />;
  }

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Nutrition Calculator</h1>
        <p className="text-gray-600">Calculate your daily nutritional needs based on your profile</p>
      </div>

      {error && (
        <ErrorDisplay
          error={error}
          title="Calculation Error"
          onRetry={calculateNutrition}
        />
      )}

      {/* Input Form */}
      <div className="bg-white rounded-lg shadow-lg p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-6">Your Profile</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Basic Information */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-gray-800">Basic Information</h3>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Age</label>
              <input
                type="number"
                min="15"
                max="100"
                value={profile.age}
                onChange={(e) => updateProfile('age', parseInt(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Gender</label>
              <select
                value={profile.gender}
                onChange={(e) => updateProfile('gender', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="male">Male</option>
                <option value="female">Female</option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Height (cm)</label>
                <input
                  type="number"
                  min="100"
                  max="250"
                  value={profile.height}
                  onChange={(e) => updateProfile('height', parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Weight (kg)</label>
                <input
                  type="number"
                  min="30"
                  max="300"
                  value={profile.weight}
                  onChange={(e) => updateProfile('weight', parseInt(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>
          </div>

          {/* Activity and Goals */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium text-gray-800">Activity & Goals</h3>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Activity Level</label>
              <select
                value={profile.activityLevel}
                onChange={(e) => updateProfile('activityLevel', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="sedentary">{getActivityDescription('sedentary')}</option>
                <option value="lightly_active">{getActivityDescription('lightly_active')}</option>
                <option value="moderately_active">{getActivityDescription('moderately_active')}</option>
                <option value="very_active">{getActivityDescription('very_active')}</option>
                <option value="extra_active">{getActivityDescription('extra_active')}</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Goal</label>
              <select
                value={profile.goal}
                onChange={(e) => updateProfile('goal', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="lose_weight">{getGoalDescription('lose_weight')}</option>
                <option value="maintain">{getGoalDescription('maintain')}</option>
                <option value="gain_weight">{getGoalDescription('gain_weight')}</option>
                <option value="gain_muscle">{getGoalDescription('gain_muscle')}</option>
              </select>
            </div>
          </div>
        </div>

        <button
          onClick={calculateNutrition}
          disabled={loading}
          className="mt-6 w-full md:w-auto px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Calculating...' : 'Calculate Nutrition'}
        </button>
      </div>

      {/* Results */}
      {results && (
        <div className="bg-white rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Your Nutrition Results</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* BMR Card */}
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">Basal Metabolic Rate</p>
                  <p className="text-2xl font-bold text-blue-900">{results.bmr}</p>
                  <p className="text-sm text-blue-700">calories/day</p>
                </div>
                <div className="text-blue-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-blue-600 mt-2">Calories needed at complete rest</p>
            </div>

            {/* TDEE Card */}
            <div className="bg-green-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-green-600">Total Daily Energy</p>
                  <p className="text-2xl font-bold text-green-900">{results.tdee}</p>
                  <p className="text-sm text-green-700">calories/day</p>
                </div>
                <div className="text-green-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M11.3 1.046A1 1 0 0112 2v5h4a1 1 0 01.82 1.573l-7 10A1 1 0 018 18v-5H4a1 1 0 01-.82-1.573l7-10a1 1 0 011.12-.38z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-green-600 mt-2">Including daily activity</p>
            </div>

            {/* Target Calories Card */}
            <div className="bg-purple-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-purple-600">Target Calories</p>
                  <p className="text-2xl font-bold text-purple-900">{results.targetCalories}</p>
                  <p className="text-sm text-purple-700">calories/day</p>
                </div>
                <div className="text-purple-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                    <path fillRule="evenodd" d="M4 5a2 2 0 012-2 1 1 0 000 2H6a2 2 0 100 4h2a2 2 0 100 4h2a1 1 0 100 2H6a2 2 0 01-2-2V5z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-purple-600 mt-2">Based on your goal</p>
            </div>

            {/* Protein Card */}
            <div className="bg-red-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-red-600">Protein</p>
                  <p className="text-2xl font-bold text-red-900">{results.macros.protein}g</p>
                  <p className="text-sm text-red-700">per day</p>
                </div>
                <div className="text-red-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-red-600 mt-2">{Math.round((results.macros.protein * 4))} calories</p>
            </div>

            {/* Carbs Card */}
            <div className="bg-yellow-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-yellow-600">Carbohydrates</p>
                  <p className="text-2xl font-bold text-yellow-900">{results.macros.carbs}g</p>
                  <p className="text-sm text-yellow-700">per day</p>
                </div>
                <div className="text-yellow-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M2 10a8 8 0 018-8v8h8a8 8 0 11-16 0z" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-yellow-600 mt-2">{Math.round((results.macros.carbs * 4))} calories</p>
            </div>

            {/* Fat Card */}
            <div className="bg-indigo-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-indigo-600">Fat</p>
                  <p className="text-2xl font-bold text-indigo-900">{results.macros.fat}g</p>
                  <p className="text-sm text-indigo-700">per day</p>
                </div>
                <div className="text-indigo-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M5 2a1 1 0 011 1v1h1a1 1 0 010 2H6v1a1 1 0 01-2 0V6H3a1 1 0 010-2h1V3a1 1 0 011-1zm0 10a1 1 0 011 1v1h1a1 1 0 110 2H6v1a1 1 0 11-2 0v-1H3a1 1 0 110-2h1v-1a1 1 0 011-1zM12 2a1 1 0 01.967.744L14.146 7.2 17.5 9.134a1 1 0 010 1.732l-3.354 1.935-1.18 4.455a1 1 0 01-1.933 0L9.854 12.8 6.5 10.866a1 1 0 010-1.732l3.354-1.935 1.18-4.455A1 1 0 0112 2z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-indigo-600 mt-2">{Math.round((results.macros.fat * 9))} calories</p>
            </div>

            {/* Fiber Card */}
            <div className="bg-green-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-green-600">Fiber</p>
                  <p className="text-2xl font-bold text-green-900">{results.macros.fiber}g</p>
                  <p className="text-sm text-green-700">per day</p>
                </div>
                <div className="text-green-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-green-600 mt-2">For digestive health</p>
            </div>

            {/* Water Card */}
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">Water Intake</p>
                  <p className="text-2xl font-bold text-blue-900">{results.water}L</p>
                  <p className="text-sm text-blue-700">per day</p>
                </div>
                <div className="text-blue-200">
                  <svg className="h-8 w-8" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 2a1 1 0 00-1 1v1.323l-3.954 1.582 2.876 5.717A3 3 0 1013 14.836l2.876-5.717L11 4.323V3a1 1 0 00-1-1zm0 12a1 1 0 100 2 1 1 0 000-2z" clipRule="evenodd" />
                  </svg>
                </div>
              </div>
              <p className="text-xs text-blue-600 mt-2">35ml per kg body weight</p>
            </div>
          </div>

          {/* Recommendations */}
          <div className="mt-8 p-4 bg-gray-50 rounded-lg">
            <h3 className="text-lg font-medium text-gray-900 mb-3">Recommendations</h3>
            <div className="space-y-2 text-sm text-gray-700">
              <p>• Eat at least <span className="font-semibold">{results.macros.protein}g protein</span> daily for muscle maintenance</p>
              <p>• Include <span className="font-semibold">{results.macros.fiber}g fiber</span> for digestive health</p>
              <p>• Drink at least <span className="font-semibold">{results.water}L water</span> throughout the day</p>
              <p>• Spread calories across <span className="font-semibold">4-6 meals</span> for better metabolism</p>
              <p>• Adjust portions based on hunger and energy levels</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default NutritionCalculator;
