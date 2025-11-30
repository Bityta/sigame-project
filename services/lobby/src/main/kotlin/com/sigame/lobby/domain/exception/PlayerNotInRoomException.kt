package com.sigame.lobby.domain.exception

import java.util.UUID

class PlayerNotInRoomException(
    userId: UUID,
    roomId: UUID
) : RuntimeException("User $userId is not in room $roomId")

