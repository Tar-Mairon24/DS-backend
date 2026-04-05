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
	imageRepo    ports.ImageRepository
}

func NewPropertyUseCase(propertyRepo ports.PropertyRepository, imageRepo ports.ImageRepository) ports.PropertyUseCase {
	return &PropertyUseCase{
		propertyRepo: propertyRepo,
		imageRepo:    imageRepo,
	}
}

func (p *PropertyUseCase) GetAllProperties(ctx context.Context) ([]models.PropertyCard, error) {
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

func (p *PropertyUseCase) GetPropertyCardByID(ctx context.Context, id uint) (*models.PropertyCard, error) {
	propertyCard, err := p.propertyRepo.GetPropertyCardByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if propertyCard == nil {
		return nil, errors.New("property not found")
	}
	return propertyCard, nil
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

	// Fetch the existing property to preserve unmodified fields (partial update support)
	existingResponse, err := p.propertyRepo.GetByID(ctx, property.ID)
	if err != nil {
		logrus.WithError(err).Warnf("Property with ID %d not found", property.ID)
		return nil, errors.New("property not found")
	}

	// Convert PropertyResponse back to Property struct to get all fields
	existing := &models.Property{
		ID:              existingResponse.ID,
		Title:           existingResponse.Title,
		Address:         existingResponse.Address,
		Neighborhood:    existingResponse.Neighborhood,
		City:            existingResponse.City,
		Zone:            existingResponse.Zone,
		Reference:       existingResponse.Reference,
		Price:           existingResponse.Price,
		ConstructionM2:  existingResponse.ConstructionM2,
		LandM2:          existingResponse.LandM2,
		IsOccupied:      existingResponse.IsOccupied,
		IsFurnished:     existingResponse.IsFurnished,
		Floors:          existingResponse.Floors,
		Bedrooms:        existingResponse.Bedrooms,
		Bathrooms:       existingResponse.Bathrooms,
		GarageSize:      existingResponse.GarageSize,
		GardenM2:        existingResponse.GardenM2,
		GasTypes:        existingResponse.GasTypes,
		Amenities:       existingResponse.Amenities,
		Extras:          existingResponse.Extras,
		Utilities:       existingResponse.Utilities,
		Notes:           existingResponse.Notes,
		Description:     existingResponse.Description, // Map PublicNotes to Description
		PropertyType:    existingResponse.PropertyType,
		TransactionType: existingResponse.TransactionType,
		Status:          existingResponse.Status,
		CreatedAt:       existingResponse.CreatedAt,
		UpdatedAt:       existingResponse.UpdatedAt,
		OwnerID:         existingResponse.OwnerID,
		UserID:          property.UserID, // Use the current user from request
	}

	// Merge: Only override fields that were explicitly provided (non-empty/non-zero values)
	if property.Title != "" {
		existing.Title = property.Title
	}
	if property.Address != "" {
		existing.Address = property.Address
	}
	if property.Neighborhood != "" {
		existing.Neighborhood = property.Neighborhood
	}
	if property.City != "" {
		existing.City = property.City
	}
	if property.Zone != "" {
		existing.Zone = property.Zone
	}
	if property.Reference != "" {
		existing.Reference = property.Reference
	}
	if property.Price > 0 {
		existing.Price = property.Price
	}
	if property.ConstructionM2 > 0 {
		existing.ConstructionM2 = property.ConstructionM2
	}
	if property.LandM2 > 0 {
		existing.LandM2 = property.LandM2
	}
	if property.Floors > 0 {
		existing.Floors = property.Floors
	}
	if property.Bedrooms >= 0 {
		existing.Bedrooms = property.Bedrooms
	}
	if property.Bathrooms >= 0 {
		existing.Bathrooms = property.Bathrooms
	}
	if property.GarageSize >= 0 {
		existing.GarageSize = property.GarageSize
	}
	if property.GardenM2 >= 0 {
		existing.GardenM2 = property.GardenM2
	}
	if len(property.GasTypes) > 0 {
		existing.GasTypes = property.GasTypes
	}
	if len(property.Amenities) > 0 {
		existing.Amenities = property.Amenities
	}
	if len(property.Extras) > 0 {
		existing.Extras = property.Extras
	}
	if len(property.Utilities) > 0 {
		existing.Utilities = property.Utilities
	}
	if property.Notes != nil && *property.Notes != "" {
		existing.Notes = property.Notes
	}
	if property.Description != nil && *property.Description != "" {
		existing.Description = property.Description
	}
	if property.PropertyType != "" {
		existing.PropertyType = property.PropertyType
	}
	if property.TransactionType != "" {
		existing.TransactionType = property.TransactionType
	}
	if property.Status != "" {
		existing.Status = property.Status
	}
	// Preserve IsOccupied and IsFurnished from request only if they were explicitly set
	if property.IsOccupied {
		existing.IsOccupied = property.IsOccupied
	}
	if property.IsFurnished {
		existing.IsFurnished = property.IsFurnished
	}

	// Validate merged data
	if existing.Address == "" {
		logrus.Error("Address cannot be empty")
		return nil, errors.New("address cannot be empty")
	}
	if existing.Price <= 0 {
		logrus.Error("Price must be greater than zero")
		return nil, errors.New("price must be greater than zero")
	}

	updatedProperty, err := p.propertyRepo.Update(ctx, existing)
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

	err = p.imageRepo.DeleteImagesByPropertyID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
