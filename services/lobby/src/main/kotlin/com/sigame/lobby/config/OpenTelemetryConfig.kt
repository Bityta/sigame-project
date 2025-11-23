package com.sigame.lobby.config

import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.common.Attributes
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter
import io.opentelemetry.sdk.OpenTelemetrySdk
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.sdk.trace.SdkTracerProvider
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor
import io.opentelemetry.semconv.ResourceAttributes
import mu.KotlinLogging
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

private val logger = KotlinLogging.logger {}

@Configuration
class OpenTelemetryConfig {

    @Value("\${otel.exporter.otlp.endpoint:http://tempo:4317}")
    private lateinit var otlpEndpoint: String

    @Value("\${otel.service.name:lobby-service}")
    private lateinit var serviceName: String

    @Bean
    fun openTelemetry(): OpenTelemetry {
        logger.info { "Initializing OpenTelemetry..." }
        logger.info { "  Service: $serviceName" }
        logger.info { "  OTLP Endpoint: $otlpEndpoint" }

        // Create resource with service information
        val resource = Resource.getDefault().merge(
            Resource.create(
                Attributes.of(
                    ResourceAttributes.SERVICE_NAME, serviceName,
                    ResourceAttributes.SERVICE_VERSION, "1.0.0"
                )
            )
        )

        // Create OTLP exporter
        val spanExporter = OtlpGrpcSpanExporter.builder()
            .setEndpoint(otlpEndpoint)
            .build()

        // Create tracer provider
        val tracerProvider = SdkTracerProvider.builder()
            .addSpanProcessor(BatchSpanProcessor.builder(spanExporter).build())
            .setResource(resource)
            .build()

        // Build OpenTelemetry SDK
        val openTelemetry = OpenTelemetrySdk.builder()
            .setTracerProvider(tracerProvider)
            .buildAndRegisterGlobal()

        logger.info { "✓ OpenTelemetry initialized for $serviceName" }

        return openTelemetry
    }
}

