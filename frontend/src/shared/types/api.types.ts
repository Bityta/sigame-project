/**
 * API Response Types
 */

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  size: number;
  total_pages: number;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

export interface ApiErrorResponse {
  error: string;
  message: string;
  statusCode: number;
  timestamp?: string;
  path?: string;
}

/**
 * Auth API Response Types
 * Структура ответов от auth-service (snake_case как на сервере)
 */

export interface AuthApiTokens {
  access_token: string;
  refresh_token: string;
}

export interface AuthApiUser {
  id: string;
  username: string;
  created_at: string;
}

export interface LoginApiResponse extends AuthApiTokens {
  user?: AuthApiUser;
}

export interface RegisterApiResponse extends AuthApiTokens {
  user: AuthApiUser;
}

