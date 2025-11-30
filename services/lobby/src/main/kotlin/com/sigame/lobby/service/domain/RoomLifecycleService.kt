package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomSettingsDto
import com.sigame.lobby.domain.dto.StartGameResponse
import com.sigame.lobby.domain.dto.UpdateRoomSettingsRequest
import com.sigame.lobby.domain.dto.UpdateRoomSettingsResponse
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.InsufficientPlayersException
import com.sigame.lobby.domain.exception.InvalidRoomStateException
import com.sigame.lobby.domain.exception.PackNotApprovedException
import com.sigame.lobby.domain.exception.PackNotFoundException
import com.sigame.lobby.domain.exception.PackNotOwnedException
import com.sigame.lobby.domain.exception.PlayerAlreadyInRoomException
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.UnauthorizedRoomActionException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
import com.sigame.lobby.service.PlayerEventData
import com.sigame.lobby.service.RoomCodeGenerator
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.external.GameServiceClient
import com.sigame.lobby.service.external.GameSettings
import com.sigame.lobby.service.mapper.RoomMapper
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
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
    private val roomCodeGenerator: RoomCodeGenerator,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val gameServiceClient: GameServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val roomMapper: RoomMapper,
    private val passwordEncoder: BCryptPasswordEncoder = BCryptPasswordEncoder(12)
) {

    @Transactional
    suspend fun createRoom(hostId: UUID, request: CreateRoomRequest): RoomDto {
        logger.info { "Creating room for host $hostId with pack ${request.packId}" }

        validatePackForRoom(request.packId, hostId)
        validateUserNotInRoom(hostId)

        val savedRoom = createGameRoom(hostId, request)
        val roomId = savedRoom.id!!

        createRoomSettings(roomId, request.settings)
        addHostAsPlayer(roomId, hostId)
        updateCacheForNewRoom(savedRoom, hostId)
        publishRoomCreatedEvent(savedRoom, hostId, request)
        recordRoomCreatedMetrics()

        return roomMapper.toDto(savedRoom, 1)
    }

    @Transactional
    suspend fun startRoom(roomId: UUID, userId: UUID): StartGameResponse {
        logger.info { "Starting room $roomId by user $userId" }

        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)

        // Проверяем, что пользователь - хост
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "start room", "Only the host can start the room")
        }

        // Проверяем состояние комнаты
        if (room.getStatusEnum() != RoomStatus.WAITING) {
            throw InvalidRoomStateException(roomId, room.getStatusEnum(), "start")
        }

        // Проверяем количество игроков
        val activePlayers = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()
        if (activePlayers.size < 2) {
            throw InsufficientPlayersException(roomId, activePlayers.size)
        }

        // Обновляем статус на STARTING
        val updatedRoom = room.copy(
            status = GameRoom.statusFromEnum(RoomStatus.STARTING),
            startedAt = LocalDateTime.now()
        )
        val savedRoom = gameRoomRepository.save(updatedRoom).awaitFirst()

        // Обновляем кэш
        roomCacheService.cacheRoomData(savedRoom, activePlayers.size)

        // Метрики
        lobbyMetrics.recordRoomStarted()

        // Создаем игровую сессию
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

            // Обновляем статус на PLAYING
            val playingRoom = savedRoom.copy(
                status = GameRoom.statusFromEnum(RoomStatus.PLAYING)
            )
            gameRoomRepository.save(playingRoom).awaitFirstOrNull()
            roomCacheService.cacheRoomData(playingRoom, activePlayers.size)

            // Публикуем событие ROOM_STARTED после успешного создания игры
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

            // Откат статуса комнаты
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

        // Проверяем, что пользователь - хост
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "update settings", "Only the host can update settings")
        }

        // Проверяем состояние комнаты
        if (room.getStatusEnum() != RoomStatus.WAITING) {
            throw InvalidRoomStateException(roomId, room.getStatusEnum(), "update settings")
        }

        // Получаем текущие настройки
        val currentSettings = roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()

        // Обновляем настройки игры (PATCH - только переданные поля)
        val updatedSettings = if (currentSettings != null) {
            currentSettings.copy(
                timeForAnswer = request.timeForAnswer ?: currentSettings.timeForAnswer,
                timeForChoice = request.timeForChoice ?: currentSettings.timeForChoice,
                allowWrongAnswer = request.allowWrongAnswer ?: currentSettings.allowWrongAnswer,
                showRightAnswer = request.showRightAnswer ?: currentSettings.showRightAnswer
            )
        } else {
            // Создаём новые настройки с дефолтами
            RoomSettings(
                roomId = roomId,
                timeForAnswer = request.timeForAnswer ?: 30,
                timeForChoice = request.timeForChoice ?: 60,
                allowWrongAnswer = request.allowWrongAnswer ?: true,
                showRightAnswer = request.showRightAnswer ?: true
            )
        }

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

        // Обновляем кэш
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

        // Проверяем, что пользователь - хост
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, "delete room", "Only the host can delete the room")
        }

        // Обновляем статус на CANCELLED
        val updatedRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()

        // Получаем всех игроков для очистки их кэша
        val players = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()

        // Очищаем кэш
        roomCacheService.clearRoomCache(roomId)
        players.forEach { player ->
            roomCacheService.deleteUserCurrentRoom(player.userId)
        }

        // Метрики
        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        players.forEach { _ ->
            lobbyMetrics.recordPlayerLeft()
        }

        // Публикуем событие
        kafkaEventPublisher.publishRoomCancelled(roomId, "manual")
    }

    private suspend fun buildRoomDto(room: GameRoom): RoomDto {
        val players = roomPlayerRepository.findActiveByRoomId(room.id!!).asFlow().toList()
        val settings = roomSettingsRepository.findByRoomId(room.id).awaitFirstOrNull()

        return roomMapper.toDto(
            room = room,
            currentPlayers = players.size,
            players = players,
            settings = settings
        )
    }

    private suspend fun validatePackForRoom(packId: UUID, userId: UUID) {
        val validation = packServiceClient.validatePack(packId, userId)
            ?: throw PackNotFoundException(packId)

        if (!validation.exists) throw PackNotFoundException(packId)
        if (!validation.isApproved) throw PackNotApprovedException(packId, validation.status)
        if (!validation.isOwner) throw PackNotOwnedException(packId, userId)
    }

    private suspend fun validateUserNotInRoom(userId: UUID) {
        val existingPlayer = roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()
        if (existingPlayer != null) {
            throw PlayerAlreadyInRoomException(userId, existingPlayer.roomId)
        }
    }

    private suspend fun createGameRoom(hostId: UUID, request: CreateRoomRequest): GameRoom {
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

    private suspend fun createRoomSettings(roomId: UUID, settings: RoomSettingsDto?) {
        roomSettingsRepository.insertRoomSettings(
            roomId = roomId,
            timeForAnswer = settings?.timeForAnswer ?: 30,
            timeForChoice = settings?.timeForChoice ?: 60,
            allowWrongAnswer = settings?.allowWrongAnswer ?: true,
            showRightAnswer = settings?.showRightAnswer ?: true
        ).awaitFirstOrNull()
    }

    private suspend fun addHostAsPlayer(roomId: UUID, hostId: UUID) {
        val hostPlayer = RoomPlayer(
            roomId = roomId,
            userId = hostId,
            role = RoomPlayer.roleFromEnum(PlayerRole.HOST)
        )
        roomPlayerRepository.save(hostPlayer).awaitFirstOrNull()
    }

    private suspend fun updateCacheForNewRoom(room: GameRoom, hostId: UUID) {
        roomCacheService.cacheRoomData(room, 1)
        roomCacheService.setUserCurrentRoom(hostId, room.id!!)
        roomCacheService.addRoomPlayer(room.id, hostId)
    }

    private suspend fun publishRoomCreatedEvent(room: GameRoom, hostId: UUID, request: CreateRoomRequest) {
        val hostInfo = authServiceClient.getUserInfo(hostId)
        val packInfo = packServiceClient.getPackInfo(request.packId)

        kafkaEventPublisher.publishRoomCreated(
            roomId = room.id!!,
            roomCode = room.roomCode,
            hostId = hostId,
            hostUsername = hostInfo?.username ?: "Unknown",
            packId = request.packId,
            packName = packInfo?.name,
            maxPlayers = request.maxPlayers,
            isPublic = request.isPublic
        )
    }

    private fun recordRoomCreatedMetrics() {
        lobbyMetrics.recordRoomCreated()
        lobbyMetrics.recordPlayerJoined()
    }
}

