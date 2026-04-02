package models

import (
	"encoding/json"
	"time"
)

type Appointment struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_time"` 
	EndDate     time.Time  `json:"end_time"`	
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
	StartDate   time.Time `json:"start_time" validate:"required"`
	EndDate     time.Time `json:"end_time" validate:"required,gtfield=StartDate"`
	Status      string    `json:"status" validate:"required,oneof=scheduled completed cancelled"`
	Notes       *string   `json:"notes,omitempty"`
	ClientID    uint      `json:"client_id" validate:"required"`
	OwnerID     uint      `json:"owner_id" validate:"required"`
	PropertyID  uint      `json:"property_id" validate:"required"`
}

type AppointmentCalendarView struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_time"`
	EndDate     time.Time  `json:"end_time"`
	Status      string     `json:"status"`
	Property 	PropertyTile   `json:"property"`
	ClientID    uint       `json:"client_id"`
	OwnerID     uint       `json:"owner_id"`
}

type AppointmentDetail struct {
    ID              uint        `json:"id"`
    Title           string      `json:"title"`
    Description     string      `json:"description"`
    StartDate       time.Time   `json:"start_date"`
    EndDate         time.Time   `json:"end_date"`
    Status          string      `json:"status"`
    Notes           *string     `json:"notes,omitempty"`
    PropertyID      uint        `json:"property_id"`
    PropertyTitle   string      `json:"property_title"`
    PropertyAddress string      `json:"property_address"`
    ClientID        uint        `json:"client_id"`
    ClientName      string      `json:"client_name"`
    ClientEmail     string      `json:"client_email"`
    ClientPhone     *string     `json:"client_phone,omitempty"`
    OwnerID         *uint       `json:"owner_id,omitempty"`
    OwnerName       *string     `json:"owner_name,omitempty"`
    OwnerEmail      *string     `json:"owner_email,omitempty"`
    Agents          []UserInfo `json:"agents"`
}

type rawAppointmentDetail struct {
    AppointmentDetail
    AgentsJSON json.RawMessage
}
