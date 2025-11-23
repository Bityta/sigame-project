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
    port           = 8001
  }

  # Lobby Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Lobby Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8002
  }

  # Game Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Game Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8003
  }

  # Game Service WebSocket
  ingress {
    protocol       = "TCP"
    description    = "Game Service WebSocket"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8083
  }

  # Pack Service HTTP
  ingress {
    protocol       = "TCP"
    description    = "Pack Service HTTP"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 8005
  }

  # Frontend (Nginx)
  ingress {
    protocol       = "TCP"
    description    = "Frontend Nginx"
    v4_cidr_blocks = ["0.0.0.0/0"]
    port           = 3001
  }

  # Allow all outgoing traffic
  egress {
    protocol       = "ANY"
    description    = "Allow all outgoing"
    v4_cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow traffic from infra subnet
  ingress {
    protocol          = "ANY"
    description       = "Allow from infrastructure subnet"
    v4_cidr_blocks    = ["10.129.0.0/24"]
  }
}

# Security Group for Infrastructure Server
resource "yandex_vpc_security_group" "infra_server_sg" {
  name        = "${var.project_name}-infra-sg-${var.environment}"
  description = "Security group for infrastructure server"
  network_id  = yandex_vpc_network.sigame_network.id

  # SSH from application server only
  ingress {
    protocol       = "TCP"
    description    = "SSH from app subnet"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 22
  }

  # PostgreSQL Auth
  ingress {
    protocol       = "TCP"
    description    = "PostgreSQL Auth"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 5432
  }

  # PostgreSQL Lobby
  ingress {
    protocol       = "TCP"
    description    = "PostgreSQL Lobby"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 5433
  }

  # PostgreSQL Packs
  ingress {
    protocol       = "TCP"
    description    = "PostgreSQL Packs"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 5434
  }

  # PostgreSQL Game
  ingress {
    protocol       = "TCP"
    description    = "PostgreSQL Game"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 5435
  }

  # Redis
  ingress {
    protocol       = "TCP"
    description    = "Redis"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 6379
  }

  # Kafka
  ingress {
    protocol       = "TCP"
    description    = "Kafka"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 9092
  }

  # MinIO API
  ingress {
    protocol       = "TCP"
    description    = "MinIO API"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 9000
  }

  # MinIO Console
  ingress {
    protocol       = "TCP"
    description    = "MinIO Console"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 9001
  }

  # Prometheus
  ingress {
    protocol       = "TCP"
    description    = "Prometheus"
    v4_cidr_blocks = ["10.128.0.0/24", "0.0.0.0/0"]
    port           = 9090
  }

  # Grafana
  ingress {
    protocol       = "TCP"
    description    = "Grafana"
    v4_cidr_blocks = ["10.128.0.0/24", "0.0.0.0/0"]
    port           = 3000
  }

  # Loki
  ingress {
    protocol       = "TCP"
    description    = "Loki"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 3100
  }

  # Tempo HTTP
  ingress {
    protocol       = "TCP"
    description    = "Tempo HTTP"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 3200
  }

  # Tempo OTLP gRPC
  ingress {
    protocol       = "TCP"
    description    = "Tempo OTLP gRPC"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 4317
  }

  # Tempo OTLP HTTP
  ingress {
    protocol       = "TCP"
    description    = "Tempo OTLP HTTP"
    v4_cidr_blocks = ["10.128.0.0/24"]
    port           = 4318
  }

  # Allow all outgoing traffic
  egress {
    protocol       = "ANY"
    description    = "Allow all outgoing"
    v4_cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow traffic from app subnet
  ingress {
    protocol       = "ANY"
    description    = "Allow from application subnet"
    v4_cidr_blocks = ["10.128.0.0/24"]
  }
}

