FROM gradle:8.5-jdk21 AS builder
WORKDIR /app

# Install protoc and grpc plugin
RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

COPY . .
RUN gradle build --no-daemon -x test

FROM amazoncorretto:21-alpine
WORKDIR /app
COPY --from=builder /app/build/libs/*.jar app.jar

EXPOSE 8082
ENV SPRING_PROFILES_ACTIVE=prod
ENTRYPOINT ["java", "-jar", "app.jar"]
