package app

import (
	"bwanews/config"
	"bwanews/internal/adapter/cloudflare"
	"bwanews/internal/adapter/handler"
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/service"
	"bwanews/lib/auth"
	"bwanews/lib/middleware"
	"bwanews/lib/pagination"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("Could not connect to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatal("Could not create temp content directory: %v", err)
		return
	}

	// Cloudflare R2
	cfgR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cfgR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)

	_ = pagination.NewPagination()

	//repository
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	employeeRepo := repository.NewEmployeeRepository(db.DB)
	tenantRepo := repository.NewTenantRepository(db.DB)
	clientRepo := repository.NewClientRepository(db.DB)
	manpowerReqRepo := repository.NewManpowerReqRepository(db.DB)

	//service
	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, r2Adapter)
	userService := service.NewUserService(userRepo)
	employeeService := service.NewEmployeeService(employeeRepo, cfg, r2Adapter)
	tenantService := service.NewTenantService(tenantRepo)
	clientService := service.NewClientService(clientRepo)
	manpowerReqService := service.NewManpowerReqService(manpowerReqRepo)

	//handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	tenantHandler := handler.NewTenantHandler(tenantService)
	clientHandler := handler.NewClientHandler(clientService)
	manpowerReqHandler := handler.NewManpowerReqHandler(manpowerReqService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} - ${method}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	}

	api := app.Group("/api")

	api.Post("/login", authHandler.Login)

	adminApp := api.Group("/v1/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	// category
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategoryByID)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryByID)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategory)

	// content
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Put("/:contentID", contentHandler.UpdateContent)
	contentApp.Get("/:contentID", contentHandler.GetContentByID)
	contentApp.Delete("/:contentID", contentHandler.DeleteContent)
	contentApp.Post("/upload-image", contentHandler.UploadImageR2)

	// User
	userApp := adminApp.Group("/users")
	userApp.Get("/profile", userHandler.GetUserByID)
	userApp.Put("/update-password", userHandler.UpdatePassword)

	// Employee
	employeeApp := adminApp.Group("/employees")
	employeeApp.Get("/", middlewareAuth.RequireRole("SUPER_ADMIN"), employeeHandler.GetEmployees)

	// FE
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentID", contentHandler.GetContentDetail)

	// Tenant
	feApp.Post("/tenants/register", tenantHandler.RegisterTenant)

	tenantApp := api.Group("/v1")
	tenantApp.Use(middlewareAuth.CheckToken())

	// employee
	tenantApp.Get("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.GetEmployees)
	tenantApp.Post("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.CreateEmployee)

	// client
	tenantApp.Get("/client", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.GetClientByTenant)
	tenantApp.Post("/client", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.CreateClient)

	// manpower request
	tenantApp.Get("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.GetManpowerReqByTenant)
	tenantApp.Post("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.CreateManpowerReq)

	// Start server
	log.Println("Starting server on port:", cfg.App.AppPort)
	if cfg.App.AppPort == "" {
		cfg.App.AppPort = os.Getenv("APP_PORT")
	}

	err = app.Listen(":" + cfg.App.AppPort)
	if err != nil {
		log.Fatal("Could not start server: %v", err)
	}

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("app listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
