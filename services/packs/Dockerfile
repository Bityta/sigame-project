# Pack Service

FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first for better caching
COPY services/packs/requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy proto files
COPY proto/ /proto/

# Copy application code
COPY services/packs/ .

# Generate gRPC code from shared proto files
RUN mkdir -p app/grpc && \
    python -m grpc_tools.protoc \
    -I/proto \
    --python_out=./app/grpc \
    --grpc_python_out=./app/grpc \
    /proto/pack/pack.proto

# Fix imports in generated files
RUN sed -i 's/import pack_pb2/from app.grpc import pack_pb2/g' app/grpc/pack_pb2_grpc.py

# Expose ports
EXPOSE 8005 50055

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8005/health')" || exit 1

# Run the application
CMD ["python", "main.py"]
