package repositoryMock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockEmailRepository struct {
	mock.Mock
}

func NewMockEmailRepository() *MockEmailRepository {
	return &MockEmailRepository{}
}

func (m *MockEmailRepository) SaveVerificationCode(ctx context.Context, toEmail string, code string, motivo string) error {
	args := m.Called(ctx, toEmail, code, motivo)
	return args.Error(0)
}

func (m *MockEmailRepository) GetIDFromEmail(ctx context.Context, toEmail string) (int, error) {
	args := m.Called(ctx, toEmail)
	return args.Int(0), args.Error(1)
}

func (m *MockEmailRepository) VerifyTokenAndUser(ctx context.Context, code string, toEmail string) (int, bool, error) {
	args := m.Called(ctx, code, toEmail)
	return args.Int(0), args.Bool(1), args.Error(2)
}

func (m *MockEmailRepository) UpdateTokenAsUsed(ctx context.Context, code string, userID int) error {
	args := m.Called(ctx, code, userID)
	return args.Error(0)
}

func (m *MockEmailRepository) GetTokenVerificationStatus(ctx context.Context, userID int) (bool, string, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockEmailRepository) GetLatestTokenInfo(ctx context.Context, userID int) (int, string, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.String(1), args.Error(2)
}

func (m *MockEmailRepository) UpdateTokenResendInfo(ctx context.Context, userID int, oldToken string, newToken string) error {
	args := m.Called(ctx, userID, oldToken, newToken)
	return args.Error(0)
}

func (m *MockEmailRepository) UpdateUserVerificationStatus(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
