import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect, useRef } from 'react';
import { useRoom, useLeaveRoom, useRoomEvents, useKickPlayer, useTransferHost, useJoinRoom, useSetReady } from '@/entities/room';
import { useCurrentUser } from '@/entities/user';
import { RoomSettingsComponent } from '@/features/room';
import { Button, Card, Spinner } from '@/shared/ui';
import { ROUTES, TEXTS } from '@/shared/config';
import './RoomPage.css';

export const RoomPage = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();
  const [copySuccess, setCopySuccess] = useState(false);
  const hasJoined = useRef(false);
  
  const { data: room, isLoading, refetch } = useRoom(roomId!);
  const { data: user } = useCurrentUser();
  const joinRoomMutation = useJoinRoom();

  // Автоматически присоединяемся к комнате при входе на страницу
  useEffect(() => {
    // Ждём загрузки данных комнаты и пользователя
    if (!roomId || !user || isLoading) return;
    
    // Проверяем только когда room загружен
    const isAlreadyInRoom = room?.players?.some(p => p.userId === user.id);
    
    // Если уже в комнате - ничего не делаем, просто показываем страницу
    if (isAlreadyInRoom) return;
    
    // Если комната существует но нас в ней нет - присоединяемся
    if (room && !hasJoined.current && !joinRoomMutation.isPending) {
      hasJoined.current = true;
      joinRoomMutation.mutate(
        { id: roomId, data: {} },
        {
          onSuccess: () => {
            refetch();
          },
          onError: (error: any) => {
            console.error('Failed to join room:', error);
            // Если ошибка "уже в комнате" (409) - просто обновляем данные
            if (error?.response?.status === 409) {
              refetch();
              return;
            }
            // Иначе возвращаемся в лобби
            navigate(ROUTES.LOBBY);
          },
        }
      );
    }
  }, [roomId, user, room, isLoading]);
  const leaveRoomMutation = useLeaveRoom();
  const kickPlayerMutation = useKickPlayer();
  const transferHostMutation = useTransferHost();
  const setReadyMutation = useSetReady();

  useRoomEvents(roomId, {
    onGameStarted: (event) => {
      navigate(ROUTES.GAME(event.gameId));
    },
    onRoomClosed: () => {
      navigate(ROUTES.LOBBY);
    },
    onPlayerReady: (event) => {
      // Room data will be invalidated automatically, game will start via SSE if all ready
      if (event.allPlayersReady) {
        // Game should start automatically
      }
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

  const handleKickPlayer = (userId: string) => {
    if (roomId) {
      kickPlayerMutation.mutate({ roomId, userId });
    }
  };

  const handleTransferHost = (newHostId: string) => {
    if (roomId && window.confirm('Вы уверены, что хотите передать роль хоста?')) {
      transferHostMutation.mutate({ roomId, newHostId });
    }
  };

  const handleReady = () => {
    if (!roomId) return;
    const currentPlayer = room?.players.find(p => p.userId === user?.id);
    const newReadyState = !(currentPlayer?.isReady);
    
    setReadyMutation.mutate(
      { roomId, isReady: newReadyState },
      {
        onSuccess: (response) => {
          if (response.gameStarted && response.gameId) {
            navigate(ROUTES.GAME(response.gameId));
          }
        },
      }
    );
  };

  const handleCopyCode = async () => {
    if (!room) return;
    
    try {
      if (navigator.clipboard && navigator.clipboard.writeText) {
        await navigator.clipboard.writeText(room.roomCode);
      } else {
        const textArea = document.createElement('textarea');
        textArea.value = room.roomCode;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
      }
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
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

  const currentPlayer = room.players.find(p => p.userId === user?.id);
  const isCurrentPlayerReady = currentPlayer?.isReady ?? false;
  const readyCount = room.players.filter(p => p.isReady).length;
  const allPlayersReady = readyCount === room.currentPlayers && room.currentPlayers >= 2;

  return (
    <div className="room-page">
      <header className="room-page__header">
        <h1 className="room-page__title">{room.name}</h1>
        <div 
          className={`room-page__code ${copySuccess ? 'room-page__code--copied' : ''}`}
          onClick={handleCopyCode}
          title="Нажмите, чтобы скопировать код"
        >
          <span className="room-page__code-label">{TEXTS.ROOM.ROOM_CODE}</span>
          <span className="room-page__code-value">{room.roomCode}</span>
        </div>
      </header>

      <div className="room-page__content">
        <aside className="room-page__sidebar">
          <Card padding="medium">
            <h2 className="room-page__subtitle">
              {TEXTS.ROOM.PLAYERS(room.currentPlayers, room.maxPlayers)}
            </h2>
            
            {/* Ready Status */}
            <div className="room-page__ready-status">
              <p className="room-page__ready-count">{readyCount} / {room.currentPlayers}</p>
              <p className="room-page__ready-label">игроков готовы</p>
            </div>
            
            <div className="room-page__players">
              {room.players.map((player) => (
                <div 
                  key={player.userId} 
                  className={`room-page__player ${player.isReady ? 'room-page__player--ready' : ''}`}
                >
                  <div className="room-page__player-info">
                    <span className="room-page__player-name">
                      {player.username}
                      {player.role === 'host' && ' 👑'}
                    </span>
                    <span className={`room-page__player-ready-badge ${player.isReady ? 'room-page__player-ready-badge--ready' : 'room-page__player-ready-badge--waiting'}`}>
                      {player.isReady ? 'Готов' : 'Ждёт'}
                    </span>
                  </div>
                  {isHost && player.userId !== user?.id && room.status === 'waiting' && (
                    <div className="room-page__player-actions">
                      <button
                        className="room-page__player-action room-page__player-action--transfer"
                        onClick={() => handleTransferHost(player.userId)}
                        disabled={transferHostMutation.isPending}
                        title="Передать хоста"
                      >
                        👑
                      </button>
                      <button
                        className="room-page__player-action room-page__player-action--kick"
                        onClick={() => handleKickPlayer(player.userId)}
                        disabled={kickPlayerMutation.isPending}
                        title="Выгнать"
                      >
                        ✕
                      </button>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </Card>

          <Card padding="medium" className="room-page__actions">
            {/* Ready Button - for all players */}
            <Button
              variant={isCurrentPlayerReady ? 'secondary' : 'primary'}
              size="large"
              fullWidth
              onClick={handleReady}
              isLoading={setReadyMutation.isPending}
              disabled={room.status !== 'waiting'}
            >
              {isCurrentPlayerReady ? '✓ Вы готовы' : 'Готов!'}
            </Button>
            
            {/* Status message */}
            {room.currentPlayers < 2 ? (
              <p className="room-page__hint">
                {TEXTS.ROOM.MIN_PLAYERS_REQUIRED}
              </p>
            ) : !allPlayersReady ? (
              <p className="room-page__hint">
                Ожидание готовности всех игроков...
              </p>
            ) : (
              <p className="room-page__hint" style={{ color: '#22c55e' }}>
                ✓ Все готовы! Игра начинается...
              </p>
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
