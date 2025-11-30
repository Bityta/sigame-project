package com.sigame.lobby.domain.exception

import java.util.UUID

class PackInfoNotFoundException(packId: UUID) : RuntimeException("Pack info not found: $packId")

