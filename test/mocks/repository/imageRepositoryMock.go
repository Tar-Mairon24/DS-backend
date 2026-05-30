package repositoryMock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockImageRepository struct {
	mock.Mock
}

func NewMockImageRepository() *MockImageRepository {
	return &MockImageRepository{}
}

func (m *MockImageRepository) SaveImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	args := m.Called(ctx, image)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageRepository) GetImageByID(ctx context.Context, id uint) (*models.Image, error) {
	args := m.Called(ctx, id)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageRepository) GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error) {
	args := m.Called(ctx, propertyID)
	if images, ok := args.Get(0).([]models.Image); ok {
		return images, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageRepository) GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error) {
	args := m.Called(ctx, propertyID)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageRepository) UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error {
	args := m.Called(ctx, propertyID, imageID)
	return args.Error(0)
}

func (m *MockImageRepository) DeleteImagesByPropertyID(ctx context.Context, propertyID uint) error {
	args := m.Called(ctx, propertyID)
	return args.Error(0)
}

func (m *MockImageRepository) UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	args := m.Called(ctx, image)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageRepository) DeleteImage(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockImageRepository) HardDeleteImage(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
