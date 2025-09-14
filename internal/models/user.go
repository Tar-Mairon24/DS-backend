package models

type User struct {
	Email    string `json:"email"`
	Nombre   string `json:"nombre"`
	Password string `json:"password"`
	Rol      string `json:"rol"`
	CreadoEn string `json:"creado_en"`
	ActualizadoEn string `json:"actualizado_en"`
}

type UserResponse struct {
	Email  string `json:"email"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		Email:  u.Email,
		Nombre: u.Nombre,
		Rol:    u.Rol,
	}
}