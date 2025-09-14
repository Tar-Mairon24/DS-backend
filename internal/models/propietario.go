package models

type Propietario struct {
	IDPropietario int    `json:"id_propietario"` // Clave primaria
	Nombre        string `json:"nombre"`         // Nombre del propietario
	ApellidoP     string `json:"apellido_p"`     // Apellido paterno
	ApellidoM     string `json:"apellido_m"`     // Apellido materno
	Telefono      string `json:"telefono"`       // Teléfono de contacto
	Correo        string `json:"correo"`         // Correo electrónico
}