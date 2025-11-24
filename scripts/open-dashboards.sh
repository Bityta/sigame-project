#!/bin/bash
# Open all monitoring dashboards in browser
# Usage: ./scripts/open-dashboards.sh

set -e

cd "$(dirname "$0")/../deployment/terraform"

echo "Opening monitoring dashboards..."
echo ""

# Get dashboard URLs from terraform output
DASHBOARDS=$(terraform output -json monitoring_dashboards 2>/dev/null)

if [ -z "$DASHBOARDS" ] || [ "$DASHBOARDS" = "null" ]; then
  echo "Error: Unable to get dashboard URLs from terraform output"
  echo "Run 'terraform apply' first"
  exit 1
fi

# Parse and open each dashboard
echo "$DASHBOARDS" | jq -r '.[]' | while read -r url; do
  echo "Opening: $url"
  open "$url" 2>/dev/null || xdg-open "$url" 2>/dev/null || echo "Could not open browser automatically. Please visit: $url"
  sleep 1
done

echo ""
echo "All dashboards opened!"

