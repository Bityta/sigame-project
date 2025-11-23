package com.sigame.lobby.config

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.cors.CorsConfiguration
import org.springframework.web.cors.reactive.CorsWebFilter
import org.springframework.web.cors.reactive.UrlBasedCorsConfigurationSource

@Configuration
class CorsConfig {

    @Bean
    fun corsWebFilter(): CorsWebFilter {
        val corsConfig = CorsConfiguration()
        
        // Allow frontend origins (development + production)
        corsConfig.allowedOrigins = listOf(
            "http://localhost:3000",
            "http://localhost:3001",
            "http://localhost:5173", // Vite dev server default port
            "http://127.0.0.1:3000",
            "http://127.0.0.1:3001",
            "http://127.0.0.1:5173",
            // Docker/Production origins
            "http://frontend:80",
            "http://frontend",
            // Production server IP
            "http://89.169.139.21",
            "http://89.169.139.21:80"
        )
        
        // Allow all HTTP methods
        corsConfig.allowedMethods = listOf("GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS")
        
        // Allow all headers
        corsConfig.allowedHeaders = listOf("*")
        
        // Allow credentials (cookies, authorization headers)
        corsConfig.allowCredentials = true
        
        // Max age for preflight requests
        corsConfig.maxAge = 3600L

        val source = UrlBasedCorsConfigurationSource()
        source.registerCorsConfiguration("/**", corsConfig)

        return CorsWebFilter(source)
    }
}

