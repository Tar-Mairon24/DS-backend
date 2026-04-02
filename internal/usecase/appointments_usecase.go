package usecase

import (
	"context"
	"errors"
	"time"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type AppointmentUseCase struct {
	repo ports.AppointmentRepository
}

func NewAppointmentUseCase(repo ports.AppointmentRepository) ports.AppointmentUseCase {
	return &AppointmentUseCase{
		repo: repo,
	}
}

func (uc *AppointmentUseCase) GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *AppointmentUseCase) GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *AppointmentUseCase) GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error) {
	return uc.repo.GetByDay(ctx, day)
}

func (uc *AppointmentUseCase) GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error) {
	return uc.repo.GetByWeek(ctx, from)
}

func (uc *AppointmentUseCase) GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error) {
	return uc.repo.GetByMonth(ctx, year, month)
}

func (uc *AppointmentUseCase) Create(ctx context.Context, req *models.AppointmentRequest) (*models.Appointment, error) {
    overlap, err := uc.repo.ClientHasOverlap(ctx, req.ClientID, req.StartDate, req.EndDate, 0)
    if err != nil {
        return nil, err
    }
    if overlap {
        return nil, errors.New("client already has an appointment during this time")
    }

    appointment := &models.Appointment{
        Title:       req.Title,
        Description: req.Description,
        StartDate:   req.StartDate,
        EndDate:     req.EndDate,
        Status:      req.Status,
        Notes:       req.Notes,
        ClientID:    req.ClientID,
        PropertyID:  req.PropertyID,
    }

    return uc.repo.Create(ctx, appointment)
}

func (uc *AppointmentUseCase) Update(ctx context.Context, id uint, req *models.AppointmentRequest) (*models.Appointment, error) {
    overlap, err := uc.repo.ClientHasOverlap(ctx, req.ClientID, req.StartDate, req.EndDate, id)
    if err != nil {
        return nil, err
    }
    if overlap {
        return nil, errors.New("client already has an appointment during this time")
    }

    appointment := &models.Appointment{
        ID:          id,
        Title:       req.Title,
        Description: req.Description,
        StartDate:   req.StartDate,
        EndDate:     req.EndDate,
        Status:      req.Status,
        Notes:       req.Notes,
        ClientID:    req.ClientID,
        PropertyID:  req.PropertyID,
    }
    return uc.repo.Update(ctx, appointment)
}

func (uc *AppointmentUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
