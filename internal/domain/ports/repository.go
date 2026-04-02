package ports

import (
	"context"
	"time"

	"ds-backend/internal/domain/models"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.UserResponse, error)
	GetByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ConsultPassword(ctx context.Context, email string) (string, error)
	Create(ctx context.Context, user *models.User) (*models.UserResponse, error)
	Update(ctx context.Context, user *models.User) (*models.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

type PropertyRepository interface {
	GetAll(ctx context.Context) ([]models.PropertyCard, error)
	GetByID(ctx context.Context, id uint) (*models.PropertyResponse, error)
	Create(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	Update(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	Delete(ctx context.Context, id uint) error
}

type TokenRepository interface {
	SaveToken(ctx context.Context, token *models.RefreshToken) error
	DeleteToken(ctx context.Context, tokenID string) error
	GetTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	GetTokenByUserID(ctx context.Context, userID uint) (*models.RefreshToken, error)
	DeleteExpiredTokensByUserID(ctx context.Context, userID uint) error
}

type EmailRepository interface {
	SaveVerificationCode(ctx context.Context, toEmail string, code string, motivo string) error
	GetIDFromEmail(ctx context.Context, toEmail string) (int, error)
	VerifyTokenAndUser(ctx context.Context, code string, toEmail string) (int, bool, error)
	UpdateUserVerificationStatus(ctx context.Context, userID int) error
	UpdateTokenAsUsed(ctx context.Context, code string, userID int) error
	GetTokenVerificationStatus(ctx context.Context, userID int) (bool, string, error)
	GetLatestTokenInfo(ctx context.Context, userID int) (int, string, error)
	UpdateTokenResendInfo(ctx context.Context, userID int, oldToken string, newToken string) error
}

type ImageRepository interface {
	SaveImage(ctx context.Context, image *models.Image) (*models.Image, error)
	GetImageByID(ctx context.Context, id uint) (*models.Image, error)
	GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error)
	GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error)
	UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error
	DeleteImagesByPropertyID(ctx context.Context, propertyID uint) error
	UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error)
	DeleteImage(ctx context.Context, id uint) error
	HardDeleteImage(ctx context.Context, id uint) error
}

type AppointmentRepository interface {
    GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error)
    GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error)
    GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error)
    GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error)
    GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error)
	ClientHasOverlap(ctx context.Context, clientID uint, start, end time.Time, excludeID uint) (bool, error)
    Create(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error)
    Update(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error)
    Delete(ctx context.Context, id uint) error
}
