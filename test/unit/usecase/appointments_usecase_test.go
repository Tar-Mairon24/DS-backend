package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/usecase"
	repositoryMock "ds-backend/test/mocks/repository"
)

// TestRescheduleAppointment tests the RescheduleAppointment method
func TestRescheduleAppointment(t *testing.T) {
	futureStart := time.Now().Add(24 * time.Hour)
	futureEnd := time.Now().Add(25 * time.Hour)

	tests := []struct {
		name             string
		appointmentID    uint
		newStart         time.Time
		newEnd           time.Time
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:          "success: reschedule appointment",
			appointmentID: 1,
			newStart:      futureStart,
			newEnd:        futureEnd,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetClientIDByAppointmentID", context.Background(), uint(1)).Return(uint(1), nil)
				mockRepo.On("ClientHasOverlap", context.Background(), uint(1), futureStart, futureEnd, uint(1)).Return(false, nil)
				mockRepo.On("UpdateTime", context.Background(), uint(1), futureStart, futureEnd).Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "error: appointment not found",
			appointmentID: 999,
			newStart:      futureStart,
			newEnd:        futureEnd,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetClientIDByAppointmentID", context.Background(), uint(999)).
					Return(uint(0), errors.New("appointment not found"))
			},
			expectedError:    true,
			expectedErrorMsg: "appointment not found",
		},
		{
			name:          "error: conflict with existing appointment",
			appointmentID: 1,
			newStart:      futureStart,
			newEnd:        futureEnd,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetClientIDByAppointmentID", context.Background(), uint(1)).Return(uint(1), nil)
				mockRepo.On("ClientHasOverlap", context.Background(), uint(1), futureStart, futureEnd, uint(1)).Return(true, nil)
			},
			expectedError:    true,
			expectedErrorMsg: "client already has an appointment during this time",
		},
		{
			name:          "error: repository error on overlap check",
			appointmentID: 1,
			newStart:      futureStart,
			newEnd:        futureEnd,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetClientIDByAppointmentID", context.Background(), uint(1)).Return(uint(1), nil)
				mockRepo.On("ClientHasOverlap", context.Background(), uint(1), futureStart, futureEnd, uint(1)).Return(false, errors.New("database error"))
			},
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			err := uc.Reschedule(context.Background(), tt.appointmentID, tt.newStart, tt.newEnd)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateAppointmentStatus tests the UpdateStatus method
func TestUpdateAppointmentStatus(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	pastTime := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name             string
		appointmentID    uint
		newStatus        string
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:          "success: update status to completed",
			appointmentID: 1,
			newStatus:     "completed",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("UpdateStatus", context.Background(), uint(1), "completed").Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "success: update status to cancelled",
			appointmentID: 1,
			newStatus:     "cancelled",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("UpdateStatus", context.Background(), uint(1), "cancelled").Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "success: update status to no-show",
			appointmentID: 1,
			newStatus:     "no-show",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("UpdateStatus", context.Background(), uint(1), "no-show").Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "success: update status to archived",
			appointmentID: 1,
			newStatus:     "archived",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("UpdateStatus", context.Background(), uint(1), "archived").Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "success: reschedule to future time",
			appointmentID: 1,
			newStatus:     "scheduled",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByID", context.Background(), uint(1)).Return(&models.AppointmentDetail{
					ID:        1,
					StartDate: futureTime,
				}, nil)
				mockRepo.On("UpdateStatus", context.Background(), uint(1), "scheduled").Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "error: appointment not found",
			appointmentID: 999,
			newStatus:     "completed",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("UpdateStatus", context.Background(), uint(999), "completed").
					Return(errors.New("appointment not found"))
			},
			expectedError:    true,
			expectedErrorMsg: "appointment not found",
		},
		{
			name:          "error: cannot schedule past appointment",
			appointmentID: 1,
			newStatus:     "scheduled",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByID", context.Background(), uint(1)).Return(&models.AppointmentDetail{
					ID:        1,
					StartDate: pastTime,
				}, nil)
			},
			expectedError:    true,
			expectedErrorMsg: "cannot schedule an appointment for a time that has already passed",
		},
		{
			name:          "error: repository error on get",
			appointmentID: 1,
			newStatus:     "scheduled",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByID", context.Background(), uint(1)).
					Return(nil, errors.New("database error"))
			},
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			err := uc.UpdateStatus(context.Background(), tt.appointmentID, tt.newStatus)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAllAppointments tests retrieving all appointments
func TestGetAllAppointments(t *testing.T) {
	tests := []struct {
		name             string
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedCount    int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "success: get all appointments",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				appointments := []models.AppointmentCalendarView{
					{
						ID:       1,
						Title:    "Appointment 1",
						ClientID: 1,
						Status:   "scheduled",
					},
					{
						ID:       2,
						Title:    "Appointment 2",
						ClientID: 2,
						Status:   "completed",
					},
				}
				mockRepo.On("GetAll", context.Background()).Return(appointments, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "success: no appointments found",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetAll", context.Background()).Return([]models.AppointmentCalendarView{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "error: repository error",
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetAll", context.Background()).Return(nil, errors.New("database error"))
			},
			expectedCount:    0,
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			result, err := uc.GetAll(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAppointmentByID tests retrieving a single appointment
func TestGetAppointmentByID(t *testing.T) {
	tests := []struct {
		name             string
		appointmentID    uint
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedError    bool
		expectedErrorMsg string
		expectedID       uint
	}{
		{
			name:          "success: get appointment by id",
			appointmentID: 1,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByID", context.Background(), uint(1)).Return(&models.AppointmentDetail{
					ID:         1,
					Title:      "Test Appointment",
					ClientID:   1,
					PropertyID: 1,
					Status:     "scheduled",
				}, nil)
			},
			expectedError: false,
			expectedID:    1,
		},
		{
			name:          "error: appointment not found",
			appointmentID: 999,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByID", context.Background(), uint(999)).
					Return(nil, errors.New("appointment not found"))
			},
			expectedError:    true,
			expectedErrorMsg: "appointment not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			result, err := uc.GetByID(context.Background(), tt.appointmentID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedID, result.ID)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAppointmentsByDay tests retrieving appointments for a specific day
func TestGetAppointmentsByDay(t *testing.T) {
	day := time.Now()

	tests := []struct {
		name             string
		day              time.Time
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedCount    int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "success: get appointments for day",
			day:  day,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				appointments := []models.AppointmentCalendarView{
					{
						ID:       1,
						Title:    "Morning Appointment",
						ClientID: 1,
						Status:   "scheduled",
					},
					{
						ID:       2,
						Title:    "Afternoon Appointment",
						ClientID: 2,
						Status:   "scheduled",
					},
				}
				mockRepo.On("GetByDay", context.Background(), day).Return(appointments, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "success: no appointments for day",
			day:  day,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByDay", context.Background(), day).Return([]models.AppointmentCalendarView{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "error: repository error",
			day:  day,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByDay", context.Background(), day).Return(nil, errors.New("database error"))
			},
			expectedCount:    0,
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			result, err := uc.GetByDay(context.Background(), tt.day)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestDeleteAppointment tests deleting an appointment
func TestDeleteAppointment(t *testing.T) {
	tests := []struct {
		name             string
		appointmentID    uint
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:          "success: delete appointment",
			appointmentID: 1,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("Delete", context.Background(), uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "error: appointment not found",
			appointmentID: 999,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("Delete", context.Background(), uint(999)).
					Return(errors.New("appointment not found"))
			},
			expectedError:    true,
			expectedErrorMsg: "appointment not found",
		},
		{
			name:          "error: repository error",
			appointmentID: 1,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("Delete", context.Background(), uint(1)).
					Return(errors.New("database error"))
			},
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			err := uc.Delete(context.Background(), tt.appointmentID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAppointmentsByWeek tests retrieving appointments for a week
func TestGetAppointmentsByWeek(t *testing.T) {
	weekStart := time.Now()

	tests := []struct {
		name             string
		weekStart        time.Time
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedCount    int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:      "success: get appointments for week",
			weekStart: weekStart,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				appointments := []models.AppointmentCalendarView{
					{
						ID:       1,
						Title:    "Monday Appointment",
						ClientID: 1,
						Status:   "scheduled",
					},
					{
						ID:       2,
						Title:    "Wednesday Appointment",
						ClientID: 2,
						Status:   "scheduled",
					},
					{
						ID:       3,
						Title:    "Friday Appointment",
						ClientID: 3,
						Status:   "completed",
					},
				}
				mockRepo.On("GetByWeek", context.Background(), weekStart).Return(appointments, nil)
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:      "success: no appointments for week",
			weekStart: weekStart,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByWeek", context.Background(), weekStart).Return([]models.AppointmentCalendarView{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:      "error: repository error",
			weekStart: weekStart,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByWeek", context.Background(), weekStart).Return(nil, errors.New("database error"))
			},
			expectedCount:    0,
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			result, err := uc.GetByWeek(context.Background(), tt.weekStart)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}

// TestGetAppointmentsByMonth tests retrieving appointments for a month
func TestGetAppointmentsByMonth(t *testing.T) {
	tests := []struct {
		name             string
		year             int
		month            int
		setupMocks       func(*repositoryMock.MockAppointmentRepository)
		expectedCount    int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:  "success: get appointments for month",
			year:  2024,
			month: 1,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				appointments := []models.AppointmentCalendarView{
					{
						ID:       1,
						Title:    "January Appointment 1",
						ClientID: 1,
						Status:   "scheduled",
					},
					{
						ID:       2,
						Title:    "January Appointment 2",
						ClientID: 2,
						Status:   "completed",
					},
				}
				mockRepo.On("GetByMonth", context.Background(), 2024, 1).Return(appointments, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:  "success: no appointments for month",
			year:  2024,
			month: 2,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByMonth", context.Background(), 2024, 2).Return([]models.AppointmentCalendarView{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:  "error: repository error",
			year:  2024,
			month: 3,
			setupMocks: func(mockRepo *repositoryMock.MockAppointmentRepository) {
				mockRepo.On("GetByMonth", context.Background(), 2024, 3).Return(nil, errors.New("database error"))
			},
			expectedCount:    0,
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAppointmentRepo := repositoryMock.NewMockAppointmentRepository()
			mockUserRepo := repositoryMock.NewMockUserRepo()
			mockPropertyRepo := repositoryMock.NewMockPropertyRepository()

			tt.setupMocks(mockAppointmentRepo)

			uc := usecase.NewAppointmentUseCase(mockAppointmentRepo, mockUserRepo, mockPropertyRepo)

			result, err := uc.GetByMonth(context.Background(), tt.year, tt.month)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrorMsg, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedCount)
			}

			mockAppointmentRepo.AssertExpectations(t)
		})
	}
}
