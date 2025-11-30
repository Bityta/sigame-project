package com.sigame.lobby.controller

import com.sigame.lobby.domain.dto.*
import com.sigame.lobby.security.AuthenticatedUser
import com.sigame.lobby.security.CurrentUser
import com.sigame.lobby.service.domain.RoomLifecycleService
import jakarta.validation.Valid
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
@RequestMapping(ApiRoutes.BASE)
class RoomLifecycleController(
    private val roomLifecycleService: RoomLifecycleService
) {

    @PostMapping(ApiRoutes.Rooms.LIST)
    suspend fun createRoom(
        @Valid @RequestBody request: CreateRoomRequest,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        val room = roomLifecycleService.createRoom(user.userId, request)
        return ResponseEntity.status(HttpStatus.CREATED).body(room)
    }

    @PostMapping(ApiRoutes.Rooms.START)
    suspend fun startRoom(
        @PathVariable id: UUID,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<StartGameResponse> {
        return ResponseEntity.ok(roomLifecycleService.startRoom(id, user.userId))
    }

    @PatchMapping(ApiRoutes.Rooms.SETTINGS)
    suspend fun updateRoomSettings(
        @PathVariable id: UUID,
        @Valid @RequestBody request: UpdateRoomSettingsRequest,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<UpdateRoomSettingsResponse> {
        return ResponseEntity.ok(roomLifecycleService.updateRoomSettings(id, user.userId, request))
    }

    @DeleteMapping(ApiRoutes.Rooms.BY_ID)
    suspend fun deleteRoom(
        @PathVariable id: UUID,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        roomLifecycleService.deleteRoom(id, user.userId)
        return ResponseEntity.noContent().build()
    }
}

