package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)
// JwtClaims represents the JWT claims used for authentication.
type JWTClaims struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// LoginResponse represents the response returned upon successful login.
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
}

// LoginData represents the data required for user login.
type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshToken represents a refresh token used for obtaining new JWT tokens.
type RefreshToken struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	UserID    uint      `json:"user_id"`
	ExpiresAt int64     `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// RefreshTokenData represents the data required to refresh a JWT token.
type RefreshTokenData struct {
	JwtToken     string `json:"jwt_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
