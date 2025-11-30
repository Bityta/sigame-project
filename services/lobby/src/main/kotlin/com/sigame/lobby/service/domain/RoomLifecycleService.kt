package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.dto.StartGameResponse
import com.sigame.lobby.domain.dto.UpdateRoomSettingsRequest
import com.sigame.lobby.domain.dto.UpdateRoomSettingsResponse
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.InsufficientPlayersException
import com.sigame.lobby.domain.exception.InvalidRoomStateException
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.UnauthorizedRoomActionException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
import com.sigame.lobby.service.PlayerEventData
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.external.GameServiceClient
import com.sigame.lobby.service.external.GameSettings
import com.sigame.lobby.service.mapper.RoomMapper
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import java.time.LocalDateTime
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class RoomLifecycleService(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val authServiceClient: AuthServiceClient,
    private val gameServiceClient: GameServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val roomMapper: RoomMapper,
    private val helper: RoomLifecycleHelper
) {

    @Transactional
    suspend fun createRoom(hostId: UUID, request: CreateRoomRequest): RoomDto = coroutineScope {
        logger.info { "Creating room for host $hostId with pack ${request.packId}" }

        val packValidation = async { helper.validatePackForRoom(request.packId, hostId) }
        val userValidation = async { helper.validateUserNotInRoom(hostId) }
        val hostInfoDeferred = async { helper.fetchRequiredUserInfo(hostId) }
        val packInfoDeferred = async { helper.fetchRequiredPackInfo(request.packId) }

        packValidation.await()
        userValidation.await()
        val hostInfo = hostInfoDeferred.await()
        val packInfo = packInfoDeferred.await()

        val savedRoom = helper.createGameRoom(hostId, request)
        val hostPlayer = awaitAll(
            async { helper.createRoomSettings(
                roomId = savedRoom.id,
                settings = request.settings
            ) },
            async { helper.createHostPlayer(
                roomId = savedRoom.id,
                hostId = hostId,
                username = hostInfo.username,
                avatarUrl = hostInfo.avatarUrl
            ) }
        ).last() as RoomPlayer

        launch { helper.updateCacheForNewRoom(savedRoom, hostId) }
        launch { helper.publishRoomCreatedEvent(savedRoom, hostId, hostInfo, packInfo.name, request) }
        helper.recordRoomCreatedMetrics()

        roomMapper.toDto(
            room = savedRoom,
            currentPlayers = 1,
            players = listOf(hostPlayer),
            packName = packInfo.name
        )
    }

    @Transactional
    suspend fun startRoom(roomId: UUID, userId: UUID): StartGameResponse {
        logger.info { "Starting room $roomId by user $userId" }

        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)

        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "start room", "Only the host can start the room")
        }

        if (room.getStatusEnum() != RoomStatus.WAITING) {
            throw InvalidRoomStateException(roomId, room.getStatusEnum(), "start")
        }

        val activePlayers = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()
        if (activePlayers.size < 2) {
            throw InsufficientPlayersException(roomId, activePlayers.size)
        }

        val updatedRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.STARTING),
            startedAt = LocalDateTime.now()
        )
        val savedRoom = gameRoomRepository.save(updatedRoom).awaitFirst()

        roomCacheService.cacheRoomData(savedRoom, activePlayers.size)
        lobbyMetrics.recordRoomStarted()

        try {
            val settings = roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()
            val gameSettings = GameSettings(
                timeForAnswer = settings?.timeForAnswer ?: 30,
                timeForChoice = settings?.timeForChoice ?: 60,
                allowWrongAnswer = settings?.allowWrongAnswer ?: true,
                showRightAnswer = settings?.showRightAnswer ?: true
            )

            val gameSessionResponse = gameServiceClient.createGameSession(
                roomId = roomId,
                packId = room.packId,
                players = activePlayers,
                settings = gameSettings
            )

            val playingRoom = savedRoom.copy(
                status = GameRoom.statusFromEnum(RoomStatus.PLAYING)
            )
            gameRoomRepository.save(playingRoom).awaitFirstOrNull()
            roomCacheService.cacheRoomData(playingRoom, activePlayers.size)

            val playerEventDataList = activePlayers.map { player ->
                val userInfo = authServiceClient.getUserInfo(player.userId)
                PlayerEventData(
                    user_id = player.userId.toString(),
                    username = userInfo?.username ?: "Unknown",
                    role = player.role
                )
            }
            kafkaEventPublisher.publishRoomStarted(
                roomId = roomId,
                gameId = gameSessionResponse.gameSessionId,
                packId = room.packId,
                players = playerEventDataList
            )

            return StartGameResponse(
                gameId = gameSessionResponse.gameSessionId,
                websocketUrl = gameSessionResponse.wsUrl
            )
        } catch (e: Exception) {
            logger.error(e) { "Failed to create game session, rolling back room status" }

            val rollbackRoom = savedRoom.copy(
                status = GameRoom.statusFromEnum(RoomStatus.WAITING),
                startedAt = null
            )
            gameRoomRepository.save(rollbackRoom).awaitFirstOrNull()
            roomCacheService.cacheRoomData(rollbackRoom, activePlayers.size)

            throw e
        }
    }

    @Transactional
    suspend fun updateRoomSettings(
        roomId: UUID,
        userId: UUID,
        request: UpdateRoomSettingsRequest
    ): UpdateRoomSettingsResponse {
        logger.info { "Updating room $roomId settings by user $userId" }

        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)

        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "update settings", "Only the host can update settings")
        }

        if (room.getStatusEnum() != RoomStatus.WAITING) {
            throw InvalidRoomStateException(roomId, room.getStatusEnum(), "update settings")
        }

        val currentSettings = roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()

        val updatedSettings = currentSettings?.copy(
            timeForAnswer = request.timeForAnswer ?: currentSettings.timeForAnswer,
            timeForChoice = request.timeForChoice ?: currentSettings.timeForChoice,
            allowWrongAnswer = request.allowWrongAnswer ?: currentSettings.allowWrongAnswer,
            showRightAnswer = request.showRightAnswer ?: currentSettings.showRightAnswer
        ) ?: RoomSettings(
            roomId = roomId,
            timeForAnswer = request.timeForAnswer ?: 30,
            timeForChoice = request.timeForChoice ?: 60,
            allowWrongAnswer = request.allowWrongAnswer ?: true,
            showRightAnswer = request.showRightAnswer ?: true
        )

        if (currentSettings != null) {
            roomSettingsRepository.save(updatedSettings).awaitFirstOrNull()
        } else {
            roomSettingsRepository.insertRoomSettings(
                roomId = roomId,
                timeForAnswer = updatedSettings.timeForAnswer,
                timeForChoice = updatedSettings.timeForChoice,
                allowWrongAnswer = updatedSettings.allowWrongAnswer,
                showRightAnswer = updatedSettings.showRightAnswer
            ).awaitFirstOrNull()
        }

        val currentPlayers = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        roomCacheService.cacheRoomData(room, currentPlayers)

        return UpdateRoomSettingsResponse(
            settings = RoomSettingsDto(
                timeForAnswer = updatedSettings.timeForAnswer,
                timeForChoice = updatedSettings.timeForChoice,
                allowWrongAnswer = updatedSettings.allowWrongAnswer,
                showRightAnswer = updatedSettings.showRightAnswer
            )
        )
    }

    @Transactional
    suspend fun deleteRoom(roomId: UUID, userId: UUID) {
        logger.info { "Deleting room $roomId by user $userId" }

        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)

        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "delete room", "Only the host can delete the room")
        }

        val updatedRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()

        val players = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()

        roomCacheService.clearRoomCache(roomId)
        players.forEach { player ->
            roomCacheService.deleteUserCurrentRoom(player.userId)
        }

        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        players.forEach { _ ->
            lobbyMetrics.recordPlayerLeft()
        }

        kafkaEventPublisher.publishRoomCancelled(roomId, "manual")
    }
}
