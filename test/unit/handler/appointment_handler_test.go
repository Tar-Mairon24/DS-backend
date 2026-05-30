package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/interface/api/handler"
	usecaseMocks "ds-backend/test/mocks/usecase"
)

// Bug Condition B4 Exploration Tests
// These tests MUST FAIL on unfixed code to confirm the bug exists
// **Validates: Requirements 2.1, 2.2**

func TestGetAll_EmptyAppointments_ShouldReturn200OK(t *testing.T) {
	// Bug Condition: GET /api/v1/appointments on empty database
	// Expected: 200 OK with empty array
	// Current (buggy): 404 Not Found
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	// Repository returns empty slice (no appointments in DB)
	mockUC.On("GetAll", mock.Anything).Return([]models.AppointmentCalendarView{}, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	// EXPECTED: 200 OK with empty array
	assert.Equal(t, http.StatusOK, w.Code, "Empty appointments should return 200 OK, not 404")

	// Verify response contains empty data array
	assert.Contains(t, w.Body.String(), "\"data\":[]", "Response should contain empty data array")
	mockUC.AssertExpectations(t)
}

func TestGetAll_WithAppointments_ShouldReturn200OK(t *testing.T) {
	// Preservation: GET /api/v1/appointments with data should continue to work
	// Expected: 200 OK with appointments array
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
		{
			ID:       2,
			Title:    "Appointment 2",
			Property: models.PropertyTile{ID: 2, Title: "Property 2"},
		},
	}
	mockUC.On("GetAll", mock.Anything).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	// Should return 200 OK with appointments
	assert.Equal(t, http.StatusOK, w.Code, "Non-empty appointments should return 200 OK")
	assert.Contains(t, w.Body.String(), "\"count\":2", "Response should contain count of 2")
	mockUC.AssertExpectations(t)
}

func TestGetAll_RepositoryError_ShouldReturn500(t *testing.T) {
	// Error handling: Repository error should return 500
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("GetAll", mock.Anything).Return([]models.AppointmentCalendarView{}, errors.New("database error"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	expected := &models.AppointmentDetail{
		ID:    1,
		Title: "Test Appointment",
	}
	mockUC.On("GetByID", mock.Anything, uint(1)).Return(expected, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/appointments/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/appointments/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/appointments/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	expected := &models.Appointment{
		ID:    1,
		Title: "New Appointment",
	}
	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*models.AppointmentRequest")).Return(expected, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/appointments", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"title": "New Appointment",
		"description": "Test",
		"start_date": "2024-01-15T10:00:00Z",
		"end_date": "2024-01-15T11:00:00Z",
		"status": "scheduled",
		"client_id": 1,
		"property_id": 1
	}`))

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreate_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()
	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/appointments", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	expected := &models.Appointment{
		ID:    1,
		Title: "Updated Appointment",
	}
	mockUC.On("Update", mock.Anything, uint(1), mock.AnythingOfType("*models.AppointmentRequest")).Return(expected, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"title": "Updated Appointment",
		"description": "Test",
		"start_date": "2024-01-15T10:00:00Z",
		"end_date": "2024-01-15T11:00:00Z",
		"status": "scheduled",
		"client_id": 1,
		"property_id": 1
	}`))

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("Delete", mock.Anything, uint(1)).Return(nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/appointments/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetCalendar_Month_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
	}
	mockUC.On("GetByMonth", mock.Anything, 2024, 1).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments/calendar?view=month&year=2024&month=1", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.GetCalendar(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetCalendar_Day_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
	}
	mockUC.On("GetByDay", mock.Anything, mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == 2024 && t.Month() == 1 && t.Day() == 15
	})).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments/calendar?view=day&date=2024-01-15", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.GetCalendar(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

// Preservation Property Tests for B4
// **Validates: Requirements 3.1, 3.2**
// These tests verify that non-empty collections continue to return 200 OK with all records
// and that response structure is consistent across all collection sizes

func TestGetAll_WithOneAppointment_ShouldReturn200OK(t *testing.T) {
	// Preservation: GET /api/v1/appointments with 1 appointment should return 200 OK
	// Expected: 200 OK with array containing 1 appointment and count=1
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
	}
	mockUC.On("GetAll", mock.Anything).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code, "Non-empty appointments should return 200 OK")
	assert.Contains(t, w.Body.String(), "\"count\":1", "Response should contain count of 1")
	assert.Contains(t, w.Body.String(), "Appointment 1", "Response should contain appointment data")
	mockUC.AssertExpectations(t)
}

func TestGetAll_WithMultipleAppointments_ShouldReturn200OK(t *testing.T) {
	// Preservation: GET /api/v1/appointments with N appointments should return 200 OK with all N records
	// Expected: 200 OK with array containing all appointments and correct count
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
		{
			ID:       2,
			Title:    "Appointment 2",
			Property: models.PropertyTile{ID: 2, Title: "Property 2"},
		},
		{
			ID:       3,
			Title:    "Appointment 3",
			Property: models.PropertyTile{ID: 3, Title: "Property 3"},
		},
		{
			ID:       4,
			Title:    "Appointment 4",
			Property: models.PropertyTile{ID: 4, Title: "Property 4"},
		},
		{
			ID:       5,
			Title:    "Appointment 5",
			Property: models.PropertyTile{ID: 5, Title: "Property 5"},
		},
	}
	mockUC.On("GetAll", mock.Anything).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code, "Non-empty appointments should return 200 OK")
	assert.Contains(t, w.Body.String(), "\"count\":5", "Response should contain count of 5")
	// Verify all appointments are in response
	for i := 1; i <= 5; i++ {
		assert.Contains(t, w.Body.String(), "Appointment "+string(rune(48+i)), "Response should contain all appointments")
	}
	mockUC.AssertExpectations(t)
}

func TestGetAll_ResponseStructureConsistency(t *testing.T) {
	// Preservation: Response structure should be consistent for all collection sizes
	// Expected: Same response format (data, count, message fields) for all sizes
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	appointments := []models.AppointmentCalendarView{
		{
			ID:       1,
			Title:    "Appointment 1",
			Property: models.PropertyTile{ID: 1, Title: "Property 1"},
		},
		{
			ID:       2,
			Title:    "Appointment 2",
			Property: models.PropertyTile{ID: 2, Title: "Property 2"},
		},
		{
			ID:       3,
			Title:    "Appointment 3",
			Property: models.PropertyTile{ID: 3, Title: "Property 3"},
		},
	}
	mockUC.On("GetAll", mock.Anything).Return(appointments, nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/appointments", nil)

	h.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	// Response should contain required fields
	assert.Contains(t, w.Body.String(), "\"data\":", "Response should contain data field")
	assert.Contains(t, w.Body.String(), "\"count\":", "Response should contain count field")
	assert.Contains(t, w.Body.String(), "\"message\":", "Response should contain message field")
	mockUC.AssertExpectations(t)
}

// Reschedule Handler Tests
// **Validates: Requirements 2.7, 2.8**

func TestReschedule_Success(t *testing.T) {
	// Success: Valid reschedule data, returns 200 OK
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	newStart := time.Now().Add(24 * time.Hour)
	newEnd := newStart.Add(1 * time.Hour)

	mockUC.On("Reschedule", mock.Anything, uint(1), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newStart.Year() && t.Month() == newStart.Month() && t.Day() == newStart.Day()
	}), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newEnd.Year() && t.Month() == newEnd.Month() && t.Day() == newEnd.Day()
	})).Return(nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "` + newStart.Format(time.RFC3339) + `",
		"end_date": "` + newEnd.Format(time.RFC3339) + `"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Appointment rescheduled successfully")
	mockUC.AssertExpectations(t)
}

func TestReschedule_NotFound(t *testing.T) {
	// Error: Appointment not found, returns 404 Not Found
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	newStart := time.Now().Add(24 * time.Hour)
	newEnd := newStart.Add(1 * time.Hour)

	mockUC.On("Reschedule", mock.Anything, uint(999), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newStart.Year() && t.Month() == newStart.Month() && t.Day() == newStart.Day()
	}), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newEnd.Year() && t.Month() == newEnd.Month() && t.Day() == newEnd.Day()
	})).Return(errors.New("appointment not found"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/999/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "` + newStart.Format(time.RFC3339) + `",
		"end_date": "` + newEnd.Format(time.RFC3339) + `"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "appointment not found")
	mockUC.AssertExpectations(t)
}

func TestReschedule_InvalidDatetime(t *testing.T) {
	// Error: Invalid datetime, returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "invalid-date",
		"end_date": "also-invalid"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
	mockUC.AssertExpectations(t)
}

func TestReschedule_Conflict(t *testing.T) {
	// Error: Conflict with other appointments, returns 409 Conflict
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	newStart := time.Now().Add(24 * time.Hour)
	newEnd := newStart.Add(1 * time.Hour)

	mockUC.On("Reschedule", mock.Anything, uint(1), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newStart.Year() && t.Month() == newStart.Month() && t.Day() == newStart.Day()
	}), mock.MatchedBy(func(t time.Time) bool {
		return t.Year() == newEnd.Year() && t.Month() == newEnd.Month() && t.Day() == newEnd.Day()
	})).Return(errors.New("client already has an appointment during this time"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "` + newStart.Format(time.RFC3339) + `",
		"end_date": "` + newEnd.Format(time.RFC3339) + `"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "client already has an appointment during this time")
	mockUC.AssertExpectations(t)
}

func TestReschedule_InvalidID(t *testing.T) {
	// Error: Invalid appointment ID format, returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/invalid/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "2024-01-15T10:00:00Z",
		"end_date": "2024-01-15T11:00:00Z"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid appointment ID format")
}

func TestReschedule_MissingStartDate(t *testing.T) {
	// Error: Missing required field (start_date), returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"end_date": "2024-01-15T11:00:00Z"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestReschedule_MissingEndDate(t *testing.T) {
	// Error: Missing required field (end_date), returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/reschedule", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"start_date": "2024-01-15T10:00:00Z"
	}`))

	h.Reschedule(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

// UpdateStatus Handler Tests
// **Validates: Requirements 2.7, 2.8**

func TestUpdateStatus_Success(t *testing.T) {
	// Success: Valid status, returns 200 OK
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("UpdateStatus", mock.Anything, uint(1), "completed").Return(nil)

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"status": "completed"
	}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Appointment status updated successfully")
	mockUC.AssertExpectations(t)
}

func TestUpdateStatus_NotFound(t *testing.T) {
	// Error: Appointment not found, returns 404 Not Found
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("UpdateStatus", mock.Anything, uint(999), "completed").Return(errors.New("appointment not found"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/999/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"status": "completed"
	}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "appointment not found")
	mockUC.AssertExpectations(t)
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	// Error: Invalid status, returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"status": "invalid-status"
	}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestUpdateStatus_InvalidTransition(t *testing.T) {
	// Error: Invalid status transition, returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	mockUC.On("UpdateStatus", mock.Anything, uint(1), "scheduled").Return(errors.New("cannot schedule an appointment for a time that has already passed"))

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"status": "scheduled"
	}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "cannot schedule an appointment for a time that has already passed")
	mockUC.AssertExpectations(t)
}

func TestUpdateStatus_InvalidID(t *testing.T) {
	// Error: Invalid appointment ID format, returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/invalid/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{
		"status": "completed"
	}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid appointment ID format")
}

func TestUpdateStatus_MissingStatus(t *testing.T) {
	// Error: Missing required field (status), returns 400 Bad Request
	gin.SetMode(gin.TestMode)
	mockUC := usecaseMocks.NewMockAppointmentUseCase()

	h := handler.NewAppointmentHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/status", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{}`))

	h.UpdateStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestUpdateStatus_AllValidStatuses(t *testing.T) {
	// Success: All valid statuses should be accepted
	validStatuses := []string{"scheduled", "completed", "cancelled", "archived", "no-show"}

	for _, status := range validStatuses {
		t.Run("Status_"+status, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			mockUC := usecaseMocks.NewMockAppointmentUseCase()

			mockUC.On("UpdateStatus", mock.Anything, uint(1), status).Return(nil)

			h := handler.NewAppointmentHandler(mockUC)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "1"}}
			c.Request, _ = http.NewRequest("PUT", "/api/v1/appointments/1/status", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Body = io.NopCloser(strings.NewReader(`{
				"status": "` + status + `"
			}`))

			h.UpdateStatus(c)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "Appointment status updated successfully")
			mockUC.AssertExpectations(t)
		})
	}
}
