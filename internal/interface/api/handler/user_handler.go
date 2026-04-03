package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type UserHandler struct {
	userUsecase ports.UserUseCase
}

func NewUserHandler(userUsecase ports.UserUseCase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	logrus.Info("GetUsers endpoint called")

	users, err := h.userUsecase.GetAllUsers(c.Request.Context())
	if err != nil {
		logrus.WithError(err).Error("Failed to get users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"message": err.Error(),
		})
		return
	}

	if len(users) == 0 {
		logrus.Warn("No users found")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "No users found",
			"message": "No users available in the database",
		})
		return
	}

	logrus.Infof("Retrieved %d users", len(users))
	c.JSON(http.StatusOK, gin.H{
		"data":    users,
		"message": "Users retrieved successfully",
		"count":   len(users),
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	logrus.Infof("GetUserByID endpoint called with ID: %s", userIDStr)

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid user ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	user, err := h.userUsecase.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("User retrieved successfully: %v", user)
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	logrus.Infof("CreateUser endpoint called")
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	if req.Role != "admin" {
		req.Role = "agente"
	}

	userResponse, err := h.userUsecase.CreateUser(c.Request.Context(), &req, models.CreateContextSelf)
	if err != nil {
		logrus.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("User created successfully: %v", userResponse)
	c.JSON(http.StatusCreated, gin.H{
		"data":    userResponse,
		"message": "User created successfully",
	})
}

func (h *UserHandler) CreateOwner(c *gin.Context) {
	logrus.Infof("CreateOwner endpoint called")
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	req.Role = "owner"
	
	userResponse, err := h.userUsecase.CreateUser(c.Request.Context(), &req, models.CreateContextAgent)
	if err != nil {
		logrus.WithError(err).Error("Failed to create owner")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create owner",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Owner created successfully: %v", userResponse)
	c.JSON(http.StatusCreated, gin.H{
		"data":    userResponse,
		"message": "Owner created successfully",
	})
}

func (h *UserHandler) CreateClient(c *gin.Context) {
	logrus.Infof("CreateClient endpoint called")
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	req.Role = "client"
	
	userResponse, err := h.userUsecase.CreateUser(c.Request.Context(), &req, models.CreateContextAgent)
	if err != nil {
		logrus.WithError(err).Error("Failed to create client")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create client",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Client created successfully: %v", userResponse)
	c.JSON(http.StatusCreated, gin.H{
		"data":    userResponse,
		"message": "Client created successfully",
	})
}


func (h *UserHandler) UpdateUser(c *gin.Context) {
	logrus.Infof("UpdateUser endpoint called")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	UserResponse, err := h.userUsecase.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		logrus.WithError(err).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("User updated successfully: %v", UserResponse)
	c.JSON(http.StatusOK, gin.H{
		"data":    UserResponse,
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	logrus.Infof("DeleteUser endpoint called")
	userIDStr := c.Param("id")
	logrus.Infof("DeleteUser endpoint called with ID: %s", userIDStr)

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid user ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	if err := h.userUsecase.DeleteUser(c.Request.Context(), uint(userID)); err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("DeleteUser endpoint called with ID: %d", userID)
	c.JSON(http.StatusNoContent, nil)
}
