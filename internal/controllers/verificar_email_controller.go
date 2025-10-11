package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/models"
	"backend/internal/services"
)

type VerificarEmailController struct {
	EmailService *services.EmailService
}

func NewVerificarEmailController(emailService *services.EmailService) *VerificarEmailController {
	return &VerificarEmailController{
		EmailService: emailService,
	}
}

func (controller *VerificarEmailController) VerificarEmail(c *gin.Context) {
	verifyData := models.EmailVerification{}

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := controller.EmailService.VerifyEmail(verifyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
