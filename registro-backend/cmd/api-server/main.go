package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"registro-backend/internal/admin"
	"registro-backend/internal/attendance"
	"registro-backend/internal/auth"
	"registro-backend/internal/classes"
	"registro-backend/internal/communications"
	"registro-backend/internal/config"
	"registro-backend/internal/db"
	"registro-backend/internal/didactic_materials"
	"registro-backend/internal/documents"
	"registro-backend/internal/grades"
	"registro-backend/internal/handler"
	"registro-backend/internal/lessons"
	"registro-backend/internal/middleware"
	"registro-backend/internal/notes"
	"registro-backend/internal/orientamento"
	"registro-backend/internal/pcto"
	"registro-backend/internal/postgres"
	"registro-backend/internal/scheduling"
	"registro-backend/internal/schools"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/signatures"
	"registro-backend/internal/subjects"
	"registro-backend/internal/teachers"
	"registro-backend/internal/textbooks"
	"registro-backend/internal/timetables"
	"registro-backend/internal/users"
	"registro-backend/internal/ws"
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
	// Keys are saved to disk to persist sessions across restarts
	privateKey, publicKey, err := jwt.GetOrGenerateKeys("private_key.pem", "public_key.pem")
	if err != nil {
		log.Fatalf("Failed to load/generate RSA keys: %v", err)
	}

	tokenManager := jwt.NewTokenManager(privateKey, publicKey)
	mfaService := auth.NewMFAService("RegistroElettronico")
	wsHub := ws.NewHub()
	go wsHub.Run()

	// 5. Setup Repositories
	authRepo := auth.NewRepository(database)
	usersRepo := users.NewRepository(database)
	classesRepo := classes.NewRepository(database)
	gradesRepo := grades.NewRepository(database)
	attendanceRepo := attendance.NewRepository(database)
	docsRepo := documents.NewRepository(database)
	schedRepo := scheduling.NewRepository(database)
	pctoRepo := pcto.NewRepository(database)
	orientRepo := orientamento.NewRepository(database)
	schoolsRepo := schools.NewRepository(database)
	commsRepo := communications.NewRepository(database)
	notesRepo := notes.NewRepository(database)
	adminRepo := postgres.NewAdminRepository(database)
	timetablesRepo := timetables.NewRepository(database)
	teachersRepo := teachers.NewRepository(database)

	authMiddleware := auth.NewMiddleware(tokenManager, usersRepo)

	// 6. Setup Services
	authSvc := auth.NewService(authRepo, tokenManager, mfaService)
	usersSvc := users.NewService(usersRepo)
	classesSvc := classes.NewService(classesRepo)
	gradesSvc := grades.NewService(gradesRepo, usersRepo, database, wsHub)
	gradesAnalytics := grades.NewAnalyticsService(gradesRepo)
	attendanceSvc := attendance.NewService(attendanceRepo, wsHub)
	docsSvc := documents.NewService(docsRepo)
	schedSvc := scheduling.NewService(schedRepo, teachersRepo)
	pctoSvc := pcto.NewService(pctoRepo)
	orientSvc := orientamento.NewService(orientRepo)
	schoolsSvc := schools.NewService(schoolsRepo)
	signaturesSvc := signatures.NewService(signatures.NewRepository(database), docsSvc, usersRepo)

	commsSvc := communications.NewService(commsRepo)
	notesSvc := notes.NewService(notesRepo)
	adminSvc := admin.NewService(adminRepo)

	// 7. Setup Handlers
	authH := auth.NewHandler(authSvc)
	usersH := users.NewHandler(usersSvc)
	classesH := classes.NewHandler(classesSvc)
	gradesH := grades.NewHandler(gradesSvc, gradesAnalytics)
	attendanceH := attendance.NewHandler(attendanceSvc)
	docsH := documents.NewHandler(docsSvc)
	schedH := scheduling.NewHandler(schedSvc)
	signaturesH := signatures.NewHandler(signaturesSvc)

	notesH := notes.NewHandler(notesSvc)
	adminH := admin.NewHandler(adminSvc)
	timetablesH := timetables.NewHandler(timetablesRepo)

	wsHandler := ws.NewHandler(wsHub)

	adminMiddleware := admin.NewMiddleware()
	// ...
	healthH := handler.NewHealthHandler(database)

	// 8. Setup Router
	r := gin.New() // Use New() to control middleware order explicitly
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
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

		// WebSocket Route
		api.GET("/ws", authMiddleware.Authenticate(), func(c *gin.Context) {
			wsHandler.Listen(c)
		})

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
				usersGroup.POST("", usersH.Create)
				usersGroup.GET("", usersH.List)
				usersGroup.GET("/:id", usersH.Get)
				usersGroup.PATCH("/:id", usersH.Update)
				usersGroup.DELETE("/:id", usersH.Delete)
				usersGroup.POST("/:id/restore", usersH.Restore)
				usersGroup.POST("/bulk-import", usersH.BulkImport)
				usersGroup.POST("/bulk-delete", usersH.BulkDelete)
				usersGroup.POST("/:id/change-password", usersH.ChangePassword)
				usersGroup.POST("/:id/reset-password", usersH.ForceResetPassword)
				usersGroup.PATCH("/:id/roles", usersH.AssignRoles)
				usersGroup.GET("/:id/audit-log", usersH.GetAuditLog)
				usersGroup.POST("/:id/gdpr-export", usersH.ExportGDPR)
				usersGroup.DELETE("/:id/gdpr-delete", usersH.DeleteGDPR)
				usersGroup.GET("/search", usersH.List) // Merged into List logic
				usersGroup.PATCH("/:id/disable-mfa", usersH.DisableMFA)
				usersGroup.GET("/:id/guardians", usersH.GetGuardians)
				usersGroup.POST("/:id/guardians", usersH.AddGuardian)
				usersGroup.DELETE("/:id/guardians/:guardianId", usersH.RemoveGuardian)
			}

			gradesH.RegisterRoutes(protected)

			attendanceH.RegisterRoutes(protected)
			docsH.RegisterRoutes(protected)
			schedH.RegisterRoutes(protected)
			timetablesH.RegisterRoutes(protected)

			// New modules
			pctoH := pcto.NewHandler(pctoSvc)
			pctoH.RegisterRoutes(protected)

			lessonsRepo := lessons.NewRepository(database)
			lessonsSvc := lessons.NewService(lessonsRepo)
			lessonsH := lessons.NewHandler(lessonsSvc)
			lessonsH.RegisterRoutes(protected)

			materialsRepo := didactic_materials.NewRepository(database)
			materialsSvc := didactic_materials.NewService(materialsRepo)
			materialsH := didactic_materials.NewHandler(materialsSvc)
			materialsH.RegisterRoutes(protected)

			orientH := orientamento.NewHandler(orientSvc)
			orientH.RegisterRoutes(protected)

			schoolsH := schools.NewHandler(schoolsSvc)
			schoolsH.RegisterRoutes(protected)

			commsH := communications.NewHandler(commsSvc)
			commsH.RegisterRoutes(protected)

			textbooksRepo := textbooks.NewRepository(database)
			textbooksSvc := textbooks.NewService(textbooksRepo)
			textbooksH := textbooks.NewHandler(textbooksSvc)
			textbooksH.RegisterRoutes(protected)

			scrutinyRepo := scrutiny.NewRepository(database)
			attendanceRepo := attendance.NewRepository(database)
			scrutinySvc := scrutiny.NewService(scrutinyRepo, gradesRepo, classesRepo, usersRepo, attendanceRepo)
			scrutinyH := scrutiny.NewHandler(scrutinySvc)
			scrutinyH.RegisterRoutes(protected)

			classesH.RegisterRoutes(protected)
			notesH.RegisterRoutes(protected)

			subjectsRepo := subjects.NewRepository(database)
			subjectsSvc := subjects.NewService(subjectsRepo)
			subjectsH := subjects.NewHandler(subjectsSvc)
			subjectsH.RegisterRoutes(protected)

			teachersSvc := teachers.NewService(teachersRepo)
			teachersH := teachers.NewHandler(teachersSvc)
			teachersH.RegisterRoutes(protected)

			// Admin routes
			adminH.RegisterRoutes(protected, adminMiddleware)
			signatures := protected.Group("/signatures")
			signatures.POST("/", signaturesH.SignDocument)
			signatures.GET("/:id", signaturesH.GetSignatures)
		}
	}

	// 9. Run
	logger.Log.Infof("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Log.Fatalf("Server failed to start: %v", err)
	}
}
