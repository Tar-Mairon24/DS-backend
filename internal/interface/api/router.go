package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"ds-backend/cmd/di"
)

func SetupRouter(handlers *di.Handlers, services di.Services, middleware di.Middleware) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, handlers.HealthHandler)
		setupAuthRoutes(v1, handlers.AuthHandler)
		setupUserRoutes(v1, handlers.UserHandler, services.JwtService, middleware.AuthMiddleware)
		setupPropertyRoutes(v1, handlers.PropertyHandler, services.JwtService, middleware.AuthMiddleware)
		setupEmailRoutes(v1, handlers.EmailHandler)
		setupImageRoutes(v1, handlers.ImageHandler, services.JwtService, middleware.AuthMiddleware)
	}

	return r
}
