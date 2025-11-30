package com.sigame.lobby.controller

import com.sigame.lobby.domain.dto.*
import com.sigame.lobby.security.AuthenticatedUser
import com.sigame.lobby.security.CurrentUser
import com.sigame.lobby.service.domain.RoomLifecycleService
import com.sigame.lobby.service.domain.RoomMembershipService
import com.sigame.lobby.service.domain.RoomQueryService
import jakarta.validation.Valid
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
@RequestMapping(ApiRoutes.BASE)
class LobbyController(
    private val roomQueryService: RoomQueryService,
    private val roomLifecycleService: RoomLifecycleService,
    private val roomMembershipService: RoomMembershipService
) {
    
    @GetMapping(ApiRoutes.HEALTH)
    suspend fun health() = ResponseEntity.ok(mapOf("status" to "UP"))
    
    @PostMapping(ApiRoutes.Rooms.LIST)
    suspend fun createRoom(
        @Valid @RequestBody request: CreateRoomRequest,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        val room = roomLifecycleService.createRoom(user.userId, request)
        return ResponseEntity.status(HttpStatus.CREATED).body(room)
    }
    
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
    
    @PostMapping(ApiRoutes.Rooms.JOIN)
    suspend fun joinRoom(
        @PathVariable id: UUID,
        @RequestBody(required = false) request: JoinRoomRequest = JoinRoomRequest(),
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<RoomDto> {
        return ResponseEntity.ok(roomMembershipService.joinRoom(id, user.userId, request))
    }
    
    @DeleteMapping(ApiRoutes.Rooms.LEAVE)
    suspend fun leaveRoom(
        @PathVariable id: UUID,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        roomMembershipService.leaveRoom(id, user.userId)
        return ResponseEntity.noContent().build()
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
