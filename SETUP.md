# Setup Instructions for BEEP - Spirit Planning

## Prerequisites

1. **Go** (уже установлен ✓)
   - Версия: go1.23.4

2. **PostgreSQL** (нужно установить)
   - Опция 1: Установить Docker Desktop для Windows
     - Скачать с: https://www.docker.com/products/docker-desktop
   - Опция 2: Установить PostgreSQL напрямую
     - Скачать с: https://www.postgresql.org/download/windows/

## Quick Start

### Если используете Docker (рекомендуется):

```bash
# В директории beep-backend запустите:
docker compose up -d

# Это запустит PostgreSQL и само приложение
```

### Если используете локальный PostgreSQL:

1. Создайте базу данных:
```sql
CREATE DATABASE beep;
```

2. Запустите приложение:
```bash
cd beep-backend/cmd
go run main.go
```

3. Откройте браузер:
```
http://localhost:8080
```

## Структура

- `static/index.html` - Frontend интерфейс
- `cmd/main.go` - Entry point приложения
- `internal/` - Backend логика
  - `handlers/` - HTTP обработчики
  - `database/` - Работа с БД
  - `models/` - Модели данных
  - `router/` - Роутинг
  - `repository/` - Слой данных

## API Endpoints

- `GET /api/v1/categories` - Список категорий
- `GET /api/v1/services` - Список услуг
- `GET /api/v1/cars` - Список машин
- `GET /api/v1/masters` - Список мастеров
- `POST /api/v1/appointments` - Создать запись
- `GET /api/v1/pricing/calculate` - Рассчитать цену

## Текущий статус

- ✅ Go установлен
- ✅ Зависимости загружены
- ⏳ Ожидается установка PostgreSQL/Docker
- ✅ Frontend готов
- ✅ Backend готов


