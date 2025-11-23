/**
 * Pack Entity - API
 * API методы для работы с паками
 */

import { packApi as packApiClient } from '@/shared/api';
import { API_CONFIG } from '@/shared/config';
import type { Pack, PackListQuery, PaginatedResponse } from '@/shared/types';

export const packApi = {
  /**
   * Получить список паков
   */
  async getPacks(query?: PackListQuery): Promise<Pack[]> {
    const response = await packApiClient.get<{ packs: Pack[] }>(
      API_CONFIG.ENDPOINTS.PACK.LIST,
      { params: query }
    );
    return response.data.packs || [];
  },

  /**
   * Получить пак по ID
   */
  async getPackById(id: string): Promise<Pack> {
    const response = await packApiClient.get<Pack>(
      API_CONFIG.ENDPOINTS.PACK.BY_ID(id)
    );
    return response.data;
  },

  /**
   * Поиск паков
   */
  async searchPacks(query: string): Promise<Pack[]> {
    const response = await packApiClient.get<{ packs: Pack[] }>(
      API_CONFIG.ENDPOINTS.PACK.SEARCH,
      { params: { q: query } }
    );
    return response.data.packs || [];
  },
};

