#!/bin/bash

# Setup Script for Infrastructure Server (Server 2)
# Installs Docker, clones repository, and starts infrastructure services

set -e

echo "=========================================="
echo "  Infrastructure Server Setup (Server 2)"
echo "=========================================="

# Configuration
REPO_URL=${REPO_URL:-"https://github.com/YOUR_USERNAME/sigame-project.git"}
INSTALL_DIR="/opt/sigame"
BRANCH=${BRANCH:-"main"}
SWAP_SIZE="8G"

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
    postgresql-client

echo -e "${GREEN}✓ Tools installed${NC}"

echo -e "${BLUE}Step 4: Setting up Swap${NC}"
# Create swap file if not exists
if [ ! -f /swapfile ]; then
    fallocate -l "$SWAP_SIZE" /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
    
    # Optimize swap usage
    sysctl vm.swappiness=10
    sysctl vm.vfs_cache_pressure=50
    echo 'vm.swappiness=10' >> /etc/sysctl.conf
    echo 'vm.vfs_cache_pressure=50' >> /etc/sysctl.conf
    
    echo -e "${GREEN}✓ Swap created ($SWAP_SIZE)${NC}"
else
    echo -e "${GREEN}✓ Swap already exists${NC}"
fi

echo -e "${BLUE}Step 5: Optimizing System for Databases${NC}"
# PostgreSQL optimizations
cat > /etc/sysctl.d/99-postgresql.conf << EOF
# PostgreSQL Optimizations
vm.overcommit_memory = 2
vm.overcommit_ratio = 80

# Shared Memory (adjust based on RAM)
kernel.shmmax = 8589934592
kernel.shmall = 2097152

# Network optimizations
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_keepalive_time = 300
net.ipv4.tcp_keepalive_probes = 5
net.ipv4.tcp_keepalive_intvl = 15
EOF

sysctl -p /etc/sysctl.d/99-postgresql.conf

# Increase file descriptors
cat >> /etc/security/limits.conf << EOF
* soft nofile 65535
* hard nofile 65535
* soft nproc 65535
* hard nproc 65535
EOF

echo -e "${GREEN}✓ System optimized${NC}"

echo -e "${BLUE}Step 6: Setting up Project Directory${NC}"
# Create project directory
mkdir -p "$INSTALL_DIR"
mkdir -p /backups/postgresql
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
chown -R ubuntu:ubuntu /backups

echo -e "${BLUE}Step 7: Setting up Environment${NC}"
if [ ! -f ".env.production" ]; then
    echo "Please create .env.production file with your configuration"
    
    # Create from example if exists
    if [ -f ".env.production.example" ]; then
        cp .env.production.example .env.production
        echo -e "${GREEN}✓ Created .env.production from example${NC}"
        echo "IMPORTANT: Edit .env.production with correct values!"
    fi
else
    echo -e "${GREEN}✓ .env.production exists${NC}"
fi

echo -e "${BLUE}Step 8: Setting up Backup Cron Job${NC}"
# Setup daily backup cron job
cat > /etc/cron.daily/sigame-backup << 'EOFCRON'
#!/bin/bash
/opt/sigame/deployment/scripts/backup.sh >> /var/log/sigame-backup.log 2>&1
EOFCRON

chmod +x /etc/cron.daily/sigame-backup
echo -e "${GREEN}✓ Backup cron job created${NC}"

echo -e "${BLUE}Step 9: Starting Infrastructure Services${NC}"
# Start services as ubuntu user
su - ubuntu << 'EOSU'
cd /opt/sigame

# Load environment
if [ -f .env.production ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

# Start services
docker compose -f docker-compose.infra.yml up -d

# Wait for services to be healthy
echo "Waiting for services to start..."
sleep 30

# Show status
docker compose -f docker-compose.infra.yml ps
EOSU

echo ""
echo "=========================================="
echo -e "${GREEN}✓ Infrastructure Server Setup Complete!${NC}"
echo "=========================================="
echo ""
echo "Services Status:"
docker compose -f docker-compose.infra.yml ps
echo ""
echo "Important Information:"
echo "- PostgreSQL ports: 5432-5435"
echo "- Redis port: 6379"
echo "- Kafka port: 9092"
echo "- Grafana: http://localhost:3000 (admin/admin)"
echo "- Prometheus: http://localhost:9090"
echo "- MinIO Console: http://localhost:9001"
echo ""
echo "Next Steps:"
echo "1. Configure databases (they should auto-initialize)"
echo "2. Setup Grafana dashboards"
echo "3. Configure Application Server with this server's IP"
echo "4. Check logs: docker compose -f docker-compose.infra.yml logs -f"
echo "=========================================="

