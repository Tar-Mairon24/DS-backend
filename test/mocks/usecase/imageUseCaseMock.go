package usecaseMocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockImageUseCase struct {
	mock.Mock
}

func NewMockImageUseCase() *MockImageUseCase {
	return &MockImageUseCase{}
}

func (m *MockImageUseCase) SaveImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	args := m.Called(ctx, image)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageUseCase) GeneratePath(fileName string, propertyID uint) (diskPath string, urlPath string, err error) {
	args := m.Called(fileName, propertyID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockImageUseCase) GetImageByID(ctx context.Context, id uint) (*models.Image, error) {
	args := m.Called(ctx, id)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageUseCase) GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error) {
	args := m.Called(ctx, propertyID)
	if images, ok := args.Get(0).([]models.Image); ok {
		return images, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageUseCase) GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error) {
	args := m.Called(ctx, propertyID)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageUseCase) UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error {
	args := m.Called(ctx, propertyID, imageID)
	return args.Error(0)
}

func (m *MockImageUseCase) UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	args := m.Called(ctx, image)
	if img, ok := args.Get(0).(*models.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockImageUseCase) DeleteImage(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
