package com.sigame.lobby.sse.service

import com.sigame.lobby.sse.event.RoomEvent
import mu.KotlinLogging
import org.springframework.stereotype.Service
import reactor.core.publisher.Flux
import reactor.core.publisher.Sinks
import java.util.UUID
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.atomic.AtomicLong

private val logger = KotlinLogging.logger {}

@Service
class RoomEventPublisher {

    private val roomSinks = ConcurrentHashMap<UUID, Sinks.Many<RoomEvent>>()
    private val publishedCount = AtomicLong(0)
    private val droppedCount = AtomicLong(0)
    private val failedCount = AtomicLong(0)

    private fun createSink(): Sinks.Many<RoomEvent> =
        Sinks.many().multicast().onBackpressureBuffer(100)

    fun subscribe(roomId: UUID): Flux<RoomEvent> {
        val sink = roomSinks.compute(roomId) { _, existing ->
            if (existing == null || existing.currentSubscriberCount() == 0) {
                createSink()
            } else {
                existing
            }
        }!!
        return sink.asFlux()
    }

    fun publish(event: RoomEvent) {
        val sink = roomSinks[event.roomId]
        if (sink == null) {
            droppedCount.incrementAndGet()
            return
        }

        var result = sink.tryEmitNext(event)
        
        if (result == Sinks.EmitResult.FAIL_CANCELLED || result == Sinks.EmitResult.FAIL_TERMINATED) {
            val newSink = createSink()
            roomSinks[event.roomId] = newSink
            result = newSink.tryEmitNext(event)
        }
        
        if (result.isSuccess) {
            publishedCount.incrementAndGet()
        } else {
            failedCount.incrementAndGet()
            logger.warn { "SSE emit failed: ${event.type} room=${event.roomId} result=$result" }
        }
    }

    fun closeRoom(roomId: UUID) {
        roomSinks.remove(roomId)?.tryEmitComplete()
    }

    fun getActiveRoomsCount(): Int = roomSinks.size
    fun getPublishedCount(): Long = publishedCount.get()
    fun getDroppedCount(): Long = droppedCount.get()
    fun getFailedCount(): Long = failedCount.get()
}

