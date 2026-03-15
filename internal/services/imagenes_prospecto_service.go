package services

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"

	"backend/internal/database"
	"backend/internal/models"
)

type ImagenesProspectoService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor para ImagenesService
func NewImagenesProspectoService(db *sql.DB) *ImagenesProspectoService {
	return &ImagenesProspectoService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Recupera una imagen por su ID
func (service *ImagenesProspectoService) GetImagen(id int) (*models.ImagenProspecto, error) {
	var imagen models.ImagenProspecto
	query := service.sq.Select("id_imagen", "ruta_imagen", "descripcion_imagen", "principal", "id_prospecto").From("ImagenesProspecto").Where(sq.Eq{"id_imagen": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	row := service.DB.QueryRow(sqlStr, args...)
	err = row.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDProspecto)
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

func (service *ImagenesProspectoService) GetImagenPrincipal(id int) (*models.ImagenProspecto, error) {
	var imagen models.ImagenProspecto
	query := service.sq.Select("id_imagen", "ruta_imagen", "descripcion_imagen", "principal", "id_prospecto").From("ImagenesProspecto").Where(sq.Eq{"id_prospecto": id, "principal": 1})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	row := service.DB.QueryRow(sqlStr, args...)
	err = row.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDProspecto)
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
func (service *ImagenesProspectoService) GetImagenesByProspecto(idPropiedad int) ([]*models.ImagenProspecto, error) {
	var imagenes []*models.ImagenProspecto
	query := service.sq.Select("id_imagen", "ruta_imagen", "descripcion_imagen", "principal", "id_prospecto").From("ImagenesProspecto").Where(sq.Eq{"id_prospecto": idPropiedad})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	rows, err := service.DB.Query(sqlStr, args...)
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
		var imagen models.ImagenProspecto
		err := rows.Scan(&imagen.IDImagen, &imagen.RutaImagen, &imagen.Descripcion, &imagen.Principal, &imagen.IDProspecto)
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
func (service *ImagenesProspectoService) InsertImagen(imagen *models.ImagenProspecto) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("ImagenesProspecto", "id_imagen")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	imagen.IDImagen = lastId + 1

	lastIdProspecto, err := utils.GetLastId("Citas", "id_citas")
	if err != nil {
		log.Println("Error obteniendo último ID:", err)
		return 0, err
	}
	imagen.IDProspecto = lastIdProspecto

	query := service.sq.Insert("ImagenesProspecto").
		Columns("id_imagen", "ruta_imagen", "descripcion_imagen", "principal", "id_prospecto").
		Values(imagen.IDImagen, imagen.RutaImagen, imagen.Descripcion, imagen.Principal, imagen.IDProspecto)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
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
func (service *ImagenesProspectoService) UpdateImagen(imagen *models.ImagenProspecto, id int) error {
	query := service.sq.Update("ImagenesProspecto").
		Set("ruta_imagen", imagen.RutaImagen).
		Set("descripcion_imagen", imagen.Descripcion).
		Set("principal", imagen.Principal).
		Set("id_prospecto", imagen.IDProspecto).
		Where(sq.Eq{"id_imagen": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
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
func (service *ImagenesProspectoService) DeleteImagen(id int) error {
	query := service.sq.Delete("ImagenesProspecto").Where(sq.Eq{"id_imagen": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
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
