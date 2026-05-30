package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	middlewareMock "ds-backend/test/mocks/middleware"
	repositoryMock "ds-backend/test/mocks/repository"
)

func TestUserUseCase_GetAllUsers(t *testing.T) {
	t.Run("should return all users successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		users := []models.UserResponse{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
		}

		mockRepo.On("GetAll", context.Background(), models.UserTypeAll, "").Return(users, nil)

		result, err := uc.GetAllUsers(context.Background(), models.UserTypeAll, "")

		assert.NoError(t, err)
		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		mockRepo.On("GetAll", context.Background(), models.UserTypeAll, "").Return(nil, errors.New("repo error"))

		result, err := uc.GetAllUsers(context.Background(), models.UserTypeAll, "")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "repo error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	t.Run("should return user successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("GetByID", context.Background(), uint(1)).Return(user, nil)

		result, err := uc.GetUserByID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		mockRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("user not found"))

		result, err := uc.GetUserByID(context.Background(), uint(999))

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_CreateUser(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		userResp := &models.UserResponse{
			ID:       1,
			Username: user.Username,
			Email:    user.Email,
		}

		mockHashing.On("HashPassword", "password123").Return("hashedpassword", nil)
		mockRepo.On("Create", context.Background(), mock.MatchedBy(func(u *models.User) bool {
			return u.Username == user.Username && u.Email == user.Email
		})).Return(userResp, nil)

		result, err := uc.CreateUser(context.Background(), user, "self")

		assert.NoError(t, err)
		assert.Equal(t, userResp, result)
		mockRepo.AssertExpectations(t)
		mockHashing.AssertExpectations(t)
	})

	t.Run("should return error when hashing fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockHashing.On("HashPassword", "password123").Return("", errors.New("hashing error"))

		result, err := uc.CreateUser(context.Background(), user, "self")

		assert.Nil(t, result)
		assert.Error(t, err)
		mockHashing.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockHashing.On("HashPassword", "password123").Return("hashedpassword", nil)
		mockRepo.On("Create", context.Background(), mock.MatchedBy(func(u *models.User) bool {
			return u.Username == user.Username && u.Email == user.Email
		})).Return(nil, errors.New("repo error"))

		result, err := uc.CreateUser(context.Background(), user, "self")

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockHashing.AssertExpectations(t)
	})

	t.Run("should create user with agent context", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			Username: "newagent",
			Email:    "agent@example.com",
			Password: "password123",
		}

		userResp := &models.UserResponse{
			ID:       2,
			Username: user.Username,
			Email:    user.Email,
		}

		mockHashing.On("HashPassword", "password123").Return("hashedpassword", nil)
		mockRepo.On("Create", context.Background(), mock.MatchedBy(func(u *models.User) bool {
			return u.Username == user.Username && u.Email == user.Email
		})).Return(userResp, nil)

		result, err := uc.CreateUser(context.Background(), user, "agent")

		assert.NoError(t, err)
		assert.Equal(t, userResp, result)
		mockRepo.AssertExpectations(t)
		mockHashing.AssertExpectations(t)
	})
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	t.Run("should update user successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			ID:       1,
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		userResp := &models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}

		mockRepo.On("Update", context.Background(), user).Return(userResp, nil)

		result, err := uc.UpdateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, userResp, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		user := &models.User{
			ID:       1,
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		mockRepo.On("Update", context.Background(), user).Return(nil, errors.New("update error"))

		result, err := uc.UpdateUser(context.Background(), user)

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	t.Run("should delete user successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		mockRepo.On("Delete", context.Background(), uint(1)).Return(nil)

		err := uc.DeleteUser(context.Background(), uint(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		uc := usecase.NewUserUseCase(mockRepo, mockHashing)

		mockRepo.On("Delete", context.Background(), uint(1)).Return(errors.New("delete error"))

		err := uc.DeleteUser(context.Background(), uint(1))

		assert.Error(t, err)
		assert.EqualError(t, err, "delete error")
		mockRepo.AssertExpectations(t)
	})
}
