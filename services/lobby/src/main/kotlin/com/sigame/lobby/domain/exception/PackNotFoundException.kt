package com.sigame.lobby.domain.exception

import java.util.UUID

class PackNotFoundException(packId: UUID) : RuntimeException("Pack with id $packId not found")

