package com.sigame.lobby.metrics

import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import org.springframework.stereotype.Component
import java.time.Duration
import java.util.concurrent.ConcurrentHashMap

@Component
class HttpMetrics(
    private val meterRegistry: MeterRegistry
) {
    private val timers = ConcurrentHashMap<String, Timer>()
    
    fun recordHttpRequest(method: String, endpoint: String, status: Int, durationMs: Long) {
        // Counter - simple increment
        meterRegistry.counter(
            "http_requests_total",
            "method", method,
            "endpoint", endpoint,
            "status", statusCodeGroup(status)
        ).increment()

        // Timer with histogram - cached per method+endpoint
        val timerKey = "$method:$endpoint"
        val timer = timers.computeIfAbsent(timerKey) {
            Timer.builder("http_request_duration_seconds")
                .tag("method", method)
                .tag("endpoint", endpoint)
                .publishPercentileHistogram()
                .serviceLevelObjectives(
                    Duration.ofMillis(10),
                    Duration.ofMillis(50),
                    Duration.ofMillis(100),
                    Duration.ofMillis(200),
                    Duration.ofMillis(500),
                    Duration.ofSeconds(1),
                    Duration.ofSeconds(5)
                )
                .register(meterRegistry)
        }
        timer.record(Duration.ofMillis(durationMs))
    }

    private fun statusCodeGroup(status: Int): String = when (status) {
        in 200..299 -> "2xx"
        in 300..399 -> "3xx"
        in 400..499 -> "4xx"
        in 500..599 -> "5xx"
        else -> "unknown"
    }
}
