package events

import (
	"strconv"

	handoverrepo "github.com/Kal-el21/backend/internal/domain/handover/repository"
	notificationservice "github.com/Kal-el21/backend/internal/domain/notification/service"
	projectentity "github.com/Kal-el21/backend/internal/domain/project/entity"
	memberrepo "github.com/Kal-el21/backend/internal/domain/project/repository"
	userentity "github.com/Kal-el21/backend/internal/domain/user/entity"
	userrepo "github.com/Kal-el21/backend/internal/domain/user/repository"
	"github.com/Kal-el21/backend/internal/shared/logger"
)

func RegisterNotificationSubscriber(
	bus *Bus,
	notifSvc notificationservice.NotificationService,
	userRepo userrepo.UserRepository,
	memberRepo memberrepo.MemberRepository,
	handoverRepo handoverrepo.HandoverRepository,
) {
	bus.Subscribe("project.request.submitted", func(e Event) {
		data := e.Data.(map[string]interface{})
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		admins, err := userRepo.FindBySystemRole(userentity.RoleAdmin)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification project.request.submitted: failed to fetch admins")
			return
		}
		for _, admin := range admins {
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     admin.ID,
				Type:       "REQUEST_SUBMITTED",
				Title:      "New Project Request Submitted",
				Message:    "A new request \"" + title + "\" requires your review.",
				EntityType: "PROJECT_REQUEST",
				EntityID:   &requestID,
				ActionURL:  "/project-requests/" + uint64ToStr(requestID),
			})
		}
	})

	bus.Subscribe("project.request.approved", func(e Event) {
		data := e.Data.(map[string]interface{})
		requesterID := data["requester_id"].(uint64)
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		err := notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     requesterID,
			Type:       "REQUEST_APPROVED",
			Title:      "Project Request Approved",
			Message:    "Your request \"" + title + "\" has been approved and a project has been created.",
			EntityType: "PROJECT_REQUEST",
			EntityID:   &requestID,
			ActionURL:  "/project-requests/" + uint64ToStr(requestID),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to create notification for request approval")
		}
	})

	bus.Subscribe("project.request.rejected", func(e Event) {
		data := e.Data.(map[string]interface{})
		requesterID := data["requester_id"].(uint64)
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		err := notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     requesterID,
			Type:       "REQUEST_REJECTED",
			Title:      "Project Request Rejected",
			Message:    "Your request \"" + title + "\" has been rejected.",
			EntityType: "PROJECT_REQUEST",
			EntityID:   &requestID,
			ActionURL:  "/project-requests/" + uint64ToStr(requestID),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to create notification for request rejection")
		}
	})

	bus.Subscribe("project.request.revision_requested", func(e Event) {
		data := e.Data.(map[string]interface{})
		requesterID := data["requester_id"].(uint64)
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		err := notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     requesterID,
			Type:       "REVISION_REQUESTED",
			Title:      "Project Request Needs Revision",
			Message:    "Your request \"" + title + "\" requires revision.",
			EntityType: "PROJECT_REQUEST",
			EntityID:   &requestID,
			ActionURL:  "/project-requests/" + uint64ToStr(requestID),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to create notification for revision request")
		}
	})

	bus.Subscribe("project.request.revised", func(e Event) {
		data := e.Data.(map[string]interface{})
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		admins, err := userRepo.FindBySystemRole(userentity.RoleAdmin)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification project.request.revised: failed to fetch admins")
			return
		}
		for _, admin := range admins {
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     admin.ID,
				Type:       "REQUEST_REVISED",
				Title:      "Project Request Revised",
				Message:    "The request \"" + title + "\" has been revised and is ready for review.",
				EntityType: "PROJECT_REQUEST",
				EntityID:   &requestID,
				ActionURL:  "/project-requests/" + uint64ToStr(requestID),
			})
		}
	})

	bus.Subscribe("task.assigned", func(e Event) {
		data := e.Data.(map[string]interface{})
		userID := data["user_id"].(uint64)
		taskID := data["task_id"].(uint64)

		err := notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     userID,
			Type:       "TASK_ASSIGNED",
			Title:      "New Task Assigned",
			Message:    "You have been assigned to a new task.",
			EntityType: "TASK",
			EntityID:   &taskID,
			ActionURL:  "/tasks/" + uint64ToStr(taskID),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to create notification for task assignment")
		}
	})

	bus.Subscribe("task.completed", func(e Event) {
		data := e.Data.(map[string]interface{})
		taskID := data["task_id"].(uint64)
		title := data["title"].(string)

		projectIDRaw, ok := data["project_id"]
		if !ok {
			return
		}
		projectID := projectIDRaw.(uint64)

		members, err := memberRepo.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification task.completed: failed to fetch project members")
			return
		}
		for _, m := range members {
			if m.ProjectRole != projectentity.RoleProjectManager {
				continue
			}
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "TASK_COMPLETED",
				Title:      "Task Completed",
				Message:    "Task \"" + title + "\" has been marked as done.",
				EntityType: "TASK",
				EntityID:   &taskID,
				ActionURL:  "/projects/" + uint64ToStr(projectID),
			})
		}
	})

	bus.Subscribe("task.overdue", func(e Event) {
		data := e.Data.(map[string]interface{})
		taskID := data["task_id"].(uint64)
		title := data["title"].(string)

		assigneeIDRaw, ok := data["assignee_id"]
		if !ok {
			return
		}
		assigneeID := assigneeIDRaw.(uint64)

		projectIDRaw, ok := data["project_id"]
		if !ok {
			return
		}
		projectID := projectIDRaw.(uint64)

		_ = notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     assigneeID,
			Type:       "TASK_OVERDUE",
			Title:      "Task Overdue",
			Message:    "Task \"" + title + "\" is overdue.",
			EntityType: "TASK",
			EntityID:   &taskID,
			ActionURL:  "/projects/" + uint64ToStr(projectID) + "/tasks/" + uint64ToStr(taskID),
		})
	})

	bus.Subscribe("task.due_soon", func(e Event) {
		data := e.Data.(map[string]interface{})
		taskID := data["task_id"].(uint64)
		title := data["title"].(string)

		assigneeIDRaw, ok := data["assignee_id"]
		if !ok {
			return
		}
		assigneeID := assigneeIDRaw.(uint64)

		projectIDRaw, ok := data["project_id"]
		if !ok {
			return
		}
		projectID := projectIDRaw.(uint64)

		dueDateRaw, ok := data["due_date"]
		if !ok {
			return
		}
		dueDate := dueDateRaw.(string)

		_ = notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     assigneeID,
			Type:       "TASK_DUE_SOON",
			Title:      "Task Due Soon",
			Message:    "Task \"" + title + "\" is due on " + dueDate + ".",
			EntityType: "TASK",
			EntityID:   &taskID,
			ActionURL:  "/projects/" + uint64ToStr(projectID) + "/tasks/" + uint64ToStr(taskID),
		})
	})

	bus.Subscribe("budget.warning", func(e Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)
		usagePct := data["usage_pct"].(float64)

		members, err := memberRepo.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification budget.warning: failed to fetch project members")
			return
		}
		for _, m := range members {
			if m.ProjectRole != projectentity.RoleProjectManager {
				continue
			}
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "BUDGET_WARNING",
				Title:      "Budget Warning",
				Message:    "Project budget usage has crossed the 80% threshold (" + strconv.FormatFloat(usagePct, 'f', 1, 64) + "%).",
				EntityType: "PROJECT",
				EntityID:   &projectID,
				ActionURL:  "/projects/" + uint64ToStr(projectID),
			})
		}
	})

	bus.Subscribe("budget.over_limit", func(e Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)
		usagePct := data["usage_pct"].(float64)

		members, err := memberRepo.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification budget.over_limit: failed to fetch project members")
			return
		}
		for _, m := range members {
			if m.ProjectRole != projectentity.RoleProjectManager {
				continue
			}
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "BUDGET_OVER_LIMIT",
				Title:      "Budget Over Limit",
				Message:    "Project budget usage has exceeded 100% of allocated budget (" + strconv.FormatFloat(usagePct, 'f', 1, 64) + "%).",
				EntityType: "PROJECT",
				EntityID:   &projectID,
				ActionURL:  "/projects/" + uint64ToStr(projectID),
			})
		}
	})

	bus.Subscribe("milestone.completed", func(e Event) {
		data := e.Data.(map[string]interface{})
		milestoneID := data["milestone_id"].(uint64)
		projectID := data["project_id"].(uint64)
		name := data["name"].(string)

		members, err := memberRepo.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification milestone.completed: failed to fetch project members")
			return
		}
		for _, m := range members {
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "MILESTONE_COMPLETED",
				Title:      "Milestone Completed",
				Message:    "Milestone \"" + name + "\" has been completed.",
				EntityType: "MILESTONE",
				EntityID:   &milestoneID,
				ActionURL:  "/projects/" + uint64ToStr(projectID) + "/milestones/" + uint64ToStr(milestoneID),
			})
		}
	})

	bus.Subscribe("project.completed", func(e Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)

		members, err := memberRepo.FindActiveByProject(projectID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification project.completed: failed to fetch project members")
			return
		}
		for _, m := range members {
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     m.UserID,
				Type:       "PROJECT_COMPLETED",
				Title:      "Project Completed",
				Message:    "The project has been marked as completed.",
				EntityType: "PROJECT",
				EntityID:   &projectID,
				ActionURL:  "/projects/" + uint64ToStr(projectID),
			})
		}
	})

	bus.Subscribe("handover.created", func(e Event) {
		data := e.Data.(map[string]interface{})
		handoverID := data["handover_id"].(uint64)
		senderID := data["sender_id"].(uint64)
		receiverIDRaw, hasReceiver := data["receiver_id"].(*uint64)

		if hasReceiver && receiverIDRaw != nil {
			_ = notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     *receiverIDRaw,
				Type:       "HANDOVER_CREATED",
				Title:      "New Handover Incoming",
				Message:    "A document handover has been sent to you.",
				EntityType: "HANDOVER",
				EntityID:   &handoverID,
				ActionURL:  "/handovers/" + uint64ToStr(handoverID),
			})
		}

		_ = notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     senderID,
			Type:       "HANDOVER_CREATED_CONFIRM",
			Title:      "Handover Created",
			Message:    "Your handover has been created and sent.",
			EntityType: "HANDOVER",
			EntityID:   &handoverID,
			ActionURL:  "/handovers/" + uint64ToStr(handoverID),
		})
	})

	bus.Subscribe("handover.sent", func(e Event) {
		data := e.Data.(map[string]interface{})
		handoverID := data["handover_id"].(uint64)
		receiverID, hasReceiver := data["receiver_id"].(*uint64)

		if hasReceiver && receiverID != nil {
			err := notifSvc.Create(notificationservice.CreateNotificationParams{
				UserID:     *receiverID,
				Type:       "HANDOVER_SENT",
				Title:      "New Handover Incoming",
				Message:    "A document handover has been sent to you.",
				EntityType: "HANDOVER",
				EntityID:   &handoverID,
				ActionURL:  "/handovers/" + uint64ToStr(handoverID),
			})
			if err != nil {
				logger.Log.Error().Err(err).Msg("failed to create notification for handover sent")
			}
		}
	})

	bus.Subscribe("handover.received", func(e Event) {
		data := e.Data.(map[string]interface{})
		handoverID := data["handover_id"].(uint64)

		handover, err := handoverRepo.FindByID(handoverID)
		if err != nil {
			logger.Log.Error().Err(err).Msg("notification handover.received: failed to fetch handover")
			return
		}
		_ = notifSvc.Create(notificationservice.CreateNotificationParams{
			UserID:     handover.SenderID,
			Type:       "HANDOVER_RECEIVED",
			Title:      "Handover Received",
			Message:    "Your sent handover has been marked as received.",
			EntityType: "HANDOVER",
			EntityID:   &handoverID,
			ActionURL:  "/handovers/" + uint64ToStr(handoverID),
		})
	})
}

func uint64ToStr(v uint64) string {
	if v == 0 {
		return "0"
	}
	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}
