package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ai-contract-review/internal/api"
	"ai-contract-review/internal/store"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	api.SetupRoutes(router, db)

	// Health check with service status
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Health: http://localhost:%s/health\n", port)
	fmt.Printf("Status: http://localhost:%s/status\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
