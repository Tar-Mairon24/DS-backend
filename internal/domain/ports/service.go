package ports

import (
	"context"

	"ds-backend/internal/domain/models"
)

type JWTService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*models.JWTClaims, error)
	RefreshToken(ctx context.Context, tokenString string) (string, error)
	GetUserIDFromClaims(ctx context.Context, tokenString string) (uint, error)
}

type EmailService interface {
	SendVerificationEmail(ctx context.Context, toEmail string, motivo string) error
	VerifyEmail(ctx context.Context, verificacionData models.EmailVerification) (bool, error)
	ResendVerificationEmail(ctx context.Context, toEmail string) error
}
