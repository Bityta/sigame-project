package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.RoomStatus
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("game_rooms")
data class GameRoom(
    @Id
    val id: UUID? = null,
    val roomCode: String,
    val hostId: UUID,
    val packId: UUID,
    val name: String,
    val status: String = "waiting",
    val maxPlayers: Int = 6,
    val isPublic: Boolean = true,
    val passwordHash: String? = null,
    val createdAt: LocalDateTime = LocalDateTime.now(),
    val updatedAt: LocalDateTime = LocalDateTime.now(),
    val startedAt: LocalDateTime? = null,
    val finishedAt: LocalDateTime? = null
) {
        fun getStatusEnum(): RoomStatus = RoomStatus.valueOf(status.uppercase())
    
    companion object {
                fun statusFromEnum(status: RoomStatus): String = status.name.lowercase()
    }
}

