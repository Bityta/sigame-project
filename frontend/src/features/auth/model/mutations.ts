/**
 * Auth Feature - Mutations
 * React Query мутации для аутентификации
 */

import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { UseMutationOptions } from '@tanstack/react-query';
import { authFeatureApi } from '../api/authApi';
import { userKeys } from '@/entities/user';
import type { LoginRequest, RegisterRequest } from '@/shared/types';

/**
 * Мутация: Вход
 */
export const useLogin = (
  options?: UseMutationOptions<void, Error, LoginRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, LoginRequest>({
    mutationFn: authFeatureApi.login,
    onSuccess: () => {
      // Инвалидируем запрос currentUser - React Query вызовет /auth/me
      queryClient.invalidateQueries({ queryKey: userKeys.current() });
    },
    ...options,
  });
};

/**
 * Мутация: Регистрация
 */
export const useRegister = (
  options?: UseMutationOptions<void, Error, RegisterRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, RegisterRequest>({
    mutationFn: authFeatureApi.register,
    onSuccess: () => {
      // Инвалидируем запрос currentUser - React Query вызовет /auth/me
      queryClient.invalidateQueries({ queryKey: userKeys.current() });
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

