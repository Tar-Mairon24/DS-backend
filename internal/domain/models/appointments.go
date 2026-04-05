package models

import (
	"encoding/json"
	"time"
)

// CustomTime handles ISO 8601 unmarshaling
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			ct.Time = t
			return nil
		}
	}
	return err
}

type Appointment struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Status      StatusType `json:"status"`
	Notes       *string    `json:"notes,omitempty"`
	ClientID    uint       `json:"client_id"`
	OwnerID     uint       `json:"owner_id"`
	PropertyID  uint       `json:"property_id"`
	AgentIDs    []uint     `json:"agent_ids,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"-"`
}

type AppointmentRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	StartDate   CustomTime `json:"start_date" validate:"required"`
	EndDate     CustomTime `json:"end_date" validate:"required"`
	Status      StatusType `json:"status" validate:"required,oneof=scheduled completed cancelled archived no-show"`
	Notes       *string    `json:"notes,omitempty"`
	ClientID    uint       `json:"client_id" validate:"required"`
	PropertyID  uint       `json:"property_id" validate:"required"`
	AgentIDs    []uint     `json:"agent_ids,omitempty"`
}

type AppointmentCalendarView struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	Status      StatusType   `json:"status"`
	Property    PropertyTile `json:"property"`
	ClientID    uint         `json:"client_id"`
	OwnerID     uint         `json:"owner_id"`
}

type AppointmentDetail struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Status          StatusType `json:"status"`
	Notes           *string    `json:"notes,omitempty"`
	PropertyID      uint       `json:"property_id"`
	PropertyTitle   string     `json:"property_title"`
	PropertyAddress string     `json:"property_address"`
	ClientID        uint       `json:"client_id"`
	ClientName      string     `json:"client_name"`
	ClientEmail     string     `json:"client_email"`
	ClientPhone     *string    `json:"client_phone,omitempty"`
	OwnerID         *uint      `json:"owner_id,omitempty"`
	OwnerName       *string    `json:"owner_name,omitempty"`
	OwnerEmail      *string    `json:"owner_email,omitempty"`
	Agents          []UserInfo `json:"agents"`
}

type rawAppointmentDetail struct {
	AppointmentDetail
	AgentsJSON json.RawMessage
}

type RescheduleRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type StatusUpdateRequest struct {
	Status StatusType `json:"status" binding:"required,oneof=scheduled completed cancelled archived no-show"`
}

type UpdateStatusRequest struct {
	Status StatusType `json:"status" binding:"required,oneof=scheduled completed cancelled archived no-show"`
}

type StatusType string

const (
	AppointmentStatusScheduled StatusType = "scheduled"
	AppointmentStatusCompleted StatusType = "completed"
	AppointmentStatusCancelled StatusType = "cancelled"
	AppointmentStatusArchived  StatusType = "archived"
	AppointmentStatusNoShow    StatusType = "no-show"
)
