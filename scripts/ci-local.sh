#!/bin/bash
set -e

echo "=========================================="
echo "  🔍 LOCAL CI CHECKS"
echo "=========================================="
echo ""

FAILED=0

# ==========================================
# Frontend
# ==========================================
echo "📦 [1/5] Frontend..."
cd frontend
if npm ci &>/dev/null && npm run lint &>/dev/null && npm run build &>/dev/null; then
    echo "  ✅ Frontend OK"
else
    echo "  ❌ Frontend FAILED"
    FAILED=$((FAILED + 1))
fi
cd ..

# ==========================================
# Auth Service
# ==========================================
echo "🔐 [2/5] Auth Service..."
cd services/auth
if go mod download &>/dev/null && go build -v ./cmd/server &>/dev/null; then
    echo "  ✅ Auth Service OK"
else
    echo "  ❌ Auth Service FAILED"
    FAILED=$((FAILED + 1))
fi
cd ../..

# ==========================================
# Lobby Service
# ==========================================
echo "🏛️ [3/5] Lobby Service..."
cd services/lobby
if gradle build -x test &>/dev/null; then
    echo "  ✅ Lobby Service OK"
else
    echo "  ❌ Lobby Service FAILED"
    FAILED=$((FAILED + 1))
fi
cd ../..

# ==========================================
# Game Service
# ==========================================
echo "🎮 [4/5] Game Service..."
cd services/game
if go mod download &>/dev/null && go build -v ./cmd/server &>/dev/null; then
    echo "  ✅ Game Service OK"
else
    echo "  ❌ Game Service FAILED"
    FAILED=$((FAILED + 1))
fi
cd ../..

# ==========================================
# Pack Service
# ==========================================
echo "📦 [5/5] Pack Service..."
cd services/pack
if go mod download &>/dev/null && go build -v ./cmd/server &>/dev/null; then
    echo "  ✅ Pack Service OK"
else
    echo "  ❌ Pack Service FAILED"
    FAILED=$((FAILED + 1))
fi
cd ../..

# ==========================================
# Summary
# ==========================================
echo ""
echo "=========================================="
if [ $FAILED -eq 0 ]; then
    echo "  ✅ ALL CHECKS PASSED!"
    echo "=========================================="
    echo ""
    echo "🚀 Готово к коммиту и push!"
    exit 0
else
    echo "  ❌ $FAILED CHECK(S) FAILED!"
    echo "=========================================="
    echo ""
    echo "⚠️ Пожалуйста, исправьте ошибки перед push"
    exit 1
fi

