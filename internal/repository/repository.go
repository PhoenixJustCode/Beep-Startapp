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

	// Calculate car age
	carAge := 2024 - car.Year
	price := service.BasePrice

	// Create price details array
	var priceDetails []models.PriceDetail

	// Add base price detail
	priceDetails = append(priceDetails, models.PriceDetail{
		Description: "Базовая цена услуги",
		Amount:      service.BasePrice,
		IsAddition:  false,
	})

	// Apply age-based multipliers
	if carAge > 10 {
		oldCarMultiplier := 1.2
		oldCarAmount := service.BasePrice * (oldCarMultiplier - 1)
		price *= oldCarMultiplier
		priceDetails = append(priceDetails, models.PriceDetail{
			Description: "Надбавка за старый автомобиль (+20%)",
			Amount:      oldCarAmount,
			Multiplier:  oldCarMultiplier,
			IsAddition:  true,
		})
	} else if carAge < 3 {
		newCarMultiplier := 0.9
		newCarDiscount := service.BasePrice * (1 - newCarMultiplier)
		price *= newCarMultiplier
		priceDetails = append(priceDetails, models.PriceDetail{
			Description: "Скидка за новый автомобиль (-10%)",
			Amount:      -newCarDiscount,
			Multiplier:  newCarMultiplier,
			IsAddition:  true,
		})
	}

	// Apply car type multipliers
	if car.Type == "Premium" {
		premiumMultiplier := 1.5
		premiumAmount := service.BasePrice * (premiumMultiplier - 1)
		price *= premiumMultiplier
		priceDetails = append(priceDetails, models.PriceDetail{
			Description: "Надбавка за премиум автомобиль (+50%)",
			Amount:      premiumAmount,
			Multiplier:  premiumMultiplier,
			IsAddition:  true,
		})
	} else if car.Type == "Luxury" {
		luxuryMultiplier := 2.0
		luxuryAmount := service.BasePrice * (luxuryMultiplier - 1)
		price *= luxuryMultiplier
		priceDetails = append(priceDetails, models.PriceDetail{
			Description: "Надбавка за люксовый автомобиль (+100%)",
			Amount:      luxuryAmount,
			Multiplier:  luxuryMultiplier,
			IsAddition:  true,
		})
	}

	response := &models.CalculatePriceResponse{
		ServiceID:    service.ID,
		ServiceName:  service.Name,
		CarBrand:     car.Brand,
		CarModel:     car.Model,
		CarYear:      car.Year,
		CarType:      car.Type,
		CarAge:       carAge,
		BasePrice:    service.BasePrice,
		FinalPrice:   price,
		MinPrice:     service.MinPrice,
		MaxPrice:     service.MaxPrice,
		PriceDetails: priceDetails,
	}

	return response, nil
}

// Update user profile
func (r *Repository) UpdateUserProfile(userID int, name, email, phone string) error {
	_, err := r.db.Exec("UPDATE users SET name = $1, email = $2, phone = $3, updated_at = NOW() WHERE id = $4",
		name, email, phone, userID)
	return err
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
		var specialization, photoURL, address sql.NullString
		var rating sql.NullFloat64
		var locationLat, locationLng sql.NullFloat64

		if err := rows.Scan(&master.ID, &master.Name, &master.Email, &master.Phone,
			&specialization, &rating, &photoURL, &locationLat, &locationLng,
			&address, &master.CreatedAt, &master.UpdatedAt); err != nil {
			return nil, err
		}

		// Handle nullable fields
		if specialization.Valid {
			master.Specialization = specialization.String
		}
		if rating.Valid {
			master.Rating = rating.Float64
		}
		if photoURL.Valid {
			master.PhotoURL = photoURL.String
		}
		if locationLat.Valid {
			master.LocationLat = locationLat.Float64
		}
		if locationLng.Valid {
			master.LocationLng = locationLng.Float64
		}
		if address.Valid {
			master.Address = address.String
		}

		masters = append(masters, master)
	}
	return masters, nil
}

func (r *Repository) GetMasterByID(id int) (*models.Master, error) {
	var master models.Master
	var specialization, photoURL, address sql.NullString
	var rating sql.NullFloat64
	var locationLat, locationLng sql.NullFloat64

	err := r.db.QueryRow("SELECT id, name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at FROM masters WHERE id = $1", id).
		Scan(&master.ID, &master.Name, &master.Email, &master.Phone,
			&specialization, &rating, &photoURL, &locationLat, &locationLng,
			&address, &master.CreatedAt, &master.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if specialization.Valid {
		master.Specialization = specialization.String
	}
	if rating.Valid {
		master.Rating = rating.Float64
	}
	if photoURL.Valid {
		master.PhotoURL = photoURL.String
	}
	if locationLat.Valid {
		master.LocationLat = locationLat.Float64
	}
	if locationLng.Valid {
		master.LocationLng = locationLng.Float64
	}
	if address.Valid {
		master.Address = address.String
	}

	return &master, nil
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

func (r *Repository) DeleteAppointment(appointmentID int) error {
	_, err := r.db.Exec("DELETE FROM appointments WHERE id = $1", appointmentID)
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
	var photoURL sql.NullString
	err := r.db.QueryRow("SELECT id, name, email, phone, photo_url, password_hash, created_at, updated_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &photoURL, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	return &user, nil
}

func (r *Repository) CreateUser(name, email, phone, passwordHash string) (*models.User, error) {
	var user models.User
	var photoURL sql.NullString
	err := r.db.QueryRow(
		`INSERT INTO users (name, email, phone, password_hash) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, name, email, phone, photo_url, password_hash, created_at, updated_at`,
		name, email, phone, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &photoURL, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	var photoURL sql.NullString
	err := r.db.QueryRow("SELECT id, name, email, phone, photo_url, password_hash, created_at, updated_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &photoURL, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	return &user, nil
}

// Update user photo URL
func (r *Repository) UpdateUserPhoto(userID int, photoURL string) error {
	_, err := r.db.Exec("UPDATE users SET photo_url = $1, updated_at = NOW() WHERE id = $2", photoURL, userID)
	return err
}

// Update master photo URL
func (r *Repository) UpdateMasterPhoto(masterID int, photoURL string) error {
	_, err := r.db.Exec("UPDATE masters SET photo_url = $1, updated_at = NOW() WHERE id = $2", photoURL, masterID)
	return err
}

// Master Profile Methods

// GetMasterByUserID gets master profile by user ID
func (r *Repository) GetMasterByUserID(userID int) (*models.Master, error) {
	var master models.Master
	var specialization, photoURL, address sql.NullString
	var rating sql.NullFloat64
	var locationLat, locationLng sql.NullFloat64

	err := r.db.QueryRow("SELECT id, user_id, name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at FROM masters WHERE user_id = $1", userID).
		Scan(&master.ID, &master.UserID, &master.Name, &master.Email, &master.Phone,
			&specialization, &rating, &photoURL, &locationLat, &locationLng,
			&address, &master.CreatedAt, &master.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if specialization.Valid {
		master.Specialization = specialization.String
	}
	if rating.Valid {
		master.Rating = rating.Float64
	}
	if photoURL.Valid {
		master.PhotoURL = photoURL.String
	}
	if locationLat.Valid {
		master.LocationLat = locationLat.Float64
	}
	if locationLng.Valid {
		master.LocationLng = locationLng.Float64
	}
	if address.Valid {
		master.Address = address.String
	}

	return &master, nil
}

// CreateMaster creates a new master profile
func (r *Repository) CreateMaster(userID int, name, email, phone, specialization, address string) (*models.Master, error) {
	var master models.Master
	var specializationNull, photoURLNull, addressNull sql.NullString
	var locationLatNull, locationLngNull sql.NullFloat64
	var ratingNull sql.NullFloat64

	err := r.db.QueryRow(`
		INSERT INTO masters (user_id, name, email, phone, specialization, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, user_id, name, email, phone, specialization, rating, photo_url, location_lat, location_lng, address, created_at, updated_at
	`, userID, name, email, phone, specialization, address).
		Scan(&master.ID, &master.UserID, &master.Name, &master.Email, &master.Phone,
			&specializationNull, &ratingNull, &photoURLNull, &locationLatNull, &locationLngNull,
			&addressNull, &master.CreatedAt, &master.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if specializationNull.Valid {
		master.Specialization = specializationNull.String
	}
	if ratingNull.Valid {
		master.Rating = ratingNull.Float64
	}
	if photoURLNull.Valid {
		master.PhotoURL = photoURLNull.String
	}
	if addressNull.Valid {
		master.Address = addressNull.String
	}
	if locationLatNull.Valid {
		master.LocationLat = locationLatNull.Float64
	}
	if locationLngNull.Valid {
		master.LocationLng = locationLngNull.Float64
	}

	return &master, nil
}

// UpdateMaster updates master profile
func (r *Repository) UpdateMaster(masterID int, name, email, phone, specialization, address string) error {
	_, err := r.db.Exec(`
		UPDATE masters SET name = $1, email = $2, phone = $3, specialization = $4, address = $5, updated_at = NOW()
		WHERE id = $6
	`, name, email, phone, specialization, address, masterID)
	return err
}

// DeleteMasterProfile deletes master profile by user ID
func (r *Repository) DeleteMasterProfile(userID int) error {
	_, err := r.db.Exec("DELETE FROM masters WHERE user_id = $1", userID)
	return err
}

// Master Works Methods

// GetMasterWorks gets all works for a master
func (r *Repository) GetMasterWorks(masterID int) ([]models.MasterWork, error) {
	rows, err := r.db.Query(`
		SELECT id, master_id, title, work_date, customer_name, amount, photo_urls, created_at
		FROM master_works WHERE master_id = $1 ORDER BY work_date DESC
	`, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []models.MasterWork
	for rows.Next() {
		var work models.MasterWork
		var photoURLs sql.NullString
		if err := rows.Scan(&work.ID, &work.MasterID, &work.Title, &work.WorkDate,
			&work.CustomerName, &work.Amount, &photoURLs, &work.CreatedAt); err != nil {
			return nil, err
		}
		// Parse photo URLs array (simplified - in real app would use proper array handling)
		if photoURLs.Valid && photoURLs.String != "" {
			work.PhotoURLs = []string{photoURLs.String} // Simplified for now
		}
		works = append(works, work)
	}
	return works, nil
}

// CreateMasterWork creates a new work entry
func (r *Repository) CreateMasterWork(masterID int, title string, workDate time.Time, customerName string, amount float64, photoURLs []string) (*models.MasterWork, error) {
	var work models.MasterWork
	var photoURLsStr sql.NullString

	// Handle photo URLs - convert to PostgreSQL array format
	if len(photoURLs) > 0 {
		photoURLsStr = sql.NullString{String: photoURLs[0], Valid: true} // Simplified for now
	}

	err := r.db.QueryRow(`
		INSERT INTO master_works (master_id, title, work_date, customer_name, amount, photo_urls, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, master_id, title, work_date, customer_name, amount, photo_urls, created_at
	`, masterID, title, workDate, customerName, amount, photoURLsStr).
		Scan(&work.ID, &work.MasterID, &work.Title, &work.WorkDate,
			&work.CustomerName, &work.Amount, &photoURLsStr, &work.CreatedAt)
	if err != nil {
		return nil, err
	}

	if photoURLsStr.Valid && photoURLsStr.String != "" {
		work.PhotoURLs = []string{photoURLsStr.String}
	}
	return &work, nil
}

// GetMasterWork gets a single work by ID
func (r *Repository) GetMasterWork(workID, masterID int) (*models.MasterWork, error) {
	var work models.MasterWork
	var photoURLs sql.NullString

	err := r.db.QueryRow(`
		SELECT id, master_id, title, work_date, customer_name, amount, photo_urls, created_at
		FROM master_works WHERE id = $1 AND master_id = $2
	`, workID, masterID).
		Scan(&work.ID, &work.MasterID, &work.Title, &work.WorkDate,
			&work.CustomerName, &work.Amount, &photoURLs, &work.CreatedAt)
	if err != nil {
		return nil, err
	}

	if photoURLs.Valid && photoURLs.String != "" {
		work.PhotoURLs = []string{photoURLs.String}
	}

	return &work, nil
}

// UpdateMasterWork updates a work entry
func (r *Repository) UpdateMasterWork(workID, masterID int, title string, workDate time.Time, customerName string, amount float64, photoURLs []string) error {
	var photoURLsStr sql.NullString

	if len(photoURLs) > 0 {
		photoURLsStr = sql.NullString{String: photoURLs[0], Valid: true}
	}

	result, err := r.db.Exec(`
		UPDATE master_works SET title = $1, work_date = $2, customer_name = $3, amount = $4, photo_urls = $5
		WHERE id = $6 AND master_id = $7
	`, title, workDate, customerName, amount, photoURLsStr, workID, masterID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteMasterWork deletes a work entry
func (r *Repository) DeleteMasterWork(workID, masterID int) error {
	result, err := r.db.Exec("DELETE FROM master_works WHERE id = $1 AND master_id = $2", workID, masterID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Master Payment Info Methods

// GetMasterPaymentInfo gets payment info for a master
func (r *Repository) GetMasterPaymentInfo(masterID int) (*models.MasterPaymentInfo, error) {
	var info models.MasterPaymentInfo
	var kaspiCard, freedomCard, halykCard sql.NullString

	err := r.db.QueryRow(`
		SELECT id, master_id, kaspi_card, freedom_card, halyk_card, created_at, updated_at
		FROM master_payment_info WHERE master_id = $1
	`, masterID).
		Scan(&info.ID, &info.MasterID, &kaspiCard, &freedomCard, &halykCard,
			&info.CreatedAt, &info.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if kaspiCard.Valid {
		info.KaspiCard = kaspiCard.String
	}
	if freedomCard.Valid {
		info.FreedomCard = freedomCard.String
	}
	if halykCard.Valid {
		info.HalykCard = halykCard.String
	}

	return &info, nil
}

// UpdateMasterPaymentInfo updates payment info for a master
func (r *Repository) UpdateMasterPaymentInfo(masterID int, kaspiCard, freedomCard, halykCard string) error {
	_, err := r.db.Exec(`
		INSERT INTO master_payment_info (master_id, kaspi_card, freedom_card, halyk_card, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (master_id) DO UPDATE SET
		kaspi_card = EXCLUDED.kaspi_card,
		freedom_card = EXCLUDED.freedom_card,
		halyk_card = EXCLUDED.halyk_card,
		updated_at = NOW()
	`, masterID, kaspiCard, freedomCard, halykCard)
	return err
}

// Reviews Methods

// GetMasterReviews gets all reviews for a master
func (r *Repository) GetMasterReviews(masterID int) ([]models.ReviewWithUser, error) {
	rows, err := r.db.Query(`
		SELECT r.id, r.master_id, r.user_id, r.rating, r.comment, r.created_at, u.name
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.master_id = $1 ORDER BY r.created_at DESC
	`, masterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewWithUser
	for rows.Next() {
		var review models.ReviewWithUser
		if err := rows.Scan(&review.ID, &review.MasterID, &review.UserID, &review.Rating,
			&review.Comment, &review.CreatedAt, &review.UserName); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// CreateReview creates a new review
func (r *Repository) CreateReview(masterID, userID, rating int, comment string) (*models.Review, error) {
	var review models.Review
	err := r.db.QueryRow(`
		INSERT INTO reviews (master_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, master_id, user_id, rating, comment, created_at
	`, masterID, userID, rating, comment).
		Scan(&review.ID, &review.MasterID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &review, nil
}
