package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.JoinRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import mu.KotlinLogging
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class RoomMembershipService(private val helper: RoomMembershipHelper) {

    @Transactional
    suspend fun joinRoom(roomId: UUID, userId: UUID, request: JoinRoomRequest): RoomDto {
        logger.info { "User $userId joining room $roomId" }

        val ctx = helper.fetchJoinContext(roomId, userId)
        
        val existingInThisRoom = ctx.existingPlayer?.takeIf { it.leftAt == null }
        if (existingInThisRoom != null) {
            logger.info { "User $userId already in room $roomId, returning current state" }
            return helper.buildExistingPlayerResponse(ctx)
        }
        
        helper.validateJoin(ctx, roomId, request.password)

        val savedPlayer = helper.addPlayer(roomId, userId, ctx, request.role.lowercase())
        val newCount = ctx.currentPlayers + 1

        helper.onPlayerJoined(userId, ctx.room, ctx.userInfo, newCount)

        return helper.buildJoinResponse(ctx, savedPlayer, newCount)
    }

    @Transactional
    suspend fun leaveRoom(roomId: UUID, userId: UUID) {
        logger.info { "User $userId leaving room $roomId" }

        val ctx = helper.fetchLeaveContext(roomId, userId)
        require(ctx.player.leftAt == null) { "Player has already left" }

        helper.handleLeave(ctx.room, ctx.player, ctx.currentPlayers)
    }

    @Transactional
    suspend fun kickPlayer(roomId: UUID, hostId: UUID, targetUserId: UUID) {
        logger.info { "Host $hostId kicking player $targetUserId from room $roomId" }

        val ctx = helper.fetchLeaveContext(roomId, targetUserId)
        helper.validateKick(ctx.room, hostId, targetUserId, ctx.player, roomId)

        helper.onPlayerKicked(ctx.room, ctx.player, ctx.currentPlayers)
    }

    @Transactional
    suspend fun transferHostManually(roomId: UUID, currentHostId: UUID, newHostId: UUID) {
        logger.info { "Host $currentHostId transferring host role to $newHostId in room $roomId" }

        val ctx = helper.fetchLeaveContext(roomId, newHostId)
        helper.validateTransfer(ctx.room, currentHostId, newHostId, ctx.player, roomId)

        helper.onHostTransferred(ctx.room, currentHostId, ctx.player)
    }
}
