package usecaseMocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"ds-backend/internal/domain/models"
)

type MockAppointmentUseCase struct {
	mock.Mock
}

func NewMockAppointmentUseCase() *MockAppointmentUseCase {
	return &MockAppointmentUseCase{}
}

func (m *MockAppointmentUseCase) GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentCalendarView), args.Error(1)
}

func (m *MockAppointmentUseCase) GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error) {
	args := m.Called(ctx, id)
	if appt, ok := args.Get(0).(*models.AppointmentDetail); ok {
		return appt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentUseCase) GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, day)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentCalendarView), args.Error(1)
}

func (m *MockAppointmentUseCase) GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, from)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentCalendarView), args.Error(1)
}

func (m *MockAppointmentUseCase) GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error) {
	args := m.Called(ctx, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.AppointmentCalendarView), args.Error(1)
}

func (m *MockAppointmentUseCase) Create(ctx context.Context, req *models.AppointmentRequest) (*models.Appointment, error) {
	args := m.Called(ctx, req)
	if appt, ok := args.Get(0).(*models.Appointment); ok {
		return appt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentUseCase) Update(ctx context.Context, id uint, req *models.AppointmentRequest) (*models.Appointment, error) {
	args := m.Called(ctx, id, req)
	if appt, ok := args.Get(0).(*models.Appointment); ok {
		return appt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAppointmentUseCase) Reschedule(ctx context.Context, id uint, newStart, newEnd time.Time) error {
	args := m.Called(ctx, id, newStart, newEnd)
	return args.Error(0)
}

func (m *MockAppointmentUseCase) UpdateStatus(ctx context.Context, id uint, newStatus string) error {
	args := m.Called(ctx, id, newStatus)
	return args.Error(0)
}

func (m *MockAppointmentUseCase) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
