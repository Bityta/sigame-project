package com.sigame.lobby.domain.exception

import java.util.UUID

class UserInfoNotFoundException(userId: UUID) : RuntimeException("User info not found: $userId")

