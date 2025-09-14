package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"backend/internal/models"
)

type UserService struct {
	DB *sql.DB
}

// Constructor for the UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

// Function to retrieve a user by ID
func (service *UserService) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM Usuarios WHERE id_usuario = ?"
	err := service.DB.QueryRow(query, id).Scan(&user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return user, nil
}

// Function to retrieve a user by email and password
func (service *UserService) Login(email, password string) (*models.User, error) {
	user := &models.User{}
	query := "select * from Usuarios where usuario = ? and password_usuario = ?;"
	err := service.DB.QueryRow(query, email, password).Scan(&user.Email, &user.Password)
	log.Println("User:", user.Email+" "+user.Password)
	log.Println("Error:", err)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("invalid credentials")
		}
		return user, err
	}
	return user, nil
}

func (service *UserService) CreateUser(user *models.User) (*models.UserResponse, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}
	if user.Email == "" || user.Password == "" || user.Nombre == "" || user.Rol == "" {
		return nil, errors.New("missing required user fields")
	}

	if len(user.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	if len(user.Password) > 72 {
		return nil, errors.New("password must be at most 72 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}

	query := "INSERT INTO Usuarios (usuario, nombre_usuario, password_usuario, rol_usuario, creado_en, actualizado_en) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = service.DB.Exec(query, user.Email, user.Nombre, hashedPassword, user.Rol, time.Now(), time.Now())
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}
	return user.ToResponse(), nil
}
