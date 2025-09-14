package main

import (
	"backend/internal/controllers"
	"backend/internal/database"
	"backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database.InitDB()

	// Initialize services
	userService := services.NewUserService(database.DB)
	PropiedadService := services.NewPropiedadService(database.DB)
	EstadoPropiedadService := services.NewEstadoPropiedadService(database.DB)
	PropietarioService := services.NewPropietarioService(database.DB)
	TipoPropiedadService := services.NewTipoPropiedadService(database.DB)
	citasService := services.NewCitasService(database.DB)
	prospectoService := services.NewProspectoService(database.DB)
	imagenesService := services.NewImagenesService(database.DB)
	imagenesProspectoService := services.NewImagenesProspectoService(database.DB)
	contratoService := services.NewContratosService(database.DB)
	documentosAnexosService := services.NewDocumentosAnexosService(database.DB)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	propiedadController := controllers.NewPropiedadController(PropiedadService, EstadoPropiedadService)
	estadoPropiedadController := controllers.NewEstadoPropiedadController(EstadoPropiedadService)
	propietarioController := controllers.NewPropietarioController(PropietarioService)
	tipoPropiedadController := controllers.NewTipoPropiedadController(TipoPropiedadService)
	citasController := controllers.NewCitasController(citasService)
	prospectoController := controllers.NewProspectoController(prospectoService)
	imagenesController := controllers.NewImagenesController(imagenesService)
	imagenesProspectoController := controllers.NewImagenesProspectoController(imagenesProspectoService)
	contratosController := controllers.NewContratosController(contratoService)
	documentosAnexosController := controllers.NewDocumentosAnexosController(documentosAnexosService)

	// Set up Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set the Content-Type header to application/json; charset=utf-8
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	router.RemoveExtraSlash = true

	// Define routes

	// rutas para un usuario en específico
	router.GET("/users/:id", userController.GetUser)
	// ruta para comprobar credenciales de un usuario
	router.POST("/login", userController.Login)

	// ruta para agarrar la informacion importante de todas las propiedades
	router.GET("/all/propiedades", propiedadController.GetAllPropiedades)
	// ruta para agarrar la informacion importante de todas las propiedades ordenadas por precio
	router.GET("/all/propiedadesByPrice", propiedadController.GetAllPropiedadesByPrice)
	// ruta para agarrar la informacion importante de todas las propiedades ordenadas por habitaciones
	router.GET("/all/propiedadesByBedrooms", propiedadController.GetAllPropiedadesByBedrooms)
	// ruta para agarrar toda la informacion de una propiedad en específico
	router.GET("/propiedad/:id", propiedadController.GetPropiedad)
	// ruta para agarrar toda la informacion de una propiedad en específico
	router.GET("/prospecto/:id", prospectoController.GetProspecto)
	// ruta para agarrar el estado de una propiedad en específico
	router.GET("/estadopropiedad/:id", estadoPropiedadController.GetEstadoPropiedad)
	// ruta para agarrar el propietario de una propiedad en específico
	router.GET("/propietario/:id", propietarioController.GetPropietario)
	// ruta para agarrar el tipo de propiedad
	router.GET("/tipopropiedad/:id", tipoPropiedadController.GetTipoPropiedad)
	// ruta para agarrar todas las citas de un usuario especifico
	router.GET("/all/citas/:id", citasController.GetAllCitas)
	// ruta para agarrar una cita en específico
	router.GET("/cita/:id", citasController.GetCita)
	router.GET("/all/citas/:id/:day", citasController.GetAllCitasDay)

	router.GET("/all/imagenes/propiedad/:id", imagenesController.GetImagenesByPropiedad)
	router.GET("/all/imagenes/principal/:id", imagenesController.GetImagenPrincipal)
	router.GET("/all/imagenesProspecto/principal/:id", imagenesProspectoController.GetImagenPrincipal)
	router.GET("/all/imagenes/prospecto/:id", imagenesProspectoController.GetImagenesByProspecto)

	// ruta para insertar una propiedad
	router.POST("/create/propiedad", propiedadController.CreatePropiedad)
	router.POST("/create/estadopropiedad", estadoPropiedadController.CreateEstadoPropiedad)
	router.POST("/create/propietario", propietarioController.CreatePropietario)
	router.POST("/create/tipoPropiedad", tipoPropiedadController.CreateTipoPropiedad)
	router.POST("/create/cita", citasController.InsertCita)
	router.POST("/create/prospecto", prospectoController.InsertProspecto)

	router.POST("/create/imagen", imagenesController.InsertImagen)
	router.POST("/create/imagenProspecto", imagenesProspectoController.InsertImagen)

	// ruta para actualizar una propiedad
	router.PUT("/update/propiedad/:id", propiedadController.UpdatePropiedad)
	// router.PUT("/update/estadoPropiedad/:id", estadoPropiedadController.UpdateEstadoPropiedad)
	router.PUT("/update/cita/:id", citasController.UpdateCita)
	router.PUT("/update/prospecto/:id", prospectoController.UpdateProspecto)

	// rutas para borrar una propiedad
	router.DELETE("/eliminar/propiedad/:id", propiedadController.DeletePropiedad)
	router.DELETE("/eliminar/estadoPropiedad/:id", estadoPropiedadController.DeleteEstadoPropiedad)
	router.DELETE("/eliminar/cita/:id", citasController.DeleteCita)
	router.DELETE("/eliminar/imagen/:id", imagenesController.DeleteImagen)

	router.GET("/contratos/:id", contratosController.GetContrato)
	router.GET("all/contratos", contratosController.GetContratos)
	router.GET("/contratos/propiedad/:id_propiedad", contratosController.GetContratosByPropiedad)
	router.POST("/contratos", contratosController.CreateContrato)
	router.PUT("/contratos/:id", contratosController.UpdateContrato)
	router.DELETE("/contratos/:id", contratosController.DeleteContrato)

	router.GET("/all/documentos_anexos/propiedad/:id", documentosAnexosController.GetDocumentosByPropiedad)
	router.GET("/documento_anexo/:id", documentosAnexosController.GetDocumentoAnexo)
	router.POST("/create/documento_anexo", documentosAnexosController.InsertDocumentoAnexo)

	// Start the server
	router.Run(":8080")
}
