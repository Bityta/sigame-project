package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.PlayerRole
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Column
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

@Table("room_players")
data class RoomPlayer(
    @Id
    val id: UUID? = null,
    @Column("room_id")
    val roomId: UUID,
    @Column("user_id")
    val userId: UUID,
    @Column("username")
    val username: String,
    @Column("avatar_url")
    val avatarUrl: String? = null,
    @Column("role")
    val role: String = "player",
    @Column("joined_at")
    val joinedAt: LocalDateTime = LocalDateTime.now(),
    @Column("left_at")
    val leftAt: LocalDateTime? = null
) {
    fun getRoleEnum(): PlayerRole = PlayerRole.valueOf(role.uppercase())
    
    companion object {
        fun roleFromEnum(role: PlayerRole): String = role.name.lowercase()
    }
}

