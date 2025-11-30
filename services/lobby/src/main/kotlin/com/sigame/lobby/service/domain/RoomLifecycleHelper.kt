package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.exception.PackInfoNotFoundException
import com.sigame.lobby.domain.exception.PackNotApprovedException
import com.sigame.lobby.domain.exception.PackNotFoundException
import com.sigame.lobby.domain.exception.PackNotOwnedException
import com.sigame.lobby.domain.exception.PlayerAlreadyInRoomException
import com.sigame.lobby.domain.exception.UserInfoNotFoundException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.PackInfo
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.grpc.UserInfo
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
import com.sigame.lobby.service.RoomCodeGenerator
import com.sigame.lobby.service.cache.RoomCacheService
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
import org.springframework.stereotype.Component
import java.util.UUID

@Component
class RoomLifecycleHelper(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomCodeGenerator: RoomCodeGenerator,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val passwordEncoder: BCryptPasswordEncoder = BCryptPasswordEncoder(12)
) {

    suspend fun validatePackForRoom(packId: UUID, userId: UUID) {
        val validation = packServiceClient.validatePack(packId, userId)
            ?: throw PackNotFoundException(packId)

        if (!validation.exists) throw PackNotFoundException(packId)
        if (!validation.isApproved) throw PackNotApprovedException(packId, validation.status)
        if (!validation.isOwner) throw PackNotOwnedException(packId, userId)
    }

    suspend fun validateUserNotInRoom(userId: UUID) {
        val existingPlayer = roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()
        if (existingPlayer != null) {
            throw PlayerAlreadyInRoomException(userId, existingPlayer.roomId)
        }
    }

    suspend fun fetchRequiredUserInfo(userId: UUID): UserInfo =
        authServiceClient.getUserInfo(userId)
            ?: throw UserInfoNotFoundException(userId)

    suspend fun fetchRequiredPackInfo(packId: UUID): PackInfo =
        packServiceClient.getPackInfo(packId)
            ?: throw PackInfoNotFoundException(packId)

    suspend fun createGameRoom(hostId: UUID, request: CreateRoomRequest): GameRoom {
        val roomCode = roomCodeGenerator.generateUniqueCode()
        val passwordHash = request.password?.let { passwordEncoder.encode(it) }

        val room = GameRoom(
            roomCode = roomCode,
            hostId = hostId,
            packId = request.packId,
            name = request.name,
            maxPlayers = request.maxPlayers,
            isPublic = request.isPublic,
            passwordHash = passwordHash
        )

        return gameRoomRepository.save(room).awaitFirst()
    }

    suspend fun createRoomSettings(roomId: UUID, settings: RoomSettingsDto?) {
        roomSettingsRepository.insertRoomSettings(
            roomId = roomId,
            timeForAnswer = settings?.timeForAnswer ?: 30,
            timeForChoice = settings?.timeForChoice ?: 60,
            allowWrongAnswer = settings?.allowWrongAnswer ?: true,
            showRightAnswer = settings?.showRightAnswer ?: true
        ).awaitFirstOrNull()
    }

    suspend fun createHostPlayer(
        roomId: UUID,
        hostId: UUID,
        username: String,
        avatarUrl: String?
    ): RoomPlayer {
        val hostPlayer = RoomPlayer(
            roomId = roomId,
            userId = hostId,
            username = username,
            avatarUrl = avatarUrl,
            role = RoomPlayer.roleFromEnum(PlayerRole.HOST)
        )
        return roomPlayerRepository.save(hostPlayer).awaitFirst()
    }

    suspend fun updateCacheForNewRoom(room: GameRoom, hostId: UUID) = coroutineScope {
        launch { roomCacheService.cacheRoomData(room, 1) }
        launch { roomCacheService.setUserCurrentRoom(hostId, room.id) }
        launch { roomCacheService.addRoomPlayer(room.id, hostId) }
    }

    suspend fun publishRoomCreatedEvent(
        room: GameRoom,
        hostId: UUID,
        hostInfo: UserInfo,
        packName: String,
        request: CreateRoomRequest
    ) {
        kafkaEventPublisher.publishRoomCreated(
            roomId = room.id,
            roomCode = room.roomCode,
            hostId = hostId,
            hostUsername = hostInfo.username,
            packId = request.packId,
            packName = packName,
            maxPlayers = request.maxPlayers,
            isPublic = request.isPublic
        )
    }

    fun recordRoomCreatedMetrics() {
        lobbyMetrics.recordRoomCreated()
        lobbyMetrics.recordPlayerJoined()
    }
}

