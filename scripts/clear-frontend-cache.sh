#!/bin/bash

# Скрипт для очистки кеша frontend на сервере

set -e

echo "🧹 Очистка кеша frontend..."

# Остановка и удаление контейнера
echo "1. Остановка frontend контейнера..."
docker stop synchronous_frontend 2>/dev/null || echo "   Контейнер не запущен"
docker rm synchronous_frontend 2>/dev/null || echo "   Контейнер не существует"

# Удаление старого образа
echo "2. Удаление старого Docker образа..."
docker rmi synchronous_frontend:latest 2>/dev/null || echo "   Образ не найден"

# Очистка dist директории (если есть)
echo "3. Очистка dist директории..."
cd /opt/synchronous/frontend
if [ -d "dist" ]; then
    rm -rf dist/*
    echo "   ✅ dist очищена"
else
    echo "   ⚠️  dist не найдена"
fi

# Пересборка контейнера
echo "4. Пересборка frontend контейнера..."
if [ -f "Dockerfile" ]; then
    docker build --no-cache -t synchronous_frontend:latest .
    echo "   ✅ Контейнер пересобран"
else
    echo "   ❌ Dockerfile не найден"
    exit 1
fi

# Запуск нового контейнера
echo "5. Запуск нового контейнера..."
docker run -d \
    --name synchronous_frontend \
    --network synchronous_network \
    -p 3000:80 \
    --restart unless-stopped \
    synchronous_frontend:latest

echo ""
echo "✅ Кеш очищен, контейнер перезапущен!"
echo ""
echo "Проверка статуса:"
docker ps | grep synchronous_frontend || echo "⚠️  Контейнер не запущен"

