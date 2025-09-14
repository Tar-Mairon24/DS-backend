package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
