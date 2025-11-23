output "app_server_external_ip" {
  description = "External IP address of Application Server"
  value       = yandex_compute_instance.app_server.network_interface[0].nat_ip_address
}

output "app_server_internal_ip" {
  description = "Internal IP address of Application Server"
  value       = yandex_compute_instance.app_server.network_interface[0].ip_address
}

output "infra_server_internal_ip" {
  description = "Internal IP address of Infrastructure Server"
  value       = yandex_compute_instance.infra_server.network_interface[0].ip_address
}

output "app_server_id" {
  description = "ID of Application Server"
  value       = yandex_compute_instance.app_server.id
}

output "infra_server_id" {
  description = "ID of Infrastructure Server"
  value       = yandex_compute_instance.infra_server.id
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
    
    Application Server:
      External IP: ${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}
      Internal IP: ${yandex_compute_instance.app_server.network_interface[0].ip_address}
      SSH: ssh ubuntu@${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}
    
    Infrastructure Server:
      Internal IP: ${yandex_compute_instance.infra_server.network_interface[0].ip_address}
      SSH (via App Server): ssh -J ubuntu@${yandex_compute_instance.app_server.network_interface[0].nat_ip_address} ubuntu@${yandex_compute_instance.infra_server.network_interface[0].ip_address}
    
    Services:
      Frontend:  http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:3001
      Auth:      http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8001
      Lobby:     http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8002
      Game:      http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8003
      Packs:     http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:8005
      Grafana:   http://${yandex_compute_instance.app_server.network_interface[0].nat_ip_address}:3000
    
    Next Steps:
      1. Save these IPs to .env.production
      2. Run: ./deployment/scripts/setup-server1.sh
      3. Run: ./deployment/scripts/setup-server2.sh
    ========================================
  EOT
}

