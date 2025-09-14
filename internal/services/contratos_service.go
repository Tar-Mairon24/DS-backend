package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"
)

type ContratosService struct {
	DB *sql.DB
}

// Constructor para ContratosService
func NewContratosService(db *sql.DB) *ContratosService {
	return &ContratosService{
		DB: db,
	}
}

// Recupera un contrato por su ID
func (service *ContratosService) GetContrato(id int) (*models.Contrato, error) {
	var contrato models.Contrato
	query := "SELECT id_contrato, titulo_contrato, descripcion_contrato, tipo, ruta_pdf, id_propiedad FROM Contratos WHERE id_contrato = ?"
	row := service.DB.QueryRow(query, id)
	err := row.Scan(&contrato.IDContrato, &contrato.TituloContrato, &contrato.DescripcionContrato, &contrato.Tipo, &contrato.RutaPDF, &contrato.IDPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontró el contrato")
			return nil, nil
		}
		log.Println("Error recuperando el contrato:", err)
		return nil, err
	}
	return &contrato, nil
}

func (service *ContratosService) GetContratos() ([]*models.ContratoMenu, error) {
	var contratos []*models.ContratoMenu
	query := "SELECT id_contrato, titulo_contrato, tipo, titulo FROM Contratos, Propiedades WHERE Contratos.id_propiedad = Propiedades.id_propiedad"
	rows, err := service.DB.Query(query)
	if err != nil {
		log.Println("Error recuperando contratos:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contrato models.ContratoMenu
		err := rows.Scan(&contrato.IDContrato, &contrato.TituloContrato, &contrato.Tipo, &contrato.TituloPropiedad)
		if err != nil {
			log.Println("Error procesando fila de contrato:", err)
			return nil, err
		}
		contratos = append(contratos, &contrato)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterando filas:", err)
		return nil, err
	}
	return contratos, nil
}

// Recupera todos los contratos asociados a una propiedad
func (service *ContratosService) GetContratosByPropiedad(idPropiedad int) ([]*models.Contrato, error) {
	var contratos []*models.Contrato
	query := "SELECT id_contrato, titulo_contrato, descripcion_contrato, tipo, ruta_pdf, id_propiedad FROM Contratos WHERE id_propiedad = ?"
	rows, err := service.DB.Query(query, idPropiedad)
	if err != nil {
		log.Println("Error recuperando contratos:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contrato models.Contrato
		err := rows.Scan(&contrato.IDContrato, &contrato.TituloContrato, &contrato.DescripcionContrato, &contrato.Tipo, &contrato.RutaPDF, &contrato.IDPropiedad)
		if err != nil {
			log.Println("Error procesando fila de contrato:", err)
			return nil, err
		}
		contratos = append(contratos, &contrato)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterando filas:", err)
		return nil, err
	}
	return contratos, nil
}

// Inserta un nuevo contrato en la base de datos
func (service *ContratosService) InsertContrato(contrato *models.Contrato) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Contratos", "id_contrato")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	contrato.IDContrato = lastId + 1

	lastIdPropiedad, err := utils.GetLastId("Propiedades", "id_propiedad")
	if err != nil {
		log.Println("Error obteniendo último ID de propiedad:", err)
		return 0, err
	}
	contrato.IDPropiedad = lastIdPropiedad

	query := "INSERT INTO Contratos(id_contrato, titulo_contrato, descripcion_contrato, tipo, ruta_pdf, id_propiedad) VALUES(?,?,?,?,?,?)"
	result, err := service.DB.Exec(query, contrato.IDContrato, contrato.TituloContrato, contrato.DescripcionContrato, contrato.Tipo, contrato.RutaPDF, contrato.IDPropiedad)
	if err != nil {
		log.Println("Error insertando contrato:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error insertando contrato, no se afectaron filas")
		return 0, err
	}
	return contrato.IDContrato, nil
}

// Actualiza un contrato existente por su ID
func (service *ContratosService) UpdateContrato(contrato *models.Contrato, id int) error {
	query := "UPDATE Contratos SET titulo_contrato = ?, descripcion_contrato = ?, tipo = ?, ruta_pdf = ?, id_propiedad = ? WHERE id_contrato = ?"
	result, err := service.DB.Exec(query, contrato.TituloContrato, contrato.DescripcionContrato, contrato.Tipo, contrato.RutaPDF, contrato.IDPropiedad, id)
	if err != nil {
		log.Println("Error actualizando contrato:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error actualizando contrato, no se afectaron filas")
		return err
	}
	return nil
}

// Elimina un contrato por su ID
func (service *ContratosService) DeleteContrato(id int) error {
	query := "DELETE FROM Contratos WHERE id_contrato = ?"
	result, err := service.DB.Exec(query, id)
	if err != nil {
		log.Println("Error eliminando contrato:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error eliminando contrato, no se afectaron filas")
		return err
	}
	return nil
}
