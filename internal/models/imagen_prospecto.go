package models

type ImagenProspecto struct {
	IDImagen    int    `json:"id_imagen"`
	RutaImagen  string `json:"ruta_imagen"`
	Descripcion string `json:"descripcion_imagen"`
	Principal   bool   `json:"principal"`
	IDProspecto int    `json:"id_prospecto"`
}
