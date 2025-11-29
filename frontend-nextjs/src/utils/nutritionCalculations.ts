// Nutrition calculation utilities

export interface MacroResults {
  protein: number;
  carbs: number;
  fat: number;
  fiber: number;
}

/**
 * Calculate Basal Metabolic Rate (BMR) using Mifflin-St Jeor equation
 */
export function calculateBMR(weight: number, height: number, age: number, gender: 'male' | 'female'): number {
  if (gender === 'male') {
    return 10 * weight + 6.25 * height - 5 * age + 5;
  } else {
    return 10 * weight + 6.25 * height - 5 * age - 161;
  }
}

/**
 * Calculate Total Daily Energy Expenditure (TDEE)
 */
export function calculateTDEE(bmr: number, activityLevel: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extra_active'): number {
  const activityMultipliers = {
    sedentary: 1.2,
    lightly_active: 1.375,
    moderately_active: 1.55,
    very_active: 1.725,
    extra_active: 1.9
  };

  return bmr * activityMultipliers[activityLevel];
}

/**
 * Calculate macronutrient breakdown based on calories and ratios
 */
export function calculateMacros(
  calories: number, 
  proteinRatio: number = 0.3, 
  carbsRatio: number = 0.4, 
  fatRatio: number = 0.3
): MacroResults {
  const proteinCalories = calories * proteinRatio;
  const carbsCalories = calories * carbsRatio;
  const fatCalories = calories * fatRatio;

  // Convert calories to grams (4 cal/g for protein/carbs, 9 cal/g for fat)
  const protein = Math.round(proteinCalories / 4);
  const carbs = Math.round(carbsCalories / 4);
  const fat = Math.round(fatCalories / 9);
  
  // Fiber recommendation (14g per 1000 calories)
  const fiber = Math.max(25, Math.round(calories / 1000 * 14));

  return {
    protein,
    carbs,
    fat,
    fiber
  };
}

/**
 * Calculate ideal body weight using Devine formula
 */
export function calculateIdealWeight(height: number, gender: 'male' | 'female'): number {
  if (gender === 'male') {
    return 50 + 2.3 * ((height / 2.54) - 60);
  } else {
    return 45.5 + 2.3 * ((height / 2.54) - 60);
  }
}

/**
 * Calculate Body Mass Index (BMI)
 */
export function calculateBMI(weight: number, height: number): number {
  const heightInMeters = height / 100;
  return Math.round((weight / (heightInMeters * heightInMeters)) * 10) / 10;
}

/**
 * Get BMI category
 */
export function getBMICategory(bmi: number): string {
  if (bmi < 18.5) return 'Underweight';
  if (bmi < 25) return 'Normal weight';
  if (bmi < 30) return 'Overweight';
  return 'Obese';
}

/**
 * Calculate water intake recommendation (ml per kg body weight)
 */
export function calculateWaterIntake(weight: number): number {
  // 35ml per kg for average activity
  return Math.round((weight * 35) / 1000 * 10) / 10; // in liters
}

/**
 * Calculate calories needed for weight goals
 */
export function calculateCaloriesForGoal(
  tdee: number, 
  goal: 'lose_weight' | 'maintain' | 'gain_weight' | 'gain_muscle'
): number {
  const adjustments = {
    lose_weight: -500, // 500 calorie deficit
    maintain: 0,
    gain_weight: 500, // 500 calorie surplus
    gain_muscle: 300 // 300 calorie surplus
  };

  return tdee + adjustments[goal];
}

/**
 * Calculate protein needs based on activity level and goal
 */
export function calculateProteinNeeds(
  weight: number, 
  activityLevel: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extra_active',
  goal: 'lose_weight' | 'maintain' | 'gain_weight' | 'gain_muscle'
): number {
  let multiplier = 0.8; // Base sedentary

  // Activity multipliers
  const activityMultipliers = {
    sedentary: 0.8,
    lightly_active: 1.0,
    moderately_active: 1.2,
    very_active: 1.4,
    extra_active: 1.6
  };

  multiplier = activityMultipliers[activityLevel];

  // Goal adjustments
  if (goal === 'gain_muscle') {
    multiplier *= 1.5; // Higher protein for muscle gain
  } else if (goal === 'lose_weight') {
    multiplier *= 1.2; // Higher protein to preserve muscle during weight loss
  }

  return Math.round(weight * multiplier);
}

/**
 * Calculate recommended meal frequency and timing
 */
export function calculateMealFrequency(
  totalCalories: number,
  activityLevel: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extra_active'
): { meals: number; timing: string[] } {
  let meals = 3; // Base frequency

  // Adjust based on calories and activity
  if (totalCalories > 2500 || activityLevel === 'very_active' || activityLevel === 'extra_active') {
    meals = 5;
  } else if (totalCalories > 2000 || activityLevel === 'moderately_active') {
    meals = 4;
  }

  // Calculate timing (every 3-4 hours during waking hours)
  const timing = [];
  const startHour = 7; // 7 AM start
  const endHour = 21; // 9 PM end
  const interval = Math.floor((endHour - startHour) / meals);

  for (let i = 0; i < meals; i++) {
    timing.push(`${startHour + (i * interval)}:00`);
  }

  return { meals, timing };
}

/**
 * Validate nutrition inputs
 */
export function validateNutritionInputs(
  age: number,
  height: number,
  weight: number
): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (age < 15 || age > 100) {
    errors.push('Age must be between 15 and 100');
  }

  if (height < 100 || height > 250) {
    errors.push('Height must be between 100 and 250 cm');
  }

  if (weight < 30 || weight > 300) {
    errors.push('Weight must be between 30 and 300 kg');
  }

  return {
    isValid: errors.length === 0,
    errors
  };
}

/**
 * Calculate body fat percentage estimation (using BMI method)
 */
export function estimateBodyFat(
  bmi: number,
  age: number,
  gender: 'male' | 'female'
): number {
  if (gender === 'male') {
    return (1.20 * bmi) + (0.23 * age) - 16.2;
  } else {
    return (1.20 * bmi) + (0.23 * age) - 5.4;
  }
}

/**
 * Calculate recommended daily micronutrients based on calories
 */
export function calculateMicronutrients(calories: number): {
  vitamins: Record<string, number>;
  minerals: Record<string, number>;
} {
  // Base recommendations per 2000 calories
  const baseVitamins = {
    'Vitamin A': 900, // mcg
    'Vitamin C': 90, // mg
    'Vitamin D': 20, // mcg
    'Vitamin E': 15, // mg
    'Vitamin K': 120, // mcg
    'Thiamin (B1)': 1.2, // mg
    'Riboflavin (B2)': 1.3, // mg
    'Niacin (B3)': 16, // mg
    'Vitamin B6': 1.7, // mg
    'Folate (B9)': 400, // mcg
    'Vitamin B12': 2.4, // mcg
    'Biotin (B7)': 30, // mcg
    'Pantothenic Acid (B5)': 5 // mg
  };

  const baseMinerals = {
    'Calcium': 1000, // mg
    'Iron': 8, // mg (male) / 18 (female)
    'Magnesium': 420, // mg
    'Phosphorus': 700, // mg
    'Potassium': 3500, // mg
    'Sodium': 2300, // mg
    'Zinc': 11, // mg
    'Copper': 0.9, // mg
    'Manganese': 2.3, // mg
    'Selenium': 55, // mcg
    'Iodine': 150, // mcg
    'Chromium': 35 // mcg
  };

  // Scale based on calories (linear scaling)
  const scaleFactor = calories / 2000;

  const vitamins = Object.fromEntries(
    Object.entries(baseVitamins).map(([key, value]) => [key, Math.round(value * scaleFactor)])
  );

  const minerals = Object.fromEntries(
    Object.entries(baseMinerals).map(([key, value]) => [key, Math.round(value * scaleFactor)])
  );

  return { vitamins, minerals };
}
