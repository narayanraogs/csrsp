package pipeline

import (
	"csrsp/server/session"
	"fmt"
	"log/slog"
)

// TelemetrySink runs as a goroutine, listening for telemetry events and
// updating the session store accordingly.
func TelemetrySink(telemetryChan <-chan TelemetryEvent, store *session.Store) {
	slog.Info("Telemetry sink started.")
	for event := range telemetryChan {
		switch event.EventType {
		case EventTypeRSCorrection:
			details, ok := event.Details.(RSCorrectionDetails)
			if !ok {
				slog.Error("received RSCorrection event with invalid details type")
				continue
			}

			var msg string
			if details.IsUncorrectable {
				msg = "Uncorrectable"
			} else {
				msg = fmt.Sprintf("Corrected %d errors", details.ErrorsCorrected)
			}

			rsEvent := session.RSEvent{
				LineNumber:      details.LineNumber,
				IsUncorrectable: details.IsUncorrectable,
				IsCorrectable:   !details.IsUncorrectable,
				BlockNumber:     details.BlockNumber,
				ErrorsCorrected: details.ErrorsCorrected,
				Message:         msg,
			}
			store.RSDecoding.AddEvent(rsEvent)

		// Add more cases here for other event types in the future.
		default:
			slog.Warn("received unknown telemetry event type", "type", event.EventType)
		}
	}
	slog.Info("Telemetry sink finished.")
}
