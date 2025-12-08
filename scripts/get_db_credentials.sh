#!/bin/bash

# Скрипт для получения учетных данных базы данных на сервере

echo "═══════════════════════════════════════════════════════════"
echo "  Получение данных от базы данных PostgreSQL"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Проверка наличия docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не найден. Установите Docker для использования этого скрипта."
    exit 1
fi

echo "1️⃣  Данные из PostgreSQL контейнера:"
echo "───────────────────────────────────────────────────────────"
if docker ps | grep -q synchronous_postgres; then
    echo "✓ Контейнер PostgreSQL запущен"
    echo ""
    echo "Пользователь:"
    docker exec synchronous_postgres printenv POSTGRES_USER 2>/dev/null || echo "  (не задан или используется значение по умолчанию)"
    echo ""
    echo "Пароль:"
    PASSWORD=$(docker exec synchronous_postgres printenv POSTGRES_PASSWORD 2>/dev/null)
    if [ -n "$PASSWORD" ]; then
        echo "  $PASSWORD"
    else
        echo "  (не задан, проверьте .env файл или docker-compose.yml)"
    fi
    echo ""
    echo "База данных:"
    docker exec synchronous_postgres printenv POSTGRES_DB 2>/dev/null || echo "  (не задан или используется значение по умолчанию)"
else
    echo "❌ Контейнер PostgreSQL не запущен"
fi

echo ""
echo "2️⃣  Данные из Backend контейнера:"
echo "───────────────────────────────────────────────────────────"
if docker ps | grep -q synchronous_backend; then
    echo "✓ Контейнер Backend запущен"
    echo ""
    echo "DB_DSN (строка подключения):"
    docker exec synchronous_backend printenv DB_DSN 2>/dev/null || echo "  (не задан)"
    echo ""
    echo "Отдельные переменные:"
    docker exec synchronous_backend printenv | grep -E "^DB_" | sed 's/^/  /' || echo "  (не заданы)"
else
    echo "❌ Контейнер Backend не запущен"
fi

echo ""
echo "3️⃣  Данные из файла .env:"
echo "───────────────────────────────────────────────────────────"
if [ -f .env ]; then
    echo "✓ Файл .env найден"
    echo ""
    echo "DB_USER:"
    grep "^DB_USER=" .env | cut -d'=' -f2 | sed 's/^/  /' || echo "  (не задан)"
    echo ""
    echo "DB_PASSWORD:"
    grep "^DB_PASSWORD=" .env | cut -d'=' -f2 | sed 's/^/  /' || echo "  (не задан)"
    echo ""
    echo "DB_NAME:"
    grep "^DB_NAME=" .env | cut -d'=' -f2 | sed 's/^/  /' || echo "  (не задан)"
else
    echo "⚠️  Файл .env не найден"
fi

echo ""
echo "4️⃣  Значения по умолчанию из docker-compose.yml:"
echo "───────────────────────────────────────────────────────────"
if [ -f docker-compose.yml ]; then
    echo "✓ Файл docker-compose.yml найден"
    echo ""
    DEFAULT_USER=$(grep "POSTGRES_USER:" docker-compose.yml | grep -o '\${DB_USER:-[^}]*}' | cut -d: -f2 | cut -d} -f1 || echo "synchronous_user")
    DEFAULT_PASS=$(grep "POSTGRES_PASSWORD:" docker-compose.yml | grep -o '\${DB_PASSWORD:-[^}]*}' | cut -d: -f2 | cut -d} -f1 || echo "change_me")
    DEFAULT_DB=$(grep "POSTGRES_DB:" docker-compose.yml | grep -o '\${DB_NAME:-[^}]*}' | cut -d: -f2 | cut -d} -f1 || echo "synchronous_db")
    
    echo "Пользователь (по умолчанию): $DEFAULT_USER"
    echo "Пароль (по умолчанию): $DEFAULT_PASS"
    echo "База данных (по умолчанию): $DEFAULT_DB"
else
    echo "⚠️  Файл docker-compose.yml не найден"
fi

echo ""
echo "5️⃣  Формат строки подключения (DSN):"
echo "───────────────────────────────────────────────────────────"
echo "  postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable"
echo ""
echo "  Пример (внутри Docker сети):"
echo "  postgres://synchronous_user:change_me@postgres:5432/synchronous_db?sslmode=disable"
echo ""
echo "  Пример (с сервера/хоста):"
echo "  postgres://synchronous_user:change_me@localhost:5432/synchronous_db?sslmode=disable"

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Готово!"
echo "═══════════════════════════════════════════════════════════"

