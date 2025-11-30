package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.annotation.JsonProperty
import jakarta.validation.Valid
import jakarta.validation.constraints.Max
import jakarta.validation.constraints.Min
import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.NotNull
import jakarta.validation.constraints.Size
import java.util.UUID

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonInclude(JsonInclude.Include.NON_NULL)
data class CreateRoomRequest(
    @field:NotBlank(message = "Room name is required")
    @field:Size(min = 3, max = 100, message = "Room name must be between 3 and 100 characters")
    @JsonProperty("name")
    val name: String,
    
    @field:NotNull(message = "Pack ID is required")
    @JsonProperty("packId")
    val packId: UUID,
    
    @field:Min(value = 2, message = "Minimum 2 players required")
    @field:Max(value = 12, message = "Maximum 12 players allowed")
    @JsonProperty("maxPlayers")
    val maxPlayers: Int = 6,
    
    @JsonProperty("isPublic")
    val isPublic: Boolean = true,
    
    @field:Size(min = 4, max = 50, message = "Password must be between 4 and 50 characters")
    @JsonProperty("password")
    val password: String? = null,
    
    @field:Valid
    @JsonProperty("settings")
    val settings: RoomSettingsDto? = null
)
