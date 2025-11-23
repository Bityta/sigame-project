#!/bin/bash

# Install Tools Script - Terraform and Yandex Cloud CLI for macOS
# This script installs all necessary tools for deployment

set -e

echo "=========================================="
echo "   Installing Deployment Tools"
echo "=========================================="

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Error: Homebrew is not installed"
    echo "Please install Homebrew first:"
    echo "/bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    exit 1
fi

echo "✓ Homebrew found"

# Install Terraform
if ! command -v terraform &> /dev/null; then
    echo "Installing Terraform..."
    brew tap hashicorp/tap
    brew install hashicorp/tap/terraform
    echo "✓ Terraform installed"
else
    echo "✓ Terraform already installed ($(terraform version | head -n1))"
fi

# Install Yandex Cloud CLI
if ! command -v yc &> /dev/null; then
    echo "Installing Yandex Cloud CLI..."
    curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash
    
    # Add to PATH
    export PATH="$HOME/yandex-cloud/bin:$PATH"
    
    echo "✓ Yandex Cloud CLI installed"
    echo ""
    echo "IMPORTANT: Add this to your ~/.zshrc or ~/.bash_profile:"
    echo "export PATH=\"\$HOME/yandex-cloud/bin:\$PATH\""
    echo ""
    echo "Then run: source ~/.zshrc"
else
    echo "✓ Yandex Cloud CLI already installed ($(yc version | head -n1))"
fi

# Install other useful tools
echo "Installing additional tools..."

if ! command -v jq &> /dev/null; then
    brew install jq
    echo "✓ jq installed"
else
    echo "✓ jq already installed"
fi

if ! command -v git &> /dev/null; then
    brew install git
    echo "✓ git installed"
else
    echo "✓ git already installed"
fi

echo ""
echo "=========================================="
echo "✓ All tools installed successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Initialize Yandex Cloud CLI:"
echo "   yc init"
echo ""
echo "2. Configure Terraform:"
echo "   cd deployment/terraform"
echo "   cp terraform.tfvars.example terraform.tfvars"
echo "   # Edit terraform.tfvars with your values"
echo ""
echo "3. Deploy infrastructure:"
echo "   ./deployment/scripts/quick-deploy.sh"
echo "=========================================="

