#!/bin/bash

# SSL Setup Script for Let's Encrypt
# This script automates SSL certificate generation using certbot

set -e

echo "=========================================="
echo "   SIGame SSL Setup (Let's Encrypt)"
echo "=========================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Error: Please run as root (use sudo)"
    exit 1
fi

# Ask for domain name or IP
read -p "Enter your domain name (or press Enter to skip SSL setup): " DOMAIN

if [ -z "$DOMAIN" ]; then
    echo "No domain provided. Skipping SSL setup."
    echo "You can run this script later when you have a domain."
    exit 0
fi

# Ask for email
read -p "Enter your email for Let's Encrypt notifications: " EMAIL

if [ -z "$EMAIL" ]; then
    echo "Error: Email is required for Let's Encrypt"
    exit 1
fi

# Install certbot if not installed
if ! command -v certbot &> /dev/null; then
    echo "Installing certbot..."
    apt-get update
    apt-get install -y certbot python3-certbot-nginx
fi

# Stop nginx if running
if systemctl is-active --quiet nginx; then
    echo "Stopping nginx..."
    systemctl stop nginx
fi

# Create directory for Let's Encrypt verification
mkdir -p /var/www/certbot

# Obtain certificate
echo "Obtaining SSL certificate for $DOMAIN..."
certbot certonly --standalone \
    --preferred-challenges http \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN"

if [ $? -eq 0 ]; then
    echo "✓ SSL certificate obtained successfully!"
    
    # Update nginx configuration with domain and SSL
    NGINX_CONF="/opt/sigame/deployment/nginx/nginx.conf"
    
    if [ -f "$NGINX_CONF" ]; then
        echo "Updating nginx configuration..."
        
        # Backup original config
        cp "$NGINX_CONF" "${NGINX_CONF}.bak"
        
        # Replace placeholders
        sed -i "s/YOUR_DOMAIN_OR_IP/$DOMAIN/g" "$NGINX_CONF"
        sed -i "s/YOUR_DOMAIN/$DOMAIN/g" "$NGINX_CONF"
        
        # Uncomment HTTPS server block
        sed -i 's/# server {/server {/g' "$NGINX_CONF"
        sed -i 's/#     listen/    listen/g' "$NGINX_CONF"
        sed -i 's/#     server_name/    server_name/g' "$NGINX_CONF"
        sed -i 's/#     ssl_/    ssl_/g' "$NGINX_CONF"
        sed -i 's/#     location/    location/g' "$NGINX_CONF"
        sed -i 's/#     }/    }/g' "$NGINX_CONF"
        sed -i 's/# }/}/g' "$NGINX_CONF"
        
        # Uncomment redirect to HTTPS
        sed -i 's/# location \/ {/location \/ {/g' "$NGINX_CONF"
        sed -i 's/#     return 301/    return 301/g' "$NGINX_CONF"
        
        # Comment out HTTP location block
        sed -i 's/location \/ {/# location \/ {/g' "$NGINX_CONF"
        sed -i 's/proxy_pass http:\/\/frontend/# proxy_pass http:\/\/frontend/g' "$NGINX_CONF"
        
        echo "✓ Nginx configuration updated"
    fi
    
    # Setup auto-renewal
    echo "Setting up automatic certificate renewal..."
    
    # Create renewal script
    cat > /etc/cron.daily/certbot-renew << 'EOF'
#!/bin/bash
certbot renew --quiet --post-hook "systemctl reload nginx"
EOF
    
    chmod +x /etc/cron.daily/certbot-renew
    
    echo "✓ Auto-renewal configured"
    
    # Copy nginx config to system
    cp "$NGINX_CONF" /etc/nginx/nginx.conf
    
    # Test nginx configuration
    echo "Testing nginx configuration..."
    nginx -t
    
    if [ $? -eq 0 ]; then
        # Start nginx
        echo "Starting nginx..."
        systemctl start nginx
        systemctl enable nginx
        
        echo ""
        echo "=========================================="
        echo "✓ SSL Setup Complete!"
        echo "=========================================="
        echo "Your site is now available at:"
        echo "  https://$DOMAIN"
        echo ""
        echo "Certificate will auto-renew every 90 days"
        echo "=========================================="
    else
        echo "Error: Nginx configuration test failed"
        exit 1
    fi
    
else
    echo "Error: Failed to obtain SSL certificate"
    exit 1
fi

