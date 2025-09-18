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
func (service *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	user := &models.User{}
	query := "SELECT id_usuario, usuario, nombre_usuario, role FROM Usuarios WHERE id_usuario = ?"
	err := service.DB.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Nombre, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return user.ToResponse(), nil
}

// Function to retrieve a user by email and password
func (service *UserService) Login(email string, password string) (*models.UserResponse, string, error) {
	if email == "" || password == "" {
		log.Println("Email and password must be provided")
		return nil, "", errors.New("email and password must be provided")
	}

	user := &models.User{}
	query := "select id_usuario, usuario, nombre_usuario, password_usuario, role from Usuarios where usuario = ?;"
	err := service.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Nombre, &user.Password, &user.Role)

	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, "", errors.New("no such user found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Password mismatch:", err)
		return nil, "", errors.New("invalid credentials")
	}

	token, err := GenerateToken(user)
	if err != nil {
		log.Println("Error generating token:", err)
		return nil, "", errors.New("failed to generate token")
	}

	return user.ToResponse(), token, nil
}

func (service *UserService) CreateUser(user *models.User) (*models.UserResponse, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}
	if user.Email == "" || user.Password == "" || user.Nombre == "" || user.Role == "" {
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

	query := "INSERT INTO Usuarios (usuario, nombre_usuario, password_usuario, role, creado_en, actualizado_en) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = service.DB.Exec(query, user.Email, user.Nombre, hashedPassword, user.Role, time.Now(), time.Now())
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}
	return user.ToResponse(), nil
}
