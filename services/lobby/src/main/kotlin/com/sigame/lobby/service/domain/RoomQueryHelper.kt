package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.enums.RoomStatus
import com.sigame.lobby.domain.exception.RoomNotFoundException
import com.sigame.lobby.domain.exception.RoomNotFoundByCodeException
import com.sigame.lobby.domain.model.GameRoom
import com.sigame.lobby.domain.model.RoomPlayer
import com.sigame.lobby.domain.model.RoomSettings
import com.sigame.lobby.domain.repository.GameRoomRepository
import com.sigame.lobby.domain.repository.RoomPlayerRepository
import com.sigame.lobby.domain.repository.RoomSettingsRepository
import com.sigame.lobby.grpc.pack.PackInfo
import com.sigame.lobby.service.batch.BatchOperationService
import com.sigame.lobby.service.mapper.RoomMapper
import kotlinx.coroutines.async
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.toList
import kotlinx.coroutines.reactive.asFlow
import kotlinx.coroutines.reactive.awaitFirst
import kotlinx.coroutines.reactive.awaitFirstOrNull
import mu.KotlinLogging
import org.springframework.stereotype.Component
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Component
class RoomQueryHelper(
    private val gameRoomRepository: GameRoomRepository,
    private val roomPlayerRepository: RoomPlayerRepository,
    private val roomSettingsRepository: RoomSettingsRepository,
    private val batchOperationService: BatchOperationService,
    private val roomMapper: RoomMapper
) {

    data class RoomListData(
        val rooms: List<GameRoom>,
        val total: Long,
        val playersByRoom: Map<UUID, List<RoomPlayer>>,
        val packInfoMap: Map<UUID, PackInfo?>
    )

    data class RoomDetailData(
        val room: GameRoom,
        val players: List<RoomPlayer>,
        val settings: RoomSettings?,
        val packName: String?
    )

    suspend fun fetchRoomListData(status: RoomStatus?, hasSlots: Boolean?, size: Int, offset: Int): RoomListData = coroutineScope {
        val roomsD = async { fetchRooms(status, hasSlots, size, offset) }
        val totalD = async { countRooms(status, hasSlots) }

        val rooms = roomsD.await()
        val total = totalD.await()

        if (rooms.isEmpty()) {
            return@coroutineScope RoomListData(emptyList(), 0, emptyMap(), emptyMap())
        }

        val roomIds = rooms.mapNotNull { it.id }
        val packIds = rooms.map { it.packId }

        val playersD = async { loadPlayersByRoomIds(roomIds) }
        val packsD = async { batchOperationService.getPackInfoBatch(packIds) }

        RoomListData(rooms, total, playersD.await(), packsD.await())
    }

    suspend fun fetchRoomById(roomId: UUID): RoomDetailData = coroutineScope {
        val room = gameRoomRepository.findById(roomId).awaitFirstOrNull()
            ?: throw RoomNotFoundException(roomId)
        fetchRoomDetails(room)
    }

    suspend fun fetchRoomByCode(code: String): RoomDetailData = coroutineScope {
        val room = gameRoomRepository.findByRoomCode(code).awaitFirstOrNull()
            ?: throw RoomNotFoundByCodeException(code)
        fetchRoomDetails(room)
    }

    suspend fun findActivePlayerRoom(userId: UUID): UUID? =
        roomPlayerRepository.findActiveByUserId(userId).awaitFirstOrNull()?.roomId

    fun buildRoomDto(data: RoomDetailData): RoomDto =
        roomMapper.toDto(data.room, data.players.size, data.players, data.settings, data.packName)

    fun buildRoomDtoList(data: RoomListData): List<RoomDto> =
        data.rooms.map { room ->
            val players = data.playersByRoom[room.requireId()] ?: emptyList()
            roomMapper.toDtoWithCache(room, players.size, data.packInfoMap[room.packId]?.name, players)
        }

    fun parseRoomStatus(status: String?): RoomStatus? {
        if (status.isNullOrBlank()) return null
        return try {
            RoomStatus.valueOf(status.uppercase())
        } catch (e: IllegalArgumentException) {
            logger.warn { "Invalid room status: '$status'" }
            null
        }
    }

    fun calculateTotalPages(total: Long, size: Int): Int =
        if (size > 0) ((total + size - 1) / size).toInt() else 0

    private suspend fun fetchRoomDetails(room: GameRoom): RoomDetailData = coroutineScope {
        val playersD = async { roomPlayerRepository.findActiveByRoomId(room.requireId()).asFlow().toList() }
        val settingsD = async { roomSettingsRepository.findByRoomId(room.requireId()).awaitFirstOrNull() }
        val packInfoD = async { batchOperationService.getPackInfoBatch(listOf(room.packId))[room.packId] }

        RoomDetailData(room, playersD.await(), settingsD.await(), packInfoD.await()?.name)
    }

    private suspend fun fetchRooms(status: RoomStatus?, hasSlots: Boolean?, size: Int, offset: Int): List<GameRoom> =
        when {
            status != null -> gameRoomRepository.findByStatus(status.name.lowercase(), size, offset)
            hasSlots == true -> gameRoomRepository.findPublicWaitingRooms(size, offset)
            else -> gameRoomRepository.findPublicWaitingRooms(size, offset)
        }.asFlow().toList()

    private suspend fun countRooms(status: RoomStatus?, hasSlots: Boolean?): Long =
        when {
            status != null -> gameRoomRepository.countByStatus(status.name.lowercase()).awaitFirst()
            hasSlots == true -> gameRoomRepository.countPublicWaitingRooms().awaitFirst()
            else -> gameRoomRepository.countPublicWaitingRooms().awaitFirst()
        }

    private suspend fun loadPlayersByRoomIds(roomIds: List<UUID>): Map<UUID, List<RoomPlayer>> =
        roomPlayerRepository.findActiveByRoomIds(roomIds).asFlow().toList().groupBy { it.roomId }
}

