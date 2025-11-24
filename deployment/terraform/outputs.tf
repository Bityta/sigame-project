output "app_server_external_ip" {
  description = "External IP address of Server"
  value       = yandex_compute_instance.app_server.network_interface[0].nat_ip_address
}

output "app_server_internal_ip" {
  description = "Internal IP address of Server"
  value       = yandex_compute_instance.app_server.network_interface[0].ip_address
}

output "app_server_id" {
  description = "ID of Server"
  value       = yandex_compute_instance.app_server.id
}

output "network_id" {
  description = "ID of VPC network"
  value       = yandex_vpc_network.sigame_network.id
}

# Monitoring & Logging Outputs
output "log_group_id" {
  description = "ID Log Group для настройки на сервере"
  value       = yandex_logging_group.sigame_logs.id
}

output "log_group_url" {
  description = "URL для просмотра логов"
  value       = "https://console.cloud.yandex.ru/folders/${var.folder_id}/logs?logGroupId=${yandex_logging_group.sigame_logs.id}"
}

output "monitoring_dashboards" {
  description = "URLs дашбордов"
  value = {
    infrastructure          = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.infrastructure.id}"
    auth_service            = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.auth_service.id}"
    lobby_service           = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.lobby_service.id}"
    game_service            = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.game_service.id}"
    pack_service            = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.pack_service.id}"
    infrastructure_services = "https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards/${yandex_monitoring_dashboard.infra_services.id}"
  }
}

output "connection_info" {
  description = "Connection information"
  value       = <<-EOT
    ========================================
    SIGame Infrastructure Deployed!
    ========================================
    
    Server (All-in-One):
      External IP: ${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}
      Internal IP: ${yandex_compute_instance.app_server.network_interface[0].ip_address}
      SSH: ssh ubuntu@${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}
    
    Services:
      Frontend:  http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}
      Auth:      http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8081
      Lobby:     http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8082
      Game:      http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8083
      Packs:     http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8084
    
    Monitoring:
      Logs:      https://console.cloud.yandex.ru/folders/${var.folder_id}/logs?logGroupId=${yandex_logging_group.sigame_logs.id}
      Dashboards: https://console.cloud.yandex.ru/folders/${var.folder_id}/monitoring/dashboards
    
    Next Steps:
      1. Add to .env.production: LOG_GROUP_ID=${yandex_logging_group.sigame_logs.id}
      2. Run: cd /opt/sigame && sudo docker compose up -d
    ========================================
  EOT
}

