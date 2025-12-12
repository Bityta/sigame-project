/**
 * Game Entity - API
 * API методы для работы с игрой
 */

import { API_CONFIG } from '@/shared/config';
import { gameApi } from '@/shared/api/client';

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

