package com.sigame.lobby.controller

import com.fasterxml.jackson.databind.ObjectMapper
import com.sigame.lobby.sse.service.RoomEventPublisher
import org.springframework.http.MediaType
import org.springframework.http.codec.ServerSentEvent
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import reactor.core.publisher.Flux
import java.time.Duration
import java.util.UUID

@RestController
@RequestMapping(ApiRoutes.BASE)
class RoomEventsController(
    private val roomEventPublisher: RoomEventPublisher,
    private val objectMapper: ObjectMapper
) {

    companion object {
        private val HEARTBEAT_INTERVAL = Duration.ofSeconds(30)
    }

    @GetMapping(ApiRoutes.Rooms.EVENTS, produces = [MediaType.TEXT_EVENT_STREAM_VALUE])
    fun subscribeToRoomEvents(@PathVariable roomId: UUID): Flux<ServerSentEvent<String>> {
        val heartbeat = Flux.interval(HEARTBEAT_INTERVAL)
            .map { ServerSentEvent.builder<String>().comment("ping").build() }

        val events = roomEventPublisher.subscribe(roomId)
            .map { event ->
                ServerSentEvent.builder<String>()
                    .event(event.type)
                    .data(objectMapper.writeValueAsString(event))
                    .build()
            }

        return Flux.merge(heartbeat, events)
    }
}
