@echo off
echo ========================================
echo    BEEP MVP - Автосервис
echo ========================================
echo.

echo Остановка предыдущих процессов...
taskkill /F /IM go.exe 2>nul

echo.
echo Запуск приложения на порту 8080...
echo.

set DATABASE_URL=postgres://postgres:postgres@localhost:5432/beep?sslmode=disable
set PORT=8080
set JWT_SECRET=your-secret-key

go run cmd/main.go

pause

