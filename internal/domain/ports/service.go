package ports

import "ds-backend/internal/domain/models"

type JWTService interface {
    GenerateToken(user *models.User) (string, error)
    ValidateToken(tokenString string) (*models.JWTClaims, error)
    RefreshToken(tokenString string) (string, error)
    GetUserIDFromClaims(tokenString string) (uint, error)
}

type EmailService interface {
	SendVerificationEmail(toEmail string, motivo string) error
	VerifyEmail(verificacionData models.EmailVerification) (bool, error)
	ResendVerificationEmail(toEmail string) error
}