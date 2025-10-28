# Инструкция по деплою BEEP на VDS

## Быстрый старт

### 1. Подготовка VDS (один раз)

```bash
# Обновить систему
sudo apt update && sudo apt upgrade -y

# Установить PostgreSQL
sudo apt install postgresql postgresql-contrib

# Установить Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Проверить
go version
```

### 2. Настройка PostgreSQL

```bash
sudo -u postgres psql
```

```sql
CREATE DATABASE beep_db;
CREATE USER beep_user WITH PASSWORD 'beep_password';
GRANT ALL PRIVILEGES ON DATABASE beep_db TO beep_user;
ALTER USER beep_user CREATEDB;
\q
```

### 3. Установка приложения

```bash
# Клонировать репозиторий
git clone <your-repo-url>
cd beep_mvp

# Создать .env
cat > .env << EOF
DATABASE_URL=postgres://beep_user:beep_password@localhost:5432/beep_db?sslmode=disable
JWT_SECRET=change-this-to-random-string
PORT=8080
EOF

# Загрузить зависимости
go mod download

# Собрать приложение
go build -o beep-server cmd/main.go

# Запустить (для проверки)
./beep-server
```

Сервер должен запуститься на `http://localhost:8080`

### 4. Настройка автозапуска

```bash
# Создать systemd сервис
sudo nano /etc/systemd/system/beep.service
```

Добавить:
```ini
[Unit]
Description=BEEP Server
After=network.target postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/beep_mvp
ExecStart=/home/ubuntu/beep_mvp/beep-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Запустить сервис
sudo systemctl daemon-reload
sudo systemctl enable beep
sudo systemctl start beep
sudo systemctl status beep
```

### 5. Загрузка тестовых данных (опционально)

```bash
psql -U beep_user -d beep_db -f sample_data.sql
```

Тестовые пользователи:
- Email: `ivan@example.com`, Password: `password123`
- Email: `maria@example.com`, Password: `password123`

## Команды для управления

```bash
# Перезапуск
sudo systemctl restart beep

# Остановка
sudo systemctl stop beep

# Просмотр логов
sudo journalctl -u beep -f

# Проверка статуса
sudo systemctl status beep
```

## Обновление приложения

```bash
# Остановить сервис
sudo systemctl stop beep

# Обновить код
cd beep_mvp
git pull origin main

# Пересобрать
go build -o beep-server cmd/main.go

# Запустить
sudo systemctl start beep
```

## Настройка Nginx (опционально)

### 1. Установка Nginx

```bash
sudo apt install nginx
```

### 2. Конфигурация

```bash
sudo nano /etc/nginx/sites-available/beep
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
# Активировать
sudo ln -s /etc/nginx/sites-available/beep /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. SSL (опционально)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## Резервное копирование

```bash
# Создать бэкап
pg_dump -U beep_user beep_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановить
psql -U beep_user beep_db < backup_file.sql
```

## Безопасность

```bash
# Настроить файрвол
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

## Проблемы

### Порт занят
```bash
sudo netstat -tlnp | grep :8080
sudo kill <PID>
```

### База не подключается
```bash
sudo systemctl status postgresql
sudo systemctl start postgresql
```

### Логи
```bash
sudo journalctl -u beep -n 100
tail -f /var/log/nginx/error.log
```