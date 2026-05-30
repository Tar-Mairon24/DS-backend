package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ds-backend/internal/domain/models"
)

// TestStringArrayScan tests the Scan method for StringArray type
func TestStringArrayScan(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    models.StringArray
		wantErr bool
		errMsg  string
	}{
		{
			name:    "success - nil value",
			input:   nil,
			want:    models.StringArray{},
			wantErr: false,
		},
		{
			name:    "success - valid JSON bytes",
			input:   []byte(`["gas","electric"]`),
			want:    models.StringArray{"gas", "electric"},
			wantErr: false,
		},
		{
			name:    "success - valid JSON string",
			input:   `["water","internet"]`,
			want:    models.StringArray{"water", "internet"},
			wantErr: false,
		},
		{
			name:    "success - empty array bytes",
			input:   []byte(`[]`),
			want:    models.StringArray{},
			wantErr: false,
		},
		{
			name:    "success - empty array string",
			input:   `[]`,
			want:    models.StringArray{},
			wantErr: false,
		},
		{
			name:    "success - single element",
			input:   []byte(`["pool"]`),
			want:    models.StringArray{"pool"},
			wantErr: false,
		},
		{
			name:    "error - invalid type (integer)",
			input:   123,
			want:    nil,
			wantErr: true,
			errMsg:  "cannot scan StringArray",
		},
		{
			name:    "error - invalid type (float)",
			input:   123.45,
			want:    nil,
			wantErr: true,
			errMsg:  "cannot scan StringArray",
		},
		{
			name:    "error - invalid JSON bytes",
			input:   []byte(`not a json`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error - invalid JSON string",
			input:   `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa models.StringArray
			err := sa.Scan(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, sa)
			}
		})
	}
}

// TestStringArrayValue tests the Value method for StringArray type
func TestStringArrayValue(t *testing.T) {
	tests := []struct {
		name    string
		input   models.StringArray
		wantErr bool
		verify  func(t *testing.T, val any)
	}{
		{
			name:    "success - empty array",
			input:   models.StringArray{},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				assert.Equal(t, "[]", val)
			},
		},
		{
			name:    "success - single element",
			input:   models.StringArray{"pool"},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result []string
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, []string{"pool"}, result)
			},
		},
		{
			name:    "success - multiple elements",
			input:   models.StringArray{"gas", "electric", "water"},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result []string
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, []string{"gas", "electric", "water"}, result)
			},
		},
		{
			name:    "success - with special characters",
			input:   models.StringArray{"pool & gym", "24/7 security"},
			wantErr: false,
			verify: func(t *testing.T, val any) {
				jsonBytes, ok := val.([]byte)
				require.True(t, ok)
				var result []string
				err := json.Unmarshal(jsonBytes, &result)
				require.NoError(t, err)
				assert.Equal(t, []string{"pool & gym", "24/7 security"}, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := tt.input.Value()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.verify(t, val)
			}
		})
	}
}

// TestPropertyToCard tests the ToCard method for Property model
func TestPropertyToCard(t *testing.T) {
	tests := []struct {
		name    string
		input   *models.Property
		want    *models.PropertyCard
		wantErr bool
	}{
		{
			name: "success - full property to card",
			input: &models.Property{
				ID:              1,
				Title:           "Luxury Villa",
				Price:           500000.0,
				Bedrooms:        5,
				Bathrooms:       4,
				ConstructionM2:  350,
				City:            "Springfield",
				Neighborhood:    "Green Acres",
				PropertyType:    models.TypeHouse,
				TransactionType: models.TransactionSale,
				Status:          models.StatusAvailable,
				CreatedAt:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			want: &models.PropertyCard{
				ID:              1,
				Title:           "Luxury Villa",
				Price:           500000.0,
				Bedrooms:        5,
				Bathrooms:       4,
				ConstructionM2:  350,
				City:            "Springfield",
				Neighborhood:    "Green Acres",
				PropertyType:    models.TypeHouse,
				TransactionType: models.TransactionSale,
				Status:          models.StatusAvailable,
				CreatedAt:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "success - minimal property to card",
			input: &models.Property{
				ID:    2,
				Title: "Simple Property",
			},
			want: &models.PropertyCard{
				ID:    2,
				Title: "Simple Property",
			},
			wantErr: false,
		},
		{
			name: "success - property with zero values",
			input: &models.Property{
				ID:              0,
				Title:           "",
				Price:           0.0,
				Bedrooms:        0,
				Bathrooms:       0,
				ConstructionM2:  0,
				City:            "",
				Neighborhood:    "",
				PropertyType:    "",
				TransactionType: "",
				Status:          "",
			},
			want: &models.PropertyCard{
				ID:              0,
				Title:           "",
				Price:           0.0,
				Bedrooms:        0,
				Bathrooms:       0,
				ConstructionM2:  0,
				City:            "",
				Neighborhood:    "",
				PropertyType:    "",
				TransactionType: "",
				Status:          "",
			},
			wantErr: false,
		},
		{
			name:    "error - nil property",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - property with different types",
			input: &models.Property{
				ID:              3,
				Title:           "Apartment",
				PropertyType:    models.TypeApartment,
				TransactionType: models.TransactionRental,
				Status:          models.StatusRented,
			},
			want: &models.PropertyCard{
				ID:              3,
				Title:           "Apartment",
				PropertyType:    models.TypeApartment,
				TransactionType: models.TransactionRental,
				Status:          models.StatusRented,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() {
					_ = tt.input.ToCard()
				})
			} else {
				card := tt.input.ToCard()
				require.NotNil(t, card)
				assert.Equal(t, tt.want.ID, card.ID)
				assert.Equal(t, tt.want.Title, card.Title)
				assert.Equal(t, tt.want.Price, card.Price)
				assert.Equal(t, tt.want.Bedrooms, card.Bedrooms)
				assert.Equal(t, tt.want.Bathrooms, card.Bathrooms)
				assert.Equal(t, tt.want.ConstructionM2, card.ConstructionM2)
				assert.Equal(t, tt.want.City, card.City)
				assert.Equal(t, tt.want.Neighborhood, card.Neighborhood)
				assert.Equal(t, tt.want.PropertyType, card.PropertyType)
				assert.Equal(t, tt.want.TransactionType, card.TransactionType)
				assert.Equal(t, tt.want.Status, card.Status)
				assert.Equal(t, tt.want.CreatedAt, card.CreatedAt)
			}
		})
	}
}

// TestPropertyToResponse tests the ToResponse method for Property model
func TestPropertyToResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   *models.Property
		want    *models.PropertyResponse
		wantErr bool
	}{
		{
			name: "success - full property to response",
			input: &models.Property{
				ID:              1,
				Title:           "Test Property",
				Address:         "123 Main St",
				Neighborhood:    "Downtown",
				City:            "Metropolis",
				Zone:            "Central",
				Reference:       "REF001",
				Price:           100000.0,
				ConstructionM2:  120,
				LandM2:          200,
				IsOccupied:      true,
				IsFurnished:     false,
				Floors:          2,
				Bedrooms:        3,
				Bathrooms:       2,
				GarageSize:      1,
				GardenM2:        50,
				GasTypes:        models.StringArray{"gas"},
				Amenities:       models.StringArray{"pool"},
				Extras:          models.StringArray{"alarm"},
				Utilities:       models.StringArray{"water"},
				PropertyType:    models.TypeHouse,
				TransactionType: models.TransactionSale,
				Status:          models.StatusAvailable,
				CreatedAt:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: &models.PropertyResponse{
				ID:              1,
				Title:           "Test Property",
				Address:         "123 Main St",
				Neighborhood:    "Downtown",
				City:            "Metropolis",
				Zone:            "Central",
				Reference:       "REF001",
				Price:           100000.0,
				ConstructionM2:  120,
				LandM2:          200,
				IsOccupied:      true,
				IsFurnished:     false,
				Floors:          2,
				Bedrooms:        3,
				Bathrooms:       2,
				GarageSize:      1,
				GardenM2:        50,
				GasTypes:        models.StringArray{"gas"},
				Amenities:       models.StringArray{"pool"},
				Extras:          models.StringArray{"alarm"},
				Utilities:       models.StringArray{"water"},
				PropertyType:    models.TypeHouse,
				TransactionType: models.TransactionSale,
				Status:          models.StatusAvailable,
				CreatedAt:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "success - minimal property to response",
			input: &models.Property{
				ID:    2,
				Title: "Simple Property",
			},
			want: &models.PropertyResponse{
				ID:    2,
				Title: "Simple Property",
			},
			wantErr: false,
		},
		{
			name: "success - property with nil pointers",
			input: &models.Property{
				ID:          3,
				Title:       "Property with nils",
				Notes:       nil,
				Description: nil,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
			want: &models.PropertyResponse{
				ID:          3,
				Title:       "Property with nils",
				Notes:       nil,
				Description: nil,
				UpdatedAt:   nil,
			},
			wantErr: false,
		},
		{
			name: "success - property with optional fields",
			input: &models.Property{
				ID:          4,
				Title:       "Property with optional fields",
				Notes:       stringPtr("Some notes"),
				Description: stringPtr("A detailed description"),
				UpdatedAt:   timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			},
			want: &models.PropertyResponse{
				ID:          4,
				Title:       "Property with optional fields",
				Notes:       stringPtr("Some notes"),
				Description: stringPtr("A detailed description"),
				UpdatedAt:   timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name:    "error - nil property",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "success - property with arrays",
			input: &models.Property{
				ID:        5,
				Title:     "Property with arrays",
				GasTypes:  models.StringArray{"gas", "electric"},
				Amenities: models.StringArray{"pool", "gym", "parking"},
				Extras:    models.StringArray{"alarm", "security"},
				Utilities: models.StringArray{"water", "internet", "cable"},
			},
			want: &models.PropertyResponse{
				ID:        5,
				Title:     "Property with arrays",
				GasTypes:  models.StringArray{"gas", "electric"},
				Amenities: models.StringArray{"pool", "gym", "parking"},
				Extras:    models.StringArray{"alarm", "security"},
				Utilities: models.StringArray{"water", "internet", "cable"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() {
					_ = tt.input.ToResponse()
				})
			} else {
				resp := tt.input.ToResponse()
				require.NotNil(t, resp)
				assert.Equal(t, tt.want.ID, resp.ID)
				assert.Equal(t, tt.want.Title, resp.Title)
				assert.Equal(t, tt.want.Address, resp.Address)
				assert.Equal(t, tt.want.Neighborhood, resp.Neighborhood)
				assert.Equal(t, tt.want.City, resp.City)
				assert.Equal(t, tt.want.Zone, resp.Zone)
				assert.Equal(t, tt.want.Reference, resp.Reference)
				assert.Equal(t, tt.want.Price, resp.Price)
				assert.Equal(t, tt.want.ConstructionM2, resp.ConstructionM2)
				assert.Equal(t, tt.want.LandM2, resp.LandM2)
				assert.Equal(t, tt.want.IsOccupied, resp.IsOccupied)
				assert.Equal(t, tt.want.IsFurnished, resp.IsFurnished)
				assert.Equal(t, tt.want.Floors, resp.Floors)
				assert.Equal(t, tt.want.Bedrooms, resp.Bedrooms)
				assert.Equal(t, tt.want.Bathrooms, resp.Bathrooms)
				assert.Equal(t, tt.want.GarageSize, resp.GarageSize)
				assert.Equal(t, tt.want.GardenM2, resp.GardenM2)
				assert.Equal(t, tt.want.GasTypes, resp.GasTypes)
				assert.Equal(t, tt.want.Amenities, resp.Amenities)
				assert.Equal(t, tt.want.Extras, resp.Extras)
				assert.Equal(t, tt.want.Utilities, resp.Utilities)
				assert.Equal(t, tt.want.Notes, resp.Notes)
				assert.Equal(t, tt.want.Description, resp.Description)
				assert.Equal(t, tt.want.PropertyType, resp.PropertyType)
				assert.Equal(t, tt.want.TransactionType, resp.TransactionType)
				assert.Equal(t, tt.want.Status, resp.Status)
				assert.Equal(t, tt.want.CreatedAt, resp.CreatedAt)
				assert.Equal(t, tt.want.UpdatedAt, resp.UpdatedAt)
			}
		})
	}
}

// Helper functions for pointer creation
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
