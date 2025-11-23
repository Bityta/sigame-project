package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается при попытке неавторизованного действия с комнатой
 */
class UnauthorizedRoomActionException(
    userId: UUID,
    action: String,
    reason: String = "User is not authorized to perform this action"
) : RuntimeException("User $userId is not authorized to $action: $reason")

