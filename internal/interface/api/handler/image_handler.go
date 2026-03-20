package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type ImageHandler struct {
	imageUseCase ports.ImageUseCase
}

func NewImageHandler(imageUseCase ports.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		imageUseCase: imageUseCase,
	}
}

func (h *ImageHandler) SaveImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve image"})
		return
	}

	var image models.Image
	if err := c.ShouldBind(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind image data"})
		return
	}

	// usecase generates path, creates directory, saves to DB
	savedImage, err := h.imageUseCase.SaveImage(c.Request.Context(), file.Filename, &image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// save file to disk — derive real path from the URL path stored in DB
	diskPath := filepath.Join("/app", savedImage.Path)
	if err = c.SaveUploadedFile(file, diskPath); err != nil {
		// roll back DB record if file write fails
		h.imageUseCase.DeleteImage(c.Request.Context(), savedImage.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file to disk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": savedImage})
}

func (h *ImageHandler) GetImageByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := h.imageUseCase.GetImageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	if image == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": image})
}

func (h *ImageHandler) GetImagesByPropertyID(c *gin.Context) {
	propertyID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	images, err := h.imageUseCase.GetImagesByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get images"})
		return
	}
	if len(images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No images found for this property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": images})
}

func (h *ImageHandler) GetMainImageByPropertyID(c *gin.Context) {
	propertyID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	image, err := h.imageUseCase.GetMainImageByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get main image"})
		return
	}
	if image == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No main image found for this property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": image})
}

func (h *ImageHandler) UpdateMainImageStatus(c *gin.Context) {
	propertyID, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}

	imageID, err := parseUintParam(c, "imageId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	if err := h.imageUseCase.UpdateMainImageStatus(c.Request.Context(), propertyID, imageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update main image status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Main image updated"})
}

func (h *ImageHandler) UpdateImage(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	existing, err := h.imageUseCase.GetImageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	var image models.Image
	if err := c.ShouldBindJSON(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	image.ID = id

	updated, err := h.imageUseCase.UpdateImage(c.Request.Context(), &image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (h *ImageHandler) DeleteImage(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := h.imageUseCase.GetImageByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image"})
		return
	}
	if image == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	if err := h.imageUseCase.DeleteImage(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	diskPath := filepath.Join("/app", image.Path)
	removeFile(diskPath)

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func removeFile(path string) {
	_ = os.Remove(path)
}