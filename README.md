# BEEP - Сервис записи на автомойку

Приложение для записи на услуги автомойки и детейлинга.

## Структура проекта

```
beep_mvp/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── config/              # Конфигурация
│   ├── handlers/            # HTTP обработчики
│   ├── models/              # Модели данных
│   ├── repository/          # Работа с БД
│   └── router/              # Маршрутизация
├── static/                  # Статические файлы
│   ├── uploads/            # Загруженные файлы
│   ├── script/              # JavaScript
│   └── *.html               # HTML страницы
├── sample_data.sql          # Тестовые данные для БД
├── DEPLOY.md                # Подробная инструкция по деплою
└── go.mod                   # Зависимости
```

## Быстрый старт (локальная разработка)

### 1. Установка PostgreSQL

```bash
# Ubuntu/Debian
sudo apt install postgresql

# Windows
# Скачайте с postgresql.org/download/
```

### 2. Создание базы данных

```bash
# Подключиться к PostgreSQL
sudo -u postgres psql

# Создать базу и пользователя
CREATE DATABASE beep_db;
CREATE USER beep_user WITH PASSWORD 'beep_password';
GRANT ALL PRIVILEGES ON DATABASE beep_db TO beep_user;
\q
```

### 3. Настройка окружения

Создайте файл `.env` в корне проекта:

```env
DATABASE_URL=postgres://beep_user:beep_password@localhost:5432/beep_db?sslmode=disable
JWT_SECRET=your-secret-key-change-this
PORT=8080
```

### 4. Запуск приложения

```bash
# Установить зависимости
go mod download

# Запустить сервер
go run cmd/main.go
```

Приложение будет доступно по адресу: `http://localhost:8080`

### 5. Загрузка тестовых данных (опционально)

```bash
psql -U beep_user -d beep_db -f sample_data.sql
```

**Тестовые пользователи:**
- Email: `ivan@example.com`, Пароль: `password123`
- Email: `maria@example.com`, Пароль: `password123`
- Email: `alex@example.com`, Пароль: `password123`

## Деплой на VDS

Подробная инструкция доступна в файле [DEPLOY.md](DEPLOY.md)

Основные шаги:
1. Установка PostgreSQL
2. Установка Go
3. Клонирование репозитория
4. Настройка .env
5. Сборка приложения
6. Настройка systemd для автозапуска

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход

### Пользователь
- `GET /api/v1/user/profile` - Получить профиль
- `PUT /api/v1/user/profile` - Обновить профиль
- `POST /api/v1/user/photo` - Загрузить фото

### Мастер
- `POST /api/v1/master/profile` - Создать профиль мастера
- `GET /api/v1/master/profile` - Получить профиль мастера
- `PUT /api/v1/master/profile` - Обновить профиль мастера
- `DELETE /api/v1/master/profile` - Удалить профиль мастера
- `POST /api/v1/master/photo` - Загрузить фото мастера
- `GET /api/v1/master/works` - Получить работы мастера
- `POST /api/v1/master/works` - Добавить работу
- `PUT /api/v1/master/works/:id` - Обновить работу
- `DELETE /api/v1/master/works/:id` - Удалить работу

### Записи
- `GET /api/v1/appointments` - Получить записи
- `POST /api/v1/appointments` - Создать запись
- `PUT /api/v1/appointments/:id` - Обновить запись
- `PUT /api/v1/appointments/:id/cancel` - Отменить запись
- `DELETE /api/v1/appointments/:id` - Удалить запись

### Отзывы
- `GET /api/v1/master/reviews` - Получить отзывы мастера
- `POST /api/v1/reviews` - Добавить отзыв

## Технологии

- **Backend:** Go, Gin
- **Database:** PostgreSQL
- **Frontend:** HTML, JavaScript
- **Authentication:** JWT

## Лицензия

MIT
