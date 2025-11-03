-- Тестовые данные для категорий и услуг
-- Запускать после создания всех таблиц через миграции

-- Добавление категорий услуг
INSERT INTO categories (name, description, created_at) VALUES
('Автомойка', 'Полная мойка и очистка автомобиля', NOW()),
('Детейлинг', 'Детальная обработка кузова и салона', NOW()),
('Полировка', 'Полировка кузова и фар', NOW()),
('Химчистка салона', 'Глубокая чистка салона автомобиля', NOW()),
('Защитные покрытия', 'Нанесение защитных покрытий (керамика, пленка)', NOW()),
('ТО', 'Техническое обслуживание', NOW()),
('Ремонт', 'Ремонтные работы', NOW()),
('Шиномонтаж', 'Замена и ремонт шин', NOW()),
('Диагностика', 'Диагностические работы', NOW())
ON CONFLICT DO NOTHING;

-- Получаем ID категорий для услуг
-- Автомойка
DO $$
DECLARE
    category_id_autowash INTEGER;
    category_id_detail INTEGER;
    category_id_polish INTEGER;
    category_id_cleaning INTEGER;
    category_id_protection INTEGER;
    category_id_to INTEGER;
    category_id_repair INTEGER;
    category_id_tire INTEGER;
    category_id_diag INTEGER;
BEGIN
    -- Получаем ID категорий
    SELECT id INTO category_id_autowash FROM categories WHERE name = 'Автомойка' LIMIT 1;
    SELECT id INTO category_id_detail FROM categories WHERE name = 'Детейлинг' LIMIT 1;
    SELECT id INTO category_id_polish FROM categories WHERE name = 'Полировка' LIMIT 1;
    SELECT id INTO category_id_cleaning FROM categories WHERE name = 'Химчистка салона' LIMIT 1;
    SELECT id INTO category_id_protection FROM categories WHERE name = 'Защитные покрытия' LIMIT 1;
    SELECT id INTO category_id_to FROM categories WHERE name = 'ТО' LIMIT 1;
    SELECT id INTO category_id_repair FROM categories WHERE name = 'Ремонт' LIMIT 1;
    SELECT id INTO category_id_tire FROM categories WHERE name = 'Шиномонтаж' LIMIT 1;
    SELECT id INTO category_id_diag FROM categories WHERE name = 'Диагностика' LIMIT 1;

    -- Услуги категории Автомойка
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_autowash, 'Базовая мойка', 'Мойка кузова и колес', 2000, 1500, 3000, 30, NOW()),
    (category_id_autowash, 'Полная мойка', 'Мойка кузова, колес, салона', 3500, 2500, 5000, 60, NOW()),
    (category_id_autowash, 'Премиум мойка', 'Полная мойка + обработка салона', 5000, 4000, 7000, 90, NOW()),
    (category_id_autowash, 'Мойка двигателя', 'Чистка подкапотного пространства', 3000, 2000, 4000, 45, NOW()),
    (category_id_autowash, 'Мойка днища', 'Чистка днища автомобиля', 2500, 2000, 3500, 40, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Детейлинг
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_detail, 'Детейлинг кузова', 'Детальная обработка кузова', 8000, 6000, 12000, 180, NOW()),
    (category_id_detail, 'Детейлинг салона', 'Детальная обработка салона', 6000, 4500, 9000, 150, NOW()),
    (category_id_detail, 'Полный детейлинг', 'Комплексная обработка кузова и салона', 12000, 10000, 18000, 300, NOW()),
    (category_id_detail, 'Чистка двигателя', 'Детальная чистка двигателя', 5000, 4000, 7000, 120, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Полировка
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_polish, 'Полировка кузова', 'Полировка лакокрасочного покрытия', 12000, 10000, 18000, 240, NOW()),
    (category_id_polish, 'Полировка фар', 'Полировка фар и фонарей', 3000, 2000, 5000, 60, NOW()),
    (category_id_polish, 'Восстановительная полировка', 'Глубокая полировка с удалением дефектов', 18000, 15000, 25000, 360, NOW()),
    (category_id_polish, 'Защитная полировка', 'Полировка с нанесением воска', 15000, 12000, 22000, 300, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Химчистка салона
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_cleaning, 'Химчистка салона', 'Полная химчистка салона', 4000, 3000, 6000, 120, NOW()),
    (category_id_cleaning, 'Химчистка сидений', 'Чистка передних и задних сидений', 2500, 2000, 4000, 90, NOW()),
    (category_id_cleaning, 'Химчистка ковриков', 'Глубокая чистка ковриков', 1500, 1000, 2500, 60, NOW()),
    (category_id_cleaning, 'Химчистка багажника', 'Чистка багажного отделения', 2000, 1500, 3000, 45, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Защитные покрытия
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_protection, 'Керамическое покрытие', 'Нанесение керамического покрытия', 25000, 20000, 35000, 480, NOW()),
    (category_id_protection, 'Антигравийная пленка', 'Нанесение антигравийной пленки', 15000, 12000, 25000, 360, NOW()),
    (category_id_protection, 'Жидкое стекло', 'Нанесение защитного покрытия "жидкое стекло"', 8000, 6000, 12000, 180, NOW()),
    (category_id_protection, 'Бронирование фар', 'Защитная пленка на фары', 5000, 4000, 8000, 120, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории ТО
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_to, 'Замена масла', 'Замена моторного масла и фильтра', 2000, 1500, 3000, 60, NOW()),
    (category_id_to, 'Замена фильтров', 'Замена воздушного, салонного фильтров', 1500, 1000, 2500, 45, NOW()),
    (category_id_to, 'Замена свечей', 'Замена свечей зажигания', 3000, 2000, 5000, 90, NOW()),
    (category_id_to, 'Проверка и доливка жидкостей', 'Проверка всех технических жидкостей', 1000, 500, 2000, 30, NOW()),
    (category_id_to, 'Регулировка фар', 'Регулировка света фар', 2000, 1500, 3000, 30, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Ремонт
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_repair, 'Ремонт двигателя', 'Диагностика и ремонт двигателя', 15000, 10000, 30000, 480, NOW()),
    (category_id_repair, 'Ремонт тормозной системы', 'Замена колодок, дисков, тормозной жидкости', 8000, 5000, 15000, 240, NOW()),
    (category_id_repair, 'Ремонт подвески', 'Диагностика и ремонт элементов подвески', 10000, 7000, 20000, 300, NOW()),
    (category_id_repair, 'Ремонт системы охлаждения', 'Ремонт радиатора, помпы, термостата', 6000, 4000, 12000, 180, NOW()),
    (category_id_repair, 'Ремонт электрооборудования', 'Ремонт электрики автомобиля', 5000, 3000, 10000, 150, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Шиномонтаж
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_tire, 'Замена шин', 'Демонтаж и монтаж колес', 5000, 3000, 8000, 120, NOW()),
    (category_id_tire, 'Балансировка колес', 'Балансировка всех колес', 1500, 1000, 2500, 60, NOW()),
    (category_id_tire, 'Ремонт прокола', 'Ремонт прокола шины', 1000, 500, 2000, 30, NOW()),
    (category_id_tire, 'Перебортовка', 'Демонтаж и монтаж шины на диск', 2000, 1500, 3000, 45, NOW()),
    (category_id_tire, 'Хранение шин', 'Сезонное хранение шин', 2000, 1500, 3000, 15, NOW())
    ON CONFLICT DO NOTHING;

    -- Услуги категории Диагностика
    INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at) VALUES
    (category_id_diag, 'Компьютерная диагностика', 'Диагностика всех систем автомобиля', 2000, 1500, 3500, 90, NOW()),
    (category_id_diag, 'Диагностика двигателя', 'Проверка работы двигателя', 2500, 2000, 4000, 60, NOW()),
    (category_id_diag, 'Диагностика подвески', 'Проверка элементов подвески', 1500, 1000, 2500, 45, NOW()),
    (category_id_diag, 'Диагностика АКПП', 'Проверка автоматической коробки передач', 3000, 2500, 5000, 120, NOW()),
    (category_id_diag, 'Диагностика кондиционера', 'Проверка системы кондиционирования', 2000, 1500, 3500, 60, NOW())
    ON CONFLICT DO NOTHING;
END $$;

