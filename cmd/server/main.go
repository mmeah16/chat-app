package main

import (
	"chat/internal/config"
	"chat/internal/db"
	"chat/internal/handlers"
	"chat/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		zap.S().Infow("failed to load config", "error", err)
	}

  	database, err := db.Connect(cfg)
	if err != nil {
		zap.S().Fatalw("failed to connect to database", "error", err)
	}
	defer database.Close()
	
	// Create a Gin router without default middleware (logger and recovery) since we use Zap instead
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	r.Use(func(c *gin.Context) {
		c.Set("db", database)
		c.Next()
	})

	// Define a simple GET endpoint
	r.GET("/health", healthCheck)

	// Start server on port 8080 (default)
	handlers.RegisterRoutes(r)

	r.Run()
}

func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
}