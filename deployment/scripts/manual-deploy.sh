#!/bin/bash

# Manual Deploy Script
# Use this to deploy changes manually without git push

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Manual Deploy to Yandex Cloud${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
APP_SERVER_IP="89.169.139.21"
INFRA_SERVER_IP="10.129.0.26"

# Step 1: Deploy to Application Server
echo -e "${BLUE}Step 1: Deploying to Application Server...${NC}"

# Copy files
echo "Copying files..."
rsync -avz --exclude '.git' --exclude 'node_modules' --exclude '.terraform' --exclude 'frontend/dist' \
  -e ssh \
  ./ ubuntu@$APP_SERVER_IP:/opt/sigame/

# Restart services
echo "Restarting services..."
ssh ubuntu@$APP_SERVER_IP << 'EOF'
  cd /opt/sigame
  
  # Load environment
  export $(grep -v '^#' .env.production | xargs)
  
  # Build and restart
  docker compose -f docker-compose.app.yml build
  docker compose -f docker-compose.app.yml up -d --no-deps
  
  # Clean old images
  docker image prune -f
  
  echo "✅ Application Server updated"
EOF

echo -e "${GREEN}✓ Application Server deployed${NC}"
echo ""

# Step 2: Health checks
echo -e "${BLUE}Step 2: Running health checks...${NC}"
sleep 10

# Check services
echo "Checking Frontend..."
curl -f http://$APP_SERVER_IP:3001 > /dev/null 2>&1 && echo -e "${GREEN}✓ Frontend OK${NC}" || echo -e "${YELLOW}⚠ Frontend not responding${NC}"

echo "Checking Auth Service..."
curl -f http://$APP_SERVER_IP:8001/health > /dev/null 2>&1 && echo -e "${GREEN}✓ Auth Service OK${NC}" || echo -e "${YELLOW}⚠ Auth Service not responding${NC}"

echo "Checking Lobby Service..."
curl -f http://$APP_SERVER_IP:8002/api/lobby/health > /dev/null 2>&1 && echo -e "${GREEN}✓ Lobby Service OK${NC}" || echo -e "${YELLOW}⚠ Lobby Service not responding${NC}"

echo "Checking Game Service..."
curl -f http://$APP_SERVER_IP:8003/health > /dev/null 2>&1 && echo -e "${GREEN}✓ Game Service OK${NC}" || echo -e "${YELLOW}⚠ Game Service not responding${NC}"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Deployment Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Your application is available at:"
echo "  Frontend:  http://$APP_SERVER_IP:3001"
echo "  Grafana:   http://$APP_SERVER_IP:3000"
echo ""
echo "To check logs:"
echo "  ssh ubuntu@$APP_SERVER_IP 'cd /opt/sigame && docker compose -f docker-compose.app.yml logs -f'"
echo ""

