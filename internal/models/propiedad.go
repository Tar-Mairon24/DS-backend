package models

type Propiedad struct {
	IDPropiedad     int      `json:"id_propiedad"`      // Clave primaria
	Titulo          string   `json:"titulo"`            // Título de la propiedad
	FechaAlta       *string  `json:"fecha_alta"`        // Fecha en que se dio de alta la propiedad
	Direccion       string   `json:"direccion"`         // Dirección de la propiedad
	Colonia         string   `json:"colonia"`           // Colonia
	Ciudad          string   `json:"ciudad"`            // Ciudad
	Referencia      string   `json:"referencia"`        // Punto de referencia
	Precio          float64  `json:"precio"`            // Precio de la propiedad
	MtsConstruccion int      `json:"mts_construccion"`  // Metros cuadrados de construcción
	MtsTerreno      int      `json:"mts_terreno"`       // Metros cuadrados de terreno
	Habitada        bool     `json:"habitada"`          // Indica si la propiedad está habitada
	Amueblada       bool     `json:"amueblada"`         // Indica si la propiedad está amueblada
	NumPlantas      int      `json:"num_plantas"`       // Número de plantas
	NumRecamaras    int      `json:"num_recamaras"`     // Número de recámaras
	NumBanos        int      `json:"num_banos"`         // Número de baños
	SizeCochera     int      `json:"size_cochera"`      // Tamaño de la cochera (en número de autos)
	MtsJardin       int      `json:"mts_jardin"`        // Tamaño del jardín (en metros cuadrados)
	Gas             []string `json:"gas"`               // Tipo de gas ('estacionario', 'natural')
	Comodidades     []string `json:"comodidades"`       // Comodidades ('clima', 'calefaccion', 'hidroneumatico', 'aljibe', 'tinaco')
	Extras          []string `json:"extras"`            // Extras ('alberca', 'jardin', 'techada', 'cocineta', 'cuarto_servicio')
	Utilidades      []string `json:"utilidades"`        // Utilidades ('agua', 'luz', 'internet')
	Observaciones   string   `json:"observaciones"`     // Comentarios adicionales
	IDTipoPropiedad int      `json:"id_tipo_propiedad"` // Clave foránea que referencia a Tipo_Propiedad
	IDPropietario   int      `json:"id_propietario"`    // Clave foránea que referencia a Propietario
	IDUsuario       string   `json:"usuario"`           // Clave foránea que referencia a Usuarios
}

type MenuPropiedades struct {
	IDPropiedad     int     `json:"id_propiedad"` // Clave primaria
	Titulo          string  `json:"titulo"`       // Título de la propiedad
	Precio          float64 `json:"precio"`       // Precio de la propiedad
	Habitaciones    int     `json:"num_recamaras"`
	TipoTransaccion string  `json:"tipo_transaccion"` // Tipo de transacción ('venta', 'renta')
	Estado          string  `json:"estado"`           // Estado de la propiedad ('disponible', 'vendida', 'rentada')
}
