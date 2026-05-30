package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/infrastructure/service"
	repositoryMock "ds-backend/test/mocks/repository"
)

// TestVerifyEmail_Success tests successful email verification
func TestVerifyEmail_Success(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	userID := 1

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(userID, false, nil)
	mockEmailRepo.On("UpdateUserVerificationStatus", ctx, userID).Return(nil)
	mockEmailRepo.On("UpdateTokenAsUsed", ctx, verifyData.Code, userID).Return(nil)

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.NoError(t, err)
	assert.True(t, verified)
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_InvalidCode tests verification with invalid code
func TestVerifyEmail_InvalidCode(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "invalid",
		Email: "test@example.com",
	}

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(0, false, errors.New("invalid verification code or email"))

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "invalid verification code or email", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_InvalidEmail tests verification with invalid email
func TestVerifyEmail_InvalidEmail(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "invalid-email",
	}

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(0, false, errors.New("invalid verification code or email"))

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_CodeAlreadyUsed tests verification with already used code
func TestVerifyEmail_CodeAlreadyUsed(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	userID := 1

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(userID, true, nil)

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "verification code has already been used", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_EmptyCode tests verification with empty code
func TestVerifyEmail_EmptyCode(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "",
		Email: "test@example.com",
	}

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "verification code and email must be provided", err.Error())
	mockEmailRepo.AssertNotCalled(t, "VerifyTokenAndUser")
}

// TestVerifyEmail_EmptyEmail tests verification with empty email
func TestVerifyEmail_EmptyEmail(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "",
	}

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "verification code and email must be provided", err.Error())
	mockEmailRepo.AssertNotCalled(t, "VerifyTokenAndUser")
}

// TestVerifyEmail_RepositoryFails tests when repository fails during verification
func TestVerifyEmail_RepositoryFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(0, false, errors.New("database error"))

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_UpdateUserVerificationStatusFails tests when updating user verification status fails
func TestVerifyEmail_UpdateUserVerificationStatusFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	userID := 1

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(userID, false, nil)
	mockEmailRepo.On("UpdateUserVerificationStatus", ctx, userID).
		Return(errors.New("database error"))

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestVerifyEmail_UpdateTokenAsUsedFails tests when updating token as used fails
func TestVerifyEmail_UpdateTokenAsUsedFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	userID := 1

	mockEmailRepo.On("VerifyTokenAndUser", ctx, verifyData.Code, verifyData.Email).
		Return(userID, false, nil)
	mockEmailRepo.On("UpdateUserVerificationStatus", ctx, userID).Return(nil)
	mockEmailRepo.On("UpdateTokenAsUsed", ctx, verifyData.Code, userID).
		Return(errors.New("database error"))

	verified, err := emailService.VerifyEmail(ctx, verifyData)

	assert.Error(t, err)
	assert.False(t, verified)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestResendVerificationEmail_UserNotFound tests resend when user not found
func TestResendVerificationEmail_UserNotFound(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "nonexistent@example.com"

	mockEmailRepo.On("GetIDFromEmail", ctx, toEmail).
		Return(0, errors.New("user not found"))

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestResendVerificationEmail_MaxResendReached tests resend when max resends reached
func TestResendVerificationEmail_MaxResendReached(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "test@example.com"
	userID := 1
	resendCount := 3 // Max resends reached
	oldToken := "old123456"

	mockEmailRepo.On("GetIDFromEmail", ctx, toEmail).Return(userID, nil)
	mockEmailRepo.On("GetTokenVerificationStatus", ctx, userID).Return(false, "registration", nil)
	mockEmailRepo.On("GetLatestTokenInfo", ctx, userID).Return(resendCount, oldToken, nil)

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "maximum number of resends reached", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestResendVerificationEmail_EmptyEmail tests resend with empty email
func TestResendVerificationEmail_EmptyEmail(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := ""

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "email must be provided", err.Error())
	mockEmailRepo.AssertNotCalled(t, "GetIDFromEmail")
}

// TestResendVerificationEmail_RepositoryFails tests when repository fails during resend
func TestResendVerificationEmail_RepositoryFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "test@example.com"

	mockEmailRepo.On("GetIDFromEmail", ctx, toEmail).
		Return(0, errors.New("database error"))

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestResendVerificationEmail_GetTokenVerificationStatusFails tests when getting token verification status fails
func TestResendVerificationEmail_GetTokenVerificationStatusFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "test@example.com"
	userID := 1

	mockEmailRepo.On("GetIDFromEmail", ctx, toEmail).Return(userID, nil)
	mockEmailRepo.On("GetTokenVerificationStatus", ctx, userID).
		Return(false, "", errors.New("database error"))

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestResendVerificationEmail_GetLatestTokenInfoFails tests when getting latest token info fails
func TestResendVerificationEmail_GetLatestTokenInfoFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "test@example.com"
	userID := 1

	mockEmailRepo.On("GetIDFromEmail", ctx, toEmail).Return(userID, nil)
	mockEmailRepo.On("GetTokenVerificationStatus", ctx, userID).Return(false, "registration", nil)
	mockEmailRepo.On("GetLatestTokenInfo", ctx, userID).
		Return(0, "", errors.New("database error"))

	err := emailService.ResendVerificationEmail(ctx, toEmail)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}

// TestSendVerificationEmail_RepositoryFails tests when repository fails to save verification code
func TestSendVerificationEmail_RepositoryFails(t *testing.T) {
	mockEmailRepo := repositoryMock.NewMockEmailRepository()
	emailService := service.NewEmailService(mockEmailRepo)

	ctx := context.Background()
	toEmail := "test@example.com"
	motivo := "registration"

	mockEmailRepo.On("SaveVerificationCode", ctx, toEmail, mock.MatchedBy(func(code string) bool {
		return len(code) == 6
	}), motivo).Return(errors.New("database error"))

	err := emailService.SendVerificationEmail(ctx, toEmail, motivo)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockEmailRepo.AssertExpectations(t)
}
