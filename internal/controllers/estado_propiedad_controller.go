package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EstadoPropiedadController struct {
	EstadoPropiedadService *services.EstadoPropiedadService
}

func NewEstadoPropiedadController(estadoPropiedadService *services.EstadoPropiedadService) *EstadoPropiedadController {
	return &EstadoPropiedadController{
		EstadoPropiedadService: estadoPropiedadService,
	}
}

// GET /estadopropiedad/:id
func (ctrl *EstadoPropiedadController) GetEstadoPropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propiedad ID"})
		return
	}

	estadoPropiedad, err := ctrl.EstadoPropiedadService.GetEstadoPropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve estado propiedad"})
		return
	}

	if estadoPropiedad == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No estado propiedad found"})
		return
	}

	c.JSON(http.StatusOK, estadoPropiedad)
}

// POST /estadopropiedad/
func (ctrl *EstadoPropiedadController) CreateEstadoPropiedad(c *gin.Context) {
	var estadoPropiedad models.EstadoPropiedades
	if err := c.ShouldBindJSON(&estadoPropiedad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := ctrl.EstadoPropiedadService.CreateEstadoPropiedad(&estadoPropiedad)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert estado propiedad"})
		return
	}
	log.Printf("Created estado propiedad with ID: %d", id)

	c.JSON(http.StatusCreated, estadoPropiedad)
}

// PUT /Update/EstadoPropiedad/:id
// Function that updates the state of the property
// func (ctrl *EstadoPropiedadController) UpdateEstadoPropiedad(c *gin.Context) {
// 	var estadoPropiedad models.EstadoPropiedades
// 	if err := c.ShouldBindJSON(&estadoPropiedad); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := ctrl.EstadoPropiedadService.UpdateEstadoPropiedad(&estadoPropiedad)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update estado propiedad"})
// 		return
// 	}
// 	log.Printf("Updated estado propiedad with ID: %d", estadoPropiedad.IDEstadoPropiedades)

// 	c.JSON(http.StatusOK, estadoPropiedad)
// }

// DELETE /eliminar/estadoPropiedad
// Function that deletes a EstadoPropiedad from the database
func (ctrl *EstadoPropiedadController) DeleteEstadoPropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid estado propiedad ID"})
		return
	}

	err = ctrl.EstadoPropiedadService.DeleteEstadoPropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete estado propiedad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado propiedad deleted successfully"})
}
