package com.sigame.lobby.domain.dto

import jakarta.validation.Valid
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.NotNull
import jakarta.validation.constraints.Size
import java.util.UUID

/**
 * Запрос на создание комнаты
 */
data class CreateRoomRequest(
    @field:NotBlank(message = "Room name is required")
    @field:Size(min = 3, max = 100, message = "Room name must be between 3 and 100 characters")
    val name: String,
    
    @field:NotNull(message = "Pack ID is required")
    val packId: UUID,
    
    @field:Min(value = 2, message = "Minimum 2 players required")
    @field:Max(value = 12, message = "Maximum 12 players allowed")
    val maxPlayers: Int = 6,
    
    val isPublic: Boolean = true,
    
    @field:Size(min = 4, max = 50, message = "Password must be between 4 and 50 characters")
    val password: String? = null,
    
    @field:Valid
    val settings: RoomSettingsDto? = null
)

