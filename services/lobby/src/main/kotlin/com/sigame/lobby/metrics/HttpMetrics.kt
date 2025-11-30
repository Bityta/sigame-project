package com.sigame.lobby.metrics

import io.micrometer.core.instrument.Counter
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import org.springframework.stereotype.Component
import java.time.Duration

@Component
class HttpMetrics(
    private val meterRegistry: MeterRegistry
) {

    fun recordHttpRequest(method: String, endpoint: String, status: Int, durationMs: Long) {
        Counter.builder("http_requests_total")
            .tag("method", method)
            .tag("endpoint", endpoint)
            .tag("status", statusCodeGroup(status))
            .register(meterRegistry)
            .increment()

        Timer.builder("http_request_duration_seconds")
            .tag("method", method)
            .tag("endpoint", endpoint)
            .register(meterRegistry)
            .record(Duration.ofMillis(durationMs))
    }

    private fun statusCodeGroup(status: Int): String = when (status) {
        in 200..299 -> "2xx"
        in 300..399 -> "3xx"
        in 400..499 -> "4xx"
        in 500..599 -> "5xx"
        else -> "unknown"
    }
}
