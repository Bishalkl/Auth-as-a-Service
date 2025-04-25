package services

import (
	"errors"
	"fmt"

	"github.com/bishalcode869/Auth-as-a-Service.git/internal/models"
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/repositories"
	"github.com/bishalcode869/Auth-as-a-Service.git/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	RegisterUser(username, password, email string) (*models.User, string, error)
	LoginUser(username, password, email string) (*models.User, string, error)
}

type AuthServiceImpl struct {
	AuthRepo        repositories.AuthRepository
	HashPassoword   func(string) (string, error)
	ComparePassword func(string, string) bool
	GenerateToken   func(string, string) (string, error)
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo:        authRepo,
		HashPassoword:   utils.HashPassword,
		ComparePassword: utils.ComparePassword,
		GenerateToken:   utils.GenerateAccessToken,
	}
}

// userExists checks if a user exists by username or email (internal helper)
func (s *AuthServiceImpl) userExists(username, email string) error {
	if username != "" {
		user, err := s.AuthRepo.GetUserByUsername(username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking username: %v", err)
		}
		if user != nil {
			return errors.New(("username already taken"))
		}
	}

	if email != "" {
		email, err := s.AuthRepo.GetUserByEmail(email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking email: %v", err)
		}
		if email != nil {
			return errors.New(("email already taken"))
		}
	}
	return nil
}

// RegisterUser registers a new user, hashes their password, and returns a JWT token
func (s *AuthServiceImpl) RegisterUser(username, password, email string) (*models.User, string, error) {
	// ensure username/email are not already in use
	if err := s.userExists(username, email); err != nil {
		return nil, "", err
	}

	// hash password
	hashedPassword, err := s.HashPassoword(password)
	if err != nil {
		return nil, " ", fmt.Errorf("failed to hashpassword: %v", err)
	}

	// create and persist user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	createdUser, err := s.AuthRepo.CreateUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %v", err)
	}

	// Generate JWT token
	token, err := s.GenerateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, "", fmt.Errorf("falied to generate token: %v", err)
	}
	return createdUser, token, nil

}

// LoginUser checks the user credentials, compares the password, and returns a token if valid
func (s *AuthServiceImpl) LoginUser(username, password, email string) (*models.User, string, error) {
	// Retrieve the user by email (or username)
	var user *models.User
	var err error

	switch {
	case username != "":
		user, err = s.AuthRepo.GetUserByUsername(username)
	case email != "":
		user, err = s.AuthRepo.GetUserByEmail(email)
	default:
		return nil, "", errors.New("username or email is required")
	}

	// Compare the password
	if !s.ComparePassword(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	// Optional: email verification
	if !user.IsVerified {
		return nil, "", errors.New("please verify your email before logging in")
	}

	// Generate a JWT token
	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, token, nil
}
