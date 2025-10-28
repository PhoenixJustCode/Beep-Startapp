-- Тестовые данные для базы данных BEEP
-- Запускать после создания всех таблиц

-- Тестовые пользователи
INSERT INTO users (name, email, phone, password_hash, photo_url, created_at, updated_at) VALUES
('Иван Петров', 'ivan@example.com', '+7-777-123-4567', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NULL, NOW(), NOW()),
('Мария Сидорова', 'maria@example.com', '+7-777-234-5678', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NULL, NOW(), NOW()),
('Алексей Козлов', 'alex@example.com', '+7-777-345-6789', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NULL, NOW(), NOW()),
('Елена Волкова', 'elena@example.com', '+7-777-456-7890', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NULL, NOW(), NOW()),
('Дмитрий Морозов', 'dmitry@example.com', '+7-777-567-8901', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NULL, NOW(), NOW());

-- Тестовые категории услуг
INSERT INTO categories (name, description, created_at, updated_at) VALUES
('Автомойка', 'Полная мойка автомобиля', NOW(), NOW()),
('Детейлинг', 'Детальная обработка кузова', NOW(), NOW()),
('Полировка', 'Полировка кузова и фар', NOW(), NOW()),
('Химчистка салона', 'Чистка салона автомобиля', NOW(), NOW()),
('Защитные покрытия', 'Нанесение защитных покрытий', NOW(), NOW());

-- Тестовые услуги
INSERT INTO services (category_id, name, description, base_price, created_at, updated_at) VALUES
(1, 'Базовая мойка', 'Мойка кузова и колес', 2000, NOW(), NOW()),
(1, 'Полная мойка', 'Мойка кузова, колес, салона', 3500, NOW(), NOW()),
(1, 'Премиум мойка', 'Полная мойка + обработка', 5000, NOW(), NOW()),
(2, 'Детейлинг кузова', 'Детальная обработка кузова', 8000, NOW(), NOW()),
(2, 'Детейлинг салона', 'Детальная обработка салона', 6000, NOW(), NOW()),
(3, 'Полировка кузова', 'Полировка лакокрасочного покрытия', 12000, NOW(), NOW()),
(3, 'Полировка фар', 'Полировка фар и фонарей', 3000, NOW(), NOW()),
(4, 'Химчистка салона', 'Полная химчистка салона', 4000, NOW(), NOW()),
(5, 'Керамическое покрытие', 'Нанесение керамического покрытия', 25000, NOW(), NOW()),
(5, 'Антигравийная пленка', 'Нанесение антигравийной пленки', 15000, NOW(), NOW());

-- Тестовые автомобили
INSERT INTO cars (brand, model, year, color, license_plate, created_at, updated_at) VALUES
('Toyota', 'Camry', 2020, 'Белый', '123ABC01', NOW(), NOW()),
('BMW', 'X5', 2021, 'Черный', '456DEF02', NOW(), NOW()),
('Mercedes', 'E-Class', 2019, 'Серебристый', '789GHI03', NOW(), NOW()),
('Audi', 'A4', 2022, 'Синий', '012JKL04', NOW(), NOW()),
('Lexus', 'RX', 2020, 'Красный', '345MNO05', NOW(), NOW()),
('Volkswagen', 'Tiguan', 2021, 'Серый', '678PQR06', NOW(), NOW()),
('Hyundai', 'Tucson', 2019, 'Белый', '901STU07', NOW(), NOW()),
('Kia', 'Sportage', 2022, 'Черный', '234VWX08', NOW(), NOW());

-- Тестовые мастера
INSERT INTO masters (user_id, name, email, phone, specialization, address, photo_url, rating, created_at, updated_at) VALUES
(1, 'Иван Петров', 'ivan@example.com', '+7-777-123-4567', 'Автомойка и детейлинг', 'ул. Примерная, 1', NULL, 4.8, NOW(), NOW()),
(2, 'Мария Сидорова', 'maria@example.com', '+7-777-234-5678', 'Полировка и защитные покрытия', 'ул. Тестовая, 2', NULL, 4.9, NOW(), NOW()),
(3, 'Алексей Козлов', 'alex@example.com', '+7-777-345-6789', 'Химчистка салона', 'ул. Образцовая, 3', NULL, 4.7, NOW(), NOW());

-- Тестовые записи на услуги
INSERT INTO appointments (user_id, master_id, service_id, car_id, appointment_date, appointment_time, status, comment, total_price, created_at, updated_at) VALUES
(4, 1, 1, 1, '2025-10-29', '10:00', 'confirmed', 'Первая запись', 2000, NOW(), NOW()),
(5, 2, 6, 2, '2025-10-30', '14:00', 'confirmed', 'Полировка кузова', 12000, NOW(), NOW()),
(4, 3, 8, 3, '2025-10-31', '11:00', 'pending', 'Химчистка салона', 4000, NOW(), NOW()),
(5, 1, 3, 4, '2025-11-01', '16:00', 'confirmed', 'Премиум мойка', 5000, NOW(), NOW());

-- Тестовые отзывы
INSERT INTO reviews (user_id, master_id, rating, comment, created_at, updated_at) VALUES
(4, 1, 5, 'Отличная работа! Машина как новая.', NOW(), NOW()),
(5, 2, 5, 'Профессиональный подход, рекомендую!', NOW(), NOW()),
(4, 3, 4, 'Хорошее качество, быстрая работа.', NOW(), NOW()),
(5, 1, 5, 'Очень доволен результатом.', NOW(), NOW());

-- Тестовые работы мастеров
INSERT INTO master_works (master_id, title, work_date, customer_name, amount, photo_urls, created_at, updated_at) VALUES
(1, 'Полная мойка BMW X5', '2025-10-25', 'Анна Иванова', 3500, NULL, NOW(), NOW()),
(1, 'Детейлинг Mercedes E-Class', '2025-10-24', 'Петр Сидоров', 8000, NULL, NOW(), NOW()),
(2, 'Полировка Toyota Camry', '2025-10-23', 'Елена Козлова', 12000, NULL, NOW(), NOW()),
(2, 'Керамическое покрытие Audi A4', '2025-10-22', 'Михаил Волков', 25000, NULL, NOW(), NOW()),
(3, 'Химчистка Lexus RX', '2025-10-21', 'Ольга Морозова', 4000, NULL, NOW(), NOW()),
(3, 'Детейлинг салона Volkswagen Tiguan', '2025-10-20', 'Сергей Новиков', 6000, NULL, NOW(), NOW());

-- Тестовые зоны ценообразования
INSERT INTO price_zones (name, description, created_at, updated_at) VALUES
('Центр города', 'Центральная часть города', NOW(), NOW()),
('Спальные районы', 'Жилые районы', NOW(), NOW()),
('Промышленная зона', 'Промышленные районы', NOW(), NOW());

-- Тестовые правила ценообразования
INSERT INTO pricing_rules (service_id, price_zone_id, multiplier, created_at, updated_at) VALUES
(1, 1, 1.2, NOW(), NOW()),  -- Базовая мойка в центре +20%
(1, 2, 1.0, NOW(), NOW()),  -- Базовая мойка в спальных районах
(1, 3, 0.9, NOW(), NOW()),  -- Базовая мойка в промзоне -10%
(2, 1, 1.1, NOW(), NOW()),  -- Полная мойка в центре +10%
(2, 2, 1.0, NOW(), NOW()),  -- Полная мойка в спальных районах
(2, 3, 0.95, NOW(), NOW()), -- Полная мойка в промзоне -5%
(6, 1, 1.3, NOW(), NOW()),  -- Полировка в центре +30%
(6, 2, 1.1, NOW(), NOW()),  -- Полировка в спальных районах +10%
(6, 3, 1.0, NOW(), NOW());  -- Полировка в промзоне

-- Тестовое расписание мастеров
INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_available, created_at, updated_at) VALUES
(1, 1, '09:00', '18:00', true, NOW(), NOW()),  -- Понедельник
(1, 2, '09:00', '18:00', true, NOW(), NOW()),  -- Вторник
(1, 3, '09:00', '18:00', true, NOW(), NOW()),  -- Среда
(1, 4, '09:00', '18:00', true, NOW(), NOW()),  -- Четверг
(1, 5, '09:00', '18:00', true, NOW(), NOW()),  -- Пятница
(1, 6, '10:00', '16:00', true, NOW(), NOW()),  -- Суббота
(1, 0, '10:00', '14:00', true, NOW(), NOW()),  -- Воскресенье

(2, 1, '10:00', '19:00', true, NOW(), NOW()),  -- Понедельник
(2, 2, '10:00', '19:00', true, NOW(), NOW()),  -- Вторник
(2, 3, '10:00', '19:00', true, NOW(), NOW()),  -- Среда
(2, 4, '10:00', '19:00', true, NOW(), NOW()),  -- Четверг
(2, 5, '10:00', '19:00', true, NOW(), NOW()),  -- Пятница
(2, 6, '11:00', '17:00', true, NOW(), NOW()),  -- Суббота
(2, 0, '11:00', '15:00', true, NOW(), NOW()),  -- Воскресенье

(3, 1, '08:00', '17:00', true, NOW(), NOW()),  -- Понедельник
(3, 2, '08:00', '17:00', true, NOW(), NOW()),  -- Вторник
(3, 3, '08:00', '17:00', true, NOW(), NOW()),  -- Среда
(3, 4, '08:00', '17:00', true, NOW(), NOW()),  -- Четверг
(3, 5, '08:00', '17:00', true, NOW(), NOW()),  -- Пятница
(3, 6, '09:00', '15:00', true, NOW(), NOW()),  -- Суббота
(3, 0, '09:00', '13:00', true, NOW(), NOW());  -- Воскресенье

-- Информация о платежах мастеров
INSERT INTO master_payment_info (master_id, bank_name, account_number, created_at, updated_at) VALUES
(1, 'Kaspi Bank', 'KZ123456789012345678', NOW(), NOW()),
(2, 'Halyk Bank', 'KZ987654321098765432', NOW(), NOW()),
(3, 'Jusan Bank', 'KZ555555555555555555', NOW(), NOW());

-- Пароли для тестовых пользователей (все: password123)
-- Хеш: $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
