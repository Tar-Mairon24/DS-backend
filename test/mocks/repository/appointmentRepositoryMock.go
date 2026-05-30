package repositoryMock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockAppointmentRepository struct {
	mock.Mock
}

func NewMockAppointmentRepository() *MockAppointmentRepository {
	return &MockAppointmentRepository{}
}

func (m *MockAppointmentRepository) GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx)
	if appointments, ok := args.Get(0).([]models.AppointmentCalendarView); ok {
		return appointments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error) {
	args := m.Called(ctx, id)
	if appointment, ok := args.Get(0).(*models.AppointmentDetail); ok {
		return appointment, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, day)
	if appointments, ok := args.Get(0).([]models.AppointmentCalendarView); ok {
		return appointments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, from)
	if appointments, ok := args.Get(0).([]models.AppointmentCalendarView); ok {
		return appointments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, year, month)
	if appointments, ok := args.Get(0).([]models.AppointmentCalendarView); ok {
		return appointments, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) GetClientIDByAppointmentID(ctx context.Context, id uint) (uint, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAppointmentRepository) ClientHasOverlap(ctx context.Context, clientID uint, start, end time.Time, excludeID uint) (bool, error) {
	args := m.Called(ctx, clientID, start, end, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAppointmentRepository) Create(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	args := m.Called(ctx, appointment)
	if apt, ok := args.Get(0).(*models.Appointment); ok {
		return apt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) Update(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	args := m.Called(ctx, appointment)
	if apt, ok := args.Get(0).(*models.Appointment); ok {
		return apt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentRepository) UpdateTime(ctx context.Context, id uint, newStart, newEnd time.Time) error {
	args := m.Called(ctx, id, newStart, newEnd)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateStatus(ctx context.Context, id uint, newStatus string) error {
	args := m.Called(ctx, id, newStatus)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateAgents(ctx context.Context, appointmentID uint, agentIDs []uint) error {
	args := m.Called(ctx, appointmentID, agentIDs)
	return args.Error(0)
}

func (m *MockAppointmentRepository) AddAgents(ctx context.Context, appointmentID uint, agentIDs []uint) error {
	args := m.Called(ctx, appointmentID, agentIDs)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
