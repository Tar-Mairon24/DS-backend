package ports

import (
	"context"
	"time"

	"ds-backend/internal/domain/models"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error)
	CreateUser(ctx context.Context, user *models.User, createBy string) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type PropertyUseCase interface {
	GetAllProperties(ctx context.Context) ([]models.PropertyCard, error)
	GetPropertyByID(ctx context.Context, id uint) (*models.PropertyResponse, error)
	CreateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	UpdateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error)
	DeleteProperty(ctx context.Context, id uint) error
}

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (*models.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, data models.RefreshTokenData) (*models.RefreshTokenData, error)
	GetStatus(ctx context.Context, userID uint, refreshToken string) error
}

type ImageUseCase interface {
	SaveImage(ctx context.Context, image *models.Image) (*models.Image, error)
	GeneratePath(fileName string, propertyID uint) (diskPath string, urlPath string, err error)
	GetImageByID(ctx context.Context, id uint) (*models.Image, error)
	GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error)
	GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error)
	UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error
	UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error)
	DeleteImage(ctx context.Context, id uint) error
}

type AppointmentUseCase interface {
    GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error)
    GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error)
    GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error)
    GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error)
    GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error)
    Create(ctx context.Context, req *models.AppointmentRequest) (*models.Appointment, error) 
    Update(ctx context.Context, id uint, req *models.AppointmentRequest) (*models.Appointment, error)
    Delete(ctx context.Context, id uint) error
}
