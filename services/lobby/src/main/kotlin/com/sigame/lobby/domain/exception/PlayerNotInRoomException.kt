package com.sigame.lobby.domain.exception

import java.util.UUID

/**
 * Исключение выбрасывается когда игрок не находится в указанной комнате
 */
class PlayerNotInRoomException(
    userId: UUID,
    roomId: UUID
) : RuntimeException("User $userId is not in room $roomId")

