package router

import (
	"beep-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handlers.Handlers) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}

		// User profile
		user := v1.Group("/user")
		{
			user.GET("/profile", h.GetUserProfile)
			user.PUT("/profile", h.UpdateUserProfile)
			user.POST("/photo", h.UploadProfilePhoto)
		}

		// Master profile
		master := v1.Group("/master")
		{
			master.GET("/profile", h.GetMasterProfile)
			master.POST("/profile", h.CreateMasterProfile)
			master.PUT("/profile", h.UpdateMasterProfile)
			master.DELETE("/profile", h.DeleteMasterProfile)
			master.POST("/photo", h.UploadMasterPhoto)
			master.GET("/works", h.GetMasterWorks)
			master.POST("/works", h.CreateMasterWork)
			master.DELETE("/works/:id", h.DeleteMasterWork)
			master.GET("/payment-info", h.GetMasterPaymentInfo)
			master.PUT("/payment-info", h.UpdateMasterPaymentInfo)
			master.GET("/reviews", h.GetMasterReviews)
			master.POST("/reviews", h.CreateReview)
		}

		// Categories
		categories := v1.Group("/categories")
		{
			categories.GET("", h.GetCategories)
			categories.GET("/:id", h.GetCategoryByID)
		}

		// Services
		services := v1.Group("/services")
		{
			services.GET("", h.GetServices)
			services.GET("/:id", h.GetServiceByID)
		}

		// Cars
		cars := v1.Group("/cars")
		{
			cars.GET("", h.GetCars)
			cars.GET("/:id", h.GetCarByID)
		}

		// Pricing
		pricing := v1.Group("/pricing")
		{
			pricing.POST("/calculate", h.CalculatePrice)
		}

		// Masters
		masters := v1.Group("/masters")
		{
			masters.GET("", h.GetMasters)
			masters.GET("/:id", h.GetMasterByID)
			masters.GET("/:id/reviews", h.GetMasterReviews)
			masters.GET("/:id/schedule", h.GetMasterSchedule)
			masters.GET("/:id/available-slots", h.GetAvailableSlots)
		}

		// Reviews
		reviews := v1.Group("/reviews")
		{
			reviews.POST("", h.CreateReview)
		}

		// Appointments
		appointments := v1.Group("/appointments")
		{
			appointments.POST("", h.CreateAppointment)
			appointments.GET("", h.GetUserAppointments)
			appointments.GET("/:id", h.GetAppointmentByID)
			appointments.PUT("/:id", h.UpdateAppointment)
			appointments.DELETE("/:id", h.CancelAppointment)
		}
	}

	// Serve static files (after API routes)
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	r.GET("/profile", func(c *gin.Context) {
		c.File("./static/profile.html")
	})
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
