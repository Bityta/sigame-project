#!/bin/bash

# Script to generate Python code from shared proto files

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTO_DIR="$SCRIPT_DIR/../../proto"
OUTPUT_DIR="$SCRIPT_DIR/app/grpc"

echo "Generating gRPC code from shared proto files..."

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate proto code for pack
python -m grpc_tools.protoc \
    --proto_path="$PROTO_DIR" \
    --python_out="$OUTPUT_DIR" \
    --grpc_python_out="$OUTPUT_DIR" \
    "$PROTO_DIR/pack/pack.proto"

# Fix imports in generated files (Python grpc_tools issue)
sed -i.bak 's/import pack_pb2/from app.grpc import pack_pb2/g' "$OUTPUT_DIR/pack_pb2_grpc.py" 2>/dev/null || \
sed -i '' 's/import pack_pb2/from app.grpc import pack_pb2/g' "$OUTPUT_DIR/pack_pb2_grpc.py"

rm -f "$OUTPUT_DIR"/*.bak

echo "✓ Proto files generated successfully in $OUTPUT_DIR"

