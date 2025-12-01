package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.CreateRoomRequest
import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.StartGameResponse
import com.sigame.lobby.domain.dto.UpdateRoomSettingsRequest
import com.sigame.lobby.domain.dto.UpdateRoomSettingsResponse
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.service.cache.RoomCacheService
import com.sigame.lobby.service.mapper.RoomMapper
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import kotlinx.coroutines.reactive.awaitFirst
import mu.KotlinLogging
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class RoomLifecycleService(
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomCacheService: RoomCacheService,
    private val roomMapper: RoomMapper,
    private val helper: RoomLifecycleHelper
) {

    @Transactional
    suspend fun createRoom(hostId: UUID, request: CreateRoomRequest): RoomDto = coroutineScope {
        logger.info { "Creating room for host $hostId with pack ${request.packId}" }

        val packValidation = async { helper.validatePackForRoom(request.packId, hostId) }
        val userValidation = async { helper.validateUserNotInRoom(hostId) }
        val hostInfoDeferred = async { helper.fetchRequiredUserInfo(hostId) }
        val packInfoDeferred = async { helper.fetchRequiredPackInfo(request.packId) }

        packValidation.await()
        userValidation.await()
        val hostInfo = hostInfoDeferred.await()
        val packInfo = packInfoDeferred.await()

        val savedRoom = helper.createGameRoom(hostId, request)
        val hostPlayer = awaitAll(
            async {
                helper.createRoomSettings(
                    roomId = savedRoom.requireId(),
                    settings = request.settings
                )
            },
            async {
                helper.createHostPlayer(
                    roomId = savedRoom.requireId(),
                    hostId = hostId,
                    username = hostInfo.username,
                    avatarUrl = hostInfo.avatarUrl
                )
            }
        ).last() as RoomPlayer

        launch { helper.updateCacheForNewRoom(savedRoom, hostId) }
        helper.recordRoomCreatedMetrics()

        roomMapper.toDto(
            room = savedRoom,
            currentPlayers = 1,
            players = listOf(hostPlayer),
            packName = packInfo.name
        )
    }

    @Transactional
    suspend fun startRoom(roomId: UUID, userId: UUID): StartGameResponse = coroutineScope {
        logger.info { "Starting room $roomId by user $userId" }

        val roomDeferred = async { helper.findRoomOrThrow(roomId) }
        val playersDeferred = async { helper.getActivePlayersWithMinCheck(roomId) }

        val room = roomDeferred.await()
        helper.validateHostAccess(room, userId, "start room")
        helper.validateRoomInStatus(room, RoomStatus.WAITING, "start")

        val activePlayers = playersDeferred.await()

        val savedRoomDeferred = async { helper.updateRoomToStarting(room, activePlayers.size) }
        val gameSettingsDeferred = async { helper.buildGameSettings(roomId) }

        val savedRoom = savedRoomDeferred.await()
        val gameSettings = gameSettingsDeferred.await()

        helper.executeGameStart(
            roomId = roomId,
            room = room,
            savedRoom = savedRoom,
            activePlayers = activePlayers,
            gameSettings = gameSettings
        )
    }

    @Transactional
    suspend fun updateRoomSettings(
        roomId: UUID,
        userId: UUID,
        request: UpdateRoomSettingsRequest
    ): UpdateRoomSettingsResponse = coroutineScope {
        logger.info { "Updating room $roomId settings by user $userId" }

        val roomDeferred = async { helper.findRoomOrThrow(roomId) }
        val settingsDeferred = async { helper.findSettingsByRoomId(roomId) }

        val room = roomDeferred.await()
        helper.validateHostAccess(room, userId, "update settings")
        helper.validateRoomInStatus(room, RoomStatus.WAITING, "update settings")

        val currentSettings = settingsDeferred.await()
        val updatedSettings = helper.mergeSettings(roomId, currentSettings, request)

        helper.saveSettings(updatedSettings, isNew = currentSettings == null)

        val currentPlayers = roomPlayerRepository.countActiveByRoomId(roomId).awaitFirst().toInt()
        launch { roomCacheService.cacheRoomData(room, currentPlayers) }

        UpdateRoomSettingsResponse(settings = helper.toSettingsDto(updatedSettings))
    }

    @Transactional
    suspend fun deleteRoom(roomId: UUID, userId: UUID): Unit = coroutineScope {
        logger.info { "Deleting room $roomId by user $userId" }

        val roomDeferred = async { helper.findRoomOrThrow(roomId) }
        val playersDeferred = async { helper.getActivePlayers(roomId) }

        val room = roomDeferred.await()
        helper.validateHostAccess(room, userId, "delete room")

        val players = playersDeferred.await()

        helper.cancelRoom(room)
        helper.recordRoomCancelledMetrics(room, players.size)

        launch { helper.clearRoomAndPlayersCache(roomId, room.roomCode, players) }
    }
}
