/**
 * Public Route
 * Публичный роут, редиректит авторизованных пользователей
 */

import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/features/auth';
import { ROUTES } from '@/shared/config';

interface PublicRouteProps {
  children: ReactNode;
}

export const PublicRoute = ({ children }: PublicRouteProps) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (isAuthenticated) {
    return <Navigate to={ROUTES.LOBBY} replace />;
  }

  return <>{children}</>;
};

