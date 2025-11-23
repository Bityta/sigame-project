/**
 * User Entity - Queries
 * React Query хуки для пользователей
 */

import { useQuery } from '@tanstack/react-query';
import type { UseQueryOptions } from '@tanstack/react-query';
import { userApi } from '../api/userApi';
import type { User } from '@/shared/types';

// Query Keys
export const userKeys = {
  all: ['user'] as const,
  current: () => [...userKeys.all, 'current'] as const,
  checkUsername: (username: string) => [...userKeys.all, 'check', username] as const,
};

/**
 * Хук для получения текущего пользователя
 */
export const useCurrentUser = (
  options?: Omit<UseQueryOptions<User, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<User, Error>({
    queryKey: userKeys.current(),
    queryFn: userApi.getCurrentUser,
    ...options,
  });
};

/**
 * Хук для проверки доступности username
 */
export const useCheckUsername = (
  username: string,
  options?: Omit<UseQueryOptions<boolean, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<boolean, Error>({
    queryKey: userKeys.checkUsername(username),
    queryFn: () => userApi.checkUsername(username),
    enabled: username.length >= 3, // Проверяем только если длина >= 3
    ...options,
  });
};

