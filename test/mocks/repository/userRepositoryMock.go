package repositoryMock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepo() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	args := m.Called(ctx, user)
	if userResponse, ok := args.Get(0).(*models.UserResponse); ok {
		return userResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) ConsultPassword(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) GetAll(ctx context.Context, userType string, search string) ([]models.UserResponse, error) {
	args := m.Called(ctx, userType, search)
	if users, ok := args.Get(0).([]models.UserResponse); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) Update(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	args := m.Called(ctx, user)
	if userResponse, ok := args.Get(0).(*models.UserResponse); ok {
		return userResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
