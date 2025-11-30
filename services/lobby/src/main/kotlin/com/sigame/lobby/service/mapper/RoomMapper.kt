package com.sigame.lobby.service.mapper

import com.sigame.lobby.domain.dto.PlayerDto
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.PackInfo
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.grpc.UserInfo
import com.sigame.lobby.service.batch.BatchOperationService
import org.springframework.stereotype.Component
import java.util.UUID

/**
 * Mapper для преобразования доменных объектов в DTO
 */
@Component
class RoomMapper(
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val batchOperationService: BatchOperationService
) {
    
    /**
     * Преобразует GameRoom в RoomDto с дополнительной информацией
     */
    suspend fun toDto(
        room: GameRoom,
        currentPlayers: Int,
        players: List<RoomPlayer>? = null,
        settings: RoomSettings? = null
    ): RoomDto {
        val hostInfo = authServiceClient.getUserInfo(room.hostId)
        val packInfo = packServiceClient.getPackInfo(room.packId)
        
        val playerDtos = players?.let { buildPlayerDtos(it) }
        val settingsDto = settings?.let { toSettingsDto(it) }
        
        return RoomDto(
            id = room.id!!,
            roomCode = room.roomCode,
            name = room.name,
            hostId = room.hostId,
            hostUsername = hostInfo?.username,
            packId = room.packId,
            packName = packInfo?.name,
            status = room.status,
            maxPlayers = room.maxPlayers,
            currentPlayers = currentPlayers,
            isPublic = room.isPublic,
            hasPassword = room.passwordHash != null,
            createdAt = room.createdAt,
            players = playerDtos,
            settings = settingsDto
        )
    }
    
    /**
     * Преобразует GameRoom в RoomDto используя кэш данных (для batch операций)
     */
    suspend fun toDtoWithCache(
        room: GameRoom,
        currentPlayers: Int,
        userInfoCache: Map<UUID, UserInfo?> = emptyMap(),
        packInfoCache: Map<UUID, PackInfo?> = emptyMap(),
        players: List<RoomPlayer>? = null,
        settings: RoomSettings? = null
    ): RoomDto {
        val hostInfo = userInfoCache[room.hostId]
        val packInfo = packInfoCache[room.packId]
        
        val playerDtos = players?.let { buildPlayerDtosWithCache(it, userInfoCache) }
        val settingsDto = settings?.let { toSettingsDto(it) }
        
        return RoomDto(
            id = room.id!!,
            roomCode = room.roomCode,
            name = room.name,
            hostId = room.hostId,
            hostUsername = hostInfo?.username,
            packId = room.packId,
            packName = packInfo?.name,
            status = room.status,
            maxPlayers = room.maxPlayers,
            currentPlayers = currentPlayers,
            isPublic = room.isPublic,
            hasPassword = room.passwordHash != null,
            createdAt = room.createdAt,
            players = playerDtos,
            settings = settingsDto
        )
    }
    
    /**
     * Создает список PlayerDto из RoomPlayer
     */
    private suspend fun buildPlayerDtos(players: List<RoomPlayer>): List<PlayerDto> {
        // Используем batch операцию для оптимизации
        val userIds = players.map { it.userId }
        val userInfoMap = batchOperationService.getUserInfoBatch(userIds)
        
        return players.map { player ->
            val userInfo = userInfoMap[player.userId]
            PlayerDto(
                userId = player.userId,
                username = userInfo?.username ?: "Unknown",
                avatarUrl = userInfo?.avatarUrl,
                role = player.role
            )
        }
    }
    
    /**
     * Создает список PlayerDto используя кэш
     */
    private fun buildPlayerDtosWithCache(
        players: List<RoomPlayer>,
        userInfoCache: Map<UUID, UserInfo?>
    ): List<PlayerDto> {
        return players.map { player ->
            val userInfo = userInfoCache[player.userId]
            PlayerDto(
                userId = player.userId,
                username = userInfo?.username ?: "Unknown",
                avatarUrl = userInfo?.avatarUrl,
                role = player.role
            )
        }
    }
    
    /**
     * Преобразует RoomSettings (Entity) в RoomSettingsDto
     */
    private fun toSettingsDto(settings: RoomSettings): RoomSettingsDto {
        return RoomSettingsDto(
            timeForAnswer = settings.timeForAnswer,
            timeForChoice = settings.timeForChoice,
            allowWrongAnswer = settings.allowWrongAnswer,
            showRightAnswer = settings.showRightAnswer
        )
    }
}

