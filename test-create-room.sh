#!/bin/bash

# Тестовый скрипт для создания комнаты

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Тест создания комнаты ===${NC}"

# 1. Получаем токен (логинимся)
echo -e "\n${YELLOW}1. Логинимся для получения токена...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }')

echo "Login Response: $LOGIN_RESPONSE"

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
  echo -e "${RED}Ошибка: Не удалось получить токен${NC}"
  echo "Попробуем зарегистрировать пользователя..."
  
  REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8001/auth/register \
    -H "Content-Type: application/json" \
    -d '{
      "username": "testuser",
      "password": "testpass123"
    }')
  
  echo "Register Response: $REGISTER_RESPONSE"
  ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.access_token')
fi

if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
  echo -e "${RED}Ошибка: Не удалось получить токен${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Токен получен${NC}"

# 2. Получаем список паков
echo -e "\n${YELLOW}2. Получаем список паков...${NC}"
PACKS_RESPONSE=$(curl -s -X GET http://localhost:8004/api/packs)
echo "Packs Response: $PACKS_RESPONSE"

PACK_ID=$(echo $PACKS_RESPONSE | jq -r '.packs[0].id')

if [ "$PACK_ID" == "null" ] || [ -z "$PACK_ID" ]; then
  echo -e "${RED}Ошибка: Паки не найдены${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Pack ID: $PACK_ID${NC}"

# 3. Создаем комнату
echo -e "\n${YELLOW}3. Создаем комнату...${NC}"

# Правильный JSON запрос
REQUEST_BODY=$(cat <<EOF
{
  "name": "Test Room",
  "packId": "$PACK_ID",
  "maxPlayers": 4,
  "isPublic": true,
  "settings": {
    "timeForAnswer": 30,
    "timeForChoice": 10
  }
}
EOF
)

echo "Request Body:"
echo "$REQUEST_BODY" | jq '.'

ROOM_RESPONSE=$(curl -v -X POST http://localhost:8002/api/lobby/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "$REQUEST_BODY" 2>&1)

echo -e "\n${YELLOW}Response:${NC}"
echo "$ROOM_RESPONSE"

# Проверяем статус код
if echo "$ROOM_RESPONSE" | grep -q "201 Created"; then
  echo -e "\n${GREEN}✓ Комната успешно создана!${NC}"
elif echo "$ROOM_RESPONSE" | grep -q "400 Bad Request"; then
  echo -e "\n${RED}✗ Ошибка 400: Неверный запрос${NC}"
  echo "$ROOM_RESPONSE" | grep -A 50 "{"
elif echo "$ROOM_RESPONSE" | grep -q "401 Unauthorized"; then
  echo -e "\n${RED}✗ Ошибка 401: Не авторизован${NC}"
else
  echo -e "\n${RED}✗ Неизвестная ошибка${NC}"
fi

echo -e "\n${YELLOW}=== Проверьте логи сервисов для деталей ===${NC}"
echo "docker logs sigame-lobby-1"

