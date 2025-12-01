/**
 * API Configuration
 * Централизованная конфигурация для всех API endpoints
 */

const getEnvVar = (key: string, defaultValue: string): string => {
  if (typeof window === 'undefined') return defaultValue;
  return (window as any)[key] || import.meta.env[key] || defaultValue;
};

export const API_CONFIG = {
  // Base URLs для микросервисов
  AUTH_BASE_URL: getEnvVar('VITE_AUTH_API_URL', 'http://localhost:8001'),
  LOBBY_BASE_URL: getEnvVar('VITE_LOBBY_API_URL', 'http://localhost:8002'),
  GAME_BASE_URL: getEnvVar('VITE_GAME_API_URL', 'http://localhost:8003'),
  PACK_BASE_URL: getEnvVar('VITE_PACK_API_URL', 'http://localhost:8004'),
  
  // Endpoints
  ENDPOINTS: {
    // Auth
    AUTH: {
      LOGIN: '/auth/login',
      REGISTER: '/auth/register',
      LOGOUT: '/auth/logout',
      REFRESH: '/auth/refresh',
      ME: '/auth/me',
      CHECK_USERNAME: '/auth/check-username',
    },
    // Lobby
    LOBBY: {
      ROOMS: '/api/lobby/rooms',
      MY_ROOMS: '/api/lobby/rooms/my',
      ROOM_BY_ID: (id: string) => `/api/lobby/rooms/${id}`,
      ROOM_BY_CODE: (code: string) => `/api/lobby/rooms/code/${code}`,
      JOIN_ROOM: (id: string) => `/api/lobby/rooms/${id}/join`,
      LEAVE_ROOM: (id: string) => `/api/lobby/rooms/${id}/leave`,
      START_ROOM: (id: string) => `/api/lobby/rooms/${id}/start`,
      UPDATE_SETTINGS: (id: string) => `/api/lobby/rooms/${id}/settings`,
      KICK_PLAYER: (id: string) => `/api/lobby/rooms/${id}/kick`,
      TRANSFER_HOST: (id: string) => `/api/lobby/rooms/${id}/transfer-host`,
    },
    // Game
    GAME: {
      SESSION: (id: string) => `/api/game/${id}`,
      WS: (id: string) => `/api/game/${id}/ws`,
    },
    // Packs
    PACK: {
      LIST: '/api/packs',
      BY_ID: (id: string) => `/api/packs/${id}`,
      SEARCH: '/api/packs/search',
    },
  },
  
  // Таймауты
  TIMEOUT: {
    DEFAULT: 10000, // 10 секунд
    UPLOAD: 60000,  // 60 секунд для загрузки файлов
  },
} as const;

