/**
 * Auth Feature - Mutations
 * React Query мутации для аутентификации
 */

import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { UseMutationOptions } from '@tanstack/react-query';
import { authFeatureApi } from '../api/authApi';
import { userKeys } from '@/entities/user';
import type { User, LoginRequest, RegisterRequest } from '@/shared/types';

/**
 * Мутация: Вход
 */
export const useLogin = (
  options?: UseMutationOptions<User, Error, LoginRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation<User, Error, LoginRequest>({
    mutationFn: authFeatureApi.login,
    onSuccess: (user) => {
      // Сохраняем пользователя в кеш
      queryClient.setQueryData(userKeys.current(), user);
    },
    ...options,
  });
};

/**
 * Мутация: Регистрация
 */
export const useRegister = (
  options?: UseMutationOptions<User, Error, RegisterRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation<User, Error, RegisterRequest>({
    mutationFn: authFeatureApi.register,
    onSuccess: (user) => {
      // Сохраняем пользователя в кеш
      queryClient.setQueryData(userKeys.current(), user);
    },
    ...options,
  });
};

/**
 * Мутация: Выход
 */
export const useLogout = (options?: UseMutationOptions<void, Error, void>) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, void>({
    mutationFn: authFeatureApi.logout,
    onSuccess: () => {
      // Очищаем весь кеш при выходе
      queryClient.clear();
    },
    ...options,
  });
};

