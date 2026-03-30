package repositoryMock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockPropertyRepository struct {
	mock.Mock
}

func NewMockPropertyRepository() *MockPropertyRepository {
	return &MockPropertyRepository{}
}

func (m *MockPropertyRepository) GetAll(ctx context.Context) ([]models.PropertyCard, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if properties, ok := args.Get(0).([]models.PropertyCard); ok {
		return properties, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) GetByID(ctx context.Context, id uint) (*models.PropertyResponse, error) {
	args := m.Called(ctx, id)
	if property, ok := args.Get(0).(*models.PropertyResponse); ok {
		return property, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Create(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(ctx, property)
	if propertyResponse, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propertyResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Update(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(ctx, property)
	if propertyResponse, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propertyResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
