package services

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"backend/internal/database"
	"backend/internal/models"
)

type ProspectoService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the CitasService
func NewProspectoService(db *sql.DB) *ProspectoService {
	return &ProspectoService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

func (service *ProspectoService) GetProspecto(id int) (*models.Prospecto, error) {
	var prospecto models.Prospecto
	query := service.sq.Select("*").From("Prospecto").Where(sq.Eq{"id_cliente": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	row := service.DB.QueryRow(sqlStr, args...)
	err = row.Scan(&prospecto.IdCliente, &prospecto.Nombre, &prospecto.ApellidoP, &prospecto.ApellidoM, &prospecto.Telefono, &prospecto.Correo)
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

	query := service.sq.Insert("Prospecto").
		Columns("id_cliente", "nombre_prospecto", "apellido_paterno_prospecto", "apellido_materno_prospecto", "telefono_prospecto", "correo_prospecto").
		Values(prospecto.IdCliente, prospecto.Nombre, prospecto.ApellidoP, prospecto.ApellidoM, prospecto.Telefono, prospecto.Correo)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
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
	query := service.sq.Update("Prospecto").
		Set("nombre_prospecto", prospecto.Nombre).
		Set("apellido_paterno_prospecto", prospecto.ApellidoP).
		Set("apellido_materno_prospecto", prospecto.ApellidoM).
		Set("telefono_prospecto", prospecto.Telefono).
		Set("correo_prospecto", prospecto.Correo).
		Where(sq.Eq{"id_cliente": id}).
		Where(sq.Eq{"borrado_en": nil})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
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

func (service *ProspectoService) DeleteProspecto(id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Prospecto", "id_cliente")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	if id <= 0 || id > lastId {
		log.Println("Invalid prospecto ID:", id)
		return nil
	}

	query := service.sq.Update("Prospecto").Set("borrado_en", time.Now()).Where(sq.Eq{"id_cliente": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error deleting prospecto:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows != 1 {
		log.Println("Error deleting prospecto: no rows affected")
		return err
	}
	return nil
}

func (service *ProspectoService) GetAllProspectos() ([]*models.Prospecto, error) {
	var prospectos []*models.Prospecto
	query := service.sq.Select("*").From("Prospecto").Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	rows, err := service.DB.Query(sqlStr, args...)
	if err != nil {
		log.Println("Error fetching prospectos:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prospecto models.Prospecto
		err := rows.Scan(&prospecto.IdCliente, &prospecto.Nombre, &prospecto.ApellidoP, &prospecto.ApellidoM, &prospecto.Telefono, &prospecto.Correo)
		if err != nil {
			log.Println("Error scanning prospecto:", err)
			return nil, err
		}
		prospectos = append(prospectos, &prospecto)
	}
	if err = rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}
	return prospectos, nil
}
