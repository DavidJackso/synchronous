#!/bin/bash

set -e

# Скрипт для удаленного развертывания
# Используется GitHub Actions или для ручного деплоя

echo "🚀 Удаленное развертывание приложения..."

APP_NAME="synchronous"
APP_DIR="/opt/${APP_NAME}/backend"
SERVICE_NAME="synchronous"

# Проверка наличия бинарника
if [ ! -f "${APP_DIR}/app" ]; then
    echo "❌ Ошибка: бинарник не найден в ${APP_DIR}/app"
    exit 1
fi

# Установка прав
chmod +x "${APP_DIR}/app"

# Применение миграций (если указан DB_DSN)
if [ -n "$DB_DSN" ]; then
    echo "📦 Применение миграций..."
    cd "${APP_DIR}"
    export DB_DSN
    make migrate-up || echo "⚠️  Миграции не применены (возможно, уже применены)"
else
    echo "⚠️  DB_DSN не указан, миграции пропущены"
fi

# Перезапуск сервиса
echo "🔄 Перезапуск сервиса..."
sudo systemctl restart "${SERVICE_NAME}"

# Проверка статуса
sleep 2
if systemctl is-active --quiet "${SERVICE_NAME}"; then
    echo "✅ Сервис успешно запущен"
    systemctl status "${SERVICE_NAME}" --no-pager -l
else
    echo "❌ Ошибка: сервис не запущен"
    systemctl status "${SERVICE_NAME}" --no-pager -l
    journalctl -u "${SERVICE_NAME}" -n 50 --no-pager
    exit 1
fi

# Health check
echo "🏥 Проверка здоровья приложения..."
sleep 3
if curl -f -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "✅ Приложение отвечает на health check"
else
    echo "⚠️  Health check не прошел, но сервис запущен"
fi

echo "✅ Развертывание завершено!"

