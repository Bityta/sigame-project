package com.sigame.lobby.service.cache

import com.sigame.lobby.config.RoomConfig
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.model.GameRoom
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.data.redis.core.ReactiveRedisTemplate
import org.springframework.stereotype.Service
import java.time.Duration
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class RoomCacheService(
    private val redisTemplate: ReactiveRedisTemplate<String, String>,
    roomConfig: RoomConfig
) {

    private val cacheTtl = Duration.ofSeconds(roomConfig.cacheTtl)

    suspend fun cacheRoomData(room: GameRoom, currentPlayers: Int) {
        try {
            val metadata = buildMap<String, String>(7) {
                put("room_code", room.roomCode)
                put("host_id", room.hostId.toString())
                put("pack_id", room.packId.toString())
                put("name", room.name)
                put("status", room.status)
                put("player_count", currentPlayers.toString())
                put("max_players", room.maxPlayers.toString())
            }

            setRoomMeta(room.id, metadata)

            if (room.getStatusEnum() == RoomStatus.WAITING) {
                addActiveRoom(room.id, room.createdAt.toEpochSecond(java.time.ZoneOffset.UTC))
            } else {
                removeActiveRoom(room.id)
            }
        } catch (e: Exception) {
            logger.error(e) { "Failed to cache room data for room ${room.id}" }
        }
    }

    suspend fun addActiveRoom(roomId: UUID, createdAt: Long) {
        try {
            redisTemplate.opsForZSet()
                .add("active_rooms", roomId.toString(), createdAt.toDouble())
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error adding active room to Redis" }
        }
    }

    suspend fun removeActiveRoom(roomId: UUID) {
        try {
            redisTemplate.opsForZSet()
                .remove("active_rooms", roomId.toString())
                .awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error removing active room from Redis" }
        }
    }

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

    suspend fun deleteRoomMeta(roomId: UUID) {
        try {
            redisTemplate.delete("room:$roomId:meta").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting room metadata from Redis" }
        }
    }

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

    suspend fun deleteUserCurrentRoom(userId: UUID) {
        try {
            redisTemplate.delete("user:$userId:current_room").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting user current room from Redis" }
        }
    }

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

    suspend fun deleteRoomPlayers(roomId: UUID) {
        try {
            redisTemplate.delete("room:$roomId:players").awaitFirstOrNull()
        } catch (e: Exception) {
            logger.error(e) { "Error deleting room players from Redis" }
        }
    }

    suspend fun clearRoomCache(roomId: UUID) = coroutineScope {
        launch { removeActiveRoom(roomId) }
        launch { deleteRoomMeta(roomId) }
        launch { deleteRoomPlayers(roomId) }
    }
}

