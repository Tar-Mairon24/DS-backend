package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"
)

type ImagenesService struct {
	DB *sql.DB
}

// Constructor para ImagenesService
func NewImagenesService(db *sql.DB) *ImagenesService {
	return &ImagenesService{
		DB: db,
	}
}

// Recupera una imagen por su ID
func (service *ImagenesService) GetImagen(id int) (*models.Imagen, error) {
	var imagen models.Imagen
	query := "SELECT id_imagen, ruta_imagen, descripcion_imagen, principal, id_propiedad FROM Imagenes WHERE id_imagen = ?"
	row := service.DB.QueryRow(query, id)
	err := row.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontró la imagen")
			return nil, nil
		}
		log.Println("Error recuperando la imagen:", err)
		return nil, err
	}
	return &imagen, nil
}

func (service *ImagenesService) GetImagenPrincipal(id int) (*models.Imagen, error) {
	var imagen models.Imagen
	query := "SELECT id_imagen, ruta_imagen, descripcion_imagen, principal, id_propiedad FROM Imagenes WHERE id_propiedad = ? AND principal = 1"
	row := service.DB.QueryRow(query, id)
	err := row.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontró la imagen")
			return nil, nil
		}
		log.Println("Error recuperando la imagen:", err)
		return nil, err
	}
	return &imagen, nil
}

// Recupera todas las imágenes de una propiedad
func (service *ImagenesService) GetImagenesByPropiedad(idPropiedad int) ([]*models.Imagen, error) {
	var imagenes []*models.Imagen
	query := "SELECT id_imagen, ruta_imagen, descripcion_imagen, principal, id_propiedad FROM Imagenes WHERE id_propiedad = ?"
	rows, err := service.DB.Query(query, idPropiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No se encontraron imágenes")
			return nil, nil
		}
		log.Println("Error recuperando imágenes:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var imagen models.Imagen
		err := rows.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDPropiedad)
		if err != nil {
			log.Println("Error procesando fila de imagen:", err)
			return nil, err
		}
		imagenes = append(imagenes, &imagen)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterando filas:", err)
		return nil, err
	}
	return imagenes, nil
}

// Inserta una nueva imagen en la base de datos
func (service *ImagenesService) InsertImagen(imagen *models.Imagen) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Imagenes", "id_imagen")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	imagen.IDImagen = lastId + 1

	lastIdPropiedad, err := utils.GetLastId("Propiedades", "id_propiedad")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	imagen.IDPropiedad = lastIdPropiedad

	query := "INSERT INTO Imagenes(id_imagen, ruta_imagen, descripcion_imagen, principal, id_propiedad) VALUES(?,?,?,?,?)"
	result, err := service.DB.Exec(query, imagen.IDImagen, imagen.RutaImagen, imagen.Descripcion, imagen.Principal, imagen.IDPropiedad)
	if err != nil {
		log.Println("Error insertando imagen:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error insertando imagen, no se afectaron filas")
		return 0, err
	}
	return imagen.IDImagen, nil
}

// Actualiza una imagen existente por su ID
func (service *ImagenesService) UpdateImagen(imagen *models.Imagen, id int) error {
	query := "UPDATE Imagenes SET ruta_imagen = ?, descripcion_imagen = ?, principal = ?, id_propiedad = ? WHERE id_imagen = ?"
	result, err := service.DB.Exec(query, imagen.RutaImagen, imagen.Descripcion, imagen.Principal, imagen.IDPropiedad, id)
	if err != nil {
		log.Println("Error actualizando imagen:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error actualizando imagen, no se afectaron filas")
		return err
	}
	return nil
}

// Elimina una imagen por su ID
func (service *ImagenesService) DeleteImagen(id int) error {
	query := "DELETE FROM Imagenes WHERE id_imagen = ?"
	result, err := service.DB.Exec(query, id)
	if err != nil {
		log.Println("Error eliminando imagen:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error verificando filas afectadas:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error eliminando imagen, no se afectaron filas")
		return err
	}
	return nil
}
