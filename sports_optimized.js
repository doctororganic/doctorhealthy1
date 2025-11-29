/**
 * Comprehensive Sports Nutrition Data
 * Contains detailed nutritional guidelines for various sports
 * Last Updated: 2025-09-13
 */

const sportsNutritionData = {
  "general_guidelines": {
    "macronutrient_needs": {
      "endurance_sports": {
        "carbohydrates": "8-12g/kg/day",
        "protein": "1.6-1.8g/kg/day",
        "fat": "25-30% of total calories"
      },
      "strength_sports": {
        "carbohydrates": "4-7g/kg/day",
        "protein": "1.6-2.2g/kg/day",
        "fat": "25-30% of total calories"
      },
      "team_sports": {
        "carbohydrates": "6-8g/kg/day",
        "protein": "1.8-2.2g/kg/day",
        "fat": "20-25% of total calories"
      }
    },
    "nutrient_timing": {
      "pre_exercise_2_4hrs": {
        "carbohydrates": "4-7g/kg",
        "protein": "0.3g/kg",
        "fat_fiber": "Low"
      },
      "pre_exercise_30_60min": {
        "carbohydrates": "30-60g",
        "protein": "10-20g",
        "fluids": "400-600ml"
      },
      "during_exercise": {
        "carbohydrates": "30-60g/hr (up to 90g/hr for >2.5hrs)",
        "fluids": "400-800ml/hr",
        "electrolytes": "Essential"
      },
      "post_exercise_0_2hrs": {
        "carbohydrates": "1.2g/kg",
        "protein": "0.3g/kg",
        "ratio": "3:1 carb:protein"
      },
      "post_exercise_2_4hrs": {
        "meal": "Balanced meal with carbs, protein, and fats",
        "fluids": "1.5L per kg body weight lost"
      }
    }
  },

  "sport_specific_recommendations": [
    {
      "id": "football_soccer",
      "sport": "Football (Soccer)",
      "category": "Team Sport",
      "intensity": "High",
      "training_frequency": "5-6 days/week",
      "energy_demands": "90 min intermittent high-intensity running, sprints, jumps. Average 1,200-1,500 kcal/game",
      "macronutrient_needs": {
        "carbohydrates": "7-10g/kg/day",
        "protein": "1.6-2.2g/kg/day",
        "fat": "20-25% of total calories"
      },
      "hydration_strategy": {
        "pre_game": ["500ml 2hrs before", "250ml with electrolytes 15min before"],
        "during": ["150-250ml every 15-20min (6-8% carb + electrolytes)"],
        "post": ["1.5L per kg weight lost"]
      },
      "supplements_enhancers": [
        {
          "supplement": "Caffeine",
          "dose": "3-6mg/kg 60min pre-game",
          "benefits": "Endurance, focus"
        },
        {
          "supplement": "Sodium Bicarbonate",
          "dose": "0.3g/kg 60-90min pre-game",
          "benefits": "Buffering capacity"
        },
        {
          "supplement": "Beta-Alanine",
          "dose": "3-6g/day (split doses)",
          "benefits": "Delays fatigue in high-intensity efforts"
        }
      ],
      "sample_diet_plan": [
        {
          "meal": "Breakfast",
          "components": "Oatmeal (80g dry) + Whey protein (30g) + Banana (1)",
          "additional": "Almond butter (1 tbsp) + Cinnamon"
        },
        {
          "meal": "Pre-Game (2 hrs before)",
          "components": "Chicken breast (150g) + White rice (1.5 cups) + Steamed vegetables",
          "additional": "Banana (1) + Sports drink (500ml)"
        },
        {
          "meal": "During Game",
          "components": "Sports drink (500-1000ml) + Energy gel (30-60g carbs)",
          "additional": "Banana (1/2) at halftime"
        },
        {
          "meal": "Post-Game (within 30min)",
          "components": "Protein shake (30g) + Dextrose (50g)",
          "additional": "Banana (1) + Electrolyte drink"
        },
        {
          "meal": "Dinner",
          "components": "Salmon (200g) + Sweet potato (1.5 cups) + Broccoli (2 cups)",
          "additional": "Olive oil (1 tbsp) + Mixed salad"
        }
      ]
    },
    
    {
      "id": "basketball",
      "sport": "Basketball",
      "category": "Team Sport",
      "intensity": "High",
      "training_frequency": "5-6 days/week",
      "energy_demands": "Frequent jumps, sprints, direction changes. High-intensity intervals with short recovery",
      "macronutrient_needs": {
        "carbohydrates": "6-8g/kg/day",
        "protein": "1.8-2.4g/kg/day",
        "fat": "20-25% of total calories"
      },
      "hydration_strategy": {
        "pre_game": ["500ml 2hrs before", "250ml with electrolytes 15min before"],
        "during": ["200-300ml every timeout (6-8% carb + electrolytes)"],
        "post": ["1.5L per kg weight lost"]
      }
    },
    
    {
      "id": "running_distance",
      "sport": "Running (Distance)",
      "category": "Endurance",
      "intensity": "Moderate to High",
      "training_frequency": "5-7 days/week",
      "energy_demands": "Prolonged aerobic activity. Glycogen depletion, muscle damage",
      "macronutrient_needs": {
        "carbohydrates": "8-12g/kg/day",
        "protein": "1.6-1.8g/kg/day",
        "fat": "25-30% of total calories"
      },
      "hydration_strategy": {
        "pre_run": ["500ml 2hrs before", "250ml 15min before"],
        "during": ["150-350ml every 20min (6-8% carb + electrolytes)"],
        "post": ["1.5L per kg lost + sodium (prevents hyponatremia)"]
      }
    },
    
    {
      "id": "swimming",
      "sport": "Swimming",
      "category": "Endurance/Strength",
      "intensity": "High",
      "training_frequency": "6-10 sessions/week",
      "energy_demands": "Full-body resistance training in water. High energy expenditure (700-900 kcal/hr)",
      "macronutrient_needs": {
        "carbohydrates": "8-10g/kg/day",
        "protein": "1.6-2.0g/kg/day",
        "fat": "25-30% of total calories"
      },
      "hydration_strategy": {
        "pre_swim": ["500ml 2hrs before", "250ml 15min before"],
        "during": ["200-400ml every 15-20min (6-8% carb + electrolytes)"],
        "post": ["1.5L per kg lost + sodium"]
      },
      "supplements_enhancers": [
        {
          "supplement": "Caffeine",
          "dose": "3-6mg/kg 60min pre-session",
          "benefits": "Endurance, focus"
        },
        {
          "supplement": "Beta-Alanine",
          "dose": "3-6g/day (split doses)",
          "benefits": "Buffers muscle acidity"
        },
        {
          "supplement": "Omega-3",
          "dose": "1-2g EPA+DHA daily",
          "benefits": "Reduces inflammation"
        }
      ],
      "sample_diet_plan": [
        {
          "meal": "Breakfast",
          "components": "Oatmeal (120g dry) + Whey protein (30g) + Flaxseeds (2 tbsp)",
          "additional": "Mixed berries (1.5 cups) + Almond milk (300ml)"
        },
        {
          "meal": "Pre-Training (2 hrs before)",
          "components": "Banana (1) + Rice cakes (2) + Peanut butter (1 tbsp)",
          "additional": "Honey (1 tsp) + Cinnamon"
        },
        {
          "meal": "During Training (>1hr)",
          "components": "Sports drink (500-750ml) + Energy gel (30g carbs)",
          "additional": "Electrolyte tablets as needed"
        },
        {
          "meal": "Post-Training (within 30min)",
          "components": "Whey protein (30g) + Dextrose (40g)",
          "additional": "Banana (1) + Electrolyte drink"
        },
        {
          "meal": "Dinner",
          "components": "Salmon (200g) + Quinoa (1.5 cups) + Steamed vegetables",
          "additional": "Olive oil (1 tbsp) + Lemon juice"
        }
      ]
    },
    {
      "id": "cycling",
      "sport": "Cycling",
      "category": "Endurance",
      "intensity": "Moderate to High",
      "training_frequency": "5-7 days/week",
      "energy_demands": "Long duration, variable intensity. High energy expenditure (500-1000 kcal/hr)",
      "macronutrient_needs": {
        "carbohydrates": "8-12g/kg/day",
        "protein": "1.6-2.0g/kg/day",
        "fat": "25-30% of total calories"
      },
      "hydration_strategy": {
        "pre_ride": ["500ml 2hrs before", "250ml 15min before"],
        "during": ["500-1000ml/hr + 30-60g carbs/hr + electrolytes"],
        "post": ["1.5L per kg lost + sodium"]
      },
      "supplements_enhancers": [
        {
          "supplement": "Caffeine",
          "dose": "3-6mg/kg 60min pre-ride",
          "benefits": "Endurance, focus"
        },
        {
          "supplement": "Beta-Alanine",
          "dose": "3-6g/day (split doses)",
          "benefits": "Delays fatigue"
        },
        {
          "supplement": "Nitrate (Beetroot)",
          "dose": "300-600mg 2-3hrs before",
          "benefits": "Improves endurance"
        }
      ]
    },
    {
      "id": "weightlifting",
      "sport": "Weightlifting/Strength Training",
      "energy_demands": "Short, intense bursts. Focus on muscle protein synthesis",
      "macronutrient_needs": {
        "carbohydrates": "4-7g/kg/day",
        "protein": "1.6-2.2g/kg/day",
        "fat": "25-30% of total calories"
      },
      "hydration_strategy": {
        "pre_workout": ["500ml 2hrs before", "250ml 15min before"],
        "during": ["200-300ml every 15-20min"],
        "post": ["1L per kg lost + electrolytes"]
      }
    }
  ],
  
  "performance_metrics": {
    "recovery_indicators": [
      "Muscle Soreness (1-10 scale)",
      "Sleep Quality (1-10 scale)",
      "Resting Heart Rate (bpm)",
      "HRV (Heart Rate Variability)",
      "Perceived Recovery Status (1-10 scale)"
    ],
    "common_injuries": {
      "prevention": [
        "Adequate protein intake (1.6-2.2g/kg/day)",
        "Sufficient omega-3 intake (1-2g EPA+DHA)",
        "Vitamin D3 (1000-4000 IU/day)",
        "Collagen peptides (15g/day) + Vitamin C"
      ]
    },
    "performance_enhancement": [
      "Caffeine: 3-6mg/kg 60min before training",
      "Creatine: 5g/day (loading: 20g/day for 5-7 days)",
      "Beta-Alanine: 3-6g/day (split doses)",
      "Sodium Bicarbonate: 0.3g/kg 60-90min before"
    ]
  },

  "sport_specific_tips": {
    "team_sports": [
      "Focus on carb loading 24-48 hours before competition",
      "Prioritize quick-digesting carbs during halftime/intermissions",
      "Emphasize protein for muscle recovery between games"
    ],
    "endurance_sports": [
      "Practice your race-day nutrition strategy during training",
      "Aim for 30-90g carbs per hour during long sessions",
      "Monitor sweat rate to personalize hydration"
    ],
    "strength_sports": [
      "Time protein intake around workouts (before and after)",
      "Ensure adequate carbs to fuel high-intensity efforts",
      "Consider creatine supplementation (5g/day)"
    ]
  },

  "general_tips": {
    "nutrition_optimization": [
      "Periodize Nutrition: Match intake to training phases (build vs. competition)",
      "Food First: Get nutrients from whole foods before supplements",
      "Hydration is Key: Monitor urine color (pale yellow = hydrated)",
      "Recovery Window: Consume protein + carbs within 30-60min post-exercise",
      "Individualize: Needs vary by sport, position, and athlete characteristics"
    ],
    "supplement_guidance": [
      "Third-Party Tested: Choose NSF/Informed Sport certified products",
      "Avoid Banned Substances: Check with WADA/USADA lists",
      "More Isn't Better: Follow evidence-based doses",
      "Cycle When Appropriate: e.g., 8 weeks on creatine, 4 weeks off",
      "Consult Professionals: Sports dietitian/nutritionist guidance"
    ]
  }
};

// Export the data for use in other modules
if (typeof module !== 'undefined' && module.exports) {
  module.exports = sportsNutritionData;
}
