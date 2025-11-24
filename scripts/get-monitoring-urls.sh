#!/bin/bash
# Get all monitoring URLs
# Usage: ./scripts/get-monitoring-urls.sh

set -e

cd "$(dirname "$0")/../deployment/terraform"

echo "=========================================="
echo "  Yandex Cloud Monitoring URLs"
echo "=========================================="
echo ""

echo "📝 Logs:"
LOG_URL=$(terraform output -raw log_group_url 2>/dev/null)
if [ -n "$LOG_URL" ]; then
  echo "   $LOG_URL"
else
  echo "   Not configured yet. Run 'terraform apply' first."
fi

echo ""
echo "📊 Dashboards:"
DASHBOARDS=$(terraform output -json monitoring_dashboards 2>/dev/null)
if [ -n "$DASHBOARDS" ] && [ "$DASHBOARDS" != "null" ]; then
  echo "$DASHBOARDS" | jq -r 'to_entries[] | "   \(.key): \(.value)"'
else
  echo "   Not configured yet. Run 'terraform apply' first."
fi

echo ""
echo "=========================================="

