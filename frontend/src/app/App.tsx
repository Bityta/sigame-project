/**
 * App
 * Главный компонент приложения
 */

import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryProvider, AuthProvider } from './providers';
import { ProtectedRoute, PublicRoute } from './routes';
import { ErrorBanner } from '@/shared/ui';
import { ROUTES } from '@/shared/config';

// Pages
import { LoginPage } from '@/pages/login';
import { RegisterPage } from '@/pages/register';
import { LobbyPage } from '@/pages/lobby';
import { RoomPage } from '@/pages/room';
import { GamePage } from '@/pages/game';
import { CreateRoomForm } from '@/features/room';

import './styles/global.css';

export const App = () => {
  return (
    <QueryProvider>
      <AuthProvider>
        <BrowserRouter>
          {/* Глобальный ErrorBanner для всех ошибок API */}
          <ErrorBanner />

          <Routes>
            {/* Публичные роуты */}
            <Route
              path={ROUTES.LOGIN}
              element={
                <PublicRoute>
                  <LoginPage />
                </PublicRoute>
              }
            />
            <Route
              path={ROUTES.REGISTER}
              element={
                <PublicRoute>
                  <RegisterPage />
                </PublicRoute>
              }
            />

            {/* Защищенные роуты */}
            <Route
              path={ROUTES.LOBBY}
              element={
                <ProtectedRoute>
                  <LobbyPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/lobby/create"
              element={
                <ProtectedRoute>
                  <div style={{
                    minHeight: '100vh',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '2rem',
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  }}>
                    <CreateRoomForm />
                  </div>
                </ProtectedRoute>
              }
            />
            <Route
              path="/room/:roomId"
              element={
                <ProtectedRoute>
                  <RoomPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/game/:gameId"
              element={
                <ProtectedRoute>
                  <GamePage />
                </ProtectedRoute>
              }
            />

            {/* Default redirect */}
            <Route path={ROUTES.HOME} element={<Navigate to={ROUTES.LOBBY} replace />} />
            <Route path="*" element={<Navigate to={ROUTES.LOBBY} replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </QueryProvider>
  );
};

