#!/bin/bash

# Setup Script for Application Server (Server 1)
# Installs Docker, clones repository, and starts application services

set -e

echo "=========================================="
echo "  Application Server Setup (Server 1)"
echo "=========================================="

# Configuration
REPO_URL=${REPO_URL:-"https://github.com/YOUR_USERNAME/sigame-project.git"}
INSTALL_DIR="/opt/sigame"
BRANCH=${BRANCH:-"main"}

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi

echo -e "${BLUE}Step 1: System Update${NC}"
apt-get update
apt-get upgrade -y

echo -e "${BLUE}Step 2: Installing Docker${NC}"
# Install Docker
if ! command -v docker &> /dev/null; then
    apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release

    mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Add ubuntu user to docker group
    usermod -aG docker ubuntu

    echo -e "${GREEN}✓ Docker installed${NC}"
else
    echo -e "${GREEN}✓ Docker already installed${NC}"
fi

echo -e "${BLUE}Step 3: Installing Additional Tools${NC}"
apt-get install -y \
    git \
    make \
    htop \
    curl \
    wget \
    vim \
    nginx \
    certbot \
    python3-certbot-nginx

echo -e "${GREEN}✓ Tools installed${NC}"

echo -e "${BLUE}Step 4: Configuring System${NC}"
# Increase file descriptors
cat >> /etc/security/limits.conf << EOF
* soft nofile 65535
* hard nofile 65535
EOF

# Optimize network settings
cat > /etc/sysctl.d/99-sigame.conf << EOF
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_fin_timeout = 30
EOF

sysctl -p /etc/sysctl.d/99-sigame.conf

echo -e "${GREEN}✓ System configured${NC}"

echo -e "${BLUE}Step 5: Setting up Project Directory${NC}"
# Create project directory
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

# Clone repository (if not exists)
if [ ! -d ".git" ]; then
    echo "Enter your Git repository URL (or press Enter to skip):"
    read USER_REPO_URL
    
    if [ -n "$USER_REPO_URL" ]; then
        REPO_URL="$USER_REPO_URL"
        git clone -b "$BRANCH" "$REPO_URL" .
        echo -e "${GREEN}✓ Repository cloned${NC}"
    else
        echo "Skipping repository clone. Please copy your project files manually to $INSTALL_DIR"
    fi
else
    echo -e "${GREEN}✓ Repository already exists${NC}"
    git pull origin "$BRANCH"
fi

# Set ownership
chown -R ubuntu:ubuntu "$INSTALL_DIR"

echo -e "${BLUE}Step 6: Setting up Environment${NC}"
if [ ! -f ".env.production" ]; then
    echo "Please create .env.production file with your configuration"
    echo "Example provided in .env.production.example"
    
    # Create from example if exists
    if [ -f ".env.production.example" ]; then
        cp .env.production.example .env.production
        echo -e "${GREEN}✓ Created .env.production from example${NC}"
        echo "IMPORTANT: Edit .env.production with correct values!"
    fi
else
    echo -e "${GREEN}✓ .env.production exists${NC}"
fi

echo -e "${BLUE}Step 7: Setting up Nginx${NC}"
# Copy nginx configuration
if [ -f "deployment/nginx/nginx.conf" ]; then
    cp deployment/nginx/nginx.conf /etc/nginx/nginx.conf
    nginx -t && systemctl restart nginx
    systemctl enable nginx
    echo -e "${GREEN}✓ Nginx configured${NC}"
else
    echo "Warning: nginx.conf not found"
fi

echo -e "${BLUE}Step 8: Building and Starting Services${NC}"
# Build and start services as ubuntu user
su - ubuntu << 'EOSU'
cd /opt/sigame

# Load environment
if [ -f .env.production ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

# Build images
docker compose -f docker-compose.app.yml build

# Start services
docker compose -f docker-compose.app.yml up -d

# Show status
docker compose -f docker-compose.app.yml ps
EOSU

echo ""
echo "=========================================="
echo -e "${GREEN}✓ Application Server Setup Complete!${NC}"
echo "=========================================="
echo ""
echo "Services Status:"
docker compose -f docker-compose.app.yml ps
echo ""
echo "Next Steps:"
echo "1. Edit .env.production with correct INFRA_SERVER_IP"
echo "2. Setup SSL: sudo ./deployment/nginx/ssl-setup.sh"
echo "3. Check logs: docker compose -f docker-compose.app.yml logs -f"
echo "=========================================="

