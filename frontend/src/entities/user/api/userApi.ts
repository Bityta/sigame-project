/**
 * User Entity - API
 * API методы для работы с пользователями
 */

import { authApi } from '@/shared/api';
import { API_CONFIG } from '@/shared/config';
import { TokenStorage } from '@/shared/lib';
import type { User } from '@/shared/types';

export const userApi = {
  /**
   * Получить информацию о текущем пользователе
   * Берём из localStorage, так как сервер не возвращает /auth/me
   */
  async getCurrentUser(): Promise<User> {
    const user = TokenStorage.getUser();
    
    if (!user) {
      throw new Error('User not found in storage');
    }
    
    return user;
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

