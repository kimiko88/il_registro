package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"registro-backend/internal/admin"
	"registro-backend/internal/agenda"
	"registro-backend/internal/attendance"
	"registro-backend/internal/auth"
	"registro-backend/internal/classes"
	"registro-backend/internal/colloqui"
	"registro-backend/internal/communications"
	"registro-backend/internal/config"
	"registro-backend/internal/db"
	"registro-backend/internal/didactic_materials"
	"registro-backend/internal/documents"
	"registro-backend/internal/extracurricular"
	"registro-backend/internal/grades"
	"registro-backend/internal/groups"
	"registro-backend/internal/handler"
	"registro-backend/internal/lessons"
	"registro-backend/internal/middleware"
	"registro-backend/internal/notes"
	"registro-backend/internal/notifications"
	"registro-backend/internal/orientamento"
	"registro-backend/internal/parents"
	"registro-backend/internal/pcto"
	"registro-backend/internal/postgres"
	"registro-backend/internal/reports"
	"registro-backend/internal/scheduling"
	"registro-backend/internal/schools"
	"registro-backend/internal/schoolsettings"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/search"
	"registro-backend/internal/signatures"
	"registro-backend/internal/subjects"
	"registro-backend/internal/teachers"
	"registro-backend/internal/tenants"
	"registro-backend/internal/textbooks"
	"registro-backend/internal/timetables"
	"registro-backend/internal/trips"
	"registro-backend/internal/users"
	"registro-backend/internal/verbali"
	"registro-backend/internal/ws"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"
	"registro-backend/pkg/upload"
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
	agendaRepo := agenda.NewRepository(database)
	colloquiRepo := colloqui.NewRepository(database)
	verbaliRepo := verbali.NewRepository(database)
	extraRepo := extracurricular.NewRepository(database)
	notifRepo := notifications.NewRepository(database)
	tripsRepo := trips.NewRepository(database)

	authMiddleware := auth.NewMiddleware(tokenManager, usersRepo)

	// 6. Setup Services
	authSvc := auth.NewService(authRepo, tokenManager, mfaService)
	usersSvc := users.NewService(usersRepo)
	classesSvc := classes.NewService(classesRepo)
	gradesSvc := grades.NewService(gradesRepo, usersRepo, database, wsHub)
	gradesAnalytics := grades.NewAnalyticsService(gradesRepo)
	// wsHub satisfies attendance.EventBroadcaster (BroadcastToUser + BroadcastToSchool).
	// CalendarService is nil: the attendance service falls back to CountDistinctDays.
	attendanceSvc := attendance.NewService(attendanceRepo, usersRepo, wsHub, nil)
	docsSvc := documents.NewService(docsRepo)
	schedSvc := scheduling.NewService(schedRepo, teachersRepo)
	pctoSvc := pcto.NewService(pctoRepo)
	orientSvc := orientamento.NewService(orientRepo)
	schoolsSvc := schools.NewService(schoolsRepo)
	signaturesSvc := signatures.NewService(signatures.NewRepository(database), docsSvc, usersRepo)

	commsSvc := communications.NewService(commsRepo)
	notesSvc := notes.NewService(notesRepo, usersRepo)
	adminSvc := admin.NewService(adminRepo)
	agendaSvc := agenda.NewService(agendaRepo)
	colloquiSvc := colloqui.NewService(colloquiRepo)
	verbaliSvc := verbali.NewService(verbaliRepo)
	extraSvc := extracurricular.NewService(extraRepo)
	notifSvc := notifications.NewService(notifRepo)
	tripsSvc := trips.NewService(tripsRepo)

	// 7. Setup Handlers
	authH := auth.NewHandler(authSvc)
	usersH := users.NewHandler(usersSvc)
	classesH := classes.NewHandler(classesSvc)
	gradesH := grades.NewHandler(gradesSvc, gradesAnalytics)
	attendanceH := attendance.NewHandler(attendanceSvc)
	docsUploader := upload.NewSupabaseUploader(cfg.Supabase.URL, cfg.Supabase.Key, cfg.Supabase.Bucket)
	docsH := documents.NewHandler(docsSvc, docsUploader)
	schedH := scheduling.NewHandler(schedSvc)
	signaturesH := signatures.NewHandler(signaturesSvc)

	notesH := notes.NewHandler(notesSvc)
	adminH := admin.NewHandler(adminSvc)
	timetablesH := timetables.NewHandler(timetablesRepo)
	agendaH := agenda.NewHandler(agendaSvc)
	colloquiH := colloqui.NewHandler(colloquiSvc)
	verbaliH := verbali.NewHandler(verbaliSvc)
	extraH := extracurricular.NewHandler(extraSvc)
	notifH := notifications.NewHandler(notifSvc)
	tripsH := trips.NewHandler(tripsSvc)

	wsHandler := ws.NewHandler(wsHub)

	adminMiddleware := admin.NewMiddleware()
	healthH := handler.NewHealthHandler(database)

	// 8. Setup Router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	middleware.InitCircuitBreaker()

	api := r.Group("/api/v1")
	{
		r.GET("/health", healthH.Health)
		r.GET("/ready", healthH.Ready)
		r.GET("/metrics", healthH.Metrics)

		authH.RegisterRoutes(api, authMiddleware)

		api.GET("/ws", authMiddleware.Authenticate(), func(c *gin.Context) {
			wsHandler.Listen(c)
		})

		protected := api.Group("/")
		protected.Use(authMiddleware.Authenticate())
		{
			usersGroup := protected.Group("/users")
			{
				usersGroup.GET("/me/children", usersH.GetMyChildren)
				usersGroup.POST("", usersH.Create)
				usersGroup.GET("", usersH.List)
				usersGroup.GET("/:id", usersH.Get)
				usersGroup.PATCH("/:id", usersH.Update)
				usersGroup.DELETE("/:id", usersH.Delete)
				usersGroup.POST("/:id/restore", usersH.Restore)
				usersGroup.POST("/bulk-import", usersH.BulkImport)
				usersGroup.POST("/bulk-delete", adminMiddleware.RequireAdminOrSuperAdmin(), usersH.BulkDelete)
				usersGroup.POST("/:id/change-password", usersH.ChangePassword)
				usersGroup.POST("/:id/reset-password", usersH.ForceResetPassword)
				usersGroup.PATCH("/:id/roles", adminMiddleware.RequireAdminOrSuperAdmin(), usersH.AssignRoles)
				usersGroup.GET("/:id/audit-log", usersH.GetAuditLog)
				usersGroup.POST("/:id/gdpr-export", usersH.ExportGDPR)
				usersGroup.DELETE("/:id/gdpr-delete", adminMiddleware.RequireAdminOrSuperAdmin(), usersH.DeleteGDPR)
				usersGroup.GET("/search", usersH.List)
				usersGroup.PATCH("/:id/disable-mfa", usersH.DisableMFA)
				usersGroup.GET("/:id/guardians", usersH.GetGuardians)
				usersGroup.POST("/:id/guardians", usersH.AddGuardian)
				usersGroup.DELETE("/:id/guardians/:guardianId", usersH.RemoveGuardian)
				usersGroup.POST("/me/switch-child/:studentId", usersH.SwitchChild)
			}

			gradesH.RegisterRoutes(protected)
			attendanceH.RegisterRoutes(protected)
			docsH.RegisterRoutes(protected)
			schedH.RegisterRoutes(protected)
			timetablesH.RegisterRoutes(protected)
			agendaH.RegisterRoutes(protected)
			colloquiH.RegisterRoutes(protected)
			verbaliH.RegisterRoutes(protected)
			extraH.RegisterRoutes(protected)
			notifH.RegisterRoutes(protected)
			tripsH.RegisterRoutes(protected)

			pctoH := pcto.NewHandler(pctoSvc)
			pctoH.RegisterRoutes(protected)

			lessonsRepo := lessons.NewRepository(database)
			lessonsSvc := lessons.NewService(lessonsRepo)
			lessonsH := lessons.NewHandler(lessonsSvc)
			lessonsH.RegisterRoutes(protected)

			materialsRepo := didactic_materials.NewRepository(database)
			materialsSvc := didactic_materials.NewService(materialsRepo, usersRepo)
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

			groupsRepo := groups.NewRepository(database)
			groupsSvc := groups.NewService(groupsRepo)
			groupsH := groups.NewHandler(groupsSvc)
			groupsH.RegisterRoutes(protected)

			schoolSettingsRepo := schoolsettings.NewRepository(database)
			schoolSettingsSvc := schoolsettings.NewService(schoolSettingsRepo)
			schoolSettingsH := schoolsettings.NewHandler(schoolSettingsSvc)
			schoolSettingsH.RegisterRoutes(protected)

			reportsSvc := reports.NewService(scrutinySvc)
			reportsH := reports.NewHandler(reportsSvc)
			reportsH.RegisterRoutes(protected)

			searchRepo := search.NewRepository(database)
			searchSvc := search.NewService(searchRepo)
			searchH := search.NewHandler(searchSvc)
			searchH.RegisterRoutes(protected)

			parentsRepo := parents.NewRepository(database)
			parentsSvc := parents.NewService(parentsRepo, usersRepo, gradesRepo, attendanceRepo, commsRepo)
			parentsH := parents.NewHandler(parentsSvc)
			parentsH.RegisterRoutes(protected)

			tenantsRepo := tenants.NewRepository(database)
			tenantsSvc := tenants.NewService(tenantsRepo)
			tenantsH := tenants.NewHandler(tenantsSvc)
			tenantsH.RegisterRoutes(protected)

			adminH.RegisterRoutes(protected, adminMiddleware)
			signaturesGroup := protected.Group("/signatures")
			signaturesGroup.POST("/", signaturesH.SignDocument)
			signaturesGroup.GET("/:id", signaturesH.GetSignatures)
		}
	}

	// 9. Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Log.Infof("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	logger.Log.Info("Shutdown signal received — draining in-flight requests (max 15s)...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Errorf("Graceful shutdown failed: %v", err)
		os.Exit(1)
	}

	logger.Log.Info("Server stopped cleanly.")
}
