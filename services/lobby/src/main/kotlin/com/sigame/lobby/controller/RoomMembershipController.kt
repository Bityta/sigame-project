package com.sigame.lobby.controller

import com.sigame.lobby.domain.dto.JoinRoomRequest
import com.sigame.lobby.domain.dto.KickPlayerRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.TransferHostRequest
import com.sigame.lobby.security.AuthenticatedUser
import com.sigame.lobby.security.CurrentUser
import com.sigame.lobby.service.domain.RoomMembershipService
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import java.util.UUID

@RestController
@RequestMapping(ApiRoutes.BASE)
class RoomMembershipController(
    private val roomMembershipService: RoomMembershipService
) {

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

    @PostMapping(ApiRoutes.Rooms.KICK)
    suspend fun kickPlayer(
        @PathVariable id: UUID,
        @RequestBody request: KickPlayerRequest,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        roomMembershipService.kickPlayer(id, user.userId, request.userId)
        return ResponseEntity.noContent().build()
    }

    @PostMapping(ApiRoutes.Rooms.TRANSFER_HOST)
    suspend fun transferHost(
        @PathVariable id: UUID,
        @RequestBody request: TransferHostRequest,
        @CurrentUser user: AuthenticatedUser
    ): ResponseEntity<Void> {
        roomMembershipService.transferHostManually(id, user.userId, request.newHostId)
        return ResponseEntity.noContent().build()
    }
}

