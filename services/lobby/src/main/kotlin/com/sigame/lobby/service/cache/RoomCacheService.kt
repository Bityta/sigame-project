package com.sigame.lobby.service.cache

import com.fasterxml.jackson.databind.ObjectMapper
import com.sigame.lobby.config.RoomConfig
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.model.GameRoom
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.data.redis.core.ReactiveRedisTemplate
import org.springframework.stereotype.Service
import java.time.Duration
import java.util.UUID

private val logger = KotlinLogging.logger {}

/**
 * Сервис для работы с кэшем комнат в Redis
 */
@Service
class RoomCacheService(
    private val redisTemplate: ReactiveRedisTemplate<String, String>,
    private val roomConfig: RoomConfig,
    private val objectMapper: ObjectMapper = ObjectMapper()
) {
    
    private val cacheTtl = Duration.ofSeconds(roomConfig.cacheTtl)
    
    /**
     * Кэширует данные комнаты
     */
    suspend fun cacheRoomData(room: GameRoom, currentPlayers: Int) {
        try {
            val metadata = mapOf(
                "room_code" to room.roomCode,
                "host_id" to room.hostId.toString(),
                "pack_id" to room.packId.toString(),
                "name" to room.name,
                "status" to room.status,
                "player_count" to currentPlayers.toString(),
                "max_players" to room.maxPlayers.toString()
            )
            
            setRoomMeta(room.id!!, metadata)
            
            // Добавляем в sorted set активных комнат только если статус WAITING
            if (room.getStatusEnum() == RoomStatus.WAITING) {
                addActiveRoom(room.id, room.createdAt.toEpochSecond(java.time.ZoneOffset.UTC))
            } else {
                removeActiveRoom(room.id)
            }
        } catch (e: Exception) {
            logger.error(e) { "Failed to cache room data for room ${room.id}" }
        }
    }
    
    /**
     * Добавляет комнату в sorted set активных комнат
     */
    suspend fun addActiveRoom(roomId: UUID, createdAt: Long) {
        try {
            redisTemplate.opsForZSet()
                .add("active_rooms", roomId.toString(), createdAt.toDouble())
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error adding active room to Redis" }
        }
    }
    
    /**
     * Удаляет комнату из sorted set активных комнат
     */
    suspend fun removeActiveRoom(roomId: UUID) {
        try {
            redisTemplate.opsForZSet()
                .remove("active_rooms", roomId.toString())
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error removing active room from Redis" }
        }
    }
    
    /**
     * Сохраняет метаданные комнаты
     */
    suspend fun setRoomMeta(roomId: UUID, metadata: Map<String, String>) {
        try {
            val key = "room:$roomId:meta"
            redisTemplate.opsForHash<String, String>()
                .putAll(key, metadata)
                .awaitFirstOrNull()
            
            redisTemplate.expire(key, cacheTtl).awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error setting room metadata in Redis" }
        }
    }
    
    /**
     * Получает метаданные комнаты
     */
    suspend fun getRoomMeta(roomId: UUID): Map<String, String>? {
        return try {
            val key = "room:$roomId:meta"
            redisTemplate.opsForHash<String, String>()
                .entries(key)
                .collectMap({ it.key }, { it.value })
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error getting room metadata from Redis" }
            null
        }
    }
    
    /**
     * Удаляет метаданные комнаты
     */
    suspend fun deleteRoomMeta(roomId: UUID) {
        try {
            redisTemplate.delete("room:$roomId:meta").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting room metadata from Redis" }
        }
    }
    
    /**
     * Устанавливает текущую комнату пользователя
     */
    suspend fun setUserCurrentRoom(userId: UUID, roomId: UUID) {
        try {
            val key = "user:$userId:current_room"
            redisTemplate.opsForValue()
                .set(key, roomId.toString(), cacheTtl)
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error setting user current room in Redis" }
        }
    }
    
    /**
     * Получает текущую комнату пользователя
     */
    suspend fun getUserCurrentRoom(userId: UUID): UUID? {
        return try {
            val key = "user:$userId:current_room"
            val value = redisTemplate.opsForValue().get(key).awaitFirstOrNull()
            value?.let { UUID.fromString(it) }
        } catch (e: Exception) {
            logger.error(e) { "Error getting user current room from Redis" }
            null
        }
    }
    
    /**
     * Удаляет текущую комнату пользователя
     */
    suspend fun deleteUserCurrentRoom(userId: UUID) {
        try {
            redisTemplate.delete("user:$userId:current_room").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting user current room from Redis" }
        }
    }
    
    /**
     * Добавляет игрока в комнату
     */
    suspend fun addRoomPlayer(roomId: UUID, userId: UUID) {
        try {
            val key = "room:$roomId:players"
            redisTemplate.opsForSet()
                .add(key, userId.toString())
                .awaitFirstOrNull()
            
            redisTemplate.expire(key, cacheTtl).awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error adding room player to Redis" }
        }
    }
    
    /**
     * Удаляет игрока из комнаты
     */
    suspend fun removeRoomPlayer(roomId: UUID, userId: UUID) {
        try {
            val key = "room:$roomId:players"
            redisTemplate.opsForSet()
                .remove(key, userId.toString())
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error removing room player from Redis" }
        }
    }
    
    /**
     * Получает список игроков в комнате
     */
    suspend fun getRoomPlayers(roomId: UUID): Set<UUID> {
        return try {
            val key = "room:$roomId:players"
            redisTemplate.opsForSet()
                .members(key)
                .collectList()
                .awaitFirstOrNull()
                ?.mapNotNull { 
                    try {
                        UUID.fromString(it)
                    } catch (e: Exception) {
                        null
                    }
                }
                ?.toSet() ?: emptySet()
        } catch (e: Exception) {
            logger.error(e) { "Error getting room players from Redis" }
            emptySet()
        }
    }
    
    /**
     * Удаляет всех игроков комнаты
     */
    suspend fun deleteRoomPlayers(roomId: UUID) {
        try {
            redisTemplate.delete("room:$roomId:players").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting room players from Redis" }
        }
    }
    
    /**
     * Очищает весь кэш комнаты
     */
    suspend fun clearRoomCache(roomId: UUID) {
        removeActiveRoom(roomId)
        deleteRoomMeta(roomId)
        deleteRoomPlayers(roomId)
    }
}

