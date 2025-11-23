package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomListResponse
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.RoomNotFoundByCodeException
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.service.batch.BatchOperationService
import com.sigame.lobby.service.mapper.RoomMapper
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.stereotype.Service
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Сервис для получения информации о комнатах
 */
@Service
class RoomQueryService(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomMapper: RoomMapper,
    private val batchOperationService: BatchOperationService
) {
    
    /**
     * Получает список комнат с фильтрацией и пагинацией
     * Оптимизировано с помощью batch операций для избежания N+1 проблемы
     */
    suspend fun getRooms(
        page: Int,
        size: Int,
        status: RoomStatus?,
        hasSlots: Boolean?
    ): RoomListResponse {
        val offset = page * size
        
        val rooms = when {
            status != null -> gameRoomRepository.findByStatus(status.name.lowercase(), size, offset)
            hasSlots == true -> gameRoomRepository.findPublicWaitingRooms(size, offset)
            else -> gameRoomRepository.findPublicWaitingRooms(size, offset)
        }.asFlow().toList()
        
        val total = when {
            status != null -> gameRoomRepository.countByStatus(status.name.lowercase()).awaitFirst()
            hasSlots == true -> gameRoomRepository.countPublicWaitingRooms().awaitFirst()
            else -> gameRoomRepository.countPublicWaitingRooms().awaitFirst()
        }
        
        // Batch загрузка данных для оптимизации N+1 проблемы
        val hostIds = rooms.map { it.hostId }
        val packIds = rooms.map { it.packId }
        
        val userInfoMap = batchOperationService.getUserInfoBatch(hostIds)
        val packInfoMap = batchOperationService.getPackInfoBatch(packIds)
        
        val roomDtos = rooms.map { room ->
            val playerCount = roomPlayerRepository.countActiveByRoomId(room.id!!).awaitFirst().toInt()
            roomMapper.toDtoWithCache(room, playerCount, userInfoMap, packInfoMap)
        }
        
        return RoomListResponse(
            rooms = roomDtos,
            total = total,
            page = page,
            size = size
        )
    }
    
    /**
     * Получает комнату по ID
     */
    suspend fun getRoomById(roomId: UUID): RoomDto {
        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)
        
        return buildDetailedRoomDto(room)
    }
    
    /**
     * Получает комнату по коду
     */
    suspend fun getRoomByCode(code: String): RoomDto {
        val room = gameRoomRepository.findByRoomCode(code).awaitFirstOrNull()
            ?: throw RoomNotFoundByCodeException(code)
        
        return buildDetailedRoomDto(room)
    }
    
    /**
     * Получает текущую комнату пользователя
     */
    suspend fun getUserCurrentRoom(userId: UUID): RoomDto? {
        val player = roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()
            ?: return null
        
        return try {
            getRoomById(player.roomId)
        } catch (e: RoomNotFoundException) {
            logger.warn { "Player $userId has active player record but room ${player.roomId} not found" }
            null
        }
    }
    
    /**
     * Строит подробный DTO комнаты со всеми игроками и настройками
     */
    private suspend fun buildDetailedRoomDto(room: com.sigame.lobby.domain.model.GameRoom): RoomDto {
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

