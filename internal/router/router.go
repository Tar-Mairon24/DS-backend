package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/controllers"
	"backend/internal/database"
	"backend/internal/services"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Initialize services
	userService := services.NewUserService(database.DB)
	propiedadService := services.NewPropiedadService(database.DB)
	propietarioService := services.NewPropietarioService(database.DB)
	tipoPropiedadService := services.NewTipoPropiedadService(database.DB)
	citasService := services.NewCitasService(database.DB)
	prospectoService := services.NewProspectoService(database.DB)
	imagenesService := services.NewImagenesService(database.DB)
	imagenesProspectoService := services.NewImagenesProspectoService(database.DB)
	contratoService := services.NewContratosService(database.DB)
	documentosAnexosService := services.NewDocumentosAnexosService(database.DB)
	estadoPropiedadService := services.NewEstadoPropiedadService(database.DB)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	propiedadController := controllers.NewPropiedadController(propiedadService, estadoPropiedadService)
	estadoPropiedadController := controllers.NewEstadoPropiedadController(estadoPropiedadService)
	propietarioController := controllers.NewPropietarioController(propietarioService)
	tipoPropiedadController := controllers.NewTipoPropiedadController(tipoPropiedadService)
	citasController := controllers.NewCitasController(citasService)
	prospectoController := controllers.NewProspectoController(prospectoService)
	imagenesController := controllers.NewImagenesController(imagenesService)
	imagenesProspectoController := controllers.NewImagenesProspectoController(imagenesProspectoService)
	contratosController := controllers.NewContratosController(contratoService)
	documentosAnexosController := controllers.NewDocumentosAnexosController(documentosAnexosService)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	router.RemoveExtraSlash = true

	v1 := router.Group("/api/v1")

	authRoutes(v1, userController)
	userRoutes(v1, userController)
	propiedadRoutes(v1, propiedadController)
	propietarioRoutes(v1, propietarioController)
	tipoPropiedadRoutes(v1, tipoPropiedadController)
	estadoPropiedadRoutes(v1, estadoPropiedadController)
	prospectoRoutes(v1, prospectoController)
	imagenesProspectoRoutes(v1, imagenesProspectoController)
	citasRoutes(v1, citasController)
	contratosRoutes(v1, contratosController)
	imagenesRoutes(v1, imagenesController)
	documentosAnexosRoutes(v1, documentosAnexosController)

	return router
}

func authRoutes(group *gin.RouterGroup, userController *controllers.UserController) {
	group.POST("/login", userController.Login)
}

func userRoutes(group *gin.RouterGroup, userController *controllers.UserController) {
	users := group.Group("/users")
	{
		users.GET(":id", userController.GetUser)
		users.POST("/create", userController.CreateUser)
	}
}

func propiedadRoutes(group *gin.RouterGroup, propiedadController *controllers.Propiedad_Controller) {
	propiedades := group.Group("/propiedades")
	{
		propiedades.GET("/all", propiedadController.GetAllPropiedades)
		propiedades.GET("/all/propiedadesByPrice", propiedadController.GetAllPropiedadesByPrice)
		propiedades.GET("/all/propiedadesByBedrooms", propiedadController.GetAllPropiedadesByBedrooms)
		propiedades.GET("/:id", propiedadController.GetPropiedad)
		propiedades.POST("/create", propiedadController.CreatePropiedad)
		propiedades.PUT("/update/:id", propiedadController.UpdatePropiedad)
		propiedades.DELETE("/eliminar/:id", propiedadController.DeletePropiedad)
	}
}

func propietarioRoutes(group *gin.RouterGroup, propietarioController *controllers.PropietarioController) {
	propietarios := group.Group("/propietarios")
	{
		propietarios.GET("/:id", propietarioController.GetPropietario)
		propietarios.POST("/create", propietarioController.CreatePropietario)
	}
}

func prospectoRoutes(group *gin.RouterGroup, prospectoController *controllers.ProspectoController) {
	prospectos := group.Group("/prospectos")
	{
		prospectos.GET("/:id", prospectoController.GetProspecto)
		prospectos.POST("/create", prospectoController.InsertProspecto)
		prospectos.PUT("/update/:id", prospectoController.UpdateProspecto)
	}
}

func tipoPropiedadRoutes(group *gin.RouterGroup, tipoPropiedadController *controllers.TipoPropiedadController) {
	tipos := group.Group("/tipopropiedad")
	{
		tipos.GET("/:id", tipoPropiedadController.GetTipoPropiedad)
		tipos.POST("/create", tipoPropiedadController.CreateTipoPropiedad)
	}
}

func estadoPropiedadRoutes(group *gin.RouterGroup, estadoPropiedadController *controllers.EstadoPropiedadController) {
	estados := group.Group("/estadopropiedad")
	{
		estados.GET("/:id", estadoPropiedadController.GetEstadoPropiedad)
		estados.POST("/create", estadoPropiedadController.CreateEstadoPropiedad)
		estados.DELETE("/eliminar/:id", estadoPropiedadController.DeleteEstadoPropiedad)
	}
}
func imagenesProspectoRoutes(group *gin.RouterGroup, imagenesProspectoController *controllers.ImagenesProspectoController) {
	imagenes := group.Group("/imagenesProspecto")
	{
		imagenes.GET("/principal/:id", imagenesProspectoController.GetImagenPrincipal)
		imagenes.GET("/prospecto/:id", imagenesProspectoController.GetImagenesByProspecto)
		imagenes.POST("/create", imagenesProspectoController.InsertImagen)
	}
}
func citasRoutes(group *gin.RouterGroup, citasController *controllers.CitasController) {
	citas := group.Group("/citas")
	{
		citas.GET("/all/:id", citasController.GetAllCitas)
		citas.GET("/:id", citasController.GetCita)
		citas.GET("/all/:id/:day", citasController.GetAllCitasDay)
		citas.POST("/create", citasController.InsertCita)
		citas.PUT("/update/:id", citasController.UpdateCita)
		citas.DELETE("/eliminar/:id", citasController.DeleteCita)
	}
}
func contratosRoutes(group *gin.RouterGroup, contratosController *controllers.ContratosController) {
	contratos := group.Group("/contratos")
	{
		contratos.GET("/:id", contratosController.GetContrato)
		contratos.GET("/all", contratosController.GetContratos)
		contratos.GET("/propiedad/:id_propiedad", contratosController.GetContratosByPropiedad)
		contratos.POST("/", contratosController.CreateContrato)
		contratos.PUT("/:id", contratosController.UpdateContrato)
		contratos.DELETE("/:id", contratosController.DeleteContrato)
	}
}
func imagenesRoutes(group *gin.RouterGroup, imagenesController *controllers.ImagenesController) {
	imagenes := group.Group("/imagenes")
	{
		imagenes.GET("/all/propiedad/:id", imagenesController.GetImagenesByPropiedad)
		imagenes.GET("/all/principal/:id", imagenesController.GetImagenPrincipal)
		imagenes.POST("/create", imagenesController.InsertImagen)
		imagenes.DELETE("/eliminar/:id", imagenesController.DeleteImagen)
	}
}

func documentosAnexosRoutes(group *gin.RouterGroup, documentosAnexosController *controllers.DocumentosAnexosController){
	documentos := group.Group("/documentos_anexos")
	{
		documentos.GET("/all/propiedad/:id", documentosAnexosController.GetDocumentosByPropiedad)
		documentos.GET("/:id", documentosAnexosController.GetDocumentoAnexo)
		documentos.POST("/create", documentosAnexosController.InsertDocumentoAnexo)
	}
}