-- Добавляем мастеров
INSERT INTO masters (name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at) VALUES
('Иван Петров', 'ivan@beep.com', '+7 (900) 123-4567', 'Механик', 4.8, '', 55.7558, 37.6173, 'Москва, ул. Тверская, 10', NOW(), NOW()),
('Алексей Смирнов', 'alex@beep.com', '+7 (900) 234-5678', 'Электрик', 4.9, '', 55.7558, 37.6173, 'Москва, ул. Арбат, 25', NOW(), NOW()),
('Дмитрий Кузнецов', 'dmitry@beep.com', '+7 (900) 345-6789', 'Диагност', 4.7, '', 55.7558, 37.6173, 'Москва, ул. Ленина, 5', NOW(), NOW()),
('Сергей Волков', 'sergey@beep.com', '+7 (900) 456-7890', 'Маляр', 4.6, '', 55.7558, 37.6173, 'Москва, пр. Мира, 15', NOW(), NOW());

-- Добавляем расписание для каждого мастера (Пн-Пт с 8:00 до 18:00)
-- day_of_week: 0=Monday, 1=Tuesday, 2=Wednesday, 3=Thursday, 4=Friday, 5=Saturday, 6=Sunday
INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at) 
SELECT m.id, d.day, '08:00', '18:00', true, NOW()
FROM masters m
CROSS JOIN (
    SELECT 0 AS day UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
) d;

-- Добавляем субботу (9:00-15:00) для первых двух мастеров
INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at)
SELECT id, 5, '09:00', '15:00', true, NOW()
FROM masters
WHERE id <= 2;

