package com.sigame.lobby.grpc.auth

import java.util.UUID

data class UserInfo(
    val userId: UUID,
    val username: String,
    val avatarUrl: String?
)

