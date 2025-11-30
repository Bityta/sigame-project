/**
 * Room & Lobby Domain Types
 */

export type RoomStatus = 'waiting' | 'starting' | 'playing' | 'finished' | 'cancelled';
export type PlayerRole = 'host' | 'player' | 'spectator';

export interface RoomSettings {
  timeForAnswer: number;
  timeForChoice: number;
  allowWrongAnswer: boolean;
  showRightAnswer: boolean;
}

export interface RoomPlayer {
  id?: string;
  userId: string;
  username: string;
  role: PlayerRole;
  joinedAt: string;
}

export interface GameRoom {
  id: string;
  roomCode: string;
  hostId: string;
  packId: string;
  name: string;
  status: RoomStatus;
  maxPlayers: number;
  currentPlayers: number;
  isPublic: boolean;
  hasPassword: boolean;
  settings?: RoomSettings;
  players: RoomPlayer[];
  createdAt: string;
  startedAt?: string;
  finishedAt?: string;
}

export interface CreateRoomRequest {
  name: string;
  packId: string;
  maxPlayers: number;
  isPublic: boolean;
  password?: string;
  settings?: Partial<RoomSettings>;
}

export interface JoinRoomRequest {
  role?: PlayerRole;
  password?: string;
}

export interface RoomListQuery {
  page?: number;
  size?: number;
  status?: RoomStatus;
  has_slots?: boolean;
}

export interface StartGameResponse {
  gameId: string;
  websocketUrl: string;
}

export interface KickPlayerRequest {
  targetUserId: string;
}

export interface TransferHostRequest {
  newHostId: string;
}

