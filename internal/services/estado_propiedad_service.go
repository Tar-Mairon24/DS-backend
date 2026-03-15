package services

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"

	"backend/internal/database"
	"backend/internal/models"
)

type EstadoPropiedadService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the EstadoPropiedadService
func NewEstadoPropiedadService(db *sql.DB) *EstadoPropiedadService {
	return &EstadoPropiedadService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// GET /estadoPropiedad/:id_tipo_propiedad
// Funcion que recupera el estado de la propiedad dependiedo del id_tipo_propiedad que biene en el get/prpopiedad/:id
func (service *EstadoPropiedadService) GetEstadoPropiedad(id int) (*models.EstadoPropiedades, error) {
	var estado models.EstadoPropiedades
	query := service.sq.Select("*").From("Estado_Propiedades").Where(sq.Eq{"id_propiedad": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	err = service.DB.QueryRow(sqlStr, args...).Scan(&estado.IDEstadoPropiedades, &estado.TipoTransaccion, &estado.Estado, &estado.FechaTransaccion, &estado.IDPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching estado:", err)
		return nil, err
	}

	return &estado, nil
}

// POST /estadoPropiedad/
// Funcion que crea un nuevo estado de la propiedad
func (service *EstadoPropiedadService) CreateEstadoPropiedad(estado *models.EstadoPropiedades) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Estado_Propiedades", "id_estado_propiedades")
	if err != nil {
		return 0, err
	}
	estado.IDEstadoPropiedades = lastId + 1
	query := service.sq.Insert("Estado_Propiedades").
		Columns("id_estado_propiedades", "tipo_transaccion", "estado", "fecha_cambio_estado", "id_propiedad").
		Values(estado.IDEstadoPropiedades, estado.TipoTransaccion, estado.Estado, estado.FechaTransaccion, estado.IDPropiedad)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error inserting estado de la propiedad:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error inserting estado: no rows affected")
		return 0, err
	}
	return estado.IDEstadoPropiedades, nil
}

// PUT /estadoPropiedad/:id
// Funcion que actualiza el estado de la propiedad
func (service *EstadoPropiedadService) UpdateEstadoPropiedad(estado *models.EstadoPropiedades) error {
	query := service.sq.Update("Estado_Propiedades").
		Set("tipo_transaccion", estado.TipoTransaccion).
		Set("estado", estado.Estado).
		Set("fecha_transaccion", estado.FechaTransaccion).
		Set("id_propiedad", estado.IDPropiedad).
		Where(sq.Eq{"id_estado_propiedades": estado.IDEstadoPropiedades})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error updating estado de la propiedad:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows != 1 {
		log.Println("Error updating estado: no rows affected")
		return err
	}
	return nil
}

// DELETE /eliminar/estadoPropiedad
// Function that deletes a EstadoPropiedad from the database
func (service *EstadoPropiedadService) DeleteEstadoPropiedad(id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Estado_Propiedades", "id_estado_propiedades")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	if id <= 0 || id > lastId {
		log.Println("Invalid estado ID:", id)
		return err
	}
	query := service.sq.Delete("Estado_Propiedades").Where(sq.Eq{"id_estado_propiedades": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error deleting estado de la propiedad:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error deleting estado: no rows affected")
		return err
	}
	return nil
}
