package api

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"ds-backend/internal/domain/ports"
	"ds-backend/internal/interface/api/handler"
	"ds-backend/middleware"
)

func setupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHandler.UserLogin)            // POST /api/v1/auth/login
		auth.POST("/refresh-token", authHandler.RefreshToken) // POST /api/v1/auth/refresh-token
		auth.POST("/logout/:id", authHandler.UserLogout)      // POST /api/v1/auth/logout
		auth.GET("/status", authHandler.GetStatus)            // POST /api/v1/auth/status
	}
}

func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
	users := rg.Group("/users")
	users.POST("", userHandler.CreateUser) // POST /api/v1/users
	users.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		users.GET("", userHandler.GetUsers)              // GET /api/v1/users
		users.GET("/:id", userHandler.GetUserByID)       // GET /api/v1/users/:id
		users.POST("/clients", userHandler.CreateClient) // POST /api/v1/users/clients
		users.POST("/owners", userHandler.CreateOwner)   // POST /api/v1/users/owners
		users.PUT("/:id", userHandler.UpdateUser)        // PUT /api/v1/users
		users.DELETE("/:id", userHandler.DeleteUser)     // DELETE /api/v1/users/:id
	}
}

func setupPropertyRoutes(rg *gin.RouterGroup, propertyHandler *handler.PropertyHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
	properties := rg.Group("/properties")
	properties.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		properties.GET("", propertyHandler.GetProperties)         // GET /api/v1/properties
		properties.GET("/:id", propertyHandler.GetPropertyByID)   // GET /api/v1/properties/:id
		properties.GET("/:id/card", propertyHandler.GetPropertyCardByID) // GET /api/v1/properties/:id/card
		properties.POST("", propertyHandler.CreateProperty)       // POST /api/v1/properties
		properties.PUT("/:id", propertyHandler.UpdateProperty)    // PUT /api/v1/properties/:id
		properties.DELETE("/:id", propertyHandler.DeleteProperty) // DELETE /api/v1/properties/:id
	}
}

func setupEmailRoutes(rg *gin.RouterGroup, emailHandler *handler.EmailHandler) {
	email := rg.Group("/email")
	{
		email.POST("/send-verification", emailHandler.SendVerificationEmail)    // POST /api/v1/email/send-verification
		email.POST("/verify", emailHandler.VerifyEmail)                         // POST /api/v1/email/verify
		email.POST("/resend-verification", emailHandler.ResendVerificationCode) // POST /api/v1/email/resend-verification
	}
}

func setupImageRoutes(rg *gin.RouterGroup, imageHandler *handler.ImageHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
	propertyImages := rg.Group("properties/:id/images")
	propertyImages.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		propertyImages.POST("", imageHandler.SaveImage)                            // POST /api/v1/properties/:id/images
		propertyImages.GET("", imageHandler.GetImagesByPropertyID)                 // GET /api/v1/properties/:id/images
		propertyImages.GET("/main", imageHandler.GetMainImageByPropertyID)         // GET /api/v1/properties/:id/images/main
		propertyImages.PATCH("/:imageId/main", imageHandler.UpdateMainImageStatus) // PATCH /api/v1/properties/:id/images/:imageId/main
	}

	images := rg.Group("images")
	images.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		images.GET("/:id", imageHandler.GetImageByID)   // GET /api/v1/images/:id
		images.PUT("/:id", imageHandler.UpdateImage)    // PUT /api/v1/images/:id
		images.DELETE("/:id", imageHandler.DeleteImage) // DELETE /api/v1/images/:id
	}
}

func setupUploadServerRoutes(router *gin.Engine, jwtService ports.JWTService, authMiddleware middleware.AuthMiddlewareInterface) {
	uploads := router.Group("/uploads")
	uploads.Use(authMiddleware.JWTAuthMiddleware(jwtService))
	{
		uploads.GET("/*filepath", func(c *gin.Context) {
			filePath := filepath.Clean(c.Param("filepath"))

			if strings.Contains(filePath, "..") {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid path"})
				return
			}

			fs := middleware.SafeFS{Root: http.Dir("/app/uploads")}
			f, err := fs.Open(filePath)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
				return
			}
			defer f.Close()

			http.ServeContent(c.Writer, c.Request, filePath, time.Time{}, f.(io.ReadSeeker))
		})
	}
}

func setupAppointmentRoutes(rg *gin.RouterGroup, appointmentHandler *handler.AppointmentHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
	appointments := rg.Group("/appointments")
	appointments.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		appointments.GET("", appointmentHandler.GetAll)               // GET /api/v1/appointments
		appointments.GET("/:id", appointmentHandler.GetByID)          // GET /api/v1/appointments/:id
		appointments.GET("/calendar", appointmentHandler.GetCalendar) // GET /api/v1/appointments/calendar?view=day|week|month
		appointments.POST("", appointmentHandler.Create)              // POST /api/v1/appointments
		appointments.PUT("/:id", appointmentHandler.Update)           // PUT /api/v1/appointments/:id
		appointments.PATCH("/:id/reschedule", appointmentHandler.Reschedule) // PATCH /api/v1/appointments/:id/reschedule
		appointments.PATCH("/:id/status", appointmentHandler.UpdateStatus)   // PATCH /api/v1/appointments/:id/status
		appointments.DELETE("/:id", appointmentHandler.Delete)        // DELETE /api/v1/appointments/:id
	}
}

func setupHealthRoutes(rg *gin.RouterGroup, healthHandler *handler.HealthHandler) {
	health := rg.Group("/health")
	{
		health.Match([]string{"GET", "HEAD"}, "", healthHandler.RegisterHealthRoutes)                 // GET api/v1/health
		health.Match([]string{"GET", "HEAD"}, "/detailed", healthHandler.RegisterDetailedHealthRoute) // GET api/v1/health/detailed
		health.Match([]string{"GET", "HEAD"}, "/ping", healthHandler.RegisterPingRoute)               // GET api/v1/health/ping
	}
}
