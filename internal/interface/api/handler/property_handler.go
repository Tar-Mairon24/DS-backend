package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type PropertyHandler struct {
	propertyUsecase ports.PropertyUseCase
}

func NewPropertyHandler(propertyUsecase ports.PropertyUseCase) *PropertyHandler {
	return &PropertyHandler{
		propertyUsecase: propertyUsecase,
	}
}

func (h *PropertyHandler) GetProperties(c *gin.Context) {
	logrus.Info("GetProperties endpoint called")

	properties, err := h.propertyUsecase.GetAllProperties(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve properties",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Retrieved %d properties", len(properties))
	c.JSON(http.StatusOK, properties)
}

func (h *PropertyHandler) GetPropertyCardByID(c *gin.Context) {
	logrus.Info("GetPropertyCardByID endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}

	propertyCard, err := h.propertyUsecase.GetPropertyCardByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve property card",
			"message": err.Error(),
		})
		return
	}

	if propertyCard == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Property not found",
			"message": "No property found with the given ID",
		})
		return
	}

	logrus.Infof("Retrieved property card with ID: %d", propertyCard.ID)
	c.JSON(http.StatusOK, propertyCard)
}

func (h *PropertyHandler) GetPropertyByID(c *gin.Context) {
	logrus.Info("GetPropertyByID endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}

	property, err := h.propertyUsecase.GetPropertyByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve property",
			"message": err.Error(),
		})
		return
	}

	if property == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Property not found",
			"message": "No property found with the given ID",
		})
		return
	}

	logrus.Infof("Retrieved property with ID: %d", id)
	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	logrus.Info("CreateProperty endpoint called")

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid property data",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	loggedUserID := userID.(uint)

	// If no agents provided, use logged-in user
	if len(property.UserID) == 0 {
		property.UserID = []uint{loggedUserID}
	}

	if property.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "owner_id is required",
			"message": "A property must be linked to an existing owner",
		})
		return
	}

	newProperty, err := h.propertyUsecase.CreateProperty(c.Request.Context(), &property)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create property",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Property created successfully with ID: %d", newProperty.ID)
	c.JSON(http.StatusCreated, newProperty)
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	logrus.Info("UpdateProperty endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}
	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid property data",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	loggedUserID := userID.(uint)

	property.ID = uint(id)
	
	// If no agents provided, use logged-in user
	if len(property.UserID) == 0 {
		property.UserID = []uint{loggedUserID}
	}

	updatedProperty, err := h.propertyUsecase.UpdateProperty(c.Request.Context(), &property)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update property",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Property updated successfully with ID: %d", updatedProperty.ID)
	c.JSON(http.StatusOK, updatedProperty)
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	logrus.Info("DeleteProperty endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}

	err = h.propertyUsecase.DeleteProperty(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete property",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Property with ID: %d deleted successfully", id)
	c.JSON(http.StatusNoContent, nil)
}
