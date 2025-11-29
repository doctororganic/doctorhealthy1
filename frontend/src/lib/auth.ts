import { jwtDecode } from 'jose';
import { User } from '@/types';

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface AuthUser extends User {
  profile?: any;
  preferences?: any;
}

class AuthManager {
  private static instance: AuthManager;
  private user: AuthUser | null = null;
  private tokens: AuthTokens | null = null;
  private refreshTimer: NodeJS.Timeout | null = null;

  private constructor() {
    this.initializeFromStorage();
  }

  public static getInstance(): AuthManager {
    if (!AuthManager.instance) {
      AuthManager.instance = new AuthManager();
    }
    return AuthManager.instance;
  }

  private initializeFromStorage(): void {
    if (typeof window === 'undefined') return;

    try {
      const accessToken = localStorage.getItem('access_token');
      const refreshToken = localStorage.getItem('refresh_token');
      const userStr = localStorage.getItem('user');

      if (accessToken && refreshToken && userStr) {
        this.tokens = {
          accessToken,
          refreshToken,
          expiresIn: this.getTokenExpirationTime(accessToken),
        };
        this.user = JSON.parse(userStr);
        this.scheduleTokenRefresh();
      }
    } catch (error) {
      console.error('Failed to initialize auth from storage:', error);
      this.clearStorage();
    }
  }

  private getTokenExpirationTime(token: string): number {
    try {
      const decoded = jwtDecode(token);
      const exp = decoded.exp;
      if (exp) {
        return exp * 1000 - Date.now(); // Convert to milliseconds and subtract current time
      }
    } catch (error) {
      console.error('Failed to decode token:', error);
    }
    return 0;
  }

  private scheduleTokenRefresh(): void {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
    }

    if (!this.tokens) return;

    // Refresh token 5 minutes before it expires
    const refreshTime = Math.max(this.tokens.expiresIn - 5 * 60 * 1000, 0);

    this.refreshTimer = setTimeout(() => {
      this.refreshTokens();
    }, refreshTime);
  }

  private async refreshTokens(): Promise<void> {
    if (!this.tokens?.refreshToken) {
      this.logout();
      return;
    }

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          refreshToken: this.tokens.refreshToken,
        }),
      });

      if (!response.ok) {
        throw new Error('Token refresh failed');
      }

      const data = await response.json();
      const newTokens = data.data.tokens;

      this.setTokens(newTokens);
    } catch (error) {
      console.error('Failed to refresh tokens:', error);
      this.logout();
    }
  }

  private setTokens(tokens: AuthTokens): void {
    this.tokens = tokens;
    
    if (typeof window !== 'undefined') {
      localStorage.setItem('access_token', tokens.accessToken);
      localStorage.setItem('refresh_token', tokens.refreshToken);
    }

    this.scheduleTokenRefresh();
  }

  private clearStorage(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
    }
  }

  public authenticate(user: AuthUser, tokens: AuthTokens): void {
    this.user = user;
    this.setTokens(tokens);
    
    if (typeof window !== 'undefined') {
      localStorage.setItem('user', JSON.stringify(user));
    }

    // Dispatch custom event for other components to listen to
    window.dispatchEvent(new CustomEvent('auth:login', { detail: { user, tokens } }));
  }

  public logout(): void {
    if (this.refreshTimer) {
      clearTimeout(this.refreshTimer);
      this.refreshTimer = null;
    }

    this.user = null;
    this.tokens = null;
    this.clearStorage();

    // Dispatch custom event for other components to listen to
    window.dispatchEvent(new CustomEvent('auth:logout'));
  }

  public updateUser(user: Partial<AuthUser>): void {
    if (this.user) {
      this.user = { ...this.user, ...user };
      
      if (typeof window !== 'undefined') {
        localStorage.setItem('user', JSON.stringify(this.user));
      }

      // Dispatch custom event for other components to listen to
      window.dispatchEvent(new CustomEvent('auth:user-updated', { detail: { user: this.user } }));
    }
  }

  public isAuthenticated(): boolean {
    return !!this.user && !!this.tokens && this.tokens.expiresIn > 0;
  }

  public getUser(): AuthUser | null {
    return this.user;
  }

  public getTokens(): AuthTokens | null {
    return this.tokens;
  }

  public getAccessToken(): string | null {
    return this.tokens?.accessToken || null;
  }

  public getRefreshToken(): string | null {
    return this.tokens?.refreshToken || null;
  }

  public hasRole(role: 'admin' | 'user'): boolean {
    return this.user?.role === role;
  }

  public isAdmin(): boolean {
    return this.user?.role === 'admin';
  }

  public getTokenPayload(): any {
    if (!this.tokens?.accessToken) return null;
    
    try {
      return jwtDecode(this.tokens.accessToken);
    } catch (error) {
      console.error('Failed to decode access token:', error);
      return null;
    }
  }

  public isTokenExpired(): boolean {
    if (!this.tokens) return true;
    return this.tokens.expiresIn <= 0;
  }

  public async refreshTokenManually(): Promise<boolean> {
    try {
      await this.refreshTokens();
      return true;
    } catch (error) {
      console.error('Manual token refresh failed:', error);
      return false;
    }
  }
}

// Export singleton instance
export const auth = AuthManager.getInstance();

// Auth helper functions
export const getAuthHeaders = (): Record<string, string> => {
  const token = auth.getAccessToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  return headers;
};

export const withAuth = async (fetchFn: () => Promise<Response>): Promise<Response> => {
  if (!auth.isAuthenticated()) {
    throw new Error('User not authenticated');
  }

  // Check if token is expired and try to refresh
  if (auth.isTokenExpired()) {
    const refreshed = await auth.refreshTokenManually();
    if (!refreshed) {
      auth.logout();
      throw new Error('Session expired');
    }
  }

  return fetchFn();
};

// Hook for accessing auth state in components
export const useAuth = () => {
  const [user, setUser] = React.useState<AuthUser | null>(auth.getUser());
  const [isAuthenticated, setIsAuthenticated] = React.useState(auth.isAuthenticated());
  const [isLoading, setIsLoading] = React.useState(false);

  React.useEffect(() => {
    const handleLogin = (event: CustomEvent) => {
      setUser(event.detail.user);
      setIsAuthenticated(true);
      setIsLoading(false);
    };

    const handleLogout = () => {
      setUser(null);
      setIsAuthenticated(false);
      setIsLoading(false);
    };

    const handleUserUpdated = (event: CustomEvent) => {
      setUser(event.detail.user);
    };

    window.addEventListener('auth:login', handleLogin as EventListener);
    window.addEventListener('auth:logout', handleLogout);
    window.addEventListener('auth:user-updated', handleUserUpdated as EventListener);

    return () => {
      window.removeEventListener('auth:login', handleLogin as EventListener);
      window.removeEventListener('auth:logout', handleLogout);
      window.removeEventListener('auth:user-updated', handleUserUpdated as EventListener);
    };
  }, []);

  const login = React.useCallback(async (email: string, password: string, rememberMe?: boolean) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password, rememberMe }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error?.message || 'Login failed');
      }

      const data = await response.json();
      auth.authenticate(data.data.user, data.data.tokens);
      
      return data.data;
    } catch (error) {
      setIsLoading(false);
      throw error;
    }
  }, []);

  const register = React.useCallback(async (userData: {
    email: string;
    password: string;
    firstName: string;
    lastName: string;
    dateOfBirth?: string;
    gender?: string;
  }) => {
    setIsLoading(true);
    
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error?.message || 'Registration failed');
      }

      const data = await response.json();
      auth.authenticate(data.data.user, data.data.tokens);
      
      return data.data;
    } catch (error) {
      setIsLoading(false);
      throw error;
    }
  }, []);

  const logout = React.useCallback(() => {
    auth.logout();
  }, []);

  const updateUser = React.useCallback((userData: Partial<AuthUser>) => {
    auth.updateUser(userData);
  }, []);

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    updateUser,
    hasRole: auth.hasRole.bind(auth),
    isAdmin: auth.isAdmin.bind(auth),
  };
};

// React import for the hook
import React from 'react';
