package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	repositoryMock "ds-backend/test/mocks/repository"
)

func TestPropertyUseCase_GetAllProperties(t *testing.T) {
	t.Run("should return all properties successfully", func(t *testing.T) {
		// Arrange
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		expectedProperties := []models.PropertyCard{
			{ID: 1, Title: "Property 1", Price: 100000},
			{ID: 2, Title: "Property 2", Price: 200000},
		}

		mockPropertyRepo.On("GetAll", context.Background()).Return(expectedProperties, nil)

		// Act
		result, err := propertyUseCase.GetAllProperties(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperties, result)
		assert.Len(t, result, 2)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no properties exist", func(t *testing.T) {
		// Arrange
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		expectedProperties := []models.PropertyCard{}

		mockPropertyRepo.On("GetAll", context.Background()).Return(expectedProperties, nil)

		// Act
		result, err := propertyUseCase.GetAllProperties(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperties, result)
		assert.Len(t, result, 0)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("GetAll", context.Background()).Return(nil, errors.New("database error"))

		// Act
		result, err := propertyUseCase.GetAllProperties(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "database error")
		mockPropertyRepo.AssertExpectations(t)
	})
}

func TestPropertyUseCase_GetPropertyByID(t *testing.T) {
	t.Run("should return property successfully", func(t *testing.T) {
		// Arrange
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		expectedProperty := &models.PropertyResponse{
			ID:    1,
			Title: "Property 1",
			Price: 100000,
		}

		mockPropertyRepo.On("GetByID", context.Background(), uint(1)).Return(expectedProperty, nil)

		// Act
		result, err := propertyUseCase.GetPropertyByID(context.Background(), uint(1))

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperty, result)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("should return error when property not found", func(t *testing.T) {
		// Arrange
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("GetByID", context.Background(), uint(999)).Return(nil, errors.New("not found"))

		// Act
		result, err := propertyUseCase.GetPropertyByID(context.Background(), uint(999))

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockPropertyRepo.AssertExpectations(t)
	})
}

func TestPropertyUseCase_DeleteProperty(t *testing.T) {
	t.Run("should delete property successfully", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("Delete", context.Background(), uint(1)).Return(nil)
		mockImageRepo.On("DeleteImagesByPropertyID", context.Background(), uint(1)).Return(nil)

		err := propertyUseCase.DeleteProperty(context.Background(), uint(1))

		assert.NoError(t, err)
		mockPropertyRepo.AssertExpectations(t)
		mockImageRepo.AssertExpectations(t)
	})

	t.Run("should return error when property ID is invalid", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		err := propertyUseCase.DeleteProperty(context.Background(), 0)

		assert.Error(t, err)
		assert.EqualError(t, err, "property ID must be provided")
	})

	t.Run("should return error when repository delete fails", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("Delete", context.Background(), uint(1)).Return(errors.New("delete failed"))

		err := propertyUseCase.DeleteProperty(context.Background(), uint(1))

		assert.Error(t, err)
		assert.EqualError(t, err, "delete failed")
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("should return error when image deletion fails", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("Delete", context.Background(), uint(1)).Return(nil)
		mockImageRepo.On("DeleteImagesByPropertyID", context.Background(), uint(1)).Return(errors.New("image delete failed"))

		err := propertyUseCase.DeleteProperty(context.Background(), uint(1))

		assert.Error(t, err)
		assert.EqualError(t, err, "image delete failed")
		mockPropertyRepo.AssertExpectations(t)
		mockImageRepo.AssertExpectations(t)
	})
}

func TestPropertyUseCase_GetPropertyCardByID(t *testing.T) {
	t.Run("should return property card successfully", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		expectedCard := &models.PropertyCard{
			ID:    1,
			Title: "Property 1",
			Price: 100000,
		}

		mockPropertyRepo.On("GetPropertyCardByID", context.Background(), uint(1)).Return(expectedCard, nil)

		result, err := propertyUseCase.GetPropertyCardByID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.Equal(t, expectedCard, result)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("should return error when property card not found", func(t *testing.T) {
		mockPropertyRepo := repositoryMock.NewMockPropertyRepository()
		mockUserRepo := repositoryMock.NewMockUserRepo()
		mockImageRepo := repositoryMock.NewMockImageRepository()
		propertyUseCase := usecase.NewPropertyUseCase(mockPropertyRepo, mockUserRepo, mockImageRepo)

		mockPropertyRepo.On("GetPropertyCardByID", context.Background(), uint(999)).Return(nil, nil)

		result, err := propertyUseCase.GetPropertyCardByID(context.Background(), uint(999))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "property not found")
		mockPropertyRepo.AssertExpectations(t)
	})
}
