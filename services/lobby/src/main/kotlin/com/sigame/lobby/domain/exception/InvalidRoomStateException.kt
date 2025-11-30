package com.sigame.lobby.domain.exception

import com.sigame.lobby.domain.enums.RoomStatus
import java.util.UUID

class InvalidRoomStateException(
    roomId: UUID,
    currentState: RoomStatus,
    action: String
) : RuntimeException("Cannot $action room $roomId in state $currentState")

