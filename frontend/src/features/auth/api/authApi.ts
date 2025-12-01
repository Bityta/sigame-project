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
  User,
} from '@/shared/types';
import type { 
  LoginApiResponse, 
  RegisterApiResponse,
  AuthApiUser 
} from '@/shared/types/api.types';

/**
 * Преобразует ответ сервера в формат User
 */
const mapApiUserToUser = (apiUser: AuthApiUser): User => {
  return {
    id: apiUser.id,
    username: apiUser.username,
    createdAt: apiUser.created_at,
    updatedAt: apiUser.created_at,
  };
};

export const authFeatureApi = {
  /**
   * Войти в систему
   */
  async login(credentials: LoginRequest): Promise<User> {
    const response = await authApi.post<LoginApiResponse>(
      API_CONFIG.ENDPOINTS.AUTH.LOGIN,
      credentials
    );
    
    // Сохраняем токены (сервер возвращает access_token и refresh_token)
    TokenStorage.setTokens(
      response.data.access_token,
      response.data.refresh_token
    );
    
    // Преобразуем и возвращаем пользователя
    const user = mapApiUserToUser(response.data.user);
    
    // Сохраняем пользователя в localStorage
    TokenStorage.setUser(user);
    
    return user;
  },

  /**
   * Зарегистрироваться
   */
  async register(data: RegisterRequest): Promise<User> {
    const response = await authApi.post<RegisterApiResponse>(
      API_CONFIG.ENDPOINTS.AUTH.REGISTER,
      data
    );
    
    // Сохраняем токены (сервер возвращает access_token и refresh_token)
    TokenStorage.setTokens(
      response.data.access_token,
      response.data.refresh_token
    );
    
    // Преобразуем и возвращаем пользователя
    const user = {
      id: response.data.user.id,
      username: response.data.user.username,
      createdAt: response.data.user.created_at,
      updatedAt: response.data.user.created_at,
    };
    
    // Сохраняем пользователя в localStorage
    TokenStorage.setUser(user);
    
    return user;
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

