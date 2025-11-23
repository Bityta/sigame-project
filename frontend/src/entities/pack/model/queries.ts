/**
 * Pack Entity - Queries
 * React Query хуки для паков
 */

import { useQuery } from '@tanstack/react-query';
import type { UseQueryOptions } from '@tanstack/react-query';
import { packApi } from '../api/packApi';
import type { Pack, PackListQuery } from '@/shared/types';

// Query Keys
export const packKeys = {
  all: ['pack'] as const,
  lists: () => [...packKeys.all, 'list'] as const,
  list: (filters?: PackListQuery) => [...packKeys.lists(), filters] as const,
  details: () => [...packKeys.all, 'detail'] as const,
  detail: (id: string) => [...packKeys.details(), id] as const,
  search: (query: string) => [...packKeys.all, 'search', query] as const,
};

/**
 * Хук для получения списка паков
 */
export const usePacks = (
  filters?: PackListQuery,
  options?: Omit<UseQueryOptions<Pack[], Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<Pack[], Error>({
    queryKey: packKeys.list(filters),
    queryFn: () => packApi.getPacks(filters),
    ...options,
  });
};

/**
 * Хук для получения пака по ID
 */
export const usePack = (
  id: string,
  options?: Omit<UseQueryOptions<Pack, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<Pack, Error>({
    queryKey: packKeys.detail(id),
    queryFn: () => packApi.getPackById(id),
    enabled: !!id,
    ...options,
  });
};

/**
 * Хук для поиска паков
 */
export const useSearchPacks = (
  query: string,
  options?: Omit<UseQueryOptions<Pack[], Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<Pack[], Error>({
    queryKey: packKeys.search(query),
    queryFn: () => packApi.searchPacks(query),
    enabled: query.length >= 2, // Ищем только если длина >= 2
    ...options,
  });
};

