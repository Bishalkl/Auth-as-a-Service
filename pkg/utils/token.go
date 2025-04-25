package utils

import (
	"fmt"
	"time"

	"github.com/bishalcode869/Auth-as-a-Service.git/configs"
	"github.com/dgrijalva/jwt-go"
)

// Claims defines the structure for JWT payload
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// GenerateAccessToken generates a short-lived access token
func GenerateAccessToken(userID string, email string) (string, error) {
	expiration := time.Duration(configs.Config.AccessTokenExpireMinutes) * time.Minute
	return generateToken(userID, email, expiration)
}

// GenerateRefreshToken generates a longer-lived refresh token
func GenerateRefreshToken(userID string, email string) (string, error) {
	expiration := time.Duration(configs.Config.RefreshTokenExpireHours) * time.Hour
	return generateToken(userID, email, expiration)
}

// Internal function for DRY token generation
func generateToken(userID string, email string, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    "Auth-as-a-Service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(configs.Config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// ValidateToken verifies token validity and returns the claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
