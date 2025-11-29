'use client';

import React, { useState } from 'react';

export default function CalculatorPage() {
  const [age, setAge] = useState(25);
  const [weight, setWeight] = useState(70);
  const [height, setHeight] = useState(175);
  const [gender, setGender] = useState('male');
  const [activityLevel, setActivityLevel] = useState('moderate');
  const [goal, setGoal] = useState('maintain');
  const [results, setResults] = useState<any>(null);

  const calculateBMR = () => {
    // Mifflin-St Jeor equation
    if (gender === 'male') {
      return 10 * weight + 6.25 * height - 5 * age + 5;
    } else {
      return 10 * weight + 6.25 * height - 5 * age - 161;
    }
  };

  const calculateTDEE = (bmr: number) => {
    const multipliers = {
      sedentary: 1.2,
      light: 1.375,
      moderate: 1.55,
      active: 1.725,
      very_active: 1.9
    };
    return bmr * (multipliers as any)[activityLevel];
  };

  const calculateTargetCalories = (tdee: number) => {
    switch (goal) {
      case 'lose': return tdee - 500;
      case 'gain': return tdee + 300;
      default: return tdee;
    }
  };

  const calculate = () => {
    const bmr = calculateBMR();
    const tdee = calculateTDEE(bmr);
    const targetCalories = calculateTargetCalories(tdee);
    const bmi = weight / ((height / 100) ** 2);
    
    setResults({
      bmr: Math.round(bmr),
      tdee: Math.round(tdee),
      targetCalories: Math.round(targetCalories),
      bmi: Math.round(bmi * 10) / 10,
      protein: Math.round(targetCalories * 0.3 / 4),
      carbs: Math.round(targetCalories * 0.4 / 4),
      fat: Math.round(targetCalories * 0.3 / 9)
    });
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">🧮 Nutrition Calculator</h1>
      
      <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
        <div className="grid md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Age</label>
              <input
                type="number"
                value={age}
                onChange={(e) => setAge(Number(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Weight (kg)</label>
              <input
                type="number"
                value={weight}
                onChange={(e) => setWeight(Number(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Height (cm)</label>
              <input
                type="number"
                value={height}
                onChange={(e) => setHeight(Number(e.target.value))}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Gender</label>
              <select
                value={gender}
                onChange={(e) => setGender(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="male">Male</option>
                <option value="female">Female</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Activity Level</label>
              <select
                value={activityLevel}
                onChange={(e) => setActivityLevel(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="sedentary">Sedentary</option>
                <option value="light">Light</option>
                <option value="moderate">Moderate</option>
                <option value="active">Active</option>
                <option value="very_active">Very Active</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Goal</label>
              <select
                value={goal}
                onChange={(e) => setGoal(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="lose">Lose Weight</option>
                <option value="maintain">Maintain</option>
                <option value="gain">Gain Weight</option>
              </select>
            </div>
          </div>
        </div>
        
        <button
          onClick={calculate}
          className="mt-6 px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 w-full md:w-auto"
        >
          Calculate Nutrition Targets
        </button>
      </div>
      
      {results && (
        <div className="bg-gradient-to-r from-blue-50 to-green-50 rounded-lg p-6 border">
          <h3 className="text-xl font-bold text-gray-900 mb-4">Your Results</h3>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg p-4 text-center shadow-sm">
              <div className="text-2xl font-bold text-blue-600">{results.bmr}</div>
              <div className="text-sm font-medium text-gray-700">BMR</div>
            </div>
            <div className="bg-white rounded-lg p-4 text-center shadow-sm">
              <div className="text-2xl font-bold text-green-600">{results.tdee}</div>
              <div className="text-sm font-medium text-gray-700">TDEE</div>
            </div>
            <div className="bg-white rounded-lg p-4 text-center shadow-sm">
              <div className="text-2xl font-bold text-purple-600">{results.targetCalories}</div>
              <div className="text-sm font-medium text-gray-700">Target Calories</div>
            </div>
            <div className="bg-white rounded-lg p-4 text-center shadow-sm">
              <div className="text-2xl font-bold text-orange-600">{results.bmi}</div>
              <div className="text-sm font-medium text-gray-700">BMI</div>
            </div>
          </div>
          
          <div className="bg-white rounded-lg p-4">
            <h4 className="font-semibold text-gray-900 mb-3">Daily Macro Targets</h4>
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <div className="text-xl font-bold text-red-500">{results.protein}g</div>
                <div className="text-sm text-gray-700">Protein</div>
              </div>
              <div>
                <div className="text-xl font-bold text-yellow-500">{results.carbs}g</div>
                <div className="text-sm text-gray-700">Carbs</div>
              </div>
              <div>
                <div className="text-xl font-bold text-blue-500">{results.fat}g</div>
                <div className="text-sm text-gray-700">Fat</div>
              </div>
            </div>
          </div>
        </div>
      )}
      
      <div className="mt-8 bg-blue-50 rounded-lg p-6 border border-blue-100">
        <h3 className="text-lg font-semibold text-blue-900 mb-3">💡 About Your Results</h3>
        <div className="text-sm text-blue-800 space-y-2">
          <p><strong>BMR:</strong> The calories your body burns at rest</p>
          <p><strong>TDEE:</strong> Your total daily energy expenditure including activity</p>
          <p><strong>Target Calories:</strong> Adjusted for your weight goal</p>
          <p><strong>Macros:</strong> Recommended daily protein, carbs, and fat intake</p>
        </div>
      </div>
    </div>
  );
}