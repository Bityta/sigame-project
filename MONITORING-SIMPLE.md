# 📊 Yandex Cloud Monitoring & Logging - Упрощенная версия

## 🎯 Текущая настройка

### ✅ Что УЖЕ работает

1. **Docker JSON Logs** - все логи сохраняются локально:
   ```bash
   # Смотреть логи любого сервиса
   ssh ubuntu@89.169.139.21
   sudo docker logs -f sigame-auth-service
   sudo docker logs -f sigame-lobby-service --since 10m
   
   # Поиск по логам
   sudo docker logs sigame-lobby-service 2>&1 | grep "ERROR"
   sudo docker logs sigame-lobby-service 2>&1 | grep "trace_id"
   ```

2. **Yandex Compute Cloud Metrics** (автоматически):
   - CPU, RAM, Disk, Network метрики VM
   - **URL**: https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/compute/instances

3. **JSON логи с trace_id** (настроены в приложениях):
   ```json
   {
     "timestamp": "2025-11-24T02:30:00Z",
     "level": "INFO",
     "trace_id": "abc123-def456",
     "service": "lobby-service",
     "message": "Room created"
   }
   ```

---

## 📝 Просмотр логов

### Логи на сервере

```bash
# SSH на сервер
ssh ubuntu@89.169.139.21

# Все сервисы
sudo docker ps --format "table {{.Names}}\t{{.Status}}"

# Логи конкретного сервиса
sudo docker logs -f --tail 100 sigame-auth-service
sudo docker logs -f --tail 100 sigame-lobby-service
sudo docker logs -f --tail 100 sigame-game-service
sudo docker logs -f --tail 100 sigame-pack-service

# Поиск ошибок за последний час
sudo docker logs --since 1h sigame-lobby-service 2>&1 | grep -E "ERROR|Exception"

# Поиск по trace_id
sudo docker logs sigame-lobby-service 2>&1 | grep "trace_id:abc123"

# Все логи всех контейнеров
sudo docker compose -f /opt/sigame/docker-compose.app.yml logs --tail 50
```

### Фильтры

```bash
# Только ошибки
sudo docker logs sigame-lobby-service 2>&1 | jq 'select(.level=="ERROR")'

# По пользователю
sudo docker logs sigame-lobby-service 2>&1 | grep "user_id:550e8400"

# За период
sudo docker logs --since "2025-11-24T00:00:00" --until "2025-11-24T02:00:00" sigame-lobby-service
```

---

## 📊 Monitoring (Yandex Cloud)

### 1. VM Метрики (встроенные)

https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/compute/instances

**Доступно:**
- CPU Usage (%)
- Memory Usage (%)
- Disk IOPS
- Network RX/TX

### 2. Custom Dashboards

Создай дашборд вручную:

1. Открой **Monitoring** → **Dashboards** → **Create Dashboard**
2. Добавь виджеты:
   - **CPU**: `compute.googleapis.com/instance/cpu/utilization`
   - **Memory**: `compute.googleapis.com/instance/memory/utilization`
   - **Disk**: `compute.googleapis.com/instance/disk/read_bytes_count`
   - **Network**: `compute.googleapis.com/instance/network/received_bytes_count`

---

## 🔍 Поиск проблем

### Типичные команды

```bash
# Последние ошибки из всех сервисов
for service in auth lobby game pack; do
  echo "=== $service-service ==="
  sudo docker logs --tail 20 sigame-$service-service 2>&1 | grep -i error
  echo ""
done

# Проверка health всех сервисов
sudo docker ps --format "table {{.Names}}\t{{.Status}}" | grep sigame

# Рестарт упавшего сервиса
sudo docker restart sigame-lobby-service

# Просмотр ресурсов
sudo docker stats --no-stream

# Дисковое пространство
df -h
sudo docker system df
```

---

## ⚠️ Алерты (ручная настройка)

### Создание алерта в Yandex Cloud

1. **Monitoring** → **Alerts** → **Create Alert**
2. Настрой условия:
   - **High CPU**: `cpu_utilization > 80%` для 5 минут
   - **Low Disk**: `disk_free_bytes < 1GB`
   - **Service Down**: `container_status != running`
3. Добавь канал уведомлений (Email / Telegram)

---

## 🚀 Полезные скрипты

### Скрипт для мониторинга логов

```bash
# Создай /opt/sigame/scripts/tail-errors.sh
#!/bin/bash
while true; do
  clear
  echo "=========================================="
  echo "  SIGAME ERRORS (last 5 min)"
  echo "=========================================="
  for svc in auth lobby game pack; do
    echo ""
    echo "[$svc-service]"
    sudo docker logs --since 5m sigame-$svc-service 2>&1 | grep -E "ERROR|Exception" | tail -5
  done
  sleep 30
done
```

### Скрипт для проверки health

```bash
# /opt/sigame/scripts/check-health.sh
#!/bin/bash
echo "Checking services health..."
curl -s http://localhost:8081/health | jq
curl -s http://localhost:8082/api/lobby/health | jq
curl -s http://localhost:8083/health | jq
curl -s http://localhost:8084/health | jq
```

---

## 📌 Важные URL

| Ресурс | URL |
|--------|-----|
| Application | http://89.169.139.21 |
| VM Metrics | https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/compute |
| Monitoring Dashboards | https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/monitoring |
| Cloud Logging (если настроишь) | https://console.cloud.yandex.ru/folders/b1g79ef2i8m53bbrbjru/logs |

---

## 💡 Рекомендации

1. **Для просмотра логов** - используй `docker logs` напрямую
2. **Для метрик VM** - используй Yandex Compute Cloud Dashboard
3. **Для custom метрик** - интегрируй Prometheus (опционально позже)
4. **Для централизованных логов** - настрой Yandex Cloud Logging через UI (вручную создай Log Group и подключи VM)

---

**Готово!** 🎉 Упрощенный, но рабочий мониторинг без сложной настройки Unified Agent.

