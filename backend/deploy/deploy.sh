#!/bin/bash

set -e

echo "🚀 Начало развертывания приложения..."

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Проверка прав root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Пожалуйста, запустите скрипт с правами root (sudo)${NC}"
    exit 1
fi

# Переменные
APP_NAME="synchronous"
APP_USER="synchronous"
APP_DIR="/opt/${APP_NAME}"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"
NGINX_CONFIG="/etc/nginx/sites-available/${APP_NAME}"

echo -e "${GREEN}Шаг 1: Обновление системы...${NC}"
apt-get update
apt-get upgrade -y

echo -e "${GREEN}Шаг 2: Установка зависимостей...${NC}"
apt-get install -y \
    curl \
    wget \
    git \
    build-essential \
    postgresql-client \
    nginx \
    certbot \
    python3-certbot-nginx

echo -e "${GREEN}Шаг 3: Установка Go...${NC}"
if ! command -v go &> /dev/null; then
    GO_VERSION="1.24.4"
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
else
    echo -e "${YELLOW}Go уже установлен${NC}"
fi

echo -e "${GREEN}Шаг 4: Создание пользователя приложения...${NC}"
if ! id "$APP_USER" &>/dev/null; then
    useradd -r -s /bin/bash -d "$APP_DIR" "$APP_USER"
    echo -e "${GREEN}Пользователь $APP_USER создан${NC}"
else
    echo -e "${YELLOW}Пользователь $APP_USER уже существует${NC}"
fi

echo -e "${GREEN}Шаг 5: Создание директории приложения...${NC}"
mkdir -p "$APP_DIR"
chown -R "$APP_USER:$APP_USER" "$APP_DIR"

echo -e "${GREEN}Шаг 6: Клонирование/обновление репозитория...${NC}"
if [ -d "$APP_DIR/.git" ]; then
    cd "$APP_DIR"
    sudo -u "$APP_USER" git pull
else
    echo -e "${YELLOW}Репозиторий не найден. Пожалуйста, склонируйте его вручную:${NC}"
    echo "sudo -u $APP_USER git clone <your-repo-url> $APP_DIR"
fi

echo -e "${GREEN}Шаг 7: Сборка приложения...${NC}"
cd "$APP_DIR/backend"
sudo -u "$APP_USER" /usr/local/go/bin/go mod download
sudo -u "$APP_USER" /usr/local/go/bin/go build -o "$APP_DIR/app" ./cmd/app

echo -e "${GREEN}Шаг 8: Создание systemd service...${NC}"
cat > "$SERVICE_FILE" << EOF
[Unit]
Description=Synchronous Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=$APP_USER
WorkingDirectory=$APP_DIR/backend
ExecStart=$APP_DIR/app
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$APP_NAME

# Environment variables
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable "$APP_NAME"

echo -e "${GREEN}Шаг 9: Настройка Nginx...${NC}"
cat > "$NGINX_CONFIG" << 'EOF'
server {
    listen 80;
    server_name _;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    location /swagger/ {
        proxy_pass http://127.0.0.1:8080/swagger/;
        proxy_set_header Host $host;
    }
}
EOF

ln -sf "$NGINX_CONFIG" /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx

echo -e "${GREEN}✅ Развертывание завершено!${NC}"
echo ""
echo -e "${YELLOW}Следующие шаги:${NC}"
echo "1. Настройте configs/config.toml в $APP_DIR/backend/configs/"
echo "2. Настройте базу данных MySQL"
echo "3. Примените миграции: cd $APP_DIR/backend && DB_DSN='...' make migrate-up"
echo "4. Запустите сервис: systemctl start $APP_NAME"
echo "5. Проверьте статус: systemctl status $APP_NAME"
echo "6. Настройте SSL: certbot --nginx -d your-domain.com"

