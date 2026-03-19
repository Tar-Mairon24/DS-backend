package ports

import (
	"context"

	"ds-backend/internal/domain/models"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error)
	CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type PropertyUseCase interface {
	GetAllProperties(ctx context.Context) ([]models.PropertyResponse, error)
	GetPropertyByID(ctx context.Context, id uint) (*models.PropertyResponse, error)
	CreateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	UpdateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	DeleteProperty(ctx context.Context, id uint) error
}

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (*models.LoginResponse, error)
	Logout(ctx context.Context, id uint) error
	RefreshToken(ctx context.Context, data models.RefreshTokenData) (*models.RefreshTokenData, error)
	GetStatus(ctx context.Context, userID uint, refreshToken string) error
}
