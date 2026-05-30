package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/interface/api/handler"
	serviceMocks "ds-backend/test/mocks/service"
)

// TestSendVerificationEmail_Success tests successful email verification sending
func TestSendVerificationEmail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("SendVerificationEmail", mock.Anything, "test@example.com", "registration").Return(nil)

	body := []byte(`{"email":"test@example.com","reason":"registration"}`)
	req, _ := http.NewRequest(http.MethodPost, "/send-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SendVerificationEmail(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Verification email sent successfully"`)
	mockEmailService.AssertExpectations(t)
}

// TestSendVerificationEmail_InvalidEmail tests sending with invalid email format
func TestSendVerificationEmail_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	body := []byte(`{"email":"invalid-email","reason":"registration"}`)
	req, _ := http.NewRequest(http.MethodPost, "/send-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SendVerificationEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Email and reason are required"`)
}

// TestSendVerificationEmail_MissingEmail tests sending without email field
func TestSendVerificationEmail_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	body := []byte(`{"reason":"registration"}`)
	req, _ := http.NewRequest(http.MethodPost, "/send-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SendVerificationEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
}

// TestSendVerificationEmail_MissingReason tests sending without reason field
func TestSendVerificationEmail_MissingReason(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/send-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SendVerificationEmail(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
}

// TestSendVerificationEmail_ServiceError tests when email service fails
func TestSendVerificationEmail_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("SendVerificationEmail", mock.Anything, "test@example.com", "registration").
		Return(errors.New("SMTP connection failed"))

	body := []byte(`{"email":"test@example.com","reason":"registration"}`)
	req, _ := http.NewRequest(http.MethodPost, "/send-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.SendVerificationEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to send verification email"`)
	assert.Contains(t, w.Body.String(), `"message":"SMTP connection failed"`)
	mockEmailService.AssertExpectations(t)
}

// TestVerifyEmail_Success tests successful email verification
func TestVerifyEmail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	mockEmailService.On("VerifyEmail", mock.Anything, verifyData).Return(true, nil)

	body := []byte(`{"code":"123456","email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.VerifyEmail(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Email verified successfully"`)
	mockEmailService.AssertExpectations(t)
}

// TestVerifyEmail_InvalidCode tests verification with invalid code
func TestVerifyEmail_InvalidCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	verifyData := models.EmailVerification{
		Code:  "invalid",
		Email: "test@example.com",
	}
	mockEmailService.On("VerifyEmail", mock.Anything, verifyData).
		Return(false, errors.New("invalid verification code or email"))

	body := []byte(`{"code":"invalid","email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.VerifyEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to verify email"`)
	assert.Contains(t, w.Body.String(), `"message":"invalid verification code or email"`)
	mockEmailService.AssertExpectations(t)
}

// TestVerifyEmail_CodeAlreadyUsed tests verification with already used code
func TestVerifyEmail_CodeAlreadyUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "test@example.com",
	}
	mockEmailService.On("VerifyEmail", mock.Anything, verifyData).
		Return(false, errors.New("verification code has already been used"))

	body := []byte(`{"code":"123456","email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.VerifyEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to verify email"`)
	assert.Contains(t, w.Body.String(), `"message":"verification code has already been used"`)
	mockEmailService.AssertExpectations(t)
}

// TestVerifyEmail_MissingCode tests verification without code field
func TestVerifyEmail_MissingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	// When code is missing, it will be empty string, so service will be called
	verifyData := models.EmailVerification{
		Code:  "",
		Email: "test@example.com",
	}
	mockEmailService.On("VerifyEmail", mock.Anything, verifyData).
		Return(false, errors.New("code is required"))

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.VerifyEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to verify email"`)
	assert.Contains(t, w.Body.String(), `"message":"code is required"`)
	mockEmailService.AssertExpectations(t)
}

// TestVerifyEmail_MissingEmail tests verification without email field
func TestVerifyEmail_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	// When email is missing, it will be empty string, so service will be called
	verifyData := models.EmailVerification{
		Code:  "123456",
		Email: "",
	}
	mockEmailService.On("VerifyEmail", mock.Anything, verifyData).
		Return(false, errors.New("email is required"))

	body := []byte(`{"code":"123456"}`)
	req, _ := http.NewRequest(http.MethodPost, "/verify-email", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.VerifyEmail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to verify email"`)
	assert.Contains(t, w.Body.String(), `"message":"email is required"`)
	mockEmailService.AssertExpectations(t)
}

// TestResendVerificationCode_Success tests successful resend of verification code
func TestResendVerificationCode_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("ResendVerificationEmail", mock.Anything, "test@example.com").Return(nil)

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Verification code resent successfully"`)
	mockEmailService.AssertExpectations(t)
}

// TestResendVerificationCode_InvalidEmail tests resend with invalid email format
func TestResendVerificationCode_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	body := []byte(`{"email":"invalid-email"}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Email is required"`)
}

// TestResendVerificationCode_MissingEmail tests resend without email field
func TestResendVerificationCode_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	body := []byte(`{}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Email is required"`)
}

// TestResendVerificationCode_UserNotFound tests resend when user not found
func TestResendVerificationCode_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("ResendVerificationEmail", mock.Anything, "nonexistent@example.com").
		Return(errors.New("user not found"))

	body := []byte(`{"email":"nonexistent@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to resend verification code"`)
	assert.Contains(t, w.Body.String(), `"message":"user not found"`)
	mockEmailService.AssertExpectations(t)
}

// TestResendVerificationCode_MaxResendReached tests resend when max resends reached
func TestResendVerificationCode_MaxResendReached(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("ResendVerificationEmail", mock.Anything, "test@example.com").
		Return(errors.New("maximum number of resends reached"))

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to resend verification code"`)
	assert.Contains(t, w.Body.String(), `"message":"maximum number of resends reached"`)
	mockEmailService.AssertExpectations(t)
}

// TestResendVerificationCode_ServiceError tests resend when email service fails
func TestResendVerificationCode_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockEmailService := serviceMocks.NewMockEmailService()
	h := handler.NewEmailHandler(mockEmailService)

	mockEmailService.On("ResendVerificationEmail", mock.Anything, "test@example.com").
		Return(errors.New("SMTP connection failed"))

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/resend-verification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ResendVerificationCode(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to resend verification code"`)
	assert.Contains(t, w.Body.String(), `"message":"SMTP connection failed"`)
	mockEmailService.AssertExpectations(t)
}
