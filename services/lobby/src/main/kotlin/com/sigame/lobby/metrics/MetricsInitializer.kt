package com.sigame.lobby.metrics

import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
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
        try {
            RoomStatus.entries.forEach { status ->
                val count = gameRoomRepository.countByStatus(status.name.lowercase()).awaitFirst()
                lobbyMetrics.setRoomsByStatus(status, count.toInt())
            }

            val activeStatuses = listOf("waiting", "starting", "playing")
            val totalActiveRooms = activeStatuses.sumOf {
                gameRoomRepository.countByStatus(it).awaitFirst().toInt()
            }
            lobbyMetrics.setActiveRooms(totalActiveRooms)

            val activePlayers = roomPlayerRepository.findByLeftAtIsNull().count().awaitFirst()
            lobbyMetrics.setActivePlayers(activePlayers.toInt())

            logger.info { "Metrics initialized: $totalActiveRooms active rooms, $activePlayers active players" }
        } catch (e: Exception) {
            logger.error(e) { "Failed to initialize metrics from database" }
        }
    }
}
