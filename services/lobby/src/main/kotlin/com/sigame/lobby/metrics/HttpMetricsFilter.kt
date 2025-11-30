package com.sigame.lobby.metrics

import org.springframework.core.Ordered
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Mono

@Component
class HttpMetricsFilter(
    private val httpMetrics: HttpMetrics
) : WebFilter, Ordered {

    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val startTime = System.currentTimeMillis()
        val method = exchange.request.method.name()
        val path = exchange.request.path.value()

        return chain.filter(exchange)
            .doFinally {
                val duration = System.currentTimeMillis() - startTime
                val status = exchange.response.statusCode?.value() ?: 500
                httpMetrics.recordHttpRequest(method, path, status, duration)
            }
    }

    override fun getOrder(): Int = Ordered.HIGHEST_PRECEDENCE + 10
}
