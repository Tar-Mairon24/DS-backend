package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
)

// Property represents a real estate property in the system
type Property struct {
    ID              uint               `json:"id"`
    Title           string             `json:"title"`
    Address         string             `json:"address"`
    Neighborhood    string             `json:"neighborhood"`
    City            string             `json:"city"`
	Zone			string             `json:"zone"`
    Reference       string             `json:"reference"`
    Price           float64            `json:"price"`
    ConstructionM2  int                `json:"construction_m2"`
    LandM2          int                `json:"land_m2"`
    IsOccupied      bool               `json:"is_occupied"`
    IsFurnished     bool               `json:"is_furnished"`
    Floors          int                `json:"floors"`
    Bedrooms        int                `json:"bedrooms"`
    Bathrooms       int                `json:"bathrooms"`
    GarageSize      int                `json:"garage_size"` // Number of cars
    GardenM2        int                `json:"garden_m2"`
    GasTypes        StringArray        `json:"gas_types"`
    Amenities       StringArray        `json:"amenities"`
    Extras          StringArray        `json:"extras"`
    Utilities       StringArray        `json:"utilities"`
    Notes           string             `json:"notes"`
    OwnerID         uint               `json:"owner_id"`
    UserID          uint               `json:"user_id"`
	PropertyType    PropertyType       `json:"property_type"`
    TransactionType TransactionType    `json:"transaction_type"`
    Status          PropertyStatus     `json:"status"`
    CreatedAt       time.Time          `json:"created_at"`
    UpdatedAt       *time.Time         `json:"updated_at,omitempty"`
    DeletedAt       *time.Time         `json:"-"`
	Owner           *User              `json:"owner,omitempty"`
    Users           []User             `json:"users,omitempty"`
}

// PropertyResponse represents the public view of a property
type PropertyResponse struct {
    ID              uint            `json:"id"`
    Title           string          `json:"title"`
    Address         string          `json:"address"`
    Neighborhood    string          `json:"neighborhood"`
    City            string          `json:"city"`
	Zone            string          `json:"zone"`
    Price           float64         `json:"price"`
    ConstructionM2  int             `json:"construction_m2"`
    LandM2          int             `json:"land_m2"`
    IsOccupied      bool            `json:"is_occupied"`
    IsFurnished     bool            `json:"is_furnished"`
    Floors          int             `json:"floors"`
    Bedrooms        int             `json:"bedrooms"`
    Bathrooms       int             `json:"bathrooms"`
    GarageSize      int             `json:"garage_size"`
    GardenM2        int             `json:"garden_m2"`
    GasTypes        StringArray     `json:"gas_types"`
    Amenities       StringArray     `json:"amenities"`
    Extras          StringArray     `json:"extras"`
    Utilities       StringArray     `json:"utilities"`
	PropertyType   	PropertyType    `json:"property_type"`
    TransactionType TransactionType `json:"transaction_type"`
    Status          PropertyStatus  `json:"status"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       *time.Time       `json:"updated_at,omitempty"`
    Agent          	*UserResponse   `json:"agent,omitempty"` // Agent handling the property
}

// PropertyCard represents a simplified property view for listings
type PropertyCard struct {
    ID              uint            `json:"id"`
    Title           string          `json:"title"`
    Price           float64         `json:"price"`
    Bedrooms        int             `json:"bedrooms"`
    Bathrooms       int             `json:"bathrooms"`
    ConstructionM2  int             `json:"construction_m2"`
    City            string          `json:"city"`
    Neighborhood    string          `json:"neighborhood"`
	PropertyType    PropertyType    `json:"property_type"`
    TransactionType TransactionType `json:"transaction_type"`
    Status          PropertyStatus  `json:"status"`
    MainImagePath   *string         `json:"main_image_path"`
    CreatedAt       time.Time       `json:"created_at"`
}

type TransactionType string

const (
    TransactionSale   TransactionType = "Venta"
    TransactionRental TransactionType = "Renta"
)

type PropertyStatus string

const (
    StatusAvailable PropertyStatus = "Disponible"
    StatusSold      PropertyStatus = "Vendido"
    StatusRented    PropertyStatus = "Alquilado"
    StatusReserved  PropertyStatus = "Reservado"
)

type PropertyType string 

const (
	TypeHouse     	PropertyType = "Casa"
	TypeApartment 	PropertyType = "Apartamento"
	TypeLand      	PropertyType = "Terreno"
	TypeCommercial 	PropertyType = "Comercial"
	TypeStorehouse 	PropertyType = "Almacén"
	TypeOffice     	PropertyType = "Oficina"
	TypeIndustrial  PropertyType = "Industrial"
	TypeOther      	PropertyType = "Otro"
)

type StringArray []string

func (sa *StringArray) Scan(value any) error {
    if value == nil {
        *sa = []string{}
        return nil
    }

    switch v := value.(type) {
    case []byte:
        return json.Unmarshal(v, sa)
    case string:
        return json.Unmarshal([]byte(v), sa)
    default:
        return errors.New("cannot scan StringArray")
    }
}

// Value implements the Valuer interface for database writing
func (sa StringArray) Value() (driver.Value, error) {
    if len(sa) == 0 {
        return "[]", nil
    }
    return json.Marshal(sa)
}

func (p *Property) ToResponse() *PropertyResponse {
    response := &PropertyResponse{
        ID:              p.ID,
        Title:           p.Title,
        Address:         p.Address,
        Neighborhood:    p.Neighborhood,
        City:            p.City,
        Zone:            p.Zone,
        Price:           p.Price,
        ConstructionM2:  p.ConstructionM2,
        LandM2:          p.LandM2,
        IsOccupied:      p.IsOccupied,
        IsFurnished:     p.IsFurnished,
        Floors:          p.Floors,
        Bedrooms:        p.Bedrooms,
        Bathrooms:       p.Bathrooms,
        GarageSize:      p.GarageSize,
        GardenM2:        p.GardenM2,
        GasTypes:        p.GasTypes,
        Amenities:       p.Amenities,
        Extras:          p.Extras,
        Utilities:       p.Utilities,
        PropertyType:    p.PropertyType,
        TransactionType: p.TransactionType,
        Status:          p.Status,
        CreatedAt:       p.CreatedAt,
        UpdatedAt:       p.UpdatedAt,
    }
    
    return response
}

func (p *Property) ToCard() *PropertyCard {
    return &PropertyCard{
        ID:              p.ID,
        Title:           p.Title,
        Price:           p.Price,
        Bedrooms:        p.Bedrooms,
        Bathrooms:       p.Bathrooms,
        ConstructionM2:  p.ConstructionM2,
        City:            p.City,
        Neighborhood:    p.Neighborhood,
        PropertyType:    p.PropertyType,
        TransactionType: p.TransactionType,
        Status:          p.Status,
        CreatedAt:       p.CreatedAt,
    }
}