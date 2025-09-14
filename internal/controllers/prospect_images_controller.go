package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ImagenesController es el controlador para el modelo Imagenes
type ImagenesProspectoController struct {
	ImagenesProspectoService *services.ImagenesProspectoService
}

// NewImagenesController es el constructor para ImagenesController
func NewImagenesProspectoController(imagenesProspectoService *services.ImagenesProspectoService) *ImagenesProspectoController {
	return &ImagenesProspectoController{
		ImagenesProspectoService: imagenesProspectoService,
	}
}

// GET /imagenes/:id
func (ctrl *ImagenesProspectoController) GetImagen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	imagen, err := ctrl.ImagenesProspectoService.GetImagen(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la imagen"})
		return
	}

	if imagen == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Imagen no encontrada"})
		return
	}

	c.JSON(http.StatusOK, imagen)
}

func (ctrl *ImagenesProspectoController) GetImagenPrincipal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	imagen, err := ctrl.ImagenesProspectoService.GetImagenPrincipal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la imagen"})
		return
	}

	if imagen == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Imagen no encontrada"})
		return
	}

	c.JSON(http.StatusOK, imagen)
}

// GET /imagenes/propiedad/:id
func (ctrl *ImagenesProspectoController) GetImagenesByProspecto(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de propiedad inválido"})
		return
	}

	imagenes, err := ctrl.ImagenesProspectoService.GetImagenesByProspecto(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las imágenes"})
		return
	}

	if len(imagenes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron imágenes para la propiedad"})
		return
	}

	c.JSON(http.StatusOK, imagenes)
}

// POST /imagenes
func (ctrl *ImagenesProspectoController) InsertImagen(c *gin.Context) {
	var imagen models.ImagenProspecto
	if err := c.ShouldBindJSON(&imagen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	id, err := ctrl.ImagenesProspectoService.InsertImagen(&imagen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar la imagen", "details": err.Error()})
		return
	}

	log.Printf("Imagen creada con ID: %d", id)
	c.JSON(http.StatusCreated, gin.H{"id_imagen": id})
}

// PUT /imagenes/:id
func (ctrl *ImagenesProspectoController) UpdateImagen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	var imagen models.ImagenProspecto
	if err := c.ShouldBindJSON(&imagen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	if err := ctrl.ImagenesProspectoService.UpdateImagen(&imagen, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la imagen", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagen actualizada correctamente"})
}

// DELETE /imagenes/:id
func (ctrl *ImagenesProspectoController) DeleteImagen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	if err := ctrl.ImagenesProspectoService.DeleteImagen(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la imagen", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagen eliminada correctamente"})
}
