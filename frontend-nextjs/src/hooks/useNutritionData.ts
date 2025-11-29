import { useCallback, useEffect, useRef, useState } from 'react';
import { nutritionService, type RecipesParams, type WorkoutsParams, type ComplaintsParams, type MetabolismParams, type DrugsNutritionParams } from '@/lib/api/services/nutrition.service';

interface HookState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

function useAsync<T>(fn: () => Promise<T>, deps: unknown[] = []): HookState<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);
  const retryRef = useRef(0);

  const run = useCallback(() => {
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setLoading(true);
    setError(null);

    const attempt = async () => {
      try {
        const result = await fn();
        if (!controller.signal.aborted) {
          setData(result);
          setLoading(false);
          setError(null);
        }
      } catch (e: any) {
        if (controller.signal.aborted) return;
        retryRef.current += 1;
        if (retryRef.current <= 3) {
          const backoff = 500 * Math.pow(2, retryRef.current - 1);
          setTimeout(attempt, backoff);
        } else {
          setLoading(false);
          setError(e?.message || 'Unexpected error');
        }
      }
    };

    attempt();
  }, deps);

  useEffect(() => {
    run();
    return () => abortRef.current?.abort();
  }, [run]);

  return { data, loading, error, refetch: run };
}

interface PaginatedHookState<T> extends HookState<T> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  } | null;
  goToPage: (page: number) => void;
  nextPage: () => void;
  prevPage: () => void;
}

function usePaginatedAsync<T>(
  fn: (params?: any) => Promise<T>,
  initialParams?: any
): PaginatedHookState<T> {
  const [params, setParams] = useState(initialParams || { page: 1, limit: 20 });
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const abortRef = useRef<AbortController | null>(null);

  const fetch = useCallback(async () => {
    abortRef.current?.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    setLoading(true);
    setError(null);

    try {
      const result = await fn(params);
      if (!controller.signal.aborted) {
        setData(result);
        setLoading(false);
      }
    } catch (e: any) {
      if (!controller.signal.aborted) {
        setLoading(false);
        setError(e?.message || 'Unexpected error');
      }
    }
  }, [fn, params]);

  useEffect(() => {
    fetch();
    return () => abortRef.current?.abort();
  }, [fetch]);

  const goToPage = useCallback((page: number) => {
    setParams((prev: any) => ({ ...prev, page }));
  }, []);

  const nextPage = useCallback(() => {
    setParams((prev: any) => ({ ...prev, page: (prev.page || 1) + 1 }));
  }, []);

  const prevPage = useCallback(() => {
    setParams((prev: any) => ({ ...prev, page: Math.max(1, (prev.page || 1) - 1) }));
  }, []);

  return {
    data,
    loading,
    error,
    refetch: fetch,
    pagination: (data as any)?.pagination || null,
    goToPage,
    nextPage,
    prevPage,
  };
}

export function useRecipes(params?: RecipesParams) {
  return usePaginatedAsync((p) => nutritionService.getRecipes(p), params);
}

export function useWorkouts(params?: WorkoutsParams) {
  return usePaginatedAsync((p) => nutritionService.getWorkouts(p), params);
}

export function useComplaints(params?: ComplaintsParams) {
  return usePaginatedAsync((p) => nutritionService.getComplaints(p), params);
}

export function useMetabolism(params?: MetabolismParams) {
  return usePaginatedAsync((p) => nutritionService.getMetabolism(p), params);
}

export function useDrugsNutrition(params?: DrugsNutritionParams) {
  return usePaginatedAsync((p) => nutritionService.getDrugsNutrition(p), params);
}
