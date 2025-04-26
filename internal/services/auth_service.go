package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bishalcode869/Auth-as-a-Service.git/internal/models"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/repositories"
	"github.com/bishalcode869/Auth-as-a-Service.git/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(username, password, email string) (*models.User, string, error)
	LoginUser(identifier, password string) (*models.User, string, error)
}

type AuthServiceImpl struct {
	AuthRepo        repositories.AuthRepository
	HashPassword    func(string) (string, error)
	ComparePassword func(string, string) bool
	GenerateToken   func(string, string) (string, error)
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo:        authRepo,
		HashPassword:    utils.HashPassword,
		ComparePassword: utils.ComparePassword,
		GenerateToken:   utils.GenerateAccessToken,
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

// RegisterUser registers a new user
func (s *AuthServiceImpl) RegisterUser(username, password, email string) (*models.User, string, error) {
	// Ensure username/email are not already in use
	if err := s.userExists(username, email); err != nil {
		return nil, "", err
	}

	// Hash password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Create and persist user
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

	// Generate JWT token
	token, err := s.GenerateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return createdUser, token, nil
}

// LoginUser checks user credentials and returns a JWT token
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

	// Compare password
	if !s.ComparePassword(user.PasswordHash, password) {
		return nil, "", errors.New("invalid username/email or password")
	}

	// Optional: check if email is verified
	if !user.IsVerified {
		return nil, "", errors.New("please verify your email before logging in")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, token, nil
}
