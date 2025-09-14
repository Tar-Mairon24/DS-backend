package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController is the controller for the User model
type Propiedad_Controller struct {
	PropiedadService       *services.PropiedadService
	EstadoPropiedadService *services.EstadoPropiedadService
}

// NewUserController is the constructor for the UserController
func NewPropiedadController(propiedadService *services.PropiedadService, estadoPropiedadService *services.EstadoPropiedadService) *Propiedad_Controller {
	return &Propiedad_Controller{
		PropiedadService:       propiedadService,
		EstadoPropiedadService: estadoPropiedadService,
	}
}

// GET /all/propiedades/
func (ctrl *Propiedad_Controller) GetAllPropiedades(c *gin.Context) {
	propiedades, err := ctrl.PropiedadService.GetAllPropiedades()
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

func (ctrl *Propiedad_Controller) GetAllPropiedadesByPrice(c *gin.Context) {
	propiedades, err := ctrl.PropiedadService.GetAllPropiedadesByPrice()
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

func (ctrl *Propiedad_Controller) GetAllPropiedadesByBedrooms(c *gin.Context) {
	propiedades, err := ctrl.PropiedadService.GetAllPropiedadesByBedrooms()
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

// GET /propiedad/:id
func (ctrl *Propiedad_Controller) GetPropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propiedad ID"})
		return
	}

	propiedad, err := ctrl.PropiedadService.GetPropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve propiedad"})
		return
	}

	if propiedad == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No propiedad found"})
		return
	}

	c.JSON(http.StatusOK, propiedad)
}

// POST /propiedad/
// POST /propiedad/
func (ctrl *Propiedad_Controller) CreatePropiedad(c *gin.Context) {
	var request struct {
		Propiedad         models.Propiedad         `json:"propiedad"`
		EstadoPropiedades models.EstadoPropiedades `json:"estado_propiedades"`
	}

	// Print the request payload
	fmt.Printf("Request payload before binding: %+v\n", request)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		print(err.Error())
		return
	}

	// Print the request payload after binding
	fmt.Printf("Request payload after binding: %+v\n", request)

	IDPropiedad, IDEstadoPropiedad, err := ctrl.PropiedadService.InsertPropiedad(&request.Propiedad, &request.EstadoPropiedades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create propiedad", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id_propiedad": IDPropiedad, "id_estado_propiedades": IDEstadoPropiedad})
}

// PUT /propiedad/:id
func (ctrl *Propiedad_Controller) UpdatePropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propiedad ID"})
		return
	}

	var propiedad models.Propiedad
	if err := c.ShouldBindJSON(&propiedad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = ctrl.PropiedadService.UpdatePropiedad(&propiedad, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update propiedad"})
		return
	}

	c.JSON(http.StatusOK, propiedad)
}

// DELETE eliminar/propiedad/:id
func (ctrl *Propiedad_Controller) DeletePropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid propiedad ID"})
		return
	}

	err = ctrl.EstadoPropiedadService.DeleteEstadoPropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete estado_propiedad"})
		return
	}

	err = ctrl.PropiedadService.DeletePropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete propiedad"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Propiedad and Estado deleted"})
}
