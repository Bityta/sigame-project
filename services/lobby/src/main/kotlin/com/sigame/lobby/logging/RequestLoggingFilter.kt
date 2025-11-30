package com.sigame.lobby.logging

import com.sigame.lobby.controller.ApiRoutes
import org.slf4j.MDC
import org.springframework.core.Ordered
import org.springframework.core.annotation.Order
import org.springframework.core.io.buffer.DataBuffer
import org.springframework.http.server.reactive.ServerHttpRequestDecorator
import org.springframework.http.server.reactive.ServerHttpResponseDecorator
import org.springframework.stereotype.Component
import org.springframework.web.server.ServerWebExchange
import org.springframework.web.server.WebFilter
import org.springframework.web.server.WebFilterChain
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import java.nio.charset.StandardCharsets
import java.util.UUID

@Component
@Order(Ordered.HIGHEST_PRECEDENCE + 1)
class RequestLoggingFilter(
    private val logWriter: AsyncLogWriter
) : WebFilter {

    companion object {
        private val SKIP_PATHS = ApiRoutes.SkipLogging.PATHS
        private const val MAX_BODY_SIZE = 10_000
        private val MDC_KEYS = listOf("request_id", "trace_id", "http_method", "http_path", "user_id")
    }

    override fun filter(exchange: ServerWebExchange, chain: WebFilterChain): Mono<Void> {
        val path = exchange.request.path.value()

        if (SKIP_PATHS.any { path.contains(it, ignoreCase = true) }) {
            return chain.filter(exchange)
        }

        val requestId = exchange.request.headers.getFirst("X-Request-ID") ?: UUID.randomUUID().toString()
        val startTime = System.currentTimeMillis()
        val requestBody = BodyCapture(MAX_BODY_SIZE)
        val responseBody = BodyCapture(MAX_BODY_SIZE)

        setupMdc(exchange, requestId)
        logWriter.logRequest(exchange.request, requestId)

        val decoratedExchange = decorateExchange(exchange, requestBody, responseBody)

        return chain.filter(decoratedExchange)
            .doOnSuccess {
                logWriter.logResponse(
                    request = decoratedExchange.request,
                    response = decoratedExchange.response,
                    requestId = requestId,
                    requestBody = requestBody.content,
                    responseBody = responseBody.content,
                    durationMs = System.currentTimeMillis() - startTime,
                    error = null
                )
            }
            .doOnError { error ->
                logWriter.logResponse(
                    request = decoratedExchange.request,
                    response = decoratedExchange.response,
                    requestId = requestId,
                    requestBody = requestBody.content,
                    responseBody = responseBody.content,
                    durationMs = System.currentTimeMillis() - startTime,
                    error = error
                )
            }
            .doFinally {
                MDC_KEYS.forEach(MDC::remove)
            }
    }

    private fun setupMdc(exchange: ServerWebExchange, requestId: String) {
        MDC.put("request_id", requestId)
        MDC.put("trace_id", requestId)
        MDC.put("http_method", exchange.request.method.name())
        MDC.put("http_path", exchange.request.path.value())
        exchange.request.headers.getFirst("X-User-ID")?.let { MDC.put("user_id", it) }
    }

    private fun decorateExchange(
        exchange: ServerWebExchange,
        requestBody: BodyCapture,
        responseBody: BodyCapture
    ): ServerWebExchange {
        val decoratedRequest = object : ServerHttpRequestDecorator(exchange.request) {
            override fun getBody(): Flux<DataBuffer> =
                super.getBody().doOnNext { requestBody.capture(it) }
        }

        val decoratedResponse = object : ServerHttpResponseDecorator(exchange.response) {
            override fun writeWith(body: org.reactivestreams.Publisher<out DataBuffer>): Mono<Void> =
                super.writeWith(Flux.from(body).doOnNext { responseBody.capture(it) })
        }

        return exchange.mutate()
            .request(decoratedRequest)
            .response(decoratedResponse)
            .build()
    }
}

private class BodyCapture(private val maxSize: Int) {
    private val buffer = StringBuilder()

    val content: String get() = buffer.toString()

    @Synchronized
    fun capture(dataBuffer: DataBuffer) {
        val readableBytes = dataBuffer.readableByteCount()
        if (readableBytes > 0 && buffer.length < maxSize) {
            val bytesToRead = minOf(readableBytes, maxSize - buffer.length)
            val bytes = ByteArray(bytesToRead)
            @Suppress("DEPRECATION")
            dataBuffer.toByteBuffer().get(bytes)
            buffer.append(String(bytes, StandardCharsets.UTF_8))
        }
    }
}
