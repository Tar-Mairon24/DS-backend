package models

type Imagen struct {
	IDImagen    int    `json:"id_imagen"`
	RutaImagen  string `json:"ruta_imagen"`
	Descripcion string `json:"descripcion_imagen"`
	Principal   bool   `json:"principal"`
	IDPropiedad int    `json:"id_propiedad"`
}
