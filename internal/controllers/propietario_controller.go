package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/services"
	"backend/internal/models"
)

// PropietarioController is the controller for the Propietario model
type PropietarioController struct {
	PropietarioService *services.PropietarioService
}

// NewPropietarioController is the constructor for the PropietarioController
func NewPropietarioController(propietarioService *services.PropietarioService) *PropietarioController {
	return &PropietarioController{
		PropietarioService: propietarioService,
	}
}

// GET /propietario/:id
func (ctrl *PropietarioController) GetPropietario(c *gin.Context) {
	idPropietarioParam := c.Param("id")
	idPropietario, err := strconv.Atoi(idPropietarioParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propietario ID"})
		return
	}

	propietario, err := ctrl.PropietarioService.GetPropietario(idPropietario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve propietario"})
		return
	}

	if propietario == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No propietario found"})
		return
	}

	c.JSON(http.StatusOK, propietario)
}

// POST /propietario
func (ctrl *PropietarioController) CreatePropietario(c *gin.Context) {
	var propietario models.Propietario
	if err := c.ShouldBindJSON(&propietario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propietario data"})
		return
	}

	idPropietario, err := ctrl.PropietarioService.CreatePropietario(&propietario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create propietario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id_propietario": idPropietario})
}

func (ctrl *PropietarioController) UpdatePropietario(c *gin.Context) {
	idPropietarioParam := c.Param("id")
	idPropietario, err := strconv.Atoi(idPropietarioParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propietario ID"})
		return
	}

	var propietario models.Propietario
	if err := c.ShouldBindJSON(&propietario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propietario data"})
		return
	}

	propietario.IDPropietario = idPropietario

	err = ctrl.PropietarioService.UpdatePropietario(&propietario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update propietario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Propietario updated successfully"})
}

func (ctrl *PropietarioController) DeletePropietario(c *gin.Context) {
	idPropietarioParam := c.Param("id")
	idPropietario, err := strconv.Atoi(idPropietarioParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propietario ID"})
		return
	}

	err = ctrl.PropietarioService.DeletePropietario(idPropietario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete propietario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Propietario deleted successfully"})
}