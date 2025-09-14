package models

type DocumentoAnexo struct {
	IDDocumentoAnexo     int    `json:"id_documento_anexo" gorm:"primaryKey;autoIncrement"`
	RutaDocumento        string `json:"ruta_documento" gorm:"size:255"`
	DescripcionDocumento string `json:"descripcion_documento_anexo" gorm:"size:45"`
	IDPropiedad          int    `json:"id_propiedad" gorm:"not null"`
}
