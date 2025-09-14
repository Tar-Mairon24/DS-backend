package models

type Prospecto struct {
	IdCliente int    `json:"id_cliente"`                 // Clave primaria
	Nombre    string `json:"nombre_prospecto"`           // Nombre del propietario
	ApellidoP string `json:"apellido_paterno_prospecto"` // Apellido paterno
	ApellidoM string `json:"apellido_materno_prospecto"` // Apellido materno
	Telefono  string `json:"telefono_prospecto"`         // Teléfono de contacto
	Correo    string `json:"correo_prospecto"`           // Correo electrónico
}
