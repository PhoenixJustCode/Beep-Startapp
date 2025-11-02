package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	PhotoURL     string    `json:"photo_url" db:"photo_url"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Category represents a service category
type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Service represents a service offered
type Service struct {
	ID              int       `json:"id" db:"id"`
	CategoryID      int       `json:"category_id" db:"category_id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	BasePrice       float64   `json:"base_price" db:"base_price"`
	MinPrice        float64   `json:"min_price" db:"min_price"`
	MaxPrice        float64   `json:"max_price" db:"max_price"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Car represents a car model
type Car struct {
	ID    int    `json:"id" db:"id"`
	Brand string `json:"brand" db:"brand"`
	Model string `json:"model" db:"model"`
	Year  int    `json:"year" db:"year"`
	Type  string `json:"type" db:"type"`
}

// PriceZone represents a price zone
type PriceZone struct {
	ID         int     `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	Multiplier float64 `json:"multiplier" db:"multiplier"`
}

// PricingRule represents a pricing rule
type PricingRule struct {
	ID            int     `json:"id" db:"id"`
	ServiceID     int     `json:"service_id" db:"service_id"`
	CarType       string  `json:"car_type" db:"car_type"`
	CarAgeMin     int     `json:"car_age_min" db:"car_age_min"`
	CarAgeMax     int     `json:"car_age_max" db:"car_age_max"`
	Multiplier    float64 `json:"multiplier" db:"multiplier"`
	FixedAddition float64 `json:"fixed_addition" db:"fixed_addition"`
}

// Master represents a service master
type Master struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	Phone          string    `json:"phone" db:"phone"`
	Specialization string    `json:"specialization" db:"specialization"`
	Rating         float64   `json:"rating" db:"rating"`
	PhotoURL       string    `json:"photo_url" db:"photo_url"`
	LocationLat    float64   `json:"location_lat" db:"location_lat"`
	LocationLng    float64   `json:"location_lng" db:"location_lng"`
	Address        string    `json:"address" db:"address"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// MasterSchedule represents a master's working schedule
type MasterSchedule struct {
	ID        int       `json:"id" db:"id"`
	MasterID  int       `json:"master_id" db:"master_id"`
	DayOfWeek int       `json:"day_of_week" db:"day_of_week"`
	StartTime string    `json:"start_time" db:"start_time"`
	EndTime   string    `json:"end_time" db:"end_time"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Appointment represents a booking
type Appointment struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	MasterID  int       `json:"master_id" db:"master_id"`
	ServiceID int       `json:"service_id" db:"service_id"`
	Date      time.Time `json:"date" db:"date"`
	Time      string    `json:"time" db:"time"`
	Status    string    `json:"status" db:"status"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AppointmentWithDetails represents an appointment with related details
type AppointmentWithDetails struct {
	Appointment
	ServiceName string `json:"service_name" db:"service_name"`
	MasterName  string `json:"master_name" db:"master_name"`
}

// Review represents a review of a master
type Review struct {
	ID        int       `json:"id" db:"id"`
	MasterID  int       `json:"master_id" db:"master_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CalculatePriceResponse represents a price calculation response
type CalculatePriceResponse struct {
	ServiceID    int           `json:"service_id"`
	ServiceName  string        `json:"service_name"`
	CarBrand     string        `json:"car_brand"`
	CarModel     string        `json:"car_model"`
	CarYear      int           `json:"car_year"`
	CarType      string        `json:"car_type"`
	CarAge       int           `json:"car_age"`
	BasePrice    float64       `json:"base_price"`
	FinalPrice   float64       `json:"final_price"`
	MinPrice     float64       `json:"min_price"`
	MaxPrice     float64       `json:"max_price"`
	PriceDetails []PriceDetail `json:"price_details"`
}

// PriceDetail represents a price calculation step
type PriceDetail struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Multiplier  float64 `json:"multiplier,omitempty"`
	IsAddition  bool    `json:"is_addition"`
}

// MasterWork represents a master's completed work
type MasterWork struct {
	ID           int       `json:"id" db:"id"`
	MasterID     int       `json:"master_id" db:"master_id"`
	Title        string    `json:"title" db:"title"`
	WorkDate     time.Time `json:"work_date" db:"work_date"`
	CustomerName string    `json:"customer_name" db:"customer_name"`
	Amount       float64   `json:"amount" db:"amount"`
	PhotoURLs    []string  `json:"photo_urls" db:"photo_urls"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// MasterPaymentInfo represents payment information for a master
type MasterPaymentInfo struct {
	ID          int       `json:"id" db:"id"`
	MasterID    int       `json:"master_id" db:"master_id"`
	KaspiCard   string    `json:"kaspi_card" db:"kaspi_card"`
	FreedomCard string    `json:"freedom_card" db:"freedom_card"`
	HalykCard   string    `json:"halyk_card" db:"halyk_card"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewWithUser represents a review with user information
type ReviewWithUser struct {
	Review
	UserName string `json:"user_name" db:"user_name"`
}

// Subscription represents a user subscription
type Subscription struct {
	ID             int        `json:"id" db:"id"`
	UserID         int        `json:"user_id" db:"user_id"`
	Plan           string     `json:"plan" db:"plan"`
	TrialStartDate *time.Time `json:"trial_start_date" db:"trial_start_date"`
	TrialEndDate   *time.Time `json:"trial_end_date" db:"trial_end_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// FavoriteMaster represents a favorite master for a user
type FavoriteMaster struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	MasterID  int       `json:"master_id" db:"master_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserCar represents a user's car
type UserCar struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Year      int       `json:"year" db:"year"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Guarantee represents a guarantee for a service
type Guarantee struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	AppointmentID int       `json:"appointment_id" db:"appointment_id"`
	ServiceName   string    `json:"service_name" db:"service_name"`
	MasterName    string    `json:"master_name" db:"master_name"`
	ServiceDate   time.Time `json:"service_date" db:"service_date"`
	ExpiryDate    time.Time `json:"expiry_date" db:"expiry_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// GuaranteeWithDetails represents a guarantee with appointment and car details
type GuaranteeWithDetails struct {
	Guarantee
	AppointmentDate   string `json:"appointment_date"`
	AppointmentTime   string `json:"appointment_time"`
	AppointmentStatus string `json:"appointment_status"`
	CarName           string `json:"car_name"`
	CarYear           int    `json:"car_year"`
}

// Notification represents a notification for a user
type Notification struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Type      string    `json:"type" db:"type"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	RelatedID int       `json:"related_id" db:"related_id"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MasterWithStatus represents a master with verified status
type MasterWithStatus struct {
	Master
	IsVerified  bool `json:"is_verified" db:"is_verified"`
	ReviewCount int  `json:"review_count" db:"review_count"`
	WorkCount   int  `json:"work_count" db:"work_count"`
}
