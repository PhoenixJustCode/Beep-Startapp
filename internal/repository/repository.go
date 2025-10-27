package repository

import (
	"beep-backend/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Categories
func (r *Repository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *Repository) GetCategoryByID(id int) (*models.Category, error) {
	var cat models.Category
	err := r.db.QueryRow("SELECT id, name, description, created_at FROM categories WHERE id = $1", id).
		Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Services
func (r *Repository) GetAllServices() ([]models.Service, error) {
	rows, err := r.db.Query("SELECT id, category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at FROM services ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var srv models.Service
		if err := rows.Scan(&srv.ID, &srv.CategoryID, &srv.Name, &srv.Description, &srv.BasePrice, &srv.MinPrice, &srv.MaxPrice, &srv.DurationMinutes, &srv.CreatedAt); err != nil {
			return nil, err
		}
		services = append(services, srv)
	}
	return services, nil
}

func (r *Repository) GetServicesByCategory(categoryID int) ([]models.Service, error) {
	rows, err := r.db.Query("SELECT id, category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at FROM services WHERE category_id = $1 ORDER BY name", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var srv models.Service
		if err := rows.Scan(&srv.ID, &srv.CategoryID, &srv.Name, &srv.Description, &srv.BasePrice, &srv.MinPrice, &srv.MaxPrice, &srv.DurationMinutes, &srv.CreatedAt); err != nil {
			return nil, err
		}
		services = append(services, srv)
	}
	return services, nil
}

func (r *Repository) GetServiceByID(id int) (*models.Service, error) {
	var srv models.Service
	err := r.db.QueryRow("SELECT id, category_id, name, description, base_price, min_price, max_price, duration_minutes, created_at FROM services WHERE id = $1", id).
		Scan(&srv.ID, &srv.CategoryID, &srv.Name, &srv.Description, &srv.BasePrice, &srv.MinPrice, &srv.MaxPrice, &srv.DurationMinutes, &srv.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &srv, nil
}

// Cars
func (r *Repository) GetAllCars() ([]models.Car, error) {
	rows, err := r.db.Query("SELECT id, brand, model, year, type FROM cars ORDER BY brand, model")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Type); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (r *Repository) GetCarByID(id int) (*models.Car, error) {
	var car models.Car
	err := r.db.QueryRow("SELECT id, brand, model, year, type FROM cars WHERE id = $1", id).
		Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Type)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

// Pricing calculation
func (r *Repository) CalculatePrice(serviceID, carID int) (*models.CalculatePriceResponse, error) {
	var service models.Service
	err := r.db.QueryRow("SELECT id, name, base_price, min_price, max_price FROM services WHERE id = $1", serviceID).
		Scan(&service.ID, &service.Name, &service.BasePrice, &service.MinPrice, &service.MaxPrice)
	if err != nil {
		return nil, err
	}

	var car models.Car
	err = r.db.QueryRow("SELECT id, brand, model, year, type FROM cars WHERE id = $1", carID).
		Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Type)
	if err != nil {
		return nil, err
	}

	// Calculate price with multipliers
	carAge := 2024 - car.Year
	price := service.BasePrice

	// Apply age-based multipliers
	if carAge > 10 {
		price *= 1.2 // Older cars are more expensive to service
	} else if carAge < 3 {
		price *= 0.9 // Newer cars are cheaper
	}

	// Apply car type multipliers
	if car.Type == "Premium" {
		price *= 1.5
	}

	response := &models.CalculatePriceResponse{
		ServiceID:   service.ID,
		ServiceName: service.Name,
		BasePrice:   service.BasePrice,
		FinalPrice:  price,
		MinPrice:    service.MinPrice,
		MaxPrice:    service.MaxPrice,
	}

	return response, nil
}

// Masters
func (r *Repository) GetAllMasters() ([]models.Master, error) {
	rows, err := r.db.Query("SELECT id, name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at FROM masters ORDER BY rating DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var masters []models.Master
	for rows.Next() {
		var master models.Master
		if err := rows.Scan(&master.ID, &master.Name, &master.Email, &master.Phone, &master.Specialization, &master.Rating, &master.PhotoURL, &master.LocationLat, &master.LocationLng, &master.Address, &master.CreatedAt, &master.UpdatedAt); err != nil {
			return nil, err
		}
		masters = append(masters, master)
	}
	return masters, nil
}

func (r *Repository) GetMasterByID(id int) (*models.Master, error) {
	var master models.Master
	err := r.db.QueryRow("SELECT id, name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at FROM masters WHERE id = $1", id).
		Scan(&master.ID, &master.Name, &master.Email, &master.Phone, &master.Specialization, &master.Rating, &master.PhotoURL, &master.LocationLat, &master.LocationLng, &master.Address, &master.CreatedAt, &master.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &master, nil
}

func (r *Repository) GetMasterReviews(masterID int) ([]models.Review, error) {
	rows, err := r.db.Query("SELECT id, master_id, user_id, rating, comment, created_at FROM reviews WHERE master_id = $1 ORDER BY created_at DESC", masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.ID, &review.MasterID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *Repository) GetMasterSchedule(masterID int) ([]models.MasterSchedule, error) {
	rows, err := r.db.Query("SELECT id, master_id, day_of_week, start_time, end_time, is_active, created_at FROM master_schedule WHERE master_id = $1 AND is_active = true", masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.MasterSchedule
	for rows.Next() {
		var schedule models.MasterSchedule
		if err := rows.Scan(&schedule.ID, &schedule.MasterID, &schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime, &schedule.IsActive, &schedule.CreatedAt); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// Get available time slots for a master on a specific date
func (r *Repository) GetAvailableSlots(masterID int, date time.Time) ([]string, error) {
	// Get master's schedule for that day of week
	dayOfWeek := int(date.Weekday())
	dayOfWeek = (dayOfWeek + 6) % 7 // PostgreSQL uses 0=Monday, 6=Sunday

	var startTime, endTime string
	err := r.db.QueryRow("SELECT start_time, end_time FROM master_schedule WHERE master_id = $1 AND day_of_week = $2 AND is_active = true", masterID, dayOfWeek).
		Scan(&startTime, &endTime)
	if err == sql.ErrNoRows {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}

	// Get existing appointments for that date
	existingAppointments, err := r.db.Query("SELECT time FROM appointments WHERE master_id = $1 AND date = $2 AND status != 'cancelled'", masterID, date.Format("2006-01-02"))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	bookedTimes := make(map[string]bool)
	if err == nil {
		for existingAppointments.Next() {
			var timeStr string
			if err := existingAppointments.Scan(&timeStr); err != nil {
				return nil, err
			}
			bookedTimes[timeStr] = true
		}
	}

	// Generate available time slots (every hour from 8:00 to 18:00)
	slots := []string{}
	for hour := 8; hour <= 18; hour++ {
		slot := fmt.Sprintf("%02d:00", hour)
		if !bookedTimes[slot] {
			slots = append(slots, slot)
		}
	}

	return slots, nil
}

// Appointments
func (r *Repository) CreateAppointment(userID, masterID, serviceID int, date time.Time, timeStr, comment string) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.QueryRow(
		`INSERT INTO appointments (user_id, master_id, service_id, date, time, status, comment) 
		 VALUES ($1, $2, $3, $4, $5, 'pending', $6) 
		 RETURNING id, user_id, master_id, service_id, date, time, status, comment, created_at, updated_at`,
		userID, masterID, serviceID, date.Format("2006-01-02"), timeStr, comment,
	).Scan(&appointment.ID, &appointment.UserID, &appointment.MasterID, &appointment.ServiceID, &appointment.Date, &appointment.Time, &appointment.Status, &appointment.Comment, &appointment.CreatedAt, &appointment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *Repository) GetUserAppointments(userID int) ([]models.AppointmentWithDetails, error) {
	query := `
		SELECT 
			a.id, a.user_id, a.master_id, a.service_id, a.date, a.time, a.status, a.comment, a.created_at, a.updated_at,
			s.name as service_name,
			m.name as master_name
		FROM appointments a
		LEFT JOIN services s ON a.service_id = s.id
		LEFT JOIN masters m ON a.master_id = m.id
		WHERE a.user_id = $1 
		ORDER BY a.date DESC, a.time DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.AppointmentWithDetails
	for rows.Next() {
		var appt models.AppointmentWithDetails
		if err := rows.Scan(&appt.ID, &appt.UserID, &appt.MasterID, &appt.ServiceID, &appt.Date, &appt.Time, &appt.Status, &appt.Comment, &appt.CreatedAt, &appt.UpdatedAt, &appt.ServiceName, &appt.MasterName); err != nil {
			return nil, err
		}
		appointments = append(appointments, appt)
	}
	return appointments, nil
}

func (r *Repository) UpdateAppointment(appointmentID int, comment string) error {
	_, err := r.db.Exec("UPDATE appointments SET comment = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", comment, appointmentID)
	return err
}

func (r *Repository) CancelAppointment(appointmentID int) error {
	_, err := r.db.Exec("UPDATE appointments SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = $1", appointmentID)
	return err
}

func (r *Repository) GetAppointmentByID(appointmentID int) (*models.Appointment, error) {
	var appt models.Appointment
	err := r.db.QueryRow("SELECT id, user_id, master_id, service_id, date, time, status, comment, created_at, updated_at FROM appointments WHERE id = $1", appointmentID).
		Scan(&appt.ID, &appt.UserID, &appt.MasterID, &appt.ServiceID, &appt.Date, &appt.Time, &appt.Status, &appt.Comment, &appt.CreatedAt, &appt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &appt, nil
}

// Users
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, name, email, phone, password_hash, created_at, updated_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(name, email, phone, passwordHash string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`INSERT INTO users (name, email, phone, password_hash) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, name, email, phone, password_hash, created_at, updated_at`,
		name, email, phone, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
