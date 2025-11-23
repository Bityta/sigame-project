# 📊 Yandex Cloud Monitoring & Logging

## 🎯 Обзор

Используем встроенные инструменты Yandex Cloud вместо Grafana/Loki/Tempo:

- **Yandex Monitoring** - метрики, дашборды, алерты
- **Yandex Cloud Logging** - централизованные логи
- **Unified Agent** - сбор метрик и логов с сервера

---

## 📈 Доступ к мониторингу

### 1. Yandex Monitoring
```
https://console.cloud.yandex.ru/folders/<folder_id>/monitoring
```

**Доступные метрики:**
- CPU, RAM, Network, Disk I/O (VM)
- Метрики приложений (через Unified Agent)
- Request rate, latency, errors (по сервисам)
- Database connections, query time
- Kafka lag, Redis operations

### 2. Yandex Cloud Logging
```
https://console.cloud.yandex.ru/folders/<folder_id>/logs
```

**Возможности:**
- Централизованные логи всех Docker контейнеров
- Поиск по trace_id, user_id, request_id
- Фильтрация по уровню (ERROR, WARN, INFO)
- Экспорт логов в Object Storage

### 3. Compute Cloud (VM метрики)
```
https://console.cloud.yandex.ru/folders/<folder_id>/compute/instances
```

**Показывает:**
- CPU usage
- Memory usage
- Disk IOPS
- Network bandwidth

---

## 🔧 Настройка Unified Agent

Unified Agent уже настроен в `docker-compose.infra.yml` и собирает:

1. **Docker логи** - все контейнеры с префиксом `sigame-*`
2. **Метрики сервисов:**
   - Auth Service: `http://localhost:8081/metrics`
   - Lobby Service: `http://localhost:8082/actuator/prometheus`
   - Game Service: `http://localhost:8083/metrics`
   - Pack Service: `http://localhost:8084/metrics`

### Переменные окружения

Добавь в `.env.production`:

```bash
# Yandex Cloud Monitoring
YC_FOLDER_ID=<твой-folder-id>
```

**Получить Folder ID:**
```bash
yc config get folder-id
```

---

## 📊 Создание дашбордов

### 1. Через Yandex Cloud Console

1. Открой **Monitoring** → **Dashboards**
2. Создай новый дашборд
3. Добавь виджеты для нужных метрик

### 2. Terraform (автоматизация)

Создам конфигурацию для автоматического создания дашбордов.

---

## 🔍 Поиск логов с Trace ID

### Формат логов

Все сервисы логируют в JSON с полями:

```json
{
  "timestamp": "2025-11-23T12:34:56.789Z",
  "level": "INFO",
  "trace_id": "abc123-def456-789",
  "span_id": "xyz789",
  "service": "lobby-service",
  "message": "Room created",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Поиск в Yandex Cloud Logging

**По Trace ID:**
```
trace_id:"abc123-def456-789"
```

**По пользователю:**
```
user_id:"550e8400-e29b-41d4-a716-446655440000"
```

**Только ошибки:**
```
level:"ERROR"
```

**Комбинированный поиск:**
```
service:"lobby-service" AND level:"ERROR" AND trace_id:"abc123*"
```

---

## ⚠️ Алерты

### Настройка алертов

1. **Monitoring** → **Alerts**
2. Создай алерт на:
   - **High error rate** - `http_requests_total{status=~"5.."} > 10`
   - **High latency** - `http_request_duration_seconds{quantile="0.95"} > 1`
   - **Low disk space** - `disk_free_bytes < 1GB`
   - **High CPU** - `cpu_usage > 80%`
   - **Database connections** - `db_connections > 80`

3. Укажи канал уведомлений (Telegram, Email)

---

## 📦 Удаление старых компонентов

Удалены следующие сервисы:
- ❌ Grafana
- ❌ Prometheus
- ❌ Loki
- ❌ Tempo
- ❌ Promtail

Освобождено ~2GB RAM и ~10GB диска.

---

## 🚀 Деплой изменений

```bash
# На сервере
cd /opt/sigame
git pull origin main
sudo docker compose -f docker-compose.infra.yml down grafana prometheus loki tempo promtail
sudo docker compose -f docker-compose.infra.yml -f docker-compose.app.yml up -d --build
```

---

## 📌 Полезные ссылки

- [Yandex Monitoring Docs](https://cloud.yandex.ru/docs/monitoring/)
- [Yandex Cloud Logging](https://cloud.yandex.ru/docs/logging/)
- [Unified Agent Configuration](https://cloud.yandex.ru/docs/monitoring/concepts/data-collection/unified-agent)
- [PromQL Query Examples](https://prometheus.io/docs/prometheus/latest/querying/examples/)

---

## 🎯 Быстрый доступ

| Сервис | URL |
|--------|-----|
| Monitoring | https://console.cloud.yandex.ru/monitoring |
| Logs | https://console.cloud.yandex.ru/logs |
| VM Metrics | https://console.cloud.yandex.ru/compute |
| Application | http://89.169.139.21 |

---

**Все готово!** 🎉 Теперь используем только инструменты Yandex Cloud.

