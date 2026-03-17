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
	candidateRepo := repository.NewCandidateRepository(db.DB)
	candidateApplicationRepo := repository.NewCandidateApplicationRepository(db.DB)

	// New repositories for payroll, attendance, interview, offer, and employee assignment
	payrollRepo := repository.NewPayrollRepository(db.DB)
	attendanceRepo := repository.NewAttendanceRepository(db.DB)

	//service
	authService := service.NewAuthService(authRepo, cfg, jwt)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, cfg, r2Adapter)
	userService := service.NewUserService(userRepo)
	employeeService := service.NewEmployeeService(employeeRepo, cfg, r2Adapter)
	tenantService := service.NewTenantService(tenantRepo)
	clientService := service.NewClientService(clientRepo)
	manpowerReqService := service.NewManpowerReqService(manpowerReqRepo)
	candidateService := service.NewCandidateService(candidateRepo)
	candidateApplicationService := service.NewCandidateApplicationService(candidateApplicationRepo)

	// New services for payroll, attendance, interview, offer, and employee assignment
	payrollService := service.NewPayrollService(payrollRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo)
	interviewService := service.NewInterviewService(candidateApplicationRepo)
	offerService := service.NewOfferService(candidateApplicationRepo)
	employeeAssignmentService := service.NewEmployeeAssignmentService(employeeRepo)

	//handler
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService)
	userHandler := handler.NewUserHandler(userService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	tenantHandler := handler.NewTenantHandler(tenantService)
	clientHandler := handler.NewClientHandler(clientService)
	manpowerReqHandler := handler.NewManpowerReqHandler(manpowerReqService)
	candidateHandler := handler.NewCandidateHandler(candidateService)
	candidateApplicationHandler := handler.NewCandidateApplicationHandler(candidateApplicationService)

	// New handlers for payroll, attendance, interview, offer, and employee assignment
	payrollHandler := handler.NewPayrollHandler(payrollService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	interviewHandler := handler.NewInterviewHandler(interviewService)
	offerHandler := handler.NewOfferHandler(offerService)
	employeeAssignmentHandler := handler.NewEmployeeAssignmentHandler(employeeAssignmentService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	}))
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

	authApp := api.Group("/v1/auth")
	authApp.Post("/login", authHandler.Login)
	authApp.Post("/logout", authHandler.Logout)
	authApp.Post("/refresh", authHandler.RefreshToken)

	adminApp := api.Group("/v1/admin")
	adminApp.Use(middlewareAuth.CheckCookieToken())

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
	tenantApp.Use(middlewareAuth.CheckCookieToken())

	// employee
	tenantApp.Get("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.GetEmployees)
	tenantApp.Post("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.CreateEmployee)

	// client
	tenantApp.Get("/clients", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.GetClientByTenant)
	tenantApp.Post("/clients", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.CreateClient)

	// manpower request
	tenantApp.Get("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.GetManpowerReqByTenant)
	tenantApp.Get("/manpower-requests/:manpowerRequestID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.GetDetailManpowerRequestByTenant)
	tenantApp.Get("/manpower_requests/:manpowerRequestID/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.GetCandidateApplicationByTenantMR)
	tenantApp.Post("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.CreateManpowerReq)

	// candidate
	tenantApp.Get("/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.GetCandidatesByTenant)
	tenantApp.Post("/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.CreateCandidate)

	// candidate application
	tenantApp.Post("/candidate-application", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.CreateCandidateApplication)

	// ========== NEW ENDPOINTS ==========
	// Payroll endpoints
	payrollApp := tenantApp.Group("/payrolls", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	payrollApp.Get("/", payrollHandler.GetPayrollsByTenant)
	payrollApp.Post("/", payrollHandler.CreatePayroll)
	payrollApp.Get("/:id", payrollHandler.GetPayrollByID)
	payrollApp.Put("/:id", payrollHandler.UpdatePayroll)
	payrollApp.Delete("/:id", payrollHandler.DeletePayroll)
	payrollApp.Get("/employee/:employeeId", payrollHandler.GetPayrollsByEmployee)
	payrollApp.Post("/:id/process", payrollHandler.ProcessPayroll)
	payrollApp.Post("/:id/mark-as-paid", payrollHandler.MarkAsPaid)

	// Attendance endpoints
	attendanceApp := tenantApp.Group("/attendances", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	attendanceApp.Get("/", attendanceHandler.GetAttendancesByTenant)
	attendanceApp.Post("/", attendanceHandler.CreateAttendance)
	attendanceApp.Get("/:id", attendanceHandler.GetAttendanceByID)
	attendanceApp.Put("/:id", attendanceHandler.UpdateAttendance)
	attendanceApp.Delete("/:id", attendanceHandler.DeleteAttendance)
	attendanceApp.Get("/employee/:employeeId", attendanceHandler.GetAttendancesByEmployee)
	attendanceApp.Post("/:id/check-in", attendanceHandler.CheckIn)
	attendanceApp.Post("/:id/check-out", attendanceHandler.CheckOut)
	attendanceApp.Get("/today", attendanceHandler.GetTodayAttendance)

	// Interview endpoints
	interviewApp := tenantApp.Group("/interviews", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	interviewApp.Get("/", interviewHandler.GetInterviewsByTenant)
	interviewApp.Post("/", interviewHandler.CreateInterview)
	interviewApp.Get("/:id", interviewHandler.GetInterviewByID)
	interviewApp.Put("/:id", interviewHandler.UpdateInterview)
	interviewApp.Delete("/:id", interviewHandler.DeleteInterview)
	interviewApp.Get("/candidate/:candidateApplicationId", interviewHandler.GetInterviewsByCandidate)
	interviewApp.Get("/interviewer/:interviewerId", interviewHandler.GetInterviewsByInterviewer)
	interviewApp.Post("/:id/reschedule", interviewHandler.RescheduleInterview)
	interviewApp.Post("/:id/cancel", interviewHandler.CancelInterview)
	interviewApp.Post("/:id/complete", interviewHandler.CompleteInterview)
	interviewApp.Post("/:id/feedback", interviewHandler.SubmitFeedback)

	// Offer endpoints
	offerApp := tenantApp.Group("/offers", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	offerApp.Get("/", offerHandler.GetOffersByTenant)
	offerApp.Post("/", offerHandler.CreateOffer)
	offerApp.Get("/:id", offerHandler.GetOfferByID)
	offerApp.Put("/:id", offerHandler.UpdateOffer)
	offerApp.Delete("/:id", offerHandler.DeleteOffer)
	offerApp.Get("/candidate/:candidateApplicationId", offerHandler.GetOffersByCandidate)
	offerApp.Post("/:id/send", offerHandler.SendOffer)
	offerApp.Post("/:id/withdraw", offerHandler.WithdrawOffer)
	offerApp.Post("/:id/accept", offerHandler.AcceptOffer)
	offerApp.Post("/:id/reject", offerHandler.RejectOffer)
	offerApp.Post("/:id/negotiate", offerHandler.NegotiateOffer)

	// Employee Assignment endpoints
	assignmentApp := tenantApp.Group("/assignments", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	assignmentApp.Get("/", employeeAssignmentHandler.GetAssignmentsByTenant)
	assignmentApp.Post("/", employeeAssignmentHandler.CreateAssignment)
	assignmentApp.Get("/:id", employeeAssignmentHandler.GetAssignmentByID)
	assignmentApp.Put("/:id", employeeAssignmentHandler.UpdateAssignment)
	assignmentApp.Delete("/:id", employeeAssignmentHandler.DeleteAssignment)
	assignmentApp.Get("/employee/:employeeId", employeeAssignmentHandler.GetAssignmentsByEmployee)
	assignmentApp.Get("/employee/:employeeId/active", employeeAssignmentHandler.GetActiveAssignmentsByEmployee)
	assignmentApp.Get("/project/:projectId", employeeAssignmentHandler.GetAssignmentsByProject)
	assignmentApp.Post("/:id/start", employeeAssignmentHandler.StartAssignment)
	assignmentApp.Post("/:id/end", employeeAssignmentHandler.EndAssignment)

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

	log.Println("Shutting down server....")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
