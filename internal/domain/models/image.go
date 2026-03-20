package models

import "time"

type Image struct {
	ID          uint      `json:"id"`
	PropertyID  uint      `json:"property_id"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	MainImage   bool      `json:"main_image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}