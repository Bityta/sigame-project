/**
 * Room Entity - Queries
 * React Query хуки для комнат
 */

import {
  useQuery,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query';
import type {
  UseQueryOptions,
  UseMutationOptions,
} from '@tanstack/react-query';
import { roomApi } from '../api/roomApi';
import type {
  GameRoom,
  CreateRoomRequest,
  JoinRoomRequest,
  RoomListQuery,
  RoomSettings,
  StartGameResponse,
} from '@/shared/types';

// Query Keys
export const roomKeys = {
  all: ['room'] as const,
  lists: () => [...roomKeys.all, 'list'] as const,
  list: (filters?: RoomListQuery) => [...roomKeys.lists(), filters] as const,
  details: () => [...roomKeys.all, 'detail'] as const,
  detail: (id: string) => [...roomKeys.details(), id] as const,
  byCode: (code: string) => [...roomKeys.all, 'code', code] as const,
  myActive: () => [...roomKeys.all, 'my-active'] as const,
};

/**
 * Хук для получения списка комнат
 */
export const useRooms = (
  filters?: RoomListQuery,
  options?: Omit<UseQueryOptions<GameRoom[], Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<GameRoom[], Error>({
    queryKey: roomKeys.list(filters),
    queryFn: () => roomApi.getRooms(filters),
    ...options,
  });
};

/**
 * Хук для получения комнаты по ID
 */
export const useRoom = (
  id: string,
  options?: Omit<UseQueryOptions<GameRoom, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<GameRoom, Error>({
    queryKey: roomKeys.detail(id),
    queryFn: () => roomApi.getRoomById(id),
    enabled: !!id,
    ...options,
  });
};

/**
 * Хук для получения комнаты по коду
 */
export const useRoomByCode = (
  code: string,
  options?: Omit<UseQueryOptions<GameRoom, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<GameRoom, Error>({
    queryKey: roomKeys.byCode(code),
    queryFn: () => roomApi.getRoomByCode(code),
    enabled: !!code,
    ...options,
  });
};

/**
 * Хук для получения активной комнаты текущего пользователя
 */
export const useMyActiveRoom = (
  options?: Omit<UseQueryOptions<GameRoom | null, Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery<GameRoom | null, Error>({
    queryKey: roomKeys.myActive(),
    queryFn: () => roomApi.getMyActiveRoom(),
    ...options,
  });
};

/**
 * Мутация: Создание комнаты
 */
export const useCreateRoom = (
  options?: UseMutationOptions<GameRoom, Error, CreateRoomRequest>
) => {
  const queryClient = useQueryClient();

  return useMutation<GameRoom, Error, CreateRoomRequest>({
    mutationFn: roomApi.createRoom,
    onSuccess: (data) => {
      // Инвалидируем список комнат
      queryClient.invalidateQueries({ queryKey: roomKeys.lists() });
      // Добавляем новую комнату в кеш
      queryClient.setQueryData(roomKeys.detail(data.id), data);
    },
    ...options,
  });
};

/**
 * Мутация: Присоединение к комнате
 */
export const useJoinRoom = (
  options?: UseMutationOptions<
    GameRoom,
    Error,
    { id: string; data: JoinRoomRequest }
  >
) => {
  const queryClient = useQueryClient();

  return useMutation<GameRoom, Error, { id: string; data: JoinRoomRequest }>({
    mutationFn: ({ id, data }) => roomApi.joinRoom(id, data),
    onSuccess: (data) => {
      // Обновляем данные комнаты
      queryClient.setQueryData(roomKeys.detail(data.id), data);
      // Инвалидируем список комнат и активную комнату
      queryClient.invalidateQueries({ queryKey: roomKeys.lists() });
      queryClient.invalidateQueries({ queryKey: roomKeys.myActive() });
    },
    ...options,
  });
};

/**
 * Мутация: Выход из комнаты
 */
export const useLeaveRoom = (
  options?: UseMutationOptions<void, Error, string>
) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: roomApi.leaveRoom,
    onSuccess: (_, id) => {
      // Удаляем данные комнаты из кеша
      queryClient.removeQueries({ queryKey: roomKeys.detail(id) });
      // Инвалидируем список комнат и активную комнату
      queryClient.invalidateQueries({ queryKey: roomKeys.lists() });
      queryClient.invalidateQueries({ queryKey: roomKeys.myActive() });
    },
    ...options,
  });
};

/**
 * Мутация: Обновление настроек комнаты
 */
export const useUpdateRoomSettings = (
  options?: UseMutationOptions<
    GameRoom,
    Error,
    { id: string; settings: Partial<RoomSettings> }
  >
) => {
  const queryClient = useQueryClient();

  return useMutation<
    GameRoom,
    Error,
    { id: string; settings: Partial<RoomSettings> }
  >({
    mutationFn: ({ id, settings }) => roomApi.updateRoomSettings(id, settings),
    onSuccess: (data) => {
      // Обновляем данные комнаты
      queryClient.setQueryData(roomKeys.detail(data.id), data);
    },
    ...options,
  });
};

/**
 * Мутация: Начало игры
 */
export const useStartGame = (
  options?: UseMutationOptions<StartGameResponse, Error, string>
) => {
  const queryClient = useQueryClient();

  return useMutation<StartGameResponse, Error, string>({
    mutationFn: roomApi.startGame,
    onSuccess: (_, roomId) => {
      // Инвалидируем данные комнаты
      queryClient.invalidateQueries({ queryKey: roomKeys.detail(roomId) });
    },
    ...options,
  });
};

/**
 * Мутация: Выгнать игрока из комнаты
 */
export const useKickPlayer = (
  options?: UseMutationOptions<void, Error, { roomId: string; userId: string }>
) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, { roomId: string; userId: string }>({
    mutationFn: ({ roomId, userId }) => roomApi.kickPlayer(roomId, userId),
    onSuccess: (_, { roomId }) => {
      queryClient.invalidateQueries({ queryKey: roomKeys.detail(roomId) });
    },
    ...options,
  });
};

/**
 * Мутация: Передать роль хоста
 */
export const useTransferHost = (
  options?: UseMutationOptions<void, Error, { roomId: string; newHostId: string }>
) => {
  const queryClient = useQueryClient();

  return useMutation<void, Error, { roomId: string; newHostId: string }>({
    mutationFn: ({ roomId, newHostId }) => roomApi.transferHost(roomId, newHostId),
    onSuccess: (_, { roomId }) => {
      // Инвалидируем данные комнаты
      queryClient.invalidateQueries({ queryKey: roomKeys.detail(roomId) });
    },
    ...options,
  });
};

