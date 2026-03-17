package di

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
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

	userRepo     ports.UserRepository
	propertyRepo ports.PropertyRepository
	tokenRepo    ports.TokenRepository
	emailRepo    ports.EmailRepository

	userUsecase     ports.UserUseCase
	propertyUsecase ports.PropertyUseCase
	authUsecase     ports.AuthUseCase

	jwtService ports.JWTService
	emailService ports.EmailService

	userHandler     *handler.UserHandler
	propertyHandler *handler.PropertyHandler
	healthHandler   *handler.HealthHandler
	authHandler     *handler.AuthHandler
	emailHandler   	*handler.EmailHandler

	hashing        middleware.HashingInterface
	authMiddleware middleware.AuthMiddlewareInterface
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

	// services
	container.jwtService = service.NewJWTService(container.userRepo)
	container.emailService = service.NewEmailService(container.emailRepo)

	// usecases
	container.userUsecase = usecase.NewUserUseCase(container.userRepo, container.hashing)
	container.propertyUsecase = usecase.NewPropertyUseCase(container.propertyRepo)
	container.authUsecase = usecase.NewAuthUseCase(container.userRepo, container.tokenRepo, container.jwtService, container.hashing, container.authMiddleware)

	// handlers
	container.userHandler = handler.NewUserHandler(container.userUsecase)
	container.propertyHandler = handler.NewPropertyHandler(container.propertyUsecase)
	container.authHandler = handler.NewAuthHandler(container.jwtService, container.authUsecase)
	container.emailHandler = handler.NewEmailHandler(container.emailService)
	container.healthHandler = handler.NewHealthHandler()

	err := container.seedUser()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to seed initial user")
	}

	logrus.Info("DI container initialized successfully")
	return container
}

type Handlers struct {
	PropertyHandler *handler.PropertyHandler
	UserHandler     *handler.UserHandler
	HealthHandler   *handler.HealthHandler
	AuthHandler     *handler.AuthHandler
	EmailHandler    *handler.EmailHandler
}

func (c *Container) GetHandlers() *Handlers {
	return &Handlers{
		PropertyHandler: c.propertyHandler,
		UserHandler:     c.userHandler,
		HealthHandler:   c.healthHandler,
		AuthHandler:     c.authHandler,
		EmailHandler:    c.emailHandler,
	}
}

type Services struct {
	JwtService ports.JWTService
	EmailService ports.EmailService
}

func (c *Container) GetServices() Services {
	return Services{
		JwtService: c.jwtService,
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

func (c *Container) seedUser() error {
	users, err := c.userRepo.GetAll()
	if err != nil {
		return err
	}
	if len(users) == 0 {
		password, err := c.hashing.HashPassword("12345678")
		if err != nil {
			return err
		}
		user := models.User{
			Username: "tarmairon",
			Email:    "tarmairon@prueba.com",
			Password: password,
		}
		_, err = c.userRepo.Create(&user)
		if err != nil {
			logrus.WithError(err).Error("Failed to seed initial user")
			return err
		}
		logrus.Info("Seeded initial user successfully")
	}
	logrus.Info("Users already exist, skipping seeding")

	return nil
}
