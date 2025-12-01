package com.sigame.lobby.service.domain

import com.sigame.lobby.domain.dto.RoomDto
import com.sigame.lobby.domain.dto.RoomListResponse
import com.sigame.lobby.domain.exception.RoomNotFoundException
import mu.KotlinLogging
import org.springframework.stereotype.Service
import java.util.UUID

private val logger = KotlinLogging.logger {}

@Service
class RoomQueryService(private val helper: RoomQueryHelper) {

    suspend fun getRooms(page: Int, size: Int, status: String?, hasSlots: Boolean?): RoomListResponse {
        val statusEnum = helper.parseRoomStatus(status)
        val offset = page * size

        val data = helper.fetchRoomListData(statusEnum, hasSlots, size, offset)

        if (data.rooms.isEmpty()) {
            return RoomListResponse(emptyList(), page, size, 0, 0)
        }

        return RoomListResponse(
            rooms = helper.buildRoomDtoList(data),
            page = page,
            size = size,
            totalElements = data.total,
            totalPages = helper.calculateTotalPages(data.total, size)
        )
    }

    suspend fun getRoomById(roomId: UUID): RoomDto {
        val data = helper.fetchRoomById(roomId)
        return helper.buildRoomDto(data)
    }

    suspend fun getRoomByCode(code: String): RoomDto {
        val data = helper.fetchRoomByCode(code)
        return helper.buildRoomDto(data)
    }

    suspend fun getMyActiveRoom(userId: UUID): RoomDto? {
        val roomData = helper.fetchActiveRoomByUserId(userId)
        return roomData?.let { helper.buildRoomDto(it) }
    }
}