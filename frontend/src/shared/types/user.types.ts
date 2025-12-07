/**
 * User & Auth Domain Types
 */

export interface User {
  id: string;
  username: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface UsernameCheckResponse {
  available: boolean;
  reason?: string;
}

export interface LoginResponse {
  data: User;
  tokens: AuthTokens;
}

export interface RegisterResponse {
  data: User;
  tokens: AuthTokens;
}

