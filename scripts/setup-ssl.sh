#!/bin/bash

set -e

echo "🔐 Настройка SSL сертификата..."

# Цвета
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Проверка наличия certbot
if ! command -v certbot &> /dev/null; then
    echo -e "${YELLOW}Certbot не установлен. Установка...${NC}"
    apt-get update
    apt-get install -y certbot python3-certbot-nginx
fi

# Запрос домена
read -p "Введите ваш домен (например: api.example.com): " DOMAIN

if [ -z "$DOMAIN" ]; then
    echo -e "${RED}Домен не указан. Выход.${NC}"
    exit 1
fi

# Очистка домена от http://, https://, слешей и пробелов
DOMAIN=$(echo "$DOMAIN" | sed 's|^[[:space:]]*https\?://||' | sed 's|/.*$||' | sed 's|[[:space:]]*$||' | sed 's|^[[:space:]]*||')

if [ -z "$DOMAIN" ]; then
    echo -e "${RED}Домен не может быть пустым после очистки. Выход.${NC}"
    exit 1
fi

echo -e "${GREEN}Очищенный домен: $DOMAIN${NC}"

echo -e "${GREEN}Настройка SSL для домена: $DOMAIN${NC}"

# Обновление конфигурации Nginx для домена с SSL
cat > /opt/synchronous/nginx/conf.d/synchronous.conf << EOF
# Редирект с HTTP на HTTPS
server {
    listen 80;
    server_name $DOMAIN;
    return 301 https://\$server_name\$request_uri;
}

# HTTPS сервер
server {
    listen 443 ssl http2;
    server_name $DOMAIN;

    # SSL сертификаты
    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
    
    # SSL настройки
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    client_max_body_size 10M;

    location / {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }

    location /swagger/ {
        proxy_pass http://backend:8080/swagger/;
        proxy_set_header Host \$host;
    }
}
EOF

# Получение сертификата через standalone (для Docker)
echo -e "${GREEN}Получение SSL сертификата...${NC}"
echo -e "${YELLOW}Временно останавливаем Nginx контейнер...${NC}"
docker compose -f /opt/synchronous/docker-compose.yml stop nginx

# Получение сертификата в standalone режиме
certbot certonly --standalone -d $DOMAIN --non-interactive --agree-tos --email admin@$DOMAIN

# Запуск Nginx обратно
echo -e "${GREEN}Запуск Nginx...${NC}"
docker compose -f /opt/synchronous/docker-compose.yml start nginx

echo -e "${GREEN}✅ SSL настроен!${NC}"
echo ""
echo "Теперь ваш сайт доступен по https://$DOMAIN"

