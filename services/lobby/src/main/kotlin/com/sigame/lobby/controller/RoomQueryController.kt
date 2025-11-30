package com.sigame.lobby.controller

import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomListResponse
import com.sigame.lobby.service.domain.RoomQueryService
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
@RequestMapping(ApiRoutes.BASE)
class RoomQueryController(
    private val roomQueryService: RoomQueryService
) {

    @GetMapping(ApiRoutes.HEALTH)
    suspend fun health() = ResponseEntity.ok(mapOf("status" to "UP"))

    @GetMapping(ApiRoutes.Rooms.LIST)
    suspend fun getRooms(
        @RequestParam(defaultValue = "0") page: Int,
        @RequestParam(defaultValue = "20") size: Int,
        @RequestParam(required = false) status: String?,
        @RequestParam(name = "has_slots", required = false) hasSlots: Boolean?
    ): ResponseEntity<RoomListResponse> {
        return ResponseEntity.ok(roomQueryService.getRooms(page, size, status, hasSlots))
    }

    @GetMapping(ApiRoutes.Rooms.BY_ID)
    suspend fun getRoomById(@PathVariable id: UUID): ResponseEntity<RoomDto> {
        return ResponseEntity.ok(roomQueryService.getRoomById(id))
    }

    @GetMapping(ApiRoutes.Rooms.BY_CODE)
    suspend fun getRoomByCode(@PathVariable code: String): ResponseEntity<RoomDto> {
        return ResponseEntity.ok(roomQueryService.getRoomByCode(code))
    }
}

