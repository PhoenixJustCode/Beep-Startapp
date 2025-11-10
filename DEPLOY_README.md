# Инструкция по деплою BEEP

## Быстрый деплой на VPS

Подробная инструкция находится в файле **VPS_DEPLOY.md**

## Минимальный набор файлов для деплоя

```
beep_mvp/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrations.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── models/
│   │   └── models.go
│   ├── repository/
│   │   └── repository.go
│   └── router/
│       └── router.go
├── static/
│   ├── index.html
│   ├── login.html
│   ├── profile.html
│   ├── reviews.html
│   ├── logo.png
│   ├── script/
│   │   ├── profile.js
│   │   ├── reviews.js
│   │   └── script.js
│   ├── style/
│   │   └── style.css
│   └── uploads/
├── database_structure.sql
├── sample_data.sql          # Тестовые данные (2 мастера, 2 машины)
├── go.mod
├── go.sum
└── .env                     # Создать вручную на сервере
```

## Краткие команды

```bash
# 1. Собрать приложение
go build -o beep-server cmd/main.go

# 2. Запустить
./beep-server

# 3. Или через systemd (см. VPS_DEPLOY.md)
sudo systemctl start beep
```

## Файлы для удаления

Эти файлы были удалены как ненужные:
- `add_masters.go` - перенесено в sample_data.sql
- `test_data.sql` - объединено с sample_data.sql
- `services_autoservice.sql` - объединено с sample_data.sql
- `DEPLOY.md` - заменен на VPS_DEPLOY.md




