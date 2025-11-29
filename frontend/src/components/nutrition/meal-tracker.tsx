'use client';

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { 
  Search, 
  Plus, 
  X, 
  Utensils,
  Clock,
  Flame,
  Dumbbell
} from "lucide-react";
import { Meal, MealFood, Food } from "@/types";

interface MealTrackerProps {
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  onMealUpdate: (meal: Partial<Meal>) => void;
}

export function MealTracker({ mealType, onMealUpdate }: MealTrackerProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<Food[]>([]);
  const [selectedFoods, setSelectedFoods] = useState<MealFood[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  const mockFoods: Food[] = [
    {
      id: "1",
      name: "Apple",
      category: "Fruits",
      servingSize: 1,
      servingUnit: "medium",
      calories: 95,
      macros: {
        protein: 0.5,
        carbs: 25,
        fat: 0.3,
        fiber: 4
      },
      verified: true,
      source: "database",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: "2",
      name: "Grilled Chicken Breast",
      category: "Protein",
      servingSize: 100,
      servingUnit: "g",
      calories: 165,
      macros: {
        protein: 31,
        carbs: 0,
        fat: 3.6,
        fiber: 0
      },
      verified: true,
      source: "database",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: "3",
      name: "Brown Rice",
      category: "Grains",
      servingSize: 1,
      servingUnit: "cup",
      calories: 216,
      macros: {
        protein: 5,
        carbs: 45,
        fat: 1.8,
        fiber: 3.5
      },
      verified: true,
      source: "database",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ];

  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }

    setIsSearching(true);
    
    // Simulate API search delay
    setTimeout(() => {
      const filtered = mockFoods.filter(food => 
        food.name.toLowerCase().includes(query.toLowerCase()) ||
        food.category.toLowerCase().includes(query.toLowerCase())
      );
      setSearchResults(filtered);
      setIsSearching(false);
    }, 500);
  };

  const addFood = (food: Food, quantity: number = 1) => {
    const mealFood: MealFood = {
      id: Date.now().toString(),
      foodId: food.id,
      quantity,
      unit: food.servingUnit,
      calories: Math.round(food.calories * quantity),
      macros: {
        protein: Math.round(food.macros.protein * quantity * 10) / 10,
        carbs: Math.round(food.macros.carbs * quantity * 10) / 10,
        fat: Math.round(food.macros.fat * quantity * 10) / 10,
        fiber: food.macros.fiber ? Math.round(food.macros.fiber * quantity * 10) / 10 : 0,
      }
    };

    setSelectedFoods([...selectedFoods, mealFood]);
    setSearchQuery("");
    setSearchResults([]);
  };

  const removeFood = (foodId: string) => {
    setSelectedFoods(selectedFoods.filter(food => food.id !== foodId));
  };

  const calculateTotals = () => {
    return selectedFoods.reduce(
      (acc, food) => ({
        calories: acc.calories + food.calories,
        protein: acc.protein + food.macros.protein,
        carbs: acc.carbs + food.macros.carbs,
        fat: acc.fat + food.macros.fat,
        fiber: acc.fiber + (food.macros.fiber || 0),
      }),
      { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 }
    );
  };

  const saveMeal = () => {
    const totals = calculateTotals();
    const meal: Partial<Meal> = {
      mealType,
      foods: selectedFoods,
      totalCalories: totals.calories,
      totalMacros: {
        protein: Math.round(totals.protein * 10) / 10,
        carbs: Math.round(totals.carbs * 10) / 10,
        fat: Math.round(totals.fat * 10) / 10,
        fiber: Math.round(totals.fiber * 10) / 10,
        sugar: 0,
        sodium: 0,
      },
      mealDate: new Date().toISOString(),
    };

    onMealUpdate(meal);
    setSelectedFoods([]);
  };

  const totals = calculateTotals();
  const mealTypeLabels = {
    breakfast: "Breakfast",
    lunch: "Lunch",
    dinner: "Dinner",
    snack: "Snack"
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Utensils className="h-5 w-5" />
          {mealTypeLabels[mealType]}
        </CardTitle>
        <CardDescription>
          Add foods to track your {mealType.toLowerCase()} nutrition
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Food Search */}
        <div className="relative">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input
              placeholder="Search for foods..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              className="pl-10"
            />
          </div>

          {/* Search Results */}
          {searchResults.length > 0 && (
            <div className="absolute z-10 w-full mt-1 bg-white border rounded-lg shadow-lg max-h-60 overflow-y-auto">
              {searchResults.map((food) => (
                <div
                  key={food.id}
                  className="p-3 hover:bg-gray-50 cursor-pointer border-b last:border-b-0"
                  onClick={() => addFood(food)}
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <p className="font-medium">{food.name}</p>
                      <p className="text-sm text-gray-600">
                        {food.calories} cal per {food.servingSize} {food.servingUnit}
                      </p>
                    </div>
                    <Button size="sm" variant="outline">
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Selected Foods */}
        {selectedFoods.length > 0 && (
          <div className="space-y-4">
            <h4 className="font-medium">Selected Foods</h4>
            <div className="space-y-2">
              {selectedFoods.map((food) => (
                <div key={food.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex-1">
                    <p className="font-medium">{food.quantity} {food.unit}</p>
                    <p className="text-sm text-gray-600">
                      {food.calories} cal • P: {food.macros.protein}g • C: {food.macros.carbs}g • F: {food.macros.fat}g
                    </p>
                  </div>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => removeFood(food.id)}
                    className="text-red-500 hover:text-red-700"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              ))}
            </div>

            {/* Meal Totals */}
            <div className="border-t pt-4">
              <h4 className="font-medium mb-3">Meal Totals</h4>
              <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
                <div className="text-center">
                  <div className="flex items-center justify-center mb-1">
                    <Flame className="h-4 w-4 text-orange-500 mr-1" />
                    <span className="text-lg font-bold">{totals.calories}</span>
                  </div>
                  <p className="text-xs text-gray-600">Calories</p>
                </div>
                <div className="text-center">
                  <div className="text-lg font-bold text-red-500">{Math.round(totals.protein * 10) / 10}g</div>
                  <p className="text-xs text-gray-600">Protein</p>
                </div>
                <div className="text-center">
                  <div className="text-lg font-bold text-yellow-500">{Math.round(totals.carbs * 10) / 10}g</div>
                  <p className="text-xs text-gray-600">Carbs</p>
                </div>
                <div className="text-center">
                  <div className="text-lg font-bold text-blue-500">{Math.round(totals.fat * 10) / 10}g</div>
                  <p className="text-xs text-gray-600">Fat</p>
                </div>
                <div className="text-center">
                  <div className="text-lg font-bold text-green-500">{Math.round(totals.fiber * 10) / 10}g</div>
                  <p className="text-xs text-gray-600">Fiber</p>
                </div>
              </div>
            </div>

            <Button onClick={saveMeal} className="w-full">
              Save {mealTypeLabels[mealType]}
            </Button>
          </div>
        )}

        {/* Empty State */}
        {selectedFoods.length === 0 && searchQuery.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <Utensils className="h-12 w-12 mx-auto mb-3 text-gray-300" />
            <p>No foods added yet</p>
            <p className="text-sm">Search for foods above to get started</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
