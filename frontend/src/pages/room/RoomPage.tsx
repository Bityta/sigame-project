import { useParams, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useRoom, useLeaveRoom, useStartGame, useRoomEvents } from '@/entities/room';
import { useCurrentUser } from '@/entities/user';
import { RoomSettingsComponent } from '@/features/room';
import { Button, Card, Spinner } from '@/shared/ui';
import { ROUTES, TEXTS } from '@/shared/config';
import './RoomPage.css';

export const RoomPage = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const [copySuccess, setCopySuccess] = useState(false);
  
  const { data: room, isLoading } = useRoom(roomId!);
  const { data: user } = useCurrentUser();
  const leaveRoomMutation = useLeaveRoom();
  const startGameMutation = useStartGame({
    onSuccess: (response) => {
      navigate(ROUTES.GAME(response.gameSessionId));
    },
  });

  useRoomEvents(roomId, {
    onGameStarted: (event) => {
      navigate(ROUTES.GAME(event.gameId));
    },
    onRoomClosed: () => {
      navigate(ROUTES.LOBBY);
    },
  });

  const isHost = room && user && room.hostId === user.id;

  const handleLeave = () => {
    if (roomId) {
      leaveRoomMutation.mutate(roomId, {
        onSuccess: () => {
          navigate(ROUTES.LOBBY);
        },
      });
    }
  };

  const handleStart = () => {
    if (roomId) {
      startGameMutation.mutate(roomId);
    }
  };

  const handleCopyCode = async () => {
    if (!room) return;
    
    try {
      await navigator.clipboard.writeText(room.roomCode);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const handleCopyLink = async () => {
    if (!room) return;
    
    const roomUrl = `${window.location.origin}${ROUTES.ROOM(room.id)}`;
    
    try {
      await navigator.clipboard.writeText(roomUrl);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  if (isLoading) {
    return (
      <div className="room-page">
        <Spinner center size="large" />
      </div>
    );
  }

  if (!room) {
    return (
      <div className="room-page">
        <Card className="room-page__error">
          <h2>{TEXTS.ROOM.ROOM_NOT_FOUND}</h2>
          <Button onClick={() => navigate(ROUTES.LOBBY)}>
            {TEXTS.ROOM.RETURN_TO_LOBBY}
          </Button>
        </Card>
      </div>
    );
  }

  const canStart = isHost && room.currentPlayers >= 2 && room.status === 'waiting';

  return (
    <div className="room-page">
      <header className="room-page__header">
        <h1 className="room-page__title">{room.name}</h1>
        <div className="room-page__share">
          <div 
            className="room-page__code"
            onClick={handleCopyCode}
            title="Нажмите, чтобы скопировать код"
          >
            {TEXTS.ROOM.ROOM_CODE} 
            <span className="room-page__code-value">
              {room.roomCode}
              {copySuccess && <span className="room-page__copy-hint"> ✓</span>}
            </span>
          </div>
          {room.isPublic && (
            <Button
              variant="secondary"
              size="small"
              onClick={handleCopyLink}
            >
              {copySuccess ? '✓ Скопировано' : '🔗 Поделиться ссылкой'}
            </Button>
          )}
        </div>
      </header>

      <div className="room-page__content">
        <aside className="room-page__sidebar">
          <Card padding="medium">
            <h2 className="room-page__subtitle">
              {TEXTS.ROOM.PLAYERS(room.currentPlayers, room.maxPlayers)}
            </h2>
            <div className="room-page__players">
              {room.players.map((player) => (
                <div key={player.userId} className="room-page__player">
                  <span className="room-page__player-name">
                    {player.username}
                    {player.role === 'host' && ' 👑'}
                  </span>
                </div>
              ))}
            </div>
          </Card>

          <Card padding="medium" className="room-page__actions">
            {isHost ? (
              <>
                <Button
                  variant="primary"
                  size="large"
                  fullWidth
                  onClick={handleStart}
                  disabled={!canStart}
                  isLoading={startGameMutation.isPending}
                >
                  {TEXTS.ROOM.START_GAME}
                </Button>
                {!canStart && room.currentPlayers < 2 && (
                  <p className="room-page__hint">
                    {TEXTS.ROOM.MIN_PLAYERS_REQUIRED}
                  </p>
                )}
              </>
            ) : (
              <div className="room-page__waiting">
                <p>{TEXTS.ROOM.WAITING_START}</p>
                <p className="room-page__hint">
                  {TEXTS.ROOM.HOST_WILL_START}
                </p>
              </div>
            )}

            <Button
              variant="danger"
              fullWidth
              onClick={handleLeave}
              isLoading={leaveRoomMutation.isPending}
            >
              {TEXTS.ROOM.LEAVE_ROOM}
            </Button>
          </Card>
        </aside>

        <main className="room-page__main">
          <RoomSettingsComponent room={room} isHost={isHost || false} />
        </main>
      </div>
    </div>
  );
};
