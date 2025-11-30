package com.sigame.lobby.domain.exception

import java.util.UUID

class CannotKickHostException(roomId: UUID) : 
    RuntimeException("Cannot kick the host from room $roomId")

