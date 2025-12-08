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

# 3. Проверьте синтаксис
sudo nginx -t

# 4. Перезагрузите nginx (reload может не подхватить изменения, используйте restart)
sudo systemctl restart nginx

# 5. Проверьте, что nginx запущен
sudo systemctl status nginx

# 6. Проверьте работу
curl -v https://tg.focus-sync.ru/api/v1/health

# 7. Если все еще 404, проверьте логи
sudo tail -50 /var/log/nginx/tg.focus-sync-error.log
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

2. **Убедитесь, что все контейнеры запущены и здоровы:**
   ```bash
   # Проверьте статус
   docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep synchronous
   
   # Если контейнеры не запущены или unhealthy, перезапустите их
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

4. **Если frontend не отвечает, проверьте логи:**
   ```bash
   docker logs synchronous_frontend --tail=50
   ```

5. **Проверьте, что порты проброшены правильно:**
   ```bash
   # Должно показать, что порт 3000 проброшен на 0.0.0.0:3000->80/tcp
   docker ps | grep frontend
   
   # Проверьте, что порт слушается на хосте
   netstat -tlnp | grep 3000
   # или
   ss -tlnp | grep 3000
   ```

6. **Если проблема сохраняется, перезапустите frontend:**
   ```bash
   cd /opt/synchronous
   docker compose restart frontend
   
   # Подождите несколько секунд и проверьте снова
   sleep 5
   curl http://localhost:3000/health
   ```

7. **Проверьте логи nginx после перезапуска:**
   ```bash
   sudo tail -f /var/log/nginx/tg.focus-sync-error.log
   ```

5. **Если контейнеры не запускаются, проверьте .env файл:**
   ```bash
   cat /opt/synchronous/.env
   ```

