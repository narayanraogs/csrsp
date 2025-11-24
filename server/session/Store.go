// Package session manages the state and data flow for individual client sessions.
package session

import (
	"fmt"
	"sync"
)

// --- Status/Error Struct Definitions ---

// AuxStatus holds the status for a single auxiliary parameter.
type AuxStatus struct {
	paramName      string
	err            bool
	warning        bool
	status         string
	value          string
	errUpdater     sync.Once // Ensures the error state is latched.
	warningUpdater sync.Once // Ensures the warning state is latched.
}

// SetError latches the error state to true.
func (as *AuxStatus) SetError() {
	as.errUpdater.Do(func() { as.err = true })
}

// SetWarning latches the warning state to true.
func (as *AuxStatus) SetWarning() {
	as.warningUpdater.Do(func() { as.warning = true })
}

// --- Domain-Specific Stores ---

// AuxStore provides a thread-safe, in-memory store for auxiliary processing data.
type AuxStore struct {
	mu     sync.RWMutex
	status map[string]*AuxStatus
	errors []AuxError
}

type AuxError struct {
	paramName     string
	lineNo        int64
	isWarning     bool
	expectedValue string
	actualValue   string
	errorMessage  string
}

// VideoStore provides a thread-safe, in-memory store for video processing data.
type VideoStore struct {
	mu     sync.RWMutex
	status map[string]*VideoStatus
}

// SARStore provides a thread-safe, in-memory store for SAR processing data.
type SARStore struct {
	mu     sync.RWMutex
	status map[string]*SARVideoStatus
}

// --- Main Store Container ---

// Store is the main container for all session-specific, in-memory data stores.
type Store struct {
	Aux        *AuxStore
	Video      *VideoStore
	Microwave  *MicrowaveStore
	RSDecoding *RSDecodingResult // Add this line
	mutex      sync.RWMutex
}

// NewStore creates a fully initialized main Store container.
func NewStore() *Store {
	return &Store{
		Aux:        NewAuxStore(),
		Video:      NewVideoStore(),
		Microwave:  NewMicrowaveStore(),
		RSDecoding: &RSDecodingResult{}, // Add this line
	}
}

// --- AuxStore Methods ---

// NewAuxStore creates an initialized store for auxiliary data.
func NewAuxStore() *AuxStore {
	return &AuxStore{
		status: make(map[string]*AuxStatus),
		errors: make([]AuxError, 0, 1000),
	}
}

// InitStatuses creates and initializes the status objects for a given set of parameters.
func (s *AuxStore) InitStatuses(key string, parameters []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// This would be called once at the start of processing for a given frame type.
	for _, paramName := range parameters {
		paramKey := key + ";" + paramName
		s.status[paramKey] = &AuxStatus{paramName: paramName}
	}
}

// GetStatus retrieves a pointer to an AuxStatus object. The returned object
// should be treated as read-only by the caller.
func (s *AuxStore) GetStatus(key string) (*AuxStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stat, ok := s.status[key]
	return stat, ok
}

// UpdateStatusValue safely updates the 'value' and 'status' fields of an AuxStatus object.
func (s *AuxStore) UpdateStatusValue(key, value, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, ok := s.status[key]
	if !ok {
		return // Or handle error
	}
	stat.value = value
	stat.status = status
}

func (s *AuxStore) UpdateStatus(key, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, ok := s.status[key]
	if !ok {
		return // Or handle error
	}
	stat.status = status
}

func (s *AuxStore) UpdateValue(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, ok := s.status[key]
	if !ok {
		return // Or handle error
	}
	stat.value = value
}

func (s *AuxStore) AddError(paramName string, warning bool, lineNo int64, expValue string, actValue string, message string) {
	var e AuxError
	e.paramName = paramName
	e.isWarning = warning
	e.lineNo = lineNo
	e.expectedValue = expValue
	e.actualValue = actValue
	e.errorMessage = message
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errors = append(s.errors, e)

}

// --- VideoStore Methods (Example) ---

// VideoStatus holds the statistical data for a single pixel over time.
type VideoStatus struct {
	mean        float64
	sd          float64
	diff        float64
	isFaulty    bool
	meanErr     bool
	sdErr       bool
	meanErrOnce sync.Once
	sdErrOnce   sync.Once
	faultyOnce  sync.Once
}

// SetMeanError latches the mean error flag to true.
func (vs *VideoStatus) SetMeanError() {
	vs.meanErrOnce.Do(func() { vs.meanErr = true })
}

// SetSDError latches the standard deviation error flag to true.
func (vs *VideoStatus) SetSDError() {
	vs.sdErrOnce.Do(func() { vs.sdErr = true })
}

// SetIsFaulty latches the faulty pixel flag to true.
func (vs *VideoStatus) SetIsFaulty(isFaulty bool) {
	if isFaulty {
		vs.faultyOnce.Do(func() { vs.isFaulty = true })
	}
}

func NewVideoStore() *VideoStore {
	return &VideoStore{
		status: make(map[string]*VideoStatus),
	}
}

// InitStatuses creates and initializes the status objects for all pixels in a frame.
func (s *VideoStore) InitStatuses(frameType, frameID string, numPixels int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < numPixels; i++ {
		key := fmt.Sprintf("%s;%s;%d", frameType, frameID, i)
		s.status[key] = &VideoStatus{}
	}
}

// GetStatus retrieves a pointer to a VideoStatus object.
func (s *VideoStore) GetStatus(key string) (*VideoStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stat, ok := s.status[key]
	return stat, ok
}

// UpdatePixelStats safely updates the statistics for a single pixel.
func (s *VideoStore) UpdatePixelStats(key string, mean, sd, diff float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, ok := s.status[key]
	if !ok {
		return
	}
	stat.mean = mean
	stat.sd = sd
	stat.diff = diff
}

// --- SARStore Methods (Example) ---

// SARVideoStatus holds data for SAR processing.
type SARVideoStatus struct {
	// ... fields for SAR status
}

func NewSARStore() *SARStore {
	return &SARStore{
		status: make(map[string]*SARVideoStatus),
	}
}

// --- MicrowaveStore Methods ---

// IQStats holds the calculated statistics for an I or Q channel.
type IQStats struct {
	Min  float64
	Max  float64
	Mean float64
	SD   float64
	RMS  float64
}

// ChirpParams holds the calculated chirp parameters.
type ChirpParams struct {
	PeakAmp    float64
	PeakAmpLoc uint64
	PSLRLeft   float64
	PSLRRight  float64
	ISLR       float64
	Resolution float64
}

// MicrowaveStatus holds the calculated results for a single PRF.
type MicrowaveStatus struct {
	IStats      IQStats
	QStats      IQStats
	ChirpParams ChirpParams
}

// MicrowaveStore provides a thread-safe, in-memory store for microwave processing data.
type MicrowaveStore struct {
	mu     sync.RWMutex
	status map[string]*MicrowaveStatus
}

// NewMicrowaveStore creates and initializes a new microwave store.
func NewMicrowaveStore() *MicrowaveStore {
	return &MicrowaveStore{
		status: make(map[string]*MicrowaveStatus),
	}
}

// InitPRFStatus initializes the status object for a new PRF.
func (s *MicrowaveStore) InitPRFStatus(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[key] = &MicrowaveStatus{}
}

// UpdatePRFStatus updates the statistics and chirp parameters for a given PRF.
func (s *MicrowaveStore) UpdatePRFStatus(key string, iStats, qStats IQStats, chirp ChirpParams) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, ok := s.status[key]
	if !ok {
		return // Or handle error appropriately
	}
	stat.IStats = iStats
	stat.QStats = qStats
	stat.ChirpParams = chirp
}
