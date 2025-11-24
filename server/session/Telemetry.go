package session

import "sync"

// RSEvent represents a single Reed-Solomon correction event.
// This mirrors the data stored in the V2 server's SingleRSDecodingStatus.
type RSEvent struct {
	LineNumber      int
	IsUncorrectable bool
	IsCorrectable   bool
	BlockNumber     int
	ErrorsCorrected int
	Message         string
}

// RSDecodingResult holds a thread-safe list of all RS correction events that have occurred.
type RSDecodingResult struct {
	mutex  sync.RWMutex
	Events []RSEvent
}

// AddEvent adds a new correction event to the list.
func (r *RSDecodingResult) AddEvent(event RSEvent) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Events = append(r.Events, event)
}

// GetEvents returns a copy of all correction events.
func (r *RSDecodingResult) GetEvents() []RSEvent {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	// Return a copy to prevent race conditions on the slice itself.
	eventsCopy := make([]RSEvent, len(r.Events))
	copy(eventsCopy, r.Events)
	return eventsCopy
}
