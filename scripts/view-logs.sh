#!/bin/bash
# View logs from Yandex Cloud Logging
# Usage: ./scripts/view-logs.sh [service-name]

set -e

SERVICE=${1:-"all"}

# Get log group ID from terraform
cd "$(dirname "$0")/../deployment/terraform"
LOG_GROUP_ID=$(terraform output -raw log_group_id 2>/dev/null)

if [ -z "$LOG_GROUP_ID" ]; then
  echo "Error: Unable to get log_group_id from terraform output"
  echo "Run 'terraform apply' first"
  exit 1
fi

echo "Viewing logs for service: $SERVICE"
echo "Log Group ID: $LOG_GROUP_ID"
echo "Press Ctrl+C to stop"
echo ""

if [ "$SERVICE" = "all" ]; then
  yc logging read \
    --group-id="$LOG_GROUP_ID" \
    --since=1h \
    --follow
else
  yc logging read \
    --group-id="$LOG_GROUP_ID" \
    --filter="service='$SERVICE'" \
    --since=1h \
    --follow
fi

