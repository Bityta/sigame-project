/**
 * Auth Provider
 * Инициализация аутентификации при загрузке приложения
 */

import { useEffect } from 'react';
import type { ReactNode } from 'react';
import { useAuthStore } from '@/features/auth';
import { Spinner } from '@/shared/ui';

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const { checkAuth, isInitialized } = useAuthStore();

  useEffect(() => {
    // Проверяем аутентификацию при загрузке
    checkAuth();
  }, [checkAuth]);

  if (!isInitialized) {
    return (
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      }}>
        <Spinner size="large" />
      </div>
    );
  }

  return <>{children}</>;
};

