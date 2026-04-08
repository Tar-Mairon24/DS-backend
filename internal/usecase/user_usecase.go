package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
	"ds-backend/middleware"
)

type UserUseCase struct {
	repo        ports.UserRepository
	hashing     middleware.HashingInterface
}

func NewUserUseCase(repo ports.UserRepository, hashing middleware.HashingInterface) ports.UserUseCase {
	return &UserUseCase{
		repo:      repo,
		hashing:   hashing,
	}
}

func (uc *UserUseCase) GetAllUsers(ctx context.Context, userType string, search string) ([]models.UserResponse, error) {
    validTypes := map[string]bool{
        models.UserTypeClient: true,
        models.UserTypeAgent:  true,
        models.UserTypeAdmin:  true,
        models.UserTypeOwner:  true,
        models.UserTypeAll:    true,
    }
    if !validTypes[userType] {
        return nil, errors.New("invalid user type")
    }
    return uc.repo.GetAll(ctx, userType, search)
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user *models.User, createBy string) (*models.UserResponse, error) {
	if user.Username == "" {
		logrus.Error("Username cannot be empty")
		return nil, errors.New("username cannot be empty")
	}
	if createBy == models.CreateContextSelf && user.Email == "" {
		logrus.Error("Email cannot be empty")
		return nil, errors.New("email cannot be empty")
	}

	if createBy == models.CreateContextSelf && user.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if user.Password != "" {
		hashedPassword, err := uc.hashing.HashPassword(user.Password)
		if err != nil {
			logrus.WithError(err).Error("Failed to hash password")
			return nil, err
		}
		user.Password = hashedPassword
	}


	return uc.repo.Create(ctx, user)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	return uc.repo.Update(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
