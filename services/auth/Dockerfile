# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies including protoc
RUN apk add --no-cache git make protoc protobuf-dev

# Install protoc plugins (versions compatible with grpc v1.60.1 in go.mod)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Set working directory
WORKDIR /app

# Copy shared proto files
COPY proto/ /proto/

# Copy service source
COPY services/auth/ .

# Generate proto from shared proto files
RUN mkdir -p proto && \
    protoc --proto_path=/proto \
    --go_out=proto --go_opt=paths=source_relative \
    --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
    /proto/auth/auth.proto && \
    mv proto/auth/*.go proto/ && \
    rmdir proto/auth

# Download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose ports
EXPOSE 8001 50051

# Healthcheck
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8001/health || exit 1

# Run the binary
CMD ["./server"]
