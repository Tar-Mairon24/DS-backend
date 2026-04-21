package usecaseMocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockAuthUseCase struct {
	mock.Mock
}

func NewMockAuthUseCase() *MockAuthUseCase {
	return &MockAuthUseCase{}
}

func (m *MockAuthUseCase) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	args := m.Called(ctx, email, password)

	if LoginResp, ok := args.Get(0).(*models.LoginResponse); ok {
		return LoginResp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthUseCase) RefreshToken(ctx context.Context, data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	args := m.Called(ctx, data)
	if resp, ok := args.Get(0).(*models.RefreshTokenData); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthUseCase) GetStatus(ctx context.Context, userID uint, refreshToken string) error {
	args := m.Called(ctx, userID, refreshToken)
	return args.Error(0)
}
