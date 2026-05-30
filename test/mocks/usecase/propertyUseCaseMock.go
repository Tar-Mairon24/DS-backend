package usecaseMocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockPropertyUseCase struct {
	mock.Mock
}

func NewMockPropertyUseCase() *MockPropertyUseCase {
	return &MockPropertyUseCase{}
}

func (m *MockPropertyUseCase) GetAllProperties(ctx context.Context) ([]models.PropertyCard, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PropertyCard), args.Error(1)
}
func (m *MockPropertyUseCase) GetPropertyByID(ctx context.Context, id uint) (*models.PropertyResponse, error) {
	args := m.Called(ctx, id)
	if prop, ok := args.Get(0).(*models.PropertyResponse); ok {
		return prop, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyUseCase) CreateProperty(ctx context.Context, p *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(ctx, p)
	if propResp, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyUseCase) UpdateProperty(ctx context.Context, p *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(ctx, p)
	if propResp, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyUseCase) DeleteProperty(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPropertyUseCase) GetPropertyCardByID(ctx context.Context, id uint) (*models.PropertyCard, error) {
	args := m.Called(ctx, id)
	if card, ok := args.Get(0).(*models.PropertyCard); ok {
		return card, args.Error(1)
	}
	return nil, args.Error(1)
}
