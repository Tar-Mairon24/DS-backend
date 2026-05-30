package models

import "time"

type Image struct {
	ID          uint      `json:"id"`
	PropertyID  uint      `json:"property_id"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	MainImage   bool      `json:"main_image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"-"`
}

type SaveImageRequest struct {
	Description string `form:"description"`
	MainImage   bool   `form:"main_image"`
}