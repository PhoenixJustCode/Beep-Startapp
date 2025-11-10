# Установка на новый сервер (DLC/VPS)

## Рекомендуемый способ: Чистая установка через git clone

**Лучше удалить старые файлы и сделать git clone заново** - это гарантирует, что все файлы будут актуальными и не будет конфликтов.

### Шаг 1: Удалить старую директорию (если есть)

```bash
# Если есть старая директория
cd ~
rm -rf beep_mvp
```

### Шаг 2: Клонировать репозиторий

```bash
# Клонировать проект в новую директорию
git clone <ваш-repo-url> beep_mvp
cd beep_mvp

# Переключиться на нужную ветку (если нужно)
git checkout sprint-3
```

### Шаг 3: Установить зависимости и настройки

```bash
# Установить Go зависимости
go mod download

# Создать файл .env
nano .env
```

Вставить в `.env`:
```env
DATABASE_URL=postgres://beep_user:your_password@localhost:5432/beep_db?sslmode=disable
JWT_SECRET=your_very_long_random_secret_key_minimum_32_characters
PORT=8080
```

### Шаг 4: Настроить базу данных

```bash
# Войти в PostgreSQL
sudo -u postgres psql

# Выполнить в PostgreSQL:
CREATE DATABASE beep_db;
CREATE USER beep_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE beep_db TO beep_user;
ALTER USER beep_user CREATEDB;
\q
```

### Шаг 5: Собрать и запустить

```bash
cd ~/beep_mvp

# Собрать приложение
go build -o beep-server cmd/main.go

# Проверить настройки (опционально)
bash check_setup.sh

# Тестовый запуск
./beep-server
```

Если видите "Server starting on port 8080" - всё работает! Остановите Ctrl+C.

### Шаг 6: Загрузить тестовые данные (опционально)

```bash
# Загрузить тестовые данные
psql -U beep_user -d beep_db -f sample_data.sql
```

### Шаг 7: Настроить systemd (для автозапуска)

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
# Активировать сервис
sudo systemctl daemon-reload
sudo systemctl enable beep
sudo systemctl start beep

# Проверить статус
sudo systemctl status beep

# Смотреть логи
sudo journalctl -u beep -f
```

## Альтернативный способ: Обновление существующего проекта

Если хотите обновить существующий проект (не рекомендуется для первого раза):

```bash
cd ~/beep_mvp

# Остановить сервис
sudo systemctl stop beep

# Обновить код
git fetch origin
git checkout sprint-3
git pull origin sprint-3

# Пересобрать
go build -o beep-server cmd/main.go

# Запустить
sudo systemctl start beep
```

## Быстрая проверка

```bash
# Проверить работу
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/categories

# Проверить логи
sudo journalctl -u beep -n 50
```

## Что делать, если что-то не работает?

1. **Проверьте логи:** `sudo journalctl -u beep -n 50`
2. **Проверьте подключение к БД:** `psql -U beep_user -d beep_db`
3. **Проверьте .env файл:** `cat ~/beep_mvp/.env`
4. **Проверьте, что порт свободен:** `sudo netstat -tlnp | grep :8080`

Подробная инструкция находится в файле **VPS_DEPLOY.md**



