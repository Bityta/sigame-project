/**
 * Room Entity - API
 * API методы для работы с комнатами
 */

import { lobbyApi } from '@/shared/api';
import { API_CONFIG } from '@/shared/config';
import type {
  GameRoom,
  CreateRoomRequest,
  JoinRoomRequest,
  RoomListQuery,
  RoomSettings,
  StartGameResponse,
  KickPlayerRequest,
  TransferHostRequest,
} from '@/shared/types';

export const roomApi = {
  /**
   * Получить список комнат
   */
  async getRooms(query?: RoomListQuery): Promise<GameRoom[]> {
    const response = await lobbyApi.get<{ rooms: GameRoom[] }>(
      API_CONFIG.ENDPOINTS.LOBBY.ROOMS,
      { params: query }
    );
    return response.data.rooms;
  },

  /**
   * Получить комнату по ID
   */
  async getRoomById(id: string): Promise<GameRoom> {
    const response = await lobbyApi.get<GameRoom>(
      API_CONFIG.ENDPOINTS.LOBBY.ROOM_BY_ID(id)
    );
    return response.data;
  },

  /**
   * Получить комнату по коду
   */
  async getRoomByCode(code: string): Promise<GameRoom> {
    const response = await lobbyApi.get<GameRoom>(
      API_CONFIG.ENDPOINTS.LOBBY.ROOM_BY_CODE(code)
    );
    return response.data;
  },

  /**
   * Создать комнату
   */
  async createRoom(data: CreateRoomRequest): Promise<GameRoom> {
    const response = await lobbyApi.post<GameRoom>(
      API_CONFIG.ENDPOINTS.LOBBY.ROOMS,
      data
    );
    return response.data;
  },

  /**
   * Присоединиться к комнате
   */
  async joinRoom(id: string, data: JoinRoomRequest): Promise<GameRoom> {
    const response = await lobbyApi.post<GameRoom>(
      API_CONFIG.ENDPOINTS.LOBBY.JOIN_ROOM(id),
      data
    );
    return response.data;
  },

  /**
   * Покинуть комнату
   */
  async leaveRoom(id: string): Promise<void> {
    await lobbyApi.delete(API_CONFIG.ENDPOINTS.LOBBY.LEAVE_ROOM(id));
  },

  /**
   * Обновить настройки комнаты
   */
  async updateRoomSettings(
    id: string,
    settings: Partial<RoomSettings>
  ): Promise<GameRoom> {
    const response = await lobbyApi.patch<GameRoom>(
      API_CONFIG.ENDPOINTS.LOBBY.UPDATE_SETTINGS(id),
      settings
    );
    return response.data;
  },

  /**
   * Начать игру
   */
  async startGame(id: string): Promise<StartGameResponse> {
    const response = await lobbyApi.post<StartGameResponse>(
      API_CONFIG.ENDPOINTS.LOBBY.START_ROOM(id)
    );
    return response.data;
  },

  /**
   * Выгнать игрока из комнаты
   */
  async kickPlayer(id: string, targetUserId: string): Promise<void> {
    await lobbyApi.post(
      API_CONFIG.ENDPOINTS.LOBBY.KICK_PLAYER(id),
      { targetUserId } as KickPlayerRequest
    );
  },

  /**
   * Передать роль хоста другому игроку
   */
  async transferHost(id: string, newHostId: string): Promise<void> {
    await lobbyApi.post(
      API_CONFIG.ENDPOINTS.LOBBY.TRANSFER_HOST(id),
      { newHostId } as TransferHostRequest
    );
  },
};

