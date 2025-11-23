package com.sigame.lobby.metrics

import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.enums.RoomStatus
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.runBlocking
import mu.KotlinLogging
import org.springframework.boot.context.event.ApplicationReadyEvent
import org.springframework.context.event.EventListener
import org.springframework.stereotype.Component

private val logger = KotlinLogging.logger {}

@Component
class MetricsInitializer(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val lobbyMetrics: LobbyMetrics
) {
    
    @EventListener(ApplicationReadyEvent::class)
    fun initializeMetrics() = runBlocking {
        logger.info { "Initializing metrics from database..." }
        
        try {
            // Count rooms by status
            RoomStatus.values().forEach { status ->
                val count = gameRoomRepository.countByStatus(status.name.lowercase()).awaitFirst()
                lobbyMetrics.setRoomsByStatus(status, count.toInt())
                logger.info { "Set rooms with status $status to $count" }
            }
            
            // Count total active rooms (waiting, starting, playing)
            val activeStatuses = listOf("waiting", "starting", "playing")
            var totalActiveRooms = 0
            activeStatuses.forEach { status ->
                totalActiveRooms += gameRoomRepository.countByStatus(status).awaitFirst().toInt()
            }
            lobbyMetrics.setActiveRooms(totalActiveRooms)
            logger.info { "Set total active rooms to $totalActiveRooms" }
            
            // Count active players (those who haven't left)
            val activePlayers = roomPlayerRepository.findByLeftAtIsNull().count().awaitFirst()
            lobbyMetrics.setActivePlayers(activePlayers.toInt())
            logger.info { "Set active players to $activePlayers" }
            
            logger.info { "Metrics initialization completed successfully" }
        } catch (e: Exception) {
            logger.error(e) { "Failed to initialize metrics from database" }
        }
    }
}

