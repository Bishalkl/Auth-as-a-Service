package main

import (
	"log"
	"net/http"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application (config, DB, controller, etc.)
	_, err := bootstrap.InitalizeApp()
	if err != nil {
		log.Fatal("❌ App initialization failed:", err)
	}

	//Initialize Gin router
	router := gin.Default()

	// Health check or welcome route
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "🚀 Hello, Auth-as-a-Service is working!"})
	})

	// Get port from config (with fallback)
	port := configs.Config.Port
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running at http://localhost:8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}

}
