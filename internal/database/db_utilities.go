package database

import (
	"database/sql"
	"fmt"
)

type Db_Utilities struct {
	db *sql.DB
}

func NewDbUtilities(db *sql.DB) *Db_Utilities {
	return &Db_Utilities{
		db: db,
	}
}

func (dbu *Db_Utilities) GetLastId(table string, idName string) (int, error) {
	var lastId sql.NullInt64 // Usamos sql.NullInt64 para manejar NULL
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", idName, table)
	err := dbu.db.QueryRow(query).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	if lastId.Valid { // Comprobamos si el valor no es NULL
		return int(lastId.Int64), nil
	}
	return 0, nil // Si es NULL, retornamos 0
}
