package repositoryMock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockTokenRepository struct {
	mock.Mock
}

func NewMockTokenRepo() *MockTokenRepository {
	return &MockTokenRepository{}
}

func (m *MockTokenRepository) SaveToken(ctx context.Context, token *models.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockTokenRepository) GetTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if t, ok := args.Get(0).(*models.RefreshToken); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokenRepository) GetTokenByUserID(ctx context.Context, userID uint) (*models.RefreshToken, error) {
	args := m.Called(ctx, userID)
	if t, ok := args.Get(0).(*models.RefreshToken); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokenRepository) DeleteExpiredTokensByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
