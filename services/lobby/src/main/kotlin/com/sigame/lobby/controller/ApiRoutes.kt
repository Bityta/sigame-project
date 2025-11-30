package com.sigame.lobby.controller

object ApiRoutes {

    const val BASE = "/api/lobby"

    const val HEALTH = "/health"

    object Rooms {
        const val BASE = "/rooms"

        const val LIST = BASE

        const val BY_ID = "$BASE/{id}"

        const val BY_CODE = "$BASE/code/{code}"

        const val JOIN = "$BASE/{id}/join"

        const val LEAVE = "$BASE/{id}/leave"

        const val START = "$BASE/{id}/start"

        const val SETTINGS = "$BASE/{id}/settings"
    }

    object Public {
        val PATHS = setOf(
            "$BASE$HEALTH",
            "/actuator/health",
            "/actuator/prometheus",
            "/actuator/metrics",
            "/actuator/info",
            "/metrics",
            "/api-docs",
            "/swagger-ui",
            "/swagger-ui.html",
            "/v3/api-docs"
        )

        val GET_PATHS = setOf(
            "$BASE${Rooms.BASE}"
        )
    }

    object SkipLogging {
        val PATHS = setOf("/health", "/metrics", "/actuator", "/swagger", "/api-docs")
    }
}

