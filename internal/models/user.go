package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Nombre   string `json:"nombre"`
	Password string `json:"password"`
	Rol      string `json:"rol"`
	CreadoEn string `json:"creado_en"`
	ActualizadoEn string `json:"actualizado_en"`
}

type UserResponse struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Nombre string `json:"nombre"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Rol      string `json:"rol"`
	jwt.RegisteredClaims
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:     u.ID,
		Email:  u.Email,
		Nombre: u.Nombre,
	}
}