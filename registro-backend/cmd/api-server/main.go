package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"github.com/gin-gonic/gin"

	"registro-backend/internal/attendance"
	"registro-backend/internal/auth"
	"registro-backend/internal/config"
	"registro-backend/internal/db"
	"registro-backend/internal/documents"
	"registro-backend/internal/grades"
	"registro-backend/internal/handler"
	"registro-backend/internal/middleware"
	"registro-backend/internal/scheduling"
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

	// 6. Setup Services
	authSvc := auth.NewService(authRepo, tokenManager, mfaService)
	usersSvc := users.NewService(usersRepo)
	gradesSvc := grades.NewService(gradesRepo)
	attendanceSvc := attendance.NewService(attendanceRepo)
	docsSvc := documents.NewService(docsRepo)
	schedSvc := scheduling.NewService(schedRepo)

	// 7. Setup Handlers
	authH := auth.NewHandler(authSvc)
	usersH := users.NewHandler(usersSvc)
	gradesH := grades.NewHandler(gradesSvc)
	attendanceH := attendance.NewHandler(attendanceSvc)
	docsH := documents.NewHandler(docsSvc)
	schedH := scheduling.NewHandler(schedSvc)
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
				usersGroup.GET("/", usersH.GetUsers)
				usersGroup.POST("/", usersH.CreateUser)
			}

			gradesGroup := protected.Group("/grades")
			{
				gradesGroup.GET("/", gradesH.GetGrades)
				gradesGroup.POST("/", gradesH.AddGrade)
			}

			api.Group("/attendance").GET("/", attendanceH.GetAttendance).POST("/", attendanceH.MarkAttendance)
			api.Group("/documents").GET("/", docsH.GetDocuments).POST("/", docsH.UploadDocument)
			api.Group("/schedule").GET("/", schedH.GetSchedule).POST("/", schedH.AddSchedule)
		}
	}

	// 9. Run
	logger.Log.Infof("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Log.Fatalf("Server failed to start: %v", err)
	}
}
