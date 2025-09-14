package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"
)

type ProspectoService struct {
	DB *sql.DB
}

// Constructor for the CitasService
func NewProspectoService(db *sql.DB) *ProspectoService {
	return &ProspectoService{
		DB: db,
	}
}

func (service *ProspectoService) GetProspecto(id int) (*models.Prospecto, error) {
	var prospecto models.Prospecto
	query := "SELECT * FROM Prospecto WHERE id_cliente = ?"
	row := service.DB.QueryRow(query, id)
	err := row.Scan(&prospecto.IdCliente, &prospecto.Nombre, &prospecto.ApellidoP, &prospecto.ApellidoM, &prospecto.Telefono, &prospecto.Correo)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching prospecto:", err)
		return nil, err
	}
	return &prospecto, nil
}

func (service *ProspectoService) InsertProspecto(prospecto *models.Prospecto) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Prospecto", "id_cliente")
	if err != nil {
		log.Println("Error getting last Id:", err)
		return 0, err
	}
	prospecto.IdCliente = lastId + 1

	query := "INSERT INTO Prospecto(id_cliente, nombre_prospecto, apellido_paterno_prospecto, apellido_materno_prospecto, telefono_prospecto, correo_prospecto) VALUES(?,?,?,?,?,?)"
	result, err := service.DB.Exec(query, prospecto.IdCliente, prospecto.Nombre, prospecto.ApellidoP, prospecto.ApellidoM, prospecto.Telefono, prospecto.Correo)
	if err != nil {
		log.Println("Error inserting prospecto:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error inserting prospecto, no rows affected")
		return 0, err
	}
	return prospecto.IdCliente, nil
}

func (service *ProspectoService) UpdateProspecto(prospecto *models.Prospecto, id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Prospecto", "id_cliente")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	println(lastId)
	if id <= 0 || id > lastId {
		log.Println("Invalid prospecto ID:", id)
		return err
	}
	println(id)
	query := "UPDATE Prospecto SET nombre_prospecto=?, apellido_paterno_prospecto=?, apellido_materno_prospecto=?, telefono_prospecto=?, correo_prospecto=? WHERE id_cliente=?"
	result, err := service.DB.Exec(query, prospecto.Nombre, prospecto.ApellidoP, prospecto.ApellidoM, prospecto.Telefono, prospecto.Correo, id)
	if err != nil {
		log.Println("Error updating prospecto:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows != 1 {
		log.Println("Error updating prospecto: no rows affected")
		return err
	}
	return nil
}
