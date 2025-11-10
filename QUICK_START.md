# Быстрый запуск на новом сервере

## 1. Настройка .env файла

Создайте файл `.env` в корне проекта:

```env
DATABASE_URL=postgres://beep_user:beep123@localhost:5432/beep?sslmode=disable
JWT_SECRET=your-secret-key
PORT=8080
```

## 2. Создание базы данных

```bash
# Войти в PostgreSQL
sudo -u postgres psql

# Создать базу данных и пользователя
CREATE DATABASE beep;
CREATE USER beep_user WITH PASSWORD 'beep123';
GRANT ALL PRIVILEGES ON DATABASE beep TO beep_user;
ALTER USER beep_user CREATEDB;
\q
```

## 3. Инициализация структуры БД

База данных создастся автоматически при первом запуске приложения (миграции запускаются автоматически).

Если нужно вручную загрузить структуру:

```bash
psql -U beep_user -d beep -f database_structure.sql
```

## 4. Загрузка тестовых данных (опционально)

```bash
psql -U beep_user -d beep -f sample_data.sql
```

## 5. Сборка и запуск

```bash
# Установить зависимости
go mod download

# Собрать приложение
go build -o beep-server cmd/main.go

# Запустить
./beep-server
```

## 6. Проверка работы

- Откройте: `http://localhost:8080`
- API health check: `http://localhost:8080/health`
- API masters: `http://localhost:8080/api/v1/masters`

## Решение проблем

### Если фильтры не работают:

1. **Проверьте, что таблица favorite_masters создана:**
   ```sql
   \dt favorite_masters
   ```
   Если нет, запустите:
   ```bash
   psql -U beep_user -d beep -f database_structure.sql
   ```

2. **Проверьте консоль браузера (F12)** - должны быть запросы к `/api/v1/masters`

3. **Проверьте ответ API:**
   ```bash
   curl http://localhost:8080/api/v1/masters
   ```
   Должны быть поля `is_verified` и `is_favorite`

### Если кнопки не работают:

1. Убедитесь, что сервер запущен
2. Откройте консоль браузера (F12) и проверьте ошибки
3. Проверьте, что API_URL правильно настроен (должен быть `window.location.origin + '/api/v1'`)

## Тестовые пользователи (после загрузки sample_data.sql)

- Email: `ivan@example.com`, Пароль: `password123` (Мастер: Иван Петров)
- Email: `maria@example.com`, Пароль: `password123` (Мастер: Мария Сидорова)




