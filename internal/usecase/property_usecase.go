package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type PropertyUseCase struct {
	propertyRepo ports.PropertyRepository
}

func NewPropertyUseCase(propertyRepo ports.PropertyRepository) ports.PropertyUseCase {
	return &PropertyUseCase{
		propertyRepo: propertyRepo,
	}
}

func (p *PropertyUseCase) GetAllProperties(ctx context.Context, ) ([]models.PropertyCard, error) {
	properties, err := p.propertyRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (p *PropertyUseCase) GetPropertyByID(ctx context.Context, id uint) (*models.PropertyResponse, error) {
	property, err := p.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, errors.New("property not found")
	}
	return property, nil
}

func (p *PropertyUseCase) CreateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	if property == nil {
		logrus.Error("Property cannot be nil")
		return nil, errors.New("property cannot be nil")
	}
	if property.Address == "" {
		logrus.Error("Address cannot be empty")
		return nil, errors.New("address cannot be empty")
	}
	if property.Price <= 0 {
		logrus.Error("Price must be greater than zero")
		return nil, errors.New("price must be greater than zero")
	}

	createdProperty, err := p.propertyRepo.Create(ctx, property)
	if err != nil {
		return nil, err
	}
	return createdProperty, nil
}

func (p *PropertyUseCase) UpdateProperty(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	if property == nil {
		logrus.Error("Property cannot be nil")
		return nil, errors.New("property cannot be nil")
	}
	if property.ID == 0 {
		logrus.Error("Property ID must be provided")
		return nil, errors.New("property ID must be provided")
	}
	if property.Address == "" {
		logrus.Error("Address cannot be empty")
		return nil, errors.New("address cannot be empty")
	}
	if property.Price <= 0 {
		logrus.Error("Price must be greater than zero")
		return nil, errors.New("price must be greater than zero")
	}

	updatedProperty, err := p.propertyRepo.Update(ctx, property)
	if err != nil {
		return nil, err
	}
	return updatedProperty, nil
}

func (p *PropertyUseCase) DeleteProperty(ctx context.Context, id uint) error {
	if id <= 0 {
		logrus.Error("Property ID must be provided")
		return errors.New("property ID must be provided")
	}

	err := p.propertyRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
