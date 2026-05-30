package handler_test

import (
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
	usecaseMocks "ds-backend/test/mocks/usecase"
)

func TestGetProperties_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	properties := []models.PropertyCard{
		{ID: 1, Title: "Prop1"},
		{ID: 2, Title: "Prop2"},
	}
	mockUC.On("GetAllProperties", mock.Anything).Return(properties, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProperties_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetAllProperties", mock.Anything).Return([]models.PropertyCard{}, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProperties_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetAllProperties", mock.Anything).Return([]models.PropertyCard{}, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	// Bug Condition B4: Empty collections should return 200 OK, not 404
	// This test verifies the FIXED behavior
	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	expected := &models.PropertyResponse{ID: 1, Title: "Prop1"}
	mockUC.On("GetPropertyByID", mock.Anything, uint(1)).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPropertyByID_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/-5", nil)
	c.Params = gin.Params{{Key: "id", Value: "-5"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPropertyByID_ErrorFromUsecase(t *testing.T) {
	var property *models.PropertyResponse = nil

	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetPropertyByID", mock.Anything, uint(2)).Return(property, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_NotFound(t *testing.T) {
	var property *models.PropertyResponse = nil
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("GetPropertyByID", mock.Anything, uint(3)).Return(property, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/3", nil)
	c.Params = gin.Params{{Key: "id", Value: "3"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}
func TestCreateProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	expected := &models.PropertyResponse{ID: 1, Title: "New Property"}
	mockUC.On("CreateProperty", mock.Anything, mock.Anything).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"New Property","owner_id":1}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateProperty_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProperty_UsecaseError(t *testing.T) {
	var property *models.PropertyResponse = nil
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("CreateProperty", mock.Anything, mock.Anything).Return(property, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"New Property","owner_id":1}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("DeleteProperty", mock.Anything, uint(1)).Return(nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/properties/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteProperty_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/properties/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteProperty_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/properties/-10", nil)
	c.Params = gin.Params{{Key: "id", Value: "-10"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteProperty_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	mockUC.On("DeleteProperty", mock.Anything, uint(2)).Return(errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/properties/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	expected := &models.PropertyResponse{ID: 1, Title: "Updated Property"}
	mockUC.On("UpdateProperty", mock.Anything, mock.Anything).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Updated Property"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateProperty_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/properties/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProperty_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/properties/-2", nil)
	c.Params = gin.Params{{Key: "id", Value: "-2"}}

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProperty_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProperty_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	var propertyResp *models.PropertyResponse = nil
	mockUC.On("UpdateProperty", mock.Anything, mock.Anything).Return(propertyResp, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	c.Request, _ = http.NewRequest("PUT", "/properties/2", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Set("user_id", uint(1))
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"Some Title"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

// Preservation Property Tests for B4
// **Validates: Requirements 3.1, 3.2**
// These tests verify that non-empty collections continue to return 200 OK with all records
// and that response structure is consistent across all collection sizes

func TestGetProperties_WithOneProperty_ShouldReturn200OK(t *testing.T) {
	// Preservation: GET /api/v1/properties with 1 property should return 200 OK
	// Expected: 200 OK with array containing 1 property
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	properties := []models.PropertyCard{
		{ID: 1, Title: "Property 1"},
	}
	mockUC.On("GetAllProperties", mock.Anything).Return(properties, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusOK, w.Code, "Non-empty properties should return 200 OK")
	assert.Contains(t, w.Body.String(), "Property 1", "Response should contain property data")
	mockUC.AssertExpectations(t)
}

func TestGetProperties_WithMultipleProperties_ShouldReturn200OK(t *testing.T) {
	// Preservation: GET /api/v1/properties with N properties should return 200 OK with all N records
	// Expected: 200 OK with array containing all properties
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	properties := []models.PropertyCard{
		{ID: 1, Title: "Property 1"},
		{ID: 2, Title: "Property 2"},
		{ID: 3, Title: "Property 3"},
		{ID: 4, Title: "Property 4"},
		{ID: 5, Title: "Property 5"},
	}
	mockUC.On("GetAllProperties", mock.Anything).Return(properties, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusOK, w.Code, "Non-empty properties should return 200 OK")
	// Verify all properties are in response
	for i := 1; i <= 5; i++ {
		assert.Contains(t, w.Body.String(), "Property "+string(rune(48+i)), "Response should contain all properties")
	}
	mockUC.AssertExpectations(t)
}

func TestGetProperties_ResponseStructureConsistency(t *testing.T) {
	// Preservation: Response structure should be consistent for all collection sizes
	// Expected: Same response format for 1 property and 3 properties
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockPropertyUseCase()
	properties := []models.PropertyCard{
		{ID: 1, Title: "Property 1"},
		{ID: 2, Title: "Property 2"},
		{ID: 3, Title: "Property 3"},
	}
	mockUC.On("GetAllProperties", mock.Anything).Return(properties, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties", nil)

	h.GetProperties(c)

	assert.Equal(t, http.StatusOK, w.Code)
	// Response should be a JSON array
	assert.True(t, w.Body.String()[0] == '[', "Response should be a JSON array")
	mockUC.AssertExpectations(t)
}
