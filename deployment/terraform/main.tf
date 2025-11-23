terraform {
  required_providers {
    yandex = {
      source  = "yandex-cloud/yandex"
      version = "~> 0.100"
    }
  }
  required_version = ">= 1.0"
}

provider "yandex" {
  token     = var.yandex_cloud_token
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.zone
}

# VPC Network
resource "yandex_vpc_network" "sigame_network" {
  name        = "${var.project_name}-network-${var.environment}"
  description = "Network for SIGame project"
}

# Subnet for Application Server
resource "yandex_vpc_subnet" "app_subnet" {
  name           = "${var.project_name}-app-subnet-${var.environment}"
  zone           = var.zone
  network_id     = yandex_vpc_network.sigame_network.id
  v4_cidr_blocks = ["10.128.0.0/24"]
  description    = "Subnet for all services (single server setup)"
}

# Application Server (all-in-one)
resource "yandex_compute_instance" "app_server" {
  name        = "${var.project_name}-server-${var.environment}"
  platform_id = "standard-v3"
  zone        = var.zone

  resources {
    cores  = var.app_server_cores
    memory = var.app_server_memory
  }

  boot_disk {
    initialize_params {
      image_id = "fd8kdq6d0p8sij7h5qe3" # Ubuntu 22.04 LTS
      size     = var.app_server_disk_size
      type     = "network-ssd"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.app_subnet.id
    nat       = true # External IP
    security_group_ids = [
      yandex_vpc_security_group.app_server_sg.id
    ]
  }

  metadata = {
    ssh-keys = "ubuntu:${file(var.ssh_public_key_path)}"
    user-data = templatefile("${path.module}/cloud-init-app.yml", {
      hostname = "${var.project_name}-server"
    })
  }

  scheduling_policy {
    preemptible = false
  }

  labels = {
    project     = var.project_name
    environment = var.environment
    role        = "all-in-one"
  }
}

