package com.sigame.lobby.metrics

import io.micrometer.core.instrument.Counter
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import org.springframework.stereotype.Component

/**
 * Standard HTTP and gRPC metrics for dashboards
 */
@Component
class HttpMetrics(
    private val meterRegistry: MeterRegistry
) {
    
    /**
     * Record HTTP request
     */
    fun recordHttpRequest(method: String, endpoint: String, status: Int, durationMs: Long) {
        val statusCode = statusCodeString(status)
        
        // Counter for total requests
        Counter.builder("http_requests_total")
            .tag("method", method)
            .tag("endpoint", endpoint)
            .tag("status", statusCode)
            .register(meterRegistry)
            .increment()
        
        // Timer for request duration
        Timer.builder("http_request_duration_seconds")
            .tag("method", method)
            .tag("endpoint", endpoint)
            .register(meterRegistry)
            .record(java.time.Duration.ofMillis(durationMs))
    }
    
    /**
     * Record gRPC request
     */
    fun recordGrpcRequest(method: String, status: String, durationMs: Long) {
        // Counter for total requests
        Counter.builder("grpc_requests_total")
            .tag("method", method)
            .tag("status", status)
            .register(meterRegistry)
            .increment()
        
        // Timer for request duration
        Timer.builder("grpc_request_duration_seconds")
            .tag("method", method)
            .register(meterRegistry)
            .record(java.time.Duration.ofMillis(durationMs))
    }
    
    private fun statusCodeString(status: Int): String {
        return when {
            status in 200..299 -> "2xx"
            status in 300..399 -> "3xx"
            status in 400..499 -> "4xx"
            status >= 500 -> "5xx"
            else -> "unknown"
        }
    }
}

