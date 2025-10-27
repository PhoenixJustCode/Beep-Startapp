# Быстрый старт - BEEP локально

## Шаг 1: Создайте базу данных

Откройте **pgAdmin** (или другой клиент PostgreSQL) и выполните:

```sql
CREATE DATABASE beep;
```

Или через командную строку:
```bash
# Вариант 1: Через psql с вашим паролем
psql -U postgres -c "CREATE DATABASE beep;"

# Вариант 2: Если у вас другой пользователь
psql -U [ваш_пользователь] -c "CREATE DATABASE beep;"
```

## Шаг 2: Отредактируйте .env файл

Откройте файл `beep-backend/.env` и измените строку DATABASE_URL:

```
DATABASE_URL=postgres://[ваш_пользователь]:[ваш_пароль]@localhost:5432/beep?sslmode=disable
```

Например, если используете стандартного пользователя postgres:
```
DATABASE_URL=postgres://postgres:ваш_пароль@localhost:5432/beep?sslmode=disable
```

## Шаг 3: Запустите приложение

```bash
cd beep-backend
go run cmd/main.go
```

Приложение создаст все таблицы и тестовые данные автоматически!

## Шаг 4: Откройте браузер

Перейдите на: **http://localhost:8080**

Всё должно работать! 🎉





