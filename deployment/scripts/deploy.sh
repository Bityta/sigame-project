#!/bin/bash

# Deployment Script - Updates and restarts services
# Use this script to deploy new versions

set -e

# Configuration
INSTALL_DIR="/opt/sigame"
SERVER_TYPE=${1:-"app"}  # app or infra

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

cd "$INSTALL_DIR"

echo "=========================================="
echo "  Deploying Updates ($SERVER_TYPE server)"
echo "=========================================="

# Pull latest code
echo -e "${BLUE}Step 1: Pulling latest code...${NC}"
git pull origin main
echo -e "${GREEN}✓ Code updated${NC}"

# Determine which compose file to use
if [ "$SERVER_TYPE" = "app" ]; then
    COMPOSE_FILE="docker-compose.app.yml"
    echo "Deploying Application Server..."
elif [ "$SERVER_TYPE" = "infra" ]; then
    COMPOSE_FILE="docker-compose.infra.yml"
    echo "Deploying Infrastructure Server..."
else
    echo -e "${YELLOW}Invalid server type. Usage: $0 [app|infra]${NC}"
    exit 1
fi

# Load environment
if [ -f .env.production ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

# Build new images
echo -e "${BLUE}Step 2: Building images...${NC}"
docker compose -f "$COMPOSE_FILE" build --no-cache
echo -e "${GREEN}✓ Images built${NC}"

# Pull latest base images
echo -e "${BLUE}Step 3: Pulling base images...${NC}"
docker compose -f "$COMPOSE_FILE" pull
echo -e "${GREEN}✓ Images pulled${NC}"

# Stop services
echo -e "${BLUE}Step 4: Stopping services...${NC}"
docker compose -f "$COMPOSE_FILE" down
echo -e "${GREEN}✓ Services stopped${NC}"

# Start services
echo -e "${BLUE}Step 5: Starting services...${NC}"
docker compose -f "$COMPOSE_FILE" up -d
echo -e "${GREEN}✓ Services started${NC}"

# Wait for services to be healthy
echo -e "${BLUE}Step 6: Waiting for services to be healthy...${NC}"
sleep 10

# Check health
echo -e "${BLUE}Step 7: Checking service health...${NC}"
docker compose -f "$COMPOSE_FILE" ps

# Cleanup old images
echo -e "${BLUE}Step 8: Cleaning up old images...${NC}"
docker image prune -f
echo -e "${GREEN}✓ Cleanup complete${NC}"

echo ""
echo "=========================================="
echo -e "${GREEN}✓ Deployment Complete!${NC}"
echo "=========================================="
echo ""
echo "Services Status:"
docker compose -f "$COMPOSE_FILE" ps
echo ""
echo "To view logs:"
echo "  docker compose -f $COMPOSE_FILE logs -f [service-name]"
echo "=========================================="

