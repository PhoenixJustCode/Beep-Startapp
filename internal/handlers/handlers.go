package handlers

import (
	"beep-backend/internal/models"
	"beep-backend/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handlers {
	return &Handlers{repo: repo}
}

// getUserIDFromContext extracts user ID from JWT token in Authorization header
func (h *Handlers) getUserIDFromContext(c *gin.Context) (int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header")
	}

	// Remove "Bearer " prefix if present
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	// Extract email from mock token (format: "mock-jwt-token-email@example.com" or just "email@example.com")
	var email string
	if strings.Contains(token, "mock-jwt-token-") {
		parts := strings.Split(token, "mock-jwt-token-")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid token format: expected 'mock-jwt-token-email', got: %s", token)
		}
		email = parts[1]
	} else {
		// Try to use token directly as email (for backward compatibility)
		email = token
	}

	if email == "" {
		return 0, fmt.Errorf("email not found in token: %s", token)
	}

	// Get user by email
	user, err := h.repo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found with email: %s. Please register first", email)
		}
		return 0, fmt.Errorf("database error: %v", err)
	}

	return user.ID, nil
}

// Categories
func (h *Handlers) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prevent caching
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.JSON(http.StatusOK, categories)
}

func (h *Handlers) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	category, err := h.repo.GetCategoryByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// Services
func (h *Handlers) GetServices(c *gin.Context) {
	categoryID := c.Query("category_id")

	var services interface{}
	var err error

	if categoryID != "" {
		id, parseErr := strconv.Atoi(categoryID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}
		services, err = h.repo.GetServicesByCategory(id)
	} else {
		services, err = h.repo.GetAllServices()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *Handlers) GetServiceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	service, err := h.repo.GetServiceByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

// Cars
func (h *Handlers) GetCars(c *gin.Context) {
	cars, err := h.repo.GetAllCars()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (h *Handlers) GetCarByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	car, err := h.repo.GetCarByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, car)
}

// Pricing
func (h *Handlers) CalculatePrice(c *gin.Context) {
	type Request struct {
		ServiceID int `json:"service_id" binding:"required"`
		CarID     int `json:"car_id" binding:"required"`
		ZoneID    int `json:"zone_id"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.repo.CalculatePrice(req.ServiceID, req.CarID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service or car not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Masters
func (h *Handlers) GetMasters(c *gin.Context) {
	masters, err := h.repo.GetAllMasters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (optional - for checking favorites)
	userID, _ := h.getUserIDFromContext(c)

	// Enhance masters with verification status and favorite status
	type MasterResponse struct {
		*models.Master `json:",inline"`
		IsVerified     bool `json:"is_verified"`
		ReviewCount    int  `json:"review_count"`
		WorkCount      int  `json:"work_count"`
		IsFavorite     bool `json:"is_favorite"`
	}

	var result []MasterResponse
	for _, master := range masters {
		// Create a copy to avoid pointer issues
		masterCopy := master
		isVerified, reviewCount, workCount, _ := h.repo.CheckMasterVerificationStatus(master.ID)
		isFavorite := false
		if userID > 0 {
			isFavorite, _ = h.repo.IsFavoriteMaster(userID, master.ID)
		}

		result = append(result, MasterResponse{
			Master:      &masterCopy, // Use copy, not original
			IsVerified:  isVerified,
			ReviewCount: reviewCount,
			WorkCount:   workCount,
			IsFavorite:  isFavorite,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) GetMasterByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	master, err := h.repo.GetMasterByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Master not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, master)
}

func (h *Handlers) GetMasterSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	schedules, err := h.repo.GetMasterSchedule(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// Get available time slots for a master on a specific date
func (h *Handlers) GetAvailableSlots(c *gin.Context) {
	masterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	slots, err := h.repo.GetAvailableSlots(masterID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"date": dateStr, "slots": slots})
}

// Appointments
func (h *Handlers) CreateAppointment(c *gin.Context) {
	type Request struct {
		MasterID  int    `json:"master_id" binding:"required"`
		ServiceID int    `json:"service_id" binding:"required"`
		Date      string `json:"date" binding:"required"`
		Time      string `json:"time" binding:"required"`
		Comment   string `json:"comment"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appointment, err := h.repo.CreateAppointment(userID, req.MasterID, req.ServiceID, date, req.Time, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (h *Handlers) GetUserAppointments(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	appointments, err := h.repo.GetUserAppointments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func (h *Handlers) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	appointment, err := h.repo.GetAppointmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func (h *Handlers) UpdateAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	type Request struct {
		Comment string `json:"comment"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateAppointment(id, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully"})
}

func (h *Handlers) DeleteAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify that the appointment belongs to the current user
	appointment, err := h.repo.GetAppointmentByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if appointment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this appointment"})
		return
	}

	if err := h.repo.DeleteAppointment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}

func (h *Handlers) CancelAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.CancelAppointment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}

// Auth handlers
func (h *Handlers) Login(c *gin.Context) {
	type Request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual authentication with JWT
	// For MVP: Simple authentication
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For MVP: Simple password check (in production use bcrypt)
	// For now, just return user (without sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"token": "mock-jwt-token-" + user.Email,
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
		},
	})
}

func (h *Handlers) Register(c *gin.Context) {
	type Request struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual password hashing with bcrypt
	user, err := h.repo.CreateUser(req.Name, req.Email, req.Phone, "hashed-"+req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": "mock-jwt-token-" + user.Email,
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"phone":      user.Phone,
			"created_at": user.CreatedAt,
		},
	})
}

// GetUserProfile gets user profile
func (h *Handlers) GetUserProfile(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handlers) UpdateUserProfile(c *gin.Context) {
	type Request struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
		Phone *string `json:"phone"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get current user data
	currentUser, err := h.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Update only provided fields
	name := currentUser.Name
	email := currentUser.Email
	phone := currentUser.Phone

	if req.Name != nil {
		name = *req.Name
	}
	if req.Email != nil {
		email = *req.Email
	}
	if req.Phone != nil {
		phone = *req.Phone
	}

	// Update user profile
	if err := h.repo.UpdateUserProfile(userID, name, email, phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// Upload profile photo
func (h *Handlers) UploadProfilePhoto(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No photo uploaded"})
		return
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().Unix(), ext)
	filepath := filepath.Join("static", "uploads", filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update user photo URL in database
	photoURL := "/static/uploads/" + filename
	if err := h.repo.UpdateUserPhoto(userID, photoURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Photo uploaded successfully",
		"photo_url": photoURL,
	})
}

// Upload master photo
func (h *Handlers) UploadMasterPhoto(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get master profile
	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Master profile not found"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No photo uploaded"})
		return
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("master_%d_%d%s", master.ID, time.Now().Unix(), ext)
	filepath := filepath.Join("static", "uploads", filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update master photo URL in database
	photoURL := "/static/uploads/" + filename
	if err := h.repo.UpdateMasterPhoto(master.ID, photoURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Master photo uploaded successfully",
		"photo_url": photoURL,
	})
}

// Upload work photo
func (h *Handlers) UploadWorkPhoto(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		log.Printf("Unauthorized work photo upload attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("User %d uploading work photo", userID)

	// Get the uploaded file
	file, err := c.FormFile("photo")
	if err != nil {
		log.Printf("Error getting work photo file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No photo uploaded"})
		return
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("work_%d%s", time.Now().Unix(), ext)
	filepath := filepath.Join("static", "uploads", filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	photoURL := "/static/uploads/" + filename

	log.Printf("Work photo uploaded successfully: %s", photoURL)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Work photo uploaded successfully",
		"photo_url": photoURL,
	})
}

// Master Profile Handlers

// GetMasterProfile gets master profile by user ID
func (h *Handlers) GetMasterProfile(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Master profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, master)
}

// CreateMasterProfile creates a new master profile
func (h *Handlers) CreateMasterProfile(c *gin.Context) {
	type Request struct {
		Name           string `json:"name" binding:"required"`
		Email          string `json:"email" binding:"required"`
		Phone          string `json:"phone" binding:"required"`
		Specialization string `json:"specialization"`
		Address        string `json:"address"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Generate unique email for master profile to avoid conflicts
	masterEmail := fmt.Sprintf("master_%d_%s", userID, req.Email)

	master, err := h.repo.CreateMaster(userID, req.Name, masterEmail, req.Phone, req.Specialization, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, master)
}

// DeleteMasterProfile deletes master profile
func (h *Handlers) DeleteMasterProfile(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.DeleteMasterProfile(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Master profile deleted successfully"})
}

// UpdateMasterProfile updates master profile
func (h *Handlers) UpdateMasterProfile(c *gin.Context) {
	type Request struct {
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Phone          *string `json:"phone"`
		Specialization *string `json:"specialization"`
		Address        *string `json:"address"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		log.Printf("Unauthorized access to master profile update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get current master data
	currentMaster, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		log.Printf("Master not found for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	log.Printf("Updating master profile for user %d, master %d", userID, currentMaster.ID)

	// Update only provided fields
	name := currentMaster.Name
	email := currentMaster.Email
	phone := currentMaster.Phone
	specialization := currentMaster.Specialization
	address := currentMaster.Address

	if req.Name != nil {
		name = *req.Name
		log.Printf("Updating name to: %s", name)
	}
	if req.Email != nil {
		email = *req.Email
		log.Printf("Updating email to: %s", email)
	}
	if req.Phone != nil {
		phone = *req.Phone
		log.Printf("Updating phone to: %s", phone)
	}
	if req.Specialization != nil {
		specialization = *req.Specialization
		log.Printf("Updating specialization to: %s", specialization)
	}
	if req.Address != nil {
		address = *req.Address
		log.Printf("Updating address to: %s", address)
	}

	// Update master profile
	if err := h.repo.UpdateMaster(currentMaster.ID, name, email, phone, specialization, address); err != nil {
		log.Printf("Error updating master profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Master profile updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Master profile updated successfully"})
}

// Master Works Handlers

// GetMasterWorks gets all works for a master
func (h *Handlers) GetMasterWorks(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	works, err := h.repo.GetMasterWorks(master.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, works)
}

// CreateMasterWork creates a new work entry
func (h *Handlers) CreateMasterWork(c *gin.Context) {
	type Request struct {
		Title        string   `json:"title" binding:"required"`
		WorkDate     string   `json:"work_date" binding:"required"`
		CustomerName string   `json:"customer_name" binding:"required"`
		Amount       float64  `json:"amount" binding:"required"`
		PhotoURLs    []string `json:"photo_urls"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	work, err := h.repo.CreateMasterWork(master.ID, req.Title, workDate, req.CustomerName, req.Amount, req.PhotoURLs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, work)
}

// GetMasterWork gets a single work by ID
func (h *Handlers) GetMasterWork(c *gin.Context) {
	workID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid work ID: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		log.Printf("Unauthorized access to work %d", workID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		log.Printf("Master not found for user %d", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	log.Printf("Getting work %d for master %d (user %d)", workID, master.ID, userID)

	work, err := h.repo.GetMasterWork(workID, master.ID)
	if err != nil {
		log.Printf("Error getting work %d: %v", workID, err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully retrieved work %d", workID)
	c.JSON(http.StatusOK, work)
}

// UpdateMasterWork updates a work entry
func (h *Handlers) UpdateMasterWork(c *gin.Context) {
	workID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work ID"})
		return
	}

	type Request struct {
		Title        string   `json:"title" binding:"required"`
		WorkDate     string   `json:"work_date" binding:"required"`
		CustomerName string   `json:"customer_name" binding:"required"`
		Amount       float64  `json:"amount" binding:"required"`
		PhotoURLs    []string `json:"photo_urls"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	workDate, err := time.Parse("2006-01-02", req.WorkDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	if err := h.repo.UpdateMasterWork(workID, master.ID, req.Title, workDate, req.CustomerName, req.Amount, req.PhotoURLs); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Work not found or does not belong to this master"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work updated successfully"})
}

// DeleteMasterWork deletes a work entry
func (h *Handlers) DeleteMasterWork(c *gin.Context) {
	workID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify that the work belongs to the current master
	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	if err := h.repo.DeleteMasterWork(workID, master.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work deleted successfully"})
}

// Master Payment Info Handlers

// GetMasterPaymentInfo gets payment info for a master
func (h *Handlers) GetMasterPaymentInfo(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	info, err := h.repo.GetMasterPaymentInfo(master.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment info not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// UpdateMasterPaymentInfo updates payment info for a master
func (h *Handlers) UpdateMasterPaymentInfo(c *gin.Context) {
	type Request struct {
		KaspiCard   string `json:"kaspi_card"`
		FreedomCard string `json:"freedom_card"`
		HalykCard   string `json:"halyk_card"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	if err := h.repo.UpdateMasterPaymentInfo(master.ID, req.KaspiCard, req.FreedomCard, req.HalykCard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment info updated successfully"})
}

// Reviews Handlers

// GetMasterReviews gets all reviews for a master
func (h *Handlers) GetMasterReviews(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	reviews, err := h.repo.GetMasterReviews(master.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview creates a new review
func (h *Handlers) CreateReview(c *gin.Context) {
	type Request struct {
		MasterID int    `json:"master_id" binding:"required"`
		Rating   int    `json:"rating" binding:"required"`
		Comment  string `json:"comment"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	review, err := h.repo.CreateReview(req.MasterID, userID, req.Rating, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// New Feature Handlers

// GetMasterVerificationStatus gets verification status for a master
func (h *Handlers) GetMasterVerificationStatus(c *gin.Context) {
	masterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	isVerified, reviewCount, workCount, err := h.repo.CheckMasterVerificationStatus(masterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_verified":  isVerified,
		"review_count": reviewCount,
		"work_count":   workCount,
		"status":       "Проверенный мастер",
	})
}

// Subscription Handlers

// GetUserSubscription gets subscription for current user
func (h *Handlers) GetUserSubscription(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	subscription, err := h.repo.GetUserSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpdateUserSubscription updates subscription plan
func (h *Handlers) UpdateUserSubscription(c *gin.Context) {
	type Request struct {
		Plan string `json:"plan" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Plan != "basic" && req.Plan != "premium" && req.Plan != "trial" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan. Must be 'basic', 'premium' or 'trial'"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "user not found") || strings.Contains(errMsg, "Please register first") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Пользователь не найден. Пожалуйста, войдите в систему или зарегистрируйтесь.",
				"code":  "USER_NOT_FOUND",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован: " + errMsg})
		}
		return
	}

	// Check current subscription to allow upgrade from trial to premium
	currentSub, err := h.repo.GetUserSubscription(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о подписке: " + err.Error()})
		return
	}

	// Allow upgrade from trial to premium or any other transition
	if currentSub.Plan == "trial" && req.Plan == "premium" {
		// This is fine - allow upgrade
	} else if currentSub.Plan == "trial" && req.Plan == "basic" {
		// Downgrade from trial to basic is also fine
	}

	if err := h.repo.UpdateUserSubscription(userID, req.Plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить подписку: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// StartTrial starts a trial period for current user
func (h *Handlers) StartTrial(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	subscription, err := h.repo.StartTrial(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// Favorite Masters Handlers

// AddFavoriteMaster adds a master to favorites
func (h *Handlers) AddFavoriteMaster(c *gin.Context) {
	type Request struct {
		MasterID int `json:"master_id" binding:"required"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		// Provide helpful error message
		errMsg := err.Error()
		if strings.Contains(errMsg, "user not found") || strings.Contains(errMsg, "Please register first") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Пользователь не найден. Пожалуйста, войдите в систему или зарегистрируйтесь.",
				"code":    "USER_NOT_FOUND",
				"details": errMsg,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован: " + errMsg})
		}
		return
	}

	if req.MasterID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master_id"})
		return
	}

	if err := h.repo.AddFavoriteMaster(userID, req.MasterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Master added to favorites"})
}

// RemoveFavoriteMaster removes a master from favorites
func (h *Handlers) RemoveFavoriteMaster(c *gin.Context) {
	masterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		// Provide helpful error message
		errMsg := err.Error()
		if strings.Contains(errMsg, "user not found") || strings.Contains(errMsg, "Please register first") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Пользователь не найден. Пожалуйста, войдите в систему или зарегистрируйтесь.",
				"code":    "USER_NOT_FOUND",
				"details": errMsg,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован: " + errMsg})
		}
		return
	}

	if masterID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid master ID"})
		return
	}

	if err := h.repo.RemoveFavoriteMaster(userID, masterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Master removed from favorites"})
}

// GetFavoriteMasters gets all favorite masters for current user
func (h *Handlers) GetFavoriteMasters(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	masters, err := h.repo.GetFavoriteMasters(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, masters)
}

// User Cars Handlers

// GetUserCars gets all cars for current user
func (h *Handlers) GetUserCars(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cars, err := h.repo.GetUserCars(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// CreateUserCar creates a new car for current user
func (h *Handlers) CreateUserCar(c *gin.Context) {
	type Request struct {
		Name    string `json:"name" binding:"required"`
		Year    int    `json:"year"`
		Comment string `json:"comment"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	car, err := h.repo.CreateUserCar(userID, req.Name, req.Year, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, car)
}

// UpdateUserCar updates a car
func (h *Handlers) UpdateUserCar(c *gin.Context) {
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
		return
	}

	type Request struct {
		Name    string `json:"name" binding:"required"`
		Year    int    `json:"year"`
		Comment string `json:"comment"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.UpdateUserCar(carID, userID, req.Name, req.Year, req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car updated successfully"})
}

// DeleteUserCar deletes a car
func (h *Handlers) DeleteUserCar(c *gin.Context) {
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid car ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.DeleteUserCar(carID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}

// Guarantees Handlers

// GetUserGuarantees gets all active guarantees for current user
func (h *Handlers) GetUserGuarantees(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	guarantees, err := h.repo.GetUserGuarantees(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guarantees)
}

// Notifications Handlers

// GetUserNotifications gets all notifications for current user
func (h *Handlers) GetUserNotifications(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifications, err := h.repo.GetUserNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetMasterNotifications gets appointments for master (for master notifications)
func (h *Handlers) GetMasterNotifications(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	master, err := h.repo.GetMasterByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Master not found"})
		return
	}

	appointments, err := h.repo.GetMasterAppointmentsForNotifications(master.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

// MarkNotificationRead marks a notification as read
func (h *Handlers) MarkNotificationRead(c *gin.Context) {
	notificationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.MarkNotificationRead(notificationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
