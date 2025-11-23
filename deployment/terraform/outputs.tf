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

output "connection_info" {
  description = "Connection information"
  value = <<-EOT
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
      Grafana:   http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:3000
    
    Next Steps:
      1. Update .env.production with this IP
      2. Run: cd /opt/sigame && sudo docker compose up -d
    ========================================
  EOT
}

