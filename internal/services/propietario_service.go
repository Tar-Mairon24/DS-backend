package services

import (
	"backend/internal/models"
	"backend/internal/database"
	"database/sql"
	"log"
)

type PropietarioService struct {
	DB *sql.DB
}

// Constructor for the PropiedadService
func NewPropietarioService(db *sql.DB) *PropietarioService {
	return &PropietarioService{
		DB: db,
	}
}

// GET /propietario/:id_propietario
// Funcion que recupera el propietario de la propiedad dependiendo del id_propietario que biene en el get/prpopiedad/:id
func (service *PropietarioService) GetPropietario(id int) (*models.Propietario, error) {
	var propietario models.Propietario
	query := "SELECT * FROM Propietario WHERE id_propietario = ?"
	err := service.DB.QueryRow(query, id).Scan(&propietario.IDPropietario, &propietario.Nombre, &propietario.ApellidoP, &propietario.ApellidoM, &propietario.Telefono, &propietario.Correo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching propietario:", err)
		return nil, err
	}
	return &propietario, nil
}

// POST /propietario
// Funcion que inserta un nuevo propietario en la base de datos
func (service *PropietarioService) CreatePropietario(propietario *models.Propietario) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Propietario", "id_propietario")
	if err != nil {
		log.Println("Error gettin last Id in Propietario table:", err)
		return 0, err
	}
	propietario.IDPropietario = lastId + 1
	query := "INSERT INTO Propietario (id_propietario, nombre, apellido_p, apellido_m, telefono, correo) VALUES (?, ?, ?, ?, ?)"
	result, err := service.DB.Exec(query, propietario.IDPropietario, propietario.Nombre, propietario.ApellidoP, propietario.ApellidoM, propietario.Telefono, propietario.Correo)
	if err != nil {
		log.Println("Error inserting propietario:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, err
	}
	if (rows != 0) {
		log.Println("No rows affected")
		return 0, nil
	}

	return propietario.IDPropietario, nil
}

// PUT /propietario
// Funcion que actualiza la informacion de un propietario en la base de datos
func (service *PropietarioService) UpdatePropietario(propietario *models.Propietario) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Propietario", "id_propietario")
	if err != nil {
		log.Println("Error gettin last Id in Propietario table:", err)
		return err
	}
	if(propietario.IDPropietario > lastId && propietario.IDPropietario < 0) {
		log.Println("Propietario ID not found")
		return nil
	}
	query := "UPDATE Propietario SET nombre = ?, apellido_p = ?, apellido_m = ?, telefono = ?, correo = ? WHERE id_propietario = ?"
	result, err := service.DB.Exec(query, propietario.Nombre, propietario.ApellidoP, propietario.ApellidoM, propietario.Telefono, propietario.Correo, propietario.IDPropietario)
	if err != nil {
		log.Println("Error updating propietario:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if (rows != 0) {
		log.Println("No rows affected")
		return nil
	}
	return nil
}

// DELETE /propietario/:id
// Funcion que elimina un propietario de la base de datos
func (service *PropietarioService) DeletePropietario(id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Propietario", "id_propietario")
	if err != nil {
		log.Println("Error gettin last Id in Propietario table:", err)
		return err
	}
	if(id > lastId && id < 0) {
		log.Println("Propietario ID not found")
		return nil
	}
	query := "DELETE FROM Propietario WHERE id_propietario = ?"
	result, err := service.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting propietario:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows == 0 {
		log.Println("No rows affected")
		return nil
	}
	return nil
}