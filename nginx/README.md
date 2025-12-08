# Конфигурация Nginx для Synchronous (поддомен tg)

## Установка конфигурации

1. Скопируйте конфигурацию на сервер:
   ```bash
   sudo cp /opt/synchronous/nginx/conf.d/tg.focus-sync.conf /etc/nginx/sites-available/tg.focus-sync.conf
   ```

2. Создайте символическую ссылку:
   ```bash
   sudo ln -s /etc/nginx/sites-available/tg.focus-sync.conf /etc/nginx/sites-enabled/
   ```

3. Удалите дефолтную конфигурацию (если есть):
   ```bash
   sudo rm -f /etc/nginx/sites-enabled/default
   ```

4. Проверьте конфигурацию:
   ```bash
   sudo nginx -t
   ```

5. Перезагрузите nginx:
   ```bash
   sudo systemctl reload nginx
   ```

## Получение SSL сертификата

**ВАЖНО:** Конфигурация настроена для работы на HTTP до получения сертификата.

1. Убедитесь, что nginx работает на HTTP:
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

2. Получите SSL сертификат:
   ```bash
   sudo certbot --nginx -d tg.focus-sync.ru
   ```

3. Certbot автоматически:
   - Получит сертификат
   - Обновит конфигурацию для HTTPS
   - Настроит редирект с HTTP на HTTPS

После получения сертификата приложение будет доступно по HTTPS.

## Структура конфигурации

- **HTTP (порт 80)**: Редирект на HTTPS + поддержка Let's Encrypt
- **HTTPS (порт 443)**: Основной сервер
  - `/api/` → Backend (порт 8080)
  - `/swagger/` → Swagger UI (порт 8081)
  - `/` → Frontend (порт 3000)

## Проверка работы

```bash
# Проверка статуса nginx
sudo systemctl status nginx

# Проверка логов
sudo tail -f /var/log/nginx/tg.focus-sync-access.log
sudo tail -f /var/log/nginx/tg.focus-sync-error.log

# Проверка подключения
curl -I https://tg.focus-sync.ru/api/v1/health
```

## Обновление конфигурации

После изменения конфигурации:
```bash
sudo nginx -t && sudo systemctl reload nginx
```

## Быстрая установка (если контейнеры уже запущены)

Если контейнеры запущены, но nginx возвращает 404:

```bash
# 1. Скопируйте конфигурацию
sudo cp /opt/synchronous/nginx/conf.d/tg.focus-sync.conf /etc/nginx/sites-available/tg.focus-sync.conf

# 2. Активируйте
sudo rm -f /etc/nginx/sites-enabled/tg.focus-sync.conf
sudo ln -s /etc/nginx/sites-available/tg.focus-sync.conf /etc/nginx/sites-enabled/

# 3. Проверьте и перезагрузите
sudo nginx -t && sudo systemctl reload nginx

# 4. Проверьте работу
curl -I https://tg.focus-sync.ru/api/v1/health
```

## Устранение неполадок

### 404 Not Found

Если получаете 404 через nginx, но backend работает локально:

1. **Проверьте текущую конфигурацию nginx:**
   ```bash
   sudo cat /etc/nginx/sites-enabled/tg.focus-sync.conf | grep -A 10 "location /api"
   ```

2. **Если location /api отсутствует, обновите конфигурацию:**
   ```bash
   # Скопируйте правильную конфигурацию
   sudo cp /opt/synchronous/nginx/conf.d/tg.focus-sync.conf /etc/nginx/sites-available/tg.focus-sync.conf
   
   # Пересоздайте ссылку
   sudo rm /etc/nginx/sites-enabled/tg.focus-sync.conf
   sudo ln -s /etc/nginx/sites-available/tg.focus-sync.conf /etc/nginx/sites-enabled/
   
   # Проверьте и перезагрузите
   sudo nginx -t
   sudo systemctl reload nginx
   ```

3. **Проверьте, что backend запущен:**
   ```bash
   docker ps | grep backend
   curl http://localhost:8080/api/v1/health
   ```

4. **Проверьте логи nginx:**
   ```bash
   sudo tail -20 /var/log/nginx/tg.focus-sync-error.log
   ```

### 502 Bad Gateway или Connection refused

Если получаете 502 или `Connection refused` в логах:

1. **Проверьте статус всех контейнеров:**
   ```bash
   docker ps -a | grep synchronous
   ```

2. **Запустите контейнеры, если они не запущены:**
   ```bash
   cd /opt/synchronous
   docker compose up -d backend frontend swagger-ui
   ```

3. **Проверьте, что контейнеры слушают на правильных портах:**
   ```bash
   # Backend должен быть на 8080
   curl http://localhost:8080/api/v1/health
   
   # Frontend должен быть на 3000
   curl http://localhost:3000/health
   
   # Swagger UI должен быть на 8081
   curl http://localhost:8081/
   ```

4. **Проверьте логи контейнеров:**
   ```bash
   docker logs synchronous_backend
   docker logs synchronous_frontend
   ```

5. **Если контейнеры не запускаются, проверьте .env файл:**
   ```bash
   cat /opt/synchronous/.env
   ```

