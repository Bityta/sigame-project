package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается когда игрок уже находится в комнате
 */
class PlayerAlreadyInRoomException(
    userId: UUID,
    roomId: UUID
) : RuntimeException("User $userId is already in room $roomId")

