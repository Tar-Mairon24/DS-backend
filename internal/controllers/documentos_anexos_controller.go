package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DocumentosAnexosController es el controlador para el modelo Documentos_Anexos
type DocumentosAnexosController struct {
	DocumentosAnexosService *services.DocumentosAnexosService
}

// NewDocumentosAnexosController es el constructor para DocumentosAnexosController
func NewDocumentosAnexosController(service *services.DocumentosAnexosService) *DocumentosAnexosController {
	return &DocumentosAnexosController{
		DocumentosAnexosService: service,
	}
}

// GET /documentos-anexos/:id
func (ctrl *DocumentosAnexosController) GetDocumentoAnexo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento anexo inválido"})
		return
	}

	documento, err := ctrl.DocumentosAnexosService.GetDocumentoAnexo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el documento anexo"})
		return
	}

	if documento == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Documento anexo no encontrado"})
		return
	}

	c.JSON(http.StatusOK, documento)
}

// GET /documentos-anexos/propiedad/:id
func (ctrl *DocumentosAnexosController) GetDocumentosByPropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de propiedad inválido"})
		return
	}

	documentos, err := ctrl.DocumentosAnexosService.GetDocumentosByPropiedad(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los documentos anexos"})
		return
	}

	if len(documentos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron documentos anexos para la propiedad"})
		return
	}

	c.JSON(http.StatusOK, documentos)
}

// POST /documentos-anexos
func (ctrl *DocumentosAnexosController) InsertDocumentoAnexo(c *gin.Context) {
	var documento models.DocumentoAnexo
	if err := c.ShouldBindJSON(&documento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	id, err := ctrl.DocumentosAnexosService.InsertDocumentoAnexo(&documento)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar el documento anexo", "details": err.Error()})
		return
	}

	log.Printf("Documento anexo creado con ID: %d", id)
	c.JSON(http.StatusCreated, gin.H{"id_documento_anexo": id})
}

// PUT /documentos-anexos/:id
func (ctrl *DocumentosAnexosController) UpdateDocumentoAnexo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento anexo inválido"})
		return
	}

	var documento models.DocumentoAnexo
	if err := c.ShouldBindJSON(&documento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	if err := ctrl.DocumentosAnexosService.UpdateDocumentoAnexo(&documento, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el documento anexo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documento anexo actualizado correctamente"})
}

// DELETE /documentos-anexos/:id
func (ctrl *DocumentosAnexosController) DeleteDocumentoAnexo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento anexo inválido"})
		return
	}

	if err := ctrl.DocumentosAnexosService.DeleteDocumentoAnexo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el documento anexo", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documento anexo eliminado correctamente"})
}
