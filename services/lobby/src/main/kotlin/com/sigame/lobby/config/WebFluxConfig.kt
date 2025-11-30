package com.sigame.lobby.config

import com.sigame.lobby.security.CurrentUserArgumentResolver
import org.springframework.context.annotation.Configuration
import org.springframework.web.reactive.config.WebFluxConfigurer
import org.springframework.web.reactive.result.method.annotation.ArgumentResolverConfigurer

@Configuration
class WebFluxConfig(
    private val currentUserArgumentResolver: CurrentUserArgumentResolver
) : WebFluxConfigurer {
    
    override fun configureArgumentResolvers(configurer: ArgumentResolverConfigurer) {
        configurer.addCustomResolver(currentUserArgumentResolver)
    }
}

