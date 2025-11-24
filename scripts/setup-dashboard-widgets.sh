#!/bin/bash
# Script to add widgets to Yandex Monitoring dashboards via API
# This creates charts for 2xx/4xx/5xx, RPS, and Latency

set -e

FOLDER_ID="b1g79ef2i8m53bbrbjru"

echo "Adding widgets to dashboards..."
echo "This requires 'yc' CLI to be configured"
echo ""

# Get dashboard IDs from terraform
cd "$(dirname "$0")/../../deployment/terraform"

AUTH_DASHBOARD_ID=$(terraform output -json monitoring_dashboards | jq -r '.auth_service' | grep -o 'fbe[a-z0-9]*$')
LOBBY_DASHBOARD_ID=$(terraform output -json monitoring_dashboards | jq -r '.lobby_service' | grep -o 'fbe[a-z0-9]*$')
GAME_DASHBOARD_ID=$(terraform output -json monitoring_dashboards | jq -r '.game_service' | grep -o 'fbe[a-z0-9]*$')
PACK_DASHBOARD_ID=$(terraform output -json monitoring_dashboards | jq -r '.pack_service' | grep -o 'fbe[a-z0-9]*$')
INFRA_DASHBOARD_ID=$(terraform output -json monitoring_dashboards | jq -r '.infrastructure' | grep -o 'fbe[a-z0-9]*$')

echo "Dashboard IDs:"
echo "  Auth: $AUTH_DASHBOARD_ID"
echo "  Lobby: $LOBBY_DASHBOARD_ID"
echo "  Game: $GAME_DASHBOARD_ID"
echo "  Pack: $PACK_DASHBOARD_ID"
echo "  Infrastructure: $INFRA_DASHBOARD_ID"
echo ""

# Note: Adding widgets requires using Yandex Cloud API directly
# The yc CLI doesn't have direct commands for dashboard widgets
# You need to use the REST API or do it manually in the UI

echo "⚠️  Adding widgets via Terraform/CLI is not supported by Yandex provider"
echo ""
echo "📝 To add widgets, please use Yandex Cloud Console:"
echo "   1. Open: https://console.cloud.yandex.ru/folders/$FOLDER_ID/monitoring/dashboards"
echo "   2. Click on a dashboard"
echo "   3. Click 'Edit' → 'Add widget'"
echo "   4. Add charts for:"
echo "      - HTTP requests by status code (2xx, 4xx, 5xx)"
echo "      - Requests per second (RPS)"
echo "      - Response latency (p50, p95, p99)"
echo ""
echo "📊 Recommended queries:"
echo ""
echo "HTTP Status 2xx:"
echo '  rate(http_requests_total{service="auth-service",status=~"2.."}[1m])'
echo ""
echo "HTTP Status 4xx:"
echo '  rate(http_requests_total{service="auth-service",status=~"4.."}[1m])'
echo ""
echo "HTTP Status 5xx:"
echo '  rate(http_requests_total{service="auth-service",status=~"5.."}[1m])'
echo ""
echo "RPS Total:"
echo '  sum(rate(http_requests_total{service="auth-service"}[1m]))'
echo ""
echo "Latency p95:"
echo '  histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{service="auth-service"}[5m]))'
echo ""

