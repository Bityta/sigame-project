package com.sigame.lobby.domain.exception

import java.util.UUID

class RoomNotFoundException(roomId: UUID) : RuntimeException("Room with id $roomId not found")

class RoomNotFoundByCodeException(roomCode: String) : RuntimeException("Room with code $roomCode not found")

