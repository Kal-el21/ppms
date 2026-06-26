package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kal-el21/backend/configs"
	"github.com/Kal-el21/backend/internal/database"
	"github.com/Kal-el21/backend/internal/events"
	"github.com/Kal-el21/backend/internal/infrastructure/membership"
	minioclient "github.com/Kal-el21/backend/internal/infrastructure/minio"
	"github.com/Kal-el21/backend/internal/middleware"
	"github.com/Kal-el21/backend/internal/shared/logger"
	"github.com/gin-gonic/gin"

	// ── Auth ─────────────────────────────────────────────────────────────────
	authhandler "github.com/Kal-el21/backend/internal/domain/auth/handler"
	authrepo "github.com/Kal-el21/backend/internal/domain/auth/repository"
	authservice "github.com/Kal-el21/backend/internal/domain/auth/service"
	emailpkg "github.com/Kal-el21/backend/internal/infrastructure/email"

	// ── Division ─────────────────────────────────────────────────────────────
	divisionhandler "github.com/Kal-el21/backend/internal/domain/division/handler"
	divisionrepo "github.com/Kal-el21/backend/internal/domain/division/repository"
	divisionservice "github.com/Kal-el21/backend/internal/domain/division/service"

	// ── User ─────────────────────────────────────────────────────────────────
	userhandler "github.com/Kal-el21/backend/internal/domain/user/handler"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	userservice "github.com/Kal-el21/backend/internal/domain/user/service"

	// ── Project Request ───────────────────────────────────────────────────────
	requesthandler "github.com/Kal-el21/backend/internal/domain/project_request/handler"
	requestrepo "github.com/Kal-el21/backend/internal/domain/project_request/repository"
	requestservice "github.com/Kal-el21/backend/internal/domain/project_request/service"

	// ── Project ───────────────────────────────────────────────────────────────
	projecthandler "github.com/Kal-el21/backend/internal/domain/project/handler"
	projectrepo "github.com/Kal-el21/backend/internal/domain/project/repository"
	projectservice "github.com/Kal-el21/backend/internal/domain/project/service"

	// ── Milestone ─────────────────────────────────────────────────────────────
	milestonehandler "github.com/Kal-el21/backend/internal/domain/milestone/handler"
	milestonerepo "github.com/Kal-el21/backend/internal/domain/milestone/repository"
	milestoneservice "github.com/Kal-el21/backend/internal/domain/milestone/service"

	// ── Task ──────────────────────────────────────────────────────────────────
	taskhandler "github.com/Kal-el21/backend/internal/domain/task/handler"
	taskrepo "github.com/Kal-el21/backend/internal/domain/task/repository"
	taskservice "github.com/Kal-el21/backend/internal/domain/task/service"

	// ── Budget ────────────────────────────────────────────────────────────────
	budgethandler "github.com/Kal-el21/backend/internal/domain/budget/handler"
	budgetrepo "github.com/Kal-el21/backend/internal/domain/budget/repository"
	budgetservice "github.com/Kal-el21/backend/internal/domain/budget/service"

	// ── Attachment ────────────────────────────────────────────────────────────
	approvalhandler "github.com/Kal-el21/backend/internal/domain/approval/handler"
	approvalrepo "github.com/Kal-el21/backend/internal/domain/approval/repository"
	approvalservice "github.com/Kal-el21/backend/internal/domain/approval/service"
	attachmenthandler "github.com/Kal-el21/backend/internal/domain/attachment/handler"
	attachmentrepo "github.com/Kal-el21/backend/internal/domain/attachment/repository"
	attachmentservice "github.com/Kal-el21/backend/internal/domain/attachment/service"

	// ── Handover ──────────────────────────────────────────────────────────────
	handoverhandler "github.com/Kal-el21/backend/internal/domain/handover/handler"
	handoverrepo "github.com/Kal-el21/backend/internal/domain/handover/repository"
	handoverservice "github.com/Kal-el21/backend/internal/domain/handover/service"

	// ── Notification ──────────────────────────────────────────────────────────
	notificationhandler "github.com/Kal-el21/backend/internal/domain/notification/handler"
	notificationrepo "github.com/Kal-el21/backend/internal/domain/notification/repository"
	notificationservice "github.com/Kal-el21/backend/internal/domain/notification/service"

	// ── Dashboard ─────────────────────────────────────────────────────────────
	dashboardhandler "github.com/Kal-el21/backend/internal/domain/dashboard/handler"
	dashboardrepo "github.com/Kal-el21/backend/internal/domain/dashboard/repository"
	dashboardservice "github.com/Kal-el21/backend/internal/domain/dashboard/service"

	// ── Audit ─────────────────────────────────────────────────────────────────
	audithandler "github.com/Kal-el21/backend/internal/domain/audit/handler"
	auditrepo "github.com/Kal-el21/backend/internal/domain/audit/repository"
	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"

	// ── Reporting ─────────────────────────────────────────────────────────────
	reportinghandler "github.com/Kal-el21/backend/internal/domain/reporting/handler"
	reportingrepo "github.com/Kal-el21/backend/internal/domain/reporting/repository"
	reportingservice "github.com/Kal-el21/backend/internal/domain/reporting/service"

	// ── Search ────────────────────────────────────────────────────────────────
	searchhandler "github.com/Kal-el21/backend/internal/domain/search/handler"
	searchrepo "github.com/Kal-el21/backend/internal/domain/search/repository"
	searchservice "github.com/Kal-el21/backend/internal/domain/search/service"
)

func main() {
	// =========================================================================
	// 1. CONFIG & LOGGER
	// =========================================================================
	cfg := configs.Load()
	logger.Init(cfg.AppEnv)

	// =========================================================================
	// 2. DATABASE
	// =========================================================================
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// =========================================================================
	// 3. EVENT BUS
	// =========================================================================
	eventBus := events.NewBus()

	// =========================================================================
	// 4. MINIO CLIENT
	// =========================================================================
	minioClient, err := minioclient.NewClient(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to initialize minio client")
	}

	// =========================================================================
	// 5. AUDIT SERVICE
	// Dibuat paling awal karena menjadi dependency banyak handler.
	// Hanya ada SATU deklarasi auditRepository & auditSvc di seluruh file ini.
	// =========================================================================
	auditRepository := auditrepo.NewAuditRepository(db)
	auditSvc := auditservice.NewAuditService(auditRepository)

	// =========================================================================
	// 6. REPOSITORIES
	// =========================================================================

	// Auth
	sessionRepository := authrepo.NewSessionRepository(db)
	otpRepository := authrepo.NewOTPRepository(db)
	otpSessionStore := authservice.NewOTPSessionStore()
	emailSvc := emailpkg.NewEmailService(cfg)

	// User & Division
	userRepository := userrepo.NewUserRepository(db)
	divisionRepository := divisionrepo.NewDivisionRepository(db)

	// Project Request
	requestRepository := requestrepo.NewRequestRepository(db)
	revisionRepository := requestrepo.NewRevisionRepository(db)
	approvalRepository := requestrepo.NewApprovalRepository(db)

	// Project & Member
	projectRepository := projectrepo.NewProjectRepository(db)
	memberRepository := projectrepo.NewMemberRepository(db)

	// Milestone
	milestoneRepository := milestonerepo.NewMilestoneRepository(db)

	// Task
	taskRepository := taskrepo.NewTaskRepository(db)
	assigneeRepository := taskrepo.NewAssigneeRepository(db)
	commentRepository := taskrepo.NewCommentRepository(db)

	// Budget
	budgetRepository := budgetrepo.NewBudgetRepository(db)
	transactionRepository := budgetrepo.NewTransactionRepository(db)

	// Attachment
	attachmentRepository := attachmentrepo.NewAttachmentRepository(db)
	approvalWorkflowRepository := approvalrepo.NewApprovalWorkflowRepository(db)
	approvalLevelRepository := approvalrepo.NewApprovalLevelRepository(db)

	// Handover
	handoverRepository := handoverrepo.NewHandoverRepository(db)

	// Notification
	notificationRepository := notificationrepo.NewNotificationRepository(db)
	preferenceRepository := notificationrepo.NewPreferenceRepository(db)

	// Dashboard / Reporting / Search
	dashboardRepository := dashboardrepo.NewDashboardRepository(db)
	reportingRepository := reportingrepo.NewReportingRepository(db)
	searchRepository := searchrepo.NewSearchRepository(db)

	// =========================================================================
	// 7. INFRASTRUCTURE ADAPTERS (Phase 7 additions)
	// =========================================================================

	// MembershipChecker: resolves project_id from any entity_type, then checks
	// whether a user is an active member of that project (closes attachment
	// ownership security gap flagged in Phase 4, implemented in Phase 7).
	membershipChecker := membership.NewMembershipChecker(db)

	// RequestOwnershipAdapter: validates that a user is the requester/owner of
	// a PROJECT_REQUEST, used exclusively by AttachmentService for PROJECT_REQUEST
	// entity_type validation without introducing a circular import.
	requestOwnershipAdapter := requestservice.NewRequestOwnershipAdapter(requestRepository)

	// OwnershipResolver: wraps membershipChecker for the AttachmentService API.
	ownershipResolver := attachmentservice.NewOwnershipResolver(membershipChecker)

	// =========================================================================
	// 8. SERVICES
	// Order matters here: taskSvc → milestoneSvc → projectSvc because of
	// the small provider interfaces (TaskProgressProvider, MilestoneProgressProvider)
	// used to avoid circular imports between domain packages.
	// =========================================================================

	// Auth / User / Division
	authSvc := authservice.NewAuthService(userRepository, sessionRepository, otpRepository, otpSessionStore, emailSvc, cfg)
	userSvc := userservice.NewUserService(userRepository)
	divisionSvc := divisionservice.NewDivisionService(divisionRepository)

	// Project Request
	provisioningSvc := projectservice.NewProvisioningService(projectRepository, memberRepository)
	requestSvc := requestservice.NewRequestService(
		requestRepository, revisionRepository, approvalRepository, userRepository, provisioningSvc, eventBus,
	)

	// Task (must be created before Milestone so it can satisfy TaskProgressProvider)
	taskSvc := taskservice.NewTaskService(
		taskRepository, assigneeRepository, commentRepository, eventBus,
	)

	// Milestone (receives taskSvc as TaskProgressProvider)
	milestoneSvc := milestoneservice.NewMilestoneService(milestoneRepository, taskSvc, eventBus)

	// Project (receives milestoneSvc as MilestoneProgressProvider)
	projectSvc := projectservice.NewProjectService(projectRepository, milestoneSvc, eventBus)
	memberSvc := projectservice.NewMemberService(memberRepository)

	// Budget
	budgetSvc := budgetservice.NewBudgetService(budgetRepository, transactionRepository)
	transactionSvc := budgetservice.NewTransactionService(transactionRepository, budgetRepository, eventBus)

	// Attachment (Phase 7: now includes ownership resolver & request checker)
	attachmentSvc := attachmentservice.NewAttachmentService(
		attachmentRepository, minioClient, ownershipResolver, requestOwnershipAdapter,
	)
	approvalWorkflowSvc := approvalservice.NewApprovalWorkflowService(approvalWorkflowRepository)
	approvalLevelSvc := approvalservice.NewApprovalLevelService(approvalLevelRepository)

	// Handover
	handoverSvc := handoverservice.NewHandoverService(handoverRepository, eventBus)

	// Notification
	notificationSvc := notificationservice.NewNotificationService(
		notificationRepository,
		preferenceRepository,
		emailSvc,       // EmailSender
		userRepository, // UserEmailProvider
	)
	preferenceSvc := notificationservice.NewPreferenceService(preferenceRepository)

	// Dashboard / Reporting / Search
	dashboardSvc := dashboardservice.NewDashboardService(dashboardRepository)
	reportingSvc := reportingservice.NewReportingService(reportingRepository)
	searchSvc := searchservice.NewSearchService(searchRepository)

	// =========================================================================
	// 9. HANDLERS
	// Write-action handlers receive auditSvc for audit trail.
	// Read-only handlers (dashboard, audit viewer, reporting, search,
	// notification) do NOT receive auditSvc.
	// =========================================================================

	authHdl := authhandler.NewAuthHandler(authSvc, auditSvc, cfg)
	userHdl := userhandler.NewUserHandler(userSvc, cfg, auditSvc, minioClient)
	divisionHdl := divisionhandler.NewDivisionHandler(divisionSvc, auditSvc)
	requestHdl := requesthandler.NewRequestHandler(requestSvc, auditSvc)
	projectHdl := projecthandler.NewProjectHandler(projectSvc, auditSvc)
	memberHdl := projecthandler.NewMemberHandler(memberSvc, auditSvc)
	milestoneHdl := milestonehandler.NewMilestoneHandler(milestoneSvc, auditSvc)
	taskHdl := taskhandler.NewTaskHandler(taskSvc, auditSvc)
	budgetHdl := budgethandler.NewBudgetHandler(budgetSvc, auditSvc)
	transactionHdl := budgethandler.NewTransactionHandler(transactionSvc, auditSvc)
	attachmentHdl := attachmenthandler.NewAttachmentHandler(attachmentSvc, auditSvc)
	approvalHdl := approvalhandler.NewApprovalHandler(approvalWorkflowSvc, approvalLevelSvc)
	handoverHdl := handoverhandler.NewHandoverHandler(handoverSvc, auditSvc)

	notificationHdl := notificationhandler.NewNotificationHandler(notificationSvc, preferenceSvc)
	dashboardHdl := dashboardhandler.NewDashboardHandler(dashboardSvc)
	auditHdl := audithandler.NewAuditHandler(auditSvc)
	reportingHdl := reportinghandler.NewReportingHandler(reportingSvc)
	searchHdl := searchhandler.NewSearchHandler(searchSvc)

	// =========================================================================
	// 10. EVENT SUBSCRIBERS
	// RegisterNotificationSubscriber wires up the base subscribers (request
	// approved/rejected, task assigned, handover sent). Cross-domain subscribers
	// that need access to multiple repositories are registered inline below so
	// they can close over the repository variables without creating circular
	// imports between domain packages and the events package.
	// =========================================================================
	events.RegisterNotificationSubscriber(eventBus, notificationSvc, userRepository, memberRepository, handoverRepository)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		dueSoonNotified := make(map[uint64]time.Time)
		for range ticker.C {
			now := time.Now()

			overdueTasks, err := taskRepository.FindOverdue()
			if err != nil {
				logger.Log.Error().Err(err).Msg("scheduler: failed to find overdue tasks")
			} else {
				for _, task := range overdueTasks {
					assignees, err := assigneeRepository.FindActiveByTaskID(task.ID)
					if err != nil {
						logger.Log.Error().Err(err).Uint64("task_id", task.ID).Msg("scheduler: failed to find assignees for overdue task")
						continue
					}
					for _, assignee := range assignees {
						eventBus.Publish(events.Event{
							Name: "task.overdue",
							Data: map[string]interface{}{
								"task_id":     task.ID,
								"title":       task.Title,
								"project_id":  task.ProjectID,
								"assignee_id": assignee.UserID,
							},
						})
					}
				}
			}

			dueSoonTasks, err := taskRepository.FindDueSoon()
			if err != nil {
				logger.Log.Error().Err(err).Msg("scheduler: failed to find due soon tasks")
			} else {
				for _, task := range dueSoonTasks {
					if lastNotified, ok := dueSoonNotified[task.ID]; ok && now.Sub(lastNotified) < 24*time.Hour {
						continue
					}

					assignees, err := assigneeRepository.FindActiveByTaskID(task.ID)
					if err != nil {
						logger.Log.Error().Err(err).Uint64("task_id", task.ID).Msg("scheduler: failed to find assignees for due soon task")
						continue
					}
					for _, assignee := range assignees {
						eventBus.Publish(events.Event{
							Name: "task.due_soon",
							Data: map[string]interface{}{
								"task_id":     task.ID,
								"title":       task.Title,
								"project_id":  task.ProjectID,
								"assignee_id": assignee.UserID,
								"due_date":    task.DueDate.Format("2006-01-02"),
							},
						})
					}
					dueSoonNotified[task.ID] = now
				}
			}
		}
	}()

	// =========================================================================
	// 11. RATE LIMITERS (Phase 7)
	// =========================================================================
	loginRateLimiter := middleware.NewRateLimiter(5, 1*time.Minute)
	refreshRateLimiter := middleware.NewRateLimiter(10, 1*time.Minute)
	transactionRateLimiter := middleware.NewRateLimiter(20, 1*time.Minute)

	// =========================================================================
	// 12. GIN ROUTER
	// =========================================================================
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORSConfig([]string{cfg.FrontendURL}))

	// ── Health / Readiness ──────────────────────────────────────────────────
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// =========================================================================
	// 13. ROUTES
	// =========================================================================
	v1 := router.Group("/api/v1")

	// ── Public: Auth ─────────────────────────────────────────────────────────
	auth := v1.Group("/auth")
	{
		auth.POST("/login", loginRateLimiter.LimitByIP(), authHdl.Login)
		auth.POST("/verify-otp", authHdl.VerifyOTP)
		auth.POST("/resend-otp", loginRateLimiter.LimitByIP(), authHdl.ResendOTP)
		auth.POST("/refresh", refreshRateLimiter.LimitByIP(), authHdl.RefreshToken)
		auth.POST("/logout", authHdl.Logout)
	}

	// ── Protected (require valid JWT) ─────────────────────────────────────────
	authMiddleware := middleware.AuthMiddleware(cfg, userRepository)

	protected := v1.Group("")
	protected.Use(authMiddleware)
	protected.Use(middleware.CSRFProtection())
	{
		// Auth (self-service, no system-role restriction)
		protected.POST("/auth/change-password", authHdl.ChangePassword)
		protected.POST("/auth/revoke-sessions", authHdl.RevokeAllSessions)

		//Profile & Settings (self-service, no system-role restriction)
		protected.GET("/me", userHdl.GetMe)
		protected.PUT("/me", userHdl.UpdateProfile)
		protected.POST("/me/photo", userHdl.UploadProfilePhoto)
		protected.POST("/me/toggle-2fa", userHdl.Toggle2FA)
		protected.POST("/me/toggle-email-notification", userHdl.ToggleEmailNotification)

		// ── Division ─────────────────────────────────────────────────────────
		divisions := protected.Group("/divisions")
		{
			divisions.GET("", divisionHdl.GetAll)
			divisions.GET("/:id", divisionHdl.GetByID)
			divisions.POST("", middleware.RequireSystemRole("ADMIN"), divisionHdl.Create)
			divisions.PUT("/:id", middleware.RequireSystemRole("ADMIN"), divisionHdl.Update)
			divisions.DELETE("/:id", middleware.RequireSystemRole("ADMIN"), divisionHdl.Delete)
		}

		// ── User Management (ADMIN only) ─────────────────────────────────────
		users := protected.Group("/users")
		users.Use(middleware.RequireSystemRole("ADMIN"))
		{
			users.GET("", userHdl.GetAll)
			users.GET("/:id", userHdl.GetByID)
			users.POST("", userHdl.Create)
			users.PUT("/:id", userHdl.Update)
			users.PATCH("/:id/role", userHdl.AssignRole)
			users.DELETE("/:id", userHdl.Deactivate)
			users.POST("/:id/restore", userHdl.Restore)
		}

		// ── Project Requests (VIEWER excluded) ───────────────────────────────
		requests := protected.Group("/project-requests")
		requests.Use(middleware.RequireSystemRole("USER"))
		{
			requests.POST("", requestHdl.CreateDraft)
			requests.GET("", requestHdl.GetList)
			requests.GET("/:id", requestHdl.GetByID)
			requests.PUT("/:id", requestHdl.UpdateDraft)
			requests.DELETE("/:id", requestHdl.DeleteDraft)
			requests.POST("/:id/submit", requestHdl.Submit)
			requests.POST("/:id/revise", requestHdl.Revise)
			requests.GET("/:id/revisions", requestHdl.GetRevisionHistory)
			requests.GET("/:id/approvals", requestHdl.GetApprovalHistory)
			requests.POST("/:id/review",
				middleware.RequireSystemRole("ADMIN"),
				requestHdl.Review,
			)
		}

		// ── Projects: list (system-role scoped) ──────────────────────────────
		projects := protected.Group("/projects")
		projects.Use(middleware.RequireSystemRole("USER"))
		{
			projects.GET("", projectHdl.GetList)
		}

		// ── Project Detail (requires project membership / ADMIN override) ────
		// All routes under /projects/:id share the ProjectContextMiddleware which
		// validates membership and injects project_role into the Gin context.
		projectDetail := protected.Group("/projects/:id")
		projectDetail.Use(middleware.ProjectContextMiddleware(memberRepository))
		{
			// Project
			projectDetail.GET("", projectHdl.GetByID)
			projectDetail.PUT("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				projectHdl.Update,
			)
			projectDetail.PATCH("/status",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				projectHdl.ChangeStatus,
			)

			// Members
			projectDetail.GET("/members", memberHdl.GetList)
			projectDetail.POST("/members",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				memberHdl.Add,
			)
			projectDetail.PATCH("/members/:memberId/role",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				memberHdl.ChangeRole,
			)
			projectDetail.DELETE("/members/:memberId",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				memberHdl.Remove,
			)

			// Milestones
			projectDetail.GET("/milestones", milestoneHdl.GetList)
			projectDetail.POST("/milestones",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				milestoneHdl.Create,
			)
			projectDetail.PATCH("/milestones/reorder",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				milestoneHdl.Reorder,
			)

			// Tasks
			projectDetail.GET("/tasks", taskHdl.GetList)
			projectDetail.POST("/tasks",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				taskHdl.Create,
			)

			// Budget
			projectDetail.GET("/budget", budgetHdl.GetByProjectID)
			projectDetail.POST("/budget",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				budgetHdl.Create,
			)

			// Handovers
			projectDetail.GET("/handovers", handoverHdl.GetList)
			projectDetail.POST("/handovers",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				handoverHdl.Create,
			)

			// Reporting: project-scoped (PM or ADMIN — ADMIN_OVERRIDE from context middleware)
			projectDetail.POST("/reports/generate",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				reportingHdl.GenerateForProject,
			)
		}

		// ── Milestone Detail ──────────────────────────────────────────────────
		milestoneDetail := protected.Group("/projects/:id/milestones/:milestoneId")
		milestoneDetail.Use(middleware.ProjectContextMiddleware(memberRepository))
		{
			milestoneDetail.PUT("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				milestoneHdl.Update,
			)
			milestoneDetail.PATCH("/status",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				milestoneHdl.ChangeStatus,
			)
			milestoneDetail.DELETE("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				milestoneHdl.Delete,
			)
		}

		// ── Task Detail ───────────────────────────────────────────────────────
		// MEMBER is allowed on status/progress/comments; the service layer applies
		// the "assigned only" restriction for MEMBERs automatically.
		taskDetail := protected.Group("/projects/:id/tasks/:taskId")
		taskDetail.Use(middleware.ProjectContextMiddleware(memberRepository))
		{
			taskDetail.GET("", taskHdl.GetByID)
			taskDetail.PUT("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				taskHdl.Update,
			)
			taskDetail.PATCH("/status",
				middleware.RequireProjectRole("PROJECT_MANAGER", "MEMBER"),
				taskHdl.ChangeStatus,
			)
			taskDetail.PATCH("/progress",
				middleware.RequireProjectRole("PROJECT_MANAGER", "MEMBER"),
				taskHdl.UpdateProgress,
			)
			taskDetail.POST("/assignees",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				taskHdl.AssignUsers,
			)
			taskDetail.DELETE("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				taskHdl.Delete,
			)
			taskDetail.GET("/comments", taskHdl.GetComments)
			taskDetail.POST("/comments",
				middleware.RequireProjectRole("PROJECT_MANAGER", "MEMBER"),
				taskHdl.AddComment,
			)
		}

		// ── Budget Detail (Transactions) ──────────────────────────────────────
		budgetDetail := protected.Group("/projects/:id/budget/:budgetId")
		budgetDetail.Use(middleware.ProjectContextMiddleware(memberRepository))
		{
			budgetDetail.PUT("",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				budgetHdl.Update,
			)
			budgetDetail.GET("/transactions", transactionHdl.GetList)
			budgetDetail.POST("/transactions",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				transactionRateLimiter.LimitByUser(),
				transactionHdl.Create,
			)
		}

		// ── Handover Detail ───────────────────────────────────────────────────
		handoverDetail := protected.Group("/projects/:id/handovers/:handoverId")
		handoverDetail.Use(middleware.ProjectContextMiddleware(memberRepository))
		{
			handoverDetail.PATCH("/received", handoverHdl.MarkReceived)
			handoverDetail.PATCH("/returned",
				middleware.RequireProjectRole("PROJECT_MANAGER"),
				handoverHdl.MarkReturned,
			)
		}

		// ── Attachments (generic, ownership validated in service layer) ───────
		attachments := protected.Group("/attachments")
		{
			attachments.POST("", attachmentHdl.Upload)
			attachments.GET("", attachmentHdl.GetByEntity)
			attachments.GET("/:id/download", attachmentHdl.GetDownloadURL)
			attachments.GET("/:id/versions", attachmentHdl.GetVersionHistory)
			attachments.DELETE("/:id", attachmentHdl.Delete)
		}

		// ── Notifications ─────────────────────────────────────────────────────
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHdl.GetList)
			notifications.PATCH("/:id/read", notificationHdl.MarkAsRead)
			notifications.PATCH("/read-all", notificationHdl.MarkAllAsRead)
			notifications.GET("/preferences", notificationHdl.GetPreferences)
			notifications.PUT("/preferences", notificationHdl.UpdatePreference)
		}

		// ── Dashboard ─────────────────────────────────────────────────────────
		protected.GET("/dashboard", dashboardHdl.GetSummary)

		// ── Audit Logs (ADMIN only) ───────────────────────────────────────────
		protected.GET("/audit-logs",
			middleware.RequireSystemRole("ADMIN"),
			auditHdl.GetList,
		)

		// ── Reporting: system-wide (ADMIN only; project-scoped is above) ─────
		approvalWorkflows := protected.Group("/approval-workflows")
		approvalWorkflows.Use(middleware.RequireSystemRole("ADMIN"))
		{
			approvalWorkflows.GET("", approvalHdl.ListWorkflows)
			approvalWorkflows.POST("", approvalHdl.CreateWorkflow)
			approvalWorkflows.GET("/:id", approvalHdl.GetWorkflow)
			approvalWorkflows.GET("/:id/levels", approvalHdl.GetLevels)
			approvalWorkflows.POST("/:id/levels", approvalHdl.CreateLevel)
		}

		protected.POST("/reports/generate",
			middleware.RequireSystemRole("ADMIN"),
			reportingHdl.Generate,
		)

		// ── Global Search ─────────────────────────────────────────────────────
		protected.GET("/search", searchHdl.Search)
	}

	// =========================================================================
	// 14. START SERVER
	// =========================================================================
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info().Msgf("PPMS backend starting on %s (env=%s)", addr, cfg.AppEnv)

	if err := router.Run(addr); err != nil {
		logger.Log.Fatal().Err(err).Msg("server failed to start")
	}
}
