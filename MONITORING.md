# 📊 Monitoring & Logging - Yandex Cloud

Полное руководство по мониторингу и логированию SIGame через Yandex Cloud.

## 🎯 Обзор

Используем встроенные инструменты Yandex Cloud:
- **Yandex Cloud Logging** - централизованные логи с trace_id
- **Yandex Monitoring** - метрики и дашборды
- **Prometheus** - сбор метрик приложений
- **Fluent Bit** - отправка логов в облако

---

## 📝 Логи

### Просмотр логов в Yandex Cloud Console

**URL**: https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/logs

### Поиск логов

**Примеры запросов:**

```
# По trace_id
trace_id:"abc123-def456-789"

# По сервису
service:"lobby-service"

# Только ошибки
level:"ERROR"

# Комбинированный поиск
service:"auth-service" AND level:"ERROR" AND trace_id:"*"

# За период времени
service:"game-service" timestamp>="2025-11-24T00:00:00Z"
```

### Просмотр через CLI

```bash
# Все логи за последний час
./scripts/view-logs.sh

# Логи конкретного сервиса
./scripts/view-logs.sh lobby-service
./scripts/view-logs.sh auth-service
./scripts/view-logs.sh game-service
./scripts/view-logs.sh pack-service

# С фильтром напрямую
yc logging read \
  --group-id=<LOG_GROUP_ID> \
  --filter="service='lobby-service' AND level='ERROR'" \
  --since=1h \
  --follow
```

---

## 📊 Дашборды

### Список дашбордов

После `terraform apply` будут созданы 6 дашбордов:

1. **SIGame - Infrastructure Overview**
   - CPU utilization (%)
   - Memory utilization (%)
   - Disk I/O (read/write bytes)
   - Network traffic (RX/TX)
   - Docker containers status

2. **SIGame - Auth Service**
   - HTTP Status Codes: 2xx, 4xx, 5xx (req/s)
   - RPS Total
   - RPS by Endpoint (/login, /register, /refresh, /validate)
   - Latency (p50, p95, p99) в миллисекундах
   - JWT operations rate
   - Database connections

3. **SIGame - Lobby Service**
   - HTTP Status Codes: 2xx, 4xx, 5xx (req/s)
   - RPS Total
   - RPS by Endpoint (/rooms, /rooms/{id}, /rooms/{id}/join)
   - Latency (p50, p95, p99)
   - Active rooms count
   - gRPC calls latency (to Auth, to Pack)
   - Database connections

4. **SIGame - Game Service**
   - HTTP Status Codes: 2xx, 4xx, 5xx (req/s)
   - RPS Total
   - RPS by Endpoint (/games/{id}, /games/{id}/state)
   - WebSocket connections count
   - Game events rate
   - Latency (p50, p95, p99)
   - gRPC calls to Pack service
   - Database connections

5. **SIGame - Pack Service**
   - HTTP Status Codes: 2xx, 4xx, 5xx (req/s)
   - RPS Total
   - RPS by Endpoint (/packs, /packs/{id})
   - Latency (p50, p95, p99)
   - gRPC calls from Lobby/Game
   - Database query latency

6. **SIGame - Infrastructure Services**
   - PostgreSQL: connections, query latency, transactions/s
   - Redis: operations/s, memory usage, hit rate
   - Kafka: messages/s, consumer lag, topics

### Открыть дашборды

```bash
# Открыть все дашборды в браузере
./scripts/open-dashboards.sh

# Получить URLs дашбордов
./scripts/get-monitoring-urls.sh
```

**Или вручную:**
https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/monitoring/dashboards

---

## 🚀 Развертывание

### 1. Terraform Apply

Создает Log Group и Dashboards:

```bash
cd deployment/terraform
terraform apply

# Получить LOG_GROUP_ID
terraform output -raw log_group_id
```

### 2. Настроить .env.production на сервере

```bash
ssh ubuntu@89.169.139.21

# Получить IAM token (действителен 12 часов)
yc iam create-token

# Добавить в .env.production
cd /opt/sigame
echo "YC_FOLDER_ID=b1g79ef2i8m53bbrbjru" >> .env.production
echo "LOG_GROUP_ID=<id_from_terraform_output>" >> .env.production
echo "YC_IAM_TOKEN=<token_from_yc_iam_create-token>" >> .env.production
```

### 3. Запустить Fluent Bit и Prometheus

```bash
cd /opt/sigame
sudo docker compose -f docker-compose.infra.yml up -d fluent-bit prometheus
```

### 4. Проверить работу

```bash
# Проверить статус контейнеров
sudo docker ps | grep -E "fluent-bit|prometheus"

# Проверить логи Fluent Bit
sudo docker logs sigame-fluent-bit --tail 50

# Проверить логи Prometheus
sudo docker logs sigame-prometheus --tail 50

# Через 5-10 минут проверить в Yandex Cloud Console
```

---

## 🔧 Troubleshooting

### Fluent Bit не отправляет логи

```bash
# Проверить конфигурацию
sudo docker exec sigame-fluent-bit cat /fluent-bit/etc/fluent-bit.conf

# Проверить переменные окружения
sudo docker exec sigame-fluent-bit env | grep YC_

# Проверить логи
sudo docker logs sigame-fluent-bit --tail 100

# Перезапустить
sudo docker restart sigame-fluent-bit
```

### Prometheus не отправляет метрики

```bash
# Проверить targets
curl http://localhost:9090/api/v1/targets | jq

# Проверить конфигурацию
sudo docker exec sigame-prometheus cat /etc/prometheus/prometheus.yml

# Проверить remote_write status
curl http://localhost:9090/api/v1/status/tsdb | jq

# Перезапустить
sudo docker restart sigame-prometheus
```

### IAM Token истек

IAM токен действителен 12 часов. Для продакшена используйте Service Account:

```bash
# Создать Service Account
yc iam service-account create --name sigame-monitoring

# Назначить роли
SA_ID=$(yc iam service-account get sigame-monitoring --format json | jq -r '.id')
yc resource-manager folder add-access-binding b1g79ef2i8m53bbrbjru \
  --role logging.writer \
  --service-account-id $SA_ID

# Создать authorized key
yc iam key create --service-account-id $SA_ID --output sa-key.json

# Использовать в docker-compose через volume
```

---

## 📈 Метрики

### Требуемые метрики в сервисах

Для корректной работы дашбордов, каждый сервис должен экспортировать:

**1. HTTP метрики:**
```
http_requests_total{method, endpoint, status}      # Counter
http_request_duration_seconds{method, endpoint}    # Histogram
```

**2. gRPC метрики (для Lobby/Game):**
```
grpc_client_requests_total{service, method, status}     # Counter
grpc_client_request_duration_seconds{service, method}   # Histogram
```

**3. Custom метрики:**
```
# Auth Service
jwt_operations_total{operation}                    # Counter (sign, validate, refresh)
active_sessions                                     # Gauge

# Lobby Service
active_rooms                                        # Gauge
players_online                                      # Gauge

# Game Service
active_games                                        # Gauge
websocket_connections                               # Gauge
game_events_total{type}                             # Counter
```

### Пример экспорта метрик (Go)

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}
```

---

## 💰 Стоимость

- **Log Group** 3 дня retention, до 50 ГБ/мес: **бесплатно**
- **Monitoring** dashboards: **бесплатно**
- **Prometheus remote_write**: **бесплатно** (в лимитах)
- **Fluent Bit**: ~100 МБ RAM

**Итого:** практически без дополнительных затрат

---

## 📚 Полезные ссылки

- [Yandex Cloud Logging Docs](https://cloud.yandex.ru/docs/logging/)
- [Yandex Monitoring Docs](https://cloud.yandex.ru/docs/monitoring/)
- [Fluent Bit Documentation](https://docs.fluentbit.io/)
- [Prometheus Remote Write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write)

---

**Готово!** 🎉 Все логи с trace_id и дашборды с 2xx/4xx/5xx, RPS, latency готовы к использованию.

