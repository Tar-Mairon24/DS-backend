package models

type Cita struct {
	IDCita      int    `json:"id_citas"`         // Clave primaria
	Titulo      string `json:"titulo_cita"`      // Título de la cita
	FechaCita   string `json:"fecha_cita"`       // Fecha de la cita
	HoraCita    int    `json:"hora_cita"`        // Hora de la cita
	Descripcion string `json:"descripcion_cita"` // Descripción de la cita
	IdUsuario   string `json:"usuario"`          // Clave foránea que referencia a Usuarios
	IdCliente   int    `json:"id_cliente"`       // Clave foránea que referencia a Clientes

}

type CitaMenu struct {
	IDCita                 int    `json:"id_citas"`                 // Clave primaria
	Titulo                 string `json:"titulo"`                   // Título de la cita
	FechaCita              string `json:"fecha_cita"`               // Fecha de la cita
	HoraCita               int    `json:"hora_cita"`                // Hora de la cita
	NombreCliente          string `json:"nombre_cliente"`           // Nombre del cliente
	ApellidoPaternoCliente string `json:"apellido_paterno_cliente"` // Apellido paterno del cliente
	ApellidoMaternoCliente string `json:"apellido_materno_cliente"` // Apellido materno del cliente
}
