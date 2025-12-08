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

func (service *UserService) GetAllUsers() ([]*models.UserResponse, error) {
	query := "SELECT id_usuario, usuario, nombre_usuario, role FROM Usuarios WHERE borrado_en IS NULL"
	rows, err := service.DB.Query(query)
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.UserResponse
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Nombre, &user.Role)
		if err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, user.ToResponse())
	}

	if err = rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	return users, nil
}

// Function to retrieve a user by ID
func (service *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	user := &models.User{}
	query := "SELECT id_usuario, usuario, nombre_usuario, role FROM Usuarios WHERE id_usuario = ? AND borrado_en IS NULL"
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
	query := "select id_usuario, usuario, nombre_usuario, password_usuario, role from Usuarios where usuario = ? AND borrado_en IS NULL;"
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
	if user.Email == "" || user.Nombre == "" || user.Role == "" || user.Password == "" {
		log.Println("Email, nombre and role must be provided")
		return nil, errors.New("email, nombre, password and role must be provided")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}
	user.Password = string(hashedPassword)

	query := "INSERT INTO Usuarios (usuario, nombre_usuario, password_usuario, role, creado_en, actualizado_en) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = service.DB.Exec(query, user.Email, user.Nombre, user.Password, user.Role, time.Now(), time.Now())
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}

	return user.ToResponse(), nil
}

func (service *UserService) DeleteUser(id int) error {
	query := "UPDATE Usuarios SET borrado_en = ? WHERE id_usuario = ? AND borrado_en IS NULL"
	_, err := service.DB.Exec(query, time.Now(), id)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

func (service *UserService) SetPasswordUser(id int, password string) (*models.UserResponse, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return nil, errors.New("password must be at most 72 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return nil, err
	}
	query := "UPDATE Usuarios SET password_usuario = ? WHERE id_usuario = ? AND borrado_en IS NULL"
	_, err = service.DB.Exec(query, hashedPassword, id)
	if err != nil {
		log.Println("Error updating user password:", err)
		return nil, err
	}

	log.Printf("Password updated for user ID %d", id)
	log.Println("Password is:", password)

	user, err := service.GetUserByID(id)
	if err != nil {
		log.Println("Error fetching user after password update:", err)
		return nil, errors.New("failed to fetch user")
	}
	if user == nil {
		log.Println("User not found after password update")
		return nil, errors.New("user not found")
	}

	return user, nil
}