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
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.UnauthorizedRoomActionException
import com.sigame.lobby.domain.exception.UserInfoNotFoundException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.PackInfo
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.grpc.UserInfo
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
import com.sigame.lobby.service.PlayerEventData
import com.sigame.lobby.service.RoomCodeGenerator
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.external.GameServiceClient
import com.sigame.lobby.service.external.GameSessionResponse
import com.sigame.lobby.service.external.GameSettings
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.asFlow
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
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomCodeGenerator: RoomCodeGenerator,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val gameServiceClient: GameServiceClient,
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

    suspend fun findRoomOrThrow(roomId: UUID): GameRoom =
        gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)

    fun validateHostAccess(room: GameRoom, userId: UUID, action: String) {
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, action, "Only the host can $action")
        }
    }

    fun validateRoomInStatus(room: GameRoom, expectedStatus: RoomStatus, action: String) {
        if (room.getStatusEnum() != expectedStatus) {
            throw InvalidRoomStateException(room.id, room.getStatusEnum(), action)
        }
    }

    suspend fun getActivePlayersWithMinCheck(roomId: UUID): List<RoomPlayer> {
        val players = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()
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

    suspend fun createHostPlayer(roomId: UUID, hostId: UUID, username: String, avatarUrl: String?): RoomPlayer {
        val hostPlayer = RoomPlayer(
            roomId = roomId,
            userId = hostId,
            username = username,
            avatarUrl = avatarUrl,
            role = RoomPlayer.roleFromEnum(PlayerRole.HOST)
        )
        return roomPlayerRepository.save(hostPlayer).awaitFirst()
    }

    suspend fun buildGameSettings(roomId: UUID): GameSettings {
        val settings = roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()
        return GameSettings(
            timeForAnswer = settings?.timeForAnswer ?: 30,
            timeForChoice = settings?.timeForChoice ?: 60,
            allowWrongAnswer = settings?.allowWrongAnswer ?: true,
            showRightAnswer = settings?.showRightAnswer ?: true
        )
    }

    suspend fun updateRoomToStarting(room: GameRoom, playersCount: Int): GameRoom {
        val updatedRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.STARTING),
            startedAt = LocalDateTime.now()
        )
        val savedRoom = gameRoomRepository.save(updatedRoom).awaitFirst()
        roomCacheService.cacheRoomData(savedRoom, playersCount)
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

            updateRoomToPlaying(savedRoom, activePlayers.size)
            publishRoomStartedEvent(roomId, gameSession, room.packId, activePlayers)

            return StartGameResponse(
                gameId = gameSession.gameSessionId,
                websocketUrl = gameSession.wsUrl
            )
        } catch (e: Exception) {
            logger.error(e) { "Failed to create game session, rolling back room status" }
            rollbackRoomToWaiting(savedRoom, activePlayers.size)
            throw e
        }
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

    suspend fun cancelRoom(room: GameRoom) {
        val updatedRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()
    }

    suspend fun getActivePlayers(roomId: UUID): List<RoomPlayer> =
        roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()

    suspend fun clearRoomAndPlayersCache(roomId: UUID, players: List<RoomPlayer>) = coroutineScope {
        launch { roomCacheService.clearRoomCache(roomId) }
        players.forEach { player ->
            launch { roomCacheService.deleteUserCurrentRoom(player.userId) }
        }
    }

    fun recordRoomCancelledMetrics(room: GameRoom, playersCount: Int) {
        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        repeat(playersCount) { lobbyMetrics.recordPlayerLeft() }
    }

    suspend fun findSettingsByRoomId(roomId: UUID) =
        roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()

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
            roomSettingsRepository.insertRoomSettings(
                roomId = settings.roomId,
                timeForAnswer = settings.timeForAnswer,
                timeForChoice = settings.timeForChoice,
                allowWrongAnswer = settings.allowWrongAnswer,
                showRightAnswer = settings.showRightAnswer
            ).awaitFirstOrNull()
        } else {
            roomSettingsRepository.save(settings).awaitFirstOrNull()
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

    private suspend fun updateRoomToPlaying(room: GameRoom, playersCount: Int) {
        val playingRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.PLAYING))
        gameRoomRepository.save(playingRoom).awaitFirstOrNull()
        roomCacheService.cacheRoomData(playingRoom, playersCount)
    }

    private suspend fun rollbackRoomToWaiting(room: GameRoom, playersCount: Int) {
        logger.error { "Rolling back room ${room.id} status to WAITING" }
        val rollbackRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.WAITING),
            startedAt = null
        )
        gameRoomRepository.save(rollbackRoom).awaitFirstOrNull()
        roomCacheService.cacheRoomData(rollbackRoom, playersCount)
    }

    private fun buildPlayerEventDataList(players: List<RoomPlayer>): List<PlayerEventData> =
        players.map { player ->
            PlayerEventData(
                user_id = player.userId.toString(),
                username = player.username,
                role = player.role
            )
        }

    private suspend fun publishRoomStartedEvent(
        roomId: UUID,
        gameSession: GameSessionResponse,
        packId: UUID,
        players: List<RoomPlayer>
    ) {
        val playerEventDataList = buildPlayerEventDataList(players)
        kafkaEventPublisher.publishRoomStarted(
            roomId = roomId,
            gameId = gameSession.gameSessionId,
            packId = packId,
            players = playerEventDataList
        )
    }
}
