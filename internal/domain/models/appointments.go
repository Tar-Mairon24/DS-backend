package models

import "time"

type Appointment struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"` 
	EndTime     time.Time  `json:"end_time"`	
	Status      string     `json:"status"`
	Notes       *string     `json:"notes,omitempty"`
	ClientID    uint       `json:"client_id"`
	OwnerID     uint       `json:"owner_id"`
	PropertyID  uint       `json:"property_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"-"`
}

type AppointmentRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Status      string    `json:"status" validate:"required,oneof=scheduled completed cancelled"`
	Notes       *string   `json:"notes,omitempty"`
	ClientID    uint      `json:"client_id" validate:"required"`
	OwnerID     uint      `json:"owner_id" validate:"required"`
	PropertyID  uint      `json:"property_id" validate:"required"`
}

type AppointmentMonthView struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Status      string     `json:"status"`
	Property 	PropertyTile   `json:"property,omitempty"`
	ClientID    uint       `json:"client_id"`
	OwnerID     uint       `json:"owner_id"`
}
