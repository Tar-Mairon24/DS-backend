package services

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"

	"backend/internal/database"
	"backend/internal/models"
)

type CitasService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the CitasService
func NewCitasService(db *sql.DB) *CitasService {
	return &CitasService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Funcion que recupera todas las citas de la base de datos
func (service *CitasService) GetAllCitasUser(IdUsuario string) ([]*models.CitaMenu, error) {
	var citas []*models.CitaMenu
	query := service.sq.Select("id_citas", "titulo_cita", "fecha_cita", "hora_cita", "nombre_prospecto", "apellido_paterno_prospecto", "apellido_materno_prospecto").
		From("Citas").
		Join("Prospecto ON Prospecto.id_cliente = Citas.id_cliente").
		Where(sq.Eq{"usuario": IdUsuario})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	rows, err := service.DB.Query(sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching all citas:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cita models.CitaMenu
		err := rows.Scan(&cita.IDCita, &cita.Titulo, &cita.FechaCita, &cita.HoraCita, &cita.NombreCliente, &cita.ApellidoPaternoCliente, &cita.ApellidoMaternoCliente)
		if err != nil {
			log.Println("Error scanning cita:", err)
			return nil, err
		}
		citas = append(citas, &cita)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return citas, nil
}

func (service *CitasService) GetAllCitasUserDay(IdUsuario string, day string) ([]*models.CitaMenu, error) {
	var citas []*models.CitaMenu
	query := service.sq.Select("id_citas", "titulo_cita", "fecha_cita", "hora_cita", "nombre_prospecto", "apellido_paterno_prospecto", "apellido_materno_prospecto").
		From("Citas").
		InnerJoin("Prospecto ON Prospecto.id_cliente = Citas.id_cliente").
		Where(sq.Eq{"usuario": IdUsuario, "fecha_cita": day})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	rows, err := service.DB.Query(sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching all citas:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cita models.CitaMenu
		err := rows.Scan(&cita.IDCita, &cita.Titulo, &cita.FechaCita, &cita.HoraCita, &cita.NombreCliente, &cita.ApellidoPaternoCliente, &cita.ApellidoMaternoCliente)
		if err != nil {
			log.Println("Error scanning cita:", err)
			return nil, err
		}
		citas = append(citas, &cita)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return citas, nil
}

func (service *CitasService) GetAllCitasUserMonth(IdUsuario int, Mes int) ([]*models.CitaMenu, error) {
	var citas []*models.CitaMenu
	// Modificación de la consulta para filtrar por mes, asumiendo que fecha_cita es un string en formato 'yyyy-mm-dd'
	query := service.sq.Select("id_citas", "titulo_cita", "fecha_cita", "hora_cita", "nombre_prospecto", "apellido_paterno_prospecto", "apellido_materno_prospecto").
		From("Citas").
		Join("Prospecto ON Prospecto.id_cliente = Citas.id_cliente").
		Where(sq.Eq{"id_usuario": IdUsuario}).
		Where(sq.Expr("MONTH(STR_TO_DATE(fecha_cita, '%Y-%m-%d')) = ?", Mes))

	// Ejecutamos la consulta, pasando el ID del usuario y el mes
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	rows, err := service.DB.Query(sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching all citas:", err)
		return nil, err
	}
	defer rows.Close()

	// Recorremos las filas devueltas por la consulta
	for rows.Next() {
		var cita models.CitaMenu
		err := rows.Scan(&cita.IDCita, &cita.Titulo, &cita.FechaCita, &cita.HoraCita,
			&cita.NombreCliente, &cita.ApellidoPaternoCliente, &cita.ApellidoMaternoCliente)
		if err != nil {
			log.Println("Error scanning cita:", err)
			return nil, err
		}
		citas = append(citas, &cita)
	}

	// Verificamos si hubo algún error mientras leíamos las filas
	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return citas, nil
}

func (service *CitasService) GetCita(id int) (*models.Cita, error) {
	var cita models.Cita
	query := service.sq.Select("*").From("Citas").Where(sq.Eq{"id_citas": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	row := service.DB.QueryRow(sqlStr, args...)
	err = row.Scan(&cita.IDCita, &cita.Titulo, &cita.FechaCita, &cita.HoraCita, &cita.Descripcion, &cita.IdUsuario, &cita.IdCliente)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching cita:", err)
		return nil, err
	}
	return &cita, nil
}
func (service *CitasService) InsertCita(cita *models.Cita) (int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Citas", "id_citas")
	if err != nil {
		log.Println("Error getting last Id:", err)
		return 0, err
	}
	cita.IDCita = lastId + 1
	cita.IdCliente = cita.IDCita

	query := service.sq.Insert("Citas").
		Columns("id_citas", "titulo_cita", "fecha_cita", "hora_cita", "descripcion_cita", "usuario", "id_cliente").
		Values(cita.IDCita, cita.Titulo, cita.FechaCita, cita.HoraCita, cita.Descripcion, cita.IdUsuario, cita.IdCliente)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error inserting cita:", err)
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, err
	}
	if rows != 1 {
		log.Println("Error inserting cita, no rows affected")
		return 0, err
	}
	return cita.IDCita, nil
}

// Funcion que actualiza una cita en la base de datos
func (service *CitasService) UpdateCita(cita *models.Cita, id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Citas", "id_citas")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	if id <= 0 || id > lastId {
		log.Println("Invalid cita ID:", id)
		return err
	}
	query := service.sq.Update("Citas").
		Set("titulo_cita", cita.Titulo).
		Set("fecha_cita", cita.FechaCita).
		Set("hora_cita", cita.HoraCita).
		Set("descripcion_cita", cita.Descripcion).
		Set("usuario", cita.IdUsuario).
		Set("id_cliente", cita.IdCliente).
		Where(sq.Eq{"id_citas": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error updating cita:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error updating cita, no rows affected")
		return err
	}
	return nil
}

// Funcion que elimina una cita de la base de datos
func (service *CitasService) DeleteCita(id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Citas", "id_citas")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	if id <= 0 || id > lastId {
		log.Println("Invalid propiedad ID:", id)
		return err
	}
	query := service.sq.Delete("Citas").Where(sq.Eq{"id_citas": id})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error deleting cita:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error deleting cita, no rows affected")
		return err
	}
	return nil
}
