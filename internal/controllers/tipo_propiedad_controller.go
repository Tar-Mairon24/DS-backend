package controllers

import (
	"backend/internal/services"
	"backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TipoPropiedadController is the controller for the TipoPropiedad model
type TipoPropiedadController struct {
	TipoPropiedadService *services.TipoPropiedadService
}

// NewTipoPropiedadController is the constructor for the TipoPropiedadController
func NewTipoPropiedadController(tipoPropiedadService *services.TipoPropiedadService) *TipoPropiedadController {
	return &TipoPropiedadController{
		TipoPropiedadService: tipoPropiedadService,
	}
}

// GET /TipoPropiedad/:id
func (ctrl *TipoPropiedadController) GetTipoPropiedad(c *gin.Context) {
	idTipoPropiedadParam := c.Param("id")
	idTipoPropiedad, err := strconv.Atoi(idTipoPropiedadParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tipo propiedad ID"})
		return
	}

	propiedades, err := ctrl.TipoPropiedadService.GetTipoPropiedad(idTipoPropiedad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve propiedades"})
		return
	}

	if propiedades == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No propiedades found"})
		return
	}

	c.JSON(http.StatusOK, propiedades)
}

//POST /tipopropiedad
//Funcion que sirve para insertar un nuevo tipo de propiedad en la base de datos
func (ctrl *TipoPropiedadController) CreateTipoPropiedad(c *gin.Context) {
	var tipoPropiedad models.TipoPropiedad
	err := c.BindJSON(&tipoPropiedad)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tipo propiedad data"})
		return
	}

	id, err := ctrl.TipoPropiedadService.CreateTipoPropiedad(&tipoPropiedad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tipo propiedad"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}