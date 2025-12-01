package com.sigame.lobby.service.data

import com.sigame.lobby.domain.model.RoomSettings
import java.util.UUID

interface SettingsRepository {
    suspend fun findByRoomId(roomId: UUID): RoomSettings?
    suspend fun save(settings: RoomSettings): RoomSettings
    suspend fun insert(settings: RoomSettings): RoomSettings
}

