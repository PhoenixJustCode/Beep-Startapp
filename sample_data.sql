-- Упрощенные тестовые данные для базы данных BEEP
-- 2 пользователя, 2 мастера, 2 машины, базовые категории и услуги
-- Пароль для всех пользователей: password123

-- Тестовые пользователи
INSERT INTO users (name, email, phone, password_hash, created_at, updated_at) VALUES
('Иван Петров', 'ivan@example.com', '+7-777-123-4567', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
('Мария Сидорова', 'maria@example.com', '+7-777-234-5678', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- Тестовые категории услуг
INSERT INTO categories (name, description, created_at) VALUES
('Автомойка', 'Полная мойка автомобиля', NOW()),
('Детейлинг', 'Детальная обработка кузова', NOW()),
('Полировка', 'Полировка кузова и фар', NOW()),
('Химчистка салона', 'Чистка салона автомобиля', NOW()),
('Защитные покрытия', 'Нанесение защитных покрытий', NOW())
ON CONFLICT DO NOTHING;

-- Тестовые услуги
DO $$
DECLARE
    cat_autowash INTEGER;
    cat_detail INTEGER;
    cat_polish INTEGER;
    cat_cleaning INTEGER;
    cat_protection INTEGER;
BEGIN
    -- Получить ID категорий
    SELECT id INTO cat_autowash FROM categories WHERE name = 'Автомойка' LIMIT 1;
    SELECT id INTO cat_detail FROM categories WHERE name = 'Детейлинг' LIMIT 1;
    SELECT id INTO cat_polish FROM categories WHERE name = 'Полировка' LIMIT 1;
    SELECT id INTO cat_cleaning FROM categories WHERE name = 'Химчистка салона' LIMIT 1;
    SELECT id INTO cat_protection FROM categories WHERE name = 'Защитные покрытия' LIMIT 1;

    -- Услуги
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (cat_autowash, 'Базовая мойка', 'Мойка кузова и колес', 2000, 1500, 3000, 30, NOW()),
    (cat_autowash, 'Полная мойка', 'Мойка кузова, колес, салона', 3500, 2500, 5000, 60, NOW()),
    (cat_detail, 'Детейлинг кузова', 'Детальная обработка кузова', 8000, 6000, 12000, 180, NOW()),
    (cat_polish, 'Полировка кузова', 'Полировка лакокрасочного покрытия', 12000, 10000, 18000, 240, NOW()),
    (cat_cleaning, 'Химчистка салона', 'Полная химчистка салона', 4000, 3000, 6000, 120, NOW()),
    (cat_protection, 'Керамическое покрытие', 'Нанесение керамического покрытия', 25000, 20000, 35000, 480, NOW())
    ON CONFLICT DO NOTHING;
END $$;

-- Тестовые автомобили (2 машины)
INSERT INTO cars (brand, model, year, type, created_at) VALUES
('Toyota', 'Camry', 2020, 'Sedan', NOW()),
('BMW', 'X5', 2021, 'SUV', NOW())
ON CONFLICT DO NOTHING;

-- Тестовые мастера (2 мастера)
-- Получаем user_id для мастеров
DO $$
DECLARE
    user_ivan_id INTEGER;
    user_maria_id INTEGER;
    master_ivan_id INTEGER;
    master_maria_id INTEGER;
BEGIN
    -- Получить ID пользователей
    SELECT id INTO user_ivan_id FROM users WHERE email = 'ivan@example.com';
    SELECT id INTO user_maria_id FROM users WHERE email = 'maria@example.com';

    -- Создать профили мастеров
    INSERT INTO masters (user_id, name, email, phone, specialization, rating, address, created_at, updated_at) VALUES
    (user_ivan_id, 'Иван Петров', 'ivan@example.com', '+7-777-123-4567', 'Автомойка и полировка', 4.8, 'ул. Абая, 150, Алматы', NOW(), NOW()),
    (user_maria_id, 'Мария Сидорова', 'maria@example.com', '+7-777-234-5678', 'Детейлинг и защитные покрытия', 4.9, 'ул. Сатпаева, 90, Алматы', NOW(), NOW())
    ON CONFLICT (user_id) DO NOTHING
    RETURNING id INTO master_ivan_id;

    -- Получить ID мастеров
    SELECT id INTO master_ivan_id FROM masters WHERE user_id = user_ivan_id;
    SELECT id INTO master_maria_id FROM masters WHERE user_id = user_maria_id;

    -- Расписание для первого мастера (Иван)
    INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at) VALUES
    (master_ivan_id, 0, '10:00', '16:00', true, NOW()),  -- Воскресенье
    (master_ivan_id, 1, '09:00', '18:00', true, NOW()),  -- Понедельник
    (master_ivan_id, 2, '09:00', '18:00', true, NOW()),  -- Вторник
    (master_ivan_id, 3, '09:00', '18:00', true, NOW()),  -- Среда
    (master_ivan_id, 4, '09:00', '18:00', true, NOW()),  -- Четверг
    (master_ivan_id, 5, '09:00', '18:00', true, NOW()),  -- Пятница
    (master_ivan_id, 6, '10:00', '16:00', true, NOW())   -- Суббота
    ON CONFLICT DO NOTHING;

    -- Расписание для второго мастера (Мария)
    INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at) VALUES
    (master_maria_id, 0, '11:00', '17:00', true, NOW()),  -- Воскресенье
    (master_maria_id, 1, '10:00', '19:00', true, NOW()),  -- Понедельник
    (master_maria_id, 2, '10:00', '19:00', true, NOW()),  -- Вторник
    (master_maria_id, 3, '10:00', '19:00', true, NOW()),  -- Среда
    (master_maria_id, 4, '10:00', '19:00', true, NOW()),  -- Четверг
    (master_maria_id, 5, '10:00', '19:00', true, NOW()),  -- Пятница
    (master_maria_id, 6, '11:00', '18:00', true, NOW())   -- Суббота
    ON CONFLICT DO NOTHING;

    -- Платежная информация для мастеров
    INSERT INTO master_payment_info (master_id, kaspi_card, freedom_card, halyk_card, created_at, updated_at) VALUES
    (master_ivan_id, 'KZ123456789012345678', NULL, NULL, NOW(), NOW()),
    (master_maria_id, NULL, 'KZ987654321098765432', NULL, NOW(), NOW())
    ON CONFLICT (master_id) DO NOTHING;

    -- Примеры работ мастеров
    INSERT INTO master_works (master_id, title, work_date, customer_name, amount, created_at) VALUES
    (master_ivan_id, 'Полная мойка BMW X5', CURRENT_DATE - INTERVAL '5 days', 'Анна Иванова', 3500, NOW()),
    (master_ivan_id, 'Полировка Toyota Camry', CURRENT_DATE - INTERVAL '3 days', 'Петр Сидоров', 12000, NOW()),
    (master_maria_id, 'Детейлинг кузова', CURRENT_DATE - INTERVAL '7 days', 'Елена Козлова', 8000, NOW()),
    (master_maria_id, 'Керамическое покрытие', CURRENT_DATE - INTERVAL '2 days', 'Михаил Волков', 25000, NOW())
    ON CONFLICT DO NOTHING;

    -- Примеры отзывов
    INSERT INTO reviews (master_id, user_id, rating, comment, created_at) VALUES
    (master_ivan_id, user_ivan_id, 5, 'Отличная работа! Машина как новая.', NOW()),
    (master_ivan_id, user_maria_id, 5, 'Профессиональный подход, рекомендую!', NOW()),
    (master_maria_id, user_ivan_id, 5, 'Очень доволен результатом.', NOW()),
    (master_maria_id, user_maria_id, 5, 'Прекрасный мастер, качественная работа!', NOW())
    ON CONFLICT DO NOTHING;
END $$;
