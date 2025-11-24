package pipeline

import "time"

// TelemetryEvent represents a single piece of operational data from a pipeline stage.
type TelemetryEvent struct {
	StageName string      // The name of the stage that generated the event.
	Timestamp time.Time
	EventType string      // A specific identifier for the type of event.
	Severity  string      // e.g., "Info", "Warning", "Error"
	Details   interface{} // A payload specific to the event type.
}

// Event Types
const (
	EventTypeRSCorrection = "RSCorrection"
	// We can add more event types here later, e.g., "BufferOverrun", "SyncLoss"
)

// Severity Levels
const (
	SeverityInfo    = "Info"
	SeverityWarning = "Warning"
	SeverityError   = "Error"
)

// RSCorrectionDetails is the specific payload for an RSCorrection event.
type RSCorrectionDetails struct {
	FrameIdentifier string
	LineNumber      int
	BlockNumber     int
	ErrorsCorrected int
	IsUncorrectable bool
}
