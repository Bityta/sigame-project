package com.sigame.lobby.domain.exception

import java.util.UUID

class CannotKickSelfException(userId: UUID) : 
    RuntimeException("User $userId cannot kick themselves, use leave instead")

