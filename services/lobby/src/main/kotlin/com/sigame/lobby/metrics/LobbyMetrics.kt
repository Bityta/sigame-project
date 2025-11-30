package com.sigame.lobby.metrics

import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.sse.service.RoomEventPublisher
import io.micrometer.core.instrument.Counter
import io.micrometer.core.instrument.Gauge
import io.micrometer.core.instrument.MeterRegistry
import mu.KotlinLogging
import org.springframework.stereotype.Component
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.atomic.AtomicInteger

private val logger = KotlinLogging.logger {}

@Component
class LobbyMetrics(
    private val meterRegistry: MeterRegistry,
    private val roomEventPublisher: RoomEventPublisher
) {
    private val activeRooms = AtomicInteger(0)
    private val activePlayers = AtomicInteger(0)
    private val roomsByStatus = ConcurrentHashMap<RoomStatus, AtomicInteger>()

    init {
        Counter.builder("lobby_rooms_created_total")
            .description("Total number of rooms created")
            .register(meterRegistry)

        Counter.builder("lobby_rooms_started_total")
            .description("Total number of rooms started")
            .register(meterRegistry)

        Counter.builder("lobby_rooms_cancelled_total")
            .description("Total number of rooms cancelled")
            .register(meterRegistry)

        Counter.builder("lobby_players_joined_total")
            .description("Total number of players joined")
            .register(meterRegistry)

        Counter.builder("lobby_players_left_total")
            .description("Total number of players left")
            .register(meterRegistry)

        Gauge.builder("lobby_active_rooms", activeRooms) { it.get().toDouble() }
            .description("Current number of active rooms")
            .register(meterRegistry)

        Gauge.builder("lobby_active_players", activePlayers) { it.get().toDouble() }
            .description("Current number of active players in rooms")
            .register(meterRegistry)

        Gauge.builder("lobby_sse_active_sinks", roomEventPublisher) { it.getActiveRoomsCount().toDouble() }
            .description("Number of active SSE sinks")
            .register(meterRegistry)

        Gauge.builder("lobby_sse_published_total", roomEventPublisher) { it.getPublishedCount().toDouble() }
            .description("Total SSE events published")
            .register(meterRegistry)

        Gauge.builder("lobby_sse_dropped_total", roomEventPublisher) { it.getDroppedCount().toDouble() }
            .description("Total SSE events dropped (no subscribers)")
            .register(meterRegistry)

        Gauge.builder("lobby_sse_failed_total", roomEventPublisher) { it.getFailedCount().toDouble() }
            .description("Total SSE events failed to emit")
            .register(meterRegistry)

        RoomStatus.entries.forEach { status ->
            roomsByStatus[status] = AtomicInteger(0)
            Gauge.builder("lobby_rooms_by_status", roomsByStatus[status]!!) { it.get().toDouble() }
                .description("Number of rooms by status")
                .tag("status", status.name)
                .register(meterRegistry)
        }

        logger.info { "Lobby Metrics initialized" }
    }

    fun recordRoomCreated() {
        meterRegistry.counter("lobby_rooms_created_total").increment()
        activeRooms.incrementAndGet()
        incrementRoomStatus(RoomStatus.WAITING)
    }

    fun recordRoomStarted() {
        meterRegistry.counter("lobby_rooms_started_total").increment()
        decrementRoomStatus(RoomStatus.WAITING)
        incrementRoomStatus(RoomStatus.STARTING)
    }

    fun recordRoomCancelled(oldStatus: RoomStatus) {
        meterRegistry.counter("lobby_rooms_cancelled_total").increment()
        activeRooms.decrementAndGet()
        decrementRoomStatus(oldStatus)
    }

    fun recordPlayerJoined() {
        meterRegistry.counter("lobby_players_joined_total").increment()
        activePlayers.incrementAndGet()
    }

    fun recordPlayerLeft() {
        meterRegistry.counter("lobby_players_left_total").increment()
        activePlayers.decrementAndGet()
    }

    fun recordGrpcCall(serviceName: String) {
        Counter.builder("lobby_grpc_calls_total")
            .tag("service", serviceName)
            .register(meterRegistry)
            .increment()
    }

    fun recordGrpcError(serviceName: String, errorType: String) {
        Counter.builder("lobby_grpc_errors_total")
            .tag("service", serviceName)
            .tag("error_type", errorType)
            .register(meterRegistry)
            .increment()
    }

    fun setActiveRooms(count: Int) {
        activeRooms.set(count)
    }

    fun setActivePlayers(count: Int) {
        activePlayers.set(count)
    }

    fun setRoomsByStatus(status: RoomStatus, count: Int) {
        roomsByStatus[status]?.set(count)
    }

    private fun incrementRoomStatus(status: RoomStatus) {
        roomsByStatus[status]?.incrementAndGet()
    }

    private fun decrementRoomStatus(status: RoomStatus) {
        roomsByStatus[status]?.let { counter ->
            if (counter.get() > 0) counter.decrementAndGet()
        }
    }
}
