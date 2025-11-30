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
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds

private val logger = KotlinLogging.logger {}

@Component
class LobbyMetrics(
    private val meterRegistry: MeterRegistry
) {
    
    // Counters
    private val roomsCreatedTotal: Counter
    private val roomsStartedTotal: Counter
    private val roomsCancelledTotal: Counter
    private val playersJoinedTotal: Counter
    private val playersLeftTotal: Counter
    private val kafkaEventsPublishedTotal: Counter
    private val grpcCallsTotal: Counter
    private val grpcErrorsTotal: Counter
    
    // Timers
    private val roomCreationTimer: Timer
    private val roomStartTimer: Timer
    private val playerJoinTimer: Timer
    private val grpcCallTimer: Timer
    
    // Gauges - using AtomicInteger for thread-safe operations
    private val activeRooms = AtomicInteger(0)
    private val activePlayers = AtomicInteger(0)
    private val roomsByStatus = ConcurrentHashMap<RoomStatus, AtomicInteger>()
    
    init {
        logger.info { "Initializing Lobby Metrics" }
        
        // Initialize counters
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
        
        // Initialize timers
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
        
        // Initialize gauges
        Gauge.builder("lobby_active_rooms", activeRooms) { it.get().toDouble() }
            .description("Current number of active rooms")
            .register(meterRegistry)
        
        Gauge.builder("lobby_active_players", activePlayers) { it.get().toDouble() }
            .description("Current number of active players in rooms")
            .register(meterRegistry)
        
        // Initialize room status gauges
        RoomStatus.values().forEach { status ->
            roomsByStatus[status] = AtomicInteger(0)
            Gauge.builder("lobby_rooms_by_status", roomsByStatus[status]!!) { it.get().toDouble() }
                .description("Number of rooms by status")
                .tag("status", status.name)
                .register(meterRegistry)
        }
        
        logger.info { "Lobby Metrics initialized successfully" }
    }
    
    // Room lifecycle methods
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
    
    fun recordRoomStatusChange(oldStatus: RoomStatus, newStatus: RoomStatus) {
        decrementRoomStatus(oldStatus)
        incrementRoomStatus(newStatus)
        logger.debug { "Room status changed from $oldStatus to $newStatus" }
    }
    
    fun recordRoomCancelled(oldStatus: RoomStatus) {
        roomsCancelledTotal.increment()
        activeRooms.decrementAndGet()
        decrementRoomStatus(oldStatus)
        logger.debug { "Recorded room cancelled. Active rooms: ${activeRooms.get()}" }
    }
    
    fun recordRoomFinished(oldStatus: RoomStatus) {
        activeRooms.decrementAndGet()
        decrementRoomStatus(oldStatus)
        logger.debug { "Recorded room finished. Active rooms: ${activeRooms.get()}" }
    }
    
    // Player methods
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
    
    // Timer methods
    fun <T> recordRoomCreationTime(block: () -> T): T {
        return roomCreationTimer.record<T> { block() }!!
    }
    
    fun <T> recordRoomStartTime(block: () -> T): T {
        return roomStartTimer.record<T> { block() }!!
    }
    
    fun <T> recordPlayerJoinTime(block: () -> T): T {
        return playerJoinTimer.record<T> { block() }!!
    }
    
    fun <T> recordGrpcCallTime(serviceName: String, block: () -> T): T {
        return Timer.builder("lobby_grpc_call_duration")
            .tag("service", serviceName)
            .register(meterRegistry)
            .record<T> { block() }!!
    }
    
    // gRPC metrics
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
    
    // Bulk set methods for initialization/sync
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
    
    // Helper methods
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
    
    // Get current values (for logging/debugging)
    fun getActiveRooms(): Int = activeRooms.get()
    fun getActivePlayers(): Int = activePlayers.get()
    fun getRoomsByStatus(status: RoomStatus): Int = roomsByStatus[status]?.get() ?: 0
    
    // Kafka events
    fun recordKafkaEventPublished() {
        kafkaEventsPublishedTotal.increment()
        logger.debug { "Recorded Kafka event published" }
    }
}

