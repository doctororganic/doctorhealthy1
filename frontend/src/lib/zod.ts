import { z } from 'zod';

// Common validation schemas
export const emailSchema = z
  .string()
  .min(1, 'Email is required')
  .email('Invalid email address');

export const passwordSchema = z
  .string()
  .min(8, 'Password must be at least 8 characters')
  .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
  .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
  .regex(/\d/, 'Password must contain at least one number')
  .regex(/[!@#$%^&*(),.?":{}|<>]/, 'Password must contain at least one special character');

export const nameSchema = z
  .string()
  .min(2, 'Name must be at least 2 characters')
  .max(50, 'Name must be less than 50 characters')
  .regex(/^[a-zA-Z\s'-]+$/, 'Name can only contain letters, spaces, hyphens, and apostrophes');

export const phoneSchema = z
  .string()
  .regex(/^[\d\s\-\+\(\)]+$/, 'Invalid phone number format')
  .min(10, 'Phone number must be at least 10 digits')
  .max(20, 'Phone number is too long');

export const dateOfBirthSchema = z
  .string()
  .refine((date) => {
    const parsedDate = new Date(date);
    const now = new Date();
    const minAge = 13;
    const maxAge = 120;

    const age = now.getFullYear() - parsedDate.getFullYear();
    const ageDiff = now.getMonth() - parsedDate.getMonth();
    const adjustedAge = ageDiff < 0 || (ageDiff === 0 && now.getDate() < parsedDate.getDate()) ? age - 1 : age;

    return !isNaN(parsedDate.getTime()) && adjustedAge >= minAge && adjustedAge <= maxAge;
  }, 'You must be between 13 and 120 years old');

export const genderSchema = z.enum(['male', 'female', 'other']);

export const heightSchema = z
  .union([
    z.number().min(50).max(250), // cm
    z.string().regex(/^\d['"]\d{1,2}$/, 'Invalid height format'), // ft'in"
  ])
  .refine((value) => {
    if (typeof value === 'number') return true;
    // Convert ft'in" to cm for validation
    const match = value.match(/^(\d)'(\d{1,2})"$/);
    if (!match) return false;
    const feet = parseInt(match[1]);
    const inches = parseInt(match[2]);
    const totalCm = (feet * 12 + inches) * 2.54;
    return totalCm >= 50 && totalCm <= 250;
  }, 'Height must be between 50cm and 250cm');

export const weightSchema = z
  .union([
    z.number().min(20).max(300), // kg
    z.number().min(44).max(661), // lbs (converted)
  ])
  .refine((value) => {
    if (typeof value === 'number') return true;
    return false; // Should be handled by union
  }, 'Weight must be between 20kg and 300kg');

export const activityLevelSchema = z.enum([
  'sedentary',
  'lightly_active',
  'moderately_active',
  'very_active',
  'extremely_active',
]);

export const unitsSchema = z.enum(['metric', 'imperial']);

export const goalTypeSchema = z.enum([
  'lose_weight',
  'gain_weight',
  'maintain',
  'build_muscle',
  'improve_health',
]);

// Auth schemas
export const loginSchema = z.object({
  email: emailSchema,
  password: z.string().min(1, 'Password is required'),
  rememberMe: z.boolean().optional().default(false),
});

export const registerSchema = z.object({
  email: emailSchema,
  password: passwordSchema,
  confirmPassword: z.string(),
  firstName: nameSchema,
  lastName: nameSchema,
  dateOfBirth: dateOfBirthSchema.optional(),
  gender: genderSchema.optional(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

export const forgotPasswordSchema = z.object({
  email: emailSchema,
});

export const resetPasswordSchema = z.object({
  token: z.string().min(1, 'Reset token is required'),
  newPassword: passwordSchema,
  confirmPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

export const changePasswordSchema = z.object({
  currentPassword: z.string().min(1, 'Current password is required'),
  newPassword: passwordSchema,
  confirmPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

// User profile schemas
export const userProfileSchema = z.object({
  firstName: nameSchema,
  lastName: nameSchema,
  dateOfBirth: dateOfBirthSchema.optional(),
  gender: genderSchema.optional(),
  height: heightSchema.optional(),
  weight: weightSchema.optional(),
  activityLevel: activityLevelSchema.optional(),
  avatar: z.string().url().optional(),
  bio: z.string().max(500, 'Bio must be less than 500 characters').optional(),
  location: z.string().max(100, 'Location must be less than 100 characters').optional(),
  phone: phoneSchema.optional(),
  emergencyContact: z.object({
    name: z.string().min(1, 'Emergency contact name is required'),
    phone: phoneSchema,
    relationship: z.string().min(1, 'Relationship is required'),
  }).optional(),
});

export const userPreferencesSchema = z.object({
  units: unitsSchema,
  language: z.string().min(2).max(5).default('en'),
  timezone: z.string().default('UTC'),
  currency: z.string().min(3).max(3).default('USD'),
  notifications: z.object({
    email: z.boolean().default(true),
    push: z.boolean().default(true),
    sms: z.boolean().default(false),
    mealReminders: z.boolean().default(true),
    workoutReminders: z.boolean().default(true),
    progressUpdates: z.boolean().default(true),
  }),
  privacy: z.object({
    profileVisibility: z.enum(['public', 'private', 'friends']).default('private'),
    shareProgress: z.boolean().default(false),
    allowWorkoutSharing: z.boolean().default(false),
    allowNutritionSharing: z.boolean().default(false),
  }),
});

// Nutrition schemas
export const foodSchema = z.object({
  name: z.string().min(1, 'Food name is required').max(100),
  brand: z.string().max(100).optional(),
  barcode: z.string().max(50).optional(),
  category: z.string().min(1, 'Category is required').max(50),
  description: z.string().max(500).optional(),
  ingredients: z.array(z.string()).optional(),
  servingSize: z.number().min(0.1, 'Serving size must be greater than 0'),
  servingUnit: z.string().min(1, 'Serving unit is required'),
  calories: z.number().min(0, 'Calories must be 0 or greater'),
  macros: z.object({
    protein: z.number().min(0, 'Protein must be 0 or greater'),
    carbs: z.number().min(0, 'Carbs must be 0 or greater'),
    fat: z.number().min(0, 'Fat must be 0 or greater'),
    fiber: z.number().min(0, 'Fiber must be 0 or greater').optional(),
    sugar: z.number().min(0, 'Sugar must be 0 or greater').optional(),
    sodium: z.number().min(0, 'Sodium must be 0 or greater').optional(),
  }),
  micros: z.object({
    vitamins: z.record(z.number()).optional(),
    minerals: z.record(z.number()).optional(),
  }).optional(),
});

export const mealSchema = z.object({
  name: z.string().min(1, 'Meal name is required').max(100),
  mealType: z.enum(['breakfast', 'lunch', 'dinner', 'snack']),
  foods: z.array(z.object({
    foodId: z.string().min(1, 'Food ID is required'),
    quantity: z.number().min(0.1, 'Quantity must be greater than 0'),
    unit: z.string().min(1, 'Unit is required'),
  })).min(1, 'At least one food is required'),
  notes: z.string().max(500).optional(),
  mealDate: z.string().datetime(),
});

export const nutritionGoalSchema = z.object({
  goalType: goalTypeSchema,
  targetWeight: weightSchema.optional(),
  targetDate: z.string().datetime().optional(),
  dailyCalories: z.number().min(800, 'Daily calories must be at least 800').max(10000, 'Daily calories must be less than 10000'),
  macros: z.object({
    protein: z.number().min(0, 'Protein must be 0 or greater'),
    carbs: z.number().min(0, 'Carbs must be 0 or greater'),
    fat: z.number().min(0, 'Fat must be 0 or greater'),
    fiber: z.number().min(0, 'Fiber must be 0 or greater').optional(),
  }),
  micros: z.object({
    vitamins: z.record(z.number()).optional(),
    minerals: z.record(z.number()).optional(),
  }).optional(),
  activityLevel: activityLevelSchema,
  dietaryRestrictions: z.array(z.string()).optional(),
  allergies: z.array(z.string()).optional(),
});

export const waterIntakeSchema = z.object({
  amount: z.number().min(1, 'Amount must be greater than 0'),
  unit: z.enum(['ml', 'oz']),
  timestamp: z.string().datetime(),
});

// Fitness schemas
export const exerciseSchema = z.object({
  name: z.string().min(1, 'Exercise name is required').max(100),
  category: z.enum(['cardio', 'strength', 'flexibility', 'balance', 'sports']),
  muscleGroups: z.array(z.string()).min(1, 'At least one muscle group is required'),
  equipment: z.array(z.string()).optional(),
  difficulty: z.enum(['beginner', 'intermediate', 'advanced']),
  instructions: z.array(z.string()).min(1, 'At least one instruction is required'),
  tips: z.array(z.string()).optional(),
  variations: z.array(z.string()).optional(),
  MET: z.number().min(0, 'MET must be 0 or greater').optional(),
});

export const workoutSchema = z.object({
  name: z.string().min(1, 'Workout name is required').max(100),
  description: z.string().max(500).optional(),
  exercises: z.array(z.object({
    exerciseId: z.string().min(1, 'Exercise ID is required'),
    sets: z.array(z.object({
      reps: z.number().min(1, 'Reps must be greater than 0').optional(),
      weight: z.number().min(0, 'Weight must be 0 or greater').optional(),
      duration: z.number().min(1, 'Duration must be greater than 0').optional(),
      distance: z.number().min(0, 'Distance must be 0 or greater').optional(),
      intensity: z.enum(['low', 'medium', 'high']).optional(),
    })).min(1, 'At least one set is required'),
    restTime: z.number().min(0, 'Rest time must be 0 or greater').optional(),
    notes: z.string().max(200).optional(),
  })).min(1, 'At least one exercise is required'),
  difficulty: z.enum(['beginner', 'intermediate', 'advanced']),
  notes: z.string().max(500).optional(),
  workoutDate: z.string().datetime(),
});

export const activitySchema = z.object({
  type: z.enum(['cardio', 'strength', 'sports', 'other']),
  name: z.string().min(1, 'Activity name is required').max(100),
  duration: z.number().min(1, 'Duration must be greater than 0'),
  caloriesBurned: z.number().min(0, 'Calories burned must be 0 or greater').optional(),
  distance: z.number().min(0, 'Distance must be 0 or greater').optional(),
  steps: z.number().min(0, 'Steps must be 0 or greater').optional(),
  heartRate: z.object({
    average: z.number().min(40, 'Average heart rate must be at least 40').max(220, 'Average heart rate must be less than 220').optional(),
    maximum: z.number().min(40, 'Maximum heart rate must be at least 40').max(220, 'Maximum heart rate must be less than 220').optional(),
    zones: z.array(z.object({
      zone: z.number().min(1).max(5),
      min: z.number().min(0),
      max: z.number().min(0),
      duration: z.number().min(0),
    })).optional(),
  }).optional(),
  intensity: z.enum(['low', 'medium', 'high']).optional(),
  notes: z.string().max(500).optional(),
  activityDate: z.string().datetime(),
});

// Progress schemas
export const weightEntrySchema = z.object({
  weight: weightSchema,
  bodyFat: z.number().min(0).max(100).optional(),
  muscleMass: z.number().min(0).optional(),
  notes: z.string().max(200).optional(),
  date: z.string().datetime(),
});

export const measurementSchema = z.object({
  type: z.enum(['chest', 'waist', 'hips', 'arms', 'thighs', 'custom']),
  value: z.number().min(1, 'Measurement must be greater than 0'),
  unit: z.enum(['cm', 'inches']),
  notes: z.string().max(200).optional(),
  date: z.string().datetime(),
});

export const progressPhotoSchema = z.object({
  url: z.string().url('Invalid photo URL'),
  category: z.enum(['front', 'side', 'back', 'custom']),
  notes: z.string().max(200).optional(),
  date: z.string().datetime(),
});

// Search and filter schemas
export const foodSearchSchema = z.object({
  query: z.string().optional(),
  category: z.string().optional(),
  brand: z.string().optional(),
  minCalories: z.number().min(0).optional(),
  maxCalories: z.number().min(0).optional(),
  verified: z.boolean().optional(),
  limit: z.number().min(1).max(100).optional(),
  offset: z.number().min(0).optional(),
});

export const workoutFiltersSchema = z.object({
  category: z.string().optional(),
  difficulty: z.enum(['beginner', 'intermediate', 'advanced']).optional(),
  muscleGroups: z.array(z.string()).optional(),
  equipment: z.array(z.string()).optional(),
  duration: z.object({
    min: z.number().min(0).optional(),
    max: z.number().min(0).optional(),
  }).optional(),
  dateRange: z.object({
    start: z.string().datetime(),
    end: z.string().datetime(),
  }).optional(),
  limit: z.number().min(1).max(100).optional(),
  offset: z.number().min(0).optional(),
});

export const progressFiltersSchema = z.object({
  dateRange: z.object({
    start: z.string().datetime(),
    end: z.string().datetime(),
  }).optional(),
  type: z.enum(['weight', 'measurements', 'photos', 'all']).optional(),
  limit: z.number().min(1).max(100).optional(),
  offset: z.number().min(0).optional(),
});

// Pagination schemas
export const paginationSchema = z.object({
  page: z.number().min(1).default(1),
  limit: z.number().min(1).max(100).default(20),
});

export const sortSchema = z.object({
  field: z.string(),
  direction: z.enum(['asc', 'desc']).default('asc'),
});

// File upload schemas
export const imageUploadSchema = z.object({
  file: z.instanceof(File).refine(
    (file) => file.size <= 5 * 1024 * 1024, // 5MB
    'File size must be less than 5MB'
  ).refine(
    (file) => ['image/jpeg', 'image/png', 'image/gif', 'image/webp'].includes(file.type),
    'File must be an image (JPEG, PNG, GIF, or WebP)'
  ),
  category: z.string().optional(),
});

export const documentUploadSchema = z.object({
  file: z.instanceof(File).refine(
    (file) => file.size <= 10 * 1024 * 1024, // 10MB
    'File size must be less than 10MB'
  ).refine(
    (file) => ['application/pdf', 'text/plain', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'].includes(file.type),
    'File must be a PDF, TXT, or Word document'
  ),
});

// Type inference
export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;
export type ChangePasswordInput = z.infer<typeof changePasswordSchema>;
export type UserProfileInput = z.infer<typeof userProfileSchema>;
export type UserPreferencesInput = z.infer<typeof userPreferencesSchema>;
export type FoodInput = z.infer<typeof foodSchema>;
export type MealInput = z.infer<typeof mealSchema>;
export type NutritionGoalInput = z.infer<typeof nutritionGoalSchema>;
export type WaterIntakeInput = z.infer<typeof waterIntakeSchema>;
export type ExerciseInput = z.infer<typeof exerciseSchema>;
export type WorkoutInput = z.infer<typeof workoutSchema>;
export type ActivityInput = z.infer<typeof activitySchema>;
export type WeightEntryInput = z.infer<typeof weightEntrySchema>;
export type MeasurementInput = z.infer<typeof measurementSchema>;
export type ProgressPhotoInput = z.infer<typeof progressPhotoSchema>;
export type FoodSearchParams = z.infer<typeof foodSearchSchema>;
export type WorkoutFilters = z.infer<typeof workoutFiltersSchema>;
export type ProgressFilters = z.infer<typeof progressFiltersSchema>;
export type PaginationInput = z.infer<typeof paginationSchema>;
export type SortInput = z.infer<typeof sortSchema>;
export type ImageUploadInput = z.infer<typeof imageUploadSchema>;
export type DocumentUploadInput = z.infer<typeof documentUploadSchema>;

// Error handling utilities
export const getZodErrors = (error: z.ZodError): Record<string, string> => {
  const errors: Record<string, string> = {};
  
  error.errors.forEach((err) => {
    const path = err.path.join('.');
    errors[path] = err.message;
  });
  
  return errors;
};

export const getFirstZodError = (error: z.ZodError): string => {
  return error.errors[0]?.message || 'Validation failed';
};

// Form validation utilities
export const validateForm = <T>(schema: z.ZodSchema<T>, data: unknown): {
  success: boolean;
  data?: T;
  errors?: Record<string, string>;
  error?: string;
} => {
  const result = schema.safeParse(data);
  
  if (result.success) {
    return {
      success: true,
      data: result.data,
    };
  } else {
    return {
      success: false,
      errors: getZodErrors(result.error),
      error: getFirstZodError(result.error),
    };
  }
};

// Async validation utilities
export const validateEmailAsync = async (email: string): Promise<boolean> => {
  try {
    const response = await fetch(`/api/validate/email?email=${encodeURIComponent(email)}`);
    const data = await response.json();
    return data.available;
  } catch (error) {
    return false;
  }
};

export const validateUsernameAsync = async (username: string): Promise<boolean> => {
  try {
    const response = await fetch(`/api/validate/username?username=${encodeURIComponent(username)}`);
    const data = await response.json();
    return data.available;
  } catch (error) {
    return false;
  }
};
