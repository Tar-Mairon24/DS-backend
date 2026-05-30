package di

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/ports"
	"ds-backend/internal/infrastructure/db"
	"ds-backend/internal/infrastructure/repository"
	"ds-backend/internal/infrastructure/service"
	"ds-backend/internal/interface/api/handler"
	"ds-backend/internal/usecase"
	"ds-backend/middleware"
)

type Container struct {
	SqlDB *sql.DB

	userRepo        ports.UserRepository
	propertyRepo    ports.PropertyRepository
	tokenRepo       ports.TokenRepository
	emailRepo       ports.EmailRepository
	imageRepo       ports.ImageRepository
	appointmentRepo ports.AppointmentRepository

	userUsecase        ports.UserUseCase
	propertyUsecase    ports.PropertyUseCase
	authUsecase        ports.AuthUseCase
	imageUsecase       ports.ImageUseCase
	appointmentUsecase ports.AppointmentUseCase

	jwtService   ports.JWTService
	emailService ports.EmailService

	userHandler        *handler.UserHandler
	propertyHandler    *handler.PropertyHandler
	healthHandler      *handler.HealthHandler
	authHandler        *handler.AuthHandler
	emailHandler       *handler.EmailHandler
	imageHandler       *handler.ImageHandler
	appointmentHandler *handler.AppointmentHandler
	hashing            middleware.HashingInterface
	authMiddleware     middleware.AuthMiddlewareInterface
}

func NewContainer() *Container {
	logrus.Info("Initializing DI container")

	container := &Container{}

	db.Init()
	container.SqlDB = db.GetDB()

	if container.SqlDB == nil {
		logrus.Fatal("Failed to initialize database connection")
	}

	// middleware
	container.hashing = middleware.NewHashing()
	container.authMiddleware = middleware.NewAuthMiddleware()

	// repos
	container.userRepo = repository.NewUserRepository(container.SqlDB)
	container.propertyRepo = repository.NewPropertyRepository(container.SqlDB)
	container.tokenRepo = repository.NewTokenRepository(container.SqlDB)
	container.emailRepo = repository.NewEmailRepository(container.SqlDB)
	container.imageRepo = repository.NewImageRepository(container.SqlDB)
	container.appointmentRepo = repository.NewAppointmentRepository(container.SqlDB)
	
	// services
	container.jwtService = service.NewJWTService(container.userRepo)
	container.emailService = service.NewEmailService(container.emailRepo)

	// usecases
	container.userUsecase = usecase.NewUserUseCase(container.userRepo, container.hashing)
	container.propertyUsecase = usecase.NewPropertyUseCase(container.propertyRepo, container.userRepo, container.imageRepo)
	container.authUsecase = usecase.NewAuthUseCase(container.userRepo, container.tokenRepo, container.jwtService, container.hashing, container.authMiddleware)
	container.imageUsecase = usecase.NewImageUseCase(container.imageRepo)
	container.appointmentUsecase = usecase.NewAppointmentUseCase(container.appointmentRepo, container.userRepo, container.propertyRepo)

	// handlers
	container.userHandler = handler.NewUserHandler(container.userUsecase)
	container.propertyHandler = handler.NewPropertyHandler(container.propertyUsecase)
	container.authHandler = handler.NewAuthHandler(container.jwtService, container.authUsecase)
	container.emailHandler = handler.NewEmailHandler(container.emailService)
	container.imageHandler = handler.NewImageHandler(container.imageUsecase)
	container.appointmentHandler = handler.NewAppointmentHandler(container.appointmentUsecase)
	container.healthHandler = handler.NewHealthHandler()

	logrus.Info("DI container initialized successfully")
	return container
}

type Handlers struct {
	PropertyHandler    *handler.PropertyHandler
	UserHandler        *handler.UserHandler
	HealthHandler      *handler.HealthHandler
	AuthHandler        *handler.AuthHandler
	AppointmentHandler *handler.AppointmentHandler
	ImageHandler       *handler.ImageHandler
	EmailHandler       *handler.EmailHandler
}

func (c *Container) GetHandlers() *Handlers {
	return &Handlers{
		PropertyHandler:    c.propertyHandler,
		UserHandler:        c.userHandler,
		HealthHandler:      c.healthHandler,
		AuthHandler:        c.authHandler,
		EmailHandler:       c.emailHandler,
		ImageHandler:       c.imageHandler,
		AppointmentHandler: c.appointmentHandler,
	}
}

type Services struct {
	JwtService   ports.JWTService
	EmailService ports.EmailService
}

func (c *Container) GetServices() Services {
	return Services{
		JwtService:   c.jwtService,
		EmailService: c.emailService,
	}
}

type Middleware struct {
	Hashing        middleware.HashingInterface
	AuthMiddleware middleware.AuthMiddlewareInterface
}

func (c *Container) GetMiddleware() Middleware {
	return Middleware{
		Hashing:        c.hashing,
		AuthMiddleware: c.authMiddleware,
	}
}
