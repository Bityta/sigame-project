package com.sigame.lobby.grpc.pack

import java.util.UUID

data class PackInfo(
    val id: UUID,
    val name: String,
    val author: String,
    val roundsCount: Int,
    val questionsCount: Int
)

