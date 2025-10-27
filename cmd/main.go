package main

import (
	"log"
	"os"

	"beep-backend/internal/config"
	"beep-backend/internal/database"
	"beep-backend/internal/handlers"
	"beep-backend/internal/repository"
	"beep-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	repos := repository.New(db)

	// Initialize handlers
	appHandlers := handlers.New(repos)

	// Setup router
	gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(appHandlers)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
