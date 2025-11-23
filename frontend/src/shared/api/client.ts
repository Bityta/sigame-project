/**
 * Axios API Client
 * Централизованный HTTP клиент с автоматической обработкой ошибок
 */

import axios, { AxiosError } from 'axios';
import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import { API_CONFIG } from '../config';
import { TokenStorage } from '../lib/token-storage';
import { useErrorStore } from '../lib/error-store';
import type { ApiErrorResponse } from '../types';

// Создание базового axios instance
const createApiClient = (baseURL: string): AxiosInstance => {
  const client = axios.create({
    baseURL,
    timeout: API_CONFIG.TIMEOUT.DEFAULT,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor - добавляем токен
  client.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const token = TokenStorage.getAccessToken();
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor - обрабатываем ошибки глобально
  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError<ApiErrorResponse>) => {
      const { response } = error;

      // Формируем читаемое сообщение об ошибке
      let errorMessage = 'Что-то пошло не так';
      let errorCode: string | undefined;

      if (response) {
        // Ошибки 4xx и 5xx
        errorCode = response.data?.error || `${response.status}`;
        
        // Используем сообщение от сервера если есть
        if (response.data?.message) {
          errorMessage = response.data.message;
        } else if (response.status === 401) {
          errorMessage = 'Требуется авторизация';
        } else if (response.status === 403) {
          errorMessage = 'Доступ запрещен';
        } else if (response.status === 404) {
          errorMessage = 'Ресурс не найден';
        } else if (response.status === 409) {
          errorMessage = 'Конфликт данных';
        } else if (response.status >= 500) {
          errorMessage = 'Ошибка сервера. Попробуйте позже';
        }
      } else if (error.request) {
        // Сеть недоступна
        errorMessage = 'Нет соединения с сервером';
        errorCode = 'NETWORK_ERROR';
      }

      // Отображаем ошибку в красной плашке
      useErrorStore.getState().setError(errorMessage, errorCode);

      // Если 401 и это не запрос на логин - очищаем токены
      if (response?.status === 401 && !error.config?.url?.includes('/login')) {
        TokenStorage.clearTokens();
        // Можем редиректить на /login если нужно
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }

      return Promise.reject(error);
    }
  );

  return client;
};

// Создаем клиенты для каждого микросервиса
export const authApi = createApiClient(API_CONFIG.AUTH_BASE_URL);
export const lobbyApi = createApiClient(API_CONFIG.LOBBY_BASE_URL);
export const gameApi = createApiClient(API_CONFIG.GAME_BASE_URL);
export const packApi = createApiClient(API_CONFIG.PACK_BASE_URL);

// Экспортируем для использования в хуках и сервисах
export { createApiClient };

