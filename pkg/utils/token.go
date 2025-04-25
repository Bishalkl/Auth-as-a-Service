package utils

import (
	"fmt"
	"time"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/dgrijalva/jwt-go"
)

// Custom claims for JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT access token for the user
func GenerateJWT(userID uint, email string) (string, error) {
	// Get expiration time from config (AccessTokenExpireMinutes)
	expirationTime := time.Duration(configs.Config.AccessTokenExpireMinutes) * time.Minute

	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
			Issuer:    "Auth-as-a-Service",
		},
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the JWT secret key
	signedToken, err := token.SignedString([]byte(configs.Config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// ValidateToken validates the JWT and returns the claims
func ValidateToken(tokenStr string) (*Claims, error) {
	// Parse and validate the token using the JWT secret
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Config.JWTSecret), nil
	})

	// Handle any errors in token parsing
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the claims are valid
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
