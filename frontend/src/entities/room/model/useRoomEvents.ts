import { useEffect, useRef, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { API_CONFIG } from '@/shared/config';
import { roomKeys } from './queries';

export interface RoomEvent {
  type: string;
  roomId: string;
  timestamp: string;
}

export interface PlayerJoinedEvent extends RoomEvent {
  type: 'player_joined';
  userId: string;
  username: string;
  currentPlayers: number;
}

export interface PlayerLeftEvent extends RoomEvent {
  type: 'player_left';
  userId: string;
  username: string;
  reason: string;
  currentPlayers: number;
}

export interface GameStartedEvent extends RoomEvent {
  type: 'game_started';
  gameId: string;
  websocketUrl: string;
}

export interface RoomClosedEvent extends RoomEvent {
  type: 'room_closed';
  reason: string;
}

type AnyRoomEvent = PlayerJoinedEvent | PlayerLeftEvent | GameStartedEvent | RoomClosedEvent;

interface UseRoomEventsOptions {
  onPlayerJoined?: (event: PlayerJoinedEvent) => void;
  onPlayerLeft?: (event: PlayerLeftEvent) => void;
  onGameStarted?: (event: GameStartedEvent) => void;
  onRoomClosed?: (event: RoomClosedEvent) => void;
  onError?: (error: Event) => void;
}

export const useRoomEvents = (roomId: string | undefined, options: UseRoomEventsOptions = {}) => {
  const eventSourceRef = useRef<EventSource | null>(null);
  const queryClient = useQueryClient();

  const invalidateRoom = useCallback(() => {
    if (roomId) {
      queryClient.invalidateQueries({ queryKey: roomKeys.detail(roomId) });
    }
  }, [queryClient, roomId]);

  useEffect(() => {
    if (!roomId) return;

    const url = `${API_CONFIG.LOBBY_BASE_URL}/api/lobby/rooms/${roomId}/events`;
    const eventSource = new EventSource(url);
    eventSourceRef.current = eventSource;

    eventSource.addEventListener('player_joined', (e) => {
      const event = JSON.parse(e.data) as PlayerJoinedEvent;
      invalidateRoom();
      options.onPlayerJoined?.(event);
    });

    eventSource.addEventListener('player_left', (e) => {
      const event = JSON.parse(e.data) as PlayerLeftEvent;
      invalidateRoom();
      options.onPlayerLeft?.(event);
    });

    eventSource.addEventListener('game_started', (e) => {
      const event = JSON.parse(e.data) as GameStartedEvent;
      options.onGameStarted?.(event);
    });

    eventSource.addEventListener('room_closed', (e) => {
      const event = JSON.parse(e.data) as RoomClosedEvent;
      options.onRoomClosed?.(event);
    });

    eventSource.onerror = (error) => {
      console.error('[RoomEvents] SSE error:', error);
      options.onError?.(error);
    };

    return () => {
      eventSource.close();
      eventSourceRef.current = null;
    };
  }, [roomId, invalidateRoom, options]);

  const close = useCallback(() => {
    eventSourceRef.current?.close();
    eventSourceRef.current = null;
  }, []);

  return { close };
};

