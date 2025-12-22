package global

import (
	"sync"
)

// Config is the single, canonical source of truth for all static
// application configuration, aggregated from files and the database.
type Config struct {
	// Paths
	RootPath               string
	AcqPath                string
	TempPath               string
	ArchPath               string
	LogFileDirectory       string
	AssetPath              string
	WebPath                string
	DevOpsPath             string
	ProcessingSequencePath string
	ResultNamesPath        string

	// Database Credentials
	DBServerIP string
	DBName     string
	DBUser     string
	DBPassword string

	// Core Parameters
	SystemName            string
	LocalSystemName       string
	SatName               string
	NumberOfFramesInBlock int64

	// Network & Ports
	WhiteListedIPs []string
	PCCList        []PCCInfo
	DasPorts       DasPorts

	// Hardware & Processing Parameters
	DasLockStatus DasLockStatus

	// Definitions
	Results *Results

	// From Database at Startup
	TestPhase string
}

// State holds dynamic, application-wide runtime state.
type State struct {
	mu                       sync.RWMutex
	logLevel                 string
	enableAutoArchival       bool
	enableParallelProcessing bool
	encryptionMode           string
	endProcessID             string
	enableBERPlotRealTime    bool
	enableFilterConfig       bool
	nextDisplayID            int
	nextResultID             int

	// Acquisition Lock
	isAcqActive       bool
	acquiringClientID string
	ConfiguredAcqMode string
	ChainLockStatus   map[string]string
}

// App is a singleton holding all global configuration and state.
var App struct {
	Config
	State
	DeveloperOptions
}

// Init initializes the global App object. This should be called once at startup.
func Init(cfg Config) {
	App.Config = cfg
	// Initialize default dynamic state values here
	App.State.nextDisplayID = 10 // Default starting display number
	App.State.ChainLockStatus = make(map[string]string)
}

// --- Thread-Safe Getters and Setters for Dynamic State ---

// TryLockAcquisition attempts to lock the acquisition resource for a client.
// It returns true if the lock was acquired, false otherwise.
func (s *State) TryLockAcquisition(clientID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isAcqActive {
		return false // Already locked
	}

	s.isAcqActive = true
	s.acquiringClientID = clientID
	return true
}

// ReleaseAcquisition releases the lock, verifying the clientID.
func (s *State) ReleaseAcquisition(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.acquiringClientID == clientID {
		s.isAcqActive = false
		s.acquiringClientID = ""
	}
}

// IsAcqActive returns true if an acquisition is in progress and the client holding the lock.
func (s *State) IsAcqActive() (bool, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isAcqActive, s.acquiringClientID
}

func (s *State) LogLevel() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.logLevel
}

func (s *State) SetLogLevel(level string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logLevel = level
}

func (s *State) GetNextDisplayID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextDisplayID++
	return s.nextDisplayID
}

func (s *State) AutoArchival() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enableAutoArchival
}

func (s *State) SetAutoArchival(archival bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enableAutoArchival = archival
}

// SetConfiguredAcqMode sets the currently configured acquisition mode.
func (s *State) SetConfiguredAcqMode(mode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ConfiguredAcqMode = mode
}

// GetConfiguredAcqMode returns the currently configured acquisition mode.
func (s *State) GetConfiguredAcqMode() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ConfiguredAcqMode
}

// GetChainLockStatus returns a copy of the current chain lock status map.
func (s *State) GetChainLockStatus() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to prevent race conditions on the map itself
	statusCopy := make(map[string]string)
	for key, value := range s.ChainLockStatus {
		statusCopy[key] = value
	}
	return statusCopy
}

// SetChainLockStatus updates the entire chain lock status map.
func (s *State) SetChainLockStatus(newStatus map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ChainLockStatus = newStatus
}

func SetDeveloperOptions(options DeveloperOptions) {
	App.DeveloperOptions = options
}
