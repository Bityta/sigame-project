# ===================================================================
# Multi-stage Dockerfile for Lobby Service
# Stage 1: Build the application
# Stage 2: Runtime image
# ===================================================================

# ===================================================================
# Build Stage
# ===================================================================
FROM gradle:8.5-jdk21-alpine AS build

WORKDIR /app

# Increase memory for Gradle daemon
ENV GRADLE_OPTS="-Xmx1g -XX:MaxMetaspaceSize=256m -Dorg.gradle.daemon=false"

# Copy gradle configuration files first for layer caching
COPY build.gradle.kts settings.gradle.kts gradle.properties ./

# Download dependencies (this layer will be cached if build files haven't changed)
RUN gradle dependencies --no-daemon || true

# Copy source code
COPY src ./src

# Build the application (skip tests for faster builds, tests should run in CI)
RUN gradle build --no-daemon -x test && \
    # Verify JAR was built
    ls -lh build/libs/ && \
    # Clean up build cache to reduce image size
    rm -rf ~/.gradle/caches/

# ===================================================================
# Runtime Stage
# ===================================================================
FROM eclipse-temurin:21-jre-alpine

LABEL maintainer="SIGame Team" \
      description="Lobby Service for SIGame" \
      version="1.0.0"

WORKDIR /app

# Install wget for healthcheck and create non-root user
RUN apk add --no-cache wget && \
    addgroup -S sigame && \
    adduser -S sigame -G sigame

# Copy the built JAR from build stage
COPY --from=build --chown=sigame:sigame /app/build/libs/*.jar app.jar

# Switch to non-root user
USER sigame

# Expose HTTP port
EXPOSE 8002

# Configure JVM options for containerized environment
ENV JAVA_OPTS="-XX:+UseContainerSupport \
               -XX:MaxRAMPercentage=75.0 \
               -XX:InitialRAMPercentage=50.0 \
               -XX:+UseG1GC \
               -XX:+UseStringDeduplication \
               -Djava.security.egd=file:/dev/./urandom"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=60s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8002/api/lobby/health || exit 1

# Run the application
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]

