'use client';

import { useState, useEffect } from 'react';
import { z } from 'zod';
import { nutritionService } from '@/lib/api/services/nutrition.service';
import { loggers } from '@/lib/logger';

// Form schema
const userProfileSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  age: z.number().min(1, 'Age must be at least 1'),
  weight: z.number().min(1, 'Weight must be greater than 0'),
  height: z.number().min(1, 'Height must be greater than 0'),
  activityLevel: z.enum(['sedentary', 'light', 'moderate', 'active', 'very_active']),
  goal: z.enum(['lose_weight', 'gain_weight', 'maintain_weight', 'gain_muscle', 'reshape']),
  metabolicRate: z.enum(['low', 'medium', 'high']),
  excludeIngredients: z.array(z.string()).default([]),
  diseases: z.array(z.string()).default([]),
  medications: z.array(z.string()).default([]),
});

type UserProfile = z.infer<typeof userProfileSchema>;

interface NutritionPlan {
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  equation: string;
  meals: Meal[];
}

interface Meal {
  name: string;
  components: string[];
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  preparation: string;
  alternative: {
    name: string;
    components: string[];
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
    preparation: string;
  };
}

export default function MealsPage() {
  const [userProfile, setUserProfile] = useState<UserProfile>({
    name: '',
    age: 30,
    weight: 70,
    height: 170,
    activityLevel: 'moderate',
    goal: 'maintain_weight',
    metabolicRate: 'medium',
    excludeIngredients: [],
    diseases: [],
    medications: [],
  });

  const [nutritionPlan, setNutritionPlan] = useState<NutritionPlan | null>(null);
  const [loading, setLoading] = useState(false);
  const [showPlan, setShowPlan] = useState(false);
  const [selectedMeal, setSelectedMeal] = useState<Meal | null>(null);

  // Calculate BMI
  const bmi = (userProfile.weight / Math.pow(userProfile.height / 100, 2)).toFixed(1);

  // Calculate nutrition requirements
  const calculateNutrition = () => {
    const { weight, height, activityLevel, goal, metabolicRate } = userProfile;
    
    let caloriesPerKg = 20; // Default
    let equation = 'Standard formula: 20 calories per kg';
    
    // BMI-based adjustment
    const bmiValue = parseFloat(bmi);
    if (bmiValue >= 18 && bmiValue <= 30) {
      caloriesPerKg = 20;
      equation = 'Standard formula: 20 calories per kg (BMI 18-30)';
    } else if (bmiValue >= 15 && bmiValue < 18) {
      caloriesPerKg = 25;
      equation = 'Weight gain formula: 25 calories per kg (BMI 15-17.9)';
    }
    
    // Goal-based adjustment
    if (goal === 'gain_muscle' || goal === 'gain_weight') {
      caloriesPerKg = 30;
      equation = 'Muscle gain formula: 30 calories per kg';
    } else if (goal === 'lose_weight') {
      caloriesPerKg = 18;
      equation = 'Weight loss formula: 18 calories per kg';
    }
    
    // Activity level adjustment
    let activityMultiplier = 1;
    switch (activityLevel) {
      case 'sedentary':
        activityMultiplier = 1.2;
        break;
      case 'light':
        activityMultiplier = 1.375;
        break;
      case 'moderate':
        activityMultiplier = 1.55;
        break;
      case 'active':
        activityMultiplier = 1.725;
        break;
      case 'very_active':
        activityMultiplier = 1.9;
        break;
    }
    
    // Calculate base calories
    const baseCalories = Math.round(weight * caloriesPerKg * activityMultiplier);
    
    // Calculate protein (1-1.5g/kg for regular, 1.5-1.7g/kg for active)
    const proteinPerKg = activityLevel === 'active' || activityLevel === 'very_active' ? 1.6 : 1.2;
    const protein = Math.round(weight * proteinPerKg);
    
    // Calculate macros (40% carbs, 30% protein, 30% fat)
    const proteinCalories = protein * 4;
    const fatCalories = baseCalories * 0.3;
    const carbCalories = baseCalories - proteinCalories - fatCalories;
    
    const fat = Math.round(fatCalories / 9);
    const carbs = Math.round(carbCalories / 4);
    
    return {
      calories: baseCalories,
      protein,
      carbs,
      fat,
      equation,
    };
  };

  // Generate meal plan
  const generateMealPlan = async () => {
    setLoading(true);
    
    try {
      // Validate user profile
      const validatedProfile = userProfileSchema.parse(userProfile);
      
      // Calculate nutrition requirements
      const nutrition = calculateNutrition();
      
      // Generate meals
      const meals = generateMeals(nutrition);
      
      setNutritionPlan({
        ...nutrition,
        meals,
      });
      
      setShowPlan(true);
      loggers.nutrition.info('Meal plan generated successfully', {
        userProfile: validatedProfile,
        nutrition,
      });
    } catch (error) {
      loggers.nutrition.error('Failed to generate meal plan', { error });
    } finally {
      setLoading(false);
    }
  };

  // Generate meals based on nutrition requirements
  const generateMeals = (nutrition: { calories: number; protein: number; carbs: number; fat: number }): Meal[] => {
    const { calories, protein, carbs, fat } = nutrition;
    const mealCalories = calories / 4; // 4 meals
    
    // Sample meal data - in a real implementation, this would come from a database
    const sampleMeals: Omit<Meal, 'calories' | 'protein' | 'carbs' | 'fat'>[] = [
      {
        name: 'Grilled Chicken Salad',
        components: ['200g grilled chicken breast', '100g mixed greens', '50g cherry tomatoes', '1 tbsp olive oil'],
        preparation: 'Grill chicken breast until cooked through. Combine with mixed greens, cherry tomatoes, and olive oil.',
        alternative: {
          name: 'Baked Salmon Bowl',
          components: ['200g baked salmon', '100g quinoa', '50g steamed broccoli', '1 tbsp lemon juice'],
          preparation: 'Bake salmon at 400°F for 15 minutes. Serve over cooked quinoa with steamed broccoli and lemon juice.',
        },
      },
      {
        name: 'Vegetarian Protein Bowl',
        components: ['150g chickpeas', '100g brown rice', '50g avocado', '50g pumpkin seeds'],
        preparation: 'Cook chickpeas until tender. Combine with cooked brown rice, sliced avocado, and pumpkin seeds.',
        alternative: {
          name: 'Tofu Stir-Fry',
          components: ['200g firm tofu', '100g mixed vegetables', '50g brown rice', '1 tbsp soy sauce'],
          preparation: 'Cube tofu and stir-fry with mixed vegetables. Serve over brown rice with soy sauce.',
        },
      },
      {
        name: 'Greek Yogurt Parfait',
        components: ['200g Greek yogurt', '50g granola', '50g mixed berries', '1 tbsp honey'],
        preparation: 'Layer Greek yogurt, granola, and mixed berries in a glass. Drizzle with honey.',
        alternative: {
          name: 'Cottage Cheese Bowl',
          components: ['200g cottage cheese', '50g sliced almonds', '50g apple slices', '1 tsp cinnamon'],
          preparation: 'Combine cottage cheese with sliced almonds and apple slices. Sprinkle with cinnamon.',
        },
      },
      {
        name: 'Lean Beef Wrap',
        components: ['150g lean ground beef', '100g whole wheat tortilla', '50g lettuce', '1 tbsp salsa'],
        preparation: 'Cook ground beef until browned. Wrap in tortilla with lettuce and salsa.',
        alternative: {
          name: 'Turkey Lettuce Wraps',
          components: ['150g sliced turkey breast', '100g large lettuce leaves', '50g cucumber', '1 tbsp hummus'],
          preparation: 'Wrap turkey slices in lettuce leaves with cucumber and hummus.',
        },
      },
    ];
    
    // Calculate nutrition for each meal
    return sampleMeals.map((meal, index) => {
      // In a real implementation, this would calculate based on the actual ingredients
      // For now, we'll distribute the total nutrition evenly
      const mealCalories = Math.round(calories / 4);
      const mealProtein = Math.round(protein / 4);
      const mealCarbs = Math.round(carbs / 4);
      const mealFat = Math.round(fat / 4);
      
      return {
        ...meal,
        calories: mealCalories,
        protein: mealProtein,
        carbs: mealCarbs,
        fat: mealFat,
      };
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-yellow-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h1 className="text-2xl font-bold text-green-600 mb-6">Meals and Body Enhancing</h1>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* User Profile Form */}
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-gray-800">Your Profile</h2>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                  type="text"
                  value={userProfile.name}
                  onChange={(e) => setUserProfile({...userProfile, name: e.target.value})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                />
              </div>
              
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Age</label>
                  <input
                    type="number"
                    value={userProfile.age}
                    onChange={(e) => {
                      const num = parseInt(e.target.value);
                      setUserProfile({...userProfile, age: isNaN(num) ? userProfile.age : Math.max(1, num)});
                    }}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Weight (kg)</label>
                  <input
                    type="number"
                    value={userProfile.weight}
                    onChange={(e) => {
                      const num = parseFloat(e.target.value);
                      setUserProfile({...userProfile, weight: isNaN(num) ? userProfile.weight : Math.max(1, num)});
                    }}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Height (cm)</label>
                  <input
                    type="number"
                    value={userProfile.height}
                    onChange={(e) => {
                      const num = parseInt(e.target.value);
                      setUserProfile({...userProfile, height: isNaN(num) ? userProfile.height : Math.max(1, num)});
                    }}
                    className="w-full p-2 border border-gray-300 rounded-md"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Activity Level</label>
                <select
                  value={userProfile.activityLevel}
                  onChange={(e) => setUserProfile({...userProfile, activityLevel: e.target.value as any})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                >
                  <option value="sedentary">Sedentary</option>
                  <option value="light">Light</option>
                  <option value="moderate">Moderate</option>
                  <option value="active">Active</option>
                  <option value="very_active">Very Active</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Goal</label>
                <select
                  value={userProfile.goal}
                  onChange={(e) => setUserProfile({...userProfile, goal: e.target.value as any})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                >
                  <option value="lose_weight">Lose Weight</option>
                  <option value="gain_weight">Gain Weight</option>
                  <option value="maintain_weight">Maintain Weight</option>
                  <option value="gain_muscle">Gain Muscle</option>
                  <option value="reshape">Reshape Body</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Metabolic Rate</label>
                <select
                  value={userProfile.metabolicRate}
                  onChange={(e) => setUserProfile({...userProfile, metabolicRate: e.target.value as any})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                >
                  <option value="low">Low</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Exclude Ingredients</label>
                <textarea
                  value={userProfile.excludeIngredients.join(', ')}
                  onChange={(e) => setUserProfile({...userProfile, excludeIngredients: e.target.value.split(',').map(i => i.trim())})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                  rows={2}
                  placeholder="e.g., peanuts, shellfish, dairy"
                />
              </div>
              
              <button
                onClick={generateMealPlan}
                disabled={loading}
                className="w-full btn-primary"
              >
                {loading ? 'Generating...' : 'Generate Meal Plan'}
              </button>
            </div>
            
            {/* BMI Display */}
            <div className="space-y-4">
              <h2 className="text-xl font-semibold text-gray-800">Your Metrics</h2>
              <div className="bg-gradient-to-r from-green-50 to-blue-50 rounded-lg p-4 border border-green-200">
                <div className="mb-2">
                  <span className="font-medium">BMI:</span> {bmi}
                </div>
                <div className="text-sm text-gray-600">
                  {parseFloat(bmi) < 18.5 && 'Underweight'}
                  {parseFloat(bmi) >= 18.5 && parseFloat(bmi) < 25 && 'Normal weight'}
                  {parseFloat(bmi) >= 25 && parseFloat(bmi) < 30 && 'Overweight'}
                  {parseFloat(bmi) >= 30 && 'Obese'}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        {/* Nutrition Plan Results */}
        {showPlan && nutritionPlan && (
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-bold text-green-600 mb-4">Your Nutrition Plan</h2>
            
            <div className="calculator-result mb-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-2">Daily Requirements</h3>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">{nutritionPlan.calories}</div>
                  <div className="text-sm text-gray-600">Calories</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-green-600">{nutritionPlan.protein}g</div>
                  <div className="text-sm text-gray-600">Protein</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-yellow-600">{nutritionPlan.carbs}g</div>
                  <div className="text-sm text-gray-600">Carbs</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-red-600">{nutritionPlan.fat}g</div>
                  <div className="text-sm text-gray-600">Fat</div>
                </div>
              </div>
              <div className="mt-2 text-sm text-gray-600 italic">
                Calculation method: {nutritionPlan.equation}
              </div>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Daily Meal Plan</h3>
              <div className="space-y-4">
                {nutritionPlan.meals.map((meal, index) => (
                  <div key={index} className="meal-card" onClick={() => setSelectedMeal(meal)}>
                    <div className="flex justify-between items-start">
                      <div>
                        <h4 className="font-semibold text-gray-800">{meal.name}</h4>
                        <div className="text-sm text-gray-600 mt-1">
                          {meal.components.join(', ')}
                        </div>
                      </div>
                      <div className="text-right text-sm">
                        <div>{meal.calories} cal</div>
                        <div className="text-gray-600">P: {meal.protein}g | C: {meal.carbs}g | F: {meal.fat}g</div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
        
        {/* Meal Details Modal */}
        {selectedMeal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-lg max-w-2xl w-full max-h-screen overflow-y-auto p-6">
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-xl font-bold text-gray-800">{selectedMeal.name}</h3>
                <button
                  onClick={() => setSelectedMeal(null)}
                  className="text-gray-500 hover:text-gray-700"
                >
                  ✕
                </button>
              </div>
              
              <div className="mb-4">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-medium">Nutrition:</span>
                  <span>{selectedMeal.calories} cal | P: {selectedMeal.protein}g | C: {selectedMeal.carbs}g | F: {selectedMeal.fat}g</span>
                </div>
                <div className="text-sm text-gray-600">
                  Ingredients: {selectedMeal.components.join(', ')}
                </div>
              </div>
              
              <div className="mb-4">
                <h4 className="font-semibold text-gray-800 mb-2">Preparation</h4>
                <p className="text-gray-600">{selectedMeal.preparation}</p>
              </div>
              
              <div className="border-t pt-4">
                <h4 className="font-semibold text-gray-800 mb-2">Alternative Option</h4>
                <h5 className="font-medium text-gray-700">{selectedMeal.alternative.name}</h5>
                <div className="flex justify-between items-center mb-2 text-sm">
                  <span>Nutrition:</span>
                  <span>{selectedMeal.alternative.calories} cal | P: {selectedMeal.alternative.protein}g | C: {selectedMeal.alternative.carbs}g | F: {selectedMeal.alternative.fat}g</span>
                </div>
                <div className="text-sm text-gray-600 mb-2">
                  Ingredients: {selectedMeal.alternative.components.join(', ')}
                </div>
                <p className="text-gray-600">{selectedMeal.alternative.preparation}</p>
              </div>
              
              <div className="mt-6 flex justify-end">
                <button
                  onClick={() => setSelectedMeal(null)}
                  className="btn-primary"
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        )}
        
        {/* Medical Disclaimer */}
        <div className="disclaimer mt-8">
          <p>
            <strong>Disclaimer:</strong> This site is to update you with information and guide you to useful advice for the purpose of education and awareness and is not a substitute for a doctor's visit.
          </p>
          <p>
            <strong>Halal Information:</strong> All meal plans are designed to avoid haram ingredients. If you have specific dietary restrictions, please ensure all ingredients meet halal certification standards.
          </p>
        </div>
      </div>
    </div>
  );
}