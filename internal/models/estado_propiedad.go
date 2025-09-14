package models

import (
	"database/sql"
)

type EstadoPropiedades struct {
	IDEstadoPropiedades int          `json:"id_estado_propiedades"` // Clave primaria
	TipoTransaccion     string       `json:"tipo_transaccion"`      // Tipo de transacción ('venta', 'renta')
	Estado              string       `json:"estado"`                // Estado de la propiedad ('disponible', 'vendida', 'rentada')
	FechaTransaccion    sql.NullTime `json:"fecha_cambio_estado"`   // Fecha en que cambió el estado de la propiedad
	IDPropiedad         int          `json:"id_propiedad"`          // Clave foránea
}
