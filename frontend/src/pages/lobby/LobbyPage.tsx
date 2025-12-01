import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useLogout, useAuthStore } from '@/features/auth';
import { RoomList } from '@/features/room';
import { useCurrentUser } from '@/entities/user';
import { roomApi, useMyActiveRoom, useLeaveRoom } from '@/entities/room';
import { Button, Card, Spinner } from '@/shared/ui';
import { ROUTES, TEXTS } from '@/shared/config';
import { useErrorStore } from '@/shared/lib/error-store';
import './LobbyPage.css';

export const LobbyPage = () => {
  const navigate = useNavigate();
  const { data: user, isLoading: userLoading } = useCurrentUser();
  const { data: activeRoom, isLoading: activeRoomLoading } = useMyActiveRoom();
  const leaveRoomMutation = useLeaveRoom();
  const setAuthenticated = useAuthStore((state) => state.setAuthenticated);
  const setError = useErrorStore((state) => state.setError);
  
  const hasActiveRoom = !!activeRoom;
  
  const logoutMutation = useLogout({
    onSuccess: () => {
      setAuthenticated(false);
      navigate(ROUTES.LOGIN);
    },
  });

  const [roomCode, setRoomCode] = useState('');
  const [isSearching, setIsSearching] = useState(false);

  const handleJoinByCode = async () => {
    if (hasActiveRoom) return;
    
    const code = roomCode.trim().toUpperCase();
    if (!code) return;

    setIsSearching(true);
    try {
      const room = await roomApi.getRoomByCode(code);
      navigate(ROUTES.ROOM(room.id));
    } catch (error: any) {
      const errorMessage = error?.response?.status === 404
        ? `Комната с кодом "${code}" не найдена`
        : 'Не удалось найти комнату. Проверьте код и попробуйте снова';
      
      setError(errorMessage, error?.response?.status?.toString());
    } finally {
      setIsSearching(false);
    }
  };

  const handleLeaveActiveRoom = () => {
    if (!activeRoom) return;
    leaveRoomMutation.mutate(activeRoom.id);
  };

  if (userLoading) {
    return (
      <div className="lobby-page">
        <Spinner center size="large" />
      </div>
    );
  }

  return (
    <div className="lobby-page">
      <header className="lobby-page__header">
        <h1 className="lobby-page__title">{TEXTS.APP_NAME}</h1>
        <div className="lobby-page__user">
          <span className="lobby-page__username">
            {TEXTS.LOBBY.WELCOME(user?.username || '')}
          </span>
          <Button
            variant="ghost"
            size="small"
            onClick={() => logoutMutation.mutate()}
            isLoading={logoutMutation.isPending}
            className="lobby-page__logout-btn"
          >
            {TEXTS.AUTH.LOGOUT}
          </Button>
        </div>
      </header>

      <div className="lobby-page__content">
        <aside className="lobby-page__sidebar">
          {hasActiveRoom && activeRoom && (
            <Card padding="medium" className="lobby-page__active-room">
              <div className="lobby-page__active-room-badge">⚠️ Активная комната</div>
              <h3 className="lobby-page__active-room-name">{activeRoom.name}</h3>
              <div className="lobby-page__active-room-info">
                <div className="lobby-page__active-room-row">
                  <span>Код:</span>
                  <span className="lobby-page__active-room-code">{activeRoom.roomCode}</span>
                </div>
                <div className="lobby-page__active-room-row">
                  <span>Игроки:</span>
                  <span>{activeRoom.currentPlayers}/{activeRoom.maxPlayers}</span>
                </div>
                <div className="lobby-page__active-room-row">
                  <span>Статус:</span>
                  <span>{activeRoom.status === 'waiting' ? 'Ожидание' : activeRoom.status}</span>
                </div>
              </div>
              <Button
                variant="primary"
                fullWidth
                size="large"
                onClick={() => navigate(ROUTES.ROOM(activeRoom.id))}
              >
                Вернуться в комнату
              </Button>
              <button
                className="lobby-page__leave-link"
                onClick={handleLeaveActiveRoom}
                disabled={leaveRoomMutation.isPending}
              >
                {leaveRoomMutation.isPending ? 'Выход...' : 'Покинуть комнату'}
              </button>
            </Card>
          )}

          <Card padding="medium">
            <h2 className="lobby-page__sidebar-title">{TEXTS.LOBBY.QUICK_ACTIONS}</h2>
            
            <div className="lobby-page__actions">
              <Button
                variant="primary"
                fullWidth
                size="large"
                onClick={() => navigate('/lobby/create')}
                disabled={hasActiveRoom}
                title={hasActiveRoom ? 'Сначала покиньте текущую комнату' : undefined}
              >
                {TEXTS.LOBBY.CREATE_ROOM}
              </Button>

              <div className="lobby-page__join-code">
                <input
                  type="text"
                  placeholder={TEXTS.LOBBY.ROOM_CODE_PLACEHOLDER}
                  value={roomCode}
                  onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
                  onKeyDown={(e) => e.key === 'Enter' && !isSearching && !hasActiveRoom && handleJoinByCode()}
                  className="lobby-page__code-input"
                  maxLength={6}
                  disabled={isSearching || hasActiveRoom}
                />
                <Button
                  variant="secondary"
                  fullWidth
                  onClick={handleJoinByCode}
                  disabled={!roomCode.trim() || isSearching || hasActiveRoom}
                  isLoading={isSearching}
                  title={hasActiveRoom ? 'Сначала покиньте текущую комнату' : undefined}
                >
                  {TEXTS.LOBBY.JOIN_BY_CODE}
                </Button>
              </div>
            </div>
          </Card>
        </aside>

        <main className="lobby-page__main">
          <RoomList hasActiveRoom={hasActiveRoom} />
        </main>
      </div>
    </div>
  );
};
