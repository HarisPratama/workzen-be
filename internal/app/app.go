package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workzen-be/config"
	"workzen-be/internal/adapter/cloudflare"
	"workzen-be/internal/adapter/handler"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/service"
	"workzen-be/lib/auth"
	"workzen-be/lib/middleware"
	"workzen-be/lib/pagination"
	"workzen-be/internal/ai"

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
		log.Fatalf("Could not connect to database: %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatalf("Could not create temp content directory: %v", err)
		return
	}

	// Cloudflare R2
	cfgR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cfgR2)
	r2Adapter := cloudflare.NewCloudflareR2Adapter(s3Client, cfg)

	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)
	quotaChecker := middleware.NewQuotaChecker(db.DB)

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
	interviewRepo := repository.NewInterviewRepository(db.DB)
	offerRepo := repository.NewOfferRepository(db.DB)
	assignmentRepo := repository.NewEmployeeAssignmentRepository(db.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(db.DB)
	overviewRepo := repository.NewOverviewRepository(db.DB)

	// AI Client
	aiAddr := os.Getenv("AI_SERVICE_ADDR")
	if aiAddr == "" {
		aiAddr = "localhost:50051"
	}
	aiClient, err := ai.NewClient(aiAddr)
	if err != nil {
		log.Printf("Warning: Could not connect to AI service: %v", err)
	}

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
	interviewService := service.NewInterviewService(interviewRepo, db.DB)
	offerService := service.NewOfferService(offerRepo)
	employeeAssignmentService := service.NewEmployeeAssignmentService(assignmentRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	overviewService := service.NewOverviewService(overviewRepo)
	hireService := service.NewHireService(db.DB)
	aiService := service.NewAIService(aiClient)
	jobPostingService := service.NewJobPostingService(db.DB, candidateRepo, candidateApplicationRepo, aiService)

	//handler
	authHandler := handler.NewAuthHandler(authService, cfg)
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
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	overviewHandler := handler.NewOverviewHandler(overviewService)
	hireHandler := handler.NewHireHandler(hireService)
	aiHandler := handler.NewAIHandler(aiService)
	jobPostingHandler := handler.NewJobPostingHandler(manpowerReqService, jobPostingService)

	app := fiber.New()
	corsOrigins := cfg.App.CorsOrigins
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
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

	// Subscription Plans Management (SUPER_ADMIN)
	subscriptionPlanApp := adminApp.Group("/subscription-plans")
	subscriptionPlanApp.Get("/", middlewareAuth.RequireRole("SUPER_ADMIN"), subscriptionHandler.GetPlans)
	subscriptionPlanApp.Post("/", middlewareAuth.RequireRole("SUPER_ADMIN"), subscriptionHandler.CreatePlan)
	subscriptionPlanApp.Get("/:planID", middlewareAuth.RequireRole("SUPER_ADMIN"), subscriptionHandler.GetPlanByID)
	subscriptionPlanApp.Put("/:planID", middlewareAuth.RequireRole("SUPER_ADMIN"), subscriptionHandler.UpdatePlan)
	subscriptionPlanApp.Delete("/:planID", middlewareAuth.RequireRole("SUPER_ADMIN"), subscriptionHandler.DeletePlan)

	// FE
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:contentID", contentHandler.GetContentDetail)

	// Tenant
	feApp.Post("/tenants/register", tenantHandler.RegisterTenant)

	// Subscription Plans (public - for pricing page)
	feApp.Get("/subscription-plans", subscriptionHandler.GetPlans)
	feApp.Get("/subscription-plans/:planID", subscriptionHandler.GetPlanByID)

	// Job Postings (public - no auth required)
	feApp.Get("/job-postings/:token", jobPostingHandler.GetJobPosting)
	feApp.Post("/job-postings/:token/apply", jobPostingHandler.ApplyToJob)

	tenantApp := api.Group("/v1")
	tenantApp.Use(middlewareAuth.CheckCookieToken())

	// overview
	tenantApp.Get("/overview", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), overviewHandler.GetOverview)

	// employee
	tenantApp.Get("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.GetEmployees)
	tenantApp.Post("/employees", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), quotaChecker.CheckQuota(middleware.ResourceEmployee), employeeHandler.CreateEmployee)
	tenantApp.Get("/employees/:employeeID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.GetEmployeeDetail)
	tenantApp.Put("/employees/:employeeID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.UpdateEmployee)
	tenantApp.Delete("/employees/:employeeID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), employeeHandler.DeleteEmployee)

	// client
	tenantApp.Get("/clients", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.GetClientByTenant)
	tenantApp.Post("/clients", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), quotaChecker.CheckQuota(middleware.ResourceClient), clientHandler.CreateClient)
	tenantApp.Get("/clients/:clientID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.GetClientDetailByTenant)
	tenantApp.Put("/clients/:clientID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.UpdateClient)
	tenantApp.Delete("/clients/:clientID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), clientHandler.DeleteClient)

	// manpower request
	tenantApp.Get("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.GetManpowerReqByTenant)
	tenantApp.Post("/manpower-requests", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), quotaChecker.CheckQuota(middleware.ResourceManpowerRequest), manpowerReqHandler.CreateManpowerReq)
	tenantApp.Get("/manpower-requests/:manpowerRequestID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.GetDetailManpowerRequestByTenant)
	tenantApp.Put("/manpower-requests/:manpowerRequestID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.UpdateManpowerReq)
	tenantApp.Delete("/manpower-requests/:manpowerRequestID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), manpowerReqHandler.DeleteManpowerReq)
	tenantApp.Post("/manpower-requests/:manpowerRequestID/generate-link", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), jobPostingHandler.GenerateLink)
	tenantApp.Get("/manpower-requests/:manpowerRequestID/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.GetCandidateApplicationByTenantMR)

	// candidate
	tenantApp.Get("/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.GetCandidatesByTenant)
	tenantApp.Post("/candidates", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.CreateCandidate)
	tenantApp.Get("/candidates/:candidateID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.GetCandidateDetailByTenant)
	tenantApp.Put("/candidates/:candidateID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.UpdateCandidate)
	tenantApp.Delete("/candidates/:candidateID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateHandler.DeleteCandidate)

	// candidate application
	tenantApp.Post("/candidate-applications", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.CreateCandidateApplication)
	tenantApp.Get("/candidate-applications/:applicationID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.GetCandidateApplicationDetail)
	tenantApp.Put("/candidate-applications/:applicationID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.UpdateCandidateApplication)
	tenantApp.Delete("/candidate-applications/:applicationID", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), candidateApplicationHandler.DeleteCandidateApplication)
	tenantApp.Post("/candidate-applications/:applicationID/hire", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"), hireHandler.HireCandidate)

	// Subscription (Tenant)
	subscriptionApp := tenantApp.Group("/subscriptions", middlewareAuth.RequireRole("TENANT_ADMIN"))
	subscriptionApp.Get("/active", subscriptionHandler.GetMySubscription)
	subscriptionApp.Get("/history", subscriptionHandler.GetSubscriptionHistory)
	subscriptionApp.Post("/subscribe", subscriptionHandler.Subscribe)
	subscriptionApp.Post("/:subscriptionID/cancel", subscriptionHandler.CancelSubscription)
	subscriptionApp.Post("/change-plan", subscriptionHandler.ChangePlan)

	// AI endpoints
	aiApp := tenantApp.Group("/ai", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	aiApp.Post("/analyze-cv", aiHandler.AnalyzeCV)
	aiApp.Post("/match-job", aiHandler.MatchJob)

	// ========== NEW ENDPOINTS ==========
	// Payroll endpoints
	payrollApp := tenantApp.Group("/payrolls", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	payrollApp.Get("/", payrollHandler.GetPayrollsByTenant)
	payrollApp.Post("/", payrollHandler.CreatePayroll)
	payrollApp.Get("/:payrollID", payrollHandler.GetPayrollByID)
	payrollApp.Put("/:payrollID", payrollHandler.UpdatePayroll)
	payrollApp.Delete("/:payrollID", payrollHandler.DeletePayroll)
	payrollApp.Get("/employee/:employeeID", payrollHandler.GetPayrollsByEmployee)
	payrollApp.Post("/:payrollID/process", payrollHandler.ProcessPayroll)
	payrollApp.Post("/:payrollID/mark-as-paid", payrollHandler.MarkAsPaid)

	// Attendance endpoints
	attendanceApp := tenantApp.Group("/attendances", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	attendanceApp.Get("/", attendanceHandler.GetAttendancesByTenant)
	attendanceApp.Post("/", attendanceHandler.CreateAttendance)
	attendanceApp.Get("/:attendanceID", attendanceHandler.GetAttendanceByID)
	attendanceApp.Put("/:attendanceID", attendanceHandler.UpdateAttendance)
	attendanceApp.Delete("/:attendanceID", attendanceHandler.DeleteAttendance)
	attendanceApp.Get("/employee/:employeeID", attendanceHandler.GetAttendancesByEmployee)
	attendanceApp.Post("/:attendanceID/check-in", attendanceHandler.CheckIn)
	attendanceApp.Post("/:attendanceID/check-out", attendanceHandler.CheckOut)
	attendanceApp.Get("/today/:employeeID", attendanceHandler.GetTodayAttendance)

	// Interview endpoints
	interviewApp := tenantApp.Group("/interviews", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	interviewApp.Get("/", interviewHandler.GetInterviews)
	interviewApp.Post("/", interviewHandler.CreateInterview)
	interviewApp.Get("/:interviewID", interviewHandler.GetInterviewByID)
	interviewApp.Put("/:interviewID", interviewHandler.UpdateInterview)
	interviewApp.Post("/:interviewID/feedback", interviewHandler.SubmitFeedback)
	interviewApp.Delete("/:interviewID", interviewHandler.DeleteInterview)

	// Offer endpoints
	offerApp := tenantApp.Group("/offers", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	offerApp.Get("/", offerHandler.GetOffers)
	offerApp.Post("/", offerHandler.CreateOffer)
	offerApp.Get("/:offerID", offerHandler.GetOfferByID)
	offerApp.Put("/:offerID", offerHandler.UpdateOffer)
	offerApp.Delete("/:offerID", offerHandler.DeleteOffer)

	// Employee Assignment endpoints
	assignmentApp := tenantApp.Group("/assignments", middlewareAuth.RequireRole("TENANT_ADMIN", "SUPERVISOR"))
	assignmentApp.Get("/", employeeAssignmentHandler.GetAssignments)
	assignmentApp.Post("/", employeeAssignmentHandler.CreateAssignment)
	assignmentApp.Get("/:assignmentID", employeeAssignmentHandler.GetAssignmentByID)
	assignmentApp.Put("/:assignmentID", employeeAssignmentHandler.UpdateAssignment)
	assignmentApp.Delete("/:assignmentID", employeeAssignmentHandler.DeleteAssignment)

	// Start server
	if cfg.App.AppPort == "" {
		cfg.App.AppPort = os.Getenv("APP_PORT")
	}
	log.Println("Starting server on port:", cfg.App.AppPort)

	go func() {
		if err := app.Listen(":" + cfg.App.AppPort); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server....")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
