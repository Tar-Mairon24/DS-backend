package services

import (
	"database/sql"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"backend/internal/database"
	"backend/internal/models"
)

type PropiedadService struct {
	DB *sql.DB
	sq sq.StatementBuilderType
}

// Constructor for the PropiedadService
func NewPropiedadService(db *sql.DB) *PropiedadService {
	return &PropiedadService{
		DB: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Funcion que recupera todas las propiedades de la base de datos, solo recupera los campos necesarios para mostrar en el menú, el resto de los campos se recuperan en otra función
func (service *PropiedadService) GetAllPropiedades() ([]*models.MenuPropiedades, error) {
	var propiedades []*models.MenuPropiedades
	query := service.sq.Select("Propiedades.id_propiedad", "Propiedades.titulo", "Propiedades.precio", "Propiedades.num_recamaras", "Estado_Propiedades.tipo_transaccion", "Estado_Propiedades.estado").
		From("Propiedades").
		Join("Estado_Propiedades ON Propiedades.id_propiedad = Estado_Propiedades.id_propiedad").
		Where(sq.Eq{"borrado_en": nil})
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
		log.Println("Error fetching all propiedades:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var propiedad models.MenuPropiedades
		err := rows.Scan(&propiedad.IDPropiedad, &propiedad.Titulo, &propiedad.Precio, &propiedad.Habitaciones, &propiedad.TipoTransaccion, &propiedad.Estado)
		if err != nil {
			log.Println("Error scanning propiedad:", err)
			return nil, err
		}
		propiedades = append(propiedades, &propiedad)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return propiedades, nil
}

func (service *PropiedadService) GetAllPropiedadesByPrice() ([]*models.MenuPropiedades, error) {
	var propiedades []*models.MenuPropiedades
	query := service.sq.Select("Propiedades.id_propiedad", "Propiedades.titulo", "Propiedades.precio", "Propiedades.num_recamaras", "Estado_Propiedades.tipo_transaccion", "Estado_Propiedades.estado").
		From("Propiedades").
		Join("Estado_Propiedades ON Propiedades.id_propiedad = Estado_Propiedades.id_propiedad").
		Where(sq.Eq{"borrado_en": nil}).
		OrderBy("Propiedades.precio DESC")
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
		log.Println("Error fetching all propiedades:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var propiedad models.MenuPropiedades
		err := rows.Scan(&propiedad.IDPropiedad, &propiedad.Titulo, &propiedad.Precio, &propiedad.Habitaciones, &propiedad.TipoTransaccion, &propiedad.Estado)
		if err != nil {
			log.Println("Error scanning propiedad:", err)
			return nil, err
		}
		propiedades = append(propiedades, &propiedad)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return propiedades, nil
}

func (service *PropiedadService) GetAllPropiedadesByBedrooms() ([]*models.MenuPropiedades, error) {
	var propiedades []*models.MenuPropiedades
	query := service.sq.Select("Propiedades.id_propiedad", "Propiedades.titulo", "Propiedades.precio", "Propiedades.num_recamaras", "Estado_Propiedades.tipo_transaccion", "Estado_Propiedades.estado").
		From("Propiedades").
		Join("Estado_Propiedades ON Propiedades.id_propiedad = Estado_Propiedades.id_propiedad").
		Where(sq.Eq{"borrado_en": nil}).
		OrderBy("Propiedades.num_recamaras DESC")
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
		log.Println("Error fetching all propiedades:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var propiedad models.MenuPropiedades
		err := rows.Scan(&propiedad.IDPropiedad, &propiedad.Titulo, &propiedad.Precio, &propiedad.Habitaciones, &propiedad.TipoTransaccion, &propiedad.Estado)
		if err != nil {
			log.Println("Error scanning propiedad:", err)
			return nil, err
		}
		propiedades = append(propiedades, &propiedad)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return propiedades, nil
}

// GET /propiedad/:id
// Funcion que recupera todos los campos de una propiedad en específico
func (service *PropiedadService) GetPropiedad(id int) (*models.Propiedad, error) {
	var propiedad models.Propiedad
	var gas, comodidades, extras, utilidades string
	query := service.sq.Select("*").From("Propiedades").Where(sq.Eq{"id_propiedad": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return nil, err
	}
	err = service.DB.QueryRow(sqlStr, args...).Scan(&propiedad.IDPropiedad, &propiedad.Titulo, &propiedad.FechaAlta,
		&propiedad.Direccion, &propiedad.Colonia, &propiedad.Ciudad,
		&propiedad.Referencia, &propiedad.Precio, &propiedad.MtsConstruccion,
		&propiedad.MtsTerreno, &propiedad.Habitada, &propiedad.Amueblada,
		&propiedad.NumPlantas, &propiedad.NumRecamaras, &propiedad.NumBanos,
		&propiedad.SizeCochera, &propiedad.MtsJardin, &gas,
		&comodidades, &extras, &utilidades,
		&propiedad.Observaciones, &propiedad.IDTipoPropiedad, &propiedad.IDPropietario, &propiedad.IDUsuario)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found")
			return nil, nil
		}
		log.Println("Error fetching propiedad:", err)
		return nil, err
	}

	propiedad.Gas = parseStringSet(gas)
	propiedad.Comodidades = parseStringSet(comodidades)
	propiedad.Extras = parseStringSet(extras)
	propiedad.Utilidades = parseStringSet(utilidades)

	return &propiedad, nil
}

func (service *PropiedadService) InsertPropiedad(propiedad *models.Propiedad, estado *models.EstadoPropiedades) (int, int, error) {
	utils := database.NewDbUtilities(service.DB)
	lastID, err := utils.GetLastId("Propiedades", "id_propiedad")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return 0, 0, err
	}
	propiedad.IDPropiedad = lastID + 1

	query := service.sq.Insert("Propiedades").
		Columns("id_propiedad", "titulo", "fecha_alta", "direccion", "colonia", "ciudad", "referencia",
			"precio", "mts_construccion", "mts_terreno", "habitada", "amueblada",
			"num_plantas", "num_recamaras", "num_banos", "size_cochera", "mts_jardin",
			"gas", "comodidades", "extras", "utilidades", "observaciones", "creado_en", "id_tipo_propiedad",
			"id_propietario", "usuario").
		Values(propiedad.IDPropiedad, propiedad.Titulo, propiedad.FechaAlta,
			propiedad.Direccion, propiedad.Colonia, propiedad.Ciudad, propiedad.Referencia,
			propiedad.Precio, propiedad.MtsConstruccion, propiedad.MtsTerreno, propiedad.Habitada, propiedad.Amueblada,
			propiedad.NumPlantas, propiedad.NumRecamaras, propiedad.NumBanos, propiedad.SizeCochera, propiedad.MtsJardin,
			strings.Join(propiedad.Gas, ","), strings.Join(propiedad.Comodidades, ","), strings.Join(propiedad.Extras, ","),
			strings.Join(propiedad.Utilidades, ","), propiedad.Observaciones, time.Now(), propiedad.IDTipoPropiedad, propiedad.IDPropietario, propiedad.IDUsuario)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, 0, err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error inserting propiedad:", err)
		return 0, 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, 0, err
	}
	if rows != 1 {

		log.Println("Error inserting estado: no rows affected")
		return 0, 0, err
	}

	estado.IDPropiedad = propiedad.IDPropiedad
	estado.IDEstadoPropiedades = lastID + 1
	query2 := service.sq.Insert("Estado_Propiedades").
		Columns("id_estado_propiedades", "tipo_transaccion", "estado", "fecha_cambio_estado", "id_propiedad").
		Values(estado.IDEstadoPropiedades, estado.TipoTransaccion, estado.Estado, estado.FechaTransaccion, estado.IDPropiedad)
	sqlStr2, args2, err := query2.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return 0, 0, err
	}
	result, err = service.DB.Exec(sqlStr2, args2...)
	if err != nil {
		log.Println("Error inserting estado de la propiedad:", err)
		return 0, 0, err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return 0, 0, err
	}
	if rows != 1 {
		log.Println("Error inserting estado: no rows affected")
		return 0, 0, err
	}
	return propiedad.IDPropiedad, estado.IDEstadoPropiedades, nil
}

// UpdatePropiedad updates a Propiedad in the database
func (service *PropiedadService) UpdatePropiedad(propiedad *models.Propiedad, id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Propiedades", "id_propiedad")
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	println(lastId)
	if id <= 0 || id > lastId {
		log.Println("Invalid propiedad ID:", id)
		return err
	}
	query := service.sq.Update("Propiedades").
		Set("titulo", propiedad.Titulo).
		Set("fecha_alta", propiedad.FechaAlta).
		Set("direccion", propiedad.Direccion).
		Set("colonia", propiedad.Colonia).
		Set("ciudad", propiedad.Ciudad).
		Set("referencia", propiedad.Referencia).
		Set("precio", propiedad.Precio).
		Set("mts_construccion", propiedad.MtsConstruccion).
		Set("mts_terreno", propiedad.MtsTerreno).
		Set("habitada", propiedad.Habitada).
		Set("amueblada", propiedad.Amueblada).
		Set("num_plantas", propiedad.NumPlantas).
		Set("num_recamaras", propiedad.NumRecamaras).
		Set("num_banos", propiedad.NumBanos).
		Set("size_cochera", propiedad.SizeCochera).
		Set("mts_jardin", propiedad.MtsJardin).
		Set("gas", strings.Join(propiedad.Gas, ",")).
		Set("comodidades", strings.Join(propiedad.Comodidades, ",")).
		Set("extras", strings.Join(propiedad.Extras, ",")).
		Set("utilidades", strings.Join(propiedad.Utilidades, ",")).
		Set("observaciones", propiedad.Observaciones).
		Set("actualizado_en", time.Now()).
		Set("id_tipo_propiedad", propiedad.IDTipoPropiedad).
		Set("id_propietario", propiedad.IDPropietario).
		Set("usuario", propiedad.IDUsuario).
		Where(sq.Eq{"id_propiedad": id}).
		Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error updating propiedad:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows != 1 {
		log.Println("Error updating propiedad: no rows affected")
		return err
	}
	return nil
}

// DeletePropiedad deletes a Propiedad from the database
func (service *PropiedadService) DeletePropiedad(id int) error {
	utils := database.NewDbUtilities(service.DB)
	lastId, err := utils.GetLastId("Propiedades", "id_propiedad")
	println(lastId)
	if err != nil {
		log.Println("Error getting last ID:", err)
		return err
	}
	if id <= 0 || id > lastId {
		log.Println("Invalid propiedad ID:", id)
		return err
	}

	query := service.sq.Update("Propiedades").Set("borrado_en", time.Now()).Where(sq.Eq{"id_propiedad": id}).Where(sq.Eq{"borrado_en": nil})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Println("Error building SQL query:", err)
		return err
	}
	result, err := service.DB.Exec(sqlStr, args...)
	if err != nil {
		log.Println("Error deleting propiedad:", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}
	if rows == 0 {
		log.Println("Error deleting propiedad: no rows affected")
		return err
	}
	return nil
}

// helper function to parse sets of strings
func parseStringSet(str string) []string {
	var set []string
	if str != "" {
		set = strings.Split(str, ",")
	}
	return set
}
