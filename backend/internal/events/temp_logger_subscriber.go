package events

import "github.com/Kal-el21/backend/internal/shared/logger"

// RegisterTempLoggerSubscriber adalah placeholder subscriber sampai
// Notification Handler penuh dibangun di Phase 5. Tujuannya memverifikasi
// event bus bekerja end-to-end sejak Phase 2.
func RegisterTempLoggerSubscriber(bus *Bus) {
	events := []string{
		"project.request.submitted",
		"project.request.approved",
		"project.request.rejected",
		"project.request.revision_requested",
	}

	for _, eventName := range events {
		bus.Subscribe(eventName, func(e Event) {
			logger.Log.Info().
				Str("event", e.Name).
				Interface("data", e.Data).
				Msg("event received")
		})
	}
}
