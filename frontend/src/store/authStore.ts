import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User, LoginRequest, RegisterRequest, AuthResponse } from '@/types';
import { apiClient } from '@/lib/api';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true, error: null });
        
        try {
          const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
          const { user, tokens } = response.data;
          
          // Set Authorization header for future requests
          apiClient.defaults.headers.common['Authorization'] = `Bearer ${tokens.accessToken}`;
          
          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error: any) {
          const errorMessage = error.response?.data?.error?.message || 'Login failed';
          set({
            isLoading: false,
            error: errorMessage,
          });
          throw new Error(errorMessage);
        }
      },

      register: async (userData: RegisterRequest) => {
        set({ isLoading: true, error: null });
        
        try {
          await apiClient.post<AuthResponse>('/auth/register', userData);
          
          set({
            isLoading: false,
            error: null,
          });
        } catch (error: any) {
          const errorMessage = error.response?.data?.error?.message || 'Registration failed';
          set({
            isLoading: false,
            error: errorMessage,
          });
          throw new Error(errorMessage);
        }
      },

      logout: () => {
        // Remove Authorization header
        delete apiClient.defaults.headers.common['Authorization'];
        
        set({
          user: null,
          isAuthenticated: false,
          error: null,
        });
      },

      refreshToken: async () => {
        const { user } = get();
        if (!user) return;

        try {
          const response = await apiClient.post<{ accessToken: string }>('/auth/refresh');
          const { accessToken } = response.data;
          
          // Update Authorization header
          apiClient.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
        } catch (error) {
          // If refresh fails, logout the user
          get().logout();
        }
      },

      clearError: () => {
        set({ error: null });
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Set Authorization header if user is authenticated
        if (state?.isAuthenticated && state?.user) {
          // Note: In a real app, you'd want to store and retrieve the access token
          // For now, we'll just set a placeholder
          apiClient.defaults.headers.common['Authorization'] = 'Bearer stored-token';
        }
      },
    }
  )
);
