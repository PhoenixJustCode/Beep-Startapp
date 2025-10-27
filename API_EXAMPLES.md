# Примеры использования BEEP API

## Обзор

Этот документ содержит примеры использования API для проекта BEEP.

Base URL: `http://localhost:8080`

---

## 1. Калькулятор цен (Спринт 1)

### Получить все категории услуг

```bash
GET /api/v1/categories
```

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "ТО",
    "description": "Техническое обслуживание",
    "created_at": "2024-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "name": "Ремонт",
    "description": "Ремонтные работы",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Получить услуги по категории

```bash
GET /api/v1/services?category_id=1
```

**Ответ:**
```json
[
  {
    "id": 1,
    "category_id": 1,
    "name": "Замена масла",
    "description": "",
    "base_price": 2000,
    "min_price": 1500,
    "max_price": 3000,
    "duration_minutes": 60,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Получить список машин

```bash
GET /api/v1/cars
```

**Ответ:**
```json
[
  {
    "id": 1,
    "brand": "Toyota",
    "model": "Camry",
    "year": 2020,
    "type": "Standard"
  },
  {
    "id": 2,
    "brand": "BMW",
    "model": "X5",
    "year": 2021,
    "type": "Premium"
  }
]
```

### Рассчитать цену услуги

```bash
POST /api/v1/pricing/calculate
Content-Type: application/json

{
  "service_id": 1,
  "car_id": 1
}
```

**Ответ:**
```json
{
  "service_id": 1,
  "service_name": "Замена масла",
  "base_price": 2000,
  "final_price": 2430,
  "min_price": 1500,
  "max_price": 3000
}
```

---

## 2. Онлайн-запись к мастеру (Спринт 2)

### Получить всех мастеров

```bash
GET /api/v1/masters
```

**Ответ:**
```json
[
  {
    "id": 1,
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "phone": "+79001234567",
    "specialization": "ТО, Ремонт",
    "rating": 4.8,
    "photo_url": "",
    "location_lat": 55.7558,
    "location_lng": 37.6173,
    "address": "Москва, ул. Примерная, 1",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Получить детали мастера

```bash
GET /api/v1/masters/1
```

### Получить отзывы мастера

```bash
GET /api/v1/masters/1/reviews
```

### Получить доступные слоты времени

```bash
GET /api/v1/masters/1/available-slots?date=2024-01-15
```

**Ответ:**
```json
{
  "date": "2024-01-15",
  "slots": [
    "09:00",
    "09:30",
    "10:00",
    "10:30",
    "11:00",
    "11:30"
  ]
}
```

### Создать запись к мастеру

```bash
POST /api/v1/appointments
Content-Type: application/json

{
  "master_id": 1,
  "service_id": 1,
  "date": "2024-01-15",
  "time": "10:00",
  "comment": "Пожалуйста, почините вовремя"
}
```

**Ответ:**
```json
{
  "id": 1,
  "user_id": 1,
  "master_id": 1,
  "service_id": 1,
  "date": "2024-01-15T00:00:00Z",
  "time": "10:00",
  "status": "pending",
  "comment": "Пожалуйста, почините вовремя",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Получить мои записи

```bash
GET /api/v1/appointments
```

**Ответ:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "master_id": 1,
    "service_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "time": "10:00",
    "status": "pending",
    "comment": "Пожалуйста, почините вовремя",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Отменить запись

```bash
DELETE /api/v1/appointments/1
```

**Ответ:**
```json
{
  "message": "Appointment cancelled successfully"
}
```

---

## Статусы кодов HTTP

- `200 OK` - успешный запрос
- `201 Created` - ресурс успешно создан
- `400 Bad Request` - невалидные данные
- `404 Not Found` - ресурс не найден
- `500 Internal Server Error` - ошибка сервера

---

## Примеры с cURL

### Расчет цены

```bash
curl -X POST http://localhost:8080/api/v1/pricing/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": 1,
    "car_id": 1
  }'
```

### Создание записи

```bash
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "master_id": 1,
    "service_id": 1,
    "date": "2024-01-15",
    "time": "10:00",
    "comment": "Пожалуйста"
  }'
```

### Получение мастеров

```bash
curl http://localhost:8080/api/v1/masters
```

---

## Примеры с JavaScript (Fetch)

### Расчет цены

```javascript
const response = await fetch('http://localhost:8080/api/v1/pricing/calculate', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    service_id: 1,
    car_id: 1
  })
});

const data = await response.json();
console.log(data);
```

### Создание записи

```javascript
const response = await fetch('http://localhost:8080/api/v1/appointments', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    master_id: 1,
    service_id: 1,
    date: '2024-01-15',
    time: '10:00',
    comment: 'Пожалуйста'
  })
});

const data = await response.json();
console.log(data);
```
