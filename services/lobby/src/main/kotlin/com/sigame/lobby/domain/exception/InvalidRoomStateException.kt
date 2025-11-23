package com.sigame.lobby.domain.exception

import com.sigame.lobby.domain.enums.RoomStatus
import java.util.UUID

/**
 * Исключение выбрасывается при попытке недопустимого действия в текущем состоянии комнаты
 */
class InvalidRoomStateException(
    roomId: UUID,
    currentState: RoomStatus,
    action: String
) : RuntimeException("Cannot $action room $roomId in state $currentState")

