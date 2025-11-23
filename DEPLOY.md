# 🚀 Быстрый деплой на сервер

## 1. Как обновить код на сервере

После любых изменений в коде просто выполни:

```bash
bash deployment/scripts/deploy.sh
```

Скрипт автоматически:
- Закоммитит и запушит изменения
- Подключится к серверу
- Подтянет последний код
- Пересоберёт и перезапустит контейнеры
- Покажет статус всех сервисов

## 2. Проверка работы

После деплоя:

1. **Фронтенд**: http://89.169.139.21
2. **API эндпоинты**:
   - Auth: http://89.169.139.21:8081
   - Lobby: http://89.169.139.21:8082
   - Game: http://89.169.139.21:8083
   - Pack: http://89.169.139.21:8084

## 3. Просмотр логов

```bash
# Все контейнеры
ssh ubuntu@89.169.139.21 'sudo docker ps'

# Логи конкретного сервиса
ssh ubuntu@89.169.139.21 'sudo docker logs sigame-auth -f'
ssh ubuntu@89.169.139.21 'sudo docker logs sigame-frontend -f'
```

## 4. Ручной перезапуск

Если что-то пошло не так:

```bash
ssh ubuntu@89.169.139.21
cd /opt/sigame

# Перезапуск всех сервисов
sudo docker compose -f docker-compose.app.yml --env-file .env.production restart

# Перезапуск конкретного сервиса
sudo docker restart sigame-auth
```

## 5. Полная переустановка

Если нужно начать с чистого листа:

```bash
ssh ubuntu@89.169.139.21
cd /opt/sigame

# Остановка всех контейнеров
sudo docker compose -f docker-compose.infra.yml down -v
sudo docker compose -f docker-compose.app.yml down -v

# Очистка
sudo docker system prune -af --volumes

# Запуск инфраструктуры
sudo docker compose -f docker-compose.infra.yml --env-file .env.production up -d postgres-auth postgres-lobby postgres-game postgres-packs redis minio

# Подождать 30 секунд

# Запуск приложений
sudo docker compose -f docker-compose.app.yml --env-file .env.production up -d --build
```

## 6. Структура проекта на сервере

```
/opt/sigame/
├── .env.production          # Переменные окружения
├── docker-compose.app.yml   # Приложения
├── docker-compose.infra.yml # Инфраструктура (БД, Redis, MinIO)
├── services/                # Исходники сервисов
├── frontend/                # Фронтенд
└── deployment/              # Скрипты деплоя
```

## 7. Конфигурация портов

- **80** - Frontend (Nginx)
- **8081** - Auth Service (HTTP)
- **8082** - Lobby Service (HTTP)
- **8083** - Game Service (HTTP + WebSocket)
- **8084** - Pack Service (HTTP)
- **5432-5435** - PostgreSQL (4 базы данных)
- **6379** - Redis
- **9000-9001** - MinIO (S3-совместимое хранилище)

## 8. Troubleshooting

### Сервис не запускается

```bash
# Проверить логи
sudo docker logs sigame-auth --tail 50

# Проверить, что БД работают
sudo docker ps | grep postgres

# Перезапустить инфраструктуру
sudo docker compose -f docker-compose.infra.yml restart
```

### Frontend показывает ошибки подключения

```bash
# Проверить, что все API сервисы запущены
sudo docker ps | grep -E "auth|lobby|game|pack"

# Проверить переменные окружения frontend
sudo docker exec sigame-frontend env | grep VITE
```

### База данных не подключается

```bash
# Проверить статус PostgreSQL
sudo docker logs sigame-postgres-auth --tail 30

# Войти в PostgreSQL
sudo docker exec -it sigame-postgres-auth psql -U authuser -d auth_db
```

## 9. Мониторинг

Доступ к метрикам и логам (если развернута полная инфраструктура):

- **Grafana**: http://89.169.139.21:3001
  - Логин: admin
  - Пароль: смотри в `.env.production`

- **Prometheus**: http://89.169.139.21:9090

---

**Сервер**: 89.169.139.21  
**SSH**: `ssh ubuntu@89.169.139.21`  
**Проект**: `/opt/sigame`

