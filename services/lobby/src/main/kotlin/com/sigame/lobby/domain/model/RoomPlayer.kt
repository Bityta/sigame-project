package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.PlayerRole
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("room_players")
data class RoomPlayer(
    @Id
    val id: UUID? = null,
    val roomId: UUID,
    val userId: UUID,
    val username: String,
    val avatarUrl: String? = null,
    val role: String = "player",
    val joinedAt: LocalDateTime = LocalDateTime.now(),
    val leftAt: LocalDateTime? = null
) {
        fun getRoleEnum(): PlayerRole = PlayerRole.valueOf(role.uppercase())
    
    companion object {
                fun roleFromEnum(role: PlayerRole): String = role.name.lowercase()
    }
}

