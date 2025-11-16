import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { toast } from 'react-hot-toast';
import { ApiResponse, PaginatedResponse } from '@/types';

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const APP_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000';

// Create axios instance
class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    });

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add auth token if available
        const token = this.getAuthToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        // Add request ID for tracing
        config.headers['X-Request-ID'] = this.generateRequestId();

        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response: AxiosResponse<ApiResponse<any>>) => {
        // Log successful responses in development
        if (process.env.NODE_ENV === 'development') {
          console.log(`API Response [${response.config.method?.toUpperCase()}] ${response.config.url}:`, response.data);
        }

        return response;
      },
      async (error) => {
        const originalRequest = error.config;

        // Handle 401 Unauthorized - try to refresh token
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            await this.refreshToken();
            const token = this.getAuthToken();
            if (token) {
              originalRequest.headers.Authorization = `Bearer ${token}`;
              return this.client(originalRequest);
            }
          } catch (refreshError) {
            // Refresh failed, redirect to login
            this.handleAuthError();
            return Promise.reject(refreshError);
          }
        }

        // Handle other errors
        this.handleError(error);
        return Promise.reject(error);
      }
    );
  }

  private getAuthToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('access_token');
    }
    return null;
  }

  private getRefreshToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('refresh_token');
    }
    return null;
  }

  private setTokens(accessToken: string, refreshToken: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', accessToken);
      localStorage.setItem('refresh_token', refreshToken);
    }
  }

  private clearTokens(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  }

  private generateRequestId(): string {
    return `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private async refreshToken(): Promise<void> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
        refreshToken,
      });

      const { accessToken, refreshToken: newRefreshToken } = response.data.data.tokens;
      this.setTokens(accessToken, newRefreshToken);
    } catch (error) {
      throw new Error('Failed to refresh token');
    }
  }

  private handleAuthError(): void {
    this.clearTokens();
    
    // Show toast notification
    toast.error('Your session has expired. Please log in again.');
    
    // Redirect to login page
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
  }

  private handleError(error: any): void {
    const message = error.response?.data?.error?.message || 'An unexpected error occurred';
    
    // Don't show toast for 401 errors (handled above)
    if (error.response?.status !== 401) {
      toast.error(message);
    }

    // Log error details in development
    if (process.env.NODE_ENV === 'development') {
      console.error('API Error:', {
        status: error.response?.status,
        message,
        details: error.response?.data?.error?.details,
        config: error.config,
      });
    }
  }

  // HTTP Methods
  public async get<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.get<ApiResponse<T>>(url, config);
    return response.data;
  }

  public async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.post<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  public async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.put<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  public async patch<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.patch<ApiResponse<T>>(url, data, config);
    return response.data;
  }

  public async delete<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.delete<ApiResponse<T>>(url, config);
    return response.data;
  }

  // Pagination helper
  public async getPaginated<T>(url: string, params?: any): Promise<PaginatedResponse<T>> {
    const response = await this.client.get<PaginatedResponse<T>>(url, { params });
    return response.data;
  }

  // File upload
  public async uploadFile<T>(url: string, file: File, onProgress?: (progress: number) => void): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);

    const config: AxiosRequestConfig = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(progress);
        }
      },
    };

    const response = await this.client.post<ApiResponse<T>>(url, formData, config);
    return response.data;
  }

  // Multiple file upload
  public async uploadFiles<T>(url: string, files: File[], onProgress?: (progress: number) => void): Promise<ApiResponse<T>> {
    const formData = new FormData();
    files.forEach((file, index) => {
      formData.append(`files[${index}]`, file);
    });

    const config: AxiosRequestConfig = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(progress);
        }
      },
    };

    const response = await this.client.post<ApiResponse<T>>(url, formData, config);
    return response.data;
  }
}

// Create singleton instance
const apiClient = new ApiClient();

// Export specific API methods
export const api = {
  // Auth
  auth: {
    login: (data: { email: string; password: string; rememberMe?: boolean }) =>
      apiClient.post('/auth/login', data),
    
    register: (data: { email: string; password: string; firstName: string; lastName: string; dateOfBirth?: string; gender?: string }) =>
      apiClient.post('/auth/register', data),
    
    logout: () =>
      apiClient.post('/auth/logout'),
    
    refresh: (refreshToken: string) =>
      apiClient.post('/auth/refresh', { refreshToken }),
    
    me: () =>
      apiClient.get('/auth/me'),
    
    changePassword: (data: { currentPassword: string; newPassword: string }) =>
      apiClient.post('/auth/change-password', data),
    
    forgotPassword: (email: string) =>
      apiClient.post('/auth/forgot-password', { email }),
    
    resetPassword: (data: { token: string; newPassword: string }) =>
      apiClient.post('/auth/reset-password', data),
  },

  // Users
  users: {
    getProfile: () =>
      apiClient.get('/users/profile'),
    
    updateProfile: (data: any) =>
      apiClient.put('/users/profile', data),
    
    updatePreferences: (data: any) =>
      apiClient.put('/users/preferences', data),
    
    getPreferences: () =>
      apiClient.get('/users/preferences'),
    
    deleteAccount: () =>
      apiClient.delete('/users/account'),
  },

  // Nutrition
  nutrition: {
    getFoods: (params?: any) =>
      apiClient.getPaginated('/nutrition/foods', params),
    
    getFood: (id: string) =>
      apiClient.get(`/nutrition/foods/${id}`),
    
    createFood: (data: any) =>
      apiClient.post('/nutrition/foods', data),
    
    updateFood: (id: string, data: any) =>
      apiClient.put(`/nutrition/foods/${id}`, data),
    
    deleteFood: (id: string) =>
      apiClient.delete(`/nutrition/foods/${id}`),
    
    searchFoods: (query: string, params?: any) =>
      apiClient.getPaginated('/nutrition/foods/search', { query, ...params }),
    
    getMeals: (params?: any) =>
      apiClient.getPaginated('/nutrition/meals', params),
    
    getMeal: (id: string) =>
      apiClient.get(`/nutrition/meals/${id}`),
    
    createMeal: (data: any) =>
      apiClient.post('/nutrition/meals', data),
    
    updateMeal: (id: string, data: any) =>
      apiClient.put(`/nutrition/meals/${id}`, data),
    
    deleteMeal: (id: string) =>
      apiClient.delete(`/nutrition/meals/${id}`),
    
    logWater: (data: { amount: number; unit: string }) =>
      apiClient.post('/nutrition/water', data),
    
    getWaterIntake: (params?: any) =>
      apiClient.getPaginated('/nutrition/water', params),
    
    getNutritionGoals: () =>
      apiClient.get('/nutrition/goals'),
    
    createNutritionGoal: (data: any) =>
      apiClient.post('/nutrition/goals', data),
    
    updateNutritionGoal: (id: string, data: any) =>
      apiClient.put(`/nutrition/goals/${id}`, data),
    
    deleteNutritionGoal: (id: string) =>
      apiClient.delete(`/nutrition/goals/${id}`),
  },

  // Fitness
  fitness: {
    getExercises: (params?: any) =>
      apiClient.getPaginated('/fitness/exercises', params),
    
    getExercise: (id: string) =>
      apiClient.get(`/fitness/exercises/${id}`),
    
    createExercise: (data: any) =>
      apiClient.post('/fitness/exercises', data),
    
    updateExercise: (id: string, data: any) =>
      apiClient.put(`/fitness/exercises/${id}`, data),
    
    deleteExercise: (id: string) =>
      apiClient.delete(`/fitness/exercises/${id}`),
    
    searchExercises: (query: string, params?: any) =>
      apiClient.getPaginated('/fitness/exercises/search', { query, ...params }),
    
    getWorkouts: (params?: any) =>
      apiClient.getPaginated('/fitness/workouts', params),
    
    getWorkout: (id: string) =>
      apiClient.get(`/fitness/workouts/${id}`),
    
    createWorkout: (data: any) =>
      apiClient.post('/fitness/workouts', data),
    
    updateWorkout: (id: string, data: any) =>
      apiClient.put(`/fitness/workouts/${id}`, data),
    
    deleteWorkout: (id: string) =>
      apiClient.delete(`/fitness/workouts/${id}`),
    
    getActivities: (params?: any) =>
      apiClient.getPaginated('/fitness/activities', params),
    
    getActivity: (id: string) =>
      apiClient.get(`/fitness/activities/${id}`),
    
    createActivity: (data: any) =>
      apiClient.post('/fitness/activities', data),
    
    updateActivity: (id: string, data: any) =>
      apiClient.put(`/fitness/activities/${id}`, data),
    
    deleteActivity: (id: string) =>
      apiClient.delete(`/fitness/activities/${id}`),
  },

  // Progress
  progress: {
    getWeightEntries: (params?: any) =>
      apiClient.getPaginated('/progress/weight', params),
    
    createWeightEntry: (data: any) =>
      apiClient.post('/progress/weight', data),
    
    updateWeightEntry: (id: string, data: any) =>
      apiClient.put(`/progress/weight/${id}`, data),
    
    deleteWeightEntry: (id: string) =>
      apiClient.delete(`/progress/weight/${id}`),
    
    getMeasurements: (params?: any) =>
      apiClient.getPaginated('/progress/measurements', params),
    
    createMeasurement: (data: any) =>
      apiClient.post('/progress/measurements', data),
    
    updateMeasurement: (id: string, data: any) =>
      apiClient.put(`/progress/measurements/${id}`, data),
    
    deleteMeasurement: (id: string) =>
      apiClient.delete(`/progress/measurements/${id}`),
    
    getProgressPhotos: (params?: any) =>
      apiClient.getPaginated('/progress/photos', params),
    
    uploadProgressPhoto: (file: File, data: any, onProgress?: (progress: number) => void) =>
      apiClient.uploadFile('/progress/photos', file, onProgress),
    
    deleteProgressPhoto: (id: string) =>
      apiClient.delete(`/progress/photos/${id}`),
    
    getStats: (params?: any) =>
      apiClient.get('/progress/stats', { params }),
    
    getCharts: (params?: any) =>
      apiClient.get('/progress/charts', { params }),
  },

  // Dashboard
  dashboard: {
    getStats: () =>
      apiClient.get('/dashboard/stats'),
    
    getRecentActivity: () =>
      apiClient.get('/dashboard/recent-activity'),
    
    getRecommendations: () =>
      apiClient.get('/dashboard/recommendations'),
  },

  // Upload
  upload: {
    uploadImage: (file: File, onProgress?: (progress: number) => void) =>
      apiClient.uploadFile('/upload/image', file, onProgress),
    
    uploadDocument: (file: File, onProgress?: (progress: number) => void) =>
      apiClient.uploadFile('/upload/document', file, onProgress),
  },
};

export default apiClient;
