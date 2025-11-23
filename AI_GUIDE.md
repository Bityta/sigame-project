# AI Assistant Guide - SIGame Project

Краткая инструкция для AI ассистентов по работе с проектом.

## 🔐 Доступ к серверу

```bash
# SSH подключение
ssh ubuntu@89.169.139.21

# Рабочая директория проекта
cd /opt/sigame
```

## 📦 Git workflow

```bash
# Проверка статуса
git status

# Получение изменений
git fetch origin
git pull origin main

# Создание новой ветки для фичи
git checkout -b feature/your-feature-name

# Добавление изменений
git add <files>
git commit -m "feat: описание изменений"

# Пуш в GitHub
git push origin feature/your-feature-name

# После мерджа в main - автоматический деплой через GitHub Actions
```

## 🚀 Деплой и Docker

### Быстрый деплой (на сервере)
```bash
cd /opt/sigame

# Получить последние изменения
git pull origin main

# Перезапустить все сервисы
sudo docker compose -f docker-compose.infra.yml -f docker-compose.app.yml --env-file .env.production down
sudo docker compose -f docker-compose.infra.yml --env-file .env.production up -d
sudo docker compose -f docker-compose.app.yml --env-file .env.production up -d
```

### Пересборка конкретного сервиса
```bash
# Пример для Lobby
sudo docker compose -f docker-compose.app.yml --env-file .env.production build lobby-service
sudo docker compose -f docker-compose.app.yml --env-file .env.production up -d --force-recreate lobby-service
```

### Проверка логов
```bash
# Все контейнеры
sudo docker ps

# Логи конкретного сервиса
sudo docker logs sigame-lobby-service --tail 50
sudo docker logs sigame-auth-service --tail 50
sudo docker logs sigame-game-service --tail 50
sudo docker logs sigame-pack-service --tail 50
sudo docker logs sigame-frontend --tail 50
```

## 🗄️ База данных

### Подключение к PostgreSQL
```bash
# Auth DB
sudo docker exec -it sigame-postgres-auth psql -U authuser -d auth_db

# Lobby DB
sudo docker exec -it sigame-postgres-lobby psql -U lobbyuser -d lobby_db

# Game DB
sudo docker exec -it sigame-postgres-game psql -U gameuser -d game_db

# Packs DB
sudo docker exec -it sigame-postgres-packs psql -U packsuser -d packs_db

# Полезные SQL команды
\dt                    # Список таблиц
\d table_name          # Структура таблицы
SELECT * FROM ...;     # Запросы
\q                     # Выход
```

### Redis
```bash
sudo docker exec -it sigame-redis redis-cli
# KEYS *
# GET key_name
# exit
```

## 🔧 Структура проекта

```
sigame-project/
├── frontend/              # React + TypeScript + Vite
├── services/
│   ├── auth/             # Go - JWT аутентификация
│   ├── lobby/            # Kotlin - управление комнатами
│   ├── game/             # Go - игровая логика + WebSocket
│   └── packs/            # Python - паки вопросов
├── infrastructure/        # Конфиги для Prometheus, Grafana, Loki, etc.
├── deployment/
│   ├── scripts/          # Bash скрипты
│   └── terraform/        # Yandex Cloud IaC
├── docker-compose.infra.yml   # PostgreSQL, Redis, Kafka, MinIO, мониторинг
├── docker-compose.app.yml     # Приложения
└── .env.production            # Продакшн переменные (НЕ в git!)
```

## 🌐 Endpoints

- **Frontend**: http://89.169.139.21/
- **Auth API**: http://89.169.139.21:8081
- **Lobby API**: http://89.169.139.21:8082
- **Game API**: http://89.169.139.21:8083
- **Pack API**: http://89.169.139.21:8084
- **Grafana**: http://89.169.139.21:3000

## 🔍 Диагностика проблем

### 1. CORS ошибки
- Проверь `services/lobby/src/main/kotlin/com/sigame/lobby/config/CorsConfig.kt`
- Проверь порядок фильтров (@Order аннотации)
- Убедись что IP сервера добавлен в allowedOrigins

### 2. gRPC DEADLINE_EXCEEDED
- **НЕ** устанавливай deadline на stub при инициализации
- Добавляй `.withDeadlineAfter()` к каждому вызову:
  ```kotlin
  val response = stub
      .withDeadlineAfter(5, TimeUnit.SECONDS)
      .methodName(request)
  ```

### 3. 401 Unauthorized
- Проверь токен в localStorage (DevTools → Application → Local Storage)
- Проверь что Auth Service работает: `curl http://89.169.139.21:8081/health`
- Проверь логи Auth: `sudo docker logs sigame-auth-service --tail 50`

### 4. Сервис не запускается
```bash
# Проверь статус
sudo docker ps -a | grep sigame-service-name

# Проверь логи
sudo docker logs sigame-service-name --tail 100

# Проверь health
sudo docker inspect sigame-service-name | grep -A 10 Health
```

### 5. База данных
```bash
# Проверь что БД запущены
sudo docker ps | grep postgres

# Проверь подключение изнутри контейнера
sudo docker exec sigame-lobby-service sh -c 'nc -zv postgres-lobby 5432'
```

## 📝 Важные замечания

1. **Автоматический деплой**: Push в `main` → GitHub Actions → автодеплой на сервер
2. **CI проверки**: Pull Request → запускается CI (линтеры, тесты)
3. **Переменные окружения**: `.env.production` на сервере (не в git!)
4. **Docker сеть**: Сервисы общаются по именам контейнеров (не по IP!)
5. **gRPC порты**:
   - Auth gRPC: 50051
   - Game gRPC: 50053
   - Pack gRPC: 50055

## 🐛 Известные проблемы

1. **Loki/Tempo отключены** - добавлены в profiles, не запускаются по умолчанию
2. **Pack Service использует моки** - gRPC возвращает данные из `mock_data.py`, не из БД
3. **Frontend на порту 80** - без Nginx, прямой доступ

## 🔗 Полезные ссылки

- **GitHub**: https://github.com/Bityta/sigame-project
- **GitHub Actions**: https://github.com/Bityta/sigame-project/actions
- **Yandex Cloud Console**: https://console.cloud.yandex.ru/

## 💡 Быстрые команды

```bash
# Статус всех сервисов
sudo docker ps --format "table {{.Names}}\t{{.Status}}" | grep sigame

# Перезапуск всего стека
cd /opt/sigame && sudo docker compose -f docker-compose.infra.yml -f docker-compose.app.yml --env-file .env.production restart

# Очистка неиспользуемых образов
sudo docker system prune -f

# Проверка места на диске
df -h

# Логи GitHub Actions деплоя
# Смотри в браузере: https://github.com/Bityta/sigame-project/actions
```

## 🔐 Пароли и секреты

Все пароли хранятся в:
- `.env.production` на сервере (не в git!)
- GitHub Secrets (Settings → Secrets and variables → Actions)

Основные секреты:
- `APP_SERVER_IP`: IP сервера
- `SSH_PRIVATE_KEY`: Приватный ключ для SSH
- `JWT_SECRET`: Секрет для JWT токенов

---

**Последнее обновление**: 24.11.2025  
**Версия сервера**: Single server (all-in-one) на Yandex Cloud

