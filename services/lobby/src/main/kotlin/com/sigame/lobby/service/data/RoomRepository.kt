package com.sigame.lobby.service.data

import com.sigame.lobby.domain.model.GameRoom
import java.util.UUID

interface RoomRepository {
    suspend fun findById(id: UUID): GameRoom
    suspend fun findByCode(code: String): GameRoom
    suspend fun save(room: GameRoom): GameRoom
}

