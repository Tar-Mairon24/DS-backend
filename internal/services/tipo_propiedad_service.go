package services

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type TipoPropiedadService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the PropiedadService
func NewTipoPropiedadService(db *sql.DB) *TipoPropiedadService {
	return &TipoPropiedadService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// GET /tipoPropiedad/:id_tipo_propiedad
// Funcion que recupera el tipo de propiedad dependiendo del id_tipo_propiedad que biene en el get/prpopiedad/:id
func (service *TipoPropiedadService) GetTipoPropiedad(id int) (*models.TipoPropiedad, error) {
	var tipo models.TipoPropiedad
	query := service.sq.Select("*").From("Tipo_Propiedad").Where(sq.Eq{"id_tipo_propiedad": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	err = service.DB.QueryRow(sqlStr, args...).Scan(&tipo.IDTipoPropiedad, &tipo.Tipo_Propiedad)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching tipo:", err)
		return nil, err
	}
	return &tipo, nil
}

// POST /tipoPropiedad
// Funcion que inserta un nuevo tipo de propiedad en la base de datos
func (service *TipoPropiedadService) CreateTipoPropiedad(tipo *models.TipoPropiedad) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Tipo_Propiedad", "id_tipo_propiedad")
	if err != nil {
		log.Println("Error gettin last Id in Tipo_Propiedad table:", err)
		return 0, err
	}
	tipo.IDTipoPropiedad = lastId + 1
	query := service.sq.Insert("Tipo_Propiedad").
		Columns("id_tipo_propiedad", "tipo_propiedad").
		Values(tipo.IDTipoPropiedad, tipo.Tipo_Propiedad)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error inserting tipo:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, err
	}
	if rows == 0 {
		log.Println("No rows affected")
	}
	return tipo.IDTipoPropiedad, nil
}
