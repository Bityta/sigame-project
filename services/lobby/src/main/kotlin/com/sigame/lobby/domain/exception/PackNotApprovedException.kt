package com.sigame.lobby.domain.exception

import java.util.UUID

class PackNotApprovedException(
    val packId: UUID,
    val status: String
) : RuntimeException("Pack $packId is not approved (status: $status)")

