package main

import (
	"chat/internal/config"
	"chat/internal/db"
	"chat/internal/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println(cfg.DBName)

  	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()
	
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Method 2: Use middleware to inject db
	r.Use(func(c *gin.Context) {
		c.Set("db", database)
		c.Next()
	})

	// Define a simple GET endpoint
	r.GET("/health", healthCheck)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	handlers.RegisterRoutes(r) 
	r.Run()
}

func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
}