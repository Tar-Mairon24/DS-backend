package models

import "time"

type User struct {
	ID            uint           `json:"id"`
	Username      string         `json:"username"`
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Role          string         `json:"role"`
	Verified      bool           `json:"verified"`
	Phone         *string        `json:"phone,omitempty"`
	MFAActivated  bool           `json:"mfa_activated"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     *time.Time     `json:"-"`
	Properties    []Property     `json:"properties,omitempty"`
	RefreshTokens []RefreshToken `json:"-"`
}

type UserResponse struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Role         string     `json:"role"`
	Verified     bool       `json:"verified"`
	Phone        *string    `json:"phone,omitempty"`
	MFAActivated bool       `json:"mfa_activated"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func (user *User) ToUserResponse() *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
