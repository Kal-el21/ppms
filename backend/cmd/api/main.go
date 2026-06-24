package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
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
		requestRepository, revisionRepository, approvalRepository, provisioningSvc, eventBus,
	)

	// Task (must be created before Milestone so it can satisfy TaskProgressProvider)
	taskSvc := taskservice.NewTaskService(
		taskRepository, assigneeRepository, commentRepository, eventBus,
	)

	// Milestone (receives taskSvc as TaskProgressProvider)
	milestoneSvc := milestoneservice.NewMilestoneService(milestoneRepository, taskSvc)

	// Project (receives milestoneSvc as MilestoneProgressProvider)
	projectSvc := projectservice.NewProjectService(projectRepository, milestoneSvc)
	memberSvc := projectservice.NewMemberService(memberRepository)

	// Budget
	budgetSvc := budgetservice.NewBudgetService(budgetRepository, transactionRepository)
	transactionSvc := budgetservice.NewTransactionService(transactionRepository, budgetRepository, eventBus)

	// Attachment (Phase 7: now includes ownership resolver & request checker)
	attachmentSvc := attachmentservice.NewAttachmentService(
		attachmentRepository, minioClient, ownershipResolver, requestOwnershipAdapter,
	)

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
	events.RegisterNotificationSubscriber(eventBus, notificationSvc)

	// Notify all ADMINs when a new project request is submitted.
	eventBus.Subscribe("project.request.submitted", func(e events.Event) {
		data := e.Data.(map[string]interface{})
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		admins, err := userRepository.FindBySystemRole(userentity.RoleAdmin)
		if err != nil {
			logger.Log.Error().Err(err).Msg("event project.request.submitted: failed to fetch admins")
			return
		}
		for _, admin := range admins {
			_ = notificationSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     admin.ID,
				Type:       "REQUEST_SUBMITTED",
				Title:      "New Project Request Submitted",
				Message:    "A new request \"" + title + "\" requires your review.",
				EntityType: "PROJECT_REQUEST",
				EntityID:   &requestID,
				ActionURL:  "/project-requests/" + strconv.FormatUint(requestID, 10),
			})
		}
	})

	// Notify all PROJECT_MANAGERs of a project when a task in that project is completed.
	eventBus.Subscribe("task.completed", func(e events.Event) {
		data := e.Data.(map[string]interface{})
		taskID := data["task_id"].(uint64)
		title := data["title"].(string)

		projectIDRaw, ok := data["project_id"]
		if !ok {
			return
		}
		projectID := projectIDRaw.(uint64)

		members, err := memberRepository.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("event task.completed: failed to fetch project members")
			return
		}
		for _, m := range members {
			if string(m.ProjectRole) != "PROJECT_MANAGER" {
				continue
			}
			_ = notificationSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "TASK_COMPLETED",
				Title:      "Task Completed",
				Message:    "Task \"" + title + "\" has been marked as done.",
				EntityType: "TASK",
				EntityID:   &taskID,
				ActionURL:  "/projects/" + strconv.FormatUint(projectID, 10),
			})
		}
	})

	// Notify PROJECT_MANAGERs when budget usage crosses the 80% warning threshold.
	eventBus.Subscribe("budget.warning", func(e events.Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)

		members, err := memberRepository.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("event budget.warning: failed to fetch project members")
			return
		}
		for _, m := range members {
			if string(m.ProjectRole) != "PROJECT_MANAGER" {
				continue
			}
			_ = notificationSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "BUDGET_WARNING",
				Title:      "Budget Warning",
				Message:    "Project budget usage has crossed the 80% threshold.",
				EntityType: "PROJECT",
				EntityID:   &projectID,
				ActionURL:  "/projects/" + strconv.FormatUint(projectID, 10),
			})
		}
	})

	// Notify PROJECT_MANAGERs when budget usage exceeds 100%.
	eventBus.Subscribe("budget.over_limit", func(e events.Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)

		members, err := memberRepository.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("event budget.over_limit: failed to fetch project members")
			return
		}
		for _, m := range members {
			if string(m.ProjectRole) != "PROJECT_MANAGER" {
				continue
			}
			_ = notificationSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "BUDGET_OVER_LIMIT",
				Title:      "Budget Over Limit",
				Message:    "Project budget usage has exceeded 100% of allocated budget.",
				EntityType: "PROJECT",
				EntityID:   &projectID,
				ActionURL:  "/projects/" + strconv.FormatUint(projectID, 10),
			})
		}
	})

	// Notify the original sender when their handover is marked as received.
	eventBus.Subscribe("handover.received", func(e events.Event) {
		data := e.Data.(map[string]interface{})
		handoverID := data["handover_id"].(uint64)

		handover, err := handoverRepository.FindByID(handoverID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("event handover.received: failed to fetch handover")
			return
		}
		_ = notificationSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     handover.SenderID,
			Type:       "HANDOVER_RECEIVED",
			Title:      "Handover Received",
			Message:    "Your sent handover has been marked as received.",
			EntityType: "HANDOVER",
			EntityID:   &handoverID,
			ActionURL:  "/handovers/" + strconv.FormatUint(handoverID, 10),
		})
	})

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
