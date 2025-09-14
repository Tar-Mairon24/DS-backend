package models

type TipoPropiedad struct {
	IDTipoPropiedad int    `json:"id_tipo_propiedad"` // Clave primaria
	Tipo_Propiedad  string `json:"descripcion"`       // Descripción del tipo de propiedad
}