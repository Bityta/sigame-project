package com.sigame.lobby.config

import com.sigame.lobby.domain.PlayerRole
import com.sigame.lobby.domain.RoomStatus
import io.r2dbc.postgresql.codec.EnumCodec
import io.r2dbc.spi.ConnectionFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.core.convert.converter.Converter
import org.springframework.data.convert.ReadingConverter
import org.springframework.data.convert.WritingConverter
import org.springframework.data.r2dbc.config.AbstractR2dbcConfiguration
import org.springframework.data.r2dbc.convert.R2dbcCustomConversions
import org.springframework.data.r2dbc.repository.config.EnableR2dbcRepositories
import org.springframework.r2dbc.core.DatabaseClient

@Configuration
@EnableR2dbcRepositories(basePackages = ["com.sigame.lobby.domain"])
class R2dbcConfig(private val connectionFactory: ConnectionFactory) : AbstractR2dbcConfiguration() {

    override fun connectionFactory(): ConnectionFactory = connectionFactory

    @Bean
    override fun r2dbcCustomConversions(): R2dbcCustomConversions {
        val converters = listOf(
            RoomStatusReadConverter(),
            RoomStatusWriteConverter(),
            PlayerRoleReadConverter(),
            PlayerRoleWriteConverter()
        )
        return R2dbcCustomConversions(storeConversions, converters)
    }

    @ReadingConverter
    class RoomStatusReadConverter : Converter<String, RoomStatus> {
        override fun convert(source: String): RoomStatus {
            return RoomStatus.valueOf(source.uppercase())
        }
    }

    @WritingConverter
    class RoomStatusWriteConverter : Converter<RoomStatus, String> {
        override fun convert(source: RoomStatus): String {
            return source.name.lowercase()
        }
    }

    @ReadingConverter
    class PlayerRoleReadConverter : Converter<String, PlayerRole> {
        override fun convert(source: String): PlayerRole {
            return PlayerRole.valueOf(source.uppercase())
        }
    }

    @WritingConverter
    class PlayerRoleWriteConverter : Converter<PlayerRole, String> {
        override fun convert(source: PlayerRole): String {
            return source.name.lowercase()
        }
    }
}

