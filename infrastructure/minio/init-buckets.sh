#!/bin/sh

# ==============================================
# SIGame 2.0 - MinIO Buckets Initialization
# ==============================================

set -e

echo "Waiting for MinIO to be ready..."
sleep 10

MINIO_HOST="minio:9000"
MINIO_ALIAS="sigame"

# Configure MinIO client
echo "Configuring MinIO client..."
mc alias set $MINIO_ALIAS http://$MINIO_HOST ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin}

echo "✓ MinIO client configured"

# Create bucket for original .siq packs
echo "Creating bucket: sigame-packs..."
mc mb --ignore-existing ${MINIO_ALIAS}/sigame-packs
echo "✓ Created bucket: sigame-packs"

# Create bucket for extracted media files
echo "Creating bucket: sigame-media..."
mc mb --ignore-existing ${MINIO_ALIAS}/sigame-media
echo "✓ Created bucket: sigame-media"

# Set anonymous read policy for media bucket (so clients can access media)
echo "Setting public read policy for sigame-media..."
mc anonymous set download ${MINIO_ALIAS}/sigame-media
echo "✓ Public read access enabled for sigame-media"

# Keep sigame-packs private (only accessible through signed URLs)
echo "Setting private policy for sigame-packs..."
mc anonymous set none ${MINIO_ALIAS}/sigame-packs
echo "✓ Private access set for sigame-packs"

echo ""
echo "MinIO buckets initialized successfully!"
echo ""
echo "Listing all buckets:"
mc ls $MINIO_ALIAS

echo ""
echo "Bucket policies:"
echo "sigame-packs (private):"
mc anonymous get ${MINIO_ALIAS}/sigame-packs
echo ""
echo "sigame-media (public read):"
mc anonymous get ${MINIO_ALIAS}/sigame-media

