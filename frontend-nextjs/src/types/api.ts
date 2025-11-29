// API Response Types
export interface APIResponse<T> {
  status: 'success' | 'error';
  data?: T;
  items?: T[];
  pagination?: PaginationMeta;
  error?: string;
  filters?: Record<string, any>;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

// Recipe Types
export interface Recipe {
  id: string | number;
  name: string;
  calories?: number;
  protein?: number;
  carbs?: number;
  fat?: number;
  ingredients?: string[];
  instructions?: string[];
  cuisine?: string;
  dietType?: string;
  prepTime?: number;
  cookTime?: number;
  servings?: number;
  isHalal?: boolean;
}

// Workout Types
export interface Workout {
  id: string | number;
  goal?: string;
  experience_level?: string;
  plan?: any;
  title?: {
    en?: string;
    ar?: string;
  };
  description?: {
    en?: string;
    ar?: string;
  };
}

// Complaint Types
export interface Complaint {
  id: string | number;
  name: string;
  description?: string;
  condition_en?: string;
  condition_ar?: string;
}

// Metabolism Types
export interface Metabolism {
  id: string | number;
  type: string;
  description?: string;
  section_id?: string;
  title_en?: string;
  title_ar?: string;
}

// Drug Nutrition Types
export interface DrugNutrition {
  id: string | number;
  drug: string;
  interaction?: string;
  effect?: string;
  recommendation?: string;
}

// Disease Types
export interface Disease {
  id: string | number;
  name: string;
  description?: string;
  nutrition_guidelines?: string[];
  foods_to_avoid?: string[];
  recommended_foods?: string[];
}

// Injury Types
export interface Injury {
  id: string | number;
  name: string;
  description?: string;
  recovery_nutrition?: string[];
  recommended_foods?: string[];
  foods_to_avoid?: string[];
}

// Vitamin Types
export interface Vitamin {
  id: string | number;
  name: string;
  description?: string;
  benefits?: string[];
  sources?: string[];
  recommended_daily_amount?: string;
  deficiency_symptoms?: string[];
}

// Mineral Types
export interface Mineral {
  id: string | number;
  name: string;
  description?: string;
  benefits?: string[];
  sources?: string[];
  recommended_daily_amount?: string;
  deficiency_symptoms?: string[];
}

// User Types
export interface User {
  id: string | number;
  email: string;
  name?: string;
  age?: number;
  gender?: 'male' | 'female' | 'other';
  weight?: number;
  height?: number;
  activity_level?: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
  goals?: string[];
  dietary_restrictions?: string[];
  created_at?: string;
  updated_at?: string;
}

// Nutrition Goal Types
export interface NutritionGoal {
  id: string | number;
  user_id: string | number;
  goal_type: 'weight_loss' | 'weight_gain' | 'muscle_gain' | 'maintenance';
  target_weight?: number;
  daily_calories?: number;
  daily_protein?: number;
  daily_carbs?: number;
  daily_fat?: number;
  target_date?: string;
  status: 'active' | 'completed' | 'paused';
  created_at?: string;
  updated_at?: string;
}

// Weight Log Types
export interface WeightLog {
  id: string | number;
  user_id: string | number;
  weight: number;
  date: string;
  notes?: string;
  created_at?: string;
}

// Exercise Types
export interface Exercise {
  id: string | number;
  name: string;
  description?: string;
  muscle_groups?: string[];
  equipment?: string[];
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  instructions?: string[];
  video_url?: string;
}

// Body Measurement Types
export interface BodyMeasurement {
  id: string | number;
  user_id: string | number;
  date: string;
  weight?: number;
  body_fat_percentage?: number;
  muscle_mass?: number;
  waist?: number;
  chest?: number;
  arms?: number;
  thighs?: number;
  notes?: string;
  created_at?: string;
}

// Action Types
export interface FitnessAction {
  id: string | number;
  type: 'workout' | 'meal' | 'measurement' | 'goal';
  title: string;
  description?: string;
  data?: Record<string, any>;
  user_id: string | number;
  status: 'pending' | 'completed' | 'skipped';
  scheduled_at?: string;
  completed_at?: string;
  created_at?: string;
}

// Progress Types
export interface ProgressSummary {
  user_id: string | number;
  period: 'week' | 'month' | 'quarter' | 'year';
  start_date: string;
  end_date: string;
  weight_change?: number;
  measurements_change?: Record<string, number>;
  workouts_completed?: number;
  calories_burned?: number;
  goals_achieved?: number;
}

// API Error Types
export interface APIError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

// Search/Filter Types
export interface SearchParams {
  q?: string;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  filters?: Record<string, any>;
}

// Response Wrapper Types
export interface ListResponse<T> {
  items: T[];
  pagination: PaginationMeta;
}

export interface SingleResponse<T> {
  data: T;
}

export interface ErrorResponse {
  error: APIError;
  status: 'error';
}
