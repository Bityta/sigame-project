# ==============================================
# Yandex Cloud Monitoring & Logging
# Auto-generated dashboards for SIGame services
# ==============================================

# Log Group for centralized logging
resource "yandex_logging_group" "sigame_logs" {
  name             = "sigame-logs-${var.environment}"
  folder_id        = var.folder_id
  retention_period = "259200s" # 3 days

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}

# Dashboard 1: Infrastructure Overview
resource "yandex_monitoring_dashboard" "infrastructure" {
  name        = "sigame-infrastructure-overview"
  description = "VM metrics: CPU, Memory, Disk, Network, Docker status"
  folder_id   = var.folder_id

  labels = {
    project = var.project_name
  }

  parametrization {}
}

# Dashboard 2: Auth Service
resource "yandex_monitoring_dashboard" "auth_service" {
  name        = "sigame-auth-service"
  description = "HTTP status codes (2xx/4xx/5xx), RPS, Latency, JWT operations, DB connections"
  folder_id   = var.folder_id

  labels = {
    project = var.project_name
    service = "auth"
  }

  parametrization {}
}

# Dashboard 3: Lobby Service
resource "yandex_monitoring_dashboard" "lobby_service" {
  name        = "sigame-lobby-service"
  description = "HTTP status codes, RPS, Latency, Active rooms, gRPC calls, DB connections"
  folder_id   = var.folder_id

  labels = {
    project = var.project_name
    service = "lobby"
  }

  parametrization {}
}

# Dashboard 4: Game Service
resource "yandex_monitoring_dashboard" "game_service" {
  name        = "sigame-game-service"
  description = "HTTP status codes, RPS, WebSocket connections, Game events, Latency, gRPC calls"
  folder_id   = var.folder_id

  labels = {
    project = var.project_name
    service = "game"
  }

  parametrization {}
}

# Dashboard 5: Pack Service
resource "yandex_monitoring_dashboard" "pack_service" {
  name        = "sigame-pack-service"
  description = "HTTP status codes, RPS, Latency, gRPC calls, DB query latency"
  folder_id   = var.folder_id

  labels = {
    project = var.project_name
    service = "pack"
  }

  parametrization {}
}

# Dashboard 6: Infrastructure Services
resource "yandex_monitoring_dashboard" "infra_services" {
  name        = "sigame-infrastructure-services"
  description = "PostgreSQL, Redis, Kafka metrics"
  folder_id   = var.folder_id

  labels = {
    project   = var.project_name
    component = "infrastructure"
  }

  parametrization {}
}

