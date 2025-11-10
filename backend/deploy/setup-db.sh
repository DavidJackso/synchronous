#!/bin/bash

set -e

echo "🗄️  Настройка базы данных PostgreSQL..."

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Проверка прав root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Пожалуйста, запустите скрипт с правами root (sudo)${NC}"
    exit 1
fi

echo -e "${GREEN}Установка PostgreSQL...${NC}"
apt-get update
apt-get install -y postgresql postgresql-contrib

echo -e "${GREEN}Запуск PostgreSQL...${NC}"
systemctl start postgresql
systemctl enable postgresql

echo -e "${YELLOW}Настройка PostgreSQL:${NC}"
read -p "Введите имя базы данных (по умолчанию: synchronous): " DB_NAME
DB_NAME=${DB_NAME:-synchronous}

read -p "Введите имя пользователя БД (по умолчанию: synchronous_user): " DB_USER
DB_USER=${DB_USER:-synchronous_user}

read -sp "Введите пароль для пользователя БД: " DB_PASS
echo ""

# Переключение на пользователя postgres
sudo -u postgres psql << EOF
-- Создание пользователя
CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASS}';

-- Создание базы данных
CREATE DATABASE ${DB_NAME} OWNER ${DB_USER} ENCODING 'UTF8';

-- Предоставление привилегий
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};

-- Включение расширений
\c ${DB_NAME}
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
GRANT ALL ON SCHEMA public TO ${DB_USER};
EOF

echo -e "${GREEN}База данных создана!${NC}"
echo ""
echo -e "${YELLOW}DSN для подключения:${NC}"
echo "postgres://${DB_USER}:${DB_PASS}@localhost:5432/${DB_NAME}?sslmode=disable"
echo ""
echo -e "${YELLOW}Для применения миграций выполните:${NC}"
echo "cd /opt/synchronous/backend"
echo "DB_DSN=\"postgres://${DB_USER}:${DB_PASS}@localhost:5432/${DB_NAME}?sslmode=disable\" make migrate-up"

