package repositories

import (
	"github.com/bishalcode869/Auth-as-a-Service.git/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

type AuthRepositoryImpl struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		DB: db,
	}
}

func (repo *AuthRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *AuthRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Where("username = ? ", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
