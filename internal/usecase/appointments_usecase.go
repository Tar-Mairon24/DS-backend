package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type AppointmentUseCase struct {
	repo         ports.AppointmentRepository
	userRepo     ports.UserRepository
	propertyRepo ports.PropertyRepository
}

func NewAppointmentUseCase(repo ports.AppointmentRepository, userRepo ports.UserRepository, propertyRepo ports.PropertyRepository) ports.AppointmentUseCase {
	return &AppointmentUseCase{
		repo:         repo,
		userRepo:     userRepo,
		propertyRepo: propertyRepo,
	}
}

func (uc *AppointmentUseCase) validateAgents(ctx context.Context, agentIDs []uint) error {
	for _, agentID := range agentIDs {
		user, err := uc.userRepo.GetByID(ctx, agentID)
		if err != nil {
			logrus.Warnf("Failed to get user %d: %v", agentID, err)
			return errors.New("agent not found")
		}
		logrus.Infof("Validating agent %d with role: '%s' (agent='%s', admin='%s')", agentID, user.Role, models.UserTypeAgent, models.UserTypeAdmin)
		if user.Role != models.UserTypeAgent && user.Role != models.UserTypeAdmin {
			logrus.Warnf("Agent validation failed for user %d with role '%s'", agentID, user.Role)
			return errors.New("user must be an agent or admin to be assigned to an appointment")
		}
	}
	return nil
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
	_, err := uc.userRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	_, err = uc.propertyRepo.GetByID(ctx, req.PropertyID)
	if err != nil {
		return nil, errors.New("property not found")
	}

	if err := uc.validateAgents(ctx, req.AgentIDs); err != nil {
		return nil, err
	}

	if req.Status == "scheduled" && req.StartDate.Time.Before(time.Now()) {
		return nil, errors.New("cannot schedule an appointment for a time that has already passed")
	}

	overlap, err := uc.repo.ClientHasOverlap(ctx, req.ClientID, req.StartDate.Time, req.EndDate.Time, 0)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("client already has an appointment during this time")
	}

	appointment := &models.Appointment{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate.Time,
		EndDate:     req.EndDate.Time,
		Status:      req.Status,
		Notes:       req.Notes,
		ClientID:    req.ClientID,
		PropertyID:  req.PropertyID,
		AgentIDs:    req.AgentIDs,
	}

	return uc.repo.Create(ctx, appointment)
}

func (uc *AppointmentUseCase) Update(ctx context.Context, id uint, req *models.AppointmentRequest) (*models.Appointment, error) {
	_, err := uc.repo.GetClientIDByAppointmentID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := uc.validateAgents(ctx, req.AgentIDs); err != nil {
		return nil, err
	}

	if req.Status == "scheduled" && req.StartDate.Time.Before(time.Now()) {
		return nil, errors.New("cannot schedule an appointment for a time that has already passed")
	}

	overlap, err := uc.repo.ClientHasOverlap(ctx, req.ClientID, req.StartDate.Time, req.EndDate.Time, id)
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
		StartDate:   req.StartDate.Time,
		EndDate:     req.EndDate.Time,
		Status:      req.Status,
		Notes:       req.Notes,
		ClientID:    req.ClientID,
		PropertyID:  req.PropertyID,
		AgentIDs:    req.AgentIDs,
	}
	return uc.repo.Update(ctx, appointment)
}

func (uc *AppointmentUseCase) Reschedule(ctx context.Context, id uint, newStart, newEnd time.Time) error {
	clientID, err := uc.repo.GetClientIDByAppointmentID(ctx, id)
	if err != nil {
		return err
	}

	overlap, err := uc.repo.ClientHasOverlap(ctx, clientID, newStart, newEnd, id)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("client already has an appointment during this time")
	}

	return uc.repo.UpdateTime(ctx, id, newStart, newEnd)
}

func (uc *AppointmentUseCase) UpdateStatus(ctx context.Context, id uint, newStatus string) error {
	if newStatus == "scheduled" {
		appointment, err := uc.repo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if appointment.StartDate.Before(time.Now()) {
			return errors.New("cannot schedule an appointment for a time that has already passed")
		}
	}
	return uc.repo.UpdateStatus(ctx, id, newStatus)
}

func (uc *AppointmentUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
