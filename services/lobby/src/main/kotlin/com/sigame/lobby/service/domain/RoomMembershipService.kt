package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.JoinRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.*
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.grpc.AuthServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.KafkaEventPublisher
import com.sigame.lobby.service.cache.RoomCacheService
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
 * Сервис для управления членством игроков в комнатах (присоединение/выход)
 */
@Service
class RoomMembershipService(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomCacheService: RoomCacheService,
    private val kafkaEventPublisher: KafkaEventPublisher,
    private val authServiceClient: AuthServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val passwordEncoder: BCryptPasswordEncoder = BCryptPasswordEncoder(12)
) {
    
    /**
     * Присоединяет пользователя к комнате
     */
    @Transactional
    suspend fun joinRoom(roomId: UUID, userId: UUID, request: JoinRoomRequest): RoomDto {
        logger.info { "User $userId joining room $roomId" }
        
        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)
        
        // Проверяем состояние комнаты
        if (room.getStatusEnum() != RoomStatus.WAITING) {
            throw InvalidRoomStateException(roomId, room.getStatusEnum(), "join")
        }
        
        // Проверяем пароль
        if (room.passwordHash != null) {
            if (request.password == null || !passwordEncoder.matches(request.password, room.passwordHash)) {
                throw InvalidPasswordException()
            }
        }
        
        // Проверяем вместимость
        val currentPlayers = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        if (currentPlayers >= room.maxPlayers) {
            throw RoomFullException(roomId)
        }
        
        // Проверяем, не находится ли пользователь уже в этой комнате
        val existingPlayer = roomPlayerRepository.findByRoomIdAndUserId(roomId, userId).awaitFirstOrNull()
        if (existingPlayer != null && existingPlayer.leftAt == null) {
            throw PlayerAlreadyInRoomException(userId, roomId)
        }
        
        // Проверяем, не находится ли пользователь в другой комнате
        val currentPlayer = roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()
        if (currentPlayer != null && currentPlayer.roomId != roomId) {
            throw PlayerAlreadyInRoomException(userId, currentPlayer.roomId)
        }
        
        val roleStr = request.role.lowercase()
        
        if (existingPlayer != null) {
            // Пользователь возвращается в комнату
            val rejoinedPlayer = existingPlayer.copy(leftAt = null, joinedAt = LocalDateTime.now())
            roomPlayerRepository.save(rejoinedPlayer).awaitFirstOrNull()
        } else {
            // Новый игрок
            val player = RoomPlayer(
                roomId = roomId,
                userId = userId,
                role = roleStr
            )
            roomPlayerRepository.save(player).awaitFirstOrNull()
        }
        
        // Обновляем кэш
        roomCacheService.setUserCurrentRoom(userId, roomId)
        roomCacheService.addRoomPlayer(roomId, userId)
        roomCacheService.cacheRoomData(room, currentPlayers + 1)
        
        // Метрики
        lobbyMetrics.recordPlayerJoined()
        
        // Получаем информацию о пользователе для события
        val userInfo = authServiceClient.getUserInfo(userId)
        
        // Публикуем событие
        kafkaEventPublisher.publishPlayerJoined(
            roomId = roomId,
            userId = userId,
            username = userInfo?.username ?: "Unknown",
            avatarUrl = userInfo?.avatarUrl,
            currentPlayers = currentPlayers + 1,
            maxPlayers = room.maxPlayers
        )
        
        // Возвращаем обновленную информацию о комнате
        // Используем прямой запрос вместо roomQueryService для избежания circular dependency
        val updatedRoom = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)
        
        val updatedPlayers = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()
        val currentPlayersCount = updatedPlayers.size
        
        return RoomDto(
            id = updatedRoom.id!!,
            roomCode = updatedRoom.roomCode,
            name = updatedRoom.name,
            hostId = updatedRoom.hostId,
            hostUsername = null,
            packId = updatedRoom.packId,
            packName = null,
            status = updatedRoom.status,
            maxPlayers = updatedRoom.maxPlayers,
            currentPlayers = currentPlayersCount,
            isPublic = updatedRoom.isPublic,
            hasPassword = updatedRoom.passwordHash != null,
            createdAt = updatedRoom.createdAt,
            players = null,
            settings = null
        )
    }
    
    /**
     * Удаляет пользователя из комнаты
     */
    @Transactional
    suspend fun leaveRoom(roomId: UUID, userId: UUID) {
        logger.info { "User $userId leaving room $roomId" }
        
        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)
        
        val player = roomPlayerRepository.findByRoomIdAndUserId(roomId, userId).awaitFirstOrNull()
            ?: throw PlayerNotInRoomException(userId, roomId)
        
        if (player.leftAt != null) {
            throw IllegalStateException("Player has already left")
        }
        
        // Отмечаем игрока как покинувшего
        val updatedPlayer = player.copy(leftAt = LocalDateTime.now())
        roomPlayerRepository.save(updatedPlayer).awaitFirstOrNull()
        
        // Обновляем кэш
        roomCacheService.deleteUserCurrentRoom(userId)
        roomCacheService.removeRoomPlayer(roomId, userId)
        
        // Метрики
        lobbyMetrics.recordPlayerLeft()
        
        // Если это был хост, передаем роль или отменяем комнату
        if (player.getRoleEnum() == PlayerRole.HOST) {
            handleHostLeaving(room, roomId)
        } else {
            handlePlayerLeaving(room, roomId)
        }
        
        // Получаем информацию о пользователе для события
        val userInfo = authServiceClient.getUserInfo(userId)
        val currentPlayersCount = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        
        // Публикуем событие
        kafkaEventPublisher.publishPlayerLeft(
            roomId = roomId,
            userId = userId,
            username = userInfo?.username ?: "Unknown",
            reason = "left",
            currentPlayers = currentPlayersCount
        )
    }
    
    /**
     * Обрабатывает выход хоста из комнаты
     */
    private suspend fun handleHostLeaving(room: GameRoom, roomId: UUID) {
        val activePlayers = roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()
        
        if (activePlayers.isEmpty()) {
            // Нет игроков - отменяем комнату
            cancelRoom(room, roomId)
        } else {
            // Передаем роль хоста первому игроку
            transferHost(room, roomId, activePlayers)
        }
    }
    
    /**
     * Обрабатывает выход обычного игрока
     */
    private suspend fun handlePlayerLeaving(room: GameRoom, roomId: UUID) {
        val currentPlayers = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        
        if (currentPlayers == 0) {
            // Нет игроков - отменяем комнату
            cancelRoom(room, roomId)
        } else {
            // Обновляем кэш
            roomCacheService.cacheRoomData(room, currentPlayers)
        }
    }
    
    /**
     * Отменяет комнату
     */
    private suspend fun cancelRoom(room: GameRoom, roomId: UUID) {
        logger.info { "No players left in room $roomId, cancelling room" }
        
        val updatedRoom = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()
        
        // Очищаем кэш
        roomCacheService.clearRoomCache(roomId)
        
        // Метрики
        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        
        // Публикуем событие
        kafkaEventPublisher.publishRoomCancelled(roomId, "no_players")
    }
    
    /**
     * Передает роль хоста другому игроку
     */
    private suspend fun transferHost(room: GameRoom, roomId: UUID, activePlayers: List<RoomPlayer>) {
        val newHost = activePlayers.first()
        logger.info { "Transferring host role in room $roomId to ${newHost.userId}" }
        
        // Обновляем роль нового хоста
        val updatedNewHost = newHost.copy(role = RoomPlayer.roleFromEnum(PlayerRole.HOST))
        roomPlayerRepository.save(updatedNewHost).awaitFirstOrNull()
        
        // Обновляем комнату
        val updatedRoom = room.copy(hostId = newHost.userId)
        gameRoomRepository.save(updatedRoom).awaitFirstOrNull()
        
        // Обновляем кэш
        roomCacheService.cacheRoomData(updatedRoom, activePlayers.size)
        
        // Получаем информацию о новом хосте
        val newHostInfo = authServiceClient.getUserInfo(newHost.userId)
        
        // Публикуем событие о смене хоста
        kafkaEventPublisher.publishPlayerJoined(
            roomId = roomId,
            userId = newHost.userId,
            username = newHostInfo?.username ?: "Unknown",
            avatarUrl = newHostInfo?.avatarUrl,
            currentPlayers = activePlayers.size,
            maxPlayers = room.maxPlayers
        )
    }
}

