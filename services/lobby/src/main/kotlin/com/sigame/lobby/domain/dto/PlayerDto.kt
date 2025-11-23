package com.sigame.lobby.domain.dto

import com.sigame.lobby.domain.enums.PlayerRole
import java.time.LocalDateTime
import java.util.UUID

/**
 * DTO игрока для API ответов
 */
data class PlayerDto(
    val userId: UUID,
    val username: String,
    val role: String,
    val joinedAt: LocalDateTime
) {
    /**
     * Получить роль игрока как enum
     */
    fun getRoleEnum(): PlayerRole = PlayerRole.valueOf(role.uppercase())
}

