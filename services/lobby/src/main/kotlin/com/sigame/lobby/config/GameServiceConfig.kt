package com.sigame.lobby.config

import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationProperties(prefix = "game-service")
data class GameServiceConfig(
    var host: String = "game-service",
    var port: Int = 8003
) {
    val baseUrl: String
        get() = "http://$host:$port"
}

