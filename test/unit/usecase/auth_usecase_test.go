package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	middlewareMock "ds-backend/test/mocks/middleware"
	repositoryMock "ds-backend/test/mocks/repository"
	serviceMocks "ds-backend/test/mocks/service"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Error("Could not load .env: ", err)
	}
	m.Run()
}

func TestAuthUseCase_Login(t *testing.T) {
	testUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "hashedpassword",
		Username: "testuser",
	}
	userResponse := models.UserResponse{
		ID:       testUser.ID,
		Email:    testUser.Email,
		Username: testUser.Username,
	}

	t.Run("empty email or password", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		resp, err := uc.Login(context.Background(), "", "password")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "email and password cannot be empty")

		resp, err = uc.Login(context.Background(), "email@example.com", "")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "email and password cannot be empty")
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), "notfound@example.com").Return(nil, errors.New("not found"))

		resp, err := uc.Login(context.Background(), "notfound@example.com", "password")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "user not found")

		userRepo.AssertExpectations(t)
	})

	t.Run("password verification failed", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "wrongpassword").Return(errors.New("invalid password"))

		resp, err := uc.Login(context.Background(), testUser.Email, "wrongpassword")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "invalid password")

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
	})

	t.Run("error cleaning expired tokens - but login continues", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
		tokenRepo.On("DeleteExpiredTokensByUserID", context.Background(), testUser.ID).Return(errors.New("db error"))

		// Login should continue despite the error
		authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
		tokenRepo.On("SaveToken", context.Background(), mock.AnythingOfType("*models.RefreshToken")).Return(nil)
		jwtService.On("GenerateToken", testUser).Return("jwtToken", nil)

		resp, err := uc.Login(context.Background(), testUser.Email, "password")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "jwtToken", resp.Token)
		assert.Equal(t, "refreshToken", resp.RefreshToken)

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
		authMiddleware.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("generate refresh token error", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
		tokenRepo.On("DeleteExpiredTokensByUserID", context.Background(), testUser.ID).Return(nil)
		authMiddleware.On("GenerateRefreshToken").Return("", "", errors.New("refresh error"))

		resp, err := uc.Login(context.Background(), testUser.Email, "password")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "refresh error")

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
		authMiddleware.AssertExpectations(t)
	})

	t.Run("save token error", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
		tokenRepo.On("DeleteExpiredTokensByUserID", context.Background(), testUser.ID).Return(nil)
		authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
		tokenRepo.On("SaveToken", context.Background(), mock.AnythingOfType("*models.RefreshToken")).Return(errors.New("save error"))

		resp, err := uc.Login(context.Background(), testUser.Email, "password")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "save error")

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
		authMiddleware.AssertExpectations(t)
	})

	t.Run("generate jwt token error", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
		tokenRepo.On("DeleteExpiredTokensByUserID", context.Background(), testUser.ID).Return(nil)
		authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
		tokenRepo.On("SaveToken", context.Background(), mock.AnythingOfType("*models.RefreshToken")).Return(nil)
		jwtService.On("GenerateToken", testUser).Return("", errors.New("jwt error"))

		resp, err := uc.Login(context.Background(), testUser.Email, "password")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "jwt error")

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
		authMiddleware.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		userRepo.On("GetByEmail", context.Background(), testUser.Email).Return(testUser, nil)
		hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
		tokenRepo.On("DeleteExpiredTokensByUserID", context.Background(), testUser.ID).Return(nil)
		authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
		tokenRepo.On("SaveToken", context.Background(), mock.AnythingOfType("*models.RefreshToken")).Return(nil)
		jwtService.On("GenerateToken", testUser).Return("jwtToken", nil)

		resp, err := uc.Login(context.Background(), testUser.Email, "password")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.User, "User should not be nil")

		assert.Equal(t, userResponse.ID, resp.User.ID)
		assert.Equal(t, userResponse.Email, resp.User.Email)
		assert.Equal(t, userResponse.Username, resp.User.Username)

		assert.Equal(t, "jwtToken", resp.Token)
		assert.Equal(t, "refreshToken", resp.RefreshToken)

		userRepo.AssertExpectations(t)
		hashing.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
		authMiddleware.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}

func TestAuthUseCase_Logout(t *testing.T) {
	t.Run("refreshToken is empty", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		err := uc.Logout(context.Background(), "")
		assert.EqualError(t, err, "refresh token cannot be empty")
	})

	t.Run("error getting token by token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), "token123").Return(nil, errors.New("db error"))

		err := uc.Logout(context.Background(), "token123")
		assert.EqualError(t, err, "no active session found for this token")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("token not found", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), "token123").Return(nil, nil)

		err := uc.Logout(context.Background(), "token123")
		assert.EqualError(t, err, "no active session found for this token")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("error deleting token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), "token123").Return(&models.RefreshToken{ID: "tokenid123", UserID: 3}, nil)
		tokenRepo.On("DeleteToken", context.Background(), "tokenid123").Return(errors.New("delete error"))

		err := uc.Logout(context.Background(), "token123")
		assert.EqualError(t, err, "delete error")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), "token456").Return(&models.RefreshToken{ID: "tokenid456", UserID: 4}, nil)
		tokenRepo.On("DeleteToken", context.Background(), "tokenid456").Return(nil)

		err := uc.Logout(context.Background(), "token456")
		assert.NoError(t, err)
		tokenRepo.AssertExpectations(t)
	})
}

func TestAuthUseCase_RefreshToken(t *testing.T) {
	testUserID := uint(1)
	validJwt := "valid.jwt.token"
	validRefresh := "valid-refresh-token"
	expiredRefresh := &models.RefreshToken{
		ID:        "tokenid",
		UserID:    testUserID,
		Token:     validRefresh,
		ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	validRefreshToken := &models.RefreshToken{
		ID:        "tokenid",
		UserID:    testUserID,
		Token:     validRefresh,
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	t.Run("empty jwt or refresh token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: "", RefreshToken: ""})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "JWT token and refresh token cannot be empty")
	})

	t.Run("error getting refresh token by token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(nil, errors.New("db error"))

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "refresh token not found")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("refresh token not found", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(nil, nil)

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "refresh token not found")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("refresh token expired", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(expiredRefresh, nil)

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "refresh token expired")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("error refreshing jwt token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(validRefreshToken, nil)
		jwtService.On("RefreshToken", context.Background(), validJwt).Return("", errors.New("jwt refresh error"))

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
		assert.Nil(t, resp)
		assert.EqualError(t, err, "jwt refresh error")
		tokenRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(validRefreshToken, nil)
		jwtService.On("RefreshToken", context.Background(), validJwt).Return("new.jwt.token", nil)

		resp, err := uc.RefreshToken(context.Background(), models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "new.jwt.token", resp.JwtToken)
		assert.Equal(t, validRefresh, resp.RefreshToken)
		tokenRepo.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}

func TestAuthUseCase_GetStatus(t *testing.T) {
	testUserID := uint(1)
	validRefresh := "valid-refresh-token"
	expiredRefresh := &models.RefreshToken{
		ID:        "tokenid",
		UserID:    testUserID,
		Token:     validRefresh,
		ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	validRefreshToken := &models.RefreshToken{
		ID:        "tokenid",
		UserID:    testUserID,
		Token:     validRefresh,
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	t.Run("empty userID or refreshToken", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		err := uc.GetStatus(context.Background(), 0, "")
		assert.EqualError(t, err, "user ID and refresh token cannot be empty")
		err = uc.GetStatus(context.Background(), testUserID, "")
		assert.EqualError(t, err, "user ID and refresh token cannot be empty")
		err = uc.GetStatus(context.Background(), 0, validRefresh)
		assert.EqualError(t, err, "user ID and refresh token cannot be empty")
	})

	t.Run("error getting refresh token by token", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(nil, errors.New("db error"))

		err := uc.GetStatus(context.Background(), testUserID, validRefresh)
		assert.EqualError(t, err, "refresh token not found")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("refresh token not found", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(nil, nil)

		err := uc.GetStatus(context.Background(), testUserID, validRefresh)
		assert.EqualError(t, err, "refresh token not found")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("refresh token expired", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(expiredRefresh, nil)

		err := uc.GetStatus(context.Background(), testUserID, validRefresh)
		assert.EqualError(t, err, "refresh token expired")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("invalid refresh token - user mismatch", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		// Token belongs to different user
		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(&models.RefreshToken{
			ID:        "tokenid",
			UserID:    uint(999),
			Token:     validRefresh,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			CreatedAt: time.Now().Add(-24 * time.Hour),
		}, nil)

		err := uc.GetStatus(context.Background(), testUserID, validRefresh)
		assert.EqualError(t, err, "invalid refresh token")
		tokenRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		userRepo := repositoryMock.NewMockUserRepo()
		tokenRepo := repositoryMock.NewMockTokenRepo()
		jwtService := serviceMocks.NewMockJWTService()
		hashing := middlewareMock.NewMockHashing()
		authMiddleware := middlewareMock.NewMockMiddleware()

		uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

		tokenRepo.On("GetTokenByToken", context.Background(), validRefresh).Return(validRefreshToken, nil)

		err := uc.GetStatus(context.Background(), testUserID, validRefresh)
		assert.NoError(t, err)
		tokenRepo.AssertExpectations(t)
	})
}
