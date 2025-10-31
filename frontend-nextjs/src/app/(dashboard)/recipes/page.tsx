'use client';

import { useState } from 'react';

// Types
interface UserProfile {
  name: string;
  cuisine: string;
  dietType: string;
  excludeIngredients: string[];
  calories: number;
}

interface Recipe {
  id: string;
  name: string;
  cuisine: string;
  dietType: string;
  ingredients: string[];
  instructions: string[];
  prepTime: number;
  cookTime: number;
  servings: number;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  isHalal: boolean;
  alternatives: {
    ingredient: string;
    halalAlternative: string;
  }[];
}

// Cuisine options
const cuisineOptions = [
  'American',
  'Italian',
  'Mexican',
  'Chinese',
  'Japanese',
  'Indian',
  'Mediterranean',
  'French',
  'Thai',
  'Middle Eastern',
  'Korean',
  'Vietnamese',
  'Greek',
  'Spanish',
  'Brazilian',
  'Moroccan'
];

// Diet types
const dietTypes = [
  'balanced',
  'low_carb',
  'keto',
  'mediterranean',
  'dash',
  'vegan',
  'vegetarian',
  'paleo',
  'anti_inflammatory',
  'high_carb'
];

// Haram ingredients and their halal alternatives
const haramIngredients: Record<string, string> = {
  'pork': 'beef',
  'lard': 'vegetable oil',
  'gelatin': 'agar-agar',
  'blood': 'iron supplements',
  'alcohol': 'fruit extracts',
  'carrion': 'halal meat',
  'carmine': 'beetroot juice',
  'shellac': 'carnauba wax',
  'wine vinegar': 'apple cider vinegar'
};

export default function RecipesPage() {
  const [userProfile, setUserProfile] = useState<UserProfile>({
    name: '',
    cuisine: '',
    dietType: 'balanced',
    excludeIngredients: [],
    calories: 2000
  });

  const [selectedCountry, setSelectedCountry] = useState<string>('');
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [selectedRecipe, setSelectedRecipe] = useState<Recipe | null>(null);
  const [showRecipes, setShowRecipes] = useState(false);
  const [loading, setLoading] = useState(false);

  // Generate recipes based on selected country and user preferences
  const generateRecipes = () => {
    setLoading(true);
    
    try {
      // Sample recipes - in a real implementation, this would come from a database
      const recipePool: Recipe[] = [
        {
          id: '1',
          name: 'Grilled Chicken Salad',
          cuisine: 'American',
          dietType: 'balanced',
          ingredients: ['chicken breast', 'lettuce', 'tomatoes', 'cucumber', 'olive oil', 'lemon juice', 'salt', 'pepper'],
          instructions: [
            'Season chicken breast with salt and pepper',
            'Grill chicken for 6-7 minutes per side',
            'Let chicken rest for 5 minutes',
            'Slice chicken into strips',
            'Combine lettuce, tomatoes, cucumber in a bowl',
            'Add chicken strips and olive oil',
            'Toss with lemon juice, salt, and pepper'
          ],
          prepTime: 15,
          cookTime: 15,
          servings: 2,
          calories: 350,
          protein: 35,
          carbs: 10,
          fat: 20,
          isHalal: true,
          alternatives: [
            { ingredient: 'chicken breast', halalAlternative: 'halal chicken breast' },
            { ingredient: 'olive oil', halalAlternative: 'vegetable oil' }
          ]
        },
        {
          id: '2',
          name: 'Spaghetti Carbonara',
          cuisine: 'Italian',
          dietType: 'balanced',
          ingredients: ['spaghetti', 'bacon', 'eggs', 'parmesan cheese', 'black pepper', 'salt', 'olive oil'],
          instructions: [
            'Cook spaghetti according to package directions',
            'Cook bacon until crispy',
            'Beat eggs with Parmesan cheese',
            'Add hot pasta to egg mixture',
            'Add bacon and olive oil',
            'Season with black pepper and salt'
          ],
          prepTime: 10,
          cookTime: 20,
          servings: 4,
          calories: 550,
          protein: 25,
          carbs: 60,
          fat: 25,
          isHalal: false,
          alternatives: [
            { ingredient: 'bacon', halalAlternative: 'turkey bacon' },
            { ingredient: 'parmesan cheese', halalAlternative: 'halal parmesan cheese' }
          ]
        },
        {
          id: '3',
          name: 'Chicken Tikka Masala',
          cuisine: 'Indian',
          dietType: 'balanced',
          ingredients: ['chicken', 'yogurt', 'tomatoes', 'onions', 'garlic', 'ginger', 'tikka masala', 'cream', 'rice', 'cilantro'],
          instructions: [
            'Marinate chicken in yogurt and tikka masala',
            'Cook chicken in a pan until golden',
            'Sauté onions, garlic, and ginger',
            'Add tomatoes and cook until soft',
            'Add chicken and cream',
            'Simmer for 20 minutes',
            'Garnish with cilantro',
            'Serve with rice'
          ],
          prepTime: 30,
          cookTime: 30,
          servings: 4,
          calories: 450,
          protein: 35,
          carbs: 40,
          fat: 15,
          isHalal: true,
          alternatives: [
            { ingredient: 'cream', halalAlternative: 'coconut cream' }
          ]
        },
        {
          id: '4',
          name: 'Sushi Roll',
          cuisine: 'Japanese',
          dietType: 'balanced',
          ingredients: ['sushi rice', 'nori', 'salmon', 'cucumber', 'avocado', 'rice vinegar', 'sugar', 'salt', 'wasabi', 'soy sauce'],
          instructions: [
            'Cook sushi rice and season with rice vinegar',
            'Place nori on bamboo mat',
            'Spread rice evenly on nori',
            'Add salmon, cucumber, and avocado',
            'Roll tightly using bamboo mat',
            'Cut into 8 pieces',
            'Serve with wasabi and soy sauce'
          ],
          prepTime: 30,
          cookTime: 10,
          servings: 2,
          calories: 300,
          protein: 20,
          carbs: 40,
          fat: 10,
          isHalal: false,
          alternatives: [
            { ingredient: 'salmon', halalAlternative: 'halal salmon' },
            { ingredient: 'soy sauce', halalAlternative: 'halal soy sauce' }
          ]
        },
        {
          id: '5',
          name: 'Mediterranean Quinoa Bowl',
          cuisine: 'Mediterranean',
          dietType: 'balanced',
          ingredients: ['quinoa', 'chickpeas', 'cucumber', 'tomatoes', 'red onion', 'feta cheese', 'olives', 'olive oil', 'lemon juice', 'oregano'],
          instructions: [
            'Cook quinoa according to package directions',
            'Rinse and drain chickpeas',
            'Chop cucumber, tomatoes, and red onion',
            'Combine quinoa, chickpeas, and vegetables',
            'Add feta cheese and olives',
            'Drizzle with olive oil and lemon juice',
            'Season with oregano'
          ],
          prepTime: 20,
          cookTime: 15,
          servings: 2,
          calories: 420,
          protein: 15,
          carbs: 55,
          fat: 18,
          isHalal: false,
          alternatives: [
            { ingredient: 'feta cheese', halalAlternative: 'halal feta cheese' }
          ]
        },
        {
          id: '6',
          name: 'Korean Bibimbap',
          cuisine: 'Korean',
          dietType: 'balanced',
          ingredients: ['rice', 'spinach', 'bean sprouts', 'carrots', 'cucumbers', 'mushrooms', 'egg', 'gochujang', 'sesame oil', 'sesame seeds'],
          instructions: [
            'Cook rice according to package directions',
            'Blanch spinach and bean sprouts',
            'Julienne carrots and cucumbers',
            'Slice mushrooms',
            'Fry egg sunny-side up',
            'Arrange vegetables over rice',
            'Top with fried egg',
            'Serve with gochujang and sesame oil',
            'Sprinkle with sesame seeds'
          ],
          prepTime: 30,
          cookTime: 20,
          servings: 2,
          calories: 480,
          protein: 18,
          carbs: 65,
          fat: 15,
          isHalal: true,
          alternatives: []
        }
      ];

      // Filter recipes by cuisine
      const filteredByCuisine = selectedCountry 
        ? recipePool.filter(recipe => recipe.cuisine === selectedCountry)
        : recipePool;

      // Filter by diet type
      const filteredByDiet = filteredByCuisine.filter(recipe => 
        recipe.dietType === userProfile.dietType
      );

      // Filter by excluded ingredients
      const filteredByIngredients = filteredByDiet.filter(recipe => {
        return !recipe.ingredients.some(ingredient => 
          userProfile.excludeIngredients.some(excluded => 
            ingredient.toLowerCase().includes(excluded.toLowerCase())
          )
        );
      });

      // Make recipes halal by substituting haram ingredients
      const halalRecipes = filteredByIngredients.map(recipe => {
        let halalRecipe = { ...recipe };
        let isHalal = true;

        // Check for haram ingredients
        recipe.ingredients.forEach(ingredient => {
          const ingredientLower = ingredient.toLowerCase();
          if (haramIngredients[ingredientLower]) {
            isHalal = false;
            
            // Replace haram ingredient with halal alternative
            halalRecipe.ingredients = halalRecipe.ingredients.map(ing => {
              return ing.toLowerCase() === ingredientLower 
                ? haramIngredients[ingredientLower]
                : ing;
            });
          }
        });

        // Check diet type for haram ingredients
        if (userProfile.dietType === 'halal' || !isHalal) {
          halalRecipe.isHalal = true;
        }

        return halalRecipe;
      });

      setRecipes(halalRecipes);
      setShowRecipes(true);
    } catch (error) {
      console.error('Failed to generate recipes:', error);
    } finally {
      setLoading(false);
    }
  };

  // Handle country selection
  const handleCountrySelect = (country: string) => {
    setSelectedCountry(country);
    setUserProfile({...userProfile, cuisine: country});
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-yellow-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h1 className="text-2xl font-bold text-green-600 mb-6">Recipes and Review</h1>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Your Name</label>
              <input
                type="text"
                value={userProfile.name}
                onChange={(e) => setUserProfile({...userProfile, name: e.target.value})}
                className="w-full p-2 border border-gray-300 rounded-md"
              />
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Diet Type</label>
                <select
                  value={userProfile.dietType}
                  onChange={(e) => setUserProfile({...userProfile, dietType: e.target.value})}
                  className="w-full p-2 border border-gray-300 rounded-md"
                >
                  {dietTypes.map(diet => (
                    <option key={diet} value={diet}>
                      {diet.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                    </option>
                  ))}
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Daily Calories</label>
                <input
                  type="number"
                  value={userProfile.calories}
                  onChange={(e) => {
                    const num = parseInt(e.target.value);
                    setUserProfile({...userProfile, calories: isNaN(num) ? userProfile.calories : Math.max(500, num)});
                  }}
                  className="w-full p-2 border border-gray-300 rounded-md"
                />
              </div>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Exclude Ingredients</label>
              <textarea
                value={userProfile.excludeIngredients.join(', ')}
                onChange={(e) => setUserProfile({...userProfile, excludeIngredients: e.target.value.split(',').map(i => i.trim())})}
                className="w-full p-2 border border-gray-300 rounded-md"
                rows={2}
                placeholder="e.g., pork, alcohol, shellfish"
              />
            </div>
            
            <button
              onClick={generateRecipes}
              disabled={loading}
              className="w-full btn-primary"
            >
              {loading ? 'Generating...' : 'Generate Recipes'}
            </button>
          </div>
        </div>
        
        {/* Cuisine Selection */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <h2 className="text-xl font-bold text-green-600 mb-4">Select Cuisine</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
            {cuisineOptions.map(country => (
              <button
                key={country}
                onClick={() => handleCountrySelect(country)}
                className={`p-3 rounded-lg border-2 transition-all ${
                  selectedCountry === country
                    ? 'border-green-500 bg-green-50'
                    : 'border-gray-200 hover:border-green-300 hover:bg-green-50'
                }`}
              >
                <div className="text-sm font-medium text-gray-800">{country}</div>
              </button>
            ))}
          </div>
          
          {selectedCountry && (
            <div className="mt-4 p-4 bg-green-50 rounded-lg border border-green-200">
              <p className="text-sm text-green-800">
                Selected Cuisine: <span className="font-semibold">{selectedCountry}</span>
              </p>
            </div>
          )}
        </div>
        
        {/* Recipes Results */}
        {showRecipes && (
          <div className="bg-white rounded-xl shadow-lg p-6">
            <h2 className="text-xl font-bold text-green-600 mb-4">
              Recipes{selectedCountry && ` from ${selectedCountry}`}
            </h2>
            
            <div className="space-y-4">
              {recipes.map((recipe) => (
                <div
                  key={recipe.id}
                  className="recipe-box cursor-pointer"
                  onClick={() => setSelectedRecipe(recipe)}
                >
                  <div className="flex justify-between items-start mb-2">
                    <div>
                      <h3 className="font-semibold text-gray-800">{recipe.name}</h3>
                      <div className="flex items-center gap-2 text-sm text-gray-600">
                        <span>{recipe.cuisine}</span>
                        <span>•</span>
                        <span>{recipe.dietType.replace('_', ' ')}</span>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {recipe.isHalal && (
                        <span className="halal-badge">Halal</span>
                      )}
                      <div className="text-right text-sm">
                        <div>{recipe.calories} cal</div>
                        <div className="text-gray-600">P: {recipe.protein}g | C: {recipe.carbs}g | F: {recipe.fat}g</div>
                      </div>
                    </div>
                  </div>
                  
                  <div className="text-sm text-gray-600">
                    <div className="mb-1">
                      <span className="font-medium">Time:</span> {recipe.prepTime + recipe.cookTime} min
                      <span className="ml-4"><span className="font-medium">Servings:</span> {recipe.servings}</span>
                    </div>
                    <div className="mb-1">
                      <span className="font-medium">Ingredients:</span> {recipe.ingredients.join(', ')}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
        
        {/* Recipe Details Modal */}
        {selectedRecipe && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-lg max-w-2xl w-full max-h-screen overflow-y-auto p-6">
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-xl font-bold text-gray-800">{selectedRecipe.name}</h3>
                <button
                  onClick={() => setSelectedRecipe(null)}
                  className="text-gray-500 hover:text-gray-700"
                >
                  ✕
                </button>
              </div>
              
              <div className="mb-4">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-medium">Nutrition:</span>
                  <span>{selectedRecipe.calories} cal | P: {selectedRecipe.protein}g | C: {selectedRecipe.carbs}g | F: {selectedRecipe.fat}g</span>
                </div>
                <div className="flex items-center gap-4 text-sm">
                  <span><span className="font-medium">Cuisine:</span> {selectedRecipe.cuisine}</span>
                  <span><span className="font-medium">Diet:</span> {selectedRecipe.dietType.replace('_', ' ')}</span>
                  <span><span className="font-medium">Time:</span> {selectedRecipe.prepTime + selectedRecipe.cookTime} min</span>
                  <span><span className="font-medium">Servings:</span> {selectedRecipe.servings}</span>
                  {selectedRecipe.isHalal && <span className="halal-badge">Halal</span>}
                </div>
              </div>
              
              <div className="mb-4">
                <h4 className="font-semibold text-gray-800 mb-2">Ingredients</h4>
                <div className="flex flex-wrap gap-2">
                  {selectedRecipe.ingredients.map((ingredient, index) => (
                    <span key={index} className="bg-gray-100 px-2 py-1 rounded text-sm">
                      {ingredient}
                    </span>
                  ))}
                </div>
              </div>
              
              <div className="mb-4">
                <h4 className="font-semibold text-gray-800 mb-2">Instructions</h4>
                <ol className="list-decimal list-inside space-y-2 text-gray-600">
                  {selectedRecipe.instructions.map((instruction, index) => (
                    <li key={index}>{instruction}</li>
                  ))}
                </ol>
              </div>
              
              {/* Halal Alternatives */}
              {selectedRecipe.alternatives.length > 0 && (
                <div className="border-t pt-4">
                  <h4 className="font-semibold text-gray-800 mb-2">Halal Alternatives</h4>
                  <div className="space-y-2">
                    {selectedRecipe.alternatives.map((alt, index) => (
                      <div key={index} className="text-sm">
                        <span className="font-medium">{alt.ingredient}:</span> {alt.halalAlternative}
                      </div>
                    ))}
                  </div>
                </div>
              )}
              
              <div className="mt-6 flex justify-end">
                <button
                  onClick={() => setSelectedRecipe(null)}
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
            <strong>Halal Information:</strong> All recipes are designed to avoid haram ingredients and provide halal alternatives when necessary. However, always verify that ingredients meet halal certification standards, especially when dining out or using pre-packaged products.
          </p>
          <p>
            <strong>Allergen Information:</strong> Recipes may contain common allergens. If you have food allergies or sensitivities, please review the ingredient list carefully and make appropriate substitutions.
          </p>
        </div>
      </div>
    </div>
  );
}