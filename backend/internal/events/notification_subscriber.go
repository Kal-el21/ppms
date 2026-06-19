package events

import (
	notificationservice "github.com/Kal-el21/backend/internal/domain/notification/service"
	"github.com/Kal-el21/backend/internal/shared/logger"
)

// RegisterNotificationSubscriber menggantikan RegisterTempLoggerSubscriber (Phase 2).
// Setiap event domain di-mapping ke notifikasi konkret sesuai SDD section 10 (Sample Events)
// dan PRD section 13 (Notification Events).
func RegisterNotificationSubscriber(bus *Bus, notifSvc notificationservice.NotificationService) {
	bus.Subscribe("project.request.submitted", func(e Event) {
		data := e.Data.(map[string]interface{})
		requestID := data["request_id"].(uint64)
		title := data["title"].(string)

		// Notifikasi dikirim ke semua ADMIN — karena daftar ADMIN dinamis,
		// idealnya di-resolve oleh UserRepository. Untuk Phase 5, kita publish
		// sebagai event tersendiri yang akan di-handle dengan query admin list
		// di handler terpisah (lihat catatan di bawah kode ini).
		logger.Log.Info().Uint64("request_id", requestID).Str("title", title).Msg("notification: request submitted, notify admins")
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
			Message:    "Your request \"" + title + "\" requires revision or has been rejected.",
			EntityType: "PROJECT_REQUEST",
			EntityID:   &requestID,
			ActionURL:  "/project-requests/" + uint64ToStr(requestID),
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("failed to create notification for request rejection")
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

		logger.Log.Info().Uint64("task_id", taskID).Str("title", title).Msg("notification: task completed, notify PM")
		// Catatan: notify PM project terkait membutuhkan lookup project_id -> PM user_id,
		// yang idealnya dilakukan via cross-domain query. Diimplementasikan penuh
		// saat wiring di main.go menggunakan closure yang punya akses ke member repository.
	})

	bus.Subscribe("budget.warning", func(e Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)
		usagePct := data["usage_pct"].(float64)

		logger.Log.Warn().Uint64("project_id", projectID).Float64("usage_pct", usagePct).Msg("notification: budget warning, notify PM")
	})

	bus.Subscribe("budget.over_limit", func(e Event) {
		data := e.Data.(map[string]interface{})
		projectID := data["project_id"].(uint64)

		logger.Log.Warn().Uint64("project_id", projectID).Msg("notification: budget over limit, notify PM")
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

		logger.Log.Info().Uint64("handover_id", handoverID).Msg("notification: handover received, notify sender")
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
