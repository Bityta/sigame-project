package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.StartGameResponse
import com.sigame.lobby.domain.dto.UpdateRoomSettingsRequest
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.*
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.PackServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
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

/**
 * Сервис для управления жизненным циклом комнат (создание, старт, удаление)
 */
@Service
class RoomLifecycleService(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomCodeGenerator: RoomCodeGenerator,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val packServiceClient: PackServiceClient,
    private val gameServiceClient: GameServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val roomMapper: RoomMapper,
    private val passwordEncoder: BCryptPasswordEncoder = BCryptPasswordEncoder(12)
) {
    
    /**
     * Создает новую комнату
     */
    @Transactional
    suspend fun createRoom(hostId: UUID, request: CreateRoomRequest): RoomDto {
        logger.info { "Creating room for host $hostId with pack ${request.packId}" }
        
        // Валидация пака
        val packExists = packServiceClient.validatePackExists(request.packId)
        if (!packExists) {
            throw PackNotFoundException(request.packId)
        }
        
        // Проверяем, не находится ли пользователь уже в комнате
        val existingPlayer = roomPlayerRepository.findActiveByUserId(hostId).awaitFirstOrNull()
        if (existingPlayer != null) {
            val existingRoomCode = gameRoomRepository.findById(existingPlayer.roomId).awaitFirstOrNull()?.roomCode
            logger.info { "User $hostId is already in room $existingRoomCode, will need to leave first" }
            throw PlayerAlreadyInRoomException(hostId, existingPlayer.roomId)
        }
        
        // Генерируем уникальный код комнаты
        val roomCode = roomCodeGenerator.generateUniqueCode()
        
        // Создаем комнату
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
        
        val savedRoom = gameRoomRepository.save(room).awaitFirst()
        val roomId = savedRoom.id!!
        
        // Создаем настройки комнаты
        val roomSettings = RoomSettings(
            roomId = roomId,
            timeForAnswer = request.settings?.timeForAnswer ?: 30,
            timeForChoice = request.settings?.timeForChoice ?: 60,
            allowWrongAnswer = request.settings?.allowWrongAnswer ?: true,
            showRightAnswer = request.settings?.showRightAnswer ?: true
        )
        // Используем прямой INSERT вместо save(), чтобы избежать попытки UPDATE
        roomSettingsRepository.insertRoomSettings(
            roomId = roomSettings.roomId,
            timeForAnswer = roomSettings.timeForAnswer,
            timeForChoice = roomSettings.timeForChoice,
            allowWrongAnswer = roomSettings.allowWrongAnswer,
            showRightAnswer = roomSettings.showRightAnswer
        ).awaitFirstOrNull()
        
        // Добавляем хоста как игрока
        val hostPlayer = RoomPlayer(
            roomId = roomId,
            userId = hostId,
            role = RoomPlayer.roleFromEnum(PlayerRole.HOST)
        )
        roomPlayerRepository.save(hostPlayer).awaitFirstOrNull()
        
        // Кэшируем данные
        roomCacheService.cacheRoomData(savedRoom, 1)
        roomCacheService.setUserCurrentRoom(hostId, roomId)
        roomCacheService.addRoomPlayer(roomId, hostId)
        
        // Публикуем событие
        kafkaEventPublisher.publishRoomCreated(
            roomId = roomId,
            hostId = hostId,
            packId = request.packId,
            maxPlayers = request.maxPlayers
        )
        
        // Метрики
        lobbyMetrics.recordRoomCreated()
        lobbyMetrics.recordPlayerJoined()
        
        return roomMapper.toDto(savedRoom, 1)
    }
    
    /**
     * Запускает игру в комнате
     */
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
        
        // Публикуем событие
        kafkaEventPublisher.publishRoomStarted(
            roomId = roomId,
            packId = room.packId,
            players = activePlayers.map { it.userId }
        )
        
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
            
            return StartGameResponse(
                gameSessionId = gameSessionResponse.gameSessionId,
                wsUrl = gameSessionResponse.wsUrl
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
    
    /**
     * Обновляет настройки комнаты
     */
    @Transactional
    suspend fun updateRoomSettings(
        roomId: UUID,
        userId: UUID,
        request: UpdateRoomSettingsRequest
    ): RoomDto {
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
        
        // Обновляем комнату
        val passwordHash = request.password?.let { passwordEncoder.encode(it) }
        val updatedRoom = room.copy(
            maxPlayers = request.maxPlayers ?: room.maxPlayers,
            isPublic = request.isPublic ?: room.isPublic,
            passwordHash = if (request.password != null) passwordHash else room.passwordHash
        )
        val savedRoom = gameRoomRepository.save(updatedRoom).awaitFirst()
        
        // Обновляем настройки игры
        if (request.settings != null) {
            val settings = roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()
            if (settings != null) {
                val updatedSettings = settings.copy(
                    timeForAnswer = request.settings.timeForAnswer,
                    timeForChoice = request.settings.timeForChoice,
                    allowWrongAnswer = request.settings.allowWrongAnswer,
                    showRightAnswer = request.settings.showRightAnswer
                )
                roomSettingsRepository.save(updatedSettings).awaitFirstOrNull()
            } else {
                // Если настроек нет (странная ситуация), создаём через INSERT
                roomSettingsRepository.insertRoomSettings(
                    roomId = roomId,
                    timeForAnswer = request.settings.timeForAnswer,
                    timeForChoice = request.settings.timeForChoice,
                    allowWrongAnswer = request.settings.allowWrongAnswer,
                    showRightAnswer = request.settings.showRightAnswer
                ).awaitFirstOrNull()
            }
        }
        
        // Обновляем кэш
        val currentPlayers = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        roomCacheService.cacheRoomData(savedRoom, currentPlayers)
        
        return buildRoomDto(savedRoom)
    }
    
    /**
     * Удаляет комнату
     */
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
        kafkaEventPublisher.publishRoomCancelled(roomId)
    }
    
    /**
     * Строит RoomDto из GameRoom
     */
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
}

