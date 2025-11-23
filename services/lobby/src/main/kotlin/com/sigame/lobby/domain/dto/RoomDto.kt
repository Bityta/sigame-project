package com.sigame.lobby.domain.dto

import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.model.RoomSettings
import java.time.LocalDateTime
import java.util.UUID

/**
 * DTO комнаты для API ответов
 */
data class RoomDto(
    val id: UUID,
    val roomCode: String,
    val name: String,
    val hostId: UUID,
    val hostUsername: String? = null,
    val packId: UUID,
    val packName: String? = null,
    val status: String,
    val maxPlayers: Int,
    val currentPlayers: Int,
    val isPublic: Boolean,
    val hasPassword: Boolean,
    val createdAt: LocalDateTime,
    val players: List<PlayerDto>? = null,
    val settings: RoomSettings? = null
) {
    /**
     * Получить статус комнаты как enum
     */
    fun getStatusEnum(): RoomStatus = RoomStatus.valueOf(status.uppercase())
}

