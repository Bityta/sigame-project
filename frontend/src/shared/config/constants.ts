/**
 * Application Constants
 * Все константы приложения в одном месте
 */

// Локальное хранилище ключи
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'sigame_access_token',
  REFRESH_TOKEN: 'sigame_refresh_token',
  USER: 'sigame_user',
  USER_PREFERENCES: 'sigame_user_prefs',
} as const;

// Настройки игры по умолчанию
export const DEFAULT_ROOM_SETTINGS = {
  timeForAnswer: 30,      // секунды
  timeForChoice: 10,      // секунды
  allowWrongAnswer: true,
  showRightAnswer: true,
} as const;

// Лимиты
export const LIMITS = {
  USERNAME: {
    MIN: 3,
    MAX: 20,
  },
  PASSWORD: {
    MIN: 6,
    MAX: 100,
  },
  ROOM_NAME: {
    MIN: 3,
    MAX: 50,
  },
  MAX_PLAYERS: {
    MIN: 2,
    MAX: 8,
  },
} as const;

// Таймауты и интервалы
export const TIMINGS = {
  ERROR_DISMISS: 5000,    // 5 секунд показываем ошибку
  ROOM_POLLING: 3000,     // 3 секунды обновление комнаты
  RECONNECT_DELAY: 1000,  // 1 секунда до переподключения WS
} as const;

// Роуты приложения
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  LOBBY: '/lobby',
  JOIN_BY_CODE: (code: string) => `/join/${code}`,
  ROOM: (id: string) => `/room/${id}`,
  GAME: (id: string) => `/game/${id}`,
} as const;

// Статусы и роли
export const ROOM_STATUS = {
  WAITING: 'waiting',
  STARTING: 'starting',
  PLAYING: 'playing',
  FINISHED: 'finished',
  CANCELLED: 'cancelled',
} as const;

export const PLAYER_ROLES = {
  HOST: 'host',
  PLAYER: 'player',
  SPECTATOR: 'spectator',
} as const;

export const GAME_STATUS = {
  WAITING: 'waiting',
  ROUND_START: 'round_start',
  QUESTION_SELECT: 'question_select',
  QUESTION_SHOW: 'question_show',
  BUTTON_PRESS: 'button_press',
  ANSWERING: 'answering',
  ANSWER_JUDGING: 'answer_judging',
  ROUND_END: 'round_end',
  GAME_END: 'game_end',
  FINISHED: 'finished',
} as const;

