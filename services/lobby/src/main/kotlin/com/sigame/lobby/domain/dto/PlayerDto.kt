package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonFormat
import com.fasterxml.jackson.annotation.JsonProperty
import com.sigame.lobby.domain.enums.PlayerRole
import java.time.LocalDateTime
import java.util.UUID

/**
 * DTO игрока для API ответов
 * Все поля обязательны
 */
data class PlayerDto(
    @JsonProperty("userId")
    val userId: UUID,
    
    @JsonProperty("username")
    val username: String,
    
    @JsonProperty("role")
    val role: String,
    
    @JsonProperty("joinedAt")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    val joinedAt: LocalDateTime
) {
    /**
     * Получить роль игрока как enum
     */
    fun getRoleEnum(): PlayerRole = PlayerRole.valueOf(role.uppercase())
}

