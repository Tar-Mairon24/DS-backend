package serviceMocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockEmailService struct {
	mock.Mock
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (m *MockEmailService) SendVerificationEmail(ctx context.Context, toEmail string, motivo string) error {
	args := m.Called(ctx, toEmail, motivo)
	return args.Error(0)
}

func (m *MockEmailService) VerifyEmail(ctx context.Context, verificacionData models.EmailVerification) (bool, error) {
	args := m.Called(ctx, verificacionData)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmailService) ResendVerificationEmail(ctx context.Context, toEmail string) error {
	args := m.Called(ctx, toEmail)
	return args.Error(0)
}
