/**
 * Global Error Store
 * Zustand store для глобальной обработки ошибок
 * Показывает красную плашку сверху при любой ошибке API
 */

import { create } from 'zustand';
import { TIMINGS } from '../config';

export interface AppError {
  message: string;
  code?: string;
  timestamp: number;
}

interface ErrorStoreState {
  error: AppError | null;
  
  // Actions
  setError: (message: string, code?: string) => void;
  clearError: () => void;
}

export const useErrorStore = create<ErrorStoreState>((set) => ({
  error: null,

  setError: (message: string, code?: string) => {
    const error: AppError = {
      message,
      code,
      timestamp: Date.now(),
    };
    
    set({ error });

    // Автоматически убираем ошибку через 5 секунд
    setTimeout(() => {
      set((state) => {
        // Убираем только если это та же ошибка (по timestamp)
        if (state.error?.timestamp === error.timestamp) {
          return { error: null };
        }
        return state;
      });
    }, TIMINGS.ERROR_DISMISS);
  },

  clearError: () => set({ error: null }),
}));

