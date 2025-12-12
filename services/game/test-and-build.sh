#!/bin/bash
set -e

echo "════════════════════════════════════════════════════════════"
echo "  КОМПИЛЯЦИЯ И ТЕСТИРОВАНИЕ GAME SERVICE"
echo "════════════════════════════════════════════════════════════"
echo ""

cd "$(dirname "$0")"

if ! command -v go &> /dev/null; then
    echo "❌ Go не найден в PATH"
    echo ""
    echo "Установи Go:"
    echo "  brew install go"
    echo "  или скачай с https://golang.org/dl/"
    exit 1
fi

echo "✓ Go найден: $(go version)"
echo ""

echo "1. Обновление зависимостей..."
go mod tidy
echo "   ✓ go mod tidy завершен"
echo ""

echo "2. Компиляция..."
if go build -o /tmp/game-server ./cmd/server 2>&1; then
    echo "   ✅ Компиляция успешна!"
    rm -f /tmp/game-server
    echo ""
else
    echo "   ❌ Ошибки компиляции (см. выше)"
    echo ""
    echo "Основные проблемы которые нужно исправить вручную:"
    echo "- Несовместимые типы между слоями"
    echo "- Отсутствующие методы в интерфейсах"
    echo "- Циклические зависимости"
    exit 1
fi

echo "3. Тесты domain..."
go test ./internal/domain/... -v 2>&1 | tail -20
echo ""

echo "4. Тесты core..."
go test ./internal/core/... -v 2>&1 | tail -20
echo ""

echo "5. Тесты infrastructure..."
go test ./internal/infrastructure/... -v 2>&1 | tail -20
echo ""

echo "6. Тесты adapter..."
go test ./internal/adapter/... -v 2>&1 | tail -20
echo ""

echo "════════════════════════════════════════════════════════════"
echo "  ✅ ВСЕ ПРОВЕРКИ ЗАВЕРШЕНЫ!"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "Для запуска сервера:"
echo "  go run ./cmd/server"
echo ""

