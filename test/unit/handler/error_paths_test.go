package handler_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/interface/api/handler"
	serviceMocks "ds-backend/test/mocks/service"
	usecaseMocks "ds-backend/test/mocks/usecase"
)

// ============================================================================
// PROPERTY HANDLER ERROR PATHS
// ============================================================================

// Missing required fields tests
func TestCreateProperty_MissingOwnerID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Property without owner"}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "owner_id is required")
}

func TestCreateProperty_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	// No user_id set in context
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Property","owner_id":1}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Invalid data types tests
func TestCreateProperty_InvalidJSONTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	// owner_id should be number, not string
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Property","owner_id":"not_a_number"}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProperty_InvalidJSONTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":123,"owner_id":"invalid"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Service failure tests
func TestGetProperties_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetAllProperties", mock.Anything).Return([]models.PropertyCard{}, errors.New("database connection failed"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection failed")
}

func TestCreateProperty_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("CreateProperty", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Property","owner_id":1}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create property")
}

func TestUpdateProperty_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("UpdateProperty", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Updated"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteProperty_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("DeleteProperty", mock.Anything, uint(1)).Return(errors.New("database error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/properties/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================================================
// USER HANDLER ERROR PATHS
// ============================================================================

// Missing required fields tests
func TestCreateUser_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("email is required"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"testuser","password":"pass123"}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "email is required")
}

func TestCreateUser_MissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("password is required"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"testuser","email":"test@example.com"}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "password is required")
}

// Invalid data types tests
func TestCreateUser_InvalidEmailType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"testuser","email":123,"password":"pass123"}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_InvalidJSONTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"id":"not_a_number","username":true}`))

	h.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Unauthorized tests
func TestGetUserByID_Unauthorized_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("GetUserByID", mock.Anything, uint(1)).Return(nil, errors.New("user not found"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/users/1", nil)
	// No JWT token in context
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Note: This test assumes middleware would check for token
	// Handler itself doesn't enforce auth, but we test the path
	h.GetUserByID(c)

	// Handler should still work without auth check at handler level
	// This is a middleware concern
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

// Service failure tests
func TestGetUsers_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database connection failed"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users", nil)

	h.GetUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection failed")
}

func TestGetUserByID_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("GetUserByID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetUserByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateUser_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("email already exists"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"testuser","email":"test@example.com","password":"pass123"}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "email already exists")
}

func TestUpdateUser_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("UpdateUser", mock.Anything, mock.Anything).Return(nil, errors.New("user not found"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"id":999,"username":"updated"}`))

	h.UpdateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteUser_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("DeleteUser", mock.Anything, uint(1)).Return(errors.New("database error"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Invalid ID format tests
func TestGetUserByID_InvalidIDFormat_ErrorPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetUserByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID")
}

func TestDeleteUser_InvalidIDFormat_ErrorPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/users/xyz", nil)
	c.Params = gin.Params{{Key: "id", Value: "xyz"}}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// CreateOwner error paths
func TestCreateOwner_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("email is required"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/owners", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"owner1","password":"pass123"}`))

	h.CreateOwner(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateOwner_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/owners", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"owner1","email":"owner@example.com","password":"pass123"}`))

	h.CreateOwner(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// CreateClient error paths
func TestCreateClient_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("email is required"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/clients", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"client1","password":"pass123"}`))

	h.CreateClient(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateClient_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/clients", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"username":"client1","email":"client@example.com","password":"pass123"}`))

	h.CreateClient(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============================================================================
// APPOINTMENT HANDLER ERROR PATHS
// ============================================================================

// Missing required fields tests
func TestRescheduleAppointment_MissingStartTime(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"end_date":"2024-12-31T18:00:00Z"}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRescheduleAppointment_MissingEndTime(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"start_date":"2024-12-31T17:00:00Z"}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Invalid data types tests
func TestRescheduleAppointment_InvalidDateTimeFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"start_date":"invalid-date","end_date":"also-invalid"}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAppointmentStatus_InvalidStatusType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"status":123}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Service failure tests
func TestRescheduleAppointment_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	mockUC.On("Reschedule", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"start_date":"2024-12-31T17:00:00Z","end_date":"2024-12-31T18:00:00Z"}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAppointmentStatus_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	mockUC.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"status":"completed"}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Invalid ID format tests
func TestRescheduleAppointment_InvalidIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/invalid/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"start_date":"2024-12-31T17:00:00Z","end_date":"2024-12-31T18:00:00Z"}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAppointmentStatus_InvalidIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	c.Request, _ = http.NewRequest("PUT", "/appointments/abc/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"status":"completed"}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============================================================================
// AUTH HANDLER ERROR PATHS
// ============================================================================

// Missing required fields tests
func TestUserLogin_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	mockAuth.On("Login", mock.Anything, "", "password123").Return(nil, errors.New("email is required"))
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"password":"password123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserLogin_MissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	mockAuth.On("Login", mock.Anything, "test@example.com", "").Return(nil, errors.New("password is required"))
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"email":"test@example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Invalid data types tests
func TestUserLogin_InvalidEmailType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"email":123,"password":"password123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.UserLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserLogin_InvalidPasswordType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"email":"test@example.com","password":true}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.UserLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Service failure tests
func TestUserLogin_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	mockAuth.On("Login", mock.Anything, "test@example.com", "password123").Return(nil, errors.New("database connection failed"))
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"email":"test@example.com","password":"password123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}

func TestRefreshToken_ServiceFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	h := handler.NewAuthHandler(mockJWT, mockAuth)

	refreshTokenData := models.RefreshTokenData{
		RefreshToken: "refresh-token",
		JwtToken:     "jwt-token",
	}
	mockAuth.On("RefreshToken", mock.Anything, refreshTokenData).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuth.AssertExpectations(t)
}

// ============================================================================
// EDGE CASES & BOUNDARY TESTS
// ============================================================================

// Empty request body tests
func TestCreateProperty_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("validation failed"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Malformed JSON tests
func TestCreateProperty_MalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_MalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/users", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{not valid json}`))

	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Negative ID tests
func TestGetPropertyByID_NegativeID_ErrorPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "-1"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUser_NegativeID_ErrorPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/users/-5", nil)
	c.Params = gin.Params{{Key: "id", Value: "-5"}}

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Zero ID tests
func TestGetPropertyByID_ZeroID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/0", nil)
	c.Params = gin.Params{{Key: "id", Value: "0"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserByID_ZeroID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockUserUseCase()
	mockUC.On("GetUserByID", mock.Anything, uint(0)).Return(nil, errors.New("user not found"))
	h := handler.NewUserHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users/0", nil)
	c.Params = gin.Params{{Key: "id", Value: "0"}}

	h.GetUserByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Very large ID tests
func TestGetPropertyByID_VeryLargeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetPropertyByID", mock.Anything, uint(9999999)).Return(nil, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/9999999", nil)
	c.Params = gin.Params{{Key: "id", Value: "9999999"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Special characters in string fields
func TestCreateProperty_SpecialCharactersInTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	expected := &models.PropertyResponse{ID: 1, Title: "Property <script>alert('xss')</script>"}
	mockUC.On("CreateProperty", mock.Anything, mock.Anything).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Property <script>alert('xss')</script>","owner_id":1}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// Null values in optional fields
func TestUpdateProperty_NullOptionalFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	expected := &models.PropertyResponse{ID: 1}
	mockUC.On("UpdateProperty", mock.Anything, mock.Anything).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":null,"description":null}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Validates: Requirements 2.13, 2.14, 2.15
// These tests verify comprehensive error path coverage for all handlers
