package ports

import "ds-backend/internal/domain/models"

type UserUseCase interface {
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	CreateUser(user *models.User) (*models.UserResponse, error)
	UpdateUser(user *models.User) (*models.UserResponse, error)
	DeleteUser(id uint) error
}

type PropertyUseCase interface {
	GetAllProperties() ([]models.PropertyResponse, error)
	GetPropertyByID(id uint) (*models.PropertyResponse, error)
	CreateProperty(property *models.Property) (*models.PropertyResponse, error)
	UpdateProperty(property *models.Property) (*models.PropertyResponse, error)
	DeleteProperty(id uint) error
}

type AuthUseCase interface {
	Login(email string, password string) (*models.LoginResponse, error)
	Logout(id uint) error
	RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error)
	GetStatus(userID uint, refreshToken string) error
}