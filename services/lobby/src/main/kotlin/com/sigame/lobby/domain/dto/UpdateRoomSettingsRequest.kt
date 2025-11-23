package com.sigame.lobby.domain.dto

import jakarta.validation.Valid
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.Size

/**
 * Запрос на обновление настроек комнаты
 */
data class UpdateRoomSettingsRequest(
    @field:Min(value = 2, message = "Minimum 2 players required")
    @field:Max(value = 12, message = "Maximum 12 players allowed")
    val maxPlayers: Int? = null,
    
    val isPublic: Boolean? = null,
    
    @field:Size(min = 4, max = 50, message = "Password must be between 4 and 50 characters")
    val password: String? = null,
    
    @field:Valid
    val settings: RoomSettingsDto? = null
)

