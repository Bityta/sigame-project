package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.CannotKickSelfException
import com.sigame.lobby.domain.exception.InvalidPasswordException
import com.sigame.lobby.domain.exception.InvalidRoomStateException
import com.sigame.lobby.domain.exception.PlayerAlreadyInRoomException
import com.sigame.lobby.domain.exception.PlayerNotInRoomException
import com.sigame.lobby.domain.exception.RoomFullException
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.UnauthorizedRoomActionException
import com.sigame.lobby.domain.exception.UserInfoNotFoundException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.auth.AuthServiceClient
import com.sigame.lobby.grpc.auth.UserInfo
import com.sigame.lobby.grpc.pack.PackServiceClient
import com.sigame.lobby.metrics.LobbyMetrics
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.mapper.RoomMapper
import com.sigame.lobby.sse.event.PlayerJoinedEvent
import com.sigame.lobby.sse.event.PlayerLeftEvent
import com.sigame.lobby.sse.event.RoomClosedEvent
import com.sigame.lobby.sse.service.RoomEventPublisher
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder
import org.springframework.stereotype.Component
import java.time.LocalDateTime
import java.util.UUID

@Component
class RoomMembershipHelper(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val roomCacheService: RoomCacheService,
    private val authServiceClient: AuthServiceClient,
    private val packServiceClient: PackServiceClient,
    private val lobbyMetrics: LobbyMetrics,
    private val roomMapper: RoomMapper,
    private val roomEventPublisher: RoomEventPublisher,
    private val passwordEncoder: BCryptPasswordEncoder
) {

    data class JoinContext(
        val room: GameRoom,
        val currentPlayers: Int,
        val existingPlayer: RoomPlayer?,
        val currentPlayer: RoomPlayer?,
        val userInfo: UserInfo,
        val players: List<RoomPlayer>,
        val settings: RoomSettings?,
        val packName: String?
    )

    data class LeaveContext(
        val room: GameRoom,
        val player: RoomPlayer,
        val currentPlayers: Int
    )

    suspend fun fetchJoinContext(roomId: UUID, userId: UUID): JoinContext = coroutineScope {
        val roomD = async { findRoomOrThrow(roomId) }
        val countD = async { countActivePlayers(roomId) }
        val existingD = async { findExistingPlayer(roomId, userId) }
        val currentD = async { findActivePlayerInAnyRoom(userId) }
        val userInfoD = async { fetchRequiredUserInfo(userId) }
        val playersD = async { getActivePlayers(roomId) }
        val settingsD = async { findSettingsByRoomId(roomId) }

        val room = roomD.await()
        val packNameD = async { packServiceClient.getPackInfo(room.packId)?.name }

        JoinContext(
            room = room,
            currentPlayers = countD.await(),
            existingPlayer = existingD.await(),
            currentPlayer = currentD.await(),
            userInfo = userInfoD.await(),
            players = playersD.await(),
            settings = settingsD.await(),
            packName = packNameD.await()
        )
    }

    suspend fun fetchLeaveContext(roomId: UUID, userId: UUID): LeaveContext = coroutineScope {
        val roomD = async { findRoomOrThrow(roomId) }
        val playerD = async { findPlayerOrThrow(roomId, userId) }
        val countD = async { countActivePlayers(roomId) }
        LeaveContext(roomD.await(), playerD.await(), countD.await())
    }

    fun validateJoin(ctx: JoinContext, roomId: UUID, password: String?) {
        validateRoomStatus(ctx.room, RoomStatus.WAITING, "join")
        validatePassword(ctx.room, password)
        validateCapacity(ctx.currentPlayers, ctx.room.maxPlayers, roomId)
        validatePlayerCanJoin(ctx.existingPlayer, ctx.currentPlayer, roomId)
    }

    fun validateKick(room: GameRoom, hostId: UUID, targetUserId: UUID, targetPlayer: RoomPlayer, roomId: UUID) {
        validateHostAccess(room, hostId, "kick player")
        validateRoomStatus(room, RoomStatus.WAITING, "kick player")
        if (targetUserId == hostId) throw CannotKickSelfException(hostId)
        if (targetPlayer.leftAt != null) throw PlayerNotInRoomException(targetUserId, roomId)
    }

    fun validateTransfer(
        room: GameRoom,
        currentHostId: UUID,
        newHostId: UUID,
        newHostPlayer: RoomPlayer,
        roomId: UUID
    ) {
        validateHostAccess(room, currentHostId, "transfer host")
        validateRoomStatus(room, RoomStatus.WAITING, "transfer host")
        if (newHostId == currentHostId) throw IllegalArgumentException("Cannot transfer host to yourself")
        if (newHostPlayer.leftAt != null) throw PlayerNotInRoomException(newHostId, roomId)
    }

    suspend fun addPlayer(roomId: UUID, userId: UUID, ctx: JoinContext, role: String): RoomPlayer {
        return if (ctx.existingPlayer != null) {
            val rejoined = ctx.existingPlayer.copy(
                leftAt = null,
                joinedAt = LocalDateTime.now(),
                username = ctx.userInfo.username,
                avatarUrl = ctx.userInfo.avatarUrl
            )
            roomPlayerRepository.save(rejoined).awaitFirst()
        } else {
            val newPlayer = RoomPlayer(
                roomId = roomId,
                userId = userId,
                username = ctx.userInfo.username,
                avatarUrl = ctx.userInfo.avatarUrl,
                role = role
            )
            roomPlayerRepository.save(newPlayer).awaitFirst()
        }
    }

    fun buildJoinResponse(ctx: JoinContext, savedPlayer: RoomPlayer, newCount: Int): RoomDto {
        val isNewPlayer = ctx.existingPlayer == null || ctx.existingPlayer.leftAt != null
        val players = if (isNewPlayer) ctx.players + savedPlayer else ctx.players
        return roomMapper.toDto(ctx.room, newCount, players, ctx.settings, ctx.packName)
    }

    fun buildExistingPlayerResponse(ctx: JoinContext, existingPlayer: RoomPlayer): RoomDto {
        return roomMapper.toDto(ctx.room, ctx.currentPlayers, ctx.players, ctx.settings, ctx.packName)
    }

    suspend fun onPlayerJoined(userId: UUID, room: GameRoom, userInfo: UserInfo, newCount: Int) = coroutineScope {
        launch { roomCacheService.setUserCurrentRoom(userId, room.requireId()) }
        launch { roomCacheService.addRoomPlayer(room.requireId(), userId) }
        launch { roomCacheService.cacheRoomData(room, newCount) }

        roomEventPublisher.publish(
            PlayerJoinedEvent(room.requireId(), userId, userInfo.username, newCount)
        )
        lobbyMetrics.recordPlayerJoined()
    }

    suspend fun handleLeave(room: GameRoom, player: RoomPlayer, currentPlayers: Int) = coroutineScope {
        removePlayer(player)
        launch { roomCacheService.deleteUserCurrentRoom(player.userId) }
        launch { roomCacheService.removeRoomPlayer(room.requireId(), player.userId) }
        lobbyMetrics.recordPlayerLeft()

        val newCount = currentPlayers - 1
        if (player.getRoleEnum() == PlayerRole.HOST) {
            handleHostLeave(room, player, newCount)
        } else {
            handleRegularLeave(room, player, newCount)
        }
    }

    suspend fun onPlayerKicked(room: GameRoom, player: RoomPlayer, currentPlayers: Int) = coroutineScope {
        removePlayer(player)
        launch { roomCacheService.deleteUserCurrentRoom(player.userId) }
        launch { roomCacheService.removeRoomPlayer(room.requireId(), player.userId) }
        lobbyMetrics.recordPlayerLeft()

        val newCount = currentPlayers - 1
        launch { roomCacheService.cacheRoomData(room, newCount) }

        roomEventPublisher.publish(
            PlayerLeftEvent(room.requireId(), player.userId, player.username, "kicked", newCount)
        )
    }

    suspend fun onHostTransferred(room: GameRoom, fromHostId: UUID, toPlayer: RoomPlayer, currentPlayers: Int) =
        coroutineScope {
            val updatedRoom = transferHost(room, fromHostId, toPlayer)
            launch { roomCacheService.cacheRoomData(updatedRoom, currentPlayers) }
        }

    private suspend fun handleHostLeave(room: GameRoom, leftPlayer: RoomPlayer, newCount: Int) = coroutineScope {
        if (newCount == 0) {
            cancelRoom(room)
        } else {
            val activePlayers = getActivePlayers(room.requireId())
            val newHost = activePlayers.first()
            val updatedRoom = transferHost(room, leftPlayer.userId, newHost)

            launch { roomCacheService.cacheRoomData(updatedRoom, newCount) }
            roomEventPublisher.publish(
                PlayerLeftEvent(room.requireId(), leftPlayer.userId, leftPlayer.username, "left", newCount)
            )
        }
    }

    private suspend fun handleRegularLeave(room: GameRoom, leftPlayer: RoomPlayer, newCount: Int) = coroutineScope {
        if (newCount == 0) {
            cancelRoom(room)
        } else {
            launch { roomCacheService.cacheRoomData(room, newCount) }
            roomEventPublisher.publish(
                PlayerLeftEvent(room.requireId(), leftPlayer.userId, leftPlayer.username, "left", newCount)
            )
        }
    }

    private suspend fun cancelRoom(room: GameRoom) = coroutineScope {
        val cancelled = room.copy(status = GameRoom.statusFromEnum(RoomStatus.CANCELLED))
        gameRoomRepository.save(cancelled).awaitFirstOrNull()
        lobbyMetrics.recordRoomCancelled(room.getStatusEnum())
        launch { roomCacheService.clearRoomCache(room.requireId()) }

        roomEventPublisher.publish(RoomClosedEvent(room.requireId(), "no_players"))
        roomEventPublisher.closeRoom(room.requireId())
    }

    private suspend fun removePlayer(player: RoomPlayer): RoomPlayer {
        val updated = player.copy(leftAt = LocalDateTime.now())
        return roomPlayerRepository.save(updated).awaitFirst()
    }

    private suspend fun transferHost(room: GameRoom, fromHostId: UUID, toPlayer: RoomPlayer): GameRoom = coroutineScope {
        val currentHost = roomPlayerRepository.findByRoomIdAndUserId(room.requireId(), fromHostId).awaitFirstOrNull()
        
        val demoteJob = if (currentHost != null) {
            launch { 
                roomPlayerRepository.save(currentHost.copy(role = RoomPlayer.roleFromEnum(PlayerRole.PLAYER)))
                    .awaitFirstOrNull() 
            }
        } else null
        
        launch { 
            roomPlayerRepository.save(toPlayer.copy(role = RoomPlayer.roleFromEnum(PlayerRole.HOST)))
                .awaitFirstOrNull() 
        }
        
        demoteJob?.join()
        
        val updatedRoom = room.copy(hostId = toPlayer.userId)
        gameRoomRepository.save(updatedRoom).awaitFirst()
    }

    private suspend fun findRoomOrThrow(roomId: UUID): GameRoom =
        gameRoomRepository.findById(roomId).awaitFirstOrNull() ?: throw RoomNotFoundException(roomId)

    private suspend fun findPlayerOrThrow(roomId: UUID, userId: UUID): RoomPlayer =
        roomPlayerRepository.findByRoomIdAndUserId(roomId, userId).awaitFirstOrNull()
            ?: throw PlayerNotInRoomException(userId, roomId)

    private suspend fun findExistingPlayer(roomId: UUID, userId: UUID): RoomPlayer? =
        roomPlayerRepository.findByRoomIdAndUserId(roomId, userId).awaitFirstOrNull()

    private suspend fun findActivePlayerInAnyRoom(userId: UUID): RoomPlayer? =
        roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()

    private suspend fun fetchRequiredUserInfo(userId: UUID): UserInfo =
        authServiceClient.getUserInfo(userId) ?: throw UserInfoNotFoundException(userId)

    private suspend fun countActivePlayers(roomId: UUID): Int =
        roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()

    private suspend fun getActivePlayers(roomId: UUID): List<RoomPlayer> =
        roomPlayerRepository.findActiveByRoomId(roomId).asFlow().toList()

    private suspend fun findSettingsByRoomId(roomId: UUID): RoomSettings? =
        roomSettingsRepository.findByRoomId(roomId).awaitFirstOrNull()

    private fun validateRoomStatus(room: GameRoom, expected: RoomStatus, action: String) {
        if (room.getStatusEnum() != expected) {
            throw InvalidRoomStateException(room.requireId(), room.getStatusEnum(), action)
        }
    }

    private fun validatePassword(room: GameRoom, password: String?) {
        if (room.passwordHash != null) {
            if (password == null || !passwordEncoder.matches(password, room.passwordHash)) {
                throw InvalidPasswordException()
            }
        }
    }

    private fun validateCapacity(current: Int, max: Int, roomId: UUID) {
        if (current >= max) throw RoomFullException(roomId)
    }

    private fun validatePlayerCanJoin(existing: RoomPlayer?, current: RoomPlayer?, roomId: UUID) {
        if (existing != null && existing.leftAt == null) {
            throw PlayerAlreadyInRoomException(existing.userId, existing.roomId)
        }
        if (current != null && current.roomId != roomId) {
            throw PlayerAlreadyInRoomException(current.userId, current.roomId)
        }
    }

    private fun validateHostAccess(room: GameRoom, userId: UUID, action: String) {
        if (room.hostId != userId) {
            throw UnauthorizedRoomActionException(userId, action, "Only the host can $action")
        }
    }
}

