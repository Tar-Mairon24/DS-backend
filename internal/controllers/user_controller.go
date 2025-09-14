package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/models"
	"backend/internal/services"
)

type UserController struct {
	UserService *services.UserService
}

// Constructor for UserController
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// GET /users/:id
func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("Error converting user ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		log.Println("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /login
// Login method for the UserController
func (ctrl *UserController) Login(c *gin.Context) {
	loginData := models.UserLoginData{}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println("Error binding login data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := ctrl.UserService.Login(loginData.Email, loginData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /create
func (ctrl *UserController) CreateUser(c *gin.Context){
	user := models.User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding user data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdUser, err := ctrl.UserService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

