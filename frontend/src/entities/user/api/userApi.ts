/**
 * User Entity - API
 * API методы для работы с пользователями
 */

import { authApi } from '@/shared/api';
import { API_CONFIG } from '@/shared/config';
import type { User } from '@/shared/types';

interface UserApiResponse {
  id: string;
  username: string;
  avatar_url?: string;
  created_at: string;
}

export const userApi = {
  /**
   * Получить информацию о текущем пользователе
   * Запрашивает данные с сервера по JWT токену
   */
  async getCurrentUser(): Promise<User> {
    const response = await authApi.get<UserApiResponse>(
      API_CONFIG.ENDPOINTS.AUTH.ME
    );
    
    return {
      id: response.data.id,
      username: response.data.username,
      avatarUrl: response.data.avatar_url,
      createdAt: response.data.created_at,
      updatedAt: response.data.created_at,
    };
  },

  /**
   * Проверить доступность username
   */
  async checkUsername(username: string): Promise<boolean> {
    const response = await authApi.get<{ available: boolean }>(
      API_CONFIG.ENDPOINTS.AUTH.CHECK_USERNAME,
      { params: { username } }
    );
    return response.data.available;
  },
};

