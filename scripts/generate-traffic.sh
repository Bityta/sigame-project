#!/bin/bash

# Скрипт для генерации тестового трафика к API

HOST="${1:-89.169.139.21}"
DURATION="${2:-3600}"  # 1 час по умолчанию

echo "🚀 Генерация тестового трафика к SIGame API"
echo "Host: $HOST"
echo "Duration: $DURATION seconds"
echo "=========================================="

# Счётчики
SUCCESS=0
ERRORS=0

# Токены для авторизованных запросов
declare -a TOKENS=()

# Функция для выполнения запроса
make_request() {
    local url=$1
    local method=${2:-GET}
    local data=$3
    local token=$4
    
    if [ -n "$token" ]; then
        if [ "$method" = "POST" ] && [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $token" \
                -d "$data" 2>/dev/null)
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
                -H "Authorization: Bearer $token" 2>/dev/null)
        fi
    else
        if [ "$method" = "POST" ] && [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
                -H "Content-Type: application/json" \
                -d "$data" 2>/dev/null)
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" 2>/dev/null)
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 400 ]; then
        SUCCESS=$((SUCCESS + 1))
        return 0
    else
        ERRORS=$((ERRORS + 1))
        return 1
    fi
}

# Регистрация и логин пользователей
echo "Регистрация тестовых пользователей..."
for i in {1..10}; do
    USERNAME="testuser_$(date +%s)_$i"
    PASSWORD="password123"
    
    # Регистрация
    make_request "http://$HOST:8081/api/auth/register" "POST" "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" > /dev/null 2>&1
    
    # Логин
    response=$(curl -s -X POST "http://$HOST:8081/api/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
    
    token=$(echo "$response" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    if [ -n "$token" ]; then
        TOKENS+=("$token")
    fi
done

echo "Зарегистрировано пользователей: ${#TOKENS[@]}"
echo ""

# Запуск генерации трафика
START_TIME=$(date +%s)
END_TIME=$((START_TIME + DURATION))

echo "Начало генерации трафика..."
echo ""

REQUEST_COUNT=0

while [ $(date +%s) -lt $END_TIME ]; do
    CURRENT_TIME=$(date +%s)
    ELAPSED=$((CURRENT_TIME - START_TIME))
    REQUEST_COUNT=$((REQUEST_COUNT + 1))
    
    # Очистка строки и вывод прогресса
    printf "\r⏱️  Прошло: ${ELAPSED}s | 📊 Запросов: $REQUEST_COUNT | ✅ Успешно: $SUCCESS | ❌ Ошибок: $ERRORS"
    
    # Выбор случайного токена
    if [ ${#TOKENS[@]} -gt 0 ]; then
        TOKEN=${TOKENS[$((RANDOM % ${#TOKENS[@]}))]}
    else
        TOKEN=""
    fi
    
    # === AUTH SERVICE ===
    # Регистрация нового пользователя
    if [ $((RANDOM % 10)) -eq 0 ]; then
        USERNAME="user_$(date +%s)_$RANDOM"
        make_request "http://$HOST:8081/api/auth/register" "POST" "{\"username\":\"$USERNAME\",\"password\":\"pass123\"}" > /dev/null 2>&1
    fi
    
    # Логин
    if [ $((RANDOM % 5)) -eq 0 ]; then
        make_request "http://$HOST:8081/api/auth/login" "POST" "{\"username\":\"testuser_123\",\"password\":\"password123\"}" > /dev/null 2>&1
    fi
    
    # Получение профиля
    if [ -n "$TOKEN" ]; then
        make_request "http://$HOST:8081/api/auth/profile" "GET" "" "$TOKEN" > /dev/null 2>&1
    fi
    
    # Обновление профиля
    if [ -n "$TOKEN" ] && [ $((RANDOM % 20)) -eq 0 ]; then
        make_request "http://$HOST:8081/api/auth/profile" "PUT" "{\"display_name\":\"User $RANDOM\"}" "$TOKEN" > /dev/null 2>&1
    fi
    
    # === LOBBY SERVICE ===
    # Получение списка комнат
    make_request "http://$HOST:8082/api/lobby/rooms" "GET" > /dev/null 2>&1
    
    # Создание комнаты
    if [ $((RANDOM % 8)) -eq 0 ]; then
        ROOM_NAME="Room_$(date +%s)"
        make_request "http://$HOST:8082/api/lobby/rooms" "POST" "{\"name\":\"$ROOM_NAME\",\"isPublic\":true,\"maxPlayers\":6,\"settings\":{\"timeForAnswer\":30,\"timeForChoice\":60,\"allowWrongAnswer\":true,\"showRightAnswer\":true}}" > /dev/null 2>&1
    fi
    
    # Получение информации о конкретной комнате (случайный UUID)
    if [ $((RANDOM % 15)) -eq 0 ]; then
        ROOM_ID="00000000-0000-0000-0000-000000000001"
        make_request "http://$HOST:8082/api/lobby/rooms/$ROOM_ID" "GET" > /dev/null 2>&1
    fi
    
    # Присоединение к комнате
    if [ $((RANDOM % 12)) -eq 0 ]; then
        ROOM_CODE="ABC123"
        make_request "http://$HOST:8082/api/lobby/rooms/join" "POST" "{\"roomCode\":\"$ROOM_CODE\"}" > /dev/null 2>&1
    fi
    
    # === PACK SERVICE ===
    # Получение списка паков
    make_request "http://$HOST:8084/api/packs" "GET" > /dev/null 2>&1
    
    # Получение конкретного пака
    if [ $((RANDOM % 5)) -eq 0 ]; then
        PACK_ID=$((RANDOM % 10 + 1))
        make_request "http://$HOST:8084/api/packs/$PACK_ID" "GET" > /dev/null 2>&1
    fi
    
    # Поиск паков
    if [ $((RANDOM % 10)) -eq 0 ]; then
        QUERY="test"
        make_request "http://$HOST:8084/api/packs/search?q=$QUERY" "GET" > /dev/null 2>&1
    fi
    
    # === GAME SERVICE ===
    # Получение статуса игры (случайный UUID)
    if [ $((RANDOM % 10)) -eq 0 ]; then
        GAME_ID="00000000-0000-0000-0000-000000000001"
        make_request "http://$HOST:8083/api/game/$GAME_ID/status" "GET" > /dev/null 2>&1
    fi
    
    # Случайная задержка между запросами (50-200ms)
    sleep 0.$(printf "%03d" $((RANDOM % 150 + 50)))
done

echo ""
echo ""
echo "=========================================="
echo "✅ Генерация трафика завершена!"
echo "Всего запросов: $((SUCCESS + ERRORS))"
echo "Успешных: $SUCCESS"
echo "Ошибок: $ERRORS"
echo "Успешность: $(awk "BEGIN {printf \"%.2f\", ($SUCCESS / ($SUCCESS + $ERRORS)) * 100}")%"
echo "=========================================="
