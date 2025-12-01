package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.dto.StartGameResponse
import com.sigame.lobby.domain.dto.UpdateRoomSettingsRequest
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.InsufficientPlayersException
import com.sigame.lobby.domain.exception.InvalidRoomStateException
import com.sigame.lobby.domain.exception.PackInfoNotFoundException
import com.sigame.lobby.domain.exception.PackNotApprovedException
import com.sigame.lobby.domain.exception.PackNotFoundException
import com.sigame.lobby.domain.exception.PackNotOwnedException
import com.sigame.lobby.domain.exception.PlayerAlreadyInRoomException
import com.sigame.lobby.domain.exception.UnauthorizedRoomActionException
import com.sigame.lobby.domain.exception.UserInfoNotFoundException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.grpc.auth.AuthServiceClient
import com.sigame.lobby.grpc.auth.UserInfo
import com.sigame.lobby.grpc.pack.PackInfo
import com.sigame.lobby.grpc.pack.PackServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.RoomCodeGenerator
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.data.PlayerRepository
import com.sigame.lobby.service.data.RoomRepository
import com.sigame.lobby.service.data.SettingsRepository
import com.sigame.lobby.service.external.GameServiceClient
import com.sigame.lobby.service.external.GameSettings
import com.sigame.lobby.sse.event.GameStartedEvent
import com.sigame.lobby.sse.service.RoomEventPublisher
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
import org.springframework.stereotype.Component
import java.time.LocalDateTime
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Component
class RoomLifecycleHelper(
    private val roomRepository: RoomRepository,
    private val playerRepository: PlayerRepository,
    private val settingsRepository: SettingsRepository,
    private val gameRoomRepository: GameRoomRepository,
    private val roomCodeGenerator: RoomCodeGenerator,
    private val roomCacheService: RoomCacheService,
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val gameServiceClient: GameServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val roomEventPublisher: RoomEventPublisher,
    private val passwordEncoder: BCryptPasswordEncoder
) {

    suspend fun validatePackForRoom(packId: UUID, userId: UUID) {
        val validation = packServiceClient.validatePack(packId, userId)
            ?: throw PackNotFoundException(packId)

        if (!validation.exists) throw PackNotFoundException(packId)
        if (!validation.isApproved) throw PackNotApprovedException(packId, validation.status)
        if (!validation.isOwner) throw PackNotOwnedException(packId, userId)
    }

    suspend fun validateUserNotInRoom(userId: UUID) {
        val activePlayer = playerRepository.findActiveByUserId(userId)
        if (activePlayer != null) {
            throw PlayerAlreadyInRoomException(userId, activePlayer.roomId)
        }
    }

    suspend fun findRoomOrThrow(roomId: UUID): GameRoom = roomRepository.findById(roomId)

    fun validateHostAccess(room: GameRoom, userId: UUID, action: String) {
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, action, "Only the host can $action")
        }
    }

    fun validateRoomInStatus(room: GameRoom, expectedStatus: RoomStatus, action: String) {
        if (room.getStatusEnum() != expectedStatus) {
            throw InvalidRoomStateException(room.requireId(), room.getStatusEnum(), action)
        }
    }

    suspend fun getActivePlayersWithMinCheck(roomId: UUID): List<RoomPlayer> {
        val players = playerRepository.findActiveByRoomId(roomId)
        if (players.size < 2) {
            throw InsufficientPlayersException(roomId, players.size)
        }
        return players
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

        return roomRepository.save(room)
    }

    suspend fun createRoomSettings(roomId: UUID, settings: RoomSettingsDto?) {
        val roomSettings = RoomSettings(
            roomId = roomId,
            timeForAnswer = settings?.timeForAnswer ?: 30,
            timeForChoice = settings?.timeForChoice ?: 60,
            allowWrongAnswer = settings?.allowWrongAnswer ?: true,
            showRightAnswer = settings?.showRightAnswer ?: true
        )
        settingsRepository.insert(roomSettings)
    }

    suspend fun createHostPlayer(roomId: UUID, hostId: UUID, username: String, avatarUrl: String?): RoomPlayer {
        val hostPlayer = RoomPlayer(
            roomId = roomId,
            userId = hostId,
            username = username,
            avatarUrl = avatarUrl,
            role = RoomPlayer.roleFromEnum(PlayerRole.HOST)
        )
        return playerRepository.save(hostPlayer)
    }

    suspend fun buildGameSettings(roomId: UUID): GameSettings {
        val settings = settingsRepository.findByRoomId(roomId)
        return GameSettings(
            timeForAnswer = settings?.timeForAnswer ?: 30,
            timeForChoice = settings?.timeForChoice ?: 60,
            allowWrongAnswer = settings?.allowWrongAnswer ?: true,
            showRightAnswer = settings?.showRightAnswer ?: true
        )
    }

    suspend fun updateRoomToStarting(room: GameRoom): GameRoom {
        val updatedRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.STARTING),
            startedAt = LocalDateTime.now()
        )
        val savedRoom = gameRoomRepository.save(updatedRoom).awaitFirst()
        lobbyMetrics.recordRoomStarted()
        return savedRoom
    }

    suspend fun executeGameStart(
        roomId: UUID,
        room: GameRoom,
        savedRoom: GameRoom,
        activePlayers: List<RoomPlayer>,
        gameSettings: GameSettings
    ): StartGameResponse {
        try {
            val gameSession = gameServiceClient.createGameSession(
                roomId = roomId,
                packId = room.packId,
                players = activePlayers,
                settings = gameSettings
            )

            updateRoomToPlaying(savedRoom)

            roomEventPublisher.publish(
                GameStartedEvent(
                    eventRoomId = roomId,
                    gameId = UUID.fromString(gameSession.gameSessionId),
                    websocketUrl = gameSession.wsUrl
                )
            )
            roomEventPublisher.closeRoom(roomId)

            return StartGameResponse(
                gameId = gameSession.gameSessionId,
                websocketUrl = gameSession.wsUrl
            )
        } catch (e: Exception) {
            logger.error(e) { "Failed to create game session, rolling back room status" }
            rollbackRoomToWaiting(savedRoom)
            throw e
        }
    }

    suspend fun cancelRoom(room: GameRoom) {
        val updatedRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()
    }

    suspend fun getActivePlayers(roomId: UUID): List<RoomPlayer> = playerRepository.findActiveByRoomId(roomId)

    suspend fun clearRoomAndPlayersCache(roomId: UUID, roomCode: String, players: List<RoomPlayer>) = coroutineScope {
        launch { roomCacheService.clearRoomCache(roomId, roomCode) }
        players.forEach { player ->
            launch { roomCacheService.deleteUserCurrentRoom(player.userId) }
        }
    }

    fun recordRoomCancelledMetrics(room: GameRoom, playersCount: Int) {
        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        repeat(playersCount) { lobbyMetrics.recordPlayerLeft() }
    }

    suspend fun findSettingsByRoomId(roomId: UUID): RoomSettings? = settingsRepository.findByRoomId(roomId)

    fun mergeSettings(roomId: UUID, current: RoomSettings?, request: UpdateRoomSettingsRequest): RoomSettings =
        current?.copy(
            timeForAnswer = request.timeForAnswer ?: current.timeForAnswer,
            timeForChoice = request.timeForChoice ?: current.timeForChoice,
            allowWrongAnswer = request.allowWrongAnswer ?: current.allowWrongAnswer,
            showRightAnswer = request.showRightAnswer ?: current.showRightAnswer
        ) ?: RoomSettings(
            roomId = roomId,
            timeForAnswer = request.timeForAnswer ?: 30,
            timeForChoice = request.timeForChoice ?: 60,
            allowWrongAnswer = request.allowWrongAnswer ?: true,
            showRightAnswer = request.showRightAnswer ?: true
        )

    suspend fun saveSettings(settings: RoomSettings, isNew: Boolean) {
        if (isNew) {
            settingsRepository.insert(settings)
        } else {
            settingsRepository.save(settings)
        }
    }

    fun toSettingsDto(settings: RoomSettings) = RoomSettingsDto(
        timeForAnswer = settings.timeForAnswer,
        timeForChoice = settings.timeForChoice,
        allowWrongAnswer = settings.allowWrongAnswer,
        showRightAnswer = settings.showRightAnswer
    )

    fun recordRoomCreatedMetrics() {
        lobbyMetrics.recordRoomCreated()
        lobbyMetrics.recordPlayerJoined()
    }

    private suspend fun updateRoomToPlaying(room: GameRoom) {
        val playingRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.PLAYING))
        gameRoomRepository.save(playingRoom).awaitFirstOrNull()
    }

    private suspend fun rollbackRoomToWaiting(room: GameRoom) {
        logger.error { "Rolling back room ${room.id} status to WAITING" }
        val rollbackRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.WAITING),
            startedAt = null
        )
        gameRoomRepository.save(rollbackRoom).awaitFirstOrNull()
    }
}
