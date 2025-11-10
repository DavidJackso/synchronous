#!/bin/bash

# Скрипт для настройки GitHub Actions secrets
# Запустите этот скрипт локально для генерации инструкций

echo "🔐 Настройка GitHub Actions Secrets"
echo ""
echo "Для настройки автовыката добавьте следующие secrets в GitHub:"
echo ""
echo "1. Перейдите в Settings → Secrets and variables → Actions"
echo "2. Добавьте следующие secrets:"
echo ""
echo "   SERVER_HOST - IP адрес или домен сервера"
echo "   SERVER_USER - имя пользователя для SSH (например: root или synchronous)"
echo "   SERVER_SSH_KEY - приватный SSH ключ для подключения к серверу"
echo "   SERVER_PORT - порт SSH (опционально, по умолчанию 22)"
echo "   DB_DSN - DSN для подключения к БД"
echo "            postgres://user:password@host:5432/database?sslmode=disable"
echo ""
echo "3. Для генерации SSH ключа (если нет):"
echo "   ssh-keygen -t ed25519 -C 'github-actions' -f ~/.ssh/github_actions"
echo "   cat ~/.ssh/github_actions.pub | ssh user@server 'cat >> ~/.ssh/authorized_keys'"
echo ""
echo "4. Скопируйте приватный ключ:"
echo "   cat ~/.ssh/github_actions"
echo ""
