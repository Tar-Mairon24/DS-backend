package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	repositoryMock "ds-backend/test/mocks/repository"
)

func TestImageUseCase_SaveImage(t *testing.T) {
	t.Run("should save image successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		image := &models.Image{
			PropertyID:  1,
			Path:        "/uploads/properties/1/image.jpg",
			Description: "Living room",
			MainImage:   true,
			CreatedAt:   time.Now(),
		}

		savedImage := &models.Image{
			ID:          1,
			PropertyID:  1,
			Path:        "/uploads/properties/1/image.jpg",
			Description: "Living room",
			MainImage:   true,
			CreatedAt:   image.CreatedAt,
		}

		mockRepo.On("SaveImage", context.Background(), image).Return(savedImage, nil)

		result, err := uc.SaveImage(context.Background(), image)

		assert.NoError(t, err)
		assert.Equal(t, savedImage, result)
		assert.Equal(t, uint(1), result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when image is nil", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		result, err := uc.SaveImage(context.Background(), nil)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "image cannot be nil", err.Error())
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		image := &models.Image{
			PropertyID:  1,
			Path:        "/uploads/properties/1/image.jpg",
			Description: "Living room",
			MainImage:   true,
		}

		mockRepo.On("SaveImage", context.Background(), image).Return(nil, errors.New("database error"))

		result, err := uc.SaveImage(context.Background(), image)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestImageUseCase_DeleteImage(t *testing.T) {
	t.Run("should delete image successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/image.jpg",
			MainImage:  false,
		}

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(image, nil)
		mockRepo.On("DeleteImage", context.Background(), uint(1)).Return(nil)

		err := uc.DeleteImage(context.Background(), uint(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mockRepo.On("GetImageByID", context.Background(), uint(999)).Return(nil, errors.New("image not found"))

		err := uc.DeleteImage(context.Background(), uint(999))

		assert.Error(t, err)
		assert.Equal(t, "image not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when image is nil", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(nil, nil)

		err := uc.DeleteImage(context.Background(), uint(1))

		assert.Error(t, err)
		assert.Equal(t, "image not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository delete fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/image.jpg",
			MainImage:  false,
		}

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(image, nil)
		mockRepo.On("DeleteImage", context.Background(), uint(1)).Return(errors.New("delete failed"))

		err := uc.DeleteImage(context.Background(), uint(1))

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should update main image status when deleting main image", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mainImage := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/main.jpg",
			MainImage:  true,
		}

		remainingImages := []models.Image{
			{
				ID:         2,
				PropertyID: 1,
				Path:       "/uploads/properties/1/image2.jpg",
				MainImage:  false,
			},
		}

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(mainImage, nil)
		mockRepo.On("DeleteImage", context.Background(), uint(1)).Return(nil)
		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return(remainingImages, nil)
		mockRepo.On("UpdateMainImageStatus", context.Background(), uint(1), uint(2)).Return(nil)

		err := uc.DeleteImage(context.Background(), uint(1))

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when updating main image status fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mainImage := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/main.jpg",
			MainImage:  true,
		}

		remainingImages := []models.Image{
			{
				ID:         2,
				PropertyID: 1,
				Path:       "/uploads/properties/1/image2.jpg",
				MainImage:  false,
			},
		}

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(mainImage, nil)
		mockRepo.On("DeleteImage", context.Background(), uint(1)).Return(nil)
		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return(remainingImages, nil)
		mockRepo.On("UpdateMainImageStatus", context.Background(), uint(1), uint(2)).Return(errors.New("update failed"))

		err := uc.DeleteImage(context.Background(), uint(1))

		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestImageUseCase_GetImagesByPropertyID(t *testing.T) {
	t.Run("should return images successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		images := []models.Image{
			{
				ID:         1,
				PropertyID: 1,
				Path:       "/uploads/properties/1/image1.jpg",
				MainImage:  true,
			},
			{
				ID:         2,
				PropertyID: 1,
				Path:       "/uploads/properties/1/image2.jpg",
				MainImage:  false,
			},
		}

		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return(images, nil)

		result, err := uc.GetImagesByPropertyID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.Equal(t, images, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no images found", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		emptyImages := []models.Image{}

		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(999)).Return(emptyImages, nil)

		result, err := uc.GetImagesByPropertyID(context.Background(), uint(999))

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(1)).Return(nil, errors.New("database error"))

		result, err := uc.GetImagesByPropertyID(context.Background(), uint(1))

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return multiple images for property", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		images := []models.Image{
			{
				ID:         1,
				PropertyID: 5,
				Path:       "/uploads/properties/5/image1.jpg",
				MainImage:  true,
			},
			{
				ID:         2,
				PropertyID: 5,
				Path:       "/uploads/properties/5/image2.jpg",
				MainImage:  false,
			},
			{
				ID:         3,
				PropertyID: 5,
				Path:       "/uploads/properties/5/image3.jpg",
				MainImage:  false,
			},
		}

		mockRepo.On("GetImagesByPropertyID", context.Background(), uint(5)).Return(images, nil)

		result, err := uc.GetImagesByPropertyID(context.Background(), uint(5))

		assert.NoError(t, err)
		assert.Equal(t, images, result)
		assert.Len(t, result, 3)
		assert.True(t, result[0].MainImage)
		assert.False(t, result[1].MainImage)
		mockRepo.AssertExpectations(t)
	})
}

func TestImageUseCase_GetImageByID(t *testing.T) {
	t.Run("should get image by id successfully", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		image := &models.Image{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/image.jpg",
			MainImage:  true,
		}

		mockRepo.On("GetImageByID", context.Background(), uint(1)).Return(image, nil)

		result, err := uc.GetImageByID(context.Background(), uint(1))

		assert.NoError(t, err)
		assert.Equal(t, image, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		mockRepo := repositoryMock.NewMockImageRepository()
		uc := usecase.NewImageUseCase(mockRepo)

		mockRepo.On("GetImageByID", context.Background(), uint(999)).Return(nil, errors.New("image not found"))

		result, err := uc.GetImageByID(context.Background(), uint(999))

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
