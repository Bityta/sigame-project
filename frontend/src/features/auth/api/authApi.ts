/**
 * Auth Feature - API
 * API методы для аутентификации
 */

import { authApi } from '@/shared/api';
import { API_CONFIG } from '@/shared/config';
import { TokenStorage } from '@/shared/lib';
import type {
  LoginRequest,
  RegisterRequest,
} from '@/shared/types';
import type { 
  LoginApiResponse, 
  RegisterApiResponse,
} from '@/shared/types/api.types';

export const authFeatureApi = {
  /**
   * Войти в систему
   */
  async login(credentials: LoginRequest): Promise<void> {
    const response = await authApi.post<LoginApiResponse>(
      API_CONFIG.ENDPOINTS.AUTH.LOGIN,
      credentials
    );
    
    // Сохраняем токены (сервер возвращает access_token и refresh_token)
    TokenStorage.setTokens(
      response.data.access_token,
      response.data.refresh_token
    );
  },

  /**
   * Зарегистрироваться
   */
  async register(data: RegisterRequest): Promise<void> {
    const response = await authApi.post<RegisterApiResponse>(
      API_CONFIG.ENDPOINTS.AUTH.REGISTER,
      data
    );
    
    // Сохраняем токены (сервер возвращает access_token и refresh_token)
    TokenStorage.setTokens(
      response.data.access_token,
      response.data.refresh_token
    );
  },

  /**
   * Выйти из системы
   */
  async logout(): Promise<void> {
    try {
      await authApi.post(API_CONFIG.ENDPOINTS.AUTH.LOGOUT);
    } finally {
      // Всегда очищаем токены, даже если запрос упал
      TokenStorage.clearTokens();
    }
  },
};

