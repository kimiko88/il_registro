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
	"registro-backend/internal/auditlog"
	"registro-backend/internal/auth"
	"registro-backend/internal/certificates"
	"registro-backend/internal/classes"
	"registro-backend/internal/colloqui"
	"registro-backend/internal/communications"
	"registro-backend/internal/competencies"
	"registro-backend/internal/config"
	"registro-backend/internal/credits"
	"registro-backend/internal/db"
	"registro-backend/internal/didactic_materials"
	"registro-backend/internal/documents"
	"registro-backend/internal/elearning"
	"registro-backend/internal/extracurricular"
	"registro-backend/internal/general_meetings"
	"registro-backend/internal/grades"
	"registro-backend/internal/groups"
	"registro-backend/internal/handler"
	"registro-backend/internal/lessons"
	"registro-backend/internal/mailer"
	"registro-backend/internal/middleware"
	"registro-backend/internal/notes"
	"registro-backend/internal/notifications"
	"registro-backend/internal/orientamento"
	"registro-backend/internal/parents"
	"registro-backend/internal/payments"
	"registro-backend/internal/pcto"
	"registro-backend/internal/pdp"
	"registro-backend/internal/postgres"
	"registro-backend/internal/recovery"
	"registro-backend/internal/reports"
	"registro-backend/internal/rubrics"
	"registro-backend/internal/scheduling"
	"registro-backend/internal/schoolcalendar"
	"registro-backend/internal/schools"
	"registro-backend/internal/schoolsettings"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/search"
	"registro-backend/internal/signatures"
	"registro-backend/internal/student_goals"
	"registro-backend/internal/students"
	"registro-backend/internal/subjects"
	"registro-backend/internal/substitutions"
	"registro-backend/internal/support"
	"registro-backend/internal/teacher_activities"
	"registro-backend/internal/teachers"
	"registro-backend/internal/tenants"
	"registro-backend/internal/textbooks"
	"registro-backend/internal/timetables"
	"registro-backend/internal/trips"
	"registro-backend/internal/uda"
	"registro-backend/internal/users"
	"registro-backend/internal/verbali"
	"registro-backend/internal/ws"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"
	"registro-backend/pkg/upload"
	"registro-backend/pkg/wsticket"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.Init(cfg.Server.Mode)

	// 2a. Init rate limiter backend (Redis-backed if REDIS_URL is set, in-memory otherwise).
	// Must be called before the first request, so we do it right after logging is up.
	middleware.InitRateLimiter(os.Getenv("REDIS_URL"))

	// 3. Connect DB
	database, err := db.Connect(cfg.Database)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 4. Setup Authentication (Keys & Managers)
	privateKey, publicKey, err := jwt.GetOrGenerateKeys("private_key.pem", "public_key.pem")
	if err != nil {
		log.Fatalf("Failed to load/generate RSA keys: %v", err)
	}

	tokenManager := jwt.NewTokenManager(privateKey, publicKey)
	mfaService := auth.NewMFAService("RegistroElettronico")
	wsHub := ws.NewHub(os.Getenv("REDIS_URL"))
	go wsHub.Run(ctx)

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
	rubricsRepo := rubrics.NewRepository(database)
	schoolCalendarRepo := schoolcalendar.NewRepository(database)

	// Repository FEQ — firma qualificata con valore legale (CAD art. 21 / eIDAS)
	feqRepo := signatures.NewFEQRepository(database)

	authMiddleware := auth.NewMiddleware(tokenManager, usersRepo)

	// 6. Setup Services
	mailerCfg := mailer.MailConfig{
		Host:     cfg.Mail.Host,
		Port:     cfg.Mail.Port,
		Username: cfg.Mail.Username,
		Password: cfg.Mail.Password,
		From:     cfg.Mail.From,
	}
	smsCfg := mailer.SMSConfig{
		Provider: "mock",
	}
	mailerSvc := mailer.NewService(mailerCfg, smsCfg)
	schoolCalendarSvc := schoolcalendar.NewService(schoolCalendarRepo)

	authSvc := auth.NewService(authRepo, tokenManager, mfaService, mailerSvc)
	usersSvc := users.NewService(usersRepo)
	classesSvc := classes.NewService(classesRepo)
	gradesSvc := grades.NewService(gradesRepo, usersRepo, database, wsHub)
	gradesAnalytics := grades.NewAnalyticsService(gradesRepo)
	attendanceSvc := attendance.NewService(attendanceRepo, usersRepo, wsHub, schoolCalendarSvc)
	docsSvc := documents.NewService(docsRepo)
	schedSvc := scheduling.NewService(schedRepo, teachersRepo, nil, nil, nil)
	pctoSvc := pcto.NewService(pctoRepo)
	orientSvc := orientamento.NewService(orientRepo)
	schoolsSvc := schools.NewService(schoolsRepo)

	// Firma FEA base (legacy)
	signaturesSvc := signatures.NewService(signatures.NewRepository(database), docsSvc, usersRepo)

	// Firma FEQ/FES con valore legale — D.Lgs. 82/2005 (CAD) + eIDAS Reg. UE 910/2014
	feqSvc := signatures.NewQualifiedService(database, feqRepo, usersRepo)

	// Export SIDI/MIUR — scrutini, presenze, certificazioni DM 742/2017
	sidiSvc := signatures.NewSidiExportService()

	commsSvc := communications.NewService(commsRepo, usersRepo)
	notesSvc := notes.NewService(notesRepo, usersRepo)
	adminSvc := admin.NewService(adminRepo)
	agendaSvc := agenda.NewService(agendaRepo)
	colloquiSvc := colloqui.NewService(colloquiRepo)
	verbaliSvc := verbali.NewService(verbaliRepo)
	extraSvc := extracurricular.NewService(extraRepo)
	notifSvc := notifications.NewService(notifRepo)
	tripsSvc := trips.NewService(tripsRepo)
	rubricsSvc := rubrics.NewService(rubricsRepo)

	wsTicketStore := wsticket.NewStore()

	// 7. Setup Handlers
	authH := auth.NewHandler(authSvc, wsTicketStore)
	usersH := users.NewHandler(usersSvc)
	schoolsH := schools.NewHandler(schoolsSvc)
	classesH := classes.NewHandler(classesSvc)
	gradesH := grades.NewHandler(gradesSvc, gradesAnalytics)
	attendanceH := attendance.NewHandler(attendanceSvc)
	docsUploader := upload.NewSupabaseUploader(cfg.Supabase.URL, cfg.Supabase.Key, cfg.Supabase.Bucket)
	docsH := documents.NewHandler(docsSvc, docsUploader)
	schedH := scheduling.NewHandler(schedSvc)

	// Handler firme: FEA base + FEQ/FES qualificata + SIDI export + CAD preservation
	signaturesH := signatures.NewHandler(signaturesSvc, feqSvc, sidiSvc)

	notesH := notes.NewHandler(notesSvc)
	adminH := admin.NewHandler(adminSvc)
	timetablesH := timetables.NewHandler(timetablesRepo)
	agendaH := agenda.NewHandler(agendaSvc)
	colloquiH := colloqui.NewHandler(colloquiSvc)
	verbaliH := verbali.NewHandler(verbaliSvc)
	extraH := extracurricular.NewHandler(extraSvc)
	notifH := notifications.NewHandler(notifSvc)
	tripsH := trips.NewHandler(tripsSvc)
	rubricsH := rubrics.NewHandler(rubricsSvc)
	schoolCalendarH := schoolcalendar.NewHandler(schoolCalendarSvc)

	elearningSvc := elearning.NewService(cfg.Elearning)
	elearningH := elearning.NewHandler(elearningSvc)

	wsHandler := ws.NewHandler(wsHub)

	adminMiddleware := admin.NewMiddleware()
	healthH := handler.NewHealthHandler(database)

	// 8. Setup Router
	r := gin.New()
	_ = r.SetTrustedProxies([]string{"127.0.0.1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
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

		api.GET("/swagger/doc.json", func(c *gin.Context) {
			c.File("../docs/openapi.yaml")
		})

		authH.RegisterRoutes(api, authMiddleware)
		api.GET("/public/schools", schoolsH.ListPublic)

		api.GET("/ws", authMiddleware.AuthenticateWSTicket(wsTicketStore), func(c *gin.Context) {
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
				usersGroup.GET("/students/:id/fascicolo", usersH.GetFascicolo)
			}

			classesH.RegisterRoutes(protected)
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
			notesH.RegisterRoutes(protected)
			rubricsH.RegisterRoutes(protected)
			schoolCalendarH.RegisterRoutes(protected)

			pdpRepo := pdp.NewRepository(database)
			pdpSvc := pdp.NewService(pdpRepo, usersRepo)
			pdpH := pdp.NewHandler(pdpSvc)
			pdpH.RegisterRoutes(protected)

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

			schoolsH.RegisterRoutes(protected)

			commsH := communications.NewHandler(commsSvc, docsUploader)
			commsH.RegisterRoutes(protected)

			textbooksRepo := textbooks.NewRepository(database)
			textbooksSvc := textbooks.NewService(textbooksRepo)
			textbooksH := textbooks.NewHandler(textbooksSvc)
			textbooksH.RegisterRoutes(protected)

			scrutinyRepo := scrutiny.NewRepository(database)
			scrutinySvc := scrutiny.NewService(scrutinyRepo, gradesRepo, classesRepo, usersRepo, attendanceRepo)
			scrutinyH := scrutiny.NewHandler(scrutinySvc)
			scrutinyH.RegisterRoutes(protected)

			subjectsRepo := subjects.NewRepository(database)
			subjectsSvc := subjects.NewService(subjectsRepo)
			subjectsH := subjects.NewHandler(subjectsSvc)
			subjectsH.RegisterRoutes(protected)

			goalsRepo := student_goals.NewRepository(database)
			goalsSvc := student_goals.NewService(goalsRepo)
			goalsH := student_goals.NewHandler(goalsSvc)
			goalsH.RegisterRoutes(protected)

			studentsFascicoloH := students.NewFascicoloHandler(database)
			studentsFascicoloH.RegisterRoutes(protected)

			subsRepo := substitutions.NewRepository(database)
			subsSvc := substitutions.NewService(subsRepo)
			subsH := substitutions.NewHandler(subsSvc)
			subsH.RegisterRoutes(protected)

			gmRepo := general_meetings.NewRepository(database)
			gmSvc := general_meetings.NewService(gmRepo)
			gmH := general_meetings.NewHandler(gmSvc)
			gmH.RegisterRoutes(protected)

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

			creditsRepo := credits.NewRepository(database)
			creditsSvc := credits.NewService(creditsRepo)
			creditsH := credits.NewHandler(creditsSvc)
			creditsH.RegisterRoutes(protected)

			recoveryRepo := recovery.NewRepository(database)
			recoverySvc := recovery.NewService(recoveryRepo)
			recoveryH := recovery.NewHandler(recoverySvc)
			recoveryH.RegisterRoutes(protected)

			supportRepo := support.NewRepository(database)
			supportSvc := support.NewService(supportRepo)
			supportH := support.NewHandler(supportSvc)
			supportH.RegisterRoutes(protected)

			protected.GET("/students/dashboard/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"total_grades":   0,
					"presence_rate":  100,
					"upcoming_tests": 0,
				})
			})
			protected.GET("/parents/dashboard/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"total_children":   1,
					"unread_messages":  0,
					"pending_payments": 0,
				})
			})

			tenantsRepo := tenants.NewRepository(database)
			tenantsSvc := tenants.NewService(tenantsRepo)
			tenantsH := tenants.NewHandler(tenantsSvc)
			tenantsH.RegisterRoutes(protected)

			certRepo := certificates.NewRepository(database)
			certSvc := certificates.NewService(certRepo, usersRepo)
			certH := certificates.NewHandler(certSvc)
			certH.RegisterRoutes(protected)

			auditRepo := auditlog.NewRepository(database)
			auditSvc := auditlog.NewService(auditRepo)
			auditH := auditlog.NewHandler(auditSvc)
			auditH.RegisterRoutes(protected)

			udaRepo := uda.NewRepository(database)
			udaSvc := uda.NewService(udaRepo)
			udaH := uda.NewHandler(udaSvc)
			udaH.RegisterRoutes(protected)

			compRepo := competencies.NewRepository(database)
			compSvc := competencies.NewService(compRepo)
			compH := competencies.NewHandler(compSvc)
			compH.RegisterRoutes(protected)

			// Attività libere docente (ore a disposizione, riunioni, gita, formazione, etc.)
			teacherActRepo := teacher_activities.NewRepository(database)
			teacherActSvc := teacher_activities.NewService(teacherActRepo)
			teacherActH := teacher_activities.NewHandler(teacherActSvc)
			teacherActH.RegisterRoutes(protected)

			elearningH.RegisterRoutes(protected)

			// Gestione Pagamenti & PagoPA
			paymentsRepo := payments.NewRepository(database)
			paymentsSvc := payments.NewService(paymentsRepo, usersRepo)
			paymentsH := payments.NewHandler(paymentsSvc)
			paymentsH.RegisterRoutes(protected)

			// Firme qualificate FEQ/FES + SIDI export + CAD preservation
			signaturesH.RegisterRoutes(protected)

			adminH.RegisterRoutes(protected, adminMiddleware)
		}
	}

	// 9. Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

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
