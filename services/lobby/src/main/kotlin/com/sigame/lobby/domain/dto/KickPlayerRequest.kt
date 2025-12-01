package com.sigame.lobby.domain.dto

import java.util.UUID

data class KickPlayerRequest(
    val targetUserId: UUID
)

