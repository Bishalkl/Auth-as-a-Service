package repositories

import (
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	UserExists(email string) (bool, error)
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
}

// struct
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// Constructor to initialize UserRepositoryImpl
func NewUserRepository(DB *gorm.DB) AuthRepository {
	return &UserRepositoryImpl{DB: DB}
}
