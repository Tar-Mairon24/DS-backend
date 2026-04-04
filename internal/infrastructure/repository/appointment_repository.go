package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type AppointmentRepository struct {
	db *sql.DB
	qb sq.StatementBuilderType
}

func NewAppointmentRepository(db *sql.DB) ports.AppointmentRepository {
	return &AppointmentRepository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// formatTimeForMySQL converts time.Time to MySQL DATETIME string format
func formatTimeForMySQL(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

const getAppointmentByIDSQL = `
	SELECT
		a.id, a.title, a.description, a.start_date, a.end_date, a.status, a.notes,
		a.id_property        AS property_id,
		p.title              AS property_title,
		p.address            AS property_address,
		a.id_client          AS client_id,
		u_client.username    AS client_name,
		u_client.email       AS client_email,
		u_client.phone       AS client_phone,
		p.owner_id,
		u_owner.username     AS owner_name,
		u_owner.email        AS owner_email,
		COALESCE(
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'id',       u_agents.id,
					'username', u_agents.username,
					'email',    u_agents.email
				)
			),
			JSON_ARRAY()
		) AS agents
	FROM appointments a
	JOIN properties p       ON p.id = a.id_property
	JOIN users u_client     ON u_client.id = a.id_client  AND u_client.deleted_at IS NULL
	LEFT JOIN users u_owner ON u_owner.id  = p.owner_id   AND u_owner.deleted_at  IS NULL
	LEFT JOIN appointment_agents aa  ON aa.appointment_id  = a.id
	LEFT JOIN users u_agents         ON u_agents.id        = aa.user_id AND u_agents.deleted_at IS NULL
	WHERE a.id         = ?
	AND a.deleted_at IS NULL
	GROUP BY
		a.id, a.title, a.description, a.start_date, a.end_date, a.status, a.notes,
		a.id_property, p.title, p.address,
		a.id_client, u_client.username, u_client.email, u_client.phone,
		p.owner_id, u_owner.username, u_owner.email
	`

func (r *AppointmentRepository) GetByID(ctx context.Context, id uint) (*models.AppointmentDetail, error) {
	row := r.db.QueryRowContext(ctx, getAppointmentByIDSQL, id)

	var d models.AppointmentDetail
	var agentsJSON []byte

	err := row.Scan(
		&d.ID,
		&d.Title,
		&d.Description,
		&d.StartDate,
		&d.EndDate,
		&d.Status,
		&d.Notes,
		&d.PropertyID,
		&d.PropertyTitle,
		&d.PropertyAddress,
		&d.ClientID,
		&d.ClientName,
		&d.ClientEmail,
		&d.ClientPhone,
		&d.OwnerID,
		&d.OwnerName,
		&d.OwnerEmail,
		&agentsJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Warnf("appointment %d not found", id)
			return nil, err
		}
		logrus.WithError(err).Error("failed to scan appointment detail")
		return nil, err
	}

	if err := json.Unmarshal(agentsJSON, &d.Agents); err != nil {
		logrus.WithError(err).Error("failed to unmarshal agents JSON")
		return nil, err
	}

	return &d, nil
}

func (r *AppointmentRepository) lightQuery() sq.SelectBuilder {
	return r.qb.Select(
		"a.id", "a.title", "a.description",
		"a.start_date", "a.end_date", "a.status",
		"a.id_property", "p.title AS property_title",
		"a.id_client", "p.owner_id",
	).
		From("appointments a").
		Join("properties p ON p.id = a.id_property").
		Where(sq.Expr("a.deleted_at IS NULL"))
}

func (r *AppointmentRepository) scanLightRows(rows *sql.Rows) ([]models.AppointmentCalendarView, error) {
	defer rows.Close()
	var list []models.AppointmentCalendarView
	for rows.Next() {
		var v models.AppointmentCalendarView
		if err := rows.Scan(
			&v.ID, &v.Title, &v.Description,
			&v.StartDate, &v.EndDate, &v.Status,
			&v.Property.ID, &v.Property.Title,
			&v.ClientID, &v.OwnerID,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	if list == nil {
		return []models.AppointmentCalendarView{}, nil
	}
	return list, rows.Err()
}

func (r *AppointmentRepository) GetAll(ctx context.Context) ([]models.AppointmentCalendarView, error) {
	sqlStr, args, err := r.lightQuery().ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return r.scanLightRows(rows)
}

func (r *AppointmentRepository) GetByMonth(ctx context.Context, year int, month int) ([]models.AppointmentCalendarView, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	sqlStr, args, err := r.lightQuery().
		Where(sq.And{
			sq.GtOrEq{"a.start_date": start},
			sq.Lt{"a.start_date": end},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return r.scanLightRows(rows)
}

func (r *AppointmentRepository) GetByWeek(ctx context.Context, from time.Time) ([]models.AppointmentCalendarView, error) {
	start := from.Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 7)

	sqlStr, args, err := r.lightQuery().
		Where(sq.And{
			sq.GtOrEq{"a.start_date": start},
			sq.Lt{"a.start_date": end},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return r.scanLightRows(rows)
}

func (r *AppointmentRepository) GetByDay(ctx context.Context, day time.Time) ([]models.AppointmentCalendarView, error) {
	start := day.Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 1)

	sqlStr, args, err := r.lightQuery().
		Where(sq.And{
			sq.GtOrEq{"a.start_date": start},
			sq.Lt{"a.start_date": end},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return r.scanLightRows(rows)
}

func (r *AppointmentRepository) ClientHasOverlap(ctx context.Context, clientID uint, start, end time.Time, excludeID uint) (bool, error) {
	q := r.qb.Select("1").
		From("appointments").
		Where(sq.Eq{"id_client": clientID}).
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Expr("start_date < ? AND end_date > ?", end, start))

	if excludeID > 0 {
		q = q.Where(sq.NotEq{"id": excludeID})
	}

	sqlStr, args, err := q.Limit(1).ToSql()
	if err != nil {
		return false, err
	}

	var exists int
	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AppointmentRepository) Create(ctx context.Context, a *models.Appointment) (*models.Appointment, error) {
	sqlStr, args, err := r.qb.Insert("appointments").
		Columns("title", "description", "start_date", "end_date", "status", "notes", "id_client", "id_property").
		Values(a.Title, a.Description, formatTimeForMySQL(a.StartDate), formatTimeForMySQL(a.EndDate), a.Status, a.Notes, a.ClientID, a.PropertyID).
		ToSql()
	if err != nil {
		return nil, err
	}
	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	a.ID = uint(id)
	a.CreatedAt = time.Now()
	return a, nil
}

func (r *AppointmentRepository) Update(ctx context.Context, a *models.Appointment) (*models.Appointment, error) {
	sqlStr, args, err := r.qb.Update("appointments").
		Set("title", a.Title).
		Set("description", a.Description).
		Set("start_date", formatTimeForMySQL(a.StartDate)).
		Set("end_date", formatTimeForMySQL(a.EndDate)).
		Set("status", a.Status).
		Set("notes", a.Notes).
		Where(sq.Eq{"id": a.ID}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return a, err
}

func (r *AppointmentRepository) Delete(ctx context.Context, id uint) error {
	sqlStr, args, err := r.qb.Update("appointments").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
