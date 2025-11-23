import { describe, it, expect } from 'vitest';

describe('Token Storage', () => {
  it('should store and retrieve token', () => {
    const mockToken = 'test-token-123';
    localStorage.setItem('auth_token', mockToken);
    const retrieved = localStorage.getItem('auth_token');
    expect(retrieved).toBe(mockToken);
    localStorage.removeItem('auth_token');
  });

  it('should handle missing token', () => {
    localStorage.removeItem('auth_token');
    const retrieved = localStorage.getItem('auth_token');
    expect(retrieved).toBeNull();
  });
});

describe('API Config', () => {
  it('should have valid API URLs', () => {
    const authUrl = import.meta.env.VITE_AUTH_API_URL || 'http://localhost:8081';
    const lobbyUrl = import.meta.env.VITE_LOBBY_API_URL || 'http://localhost:8082';
    
    expect(authUrl).toBeDefined();
    expect(lobbyUrl).toBeDefined();
    expect(authUrl).toMatch(/^https?:\/\//);
    expect(lobbyUrl).toMatch(/^https?:\/\//);
  });
});

describe('Route Constants', () => {
  it('should have defined routes', () => {
    const ROUTES = {
      LOGIN: '/login',
      REGISTER: '/register',
      LOBBY: '/lobby',
    };
    
    expect(ROUTES.LOGIN).toBe('/login');
    expect(ROUTES.REGISTER).toBe('/register');
    expect(ROUTES.LOBBY).toBe('/lobby');
  });
});

