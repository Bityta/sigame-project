# Security Group for Application Server
resource "yandex_vpc_security_group" "app_server_sg" {
  name        = "${var.project_name}-app-sg-${var.environment}"
  description = "Security group for application server"
  network_id  = yandex_vpc_network.sigame_network.id

  # SSH
  ingress {
    protocol       = "TCP"
    description    = "SSH"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 22
  }

  # HTTP
  ingress {
    protocol       = "TCP"
    description    = "HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 80
  }

  # HTTPS
  ingress {
    protocol       = "TCP"
    description    = "HTTPS"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 443
  }

  # Auth Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Auth Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8081
  }

  # Lobby Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Lobby Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8082
  }

  # Game Service HTTP & WebSocket
  ingress {
    protocol       = "TCP"
    description    = "Game Service HTTP & WebSocket"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8083
  }

  # Pack Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Pack Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8084
  }

  # Grafana
  ingress {
    protocol       = "TCP"
    description    = "Grafana"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 3000
  }

  # JVM Remote Debug (Lobby Service)
  ingress {
    protocol       = "TCP"
    description    = "JVM Remote Debug"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 5005
  }

  # Allow all outgoing traffic
  egress {
    protocol       = "ANY"
    description    = "Allow all outgoing"
    v4_cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow traffic from infra subnet
  ingress {
    protocol       = "ANY"
    description    = "Allow from infrastructure subnet"
    v4_cidr_blocks = ["10.129.0.0/24"]
  }
}


