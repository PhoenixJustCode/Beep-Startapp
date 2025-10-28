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
	var toID, repairID, tireID, diagID int
	db.QueryRow("SELECT id FROM categories WHERE name = 'ТО'").Scan(&toID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Ремонт'").Scan(&repairID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Шиномонтаж'").Scan(&tireID)
	db.QueryRow("SELECT id FROM categories WHERE name = 'Диагностика'").Scan(&diagID)

	// Insert sample services
	services := []struct {
		categoryID int
		name       string
		basePrice  float64
		minPrice   float64
		maxPrice   float64
		duration   int
	}{
		{toID, "Замена масла", 2000, 1500, 3000, 60},
		{toID, "Замена фильтров", 1500, 1000, 2500, 45},
		{repairID, "Ремонт двигателя", 15000, 10000, 30000, 480},
		{repairID, "Ремонт тормозной системы", 8000, 5000, 15000, 240},
		{tireID, "Замена шин", 5000, 3000, 8000, 120},
		{tireID, "Балансировка колес", 1500, 1000, 2500, 60},
		{diagID, "Компьютерная диагностика", 2000, 1500, 3500, 90},
	}

	for _, srv := range services {
		_, err := db.ExecContext(ctx, `
			INSERT INTO services (category_id, name, base_price, min_price, max_price, duration_minutes) 
			VALUES ($1, $2, $3, $4, $5, $6)
		`, srv.categoryID, srv.name, srv.basePrice, srv.minPrice, srv.maxPrice, srv.duration)
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

	// Insert sample masters
	_, err = db.ExecContext(ctx, `
		INSERT INTO masters (name, email, phone, specialization, rating, location_lat, location_lng, address) VALUES
		('Иван Петров', 'ivan@example.com', '+79001234567', 'ТО, Ремонт', 4.8, 55.7558, 37.6173, 'Москва, ул. Примерная, 1'),
		('Сергей Смирнов', 'sergey@example.com', '+79001234568', 'Шиномонтаж', 4.7, 55.7540, 37.6200, 'Москва, ул. Примерная, 2'),
		('Анна Козлова', 'anna@example.com', '+79001234569', 'Диагностика', 4.9, 55.7530, 37.6185, 'Москва, ул. Примерная, 3')
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}

	// Insert sample master schedules
	_, err = db.ExecContext(ctx, `
		INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time) VALUES
		(1, 1, '09:00', '18:00'),
		(1, 2, '09:00', '18:00'),
		(1, 3, '09:00', '18:00'),
		(1, 4, '09:00', '18:00'),
		(1, 5, '09:00', '18:00'),
		(2, 1, '08:00', '17:00'),
		(2, 2, '08:00', '17:00'),
		(2, 3, '08:00', '17:00'),
		(2, 4, '08:00', '17:00'),
		(2, 5, '08:00', '17:00'),
		(3, 1, '10:00', '19:00'),
		(3, 2, '10:00', '19:00'),
		(3, 3, '10:00', '19:00'),
		(3, 4, '10:00', '19:00'),
		(3, 5, '10:00', '19:00')
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}

	log.Println("Sample data inserted successfully")
	return nil
}
