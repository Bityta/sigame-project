variable "yandex_cloud_token" {
  description = "OAuth токен для доступа к Yandex Cloud"
  type        = string
  sensitive   = true
}

variable "cloud_id" {
  description = "ID облака в Yandex Cloud"
  type        = string
}

variable "folder_id" {
  description = "ID каталога в Yandex Cloud"
  type        = string
}

variable "zone" {
  description = "Зона доступности"
  type        = string
  default     = "ru-central1-a"
}

variable "ssh_public_key_path" {
  description = "Путь к публичному SSH ключу"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "project_name" {
  description = "Название проекта"
  type        = string
  default     = "sigame"
}

variable "environment" {
  description = "Окружение (dev, staging, prod)"
  type        = string
  default     = "prod"
}

variable "app_server_cores" {
  description = "Количество vCPU для Application сервера"
  type        = number
  default     = 4
}

variable "app_server_memory" {
  description = "Объем RAM для Application сервера (GB)"
  type        = number
  default     = 8
}

variable "app_server_disk_size" {
  description = "Размер диска для Application сервера (GB)"
  type        = number
  default     = 50
}

variable "infra_server_cores" {
  description = "Количество vCPU для Infrastructure сервера"
  type        = number
  default     = 4
}

variable "infra_server_memory" {
  description = "Объем RAM для Infrastructure сервера (GB)"
  type        = number
  default     = 16
}

variable "infra_server_disk_size" {
  description = "Размер диска для Infrastructure сервера (GB)"
  type        = number
  default     = 100
}

