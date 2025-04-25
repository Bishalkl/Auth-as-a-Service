package handlers

import (
	"net/http"

	"github.com/bishalcode869/Auth-as-a-Service.git/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func NewAuthHanlder(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// Register
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	// Bind JSON data to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to register the user
	createdUser, token, err := h.AuthService.RegisterUser(input.Username, input.Password, input.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user":    createdUser,
	})
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON data to input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Authenticate the user using the service layer
	user, token, err := h.AuthService.LoginUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/email or"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"message":  "Login successful",
		"token":    token,
	})
}
