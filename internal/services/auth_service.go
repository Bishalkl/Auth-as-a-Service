package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bishalcode869/Auth-as-a-Service.git/internal/models"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/repositories"
	"github.com/bishalcode869/Auth-as-a-Service.git/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(username, password, email string) (*models.User, string, error)
	LoginUser(identifier, password string) (*models.User, string, error)
	SendVerificationCode(email string) error
	VerifyOtp(email, otp string) error
}

// RedisStore defines expected methods for Redis usage
type RedisStore interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
}

type AuthServiceImpl struct {
	AuthRepo        repositories.AuthRepository
	HashPassword    func(string) (string, error)
	ComparePassword func(string, string) bool
	GenerateToken   func(string, string) (string, error)
	SendEmail       func(string, string, string) error
	OtpGenerator    func(int) (string, error)
	ValidateEmail   func(string) bool
	RedisClient     RedisStore
}

func NewAuthService(authRepo repositories.AuthRepository, redis RedisStore) AuthService {
	return &AuthServiceImpl{
		AuthRepo:        authRepo,
		HashPassword:    utils.HashPassword,
		ComparePassword: utils.ComparePassword,
		GenerateToken:   utils.GenerateAccessToken,
		SendEmail:       utils.SendVerificationEmail,
		OtpGenerator:    utils.GenerateOTP,
		ValidateEmail:   utils.IsValidEmail,
		RedisClient:     redis,
	}
}

// userExists checks if a user exists by username or email
func (s *AuthServiceImpl) userExists(username, email string) error {
	if username != "" {
		user, err := s.AuthRepo.GetUserByUsername(username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking username: %v", err)
		}
		if user != nil {
			return errors.New("username already taken")
		}
	}

	if email != "" {
		if !s.ValidateEmail(email) {
			return errors.New("invalid email format")
		}
		foundEmail, err := s.AuthRepo.GetUserByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking email: %v", err)
		}
		if foundEmail != nil {
			return errors.New("email already taken")
		}
	}
	return nil
}

// SendVerificationCode generates and sends an OTP to email, storing it in Redis
func (s *AuthServiceImpl) SendVerificationCode(email string) error {
	// Validate email format
	if !s.ValidateEmail(email) {
		return fmt.Errorf("invalid email address: %s", email)
	}

	// Rate Limiting: Ensure not exceed OTP request limits
	key := fmt.Sprintf("otp_request_count:%s", email)
	otpCountStr, err := s.RedisClient.Get(context.Background(), key)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to retrieve OTP request count: %v", err)
	}

	// Convert OTP request count to an integer
	otpCount := 0
	if otpCountStr != "" {
		otpCount, err = strconv.Atoi(otpCountStr)
		if err != nil {
			return fmt.Errorf("failed to convert OTP request count to integer: %v", err)
		}
	}

	// Limit OTP request (e.g., 5 requests per hour)
	if otpCount >= 5 {
		return fmt.Errorf("too many OTP requests, please try again later")
	}

	code, err := s.OtpGenerator(6)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %v", err)
	}

	subject := "Verify Your Email Address"
	body := fmt.Sprintf("<p>Your verification code is:</p><h2>%s</h2>", code)

	if err := s.SendEmail(email, subject, body); err != nil {
		log.Printf("failed to send email to %s: %v", email, err)
		return fmt.Errorf("failed to send verification email")
	}

	// Store the code in Redis with expiry (e.g., 10 minutes)
	ctx := context.Background()
	keyOtp := fmt.Sprintf("Verify:%s", email)
	if err := s.RedisClient.Set(ctx, keyOtp, code, 10*time.Minute); err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %v", err)
	}

	// Increment OTP request count and set expiry (1 hour)
	if _, err := s.RedisClient.Incr(ctx, key); err != nil {
		return fmt.Errorf("failed to increment OTP request count: %v", err)
	}

	// Set expiry for OTP request count (1 hour)
	if _, err := s.RedisClient.Expire(ctx, key, 1*time.Hour); err != nil {
		return fmt.Errorf("failed to set OTP request count expiry: %v", err)
	}

	return nil
}

// VerifyOtp check the OTP against Redis and verifies the user's email
func (s *AuthServiceImpl) VerifyOtp(email, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf("Verify:%s", email)

	// Get the OTP from Redis
	storedOtp, err := s.RedisClient.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to retrieve OTP from Redis: %v", err)
	}

	// Check if the OTP matches
	if otp != storedOtp {
		return errors.New("invalid or expired OTP")
	}

	// Mark the user's email as verified
	user, err := s.AuthRepo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %v", err)
	}

	user.IsVerified = true
	if _, err := s.AuthRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// OTP is valid, remove the key from Redis
	if err := s.RedisClient.Delete(ctx, key); err != nil {
		log.Printf("failed to delete OTP from Redis: %v", err)
	}

	return nil
}

// RegisterUser registers a new user with hashed password and token generation
func (s *AuthServiceImpl) RegisterUser(username, password, email string) (*models.User, string, error) {
	if err := s.userExists(username, email); err != nil {
		return nil, "", err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   false,
	}

	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", err)
	}

	token, err := s.GenerateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return createdUser, token, nil
}

// LoginUser authenticates user and returns JWT token if valid
func (s *AuthServiceImpl) LoginUser(identifier, password string) (*models.User, string, error) {
	var user *models.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = s.AuthRepo.GetUserByEmail(identifier)
	} else {
		user, err = s.AuthRepo.GetUserByUsername(identifier)
	}

	if err != nil {
		return nil, "", errors.New("invalid username/email or password")
	}

	if !s.ComparePassword(user.PasswordHash, password) {
		return nil, "", errors.New("invalid username/email or password")
	}

	if !user.IsVerified {
		return nil, "", errors.New("please verify your email before logging in")
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, token, nil
}
