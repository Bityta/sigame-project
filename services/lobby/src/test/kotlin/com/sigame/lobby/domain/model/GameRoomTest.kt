package com.sigame.lobby.domain.model

import com.sigame.lobby.domain.enums.PlayerRole
import com.sigame.lobby.domain.enums.RoomStatus
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import java.time.LocalDateTime
import java.util.*

class GameRoomTest {

    @Test
    fun `should create game room with default values`() {
        val room = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID()
        )

        assertEquals("ABC123", room.code)
        assertEquals(RoomStatus.WAITING, room.status)
        assertNotNull(room.createdAt)
    }

    @Test
    fun `should validate room code format`() {
        val validRoom = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID()
        )
        assertTrue(validRoom.code.matches(Regex("[A-Z0-9]{6}")))

        val invalidRoom = GameRoom(
            code = "abc",
            hostId = UUID.randomUUID()
        )
        assertFalse(invalidRoom.code.matches(Regex("[A-Z0-9]{6}")))
    }

    @Test
    fun `should check if room is full`() {
        val room = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID(),
            maxPlayers = 4,
            currentPlayers = 4
        )
        assertTrue(room.isFull())

        val notFullRoom = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID(),
            maxPlayers = 4,
            currentPlayers = 2
        )
        assertFalse(notFullRoom.isFull())
    }

    @Test
    fun `should check if room can be started`() {
        val readyRoom = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID(),
            status = RoomStatus.WAITING,
            currentPlayers = 2,
            minPlayers = 2
        )
        assertTrue(readyRoom.canStart())

        val notReadyRoom = GameRoom(
            code = "ABC123",
            hostId = UUID.randomUUID(),
            status = RoomStatus.WAITING,
            currentPlayers = 1,
            minPlayers = 2
        )
        assertFalse(notReadyRoom.canStart())
    }
}

