package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type AppointmentHandler struct {
	appointmentUseCase ports.AppointmentUseCase
}

func NewAppointmentHandler(appointmentUseCase ports.AppointmentUseCase) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentUseCase: appointmentUseCase,
	}
}

func (h *AppointmentHandler) GetAll(c *gin.Context) {
	logrus.Info("GetAll Appointments endpoint called")

	ctx := c.Request.Context()
	appointments, err := h.appointmentUseCase.GetAll(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get appointments")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve appointments",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("%d appointments retrieved successfully", len(appointments))
	c.JSON(http.StatusOK, gin.H{
		"data":    appointments,
		"message": "Appointments retrieved successfully",
		"count":   len(appointments),
	})
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	logrus.Info("GetByID Appointment endpoint called")

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment ID format",
			"message": "ID must be a valid unsigned integer",
		})
		return
	}

	appointment, err := h.appointmentUseCase.GetByID(ctx, uint(id))
	if err != nil {
		logrus.WithError(err).Errorf("Failed to get appointment with ID %d", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve appointment",
			"message": err.Error(),
		})
		return
	}

	if appointment == nil {
		logrus.Warnf("Appointment with ID %d not found", id)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Appointment not found",
			"message": "No appointment found with the provided ID",
		})
		return
	}

	logrus.Infof("Appointment with ID %d retrieved successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"data":    appointment,
		"message": "Appointment retrieved successfully",
	})
}

func (h *AppointmentHandler) GetCalendar(c *gin.Context) {
    logrus.Info("GetCalendar endpoint called")
    view := c.DefaultQuery("view", "month")

    result, err := h.resolveCalendarView(c, view)
    if err != nil {
        logrus.WithError(err).Errorf("Failed to get calendar for view: %s", view)
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Failed to retrieve appointments",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":    result,
        "message": "Appointments retrieved successfully",
    })
}

func (h *AppointmentHandler) resolveCalendarView(c *gin.Context, view string) ([]models.AppointmentCalendarView, error) {
    ctx := c.Request.Context()

    switch view {
    case "day", "week":
        day, err := time.Parse("2006-01-02", c.Query("date"))
        if err != nil {
            return nil, errors.New("invalid date format, use YYYY-MM-DD")
        }
        if view == "day" {
            return h.appointmentUseCase.GetByDay(ctx, day)
        }
        return h.appointmentUseCase.GetByWeek(ctx, day)

    case "month":
        year, err := strconv.Atoi(c.Query("year"))
        if err != nil {
            return nil, errors.New("invalid year, must be an integer")
        }
        month, err := strconv.Atoi(c.Query("month"))
        if err != nil || month < 1 || month > 12 {
            return nil, errors.New("invalid month, use 1-12")
        }
        return h.appointmentUseCase.GetByMonth(ctx, year, month)

	default:
		return nil, fmt.Errorf("invalid view '%s', use day, week or month", view)
	}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	logrus.Info("Create Appointment endpoint called")
	var req models.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	appointment, err := h.appointmentUseCase.Create(ctx, &req)
	if err != nil {
		logrus.WithError(err).Error("Failed to create appointment")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create appointment",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Appointment created successfully with ID %d", appointment.ID)
	c.JSON(http.StatusCreated, gin.H{
		"data":    appointment,
		"message": "Appointment created successfully",
	})
}

func (h *AppointmentHandler) Update(c *gin.Context) {
	logrus.Info("Update Appointment endpoint called")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment ID format",
			"message": "ID must be a valid unsigned integer",
		})
		return
	}

	var req models.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	updatedAppointment, err := h.appointmentUseCase.Update(ctx, uint(id), &req)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to update appointment with ID %d", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update appointment",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Appointment with ID %d updated successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"data":    updatedAppointment,
		"message": "Appointment updated successfully",
	})
}

func (h *AppointmentHandler) Reschedule(c *gin.Context) {
	logrus.Info("Reschedule Appointment endpoint called")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment ID format",
			"message": "ID must be a valid unsigned integer",
		})
		return
	}

	var req models.RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.appointmentUseCase.Reschedule(ctx, uint(id), req.StartDate, req.EndDate); err != nil {
		logrus.WithError(err).Errorf("Failed to reschedule appointment with ID %d", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reschedule appointment",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Appointment with ID %d rescheduled successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment rescheduled successfully",
	})
}

func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	logrus.Info("Update Appointment Status endpoint called")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment ID format",
			"message": "ID must be a valid unsigned integer",
		})
		return
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.appointmentUseCase.UpdateStatus(ctx, uint(id), string(req.Status)); err != nil {
		logrus.WithError(err).Errorf("Failed to update status of appointment with ID %d", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update appointment status",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Status of appointment with ID %d updated successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment status updated successfully",
	})
}

func (h *AppointmentHandler) Delete(c *gin.Context) {
	logrus.Info("Delete Appointment endpoint called")
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid appointment ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment ID format",
			"message": "ID must be a valid unsigned integer",
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.appointmentUseCase.Delete(ctx, uint(id)); err != nil {
		logrus.WithError(err).Errorf("Failed to delete appointment with ID %d", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete appointment",
			"message": err.Error(),
		})
		return
	}

	logrus.Infof("Appointment with ID %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment deleted successfully",
	})
}