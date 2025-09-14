package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController is the controller for the User model
type ProspectoController struct {
	ProspectoService *services.ProspectoService
}

// NewUserController is the constructor for the UserController
func NewProspectoController(prospectoService *services.ProspectoService) *ProspectoController {
	return &ProspectoController{
		ProspectoService: prospectoService,
	}
}

func (ctrl *ProspectoController) GetProspecto(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prospecto ID"})
		return
	}

	prospecto, err := ctrl.ProspectoService.GetProspecto(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve prospecto"})
		return
	}

	if prospecto == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No prospecto found"})
		return
	}

	c.JSON(http.StatusOK, prospecto)
}

// POST /citas
func (ctrl *ProspectoController) InsertProspecto(c *gin.Context) {
	var prospecto models.Prospecto

	// fmt.Printf("Request payload before binding: %+v\n", request)
	if err := c.ShouldBindJSON(&prospecto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		print(err.Error())
		return
	}

	// fmt.Printf("Request payload after binding: %+v\n", request)

	id, err := ctrl.ProspectoService.InsertProspecto(&prospecto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create propiedad", "details": err.Error()})
		return
	}
	log.Printf("Created prospecto with ID: %d", id)

	c.JSON(http.StatusCreated, prospecto)
}

func (ctrl *ProspectoController) UpdateProspecto(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prospecto ID"})
		return
	}

	var prospecto models.Prospecto
	// fmt.Printf("Request payload before binding: %+v\n", request)
	if err := c.ShouldBindJSON(&prospecto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		print(err.Error())
		return
	}

	// fmt.Printf("Request payload after binding: %+v\n", request)

	err = ctrl.ProspectoService.UpdateProspecto(&prospecto, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create propiedad", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prospecto)
}
