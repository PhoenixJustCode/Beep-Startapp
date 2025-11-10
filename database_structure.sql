з

-- Таблица автомобилей
CREATE TABLE IF NOT EXISTS cars (
    id SERIAL PRIMARY KEY,
    brand VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year INTEGER,
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица зон ценообразования
CREATE TABLE IF NOT EXISTS price_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    multiplier DECIMAL(5,2) DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица правил ценообразования
CREATE TABLE IF NOT EXISTS pricing_rules (
    id SERIAL PRIMARY KEY,
    service_id INTEGER REFERENCES services(id),
    car_type VARCHAR(50),
    car_age_min INTEGER,
    car_age_max INTEGER,
    multiplier DECIMAL(5,2) DEFAULT 1.0,
    fixed_addition DECIMAL(10,2) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица мастеров
CREATE TABLE IF NOT EXISTS masters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    specialization VARCHAR(100),
    rating DECIMAL(3,2) DEFAULT 0.0,
    photo_url TEXT,
    location_lat DECIMAL(10,8),
    location_lng DECIMAL(11,8),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица расписания мастеров
CREATE TABLE IF NOT EXISTS master_schedule (
    id SERIAL PRIMARY KEY,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица записей
CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    master_id INTEGER REFERENCES masters(id),
    service_id INTEGER REFERENCES services(id),
    date DATE NOT NULL,
    time TIME NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица отзывов
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица работ мастеров
CREATE TABLE IF NOT EXISTS master_works (
    id SERIAL PRIMARY KEY,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    work_date DATE NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    amount DECIMAL(10, 2),
    photo_urls TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица платежной информации мастеров
CREATE TABLE IF NOT EXISTS master_payment_info (
    id SERIAL PRIMARY KEY,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE UNIQUE,
    bank_name VARCHAR(50),
    account_number VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица сертификатов мастеров
CREATE TABLE IF NOT EXISTS master_certificates (
    id SERIAL PRIMARY KEY,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    photo_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица избранных мастеров
CREATE TABLE IF NOT EXISTS favorite_masters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, master_id)
);