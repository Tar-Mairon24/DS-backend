package controllers

import (
	"log"
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

func (controller *VerificarEmailController) EnviarEmailVerificacion(c *gin.Context) {
	requestData := models.EmailRequest{}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	err := controller.EmailService.SendVerificationEmail(requestData.Email, requestData.Reason)
	if err != nil {
		log.Println("Error sending verification email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to send verification email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Verification email sent successfully"})
}

func (controller *VerificarEmailController) VerificarEmail(c *gin.Context) {
	verifyData := models.EmailVerification{}

	if err := c.ShouldBindJSON(&verifyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	_, err := controller.EmailService.VerifyEmail(verifyData)
	if err != nil {
		log.Println("Error verifying email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to verify email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email verified successfully"})
}

func (controller *VerificarEmailController) ReenviarCodigoVerificacion(c *gin.Context) {
	resendData := models.EmailResendRequest{}

	if err := c.ShouldBindJSON(&resendData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}

	err := controller.EmailService.ResendVerificationEmail(resendData.Email)
	if err != nil {
		log.Println("Error resending verification code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to resend verification code", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Verification code resent successfully"})
}
