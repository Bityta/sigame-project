package com.sigame.lobby.metrics

import com.sigame.lobby.domain.enums.RoomStatus
import io.micrometer.core.instrument.Counter
import io.micrometer.core.instrument.Gauge
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import mu.KotlinLogging
import org.springframework.stereotype.Component
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.atomic.AtomicInteger

private val logger = KotlinLogging.logger {}

@Component
class LobbyMetrics(
    private val meterRegistry: MeterRegistry
) {

    private val roomsCreatedTotal: Counter
    private val roomsStartedTotal: Counter
    private val roomsCancelledTotal: Counter
    private val playersJoinedTotal: Counter
    private val playersLeftTotal: Counter
    private val kafkaEventsPublishedTotal: Counter
    private val grpcCallsTotal: Counter
    private val grpcErrorsTotal: Counter

    private val roomCreationTimer: Timer
    private val roomStartTimer: Timer
    private val playerJoinTimer: Timer
    private val grpcCallTimer: Timer

    private val activeRooms = AtomicInteger(0)
    private val activePlayers = AtomicInteger(0)
    private val roomsByStatus = ConcurrentHashMap<RoomStatus, AtomicInteger>()

    init {
        logger.info { "Initializing Lobby Metrics" }

        roomsCreatedTotal = Counter.builder("lobby_rooms_created_total")
            .description("Total number of rooms created")
            .register(meterRegistry)

        roomsStartedTotal = Counter.builder("lobby_rooms_started_total")
            .description("Total number of rooms started")
            .register(meterRegistry)

        roomsCancelledTotal = Counter.builder("lobby_rooms_cancelled_total")
            .description("Total number of rooms cancelled")
            .register(meterRegistry)

        playersJoinedTotal = Counter.builder("lobby_players_joined_total")
            .description("Total number of players joined")
            .register(meterRegistry)

        playersLeftTotal = Counter.builder("lobby_players_left_total")
            .description("Total number of players left")
            .register(meterRegistry)

        kafkaEventsPublishedTotal = Counter.builder("lobby_kafka_events_published_total")
            .description("Total number of Kafka events published")
            .register(meterRegistry)

        grpcCallsTotal = Counter.builder("lobby_grpc_calls_total")
            .description("Total number of gRPC calls")
            .register(meterRegistry)

        grpcErrorsTotal = Counter.builder("lobby_grpc_errors_total")
            .description("Total number of gRPC errors")
            .register(meterRegistry)

        roomCreationTimer = Timer.builder("lobby_room_creation_duration")
            .description("Time taken to create a room")
            .register(meterRegistry)

        roomStartTimer = Timer.builder("lobby_room_start_duration")
            .description("Time taken to start a room")
            .register(meterRegistry)

        playerJoinTimer = Timer.builder("lobby_player_join_duration")
            .description("Time taken for a player to join")
            .register(meterRegistry)

        grpcCallTimer = Timer.builder("lobby_grpc_call_duration")
            .description("Time taken for gRPC calls")
            .register(meterRegistry)

        Gauge.builder("lobby_active_rooms", activeRooms) { it.get().toDouble() }
            .description("Current number of active rooms")
            .register(meterRegistry)

        Gauge.builder("lobby_active_players", activePlayers) { it.get().toDouble() }
            .description("Current number of active players in rooms")
            .register(meterRegistry)

        RoomStatus.entries.forEach { status ->
            roomsByStatus[status] = AtomicInteger(0)
            Gauge.builder("lobby_rooms_by_status", roomsByStatus[status]!!) { it.get().toDouble() }
                .description("Number of rooms by status")
                .tag("status", status.name)
                .register(meterRegistry)
        }

        logger.info { "Lobby Metrics initialized successfully" }
    }

    fun recordRoomCreated() {
        roomsCreatedTotal.increment()
        activeRooms.incrementAndGet()
        incrementRoomStatus(RoomStatus.WAITING)
        logger.debug { "Recorded room created. Active rooms: ${activeRooms.get()}" }
    }

    fun recordRoomStarted() {
        roomsStartedTotal.increment()
        decrementRoomStatus(RoomStatus.WAITING)
        incrementRoomStatus(RoomStatus.STARTING)
        logger.debug { "Recorded room started" }
    }

    fun recordRoomCancelled(oldStatus: RoomStatus) {
        roomsCancelledTotal.increment()
        activeRooms.decrementAndGet()
        decrementRoomStatus(oldStatus)
        logger.debug { "Recorded room cancelled. Active rooms: ${activeRooms.get()}" }
    }

    fun recordPlayerJoined() {
        playersJoinedTotal.increment()
        activePlayers.incrementAndGet()
        logger.debug { "Recorded player joined. Active players: ${activePlayers.get()}" }
    }

    fun recordPlayerLeft() {
        playersLeftTotal.increment()
        activePlayers.decrementAndGet()
        logger.debug { "Recorded player left. Active players: ${activePlayers.get()}" }
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
        logger.debug { "Set active rooms to $count" }
    }

    fun setActivePlayers(count: Int) {
        activePlayers.set(count)
        logger.debug { "Set active players to $count" }
    }

    fun setRoomsByStatus(status: RoomStatus, count: Int) {
        roomsByStatus[status]?.set(count)
        logger.debug { "Set rooms for status $status to $count" }
    }

    private fun incrementRoomStatus(status: RoomStatus) {
        roomsByStatus[status]?.incrementAndGet()
    }

    private fun decrementRoomStatus(status: RoomStatus) {
        roomsByStatus[status]?.let { counter ->
            if (counter.get() > 0) {
                counter.decrementAndGet()
            }
        }
    }

    fun recordKafkaEventPublished() {
        kafkaEventsPublishedTotal.increment()
        logger.debug { "Recorded Kafka event published" }
    }
}
