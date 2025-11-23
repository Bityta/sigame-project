package com.sigame.lobby.config

import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationProperties(prefix = "grpc.auth-service")
data class AuthServiceConfig(
    var host: String = "auth-service",
    var port: Int = 50051
)

@Configuration
@ConfigurationProperties(prefix = "grpc.pack-service")
data class PackServiceConfig(
    var host: String = "pack-service",
    var port: Int = 50055
)

@Configuration
@ConfigurationProperties(prefix = "room")
data class RoomConfig(
    var codeLength: Int = 6,
    var codeCharset: String = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
    var maxPlayersLimit: Int = 12,
    var cacheTtl: Long = 7200
)

