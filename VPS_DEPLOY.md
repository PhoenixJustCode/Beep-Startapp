# Полная инструкция по переносу BEEP на VPS

## Важно! Проверьте все настройки перед запуском

### 1. Подготовка сервера (Ubuntu/Debian)

```bash
# Обновить систему
sudo apt update && sudo apt upgrade -y

# Установить PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Установить Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin
go version

# Установить Nginx (для proxy)
sudo apt install nginx -y
```

### 2. Настройка PostgreSQL

```bash
# Войти в PostgreSQL
sudo -u postgres psql
```

Выполнить в PostgreSQL:
```sql
-- Создать базу данных
CREATE DATABASE beep_db;

-- Создать пользователя
CREATE USER beep_user WITH PASSWORD 'your_secure_password_here';

-- Дать права
GRANT ALL PRIVILEGES ON DATABASE beep_db TO beep_user;
ALTER USER beep_user CREATEDB;

-- Выйти
\q
```

### 3. Настройка приложения на VPS

```bash
# Создать директорию для приложения
mkdir -p ~/beep_mvp
cd ~/beep_mvp

# Загрузить файлы проекта (через git или scp)
# Если через git:
git clone <your-repo-url> .
# Если через scp с локальной машины:
# scp -r ./beep_mvp/* user@vps_ip:~/beep_mvp/
```

### 4. Создание файла .env

```bash
cd ~/beep_mvp
nano .env
```

Вставить следующее (ЗАМЕНИТЬ значения на свои!):
```env
DATABASE_URL=postgres://beep_user:your_secure_password_here@localhost:5432/beep_db?sslmode=disable
JWT_SECRET=change_this_to_random_very_long_string_minimum_32_characters
PORT=8080
```

**ВАЖНО:**
- `your_secure_password_here` - пароль, который вы указали при создании пользователя PostgreSQL
- `JWT_SECRET` - случайная строка минимум 32 символа (для безопасности)

### 5. Сборка и тестирование

```bash
cd ~/beep_mvp

# Загрузить зависимости
go mod download

# Собрать приложение
go build -o beep-server cmd/main.go

# Проверить, что файл создан
ls -lh beep-server

# Тестовый запуск (для проверки)
./beep-server
```

Если видите "Server starting on port 8080" - всё хорошо. Остановите Ctrl+C.

### 6. Инициализация базы данных

При первом запуске приложение само создаст все таблицы. Но можно загрузить тестовые данные:

```bash
# Загрузить структуру (если нужно)
psql -U beep_user -d beep_db -f database_structure.sql

# Загрузить тестовые данные (опционально)
psql -U beep_user -d beep_db -f sample_data.sql
```

### 7. Настройка systemd сервиса

```bash
sudo nano /etc/systemd/system/beep.service
```

Вставить:
```ini
[Unit]
Description=BEEP Server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/root/beep_mvp
EnvironmentFile=/root/beep_mvp/.env
ExecStart=/root/beep_mvp/beep-server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**ВАЖНО:** Замените `/root/beep_mvp` на реальный путь к вашему приложению!

```bash
# Перезагрузить systemd
sudo systemctl daemon-reload

# Включить автозапуск
sudo systemctl enable beep

# Запустить
sudo systemctl start beep

# Проверить статус
sudo systemctl status beep

# Смотреть логи
sudo journalctl -u beep -f
```

### 8. Настройка Nginx (обязательно для работы фильтров!)

```bash
sudo nano /etc/nginx/sites-available/beep
```

Вставить:
```nginx
server {
    listen 80;
    server_name your_domain.com;  # Или IP адрес: 192.168.1.100

    # Максимальный размер загружаемых файлов
    client_max_body_size 10M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Важно для работы API и фильтров
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;
    }

    # Для статических файлов (опционально, можно оставить прокси)
    location /static/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        expires 1d;
        add_header Cache-Control "public, immutable";
    }
}
```

```bash
# Активировать конфигурацию
sudo ln -s /etc/nginx/sites-available/beep /etc/nginx/sites-enabled/

# Удалить дефолтную конфигурацию (если нужно)
sudo rm /etc/nginx/sites-enabled/default

# Проверить конфигурацию
sudo nginx -t

# Перезапустить Nginx
sudo systemctl restart nginx

# Проверить статус
sudo systemctl status nginx
```

### 9. Настройка файрвола

```bash
# Разрешить SSH
sudo ufw allow 22/tcp

# Разрешить HTTP
sudo ufw allow 80/tcp

# Разрешить HTTPS (если будет SSL)
sudo ufw allow 443/tcp

# Включить файрвол
sudo ufw enable

# Проверить статус
sudo ufw status
```

### 10. Проверка работы

1. **Проверьте, что сервер запущен:**
   ```bash
   sudo systemctl status beep
   curl http://localhost:8080/health
   ```

2. **Проверьте доступность через Nginx:**
   ```bash
   curl http://your-domain-or-ip/health
   ```

3. **Откройте в браузере:**
   - `http://your-domain-or-ip/` - главная страница
   - `http://your-domain-or-ip/login` - страница входа
   - `http://your-domain-or-ip/profile` - профиль

4. **Проверьте работу API:**
   ```bash
   curl http://your-domain-or-ip/api/v1/categories
   ```

### 11. Решение проблем

#### Проблема: Фильтры не работают

**Решение:**
1. Проверьте, что Nginx правильно настроен (см. пункт 8)
2. Проверьте логи: `sudo journalctl -u beep -n 50`
3. Проверьте консоль браузера (F12) на наличие ошибок
4. Убедитесь, что `API_URL` в JavaScript использует `window.location.origin` (уже настроено)

#### Проблема: 404 на странице отзывов

**Решение:**
1. Проверьте, что файл `static/reviews.html` существует
2. Проверьте путь в router.go - должен быть `/reviews/:id`
3. Проверьте логи Nginx: `sudo tail -f /var/log/nginx/error.log`

#### Проблема: База данных не подключается

**Решение:**
```bash
# Проверить статус PostgreSQL
sudo systemctl status postgresql

# Проверить подключение
psql -U beep_user -d beep_db -h localhost

# Проверить DATABASE_URL в .env
cat ~/beep_mvp/.env
```

#### Проблема: Порт занят

**Решение:**
```bash
# Найти процесс
sudo netstat -tlnp | grep :8080
# Или
sudo lsof -i :8080

# Убить процесс (замените PID на реальный)
sudo kill -9 <PID>
```

#### Проблема: Статические файлы не загружаются

**Решение:**
1. Проверьте права доступа:
   ```bash
   sudo chmod -R 755 ~/beep_mvp/static
   sudo chown -R root:root ~/beep_mvp/static
   ```

2. Проверьте путь в router.go - должен быть `./static`

#### Проблема: CORS ошибки

**Решение:**
CORS уже настроен в router.go. Если ошибки продолжаются:
1. Проверьте, что Nginx передает правильные заголовки (см. пункт 8)
2. Проверьте, что сервер запущен на правильном порту

### 12. Обновление приложения

```bash
cd ~/beep_mvp

# Остановить сервис
sudo systemctl stop beep

# Обновить код (если через git)
git pull origin main

# Пересобрать
go build -o beep-server cmd/main.go

# Запустить
sudo systemctl start beep

# Проверить
sudo systemctl status beep
```

### 13. Резервное копирование

```bash
# Создать бэкап базы данных
pg_dump -U beep_user beep_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановить из бэкапа
psql -U beep_user beep_db < backup_file.sql

# Бэкап всего проекта
tar -czf beep_backup_$(date +%Y%m%d_%H%M%S).tar.gz ~/beep_mvp
```

### 14. Мониторинг

```bash
# Логи приложения
sudo journalctl -u beep -f

# Логи Nginx
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# Использование ресурсов
htop
# Или
top
```

### 15. Важные моменты для работы фильтров

Фильтры работают через JavaScript, который делает запросы к API. Убедитесь:

1. ✅ Nginx правильно проксирует запросы к `/api/v1/*`
2. ✅ CORS настроен (уже есть в router.go)
3. ✅ JavaScript файлы загружаются (проверьте консоль браузера F12)
4. ✅ API отвечает (проверьте `/api/v1/categories`)

Если фильтры все равно не работают:

1. Откройте консоль браузера (F12)
2. Перейдите на вкладку "Network"
3. Попробуйте использовать фильтр
4. Посмотрите, какие запросы отправляются и какие ошибки есть
5. Проверьте ответы от API - они должны быть JSON

### 16. Быстрая проверка всех компонентов

```bash
# 1. PostgreSQL работает
sudo systemctl status postgresql

# 2. Приложение работает
sudo systemctl status beep

# 3. Nginx работает
sudo systemctl status nginx

# 4. Порт 8080 слушается
sudo netstat -tlnp | grep :8080

# 5. API отвечает
curl http://localhost:8080/api/v1/categories

# 6. Health check
curl http://localhost:8080/health

# 7. Через Nginx
curl http://your-domain-or-ip/api/v1/categories
```

Если все команды возвращают успешные результаты - всё настроено правильно!

---

## Чек-лист перед запуском

- [ ] PostgreSQL установлен и запущен
- [ ] База данных `beep_db` создана
- [ ] Пользователь `beep_user` создан с правильными правами
- [ ] Файл `.env` создан с правильными настройками
- [ ] Приложение собрано (`go build`)
- [ ] Systemd сервис создан и запущен
- [ ] Nginx настроен и работает
- [ ] Файрвол настроен (порты 22, 80, 443 открыты)
- [ ] Health check проходит: `curl http://localhost:8080/health`
- [ ] API отвечает: `curl http://your-domain-or-ip/api/v1/categories`
- [ ] Сайт открывается в браузере
- [ ] Фильтры работают (проверить в браузере)

---

**После выполнения всех пунктов ваше приложение должно работать так же, как на локалке!**


