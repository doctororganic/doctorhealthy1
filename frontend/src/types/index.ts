// User Types
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'user' | 'admin';
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UserProfile {
  id: string;
  userId: string;
  dateOfBirth?: string;
  gender?: 'male' | 'female' | 'other';
  height?: number; // in cm
  weight?: number; // in kg
  activityLevel?: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extremely_active';
  avatar?: string;
  bio?: string;
  location?: string;
  phone?: string;
  emergencyContact?: {
    name: string;
    phone: string;
    relationship: string;
  };
}

export interface UserPreferences {
  id: string;
  userId: string;
  units: 'metric' | 'imperial';
  language: string;
  timezone: string;
  currency: string;
  notifications: {
    email: boolean;
    push: boolean;
    sms: boolean;
    mealReminders: boolean;
    workoutReminders: boolean;
    progressUpdates: boolean;
  };
  privacy: {
    profileVisibility: 'public' | 'private' | 'friends';
    shareProgress: boolean;
    allowWorkoutSharing: boolean;
    allowNutritionSharing: boolean;
  };
}

export interface NutritionGoal {
  id: string;
  userId: string;
  goalType: 'lose_weight' | 'gain_weight' | 'maintain' | 'build_muscle' | 'improve_health';
  targetWeight?: number; // in kg
  targetDate?: string;
  dailyCalories: number;
  macros: {
    protein: number; // in grams
    carbs: number; // in grams
    fat: number; // in grams
    fiber?: number; // in grams
  };
  micros?: {
    vitamins: Record<string, number>;
    minerals: Record<string, number>;
  };
  activityLevel: 'sedentary' | 'lightly_active' | 'moderately_active' | 'very_active' | 'extremely_active';
  dietaryRestrictions: string[];
  allergies: string[];
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

// Nutrition Types
export interface Food {
  id: string;
  name: string;
  brand?: string;
  barcode?: string;
  category: string;
  description?: string;
  ingredients?: string[];
  servingSize: number;
  servingUnit: string;
  calories: number;
  macros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber?: number;
    sugar?: number;
    sodium?: number;
  };
  micros?: {
    vitamins: Record<string, number>;
    minerals: Record<string, number>;
  };
  verified: boolean;
  source: 'user' | 'database' | 'api';
  createdAt: string;
  updatedAt: string;
}

export interface Meal {
  id: string;
  userId: string;
  name: string;
  mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  foods: MealFood[];
  totalCalories: number;
  totalMacros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber: number;
    sugar: number;
    sodium: number;
  };
  notes?: string;
  mealDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface MealFood {
  id: string;
  foodId: string;
  quantity: number;
  unit: string;
  calories: number;
  macros: {
    protein: number;
    carbs: number;
    fat: number;
    fiber?: number;
    sugar?: number;
    sodium?: number;
  };
}

export interface WaterIntake {
  id: string;
  userId: string;
  amount: number; // in ml
  unit: 'ml' | 'oz';
  timestamp: string;
  createdAt: string;
}

// Fitness Types
export interface Exercise {
  id: string;
  name: string;
  category: 'cardio' | 'strength' | 'flexibility' | 'balance' | 'sports';
  muscleGroups: string[];
  equipment: string[];
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  instructions: string[];
  tips?: string[];
  variations?: string[];
  MET?: number; // Metabolic equivalent of task
  verified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Workout {
  id: string;
  userId: string;
  name: string;
  description?: string;
  exercises: WorkoutExercise[];
  duration: number; // in minutes
  caloriesBurned?: number;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  notes?: string;
  workoutDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface WorkoutExercise {
  id: string;
  exerciseId: string;
  sets: ExerciseSet[];
  restTime?: number; // in seconds
  notes?: string;
}

export interface ExerciseSet {
  id: string;
  reps?: number;
  weight?: number; // in kg
  duration?: number; // in seconds
  distance?: number; // in meters
  intensity?: 'low' | 'medium' | 'high';
}

export interface Activity {
  id: string;
  userId: string;
  type: 'cardio' | 'strength' | 'sports' | 'other';
  name: string;
  duration: number; // in minutes
  caloriesBurned?: number;
  distance?: number; // in meters
  steps?: number;
  heartRate?: {
    average?: number;
    maximum?: number;
    zones?: HeartRateZone[];
  };
  intensity?: 'low' | 'medium' | 'high';
  notes?: string;
  activityDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface HeartRateZone {
  zone: 1 | 2 | 3 | 4 | 5;
  min: number;
  max: number;
  duration: number; // in seconds
}

// Progress Types
export interface WeightEntry {
  id: string;
  userId: string;
  weight: number; // in kg
  bodyFat?: number; // percentage
  muscleMass?: number; // in kg
  notes?: string;
  date: string;
  createdAt: string;
}

export interface Measurement {
  id: string;
  userId: string;
  type: 'chest' | 'waist' | 'hips' | 'arms' | 'thighs' | 'custom';
  value: number; // in cm
  unit: 'cm' | 'inches';
  notes?: string;
  date: string;
  createdAt: string;
}

export interface ProgressPhoto {
  id: string;
  userId: string;
  url: string;
  category: 'front' | 'side' | 'back' | 'custom';
  notes?: string;
  date: string;
  createdAt: string;
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
  meta?: {
    pagination?: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
    };
    timestamp: string;
  };
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  meta: {
    pagination: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
    };
    timestamp: string;
  };
}

// Auth Types
export interface LoginRequest {
  email: string;
  password: string;
  rememberMe?: boolean;
}

export interface RegisterRequest {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  dateOfBirth?: string;
  gender?: 'male' | 'female' | 'other';
}

export interface AuthResponse {
  user: User;
  profile?: UserProfile;
  preferences?: UserPreferences;
  tokens: {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
  };
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface PasswordResetRequest {
  email: string;
}

export interface PasswordResetConfirmRequest {
  token: string;
  newPassword: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}

// Form Types
export interface FoodSearchFilters {
  query?: string;
  category?: string;
  brand?: string;
  minCalories?: number;
  maxCalories?: number;
  verified?: boolean;
  limit?: number;
  offset?: number;
}

export interface WorkoutFilters {
  category?: string;
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  muscleGroups?: string[];
  equipment?: string[];
  duration?: {
    min?: number;
    max?: number;
  };
  dateRange?: {
    start: string;
    end: string;
  };
  limit?: number;
  offset?: number;
}

export interface ProgressFilters {
  dateRange?: {
    start: string;
    end: string;
  };
  type?: 'weight' | 'measurements' | 'photos' | 'all';
  limit?: number;
  offset?: number;
}

// Dashboard Types
export interface DashboardStats {
  totalCaloriesConsumed: number;
  totalCaloriesBurned: number;
  netCalories: number;
  waterIntake: number;
  workoutsCompleted: number;
  activeStreak: number;
  currentWeight?: number;
  goalProgress: {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
  };
  recentActivities: {
    meals: Meal[];
    workouts: Workout[];
    activities: Activity[];
  };
}

// Chart Types
export interface ChartDataPoint {
  x: string | number;
  y: number;
  label?: string;
}

export interface NutritionChart {
  calories: ChartDataPoint[];
  macros: {
    protein: ChartDataPoint[];
    carbs: ChartDataPoint[];
    fat: ChartDataPoint[];
  };
  weight: ChartDataPoint[];
}

export interface FitnessChart {
  workoutsPerWeek: ChartDataPoint[];
  caloriesBurned: ChartDataPoint[];
  strengthProgress: ChartDataPoint[];
  cardioProgress: ChartDataPoint[];
}
