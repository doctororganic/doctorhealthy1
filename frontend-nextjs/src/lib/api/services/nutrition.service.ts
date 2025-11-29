import { z } from 'zod';

const API = process.env.NEXT_PUBLIC_API_BASE || '/api/v1';

// Generic fetch with retry and Zod validation
async function fetchJson<T>(path: string, schema: z.ZodType<T>, init?: RequestInit): Promise<T> {
  const url = `${API}${path}`;
  let attempt = 0;
  let lastErr: any;
  while (attempt < 3) {
    try {
      const res = await fetch(url, {
        ...init,
        headers: {
          'Content-Type': 'application/json',
          ...(init?.headers || {}),
        },
        cache: 'no-store',
      });
      if (!res.ok) {
        const msg = await res.text();
        throw new Error(`HTTP ${res.status} ${res.statusText}: ${msg}`);
      }
      const data = await res.json();
      const parsed = schema.parse(data);
      return parsed;
    } catch (e) {
      lastErr = e;
      attempt += 1;
      await new Promise(r => setTimeout(r, 300 * Math.pow(2, attempt - 1)));
    }
  }
  throw lastErr;
}

// Types and schemas
export const paginationSchema = z.object({
  page: z.number().int().min(1).default(1),
  limit: z.number().int().min(1).max(200).default(20),
  total: z.number().int().nonnegative().optional(),
});

const recipeSchema = z.object({
  id: z.string().or(z.number()).optional(),
  name: z.string(),
  calories: z.number().optional(),
  protein: z.number().optional(),
  carbs: z.number().optional(),
  fat: z.number().optional(),
  ingredients: z.array(z.any()).optional(),
});

const workoutSchema = z.object({
  id: z.string().or(z.number()).optional(),
  goal: z.string().optional(),
  experience_level: z.string().optional(),
  plan: z.any().optional(),
});

const complaintSchema = z.object({
  id: z.string().or(z.number()).optional(),
  name: z.string(),
  description: z.string().optional(),
});

const metabolismSchema = z.object({
  id: z.string().or(z.number()).optional(),
  type: z.string(),
  description: z.string().optional(),
});

const drugNutritionSchema = z.object({
  id: z.string().or(z.number()).optional(),
  drug: z.string(),
  interaction: z.string().optional(),
});

const listSchema = <T>(item: z.ZodType<T>) => z.object({
  items: z.array(item),
  page: z.number().int().default(1),
  limit: z.number().int().default(20),
  total: z.number().int().optional(),
});

export type RecipesParams = Partial<{ page: number; limit: number; q: string }>;
export type WorkoutsParams = Partial<{ page: number; limit: number; goal: string; level: string }>;
export type ComplaintsParams = Partial<{ page: number; limit: number; q: string }>;
export type MetabolismParams = Partial<{ page: number; limit: number; q: string }>;
export type DrugsNutritionParams = Partial<{ page: number; limit: number; q: string }>;

function qs(params?: Record<string, unknown>) {
  if (!params) return '';
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v === undefined || v === null || v === '') return;
    sp.set(k, String(v));
  });
  const str = sp.toString();
  return str ? `?${str}` : '';
}

export const nutritionService = {
  async getRecipes(params?: RecipesParams) {
    return fetchJson(`/nutrition-data/recipes${qs(params)}`, listSchema(recipeSchema));
  },
  async getWorkouts(params?: WorkoutsParams) {
    return fetchJson(`/nutrition-data/workouts${qs(params)}`, listSchema(workoutSchema));
  },
  async getComplaints(params?: ComplaintsParams) {
    return fetchJson(`/nutrition-data/complaints${qs(params)}`, listSchema(complaintSchema));
  },
  async getMetabolism(params?: MetabolismParams) {
    return fetchJson(`/nutrition-data/metabolism${qs(params)}`, listSchema(metabolismSchema));
  },
  async getDrugsNutrition(params?: DrugsNutritionParams) {
    return fetchJson(`/nutrition-data/drugs-nutrition${qs(params)}`, listSchema(drugNutritionSchema));
  },
};
