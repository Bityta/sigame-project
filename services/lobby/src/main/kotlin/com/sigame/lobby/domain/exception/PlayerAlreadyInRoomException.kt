package com.sigame.lobby.domain.exception

import java.util.UUID

class PlayerAlreadyInRoomException(
    userId: UUID,
    roomId: UUID
) : RuntimeException("User $userId is already in room $roomId")

