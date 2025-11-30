package com.sigame.lobby.service.mapper

import com.sigame.lobby.domain.dto.PlayerDto
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import org.springframework.stereotype.Component

@Component
class RoomMapper {

    fun toDto(
        room: GameRoom,
        currentPlayers: Int,
        players: List<RoomPlayer> = emptyList(),
        settings: RoomSettings? = null,
        packName: String? = null
    ): RoomDto {
        val hostPlayer = players.find { it.userId == room.hostId }
        
        val playerDtos = players.map { player ->
            PlayerDto(
                userId = player.userId,
                username = player.username,
                avatarUrl = player.avatarUrl,
                role = player.role
            )
        }

        return RoomDto(
            id = room.requireId(),
            roomCode = room.roomCode,
            name = room.name,
            hostId = room.hostId,
            hostUsername = hostPlayer?.username,
            packId = room.packId,
            packName = packName,
            status = room.status,
            maxPlayers = room.maxPlayers,
            currentPlayers = currentPlayers,
            isPublic = room.isPublic,
            hasPassword = room.passwordHash != null,
            createdAt = room.createdAt,
            players = playerDtos,
            settings = settings?.let { toSettingsDto(it) }
        )
    }

    fun toDtoWithCache(
        room: GameRoom,
        currentPlayers: Int,
        packName: String? = null,
        players: List<RoomPlayer>? = null,
        settings: RoomSettings? = null
    ): RoomDto {
        val hostPlayer = players?.find { it.userId == room.hostId }
        
        val playerDtos = players?.map { player ->
            PlayerDto(
                userId = player.userId,
                username = player.username,
                avatarUrl = player.avatarUrl,
                role = player.role
            )
        }

        return RoomDto(
            id = room.requireId(),
            roomCode = room.roomCode,
            name = room.name,
            hostId = room.hostId,
            hostUsername = hostPlayer?.username,
            packId = room.packId,
            packName = packName,
            status = room.status,
            maxPlayers = room.maxPlayers,
            currentPlayers = currentPlayers,
            isPublic = room.isPublic,
            hasPassword = room.passwordHash != null,
            createdAt = room.createdAt,
            players = playerDtos,
            settings = settings?.let { toSettingsDto(it) }
        )
    }

    private fun toSettingsDto(settings: RoomSettings): RoomSettingsDto {
        return RoomSettingsDto(
            timeForAnswer = settings.timeForAnswer,
            timeForChoice = settings.timeForChoice,
            allowWrongAnswer = settings.allowWrongAnswer,
            showRightAnswer = settings.showRightAnswer
        )
    }
}
