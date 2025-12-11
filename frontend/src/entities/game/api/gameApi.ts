/**
 * Game Entity - API
 * API методы для работы с игрой
 */

import axios from 'axios';
import { API_CONFIG } from '@/shared/config';
import { TokenStorage } from '@/shared/lib';

// Create game axios instance
const gameApi = axios.create({
  baseURL: API_CONFIG.GAME_BASE_URL,
  timeout: API_CONFIG.TIMEOUT.DEFAULT,
});

// Add auth interceptor
gameApi.interceptors.request.use((config) => {
  const token = TokenStorage.getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export interface MyActiveGameResponse {
  hasActiveGame: boolean;
  gameId?: string;
  websocketUrl?: string;
  status?: string;
}

export const gameApiService = {
  /**
   * Get user's active game
   */
  async getMyActiveGame(): Promise<MyActiveGameResponse> {
    const response = await gameApi.get<MyActiveGameResponse>(
      API_CONFIG.ENDPOINTS.GAME.MY_ACTIVE
    );
    return response.data;
  },
};

