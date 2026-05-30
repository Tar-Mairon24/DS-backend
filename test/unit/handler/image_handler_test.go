package handler_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/interface/api/handler"
	usecaseMocks "ds-backend/test/mocks/usecase"
)

// Helper function to create a temporary directory for tests
func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "image_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tmpDir
}

// Helper function to clean up test directory
func cleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Failed to clean up temp directory: %v", err)
	}
}

// TestUploadImage_Success tests successful image upload
func TestUploadImage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	diskPath := filepath.Join(tmpDir, "properties", "1", "uuid.jpg")
	urlPath := "/uploads/properties/1/uuid.jpg"

	mockUC.On("GeneratePath", "test.jpg", uint(1)).Return(diskPath, urlPath, nil)
	mockUC.On("SaveImage", mock.Anything, mock.MatchedBy(func(img *models.Image) bool {
		return img.PropertyID == 1 && img.Path == urlPath
	})).Return(&models.Image{
		ID:         1,
		PropertyID: 1,
		Path:       urlPath,
		CreatedAt:  time.Now(),
	}, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.WriteField("description", "Test image")
	writer.WriteField("main_image", "false")
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/1/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

// TestUploadImage_NoFile tests upload with no file provided
func TestUploadImage_NoFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create request without file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "Test image")
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/1/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUploadImage_InvalidFileType tests upload with invalid file type
func TestUploadImage_InvalidFileType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations - GeneratePath should fail for invalid type
	mockUC.On("GeneratePath", "test.txt", uint(1)).Return("", "", errors.New("unsupported file type"))

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form data with invalid file type
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.txt")
	part.Write([]byte("fake text data"))
	writer.WriteField("description", "Test image")
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/1/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertExpectations(t)
}

// TestUploadImage_FileTooLarge tests upload with file exceeding size limit
func TestUploadImage_FileTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations - GeneratePath should fail for large file
	mockUC.On("GeneratePath", "large.jpg", uint(1)).Return("", "", errors.New("file too large"))

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form data with large file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "large.jpg")
	// Write large data (simulating large file)
	largeData := make([]byte, 6*1024*1024) // 6MB
	part.Write(largeData)
	writer.WriteField("description", "Large image")
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/1/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUC.AssertExpectations(t)
}

// TestUploadImage_ServiceFails tests upload when service fails
func TestUploadImage_ServiceFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	diskPath := filepath.Join(tmpDir, "properties", "1", "uuid.jpg")
	urlPath := "/uploads/properties/1/uuid.jpg"

	mockUC.On("GeneratePath", "test.jpg", uint(1)).Return(diskPath, urlPath, nil)
	mockUC.On("SaveImage", mock.Anything, mock.MatchedBy(func(img *models.Image) bool {
		return img.PropertyID == 1
	})).Return(nil, errors.New("database error"))

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.WriteField("description", "Test image")
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/1/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

// TestDeleteImage_Success tests successful image deletion
func TestDeleteImage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImageByID", mock.Anything, uint(1)).Return(&models.Image{
		ID:         1,
		PropertyID: 1,
		Path:       "/uploads/properties/1/uuid.jpg",
		CreatedAt:  time.Now(),
	}, nil)
	mockUC.On("DeleteImage", mock.Anything, uint(1)).Return(nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteImage(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// TestDeleteImage_NotFound tests deletion of non-existent image
func TestDeleteImage_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImageByID", mock.Anything, uint(999)).Return(nil, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/images/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.DeleteImage(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

// TestDeleteImage_ServiceFails tests deletion when service fails
func TestDeleteImage_ServiceFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImageByID", mock.Anything, uint(1)).Return(&models.Image{
		ID:         1,
		PropertyID: 1,
		Path:       "/uploads/properties/1/uuid.jpg",
		CreatedAt:  time.Now(),
	}, nil)
	mockUC.On("DeleteImage", mock.Anything, uint(1)).Return(errors.New("database error"))

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteImage(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

// TestGetImage_Success tests successful image retrieval
func TestGetImage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImageByID", mock.Anything, uint(1)).Return(&models.Image{
		ID:          1,
		PropertyID:  1,
		Path:        "/uploads/properties/1/uuid.jpg",
		Description: "Test image",
		CreatedAt:   time.Now(),
	}, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/images/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetImageByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// TestGetImage_NotFound tests retrieval of non-existent image
func TestGetImage_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImageByID", mock.Anything, uint(999)).Return(nil, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/images/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetImageByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

// TestGetImagesByProperty_Success tests successful retrieval of property images
func TestGetImagesByProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	images := []models.Image{
		{
			ID:         1,
			PropertyID: 1,
			Path:       "/uploads/properties/1/uuid1.jpg",
			CreatedAt:  time.Now(),
		},
		{
			ID:         2,
			PropertyID: 1,
			Path:       "/uploads/properties/1/uuid2.jpg",
			CreatedAt:  time.Now(),
		},
	}
	mockUC.On("GetImagesByPropertyID", mock.Anything, uint(1)).Return(images, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/1/images", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetImagesByPropertyID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// TestGetImagesByProperty_NoImages tests retrieval when property has no images
func TestGetImagesByProperty_NoImages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations - empty list
	mockUC.On("GetImagesByPropertyID", mock.Anything, uint(999)).Return([]models.Image{}, nil)

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/999/images", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetImagesByPropertyID(c)

	// According to the handler implementation, empty list returns 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

// TestGetImagesByProperty_ServiceError tests retrieval when service fails
func TestGetImagesByProperty_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	// Setup expectations
	mockUC.On("GetImagesByPropertyID", mock.Anything, uint(1)).Return([]models.Image{}, errors.New("database error"))

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/1/images", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetImagesByPropertyID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

// TestUploadImage_InvalidPropertyID tests upload with invalid property ID
func TestUploadImage_InvalidPropertyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	c.Request, _ = http.NewRequest("POST", "/properties/abc/images", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.SaveImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteImage_InvalidImageID tests deletion with invalid image ID
func TestDeleteImage_InvalidImageID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/images/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteImage(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetImage_InvalidImageID tests retrieval with invalid image ID
func TestGetImage_InvalidImageID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/images/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetImageByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetImagesByProperty_InvalidPropertyID tests retrieval with invalid property ID
func TestGetImagesByProperty_InvalidPropertyID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockImageUseCase()

	h := handler.NewImageHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/properties/abc/images", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetImagesByPropertyID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
