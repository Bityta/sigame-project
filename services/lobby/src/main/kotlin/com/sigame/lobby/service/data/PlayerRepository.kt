package com.sigame.lobby.service.data

import com.sigame.lobby.domain.model.RoomPlayer
import java.util.UUID

interface PlayerRepository {
    suspend fun findActiveByUserId(userId: UUID): RoomPlayer?
    suspend fun findByRoomAndUser(roomId: UUID, userId: UUID): RoomPlayer?
    suspend fun findActiveByRoomId(roomId: UUID): List<RoomPlayer>
    suspend fun findActiveByRoomIds(roomIds: List<UUID>): Map<UUID, List<RoomPlayer>>
    suspend fun countActiveByRoomId(roomId: UUID): Int
    suspend fun save(player: RoomPlayer): RoomPlayer
}

