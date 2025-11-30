package com.sigame.lobby.domain.exception

import java.util.UUID

class PackNotOwnedException(
    val packId: UUID,
    val userId: UUID
) : RuntimeException("User $userId is not the owner of pack $packId")

