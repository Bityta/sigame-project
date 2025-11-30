package com.sigame.lobby.domain.dto

import com.fasterxml.jackson.annotation.JsonFormat
import com.fasterxml.jackson.annotation.JsonInclude
import com.fasterxml.jackson.annotation.JsonProperty
import com.sigame.lobby.domain.enums.RoomStatus
import java.time.LocalDateTime
import java.util.UUID

@JsonInclude(JsonInclude.Include.NON_NULL)
data class RoomDto(
    @JsonProperty("id")
    val id: UUID,
    
    @JsonProperty("roomCode")
    val roomCode: String,
    
    @JsonProperty("name")
    val name: String,
    
    @JsonProperty("hostId")
    val hostId: UUID,
    
    @JsonProperty("hostUsername")
    val hostUsername: String? = null,
    
    @JsonProperty("packId")
    val packId: UUID,
    
    @JsonProperty("packName")
    val packName: String? = null,
    
    @JsonProperty("status")
    val status: String,
    
    @JsonProperty("maxPlayers")
    val maxPlayers: Int,
    
    @JsonProperty("currentPlayers")
    val currentPlayers: Int,
    
    @JsonProperty("isPublic")
    val isPublic: Boolean,
    
    @JsonProperty("hasPassword")
    val hasPassword: Boolean,
    
    @JsonProperty("createdAt")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    val createdAt: LocalDateTime,
    
    @JsonProperty("players")
    val players: List<PlayerDto>? = null,
    
    @JsonProperty("settings")
    val settings: RoomSettingsDto? = null
)
