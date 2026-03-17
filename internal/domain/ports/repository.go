package ports

import "ds-backend/internal/domain/models"

type UserRepository interface {
	GetAll() ([]models.UserResponse, error)
	GetByID(id uint) (*models.UserResponse, error)
	GetByEmail(email string) (*models.User, error)
	ConsultPassword(email string) (string, error)
	Create(user *models.User) (*models.UserResponse, error)
	Update(user *models.User) (*models.UserResponse, error)
	Delete(id uint) error
}

type PropertyRepository interface {
	GetAll() ([]models.PropertyResponse, error)
	GetByID(id uint) (*models.PropertyResponse, error)
	Create(property *models.Property) (*models.PropertyResponse, error)
	Update(property *models.Property) (*models.PropertyResponse, error)
	Delete(id uint) error
}

type TokenRepository interface {
	SaveToken(token *models.RefreshToken) error
	DeleteToken(tokenID string) error
	GetTokenIDByUserID(userID uint) (string, error)
	GetTokenByUserID(userID uint) (*models.RefreshToken, error)
}

type EmailRepository interface {
	SaveVerificationCode(toEmail string, code string, motivo string) error
	GetIDFromEmail(toEmail string) (int, error)
	VerifyTokenAndUser(code string, toEmail string) (int, bool, error)
	UpdateUserVerificationStatus(userID int) error
	UpdateTokenAsUsed(code string, userID int) error
	GetTokenVerificationStatus(userID int) (bool, string, error)
	GetLatestTokenInfo(userID int) (int, string, error)
	UpdateTokenResendInfo(userID int, oldToken string, newToken string) error
}