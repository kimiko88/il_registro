package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"github.com/gin-gonic/gin"

	"registro-backend/internal/admin"
	"registro-backend/internal/attendance"
	"registro-backend/internal/auth"
	"registro-backend/internal/classes"
	"registro-backend/internal/communications"
// ...
	usersRepo := users.NewRepository(database)
    classesRepo := classes.NewRepository(database) // New
	gradesRepo := grades.NewRepository(database)
// ...
	usersSvc := users.NewService(usersRepo)
    classesSvc := classes.NewService(classesRepo) // New
	gradesSvc := grades.NewService(gradesRepo, usersRepo, database)
// ...
	usersH := users.NewHandler(usersSvc)
    classesH := classes.NewHandler(classesSvc) // New
	gradesH := grades.NewHandler(gradesSvc, gradesAnalytics)
// ...
			commsH := communications.NewHandler(commsSvc)
			commsH.RegisterRoutes(protected)

            // Classes route
            classesH.RegisterRoutes(protected)

			// Admin routes
			adminH.RegisterRoutes(protected, adminMiddleware)	"registro-backend/internal/config"
	"registro-backend/internal/db"
	"registro-backend/internal/documents"
	"registro-backend/internal/grades"
	"registro-backend/internal/handler"
	"registro-backend/internal/middleware"
	"registro-backend/internal/orientamento"
	"registro-backend/internal/pcto"
	"registro-backend/internal/postgres"
	"registro-backend/internal/scheduling"
	"registro-backend/internal/schools"
	"registro-backend/internal/users"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.Init(cfg.Server.Mode)

	// 3. Connect DB
	database, err := db.Connect(cfg.Database)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 4. Setup Authentication (Keys & Managers)
	// TODO: In production, load keys from files or KMS
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate RSA keys: %v", err)
	}
	publicKey := &privateKey.PublicKey

	tokenManager := jwt.NewTokenManager(privateKey, publicKey)
	mfaService := auth.NewMFAService("RegistroElettronico")
	authMiddleware := auth.NewMiddleware(tokenManager)

	// 5. Setup Repositories
	authRepo := auth.NewRepository(database)
	usersRepo := users.NewRepository(database)
	gradesRepo := grades.NewRepository(database)
	attendanceRepo := attendance.NewRepository(database)
	docsRepo := documents.NewRepository(database)
	schedRepo := scheduling.NewRepository(database)
	pctoRepo := pcto.NewRepository(database)
	orientRepo := orientamento.NewRepository(database)
	schoolsRepo := schools.NewRepository(database)
	commsRepo := communications.NewRepository(database)
	adminRepo := postgres.NewAdminRepository(database)

	// 6. Setup Services
	authSvc := auth.NewService(authRepo, tokenManager, mfaService)
	usersSvc := users.NewService(usersRepo)
	gradesSvc := grades.NewService(gradesRepo, usersRepo, database)
	gradesAnalytics := grades.NewAnalyticsService(gradesRepo)
	attendanceSvc := attendance.NewService(attendanceRepo)
	docsSvc := documents.NewService(docsRepo)
	schedSvc := scheduling.NewService(schedRepo)
	pctoSvc := pcto.NewService(pctoRepo)
	orientSvc := orientamento.NewService(orientRepo)
	schoolsSvc := schools.NewService(schoolsRepo)
	commsSvc := communications.NewService(commsRepo)
	adminSvc := admin.NewService(adminRepo)

	// 7. Setup Handlers
	authH := auth.NewHandler(authSvc)
	usersH := users.NewHandler(usersSvc)
	gradesH := grades.NewHandler(gradesSvc, gradesAnalytics)
	attendanceH := attendance.NewHandler(attendanceSvc)
	docsH := documents.NewHandler(docsSvc)
	schedH := scheduling.NewHandler(schedSvc)
	adminH := admin.NewHandler(adminSvc)
	adminMiddleware := admin.NewMiddleware()
	// ...
	healthH := handler.NewHealthHandler(database)

	// 8. Setup Router
	r := gin.New() // Use New() to control middleware order explicitly
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware()) // New Rate Limit

	// Initialize Circuit Breaker
	middleware.InitCircuitBreaker()

	api := r.Group("/api/v1")
	{
		// Health Checks
		r.GET("/health", healthH.Health)
		r.GET("/ready", healthH.Ready)

		// Auth Routes
		authH.RegisterRoutes(api, authMiddleware)

		// Protected routes
		protected := api.Group("/")
		protected.Use(authMiddleware.Authenticate())
		// Optional: Apply circuit breaker to specific routes if they call external services like Supabase
		// protected.Use(middleware.CircuitBreakerMiddleware())
		{
			usersGroup := protected.Group("/users")
			{
				// Admin Endpoints
				usersGroup.GET("/me/children", usersH.GetMyChildren) // New: Parent's children
				usersGroup.POST("/", usersH.Create)
				usersGroup.GET("/", usersH.List)
				usersGroup.GET("/:id", usersH.Get)
				usersGroup.PATCH("/:id", usersH.Update)
				usersGroup.DELETE("/:id", usersH.Delete)
				usersGroup.POST("/:id/restore", usersH.Restore)
				usersGroup.POST("/bulk-import", usersH.BulkImport)
				usersGroup.POST("/:id/change-password", usersH.ChangePassword)
				usersGroup.POST("/:id/reset-password", usersH.ForceResetPassword)
				usersGroup.PATCH("/:id/roles", usersH.AssignRoles)
				usersGroup.GET("/:id/audit-log", usersH.GetAuditLog)
				usersGroup.POST("/:id/gdpr-export", usersH.ExportGDPR)
				usersGroup.DELETE("/:id/gdpr-delete", usersH.DeleteGDPR)
				usersGroup.GET("/search", usersH.List) // Merged into List logic
				usersGroup.PATCH("/:id/disable-mfa", usersH.DisableMFA)
			}

			gradesH.RegisterRoutes(protected)

			attendanceH.RegisterRoutes(protected)
			docsH.RegisterRoutes(protected)
			schedH.RegisterRoutes(protected)

			// New modules
			pctoH := pcto.NewHandler(pctoSvc)
			pctoH.RegisterRoutes(protected)

			orientH := orientamento.NewHandler(orientSvc)
			orientH.RegisterRoutes(protected)

			schoolsH := schools.NewHandler(schoolsSvc)
			schoolsH.RegisterRoutes(protected)

			commsH := communications.NewHandler(commsSvc)
			commsH.RegisterRoutes(protected)

			// Admin routes
			adminH.RegisterRoutes(protected, adminMiddleware)
		}
	}

	// 9. Run
	logger.Log.Infof("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Log.Fatalf("Server failed to start: %v", err)
	}
}
