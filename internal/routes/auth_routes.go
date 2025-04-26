package routes

import (
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Group all auth related routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}
