package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"
)

type DocumentosAnexosService struct {
	DB *sql.DB
}

// Constructor para DocumentosAnexosService
func NewDocumentosAnexosService(db *sql.DB) *DocumentosAnexosService {
	return &DocumentosAnexosService{
		DB: db,
	}
}

// Recupera un documento anexo por su ID
func (service *DocumentosAnexosService) GetDocumentoAnexo(id int) (*models.DocumentoAnexo, error) {
	var documento models.DocumentoAnexo
	query := "SELECT id_documento_anexo, ruta_documento, descripcion_documento_anexo, id_propiedad FROM Documentos_Anexos WHERE id_documento_anexo = ?"
	row := service.DB.QueryRow(query, id)
	err := row.Scan(&documento.IDDocumentoAnexo, &documento.RutaDocumento, &documento.DescripcionDocumento, &documento.IDPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontró el documento anexo")
			return nil, nil
		}
		log.Println("Error recuperando el documento anexo:", err)
		return nil, err
	}
	return &documento, nil
}

// Recupera todos los documentos anexos de una propiedad
func (service *DocumentosAnexosService) GetDocumentosByPropiedad(idPropiedad int) ([]*models.DocumentoAnexo, error) {
	var documentos []*models.DocumentoAnexo
	query := "SELECT id_documento_anexo, ruta_documento, descripcion_documento_anexo, id_propiedad FROM Documentos_Anexos WHERE id_propiedad = ?"
	rows, err := service.DB.Query(query, idPropiedad)
	if err != nil {
		log.Println("Error recuperando documentos anexos:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var documento models.DocumentoAnexo
		err := rows.Scan(&documento.IDDocumentoAnexo, &documento.RutaDocumento, &documento.DescripcionDocumento, &documento.IDPropiedad)
		if err != nil {
			log.Println("Error procesando fila de documento anexo:", err)
			return nil, err
		}
		documentos = append(documentos, &documento)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterando filas:", err)
		return nil, err
	}
	return documentos, nil
}

// Inserta un nuevo documento anexo en la base de datos
func (service *DocumentosAnexosService) InsertDocumentoAnexo(documento *models.DocumentoAnexo) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Documentos_Anexos", "id_documento_anexo")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	documento.IDDocumentoAnexo = lastId + 1

	query := "INSERT INTO Documentos_Anexos(id_documento_anexo, ruta_documento, descripcion_documento_anexo, id_propiedad) VALUES(?, ?, ?, ?)"
	result, err := service.DB.Exec(query, documento.IDDocumentoAnexo, documento.RutaDocumento, documento.DescripcionDocumento, documento.IDPropiedad)
	if err != nil {
		log.Println("Error insertando documento anexo:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error insertando documento anexo, no se afectaron filas")
		return 0, nil
	}
	return documento.IDDocumentoAnexo, nil
}

// Actualiza un documento anexo existente por su ID
func (service *DocumentosAnexosService) UpdateDocumentoAnexo(documento *models.DocumentoAnexo, id int) error {
	query := "UPDATE Documentos_Anexos SET ruta_documento = ?, descripcion_documento_anexo = ?, id_propiedad = ? WHERE id_documento_anexo = ?"
	result, err := service.DB.Exec(query, documento.RutaDocumento, documento.DescripcionDocumento, documento.IDPropiedad, id)
	if err != nil {
		log.Println("Error actualizando documento anexo:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error actualizando documento anexo, no se afectaron filas")
		return nil
	}
	return nil
}

// Elimina un documento anexo por su ID
func (service *DocumentosAnexosService) DeleteDocumentoAnexo(id int) error {
	query := "DELETE FROM Documentos_Anexos WHERE id_documento_anexo = ?"
	result, err := service.DB.Exec(query, id)
	if err != nil {
		log.Println("Error eliminando documento anexo:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error eliminando documento anexo, no se afectaron filas")
		return nil
	}
	return nil
}
