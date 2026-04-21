package usecaseMocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockUserUseCase struct {
	mock.Mock
}

func NewMockUserUseCase() *MockUserUseCase {
	return &MockUserUseCase{}
}

func (m *MockUserUseCase) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockUserUseCase) GetUserByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	args := m.Called(ctx, user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) UpdateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	args := m.Called(ctx, user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) DeleteUser(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
