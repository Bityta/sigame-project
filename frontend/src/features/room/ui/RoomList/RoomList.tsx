/**
 * Room Feature - RoomList
 * Список комнат с активной комнатой первой
 */

import { useNavigate } from 'react-router-dom';
import { useRooms } from '@/entities/room';
import { Card, Button, Spinner } from '@/shared/ui';
import { ROUTES } from '@/shared/config';
import type { GameRoom } from '@/shared/types';
import './RoomList.css';

interface RoomListProps {
  activeRoom?: GameRoom | null;
  onLeaveRoom?: () => void;
  isLeavingRoom?: boolean;
}

export const RoomList = ({ activeRoom, onLeaveRoom, isLeavingRoom = false }: RoomListProps) => {
  const navigate = useNavigate();
  const hasActiveRoom = !!activeRoom;
  
  const { data: rooms, isLoading, refetch } = useRooms(
    { status: 'waiting', has_slots: true },
    { refetchInterval: 60000 }
  );

  // Фильтруем активную комнату из списка доступных
  const availableRooms = rooms?.filter(room => room.id !== activeRoom?.id) || [];

  const getRoomStatusText = (room: GameRoom): string => {
    switch (room.status) {
      case 'waiting':
        return 'Ожидание игроков';
      case 'starting':
        return 'Начинается...';
      case 'playing':
        return 'Играют';
      case 'finished':
        return 'Завершена';
      default:
        return room.status;
    }
  };

  const handleJoinRoom = (roomId: string) => {
    if (hasActiveRoom) return;
    navigate(ROUTES.ROOM(roomId));
  };

  const handleReturnToRoom = () => {
    if (activeRoom) {
      navigate(ROUTES.ROOM(activeRoom.id));
    }
  };

  if (isLoading) {
    return <Spinner center size="large" />;
  }

  const hasNoRooms = !activeRoom && availableRooms.length === 0;

  if (hasNoRooms) {
    return (
      <Card className="room-list-empty">
        <p className="room-list-empty__text">Нет доступных комнат</p>
        <Button onClick={() => refetch()}>Обновить</Button>
      </Card>
    );
  }

  return (
    <div className="room-list">
      <div className="room-list__header">
        <h2 className="room-list__title">Комнаты</h2>
        <Button size="small" onClick={() => refetch()}>
          Обновить
        </Button>
      </div>

      <div className="room-list__grid">
        {/* Активная комната — первая в списке */}
        {activeRoom && (
          <Card className="room-card room-card--active" padding="medium">
            <div className="room-card__header">
              <span className="room-card__active-badge">⚠️ Активная</span>
              <h3 className="room-card__name">{activeRoom.name}</h3>
              {activeRoom.hasPassword && (
                <span className="room-card__badge">🔒</span>
              )}
            </div>

            <div className="room-card__info">
              <div className="room-card__info-row">
                <span className="room-card__label">Код:</span>
                <span className="room-card__value room-card__value--code">{activeRoom.roomCode}</span>
              </div>
              <div className="room-card__info-row">
                <span className="room-card__label">Игроки:</span>
                <span className="room-card__value">
                  {activeRoom.currentPlayers}/{activeRoom.maxPlayers}
                </span>
              </div>
              <div className="room-card__info-row">
                <span className="room-card__label">Статус:</span>
                <span className="room-card__value">
                  {getRoomStatusText(activeRoom)}
                </span>
              </div>
            </div>

            <Button
              fullWidth
              variant="primary"
              onClick={handleReturnToRoom}
            >
              Вернуться
            </Button>
            
            <button
              className="room-card__leave-link"
              onClick={onLeaveRoom}
              disabled={isLeavingRoom}
            >
              {isLeavingRoom ? 'Выход...' : 'Покинуть комнату'}
            </button>
          </Card>
        )}

        {/* Доступные комнаты */}
        {availableRooms.map((room) => (
          <Card key={room.id} className="room-card" padding="medium">
            <div className="room-card__header">
              <h3 className="room-card__name">{room.name}</h3>
              {room.hasPassword && (
                <span className="room-card__badge">🔒</span>
              )}
            </div>

            <div className="room-card__info">
              <div className="room-card__info-row">
                <span className="room-card__label">Код:</span>
                <span className="room-card__value">{room.roomCode}</span>
              </div>
              <div className="room-card__info-row">
                <span className="room-card__label">Игроки:</span>
                <span className="room-card__value">
                  {room.currentPlayers}/{room.maxPlayers}
                </span>
              </div>
              <div className="room-card__info-row">
                <span className="room-card__label">Статус:</span>
                <span className="room-card__value">
                  {getRoomStatusText(room)}
                </span>
              </div>
            </div>

            <Button
              fullWidth
              variant="primary"
              onClick={() => handleJoinRoom(room.id)}
              disabled={room.currentPlayers >= room.maxPlayers || hasActiveRoom}
              title={hasActiveRoom ? 'Сначала покиньте текущую комнату' : undefined}
            >
              {room.hasPassword ? 'Войти с паролем' : 'Присоединиться'}
            </Button>
          </Card>
        ))}
      </div>
    </div>
  );
};
