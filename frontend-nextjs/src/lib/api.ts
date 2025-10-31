import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

// API functions
export const nutritionAPI = {
  analyze: async (data: {
    food: string;
    quantity: number;
    unit: string;
    checkHalal?: boolean;
  }) => {
    const response = await api.post('/nutrition/analyze', data);
    return response.data;
  },
};

export const recipesAPI = {
  getAll: async (params?: { meal_type?: string; medical_conditions?: string }) => {
    const response = await api.get('/recipes', { params });
    return response.data;
  },
  getById: async (id: string) => {
    const response = await api.get(`/recipes/${id}`);
    return response.data;
  },
};

export const workoutsAPI = {
  getAll: async (params?: {
    goal?: string;
    experience?: string;
    injury_location?: string;
  }) => {
    const response = await api.get('/workouts', { params });
    return response.data;
  },
};

export const mealsAPI = {
  getAll: async () => {
    const response = await api.get('/meals');
    return response.data;
  },
};

export const healthAPI = {
  check: async () => {
    const response = await axios.get(`${API_URL}/health`);
    return response.data;
  },
};
