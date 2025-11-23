package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.PlayerRole
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table
import java.time.LocalDateTime
import java.util.UUID

/**
 * Сущность игрока в комнате
 */
@Table("room_players")
data class RoomPlayer(
    @Id
    val id: UUID? = null,
    val roomId: UUID,
    val userId: UUID,
    val role: String = "player",
    val joinedAt: LocalDateTime = LocalDateTime.now(),
    val leftAt: LocalDateTime? = null
) {
    /**
     * Получить роль игрока как enum
     */
    fun getRoleEnum(): PlayerRole = PlayerRole.valueOf(role.uppercase())
    
    companion object {
        /**
         * Преобразовать enum роли в строку для БД
         */
        fun roleFromEnum(role: PlayerRole): String = role.name.lowercase()
    }
}

