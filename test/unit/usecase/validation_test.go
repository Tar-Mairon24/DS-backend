package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	middlewareMock "ds-backend/test/mocks/middleware"
	repositoryMock "ds-backend/test/mocks/repository"
)

// ============================================================================
// PROPERTY USECASE VALIDATION TESTS
// ============================================================================

func TestPropertyUseCase_CreateProperty_EmptyInputValidation(t *testing.T) {
	t.Run("should return error when property is nil", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		result, err := propertyUseCase.CreateProperty(context.Background(), nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property cannot be nil")
	})

	t.Run("should return error when address is empty", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		property := &models.Property{
			Address: "",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "address cannot be empty")
	})

	t.Run("should return error when agents list is empty", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "at least one agent must be assigned")
	})
}

func TestPropertyUseCase_CreateProperty_BoundaryValues(t *testing.T) {
	t.Run("should return error when price is zero", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		property := &models.Property{
			Address: "123 Main St",
			Price:   0,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "price must be greater than zero")
	})

	t.Run("should return error when price is negative", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		property := &models.Property{
			Address: "123 Main St",
			Price:   -50000,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "price must be greater than zero")
	})

	t.Run("should accept minimum valid price", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		owner := &models.UserResponse{ID: 1, Role: models.UserTypeOwner}
		agent := &models.UserResponse{ID: 2, Role: models.UserTypeAgent}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(owner, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(agent, nil)

		property := &models.Property{
			Address: "123 Main St",
			Price:   0.01,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		expectedResponse := &models.PropertyResponse{ID: 1, Price: 0.01}
		mockPropertyRepo.On("Create", context.Background(), mock.MatchedBy(func(p *models.Property) bool {
			return p.Address == "123 Main St" && p.Price == 0.01 && p.OwnerID == 1
		})).Return(expectedResponse, nil)

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestPropertyUseCase_CreateProperty_BusinessLogicValidation(t *testing.T) {
	t.Run("should return error when owner not found", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockUserRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("user not found"))

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 999,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "owner not found")
	})

	t.Run("should return error when owner is not an owner role", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		notOwner := &models.UserResponse{ID: 1, Role: models.UserTypeAgent}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(notOwner, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(&models.UserResponse{ID: 2, Role: models.UserTypeAgent}, nil)

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "user is not an owner")
	})

	t.Run("should return error when agent not found", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		owner := &models.UserResponse{ID: 1, Role: models.UserTypeOwner}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(owner, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("user not found"))

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{999},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "agent not found")
	})

	t.Run("should return error when agent has invalid role", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		owner := &models.UserResponse{ID: 1, Role: models.UserTypeOwner}
		invalidAgent := &models.UserResponse{ID: 2, Role: models.UserTypeClient}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(owner, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(invalidAgent, nil)

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "user is not an agent or admin")
	})

	t.Run("should accept admin role as valid agent", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		owner := &models.UserResponse{ID: 1, Role: models.UserTypeOwner}
		admin := &models.UserResponse{ID: 2, Role: models.UserTypeAdmin}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(owner, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(admin, nil)

		expectedResponse := &models.PropertyResponse{ID: 1}
		mockPropertyRepo.On("Create", context.Background(), mock.MatchedBy(func(p *models.Property) bool {
			return p.Address == "123 Main St" && p.Price == 100000 && p.OwnerID == 1
		})).Return(expectedResponse, nil)

		property := &models.Property{
			Address: "123 Main St",
			Price:   100000,
			OwnerID: 1,
			UserID:  []uint{2},
		}

		result, err := propertyUseCase.CreateProperty(context.Background(), property)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// ============================================================================
// USER USECASE VALIDATION TESTS
// ============================================================================

func TestUserUseCase_CreateUser_EmptyInputValidation(t *testing.T) {
	t.Run("should return error when username is empty", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "",
			Email:    "user@example.com",
			Password: "password123",
		}

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextSelf)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "username cannot be empty")
	})

	t.Run("should return error when email is empty for self-registration", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "testuser",
			Email:    "",
			Password: "password123",
		}

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextSelf)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "email cannot be empty")
	})

	t.Run("should return error when password is empty for self-registration", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "testuser",
			Email:    "user@example.com",
			Password: "",
		}

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextSelf)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "password cannot be empty")
	})

	t.Run("should allow empty email for admin-created users", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "testuser",
			Email:    "",
			Password: "",
		}

		expectedResponse := &models.UserResponse{ID: 1, Username: "testuser"}
		mockUserRepo.On("Create", context.Background(), user).Return(expectedResponse, nil)

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextAgent)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUserUseCase_CreateUser_PasswordHashing(t *testing.T) {
	t.Run("should hash password before storing", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "testuser",
			Email:    "user@example.com",
			Password: "password123",
		}

		hashedPassword := "hashed_password_123"
		mockHashing.On("HashPassword", "password123").Return(hashedPassword, nil)

		expectedUser := &models.User{
			Username: "testuser",
			Email:    "user@example.com",
			Password: hashedPassword,
		}
		expectedResponse := &models.UserResponse{ID: 1, Username: "testuser"}
		mockUserRepo.On("Create", context.Background(), expectedUser).Return(expectedResponse, nil)

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextSelf)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockHashing.AssertExpectations(t)
	})

	t.Run("should return error when password hashing fails", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		user := &models.User{
			Username: "testuser",
			Email:    "user@example.com",
			Password: "password123",
		}

		mockHashing.On("HashPassword", "password123").Return("", errors.New("hashing failed"))

		result, err := userUseCase.CreateUser(context.Background(), user, models.CreateContextSelf)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "hashing failed")
	})
}

func TestUserUseCase_GetAllUsers_InvalidUserType(t *testing.T) {
	t.Run("should return error for invalid user type", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		result, err := userUseCase.GetAllUsers(context.Background(), "invalid_type", "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid user type")
	})

	t.Run("should accept valid user types", func(t *testing.T) {
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockHashing := middlewareMock.NewMockHashing()
		userUseCase := usecase.NewUserUseCase(mockUserRepo, mockHashing)

		validTypes := []string{
			models.UserTypeClient,
			models.UserTypeAgent,
			models.UserTypeAdmin,
			models.UserTypeOwner,
			models.UserTypeAll,
		}

		for _, userType := range validTypes {
			expectedUsers := []models.UserResponse{{ID: 1, Username: "user1"}}
			mockUserRepo.On("GetAll", context.Background(), userType, "").Return(expectedUsers, nil)

			result, err := userUseCase.GetAllUsers(context.Background(), userType, "")

			assert.NoError(t, err)
			assert.NotNil(t, result)
		}
	})
}

// ============================================================================
// APPOINTMENT USECASE VALIDATION TESTS
// ============================================================================

func TestAppointmentUseCase_Create_EmptyInputValidation(t *testing.T) {
	t.Run("should return error when client not found", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		mockUserRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("user not found"))

		req := &models.AppointmentRequest{
			ClientID:   999,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "client not found")
	})

	t.Run("should return error when property not found", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.User{ID: 1, Role: models.UserTypeClient}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("property not found"))

		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 999,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property not found")
	})

	t.Run("should return error when agent not found", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.User{ID: 1, Role: models.UserTypeClient}
		property := &models.PropertyResponse{ID: 1}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(property, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("user not found"))

		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{999},
			Status:     "scheduled",
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "agent not found")
	})

	t.Run("should return error when agent has invalid role", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.UserResponse{ID: 1, Role: models.UserTypeClient}
		property := &models.PropertyResponse{ID: 1}
		invalidAgent := &models.UserResponse{ID: 2, Role: models.UserTypeClient}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(property, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(invalidAgent, nil)

		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "user must be an agent or admin to be assigned to an appointment")
	})
}

func TestAppointmentUseCase_Create_BoundaryValues(t *testing.T) {
	t.Run("should return error when scheduling appointment in the past", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.UserResponse{ID: 1, Role: models.UserTypeClient}
		property := &models.PropertyResponse{ID: 1}
		agent := &models.UserResponse{ID: 2, Role: models.UserTypeAgent}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(property, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(agent, nil)

		pastTime := time.Now().Add(-1 * time.Hour)
		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
			StartDate:  models.CustomTime{Time: pastTime},
			EndDate:    models.CustomTime{Time: pastTime.Add(1 * time.Hour)},
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "cannot schedule an appointment for a time that has already passed")
	})

	t.Run("should allow appointment scheduled for current time or future", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.UserResponse{ID: 1, Role: models.UserTypeClient}
		property := &models.PropertyResponse{ID: 1}
		agent := &models.UserResponse{ID: 2, Role: models.UserTypeAgent}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(property, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(agent, nil)

		futureTime := time.Now().Add(24 * time.Hour)
		mockAppointmentRepo.On("ClientHasOverlap", context.Background(), uint(1), futureTime, futureTime.Add(1*time.Hour), uint(0)).Return(false, nil)

		expectedAppointment := &models.Appointment{ID: 1}
		mockAppointmentRepo.On("Create", context.Background(), &models.Appointment{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
			StartDate:  futureTime,
			EndDate:    futureTime.Add(1 * time.Hour),
		}).Return(expectedAppointment, nil)

		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
			StartDate:  models.CustomTime{Time: futureTime},
			EndDate:    models.CustomTime{Time: futureTime.Add(1 * time.Hour)},
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestAppointmentUseCase_Create_ConflictDetection(t *testing.T) {
	t.Run("should return error when client has overlapping appointment", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		client := &models.UserResponse{ID: 1, Role: models.UserTypeClient}
		property := &models.PropertyResponse{ID: 1}
		agent := &models.UserResponse{ID: 2, Role: models.UserTypeAgent}
		mockUserRepo.On("GetByID", context.Background(), uint(1)).Return(client, nil)
		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(property, nil)
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(agent, nil)

		futureTime := time.Now().Add(24 * time.Hour)
		mockAppointmentRepo.On("ClientHasOverlap", context.Background(), uint(1), futureTime, futureTime.Add(1*time.Hour), uint(0)).Return(true, nil)

		req := &models.AppointmentRequest{
			ClientID:   1,
			PropertyID: 1,
			AgentIDs:   []uint{2},
			Status:     "scheduled",
			StartDate:  models.CustomTime{Time: futureTime},
			EndDate:    models.CustomTime{Time: futureTime.Add(1 * time.Hour)},
		}

		result, err := appointmentUseCase.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "client already has an appointment during this time")
	})
}

func TestAppointmentUseCase_UpdateStatus_Validation(t *testing.T) {
	t.Run("should return error when scheduling appointment with past date", func(t *testing.T) {
		mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		appointmentUseCase := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

		pastTime := time.Now().Add(-1 * time.Hour)
		appointment := &models.AppointmentDetail{
			ID:        1,
			StartDate: pastTime,
		}
		mockAppointmentRepo.On("GetByID", context.Background(), uint(1)).Return(appointment, nil)

		err := appointmentUseCase.UpdateStatus(context.Background(), uint(1), "scheduled")

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot schedule an appointment for a time that has already passed")
	})
}

// ============================================================================
// IMAGE USECASE VALIDATION TESTS
// ============================================================================

func TestImageUseCase_SaveImage_EmptyInputValidation(t *testing.T) {
	t.Run("should return error when image is nil", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		result, err := imageUseCase.SaveImage(context.Background(), nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "image cannot be nil")
	})

	t.Run("should save valid image successfully", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "image.jpg",
		}

		expectedImage := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "image.jpg",
		}
		mockImageRepo.On("SaveImage", context.Background(), image).Return(expectedImage, nil)

		result, err := imageUseCase.SaveImage(context.Background(), image)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedImage, result)
	})
}

func TestImageUseCase_DeleteImage_Validation(t *testing.T) {
	t.Run("should return error when image not found", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		mockImageRepo.On("GetImageByID", context.Background(), uint(999)).Return(nil, errors.New("image not found"))

		err := imageUseCase.DeleteImage(context.Background(), uint(999))

		assert.Error(t, err)
		assert.EqualError(t, err, "image not found")
	})

	t.Run("should return error when image is nil", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		mockImageRepo.On("GetImageByID", context.Background(), uint(1)).Return(nil, nil)

		err := imageUseCase.DeleteImage(context.Background(), uint(1))

		assert.Error(t, err)
		assert.EqualError(t, err, "image not found")
	})

	t.Run("should delete image successfully", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			MainImage:  false,
		}

		mockImageRepo.On("GetImageByID", context.Background(), uint(1)).Return(image, nil)
		mockImageRepo.On("DeleteImage", context.Background(), uint(1)).Return(nil)

		err := imageUseCase.DeleteImage(context.Background(), uint(1))

		assert.NoError(t, err)
		mockImageRepo.AssertExpectations(t)
	})

	t.Run("should update main image when deleting main image", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		mainImage := &models.Image{
			ID:         1,
			PropertyID: 1,
			MainImage:  true,
		}

		nextImage := &models.Image{
			ID:         2,
			PropertyID: 1,
			MainImage:  false,
		}

		mockImageRepo.On("GetImageByID", context.Background(), uint(1)).Return(mainImage, nil)
		mockImageRepo.On("DeleteImage", context.Background(), uint(1)).Return(nil)
		mockImageRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return([]models.Image{*nextImage}, nil)
		mockImageRepo.On("UpdateMainImageStatus", context.Background(), uint(1), uint(2)).Return(nil)

		err := imageUseCase.DeleteImage(context.Background(), uint(1))

		assert.NoError(t, err)
		mockImageRepo.AssertExpectations(t)
	})
}

func TestImageUseCase_GetImagesByPropertyID_Validation(t *testing.T) {
	t.Run("should return empty list when no images exist", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		mockImageRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return([]models.Image{}, nil)

		result, err := imageUseCase.GetImagesByPropertyID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("should return images for property", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		images := []models.Image{
			{ID: 1, PropertyID: 1, Path: "image1.jpg"},
			{ID: 2, PropertyID: 1, Path: "image2.jpg"},
		}

		mockImageRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return(images, nil)

		result, err := imageUseCase.GetImagesByPropertyID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, images, result)
	})
}

func TestImageUseCase_GetImageByID_Validation(t *testing.T) {
	t.Run("should return image when found", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "image.jpg",
		}

		mockImageRepo.On("GetImageByID", context.Background(), uint(1)).Return(image, nil)

		result, err := imageUseCase.GetImageByID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, image, result)
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		mockImageRepo := repositoryMock.NewMockImageRepository()
		imageUseCase := usecase.NewImageUseCase(mockImageRepo)

		mockImageRepo.On("GetImageByID", context.Background(), uint(999)).Return(nil, errors.New("image not found"))

		result, err := imageUseCase.GetImageByID(context.Background(), uint(999))

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// ============================================================================
// PROPERTY USECASE UPDATE VALIDATION TESTS
// ============================================================================

func TestPropertyUseCase_UpdateProperty_EmptyInputValidation(t *testing.T) {
	t.Run("should return error when property is nil", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		result, err := propertyUseCase.UpdateProperty(context.Background(), nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property cannot be nil")
	})

	t.Run("should return error when property ID is zero", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		property := &models.Property{
			ID:      0,
			Address: "123 Main St",
			Price:   100000,
		}

		result, err := propertyUseCase.UpdateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property ID must be provided")
	})

	t.Run("should return error when property not found", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("property not found"))

		property := &models.Property{
			ID:      999,
			Address: "123 Main St",
			Price:   100000,
		}

		result, err := propertyUseCase.UpdateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property not found")
	})
}

func TestPropertyUseCase_UpdateProperty_AgentValidation(t *testing.T) {
	t.Run("should return error when updating with invalid agent", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		existingProperty := &models.PropertyResponse{
			ID:      1,
			Address: "123 Main St",
			Price:   100000,
		}

		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(existingProperty, nil)

		invalidAgent := &models.UserResponse{ID: 999, Role: models.UserTypeClient}
		mockUserRepo.On("GetByID", context.Background(), uint(999)).Return(invalidAgent, nil)

		property := &models.Property{
			ID:     1,
			UserID: []uint{999},
		}

		result, err := propertyUseCase.UpdateProperty(context.Background(), property)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "user is not an agent or admin")
	})

	t.Run("should accept admin role when updating agents", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		existingProperty := &models.PropertyResponse{
			ID:      1,
			Address: "123 Main St",
			Price:   100000,
		}

		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(existingProperty, nil)

		admin := &models.UserResponse{ID: 2, Role: models.UserTypeAdmin}
		mockUserRepo.On("GetByID", context.Background(), uint(2)).Return(admin, nil)

		expectedResponse := &models.PropertyResponse{ID: 1}
		mockPropertyRepo.On("Update", context.Background(), mock.MatchedBy(func(p *models.Property) bool {
			return p.ID == 1 && p.Address == "123 Main St" && p.Price == 100000
		})).Return(expectedResponse, nil)

		property := &models.Property{
			ID:     1,
			UserID: []uint{2},
		}

		result, err := propertyUseCase.UpdateProperty(context.Background(), property)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestPropertyUseCase_DeleteProperty_BoundaryValues(t *testing.T) {
	t.Run("should return error when property ID is zero", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		err := propertyUseCase.DeleteProperty(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "property ID must be provided")
	})

	t.Run("should return error when property ID is negative", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		err := propertyUseCase.DeleteProperty(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "property ID must be provided")
	})
}
