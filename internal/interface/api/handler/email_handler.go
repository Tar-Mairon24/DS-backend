package handler

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"

    "ds-backend/internal/domain/models"
    "ds-backend/internal/domain/ports"
)

type EmailHandler struct {
    emailService ports.EmailService
}

func NewEmailHandler(emailService ports.EmailService) *EmailHandler {
    return &EmailHandler{
        emailService: emailService,
    }
}

func (h *EmailHandler) SendVerificationEmail(c *gin.Context) {
    var requestData struct {
        Email  string `json:"email" binding:"required,email"`
        Reason string `json:"reason" binding:"required"`
    }

    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Email and reason are required"})
        return
    }

    err := h.emailService.SendVerificationEmail(c.Request.Context(), requestData.Email, requestData.Reason)
    if err != nil {
        log.Println("Error sending verification email:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func (h *EmailHandler) VerifyEmail(c *gin.Context) {
    var verifyData models.EmailVerification

    if err := c.ShouldBindJSON(&verifyData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Code and email are required"})
        return
    }

    _, err := h.emailService.VerifyEmail(c.Request.Context(), verifyData)
    if err != nil {
        log.Println("Error verifying email:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *EmailHandler) ResendVerificationCode(c *gin.Context) {
    var resendData struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&resendData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": "Email is required"})
        return
    }

    err := h.emailService.ResendVerificationEmail(c.Request.Context(), resendData.Email)
    if err != nil {
        log.Println("Error resending verification code:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend verification code", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Verification code resent successfully"})
}