/**
 * Auth Feature - Auth Store
 * Zustand store для состояния аутентификации (UI state)
 */

import { create } from 'zustand';
import { TokenStorage } from '@/shared/lib';

interface AuthStoreState {
  isAuthenticated: boolean;
  isInitialized: boolean;
  
  // Actions
  setAuthenticated: (value: boolean) => void;
  setInitialized: (value: boolean) => void;
  checkAuth: () => void;
}

/**
 * Store для UI состояния аутентификации
 * Используется для быстрой проверки авторизации без запроса к API
 */
export const useAuthStore = create<AuthStoreState>((set) => ({
  isAuthenticated: false,
  isInitialized: false,

  setAuthenticated: (value: boolean) => set({ isAuthenticated: value }),
  
  setInitialized: (value: boolean) => set({ isInitialized: value }),

  checkAuth: () => {
    const hasTokens = TokenStorage.hasTokens();
    set({ 
      isAuthenticated: hasTokens,
      isInitialized: true,
    });
  },
}));

