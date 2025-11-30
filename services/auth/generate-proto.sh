#!/bin/bash

# Script to generate Go code from shared proto files

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="$SCRIPT_DIR/../../proto"
OUTPUT_DIR="$SCRIPT_DIR/proto"

echo "Generating gRPC code from shared proto files..."

# Install protoc-gen-go and protoc-gen-go-grpc if not installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate proto code for auth
protoc --proto_path="$PROTO_DIR" \
    --go_out="$OUTPUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$OUTPUT_DIR" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/auth/auth.proto"

# Move generated files to correct location
mv "$OUTPUT_DIR/auth/"*.go "$OUTPUT_DIR/" 2>/dev/null || true
rmdir "$OUTPUT_DIR/auth" 2>/dev/null || true

echo "✓ Proto files generated successfully in $OUTPUT_DIR"

