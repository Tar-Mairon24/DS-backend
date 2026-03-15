package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/services"
)

// ImagenesController es el controlador para el modelo Imagenes
type ImagenesController struct {
	ImagenesService *services.ImagenesService
}

// NewImagenesController es el constructor para ImagenesController
func NewImagenesController(imagenesService *services.ImagenesService) *ImagenesController {
	return &ImagenesController{
		ImagenesService: imagenesService,
	}
}

// GET /imagenes/:id
func (ctrl *ImagenesController) GetImagen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	imagen, err := ctrl.ImagenesService.GetImagen(id)
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

func (ctrl *ImagenesController) GetImagenPrincipal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	imagen, err := ctrl.ImagenesService.GetImagenPrincipal(id)
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
func (ctrl *ImagenesController) GetImagenesByPropiedad(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de propiedad inválido"})
		return
	}

	imagenes, err := ctrl.ImagenesService.GetImagenesByPropiedad(id)
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

func (ctrl *ImagenesController) InsertImagen(c *gin.Context) {
	// Get property ID from the URL
	idParam := c.Param("id")
	idPropiedad, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de propiedad inválido"})
		return
	}

	// Get the file from the form
	file, header, err := c.Request.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se proporcionó imagen"})
		return
	}
	defer file.Close()

	// Create folder for this property
	propertyDir := filepath.Join("/app/uploads/properties", idParam)
	if err := os.MkdirAll(propertyDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando directorio"})
		return
	}

	// Save file with a unique name
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	savePath := filepath.Join(propertyDir, filename)

	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error guardando imagen"})
		return
	}
	defer out.Close()
	io.Copy(out, file)

	// Get optional fields from form
	descripcion := c.PostForm("descripcion")
	principal := c.PostForm("principal") == "true"

	imagen := models.Imagen{
		RutaImagen:  fmt.Sprintf("/uploads/properties/%s/%s", idParam, filename),
		Descripcion: descripcion,
		Principal:   principal,
		IDPropiedad: idPropiedad,
	}

	id, err := ctrl.ImagenesService.InsertImagen(&imagen)
	if err != nil {
		// Clean up saved file if DB insert fails
		os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar la imagen"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id_imagen": id, "ruta": imagen.RutaImagen})
}

// PUT /imagenes/:id
func (ctrl *ImagenesController) UpdateImagen(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
		return
	}

	var imagen models.Imagen
	if err := c.ShouldBindJSON(&imagen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos", "details": err.Error()})
		return
	}

	if err := ctrl.ImagenesService.UpdateImagen(&imagen, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la imagen", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagen actualizada correctamente"})
}

func (ctrl *ImagenesController) DeleteImagen(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID de imagen inválido"})
        return
    }

    // Fetch the image first to get the path before deleting
    imagen, err := ctrl.ImagenesService.GetImagen(id)
    if err != nil || imagen == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Imagen no encontrada"})
        return
    }

    if err := ctrl.ImagenesService.DeleteImagen(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la imagen"})
        return
    }

    // Delete the actual file — ruta_imagen is "/uploads/properties/1/abc.jpg"
    // so we prepend /app to get the real path
    diskPath := filepath.Join("/app", imagen.RutaImagen)
    os.Remove(diskPath) // best effort, don't fail if file is already gone

    c.JSON(http.StatusOK, gin.H{"message": "Imagen eliminada correctamente"})
}