package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// You can use GIN_MODE=release in production
	router := gin.Default()

	// Simple health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Attach API routes
	registerRoutes(router)

	// Port from env or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
