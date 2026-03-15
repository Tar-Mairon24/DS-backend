package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/models"
)

type UserService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (service *UserService) GetAllUsers() ([]*models.UserResponse, error) {
	query := service.sq.Select("id_usuario", "usuario", "nombre_usuario", "role").From("Usuarios").Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}

	rows, err := service.DB.Query(sqlStr, args...)
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
	query := service.sq.Select("id_usuario", "usuario", "nombre_usuario", "role").From("Usuarios").Where(sq.Eq{"id_usuario": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	err = service.DB.QueryRow(sqlStr, args...).Scan(&user.ID, &user.Email, &user.Nombre, &user.Role)
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
	query := service.sq.Select("id_usuario", "usuario", "nombre_usuario", "password_usuario", "role").From("Usuarios").Where(sq.Eq{"usuario": email}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, "", errors.New("failed to build SQL query")
	}

	err = service.DB.QueryRow(sqlStr, args...).Scan(&user.ID, &user.Email, &user.Nombre, &user.Password, &user.Role)
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

	query := service.sq.Insert("Usuarios").Columns("usuario", "nombre_usuario", "password_usuario", "role", "creado_en", "actualizado_en").Values(user.Email, user.Nombre, user.Password, user.Role, time.Now(), time.Now())
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	_, err = service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}

	return user.ToResponse(), nil
}

func (service *UserService) DeleteUser(id int) error {
	query := service.sq.Update("Usuarios").Set("borrado_en", time.Now()).Where(sq.Eq{"id_usuario": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	_, err = service.DB.Exec(sqlStr, args...)
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
	query := service.sq.Update("Usuarios").Set("password_usuario", hashedPassword).Where(sq.Eq{"id_usuario": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	_, err = service.DB.Exec(sqlStr, args...)
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
