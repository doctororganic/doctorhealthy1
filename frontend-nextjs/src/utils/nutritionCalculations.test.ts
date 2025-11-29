import {
  calculateBMR,
  calculateTDEE,
  calculateTargetCalories,
  calculateMacros,
  calculateBMI,
  calculateWaterIntake,
  calculateAllNutritionMetrics,
  validateUserMetrics,
  convertToImperial,
  convertToMetric,
  type UserMetrics
} from './nutritionCalculations';

describe('Nutrition Calculations', () => {
  const testUserMale: UserMetrics = {
    age: 30,
    gender: 'male',
    weight: 80, // kg
    height: 180, // cm
    activityLevel: 'moderate',
    goal: 'maintain'
  };

  const testUserFemale: UserMetrics = {
    age: 25,
    gender: 'female',
    weight: 65, // kg
    height: 165, // cm
    activityLevel: 'light',
    goal: 'lose'
  };

  describe('calculateBMR', () => {
    it('should calculate BMR correctly for males', () => {
      const bmr = calculateBMR(testUserMale);
      // Expected: 10*80 + 6.25*180 - 5*30 + 5 = 800 + 1125 - 150 + 5 = 1780
      expect(bmr).toBe(1780);
    });

    it('should calculate BMR correctly for females', () => {
      const bmr = calculateBMR(testUserFemale);
      // Expected: 10*65 + 6.25*165 - 5*25 - 161 = 650 + 1031.25 - 125 - 161 = 1395.25
      expect(bmr).toBe(1395.25);
    });

    it('should handle edge cases', () => {
      const youngMale = { ...testUserMale, age: 18 };
      const bmr = calculateBMR(youngMale);
      expect(bmr).toBeGreaterThan(1700);

      const olderFemale = { ...testUserFemale, age: 60 };
      const bmrOlder = calculateBMR(olderFemale);
      expect(bmrOlder).toBeLessThan(1400);
    });
  });

  describe('calculateTDEE', () => {
    it('should calculate TDEE correctly for different activity levels', () => {
      const bmr = 1500;
      
      expect(calculateTDEE(bmr, 'sedentary')).toBe(1500 * 1.2);
      expect(calculateTDEE(bmr, 'light')).toBe(1500 * 1.375);
      expect(calculateTDEE(bmr, 'moderate')).toBe(1500 * 1.55);
      expect(calculateTDEE(bmr, 'active')).toBe(1500 * 1.725);
      expect(calculateTDEE(bmr, 'very_active')).toBe(1500 * 1.9);
    });
  });

  describe('calculateTargetCalories', () => {
    const tdee = 2000;

    it('should calculate target calories for weight loss', () => {
      expect(calculateTargetCalories(tdee, 'lose')).toBe(1500); // -500
    });

    it('should calculate target calories for maintenance', () => {
      expect(calculateTargetCalories(tdee, 'maintain')).toBe(2000);
    });

    it('should calculate target calories for weight gain', () => {
      expect(calculateTargetCalories(tdee, 'gain')).toBe(2300); // +300
    });
  });

  describe('calculateMacros', () => {
    const calories = 2000;
    const weight = 70;

    it('should calculate macros for weight loss goal', () => {
      const macros = calculateMacros(calories, 'lose', weight);
      
      expect(macros.calories).toBe(calories);
      expect(macros.proteinPercent).toBe(35);
      expect(macros.fatPercent).toBe(25);
      expect(macros.carbsPercent).toBe(40);
      
      // Protein should be at least 2.2g per kg for weight loss
      expect(macros.protein).toBeGreaterThanOrEqual(weight * 2.2);
    });

    it('should calculate macros for maintenance goal', () => {
      const macros = calculateMacros(calories, 'maintain', weight);
      
      expect(macros.proteinPercent).toBe(30);
      expect(macros.fatPercent).toBe(25);
      expect(macros.carbsPercent).toBe(45);
    });

    it('should calculate macros for weight gain goal', () => {
      const macros = calculateMacros(calories, 'gain', weight);
      
      expect(macros.proteinPercent).toBe(25);
      expect(macros.fatPercent).toBe(20);
      expect(macros.carbsPercent).toBe(55);
    });

    it('should calculate gram values correctly', () => {
      const macros = calculateMacros(2000, 'maintain', 70);
      
      // Protein: 30% of 2000 calories = 600 calories / 4 = 150g
      expect(macros.protein).toBeGreaterThanOrEqual(150);
      
      // Fat: 25% of 2000 calories = 500 calories / 9 = ~56g
      expect(macros.fat).toBeCloseTo(56, 0);
      
      // Carbs: 45% of 2000 calories = 900 calories / 4 = 225g
      expect(macros.carbs).toBeCloseTo(225, 0);
    });
  });

  describe('calculateBMI', () => {
    it('should calculate BMI and category correctly', () => {
      const { bmi, category } = calculateBMI(70, 175); // Normal weight
      expect(bmi).toBeCloseTo(22.9, 1);
      expect(category).toBe('Normal weight');
    });

    it('should categorize BMI correctly', () => {
      expect(calculateBMI(50, 175).category).toBe('Underweight'); // BMI ~16.3
      expect(calculateBMI(70, 175).category).toBe('Normal weight'); // BMI ~22.9
      expect(calculateBMI(85, 175).category).toBe('Overweight'); // BMI ~27.8
      expect(calculateBMI(100, 175).category).toBe('Obese'); // BMI ~32.7
    });
  });

  describe('calculateWaterIntake', () => {
    it('should calculate water intake based on weight and activity', () => {
      const weight = 70; // kg
      
      const sedentary = calculateWaterIntake(weight, 'sedentary');
      const active = calculateWaterIntake(weight, 'very_active');
      
      expect(sedentary).toBeCloseTo(2.5, 1); // 70 * 35 / 1000 = 2.45L
      expect(active).toBeGreaterThan(sedentary);
      expect(active).toBeCloseTo(3.4, 1); // 2.45 * 1.4 = 3.43L
    });
  });

  describe('calculateAllNutritionMetrics', () => {
    it('should calculate all metrics correctly for male user', () => {
      const results = calculateAllNutritionMetrics(testUserMale);
      
      expect(results.bmr).toBe(1780);
      expect(results.tdee).toBe(Math.round(1780 * 1.55)); // moderate activity
      expect(results.targetCalories).toBe(results.tdee); // maintain goal
      expect(results.macros).toBeDefined();
      expect(results.bmi).toBeCloseTo(24.7, 1);
      expect(results.bmiCategory).toBe('Normal weight');
      expect(results.waterIntake).toBeGreaterThan(2);
    });

    it('should calculate all metrics correctly for female user', () => {
      const results = calculateAllNutritionMetrics(testUserFemale);
      
      expect(results.bmr).toBe(1395.25);
      expect(results.tdee).toBe(Math.round(1395.25 * 1.375)); // light activity
      expect(results.targetCalories).toBe(results.tdee - 500); // lose goal
      expect(results.macros.proteinPercent).toBe(35); // weight loss
      expect(results.bmi).toBeCloseTo(23.9, 1);
      expect(results.bmiCategory).toBe('Normal weight');
    });
  });

  describe('validateUserMetrics', () => {
    it('should return no errors for valid metrics', () => {
      const errors = validateUserMetrics(testUserMale);
      expect(errors).toEqual([]);
    });

    it('should return errors for invalid age', () => {
      const invalidUser = { ...testUserMale, age: 10 };
      const errors = validateUserMetrics(invalidUser);
      expect(errors).toContain('Age must be between 15 and 100');
    });

    it('should return errors for invalid weight', () => {
      const invalidUser = { ...testUserMale, weight: 20 };
      const errors = validateUserMetrics(invalidUser);
      expect(errors).toContain('Weight must be between 30 and 300 kg');
    });

    it('should return errors for invalid height', () => {
      const invalidUser = { ...testUserMale, height: 100 };
      const errors = validateUserMetrics(invalidUser);
      expect(errors).toContain('Height must be between 120 and 250 cm');
    });

    it('should return errors for missing required fields', () => {
      const incompleteUser = { age: 25 } as Partial<UserMetrics>;
      const errors = validateUserMetrics(incompleteUser);
      
      expect(errors.length).toBeGreaterThan(0);
      expect(errors).toContain('Gender must be specified');
      expect(errors).toContain('Activity level must be specified');
      expect(errors).toContain('Goal must be specified');
    });
  });

  describe('Unit Conversion', () => {
    describe('convertToImperial', () => {
      it('should convert metric to imperial correctly', () => {
        const result = convertToImperial(70, 175);
        
        expect(result.weightLbs).toBeCloseTo(154.3, 1);
        expect(result.heightFeet).toBe(5);
        expect(result.heightInches).toBeCloseTo(9, 0);
      });
    });

    describe('convertToMetric', () => {
      it('should convert imperial to metric correctly', () => {
        const result = convertToMetric(154, 5, 9);
        
        expect(result.weight).toBeCloseTo(69.9, 1);
        expect(result.height).toBeCloseTo(175, 0);
      });
    });

    it('should be reversible conversions', () => {
      const originalWeight = 75;
      const originalHeight = 180;
      
      const imperial = convertToImperial(originalWeight, originalHeight);
      const backToMetric = convertToMetric(imperial.weightLbs, imperial.heightFeet, imperial.heightInches);
      
      expect(backToMetric.weight).toBeCloseTo(originalWeight, 0);
      expect(backToMetric.height).toBeCloseTo(originalHeight, 0);
    });
  });

  describe('Edge Cases and Error Handling', () => {
    it('should handle extreme but valid values', () => {
      const extremeUser: UserMetrics = {
        age: 80,
        gender: 'male',
        weight: 120,
        height: 200,
        activityLevel: 'very_active',
        goal: 'gain'
      };

      const results = calculateAllNutritionMetrics(extremeUser);
      
      expect(results.bmr).toBeGreaterThan(1000);
      expect(results.tdee).toBeGreaterThan(results.bmr);
      expect(results.targetCalories).toBeGreaterThan(results.tdee); // gain goal
      expect(results.bmi).toBeCloseTo(30, 0); // 120kg at 200cm
    });

    it('should handle minimum valid values', () => {
      const minUser: UserMetrics = {
        age: 15,
        gender: 'female',
        weight: 30,
        height: 120,
        activityLevel: 'sedentary',
        goal: 'maintain'
      };

      const results = calculateAllNutritionMetrics(minUser);
      
      expect(results.bmr).toBeGreaterThan(0);
      expect(results.tdee).toBeGreaterThan(results.bmr);
      expect(results.macros.protein).toBeGreaterThan(0);
      expect(results.waterIntake).toBeGreaterThan(0);
    });
  });
});