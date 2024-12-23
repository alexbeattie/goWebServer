// main.go
package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
"time"
	"github.com/alexbeattie/golangone/config"
	"github.com/alexbeattie/golangone/handlers"
	"github.com/alexbeattie/golangone/models"
	"github.com/alexbeattie/golangone/services"
)

func initDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.UserPreferences{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &config.Config{
		OneStepGPSAPIKey: os.Getenv("ONESTEPGPS_API_KEY"),
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		DSN:              os.Getenv("DSN"),
	}

	db, err := initDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	service := services.NewService(db, cfg)
	handler := handlers.NewHandler(service, db)

	r := gin.Default()
	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
				"http://localhost:8080",
        "http://localhost:8081",  // This is already there
        "http://192.168.68.66:8081",  // Add this new IP address
        "http://192.168.1.82:8081",

		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{  "Origin",
        "Content-Type",
        "Content-Length",
        "Accept",
        "Authorization",
        "X-Requested-With",},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,

	}))
	api := r.Group("/api/v1")
{
    api.GET("/preferences/:userId", handler.GetUserPreferences)    // Changed from GetPreferences
    api.PUT("/preferences/:userId", handler.UpdateUserPreferences) // Changed from UpdatePreferences
    api.GET("/devices", handler.GetDevices)
}
	// Add this new v3 group
	v3 := r.Group("/v3/api")
	{
		    v3.GET("/device-info", handler.GetDeviceInfo)
    v3.GET("/route/drive-stop", handler.GetDriveStopRoute)

	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
