#!/bin/bash

# Script to fix network issues on Yandex Cloud Ubuntu servers
# Solves IPv6 issues and configures alternative mirrors

set -e

echo "=========================================="
echo "  Network Configuration Fix"
echo "=========================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi

echo "Step 1: Disabling IPv6 (common issue on Yandex Cloud)"
# Disable IPv6
cat >> /etc/sysctl.conf << EOF
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
EOF

sysctl -p

echo "Step 2: Configuring APT to use IPv4"
echo 'Acquire::ForceIPv4 "true";' > /etc/apt/apt.conf.d/99force-ipv4

echo "Step 3: Updating DNS settings"
# Backup original resolv.conf
cp /etc/resolv.conf /etc/resolv.conf.backup 2>/dev/null || true

# Use Google DNS and Cloudflare DNS
cat > /etc/resolv.conf << EOF
nameserver 8.8.8.8
nameserver 8.8.4.4
nameserver 1.1.1.1
nameserver 1.0.0.1
EOF

# Make resolv.conf immutable
chattr +i /etc/resolv.conf 2>/dev/null || true

echo "Step 4: Configuring alternative Ubuntu mirrors"
# Backup original sources.list
cp /etc/apt/sources.list /etc/apt/sources.list.backup 2>/dev/null || true

# Use main Ubuntu mirrors
cat > /etc/apt/sources.list << EOF
deb http://archive.ubuntu.com/ubuntu/ focal main restricted universe multiverse
deb http://archive.ubuntu.com/ubuntu/ focal-updates main restricted universe multiverse
deb http://archive.ubuntu.com/ubuntu/ focal-backports main restricted universe multiverse
deb http://security.ubuntu.com/ubuntu focal-security main restricted universe multiverse
EOF

echo "Step 5: Configuring APT for better timeout handling"
cat > /etc/apt/apt.conf.d/99timeout << EOF
Acquire::http::Timeout "30";
Acquire::https::Timeout "30";
Acquire::ftp::Timeout "30";
Acquire::Retries "5";
EOF

echo "Step 6: Testing network connectivity"
echo "Testing Google DNS..."
ping -c 3 8.8.8.8 || echo "Warning: Cannot ping 8.8.8.8"

echo "Testing DNS resolution..."
nslookup google.com || echo "Warning: DNS resolution issues"

echo "Step 7: Updating package list"
apt-get clean
apt-get update

echo ""
echo "=========================================="
echo "✓ Network Configuration Complete!"
echo "=========================================="
echo ""
echo "You can now run the setup scripts:"
echo "  - sudo bash deployment/scripts/setup-server1.sh  # For App Server"
echo "  - sudo bash deployment/scripts/setup-server2.sh  # For Infra Server"
echo ""

