package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"
)

type EstadoPropiedadService struct {
	DB *sql.DB
}

// Constructor for the EstadoPropiedadService
func NewEstadoPropiedadService(db *sql.DB) *EstadoPropiedadService {
	return &EstadoPropiedadService{
		DB: db,
	}
}

// GET /estadoPropiedad/:id_tipo_propiedad
// Funcion que recupera el estado de la propiedad dependiedo del id_tipo_propiedad que biene en el get/prpopiedad/:id
func (service *EstadoPropiedadService) GetEstadoPropiedad(id int) (*models.EstadoPropiedades, error) {
	var estado models.EstadoPropiedades
	query := "SELECT * FROM Estado_Propiedades WHERE id_propiedad = ?"
	err := service.DB.QueryRow(query, id).Scan(&estado.IDEstadoPropiedades, &estado.TipoTransaccion, &estado.Estado, &estado.FechaTransaccion, &estado.IDPropiedad)
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
	query := "INSERT INTO Estado_Propiedades (id_estado_propiedades, tipo_transaccion, estado, fecha_cambio_estado, id_propiedad) VALUES (?, ?, ?, ?, ?)"
	result, err := service.DB.Exec(query, estado.IDEstadoPropiedades, estado.TipoTransaccion, estado.Estado, estado.FechaTransaccion, estado.IDPropiedad)
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
	query := "UPDATE Estado_Propiedades SET tipo_transaccion = ?, estado = ?, fecha_transaccion = ?, id_propiedad = ? WHERE id_estado_propiedades = ?"
	result, err := service.DB.Exec(query, estado.TipoTransaccion, estado.Estado, estado.FechaTransaccion, estado.IDPropiedad, estado.IDEstadoPropiedades)
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
	query := "DELETE FROM Estado_Propiedades WHERE id_estado_propiedades = ?"
	result, err := service.DB.Exec(query, id)
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
