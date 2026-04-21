package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
	"ds-backend/middleware"
)

type authUseCase struct {
	repo          	ports.UserRepository
	tokenRepo     	ports.TokenRepository
	jwtService     	ports.JWTService
	hashing       	middleware.HashingInterface
	authMiddleware 	middleware.AuthMiddlewareInterface
}

func NewAuthUseCase(repo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService, hashing middleware.HashingInterface, authMiddleware middleware.AuthMiddlewareInterface) ports.AuthUseCase {
	return &authUseCase{
		repo:           repo,
		tokenRepo:      tokenRepo,
		jwtService:     jwtService,
		hashing:        hashing,
		authMiddleware: authMiddleware,
	}
}

func (au *authUseCase) Login(ctx context.Context, email string, password string) (*models.LoginResponse, error) {
	if email == "" || password == "" {
		logrus.Error("Email and password cannot be empty")
		return nil, errors.New("email and password cannot be empty")
	}

	user, err := au.repo.GetByEmail(ctx, email)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user by email")
		return nil, errors.New("user not found")
	}

	if err := au.hashing.VerifyPassword(user.Password, password); err != nil {
		logrus.WithError(err).Error("Password verification failed")
		return nil, err
	}

	if err := au.tokenRepo.DeleteExpiredTokensByUserID(ctx, user.ID); err != nil {
		logrus.WithError(err).Warn("Failed to clean up expired tokens, proceeding anyway")
	}

	refreshToken, idToken, err := au.authMiddleware.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}

	err = au.tokenRepo.SaveToken(ctx, &models.RefreshToken{
		ID:        idToken,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	token, err := au.jwtService.GenerateToken(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		return nil, err
	}

	logrus.Infof("User %s login successful", user.Username)
	return &models.LoginResponse{
		User:         user.ToUserResponse(),
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (au *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		logrus.Error("Refresh token cannot be empty")
		return errors.New("refresh token cannot be empty")
	}

	token, err := au.tokenRepo.GetTokenByToken(ctx, refreshToken)
	if err != nil || token == nil {
		logrus.WithError(err).Warn("Token not found during logout")
		return errors.New("no active session found for this token")
	}

	if err := au.tokenRepo.DeleteToken(ctx, token.ID); err != nil {
		logrus.WithError(err).Error("Failed to delete token")
		return err
	}

	logrus.Infof("Session logged out successfully for user ID: %d", token.UserID)
	return nil
}

func (au *authUseCase) RefreshToken(ctx context.Context, data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	if data.JwtToken == "" || data.RefreshToken == "" {
		return nil, errors.New("JWT token and refresh token cannot be empty")
	}

	refreshToken, err := au.tokenRepo.GetTokenByToken(ctx, data.RefreshToken)
	if err != nil || refreshToken == nil {
		logrus.WithError(err).Error("Refresh token not found")
		return nil, errors.New("refresh token not found")
	}

	if refreshToken.ExpiresAt < time.Now().Unix() {
		logrus.Error("Refresh token expired")
		return nil, errors.New("refresh token expired")
	}

	newJwtToken, err := au.jwtService.RefreshToken(ctx, data.JwtToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh JWT token")
		return nil, err
	}

	return &models.RefreshTokenData{
		JwtToken:     newJwtToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (au *authUseCase) GetStatus(ctx context.Context, userID uint, refreshToken string) error {
	if userID == 0 || refreshToken == "" {
		return errors.New("user ID and refresh token cannot be empty")
	}

	savedToken, err := au.tokenRepo.GetTokenByToken(ctx, refreshToken)
	if err != nil || savedToken == nil {
		logrus.WithError(err).Error("Refresh token not found")
		return errors.New("refresh token not found")
	}

	if savedToken.UserID != userID {
		logrus.Error("Token does not belong to the requesting user")
		return errors.New("invalid refresh token")
	}

	if savedToken.ExpiresAt < time.Now().Unix() {
		logrus.Error("Refresh token expired")
		return errors.New("refresh token expired")
	}

	return nil
}