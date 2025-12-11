/**
 * Game Entity - useMyActiveGame Hook
 * Hook для получения активной игры пользователя
 */

import { useQuery } from '@tanstack/react-query';
import { gameApiService } from '../api/gameApi';
import type { MyActiveGameResponse } from '../api/gameApi';

export const gameKeys = {
  all: ['game'] as const,
  myActive: () => [...gameKeys.all, 'my-active'] as const,
};

export const useMyActiveGame = () => {
  return useQuery<MyActiveGameResponse>({
    queryKey: gameKeys.myActive(),
    queryFn: () => gameApiService.getMyActiveGame(),
    staleTime: 30000, // 30 seconds
    refetchInterval: 60000, // 1 minute
  });
};

