package main

import (
	"database/sql"
	"log"

	"beep-backend/internal/config"
	"beep-backend/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables (use default or from .env file)
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Printf("Connected to database: %s", cfg.DatabaseURL)

	// Run migrations first to ensure tables exist
	log.Println("Running migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Printf("Warning: Failed to run migrations (tables may already exist): %v", err)
	}

	// Add masters
	if err := addMasters(db); err != nil {
		log.Fatal("Failed to add masters:", err)
	}

	log.Println("Successfully added 3 masters to database")
}

func addMasters(db *sql.DB) error {
	// First, create users for masters
	userQueries := []struct {
		name, email, phone string
	}{
		{"Асхат Мухамеджанов", "askhat.master@example.com", "+7-700-111-2222"},
		{"Дамир Смагулов", "damir.master@example.com", "+7-700-333-4444"},
		{"Нурлан Абдуллаев", "nurlan.master@example.com", "+7-700-555-6666"},
	}

	passwordHash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"

	for _, user := range userQueries {
		// Insert user if not exists
		var userID int
		err := db.QueryRow(`
			INSERT INTO users (name, email, phone, password_hash, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, NOW(), NOW())
			ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
			RETURNING id
		`, user.name, user.email, user.phone, passwordHash).Scan(&userID)

		if err != nil {
			// If user exists, get its ID
			err = db.QueryRow("SELECT id FROM users WHERE email = $1", user.email).Scan(&userID)
			if err != nil {
				return err
			}
		}

		// Create master profile
		var spec string
		var rating float64
		var lat, lng float64
		var address string

		switch user.email {
		case "askhat.master@example.com":
			spec = "Автомойка и полировка"
			rating = 4.9
			lat = 43.238949
			lng = 76.889709
			address = "ул. Абая, 150, Алматы"
		case "damir.master@example.com":
			spec = "Детейлинг и защитные покрытия"
			rating = 4.8
			lat = 43.250000
			lng = 76.950000
			address = "пр. Аль-Фараби, 77, Алматы"
		case "nurlan.master@example.com":
			spec = "Химчистка и автодетейлинг"
			rating = 4.7
			lat = 43.220000
			lng = 76.870000
			address = "ул. Сатпаева, 90, Алматы"
		}

		_, err = db.Exec(`
			INSERT INTO masters (user_id, name, email, phone, specialization, rating, location_lat, location_lng, address, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
			ON CONFLICT (user_id) DO NOTHING
		`, userID, user.name, user.email, user.phone, spec, rating, lat, lng, address)

		if err != nil {
			log.Printf("Warning: Failed to insert master for %s: %v", user.email, err)
			continue
		}

		log.Printf("Added master: %s", user.name)

		// Get master ID to add schedule
		var masterID int
		err = db.QueryRow("SELECT id FROM masters WHERE user_id = $1", userID).Scan(&masterID)
		if err != nil {
			log.Printf("Warning: Failed to get master ID for %s: %v", user.email, err)
			continue
		}

		// Add schedule for each day of week (0-6)
		schedules := []struct {
			day       int
			startTime string
			endTime   string
		}{
			{0, "10:00", "16:00"}, // Sunday
			{1, "08:00", "20:00"}, // Monday
			{2, "08:00", "20:00"}, // Tuesday
			{3, "08:00", "20:00"}, // Wednesday
			{4, "08:00", "20:00"}, // Thursday
			{5, "08:00", "20:00"}, // Friday
			{6, "09:00", "18:00"}, // Saturday
		}

		// Adjust schedule based on master
		if user.email == "damir.master@example.com" {
			schedules[0] = struct {
				day                int
				startTime, endTime string
			}{0, "11:00", "17:00"}
			schedules[1] = struct {
				day                int
				startTime, endTime string
			}{1, "09:00", "21:00"}
			schedules[2] = struct {
				day                int
				startTime, endTime string
			}{2, "09:00", "21:00"}
			schedules[3] = struct {
				day                int
				startTime, endTime string
			}{3, "09:00", "21:00"}
			schedules[4] = struct {
				day                int
				startTime, endTime string
			}{4, "09:00", "21:00"}
			schedules[5] = struct {
				day                int
				startTime, endTime string
			}{5, "09:00", "21:00"}
			schedules[6] = struct {
				day                int
				startTime, endTime string
			}{6, "10:00", "19:00"}
		} else if user.email == "nurlan.master@example.com" {
			schedules[0] = struct {
				day                int
				startTime, endTime string
			}{0, "09:00", "15:00"}
			schedules[1] = struct {
				day                int
				startTime, endTime string
			}{1, "07:00", "19:00"}
			schedules[2] = struct {
				day                int
				startTime, endTime string
			}{2, "07:00", "19:00"}
			schedules[3] = struct {
				day                int
				startTime, endTime string
			}{3, "07:00", "19:00"}
			schedules[4] = struct {
				day                int
				startTime, endTime string
			}{4, "07:00", "19:00"}
			schedules[5] = struct {
				day                int
				startTime, endTime string
			}{5, "07:00", "19:00"}
			schedules[6] = struct {
				day                int
				startTime, endTime string
			}{6, "08:00", "17:00"}
		}

		for _, schedule := range schedules {
			_, err = db.Exec(`
				INSERT INTO master_schedule (master_id, day_of_week, start_time, end_time, is_active, created_at)
				VALUES ($1, $2, $3, $4, true, NOW())
				ON CONFLICT DO NOTHING
			`, masterID, schedule.day, schedule.startTime, schedule.endTime)

			if err != nil {
				log.Printf("Warning: Failed to add schedule for master %d, day %d: %v", masterID, schedule.day, err)
			}
		}
	}

	return nil
}
