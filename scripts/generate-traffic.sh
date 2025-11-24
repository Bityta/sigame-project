#!/bin/bash

# Скрипт для генерации тестового трафика к API

HOST="${1:-localhost}"
DURATION="${2:-300}"  # 5 минут по умолчанию

echo "🚀 Генерация тестового трафика к SIGame API"
echo "Host: $HOST"
echo "Duration: $DURATION seconds"
echo "=========================================="

# Счётчики
SUCCESS=0
ERRORS=0

# Функция для выполнения запроса
make_request() {
    local url=$1
    local method=${2:-GET}
    local data=$3
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
            -H "Content-Type: application/json" \
            -d "$data" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" "$url" 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -n1)
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 400 ]; then
        SUCCESS=$((SUCCESS + 1))
        echo "✓"
    else
        ERRORS=$((ERRORS + 1))
        echo "✗ ($http_code)"
    fi
}

# Запуск генерации трафика
START_TIME=$(date +%s)
END_TIME=$((START_TIME + DURATION))

echo "Начало генерации трафика..."
echo ""

while [ $(date +%s) -lt $END_TIME ]; do
    CURRENT_TIME=$(date +%s)
    ELAPSED=$((CURRENT_TIME - START_TIME))
    
    # Очистка строки и вывод прогресса
    printf "\r⏱️  Прошло: ${ELAPSED}s | ✅ Успешно: $SUCCESS | ❌ Ошибок: $ERRORS"
    
    # Auth Service endpoints
    make_request "http://$HOST:8081/health" > /dev/null 2>&1
    make_request "http://$HOST:8081/metrics" > /dev/null 2>&1
    
    # Lobby Service endpoints  
    make_request "http://$HOST:8082/api/lobby/health" > /dev/null 2>&1
    make_request "http://$HOST:8082/api/lobby/rooms" > /dev/null 2>&1
    make_request "http://$HOST:8082/actuator/health" > /dev/null 2>&1
    make_request "http://$HOST:8082/actuator/prometheus" > /dev/null 2>&1
    
    # Game Service endpoints
    make_request "http://$HOST:8083/health" > /dev/null 2>&1
    
    # Pack Service endpoints
    make_request "http://$HOST:8084/health" > /dev/null 2>&1
    make_request "http://$HOST:8084/api/packs" > /dev/null 2>&1
    make_request "http://$HOST:8084/metrics" > /dev/null 2>&1
    
    # Случайная задержка между запросами (100-500ms)
    sleep 0.$((RANDOM % 5))
done

echo ""
echo ""
echo "=========================================="
echo "✅ Генерация трафика завершена!"
echo "Всего запросов: $((SUCCESS + ERRORS))"
echo "Успешных: $SUCCESS"
echo "Ошибок: $ERRORS"
echo "=========================================="

