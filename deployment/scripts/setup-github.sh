#!/bin/bash

# GitHub Setup Script
# Помощь в создании и настройке GitHub репозитория

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   GitHub Repository Setup${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Шаг 1: Создание репозитория на GitHub
echo -e "${YELLOW}Шаг 1: Создайте репозиторий на GitHub${NC}"
echo ""
echo "1. Откройте: https://github.com/new"
echo "2. Введите название: sigame-project"
echo "3. Описание: Multiplayer quiz game platform"
echo "4. Выберите: Public или Private"
echo "5. НЕ добавляйте README, .gitignore, license (у нас уже есть)"
echo "6. Нажмите 'Create repository'"
echo ""
read -p "Нажмите Enter когда создадите репозиторий..."

# Шаг 2: Получение URL репозитория
echo ""
echo -e "${YELLOW}Шаг 2: Введите URL вашего репозитория${NC}"
echo "Например: https://github.com/username/sigame-project.git"
echo ""
read -p "GitHub URL: " REPO_URL

# Шаг 3: Добавление remote и push
echo ""
echo -e "${YELLOW}Шаг 3: Подключение к GitHub и загрузка кода...${NC}"
git remote add origin "$REPO_URL"
git branch -M main
git push -u origin main

echo -e "${GREEN}✓ Код загружен на GitHub!${NC}"
echo ""

# Шаг 4: Настройка секретов
echo -e "${YELLOW}Шаг 4: Настройка GitHub Secrets для CI/CD${NC}"
echo ""
echo "Откройте: $REPO_URL/settings/secrets/actions"
echo ""
echo "Добавьте 3 секрета:"
echo ""
echo "1. SSH_PRIVATE_KEY"
echo "   Выполните команду:"
echo -e "${GREEN}   cat ~/.ssh/id_rsa${NC}"
echo "   Скопируйте весь вывод (от -----BEGIN до -----END)"
echo ""
echo "2. APP_SERVER_IP"
echo -e "${GREEN}   89.169.139.21${NC}"
echo ""
echo "3. INFRA_SERVER_IP"
echo -e "${GREEN}   10.129.0.26${NC}"
echo ""
read -p "Нажмите Enter когда добавите все секреты..."

# Шаг 5: Проверка
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ GitHub настроен!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Ваш репозиторий: $REPO_URL"
echo ""
echo "Теперь при каждом git push автоматически:"
echo "  ✓ Код деплоится на серверы"
echo "  ✓ Docker образы пересобираются"
echo "  ✓ Сервисы перезапускаются"
echo "  ✓ Проверяется здоровье сервисов"
echo ""
echo "Попробуйте:"
echo "  1. Сделайте изменение в коде"
echo "  2. git add ."
echo "  3. git commit -m 'Test deployment'"
echo "  4. git push"
echo ""
echo "Проверьте статус деплоя:"
echo "  $REPO_URL/actions"
echo ""

