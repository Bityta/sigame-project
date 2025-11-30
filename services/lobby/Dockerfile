FROM gradle:8.5-jdk21 AS builder
WORKDIR /app

# Copy proto files first
COPY proto/ /proto/

# Copy service source
COPY services/lobby/ .

# Update proto source path in build.gradle.kts for Docker context
RUN sed -i 's|srcDir("../../../proto")|srcDir("/proto")|g' build.gradle.kts

# Build (proto generation happens automatically via Gradle plugin)
RUN gradle build --no-daemon -x test

FROM amazoncorretto:21-alpine
WORKDIR /app
COPY --from=builder /app/build/libs/*.jar app.jar

EXPOSE 8082
ENV SPRING_PROFILES_ACTIVE=prod

HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8082/api/lobby/health || exit 1

ENTRYPOINT ["java", "-jar", "app.jar"]
