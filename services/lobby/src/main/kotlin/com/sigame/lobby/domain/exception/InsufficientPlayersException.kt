package com.sigame.lobby.domain.exception

import java.util.UUID

class InsufficientPlayersException(
    roomId: UUID,
    currentPlayers: Int,
    minRequired: Int = 2
) : RuntimeException("Room $roomId has only $currentPlayers players, but requires at least $minRequired to start")

