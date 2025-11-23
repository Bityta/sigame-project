package com.sigame.lobby.service

import com.fasterxml.jackson.databind.ObjectMapper
import com.sigame.lobby.metrics.LobbyMetrics
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import mu.KotlinLogging
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.stereotype.Service
import java.util.UUID

private val logger = KotlinLogging.logger {}

enum class GameEventType {
    ROOM_CREATED,
    PLAYER_JOINED,
    PLAYER_LEFT,
    ROOM_STARTED,
    ROOM_FINISHED,
    ROOM_CANCELLED
}

data class GameEvent(
    val type: GameEventType,
    val roomId: UUID,
    val data: Map<String, Any>
)

@Service
class KafkaEventPublisher(
    private val kafkaTemplate: KafkaTemplate<String, String>,
    private val lobbyMetrics: LobbyMetrics,
    private val objectMapper: ObjectMapper = ObjectMapper()
) {
    
    private val topic = "game.events"
    
    suspend fun publishRoomCreated(roomId: UUID, hostId: UUID, packId: UUID, maxPlayers: Int) {
        val event = GameEvent(
            type = GameEventType.ROOM_CREATED,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString(),
                "host_id" to hostId.toString(),
                "pack_id" to packId.toString(),
                "max_players" to maxPlayers
            )
        )
        publishEvent(event)
    }
    
    suspend fun publishPlayerJoined(roomId: UUID, userId: UUID, role: String) {
        val event = GameEvent(
            type = GameEventType.PLAYER_JOINED,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString(),
                "user_id" to userId.toString(),
                "role" to role
            )
        )
        publishEvent(event)
    }
    
    suspend fun publishPlayerLeft(roomId: UUID, userId: UUID) {
        val event = GameEvent(
            type = GameEventType.PLAYER_LEFT,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString(),
                "user_id" to userId.toString()
            )
        )
        publishEvent(event)
    }
    
    suspend fun publishRoomStarted(roomId: UUID, packId: UUID, players: List<UUID>) {
        val event = GameEvent(
            type = GameEventType.ROOM_STARTED,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString(),
                "pack_id" to packId.toString(),
                "players" to players.map { it.toString() }
            )
        )
        publishEvent(event)
    }
    
    suspend fun publishRoomFinished(roomId: UUID, duration: Long) {
        val event = GameEvent(
            type = GameEventType.ROOM_FINISHED,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString(),
                "duration" to duration
            )
        )
        publishEvent(event)
    }
    
    suspend fun publishRoomCancelled(roomId: UUID) {
        val event = GameEvent(
            type = GameEventType.ROOM_CANCELLED,
            roomId = roomId,
            data = mapOf(
                "room_id" to roomId.toString()
            )
        )
        publishEvent(event)
    }
    
    private suspend fun publishEvent(event: GameEvent) = withContext(Dispatchers.IO) {
        try {
            val json = objectMapper.writeValueAsString(event)
            kafkaTemplate.send(topic, event.roomId.toString(), json).get()
            lobbyMetrics.recordKafkaEventPublished()
            logger.debug { "Published event: ${event.type} for room ${event.roomId}" }
        } catch (e: Exception) {
            logger.error(e) { "Error publishing event: ${event.type} for room ${event.roomId}" }
        }
    }
}

