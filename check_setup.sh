#!/bin/bash

# Скрипт для проверки настроек на новом сервере

echo "=== Проверка настроек BEEP ==="
echo ""

# Проверка .env файла
echo "1. Проверка .env файла..."
if [ -f ".env" ]; then
    echo "✅ .env файл существует"
    cat .env | grep -v "JWT_SECRET" | grep -v "password"
else
    echo "❌ .env файл не найден!"
    echo "Создайте файл .env с содержимым:"
    echo "DATABASE_URL=postgres://beep_user:beep123@localhost:5432/beep?sslmode=disable"
    echo "JWT_SECRET=your-secret-key"
    echo "PORT=8080"
    exit 1
fi

# Проверка подключения к БД
echo ""
echo "2. Проверка подключения к базе данных..."
DB_URL=$(grep DATABASE_URL .env | cut -d '=' -f2-)
if psql "$DB_URL" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Подключение к БД успешно"
else
    echo "❌ Не удалось подключиться к БД!"
    echo "Проверьте DATABASE_URL в .env файле"
    exit 1
fi

# Проверка таблиц
echo ""
echo "3. Проверка необходимых таблиц..."
REQUIRED_TABLES=("users" "masters" "favorite_masters" "reviews" "master_works" "master_certificates")
for table in "${REQUIRED_TABLES[@]}"; do
    if psql "$DB_URL" -c "\dt $table" > /dev/null 2>&1; then
        echo "✅ Таблица $table существует"
    else
        echo "❌ Таблица $table не найдена!"
        echo "Запустите миграции или загрузите database_structure.sql"
    fi
done

# Проверка компиляции
echo ""
echo "4. Проверка компиляции..."
if go build -o beep-server cmd/main.go 2>&1 | tee /tmp/build.log; then
    echo "✅ Приложение успешно скомпилировано"
    rm -f /tmp/build.log
else
    echo "❌ Ошибка компиляции!"
    cat /tmp/build.log
    exit 1
fi

echo ""
echo "=== Проверка завершена ==="
echo ""
echo "Запустите сервер командой: ./beep-server"




