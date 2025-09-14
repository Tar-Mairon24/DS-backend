package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContratosController struct {
	Service *services.ContratosService
}

// Constructor para ContratosController
func NewContratosController(service *services.ContratosService) *ContratosController {
	return &ContratosController{
		Service: service,
	}
}

// Recupera un contrato por ID
func (controller *ContratosController) GetContrato(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	contrato, err := controller.Service.GetContrato(id)
	if err != nil {
		log.Println("Error recuperando contrato:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	if contrato == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contrato no encontrado"})
		return
	}

	c.JSON(http.StatusOK, contrato)
}

func (controller *ContratosController) GetContratos(c *gin.Context) {
	contratos, err := controller.Service.GetContratos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contratos"})
		return
	}

	if contratos == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No contratos found"})
		return
	}

	c.JSON(http.StatusOK, contratos)
}

// Recupera todos los contratos asociados a una propiedad
func (controller *ContratosController) GetContratosByPropiedad(c *gin.Context) {
	idPropiedad, err := strconv.Atoi(c.Param("id_propiedad"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de propiedad inválido"})
		return
	}

	contratos, err := controller.Service.GetContratosByPropiedad(idPropiedad)
	if err != nil {
		log.Println("Error recuperando contratos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, contratos)
}

// Inserta un nuevo contrato
func (controller *ContratosController) CreateContrato(c *gin.Context) {
	var contrato models.Contrato
	if err := c.ShouldBindJSON(&contrato); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	id, err := controller.Service.InsertContrato(&contrato)
	if err != nil {
		log.Println("Error insertando contrato:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id_contrato": id})
}

// Actualiza un contrato existente
func (controller *ContratosController) UpdateContrato(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var contrato models.Contrato
	if err := c.ShouldBindJSON(&contrato); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	err = controller.Service.UpdateContrato(&contrato, id)
	if err != nil {
		log.Println("Error actualizando contrato:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, contrato)
}

// Elimina un contrato por ID
func (controller *ContratosController) DeleteContrato(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = controller.Service.DeleteContrato(id)
	if err != nil {
		log.Println("Error eliminando contrato:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
