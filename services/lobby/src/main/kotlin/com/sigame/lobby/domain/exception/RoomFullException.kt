package com.sigame.lobby.domain.exception

import java.util.UUID

class RoomFullException(roomId: UUID) : RuntimeException("Room $roomId is full")

