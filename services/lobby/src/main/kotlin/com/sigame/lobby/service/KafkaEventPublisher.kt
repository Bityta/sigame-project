package com.sigame.lobby.service

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.databind.SerializationFeature
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule
import com.sigame.lobby.metrics.LobbyMetrics
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import mu.KotlinLogging
import org.springframework.kafka.core.KafkaTemplate
import org.springframework.stereotype.Service
import java.time.Instant
import java.time.format.DateTimeFormatter
import java.util.UUID

private val logger = KotlinLogging.logger {}

enum class LobbyEventType {
    ROOM_CREATED,
    PLAYER_JOINED,
    PLAYER_LEFT,
    ROOM_STARTED,
    ROOM_CANCELLED
}

data class LobbyEvent(
    val event_type: String,
    val timestamp: String,
    val payload: Map<String, Any>
)

data class PlayerEventData(
    val user_id: String,
    val username: String,
    val role: String
)

@Service
class KafkaEventPublisher(
    private val kafkaTemplate: KafkaTemplate<String, String>,
    private val lobbyMetrics: LobbyMetrics
) {
    
    private val topic = "lobby.events"
    
    private val objectMapper = ObjectMapper().apply {
        registerModule(JavaTimeModule())
        disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS)
    }
    
    private fun currentTimestamp(): String = 
        DateTimeFormatter.ISO_INSTANT.format(Instant.now())
    
        suspend fun publishRoomCreated(
        roomId: UUID,
        roomCode: String,
        hostId: UUID,
        hostUsername: String,
        packId: UUID,
        packName: String?,
        maxPlayers: Int,
        isPublic: Boolean
    ) {
        val event = LobbyEvent(
            event_type = LobbyEventType.ROOM_CREATED.name,
            timestamp = currentTimestamp(),
            payload = mapOf(
                "room_id" to roomId.toString(),
                "room_code" to roomCode,
                "host_id" to hostId.toString(),
                "host_username" to hostUsername,
                "pack_id" to packId.toString(),
                "pack_name" to (packName ?: "Unknown"),
                "max_players" to maxPlayers,
                "is_public" to isPublic
            )
        )
        publishEvent(roomId, event)
    }
    
        suspend fun publishPlayerJoined(
        roomId: UUID,
        userId: UUID,
        username: String,
        avatarUrl: String?,
        currentPlayers: Int,
        maxPlayers: Int
    ) {
        val event = LobbyEvent(
            event_type = LobbyEventType.PLAYER_JOINED.name,
            timestamp = currentTimestamp(),
            payload = mapOf(
                "room_id" to roomId.toString(),
                "user_id" to userId.toString(),
                "username" to username,
                "avatar_url" to (avatarUrl ?: ""),
                "current_players" to currentPlayers,
                "max_players" to maxPlayers
            )
        )
        publishEvent(roomId, event)
    }
    
        suspend fun publishPlayerLeft(
        roomId: UUID,
        userId: UUID,
        username: String,
        reason: String,
        currentPlayers: Int
    ) {
        val event = LobbyEvent(
            event_type = LobbyEventType.PLAYER_LEFT.name,
            timestamp = currentTimestamp(),
            payload = mapOf(
                "room_id" to roomId.toString(),
                "user_id" to userId.toString(),
                "username" to username,
                "reason" to reason,
                "current_players" to currentPlayers
            )
        )
        publishEvent(roomId, event)
    }
    
        suspend fun publishRoomStarted(
        roomId: UUID,
        gameId: String,
        packId: UUID,
        players: List<PlayerEventData>
    ) {
        val event = LobbyEvent(
            event_type = LobbyEventType.ROOM_STARTED.name,
            timestamp = currentTimestamp(),
            payload = mapOf(
                "room_id" to roomId.toString(),
                "game_id" to gameId,
                "pack_id" to packId.toString(),
                "players" to players.map { 
                    mapOf(
                        "user_id" to it.user_id,
                        "username" to it.username,
                        "role" to it.role
                    )
                }
            )
        )
        publishEvent(roomId, event)
    }
    
        suspend fun publishRoomCancelled(roomId: UUID, reason: String) {
        val event = LobbyEvent(
            event_type = LobbyEventType.ROOM_CANCELLED.name,
            timestamp = currentTimestamp(),
            payload = mapOf(
                "room_id" to roomId.toString(),
                "reason" to reason
            )
        )
        publishEvent(roomId, event)
    }
    
    private suspend fun publishEvent(roomId: UUID, event: LobbyEvent) = withContext(Dispatchers.IO) {
        try {
            val json = objectMapper.writeValueAsString(event)
            kafkaTemplate.send(topic, roomId.toString(), json).get()
            lobbyMetrics.recordKafkaEventPublished()
            logger.debug { "Published event: ${event.event_type} for room $roomId" }
        } catch (e: Exception) {
            logger.error(e) { "Error publishing event: ${event.event_type} for room $roomId" }
        }
    }
}

