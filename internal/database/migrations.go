package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Create tables
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_users_table",
			sql: `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				phone VARCHAR(20),
				photo_url VARCHAR(255),
				password_hash VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_categories_table",
			sql: `
			CREATE TABLE IF NOT EXISTS categories (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				description TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_services_table",
			sql: `
			CREATE TABLE IF NOT EXISTS services (
				id SERIAL PRIMARY KEY,
				category_id INTEGER REFERENCES categories(id),
				name VARCHAR(200) NOT NULL,
				description TEXT,
				base_price DECIMAL(10,2),
				min_price DECIMAL(10,2),
				max_price DECIMAL(10,2),
				duration_minutes INTEGER,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_cars_table",
			sql: `
			CREATE TABLE IF NOT EXISTS cars (
				id SERIAL PRIMARY KEY,
				brand VARCHAR(50) NOT NULL,
				model VARCHAR(50) NOT NULL,
				year INTEGER,
				type VARCHAR(50),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_price_zones_table",
			sql: `
			CREATE TABLE IF NOT EXISTS price_zones (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				multiplier DECIMAL(5,2) DEFAULT 1.0,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_pricing_rules_table",
			sql: `
			CREATE TABLE IF NOT EXISTS pricing_rules (
				id SERIAL PRIMARY KEY,
				service_id INTEGER REFERENCES services(id),
				car_type VARCHAR(50),
				car_age_min INTEGER,
				car_age_max INTEGER,
				multiplier DECIMAL(5,2) DEFAULT 1.0,
				fixed_addition DECIMAL(10,2) DEFAULT 0.0,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_masters_table",
			sql: `
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
			);`,
		},
		{
			name: "create_master_schedule_table",
			sql: `
			CREATE TABLE IF NOT EXISTS master_schedule (
				id SERIAL PRIMARY KEY,
				master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
				day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
				start_time TIME NOT NULL,
				end_time TIME NOT NULL,
				is_active BOOLEAN DEFAULT true,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_appointments_table",
			sql: `
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
			);`,
		},
		{
			name: "create_reviews_table",
			sql: `
			CREATE TABLE IF NOT EXISTS reviews (
				id SERIAL PRIMARY KEY,
				master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
				user_id INTEGER REFERENCES users(id),
				rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
				comment TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_master_works_table",
			sql: `
			CREATE TABLE IF NOT EXISTS master_works (
				id SERIAL PRIMARY KEY,
				master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
				title VARCHAR(255) NOT NULL,
				work_date DATE NOT NULL,
				customer_name VARCHAR(255) NOT NULL,
				amount DECIMAL(10, 2),
				photo_urls TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_master_payment_info_table",
			sql: `
			CREATE TABLE IF NOT EXISTS master_payment_info (
				id SERIAL PRIMARY KEY,
				master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE UNIQUE,
				kaspi_card VARCHAR(20),
				freedom_card VARCHAR(20),
				halyk_card VARCHAR(20),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
	}

	for _, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration.sql); err != nil {
			log.Printf("Error running migration %s: %v", migration.name, err)
			return fmt.Errorf("migration failed: %s", migration.name)
		}
		log.Printf("Migration %s completed", migration.name)
	}

	// Add photo_url column if it doesn't exist
	addPhotoUrlMigration := struct {
		name string
		sql  string
	}{
		name: "add_photo_url_column",
		sql: `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'users' AND column_name = 'photo_url'
			) THEN
				ALTER TABLE users ADD COLUMN photo_url VARCHAR(255);
			END IF;
		END $$;`,
	}

	if _, err := db.ExecContext(ctx, addPhotoUrlMigration.sql); err != nil {
		log.Printf("Warning: Failed to add photo_url column: %v", err)
	} else {
		log.Printf("Migration %s completed", addPhotoUrlMigration.name)
	}

	// Add user_id column to masters table if it doesn't exist
	addUserIDMigration := struct {
		name string
		sql  string
	}{
		name: "add_user_id_to_masters",
		sql: `
			DO $$ 
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'masters' AND column_name = 'user_id'
				) THEN
					ALTER TABLE masters ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
				END IF;
			END $$;`,
	}

	if _, err := db.ExecContext(ctx, addUserIDMigration.sql); err != nil {
		log.Printf("Warning: Failed to add user_id column to masters: %v", err)
	} else {
		log.Printf("Migration %s completed", addUserIDMigration.name)
	}

	// Fix photo_urls column type in master_works table
	fixPhotoURLsMigration := struct {
		name string
		sql  string
	}{
		name: "fix_photo_urls_column_type",
		sql: `
			DO $$ 
			BEGIN
				IF EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'master_works' AND column_name = 'photo_urls' AND data_type = 'ARRAY'
				) THEN
					ALTER TABLE master_works ALTER COLUMN photo_urls TYPE TEXT;
				END IF;
			END $$;`,
	}

	if _, err := db.ExecContext(ctx, fixPhotoURLsMigration.sql); err != nil {
		log.Printf("Warning: Failed to fix photo_urls column type: %v", err)
	} else {
		log.Printf("Migration %s completed", fixPhotoURLsMigration.name)
	}

	// Add trial date columns to user_subscriptions if they don't exist
	addTrialDatesMigration := struct {
		name string
		sql  string
	}{
		name: "add_trial_dates_to_subscriptions",
		sql: `
			DO $$ 
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'user_subscriptions' AND column_name = 'trial_start_date'
				) THEN
					ALTER TABLE user_subscriptions ADD COLUMN trial_start_date TIMESTAMP;
				END IF;
				
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'user_subscriptions' AND column_name = 'trial_end_date'
				) THEN
					ALTER TABLE user_subscriptions ADD COLUMN trial_end_date TIMESTAMP;
				END IF;
			END $$;`,
	}

	if _, err := db.ExecContext(ctx, addTrialDatesMigration.sql); err != nil {
		log.Printf("Warning: Failed to add trial dates columns: %v", err)
	} else {
		log.Printf("Migration %s completed", addTrialDatesMigration.name)
	}

	// Create new tables for additional features
	newMigrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_user_subscriptions_table",
			sql: `
			CREATE TABLE IF NOT EXISTS user_subscriptions (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
				plan VARCHAR(20) DEFAULT 'basic' CHECK (plan IN ('basic', 'premium', 'trial')),
				trial_start_date TIMESTAMP,
				trial_end_date TIMESTAMP,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_favorite_masters_table",
			sql: `
			CREATE TABLE IF NOT EXISTS favorite_masters (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				master_id INTEGER REFERENCES masters(id) ON DELETE CASCADE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(user_id, master_id)
			);`,
		},
		{
			name: "create_user_cars_table",
			sql: `
			CREATE TABLE IF NOT EXISTS user_cars (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				name VARCHAR(255) NOT NULL,
				year INTEGER,
				comment TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_guarantees_table",
			sql: `
			CREATE TABLE IF NOT EXISTS guarantees (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				appointment_id INTEGER REFERENCES appointments(id) ON DELETE CASCADE,
				service_name VARCHAR(255) NOT NULL,
				master_name VARCHAR(255),
				service_date DATE NOT NULL,
				expiry_date DATE NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
		{
			name: "create_notifications_table",
			sql: `
			CREATE TABLE IF NOT EXISTS notifications (
				id SERIAL PRIMARY KEY,
				user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
				type VARCHAR(50) NOT NULL,
				title VARCHAR(255) NOT NULL,
				message TEXT,
				related_id INTEGER,
				is_read BOOLEAN DEFAULT false,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);`,
		},
	}

	for _, migration := range newMigrations {
		if _, err := db.ExecContext(ctx, migration.sql); err != nil {
			log.Printf("Warning: Failed to run migration %s: %v", migration.name, err)
		} else {
			log.Printf("Migration %s completed", migration.name)
		}
	}

	// Insert sample data
	if err := insertSampleData(db); err != nil {
		log.Printf("Warning: Failed to insert sample data: %v", err)
	}

	log.Println("All migrations completed successfully")
	return nil
}

func insertSampleData(db *sql.DB) error {
	ctx := context.Background()

	// Insert categories
	categories := []struct {
		name, description string
	}{
		{"Автомойка", "Полная мойка и очистка автомобиля"},
		{"Детейлинг", "Детальная обработка кузова и салона"},
		{"Полировка", "Полировка кузова и фар"},
		{"Химчистка салона", "Глубокая чистка салона автомобиля"},
		{"Защитные покрытия", "Нанесение защитных покрытий (керамика, пленка)"},
		{"ТО", "Техническое обслуживание"},
		{"Ремонт", "Ремонтные работы"},
		{"Шиномонтаж", "Замена и ремонт шин"},
		{"Диагностика", "Диагностические работы"},
	}

	for _, cat := range categories {
		_, err := db.ExecContext(ctx, `
			INSERT INTO categories (name, description) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING
		`, cat.name, cat.description)
		if err != nil {
			return err
		}
	}

	// Get category IDs for services
	var autowashID, detailID, polishID, cleaningID, protectionID, toID, repairID, tireID, diagID int
	db.QueryRow("SELECT id FROM categories WHERE name = 'Автомойка'").Scan(&autowashID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Детейлинг'").Scan(&detailID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Полировка'").Scan(&polishID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Химчистка салона'").Scan(&cleaningID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Защитные покрытия'").Scan(&protectionID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'ТО'").Scan(&toID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Ремонт'").Scan(&repairID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Шиномонтаж'").Scan(&tireID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Диагностика'").Scan(&diagID)

	// Insert sample services
	services := []struct {
		categoryID  int
		name        string
		description string
		basePrice   float64
		minPrice    float64
		maxPrice    float64
		duration    int
	}{
		// Автомойка
		{autowashID, "Базовая мойка", "Мойка кузова и колес", 2000, 1500, 3000, 30},
		{autowashID, "Полная мойка", "Мойка кузова, колес, салона", 3500, 2500, 5000, 60},
		{autowashID, "Премиум мойка", "Полная мойка + обработка салона", 5000, 4000, 7000, 90},
		{autowashID, "Мойка двигателя", "Чистка подкапотного пространства", 3000, 2000, 4000, 45},
		{autowashID, "Мойка днища", "Чистка днища автомобиля", 2500, 2000, 3500, 40},

		// Детейлинг
		{detailID, "Детейлинг кузова", "Детальная обработка кузова", 8000, 6000, 12000, 180},
		{detailID, "Детейлинг салона", "Детальная обработка салона", 6000, 4500, 9000, 150},
		{detailID, "Полный детейлинг", "Комплексная обработка кузова и салона", 12000, 10000, 18000, 300},
		{detailID, "Чистка двигателя", "Детальная чистка двигателя", 5000, 4000, 7000, 120},

		// Полировка
		{polishID, "Полировка кузова", "Полировка лакокрасочного покрытия", 12000, 10000, 18000, 240},
		{polishID, "Полировка фар", "Полировка фар и фонарей", 3000, 2000, 5000, 60},
		{polishID, "Восстановительная полировка", "Глубокая полировка с удалением дефектов", 18000, 15000, 25000, 360},
		{polishID, "Защитная полировка", "Полировка с нанесением воска", 15000, 12000, 22000, 300},

		// Химчистка салона
		{cleaningID, "Химчистка салона", "Полная химчистка салона", 4000, 3000, 6000, 120},
		{cleaningID, "Химчистка сидений", "Чистка передних и задних сидений", 2500, 2000, 4000, 90},
		{cleaningID, "Химчистка ковриков", "Глубокая чистка ковриков", 1500, 1000, 2500, 60},
		{cleaningID, "Химчистка багажника", "Чистка багажного отделения", 2000, 1500, 3000, 45},

		// Защитные покрытия
		{protectionID, "Керамическое покрытие", "Нанесение керамического покрытия", 25000, 20000, 35000, 480},
		{protectionID, "Антигравийная пленка", "Нанесение антигравийной пленки", 15000, 12000, 25000, 360},
		{protectionID, "Жидкое стекло", "Нанесение защитного покрытия \"жидкое стекло\"", 8000, 6000, 12000, 180},
		{protectionID, "Бронирование фар", "Защитная пленка на фары", 5000, 4000, 8000, 120},

		// ТО
		{toID, "Замена масла", "Замена моторного масла и фильтра", 2000, 1500, 3000, 60},
		{toID, "Замена фильтров", "Замена воздушного, салонного фильтров", 1500, 1000, 2500, 45},
		{toID, "Замена свечей", "Замена свечей зажигания", 3000, 2000, 5000, 90},
		{toID, "Проверка и доливка жидкостей", "Проверка всех технических жидкостей", 1000, 500, 2000, 30},
		{toID, "Регулировка фар", "Регулировка света фар", 2000, 1500, 3000, 30},

		// Ремонт
		{repairID, "Ремонт двигателя", "Диагностика и ремонт двигателя", 15000, 10000, 30000, 480},
		{repairID, "Ремонт тормозной системы", "Замена колодок, дисков, тормозной жидкости", 8000, 5000, 15000, 240},
		{repairID, "Ремонт подвески", "Диагностика и ремонт элементов подвески", 10000, 7000, 20000, 300},
		{repairID, "Ремонт системы охлаждения", "Ремонт радиатора, помпы, термостата", 6000, 4000, 12000, 180},
		{repairID, "Ремонт электрооборудования", "Ремонт электрики автомобиля", 5000, 3000, 10000, 150},

		// Шиномонтаж
		{tireID, "Замена шин", "Демонтаж и монтаж колес", 5000, 3000, 8000, 120},
		{tireID, "Балансировка колес", "Балансировка всех колес", 1500, 1000, 2500, 60},
		{tireID, "Ремонт прокола", "Ремонт прокола шины", 1000, 500, 2000, 30},
		{tireID, "Перебортовка", "Демонтаж и монтаж шины на диск", 2000, 1500, 3000, 45},
		{tireID, "Хранение шин", "Сезонное хранение шин", 2000, 1500, 3000, 15},

		// Диагностика
		{diagID, "Компьютерная диагностика", "Диагностика всех систем автомобиля", 2000, 1500, 3500, 90},
		{diagID, "Диагностика двигателя", "Проверка работы двигателя", 2500, 2000, 4000, 60},
		{diagID, "Диагностика подвески", "Проверка элементов подвески", 1500, 1000, 2500, 45},
		{diagID, "Диагностика АКПП", "Проверка автоматической коробки передач", 3000, 2500, 5000, 120},
		{diagID, "Диагностика кондиционера", "Проверка системы кондиционирования", 2000, 1500, 3500, 60},
	}

	for _, srv := range services {
		_, err := db.ExecContext(ctx, `
			INSERT INTO services (category_id, name, description, base_price, min_price, max_price, duration_minutes) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT DO NOTHING
		`, srv.categoryID, srv.name, srv.description, srv.basePrice, srv.minPrice, srv.maxPrice, srv.duration)
		if err != nil {
			return err
		}
	}

	// Insert sample cars
	_, err := db.ExecContext(ctx, `
		INSERT INTO cars (brand, model, year, type) VALUES
		('Toyota', 'Camry', 2020, 'Standard'),
		('BMW', 'X5', 2021, 'Premium'),
		('Lada', 'Granta', 2018, 'Standard')
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}

	// Insert sample users for masters (with password hash for "password123")
	// Hash: $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (name, email, phone, password_hash, created_at, updated_at) VALUES
		('Иван Петров', 'master1@beep.kz', '+77001234567', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
		('Сергей Смирнов', 'master2@beep.kz', '+77001234568', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW()),
		('Анна Козлова', 'master3@beep.kz', '+77001234569', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW())
		ON CONFLICT (email) DO NOTHING
	`)
	if err != nil {
		log.Printf("Warning: Failed to insert sample users: %v", err)
	}

	// Get user IDs for masters
	var userID1, userID2, userID3 int
	err = db.QueryRow("SELECT id FROM users WHERE email = 'master1@beep.kz'").Scan(&userID1)
	if err != nil {
		log.Printf("Warning: Failed to get user ID 1: %v", err)
		userID1 = 0
	}
	err = db.QueryRow("SELECT id FROM users WHERE email = 'master2@beep.kz'").Scan(&userID2)
	if err != nil {
		log.Printf("Warning: Failed to get user ID 2: %v", err)
		userID2 = 0
	}
	err = db.QueryRow("SELECT id FROM users WHERE email = 'master3@beep.kz'").Scan(&userID3)
	if err != nil {
		log.Printf("Warning: Failed to get user ID 3: %v", err)
		userID3 = 0
	}

	// Insert sample masters with user_id
	var masterID1, masterID2, masterID3 int
	if userID1 > 0 {
		err = db.QueryRow(`
			INSERT INTO masters (user_id, name, email, phone, specialization, rating, location_lat, location_lng, address, created_at, updated_at)
			VALUES ($1, 'Иван Петров', 'master1@beep.kz', '+77001234567', 'ТО, Ремонт двигателя', 4.8, 43.2220, 76.8512, 'Алматы, ул. Абая, 1', NOW(), NOW())
			ON CONFLICT (email) DO UPDATE SET user_id = EXCLUDED.user_id
			RETURNING id
		`, userID1).Scan(&masterID1)
		if err != nil {
			log.Printf("Warning: Failed to insert master 1: %v", err)
		}
	}
	if userID2 > 0 {
		err = db.QueryRow(`
			INSERT INTO masters (user_id, name, email, phone, specialization, rating, location_lat, location_lng, address, created_at, updated_at)
			VALUES ($1, 'Сергей Смирнов', 'master2@beep.kz', '+77001234568', 'Шиномонтаж, Замена шин', 4.7, 43.2200, 76.8500, 'Алматы, ул. Достык, 2', NOW(), NOW())
			ON CONFLICT (email) DO UPDATE SET user_id = EXCLUDED.user_id
			RETURNING id
		`, userID2).Scan(&masterID2)
		if err != nil {
			log.Printf("Warning: Failed to insert master 2: %v", err)
		}
	}
	if userID3 > 0 {
		err = db.QueryRow(`
			INSERT INTO masters (user_id, name, email, phone, specialization, rating, location_lat, location_lng, address, created_at, updated_at)
			VALUES ($1, 'Анна Козлова', 'master3@beep.kz', '+77001234569', 'Диагностика, Ремонт электрооборудования', 4.9, 43.2180, 76.8488, 'Алматы, пр. Абая, 3', NOW(), NOW())
			ON CONFLICT (email) DO UPDATE SET user_id = EXCLUDED.user_id
			RETURNING id
		`, userID3).Scan(&masterID3)
		if err != nil {
			log.Printf("Warning: Failed to insert master 3: %v", err)
		}
	}

	// Insert sample master schedules
	_, err = db.ExecContext(ctx, `
		INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at) VALUES
		($1, 1, '09:00', '18:00', true, NOW()),
		($1, 2, '09:00', '18:00', true, NOW()),
		($1, 3, '09:00', '18:00', true, NOW()),
		($1, 4, '09:00', '18:00', true, NOW()),
		($1, 5, '09:00', '18:00', true, NOW()),
		($2, 1, '08:00', '17:00', true, NOW()),
		($2, 2, '08:00', '17:00', true, NOW()),
		($2, 3, '08:00', '17:00', true, NOW()),
		($2, 4, '08:00', '17:00', true, NOW()),
		($2, 5, '08:00', '17:00', true, NOW()),
		($3, 1, '10:00', '19:00', true, NOW()),
		($3, 2, '10:00', '19:00', true, NOW()),
		($3, 3, '10:00', '19:00', true, NOW()),
		($3, 4, '10:00', '19:00', true, NOW()),
		($3, 5, '10:00', '19:00', true, NOW())
		ON CONFLICT DO NOTHING
	`, masterID1, masterID2, masterID3)
	if err != nil {
		log.Printf("Warning: Failed to insert master schedules: %v", err)
	}

	// Insert sample works for masters
	if masterID1 > 0 {
		_, err = db.ExecContext(ctx, `
			INSERT INTO master_works (master_id, title, work_date, customer_name, amount, photo_urls, created_at) VALUES
			($1, 'Замена масла и фильтров', NOW() - INTERVAL '5 days', 'Александр Иванов', 3000, NULL, NOW()),
			($1, 'Ремонт системы охлаждения', NOW() - INTERVAL '10 days', 'Мария Петрова', 8000, NULL, NOW()),
			($1, 'Замена ремня ГРМ', NOW() - INTERVAL '15 days', 'Дмитрий Сидоров', 15000, NULL, NOW())
			ON CONFLICT DO NOTHING
		`, masterID1)
		if err != nil {
			log.Printf("Warning: Failed to insert works for master 1: %v", err)
		}
	}
	if masterID2 > 0 {
		_, err = db.ExecContext(ctx, `
			INSERT INTO master_works (master_id, title, work_date, customer_name, amount, photo_urls, created_at) VALUES
			($1, 'Балансировка колес', NOW() - INTERVAL '3 days', 'Елена Козлова', 2000, NULL, NOW()),
			($1, 'Замена шин', NOW() - INTERVAL '7 days', 'Игорь Морозов', 5000, NULL, NOW()),
			($1, 'Ремонт прокола', NOW() - INTERVAL '12 days', 'Ольга Волкова', 1000, NULL, NOW())
			ON CONFLICT DO NOTHING
		`, masterID2)
		if err != nil {
			log.Printf("Warning: Failed to insert works for master 2: %v", err)
		}
	}
	if masterID3 > 0 {
		_, err = db.ExecContext(ctx, `
			INSERT INTO master_works (master_id, title, work_date, customer_name, amount, photo_urls, created_at) VALUES
			($1, 'Компьютерная диагностика', NOW() - INTERVAL '2 days', 'Антон Новиков', 2500, NULL, NOW()),
			($1, 'Ремонт генератора', NOW() - INTERVAL '8 days', 'Татьяна Лебедева', 8000, NULL, NOW()),
			($1, 'Замена аккумулятора', NOW() - INTERVAL '14 days', 'Вадим Орлов', 5000, NULL, NOW())
			ON CONFLICT DO NOTHING
		`, masterID3)
		if err != nil {
			log.Printf("Warning: Failed to insert works for master 3: %v", err)
		}
	}

	// Insert sample reviews for masters (masters review each other)
	if masterID1 > 0 && masterID2 > 0 && masterID3 > 0 {
		_, err = db.ExecContext(ctx, `
			INSERT INTO reviews (master_id, user_id, rating, comment, created_at) VALUES
			($1, $2, 5, 'Отличный мастер! Быстро и качественно выполнил работу.', NOW() - INTERVAL '5 days'),
			($1, $3, 5, 'Профессиональный подход, рекомендую!', NOW() - INTERVAL '10 days'),
			($2, $1, 5, 'Отлично сделал балансировку, машина теперь не трясется!', NOW() - INTERVAL '3 days'),
			($2, $3, 5, 'Быстрая замена шин, все отлично!', NOW() - INTERVAL '7 days'),
			($3, $1, 5, 'Отличная диагностика! Нашла все проблемы.', NOW() - INTERVAL '2 days'),
			($3, $2, 5, 'Профессионал своего дела!', NOW() - INTERVAL '8 days')
			ON CONFLICT DO NOTHING
		`, masterID1, userID2, userID3, masterID2, userID1, userID3, masterID3, userID1, userID2)
		if err != nil {
			log.Printf("Warning: Failed to insert reviews: %v", err)
		}
	}

	// Update master ratings based on reviews (only if masters were created)
	if masterID1 > 0 || masterID2 > 0 || masterID3 > 0 {
		var masterIDs []interface{}
		if masterID1 > 0 {
			masterIDs = append(masterIDs, masterID1)
		}
		if masterID2 > 0 {
			masterIDs = append(masterIDs, masterID2)
		}
		if masterID3 > 0 {
			masterIDs = append(masterIDs, masterID3)
		}

		if len(masterIDs) > 0 {
			query := "UPDATE masters SET rating = (SELECT COALESCE(AVG(rating)::decimal, 0.0) FROM reviews WHERE reviews.master_id = masters.id) WHERE id IN ("
			for i := range masterIDs {
				if i > 0 {
					query += ","
				}
				query += fmt.Sprintf("$%d", i+1)
			}
			query += ")"

			_, err = db.ExecContext(ctx, query, masterIDs...)
			if err != nil {
				log.Printf("Warning: Failed to update master ratings: %v", err)
			}
		}
	}

	log.Println("Sample data inserted successfully")
	return nil
}
