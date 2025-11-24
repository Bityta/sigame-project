package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.Valid
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.Size

/**
 * Запрос на обновление настроек комнаты
 * Все поля необязательны - обновляются только переданные
 */
@JsonIgnoreProperties(ignoreUnknown = true)
@JsonInclude(JsonInclude.Include.NON_NULL)
data class UpdateRoomSettingsRequest(
    @field:Min(value = 2, message = "Minimum 2 players required")
    @field:Max(value = 12, message = "Maximum 12 players allowed")
    @JsonProperty("maxPlayers")
    val maxPlayers: Int? = null,
    
    @JsonProperty("isPublic")
    val isPublic: Boolean? = null,
    
    @field:Size(min = 4, max = 50, message = "Password must be between 4 and 50 characters")
    @JsonProperty("password")
    val password: String? = null,
    
    @field:Valid
    @JsonProperty("settings")
    val settings: RoomSettingsDto? = null
)

