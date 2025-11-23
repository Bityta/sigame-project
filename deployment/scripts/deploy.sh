#!/bin/bash
set -e

echo "=========================================="
echo "  🚀 Быстрый деплой на сервер"
echo "=========================================="
echo ""

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

APP_SERVER="ubuntu@89.169.139.21"
PROJECT_DIR="/opt/sigame"

echo -e "${YELLOW}1. Коммит и пуш изменений...${NC}"
git add -A
if git diff --staged --quiet; then
    echo "Нет изменений для коммита"
else
    read -p "Введите сообщение коммита: " commit_msg
    git commit -m "$commit_msg"
fi

git push origin feature/deployment-setup

echo ""
echo -e "${YELLOW}2. Деплой на сервер...${NC}"
ssh $APP_SERVER << 'ENDSSH'
set -e

cd /opt/sigame

echo "📥 Получение последних изменений..."
git pull origin feature/deployment-setup

echo ""
echo "🛑 Остановка контейнеров..."
sudo docker compose -f docker-compose.app.yml --env-file .env.production down

echo ""
echo "🔨 Пересборка изменённых сервисов..."
sudo docker compose -f docker-compose.app.yml --env-file .env.production build

echo ""
echo "🚀 Запуск сервисов..."
sudo docker compose -f docker-compose.app.yml --env-file .env.production up -d

echo ""
echo "⏳ Ожидание запуска (10 секунд)..."
sleep 10

echo ""
echo "📊 Статус контейнеров:"
sudo docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep sigame || true

echo ""
echo "✅ Деплой завершён!"
ENDSSH

echo ""
echo -e "${GREEN}=========================================="
echo "  ✅ ДЕПЛОЙ ЗАВЕРШЁН УСПЕШНО"
echo "==========================================${NC}"
echo ""
echo "🌐 Проверьте: http://89.169.139.21"
echo ""
