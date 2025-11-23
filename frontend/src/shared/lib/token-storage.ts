/**
 * Token Storage Utility
 * Управление токенами и данными пользователя в localStorage
 */

import { STORAGE_KEYS } from '../config';
import type { User } from '../types';

export class TokenStorage {
  static getAccessToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  }

  static getRefreshToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
  }

  static setTokens(accessToken: string, refreshToken: string): void {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
  }

  static clearTokens(): void {
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER);
  }

  static hasTokens(): boolean {
    return !!this.getAccessToken() && !!this.getRefreshToken();
  }

  /**
   * Сохранить данные пользователя
   */
  static setUser(user: User): void {
    localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
  }

  /**
   * Получить данные пользователя из localStorage
   */
  static getUser(): User | null {
    const userData = localStorage.getItem(STORAGE_KEYS.USER);
    if (!userData) return null;
    
    try {
      return JSON.parse(userData);
    } catch {
      return null;
    }
  }
}

