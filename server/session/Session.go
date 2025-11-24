package session

import (
	"csrsp/server/global"

	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// --- Per-Session Structs ---

// Config holds the read-only parameters for a specific job within a session.
// This struct is typically populated once at the beginning of a job.
type Config struct {
	AcqMode        string
	ConfigName     string
	Remarks        string
	NumberOfFrames int64
	IsTimeBased    bool
	TimeInSec      int64
	IsBlindAcq     bool

	// Offline/File-based Job Parameters
	OfflineDate string
	OfflineTime string

	// User/Result Context
	SelectedProfile string
}

// State holds the dynamic, read-write state for a specific job.
// All access to this struct's fields must be protected by its mutex.
type State struct {
	mu                 sync.RWMutex
	isProcessingActive bool
	isProcessingCutOff bool
	userID             int64
	privileges         []string
	acqDate            string
	acqTime            string
	noOfChains         int
	acqStartTime       time.Time

	// Live acquisition status
	acqChainStatus map[int]bool
	acqFrameCounts map[int]int64
}

// Service is the core object representing a single client session.
// It holds the session's configuration, dynamic state, and data pipelines.
type Service struct {
	ClientID   string
	Config     Config
	State      State
	ChannelMap map[int]*Channel

	// Result and Display Management
	//ResultMap  map[int]*results.Format
	//Queue      chan results.Result
	Closer     chan bool
	displayNos []int
	Store      *Store

	channelMutex sync.Mutex
	aliveTimer   int64 // Unix timestamp of the last heartbeat
}

// --- Session Management ---

var (
	serviceMapper map[string]*Service
	mapMutex      sync.Mutex
)

func init() {
	serviceMapper = make(map[string]*Service)
}

// GetService retrieves an existing Service instance for the given clientID or creates a new one.
// This function is thread-safe.
func GetService(clientID string) *Service {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if service, exists := serviceMapper[clientID]; exists {
		service.heartbeat() // Refresh heartbeat on access
		return service
	}

	// Create and initialize a new service
	newService := &Service{
		ClientID:   clientID,
		ChannelMap: make(map[int]*Channel),
		//ResultMap:  make(map[int]*results.Format),
		displayNos: make([]int, 0),
		Store:      NewStore(),
	}
	newService.State.acqChainStatus = make(map[int]bool)
	newService.State.acqFrameCounts = make(map[int]int64)
	newService.heartbeat()

	// Create dedicated temp directories
	// Note: Error handling should be added here in a real implementation
	_ = newService.createSessionDirs()

	// Start the heartbeat monitor for this service
	go newService.monitorHeartbeat()

	serviceMapper[clientID] = newService
	return newService
}

// RemoveService deletes a service instance, allowing for cleanup.
func RemoveService(clientID string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if s, ok := serviceMapper[clientID]; ok {
		s.ResetDisplay()
		delete(serviceMapper, clientID)
		fmt.Println("Service Removed for ClientID= " + clientID)
	}
}

// --- Service Methods ---

// heartbeat updates the last-seen time for the session.
func (s *Service) heartbeat() {
	s.aliveTimer = time.Now().Unix()
}

// monitorHeartbeat checks periodically if the client is still alive.
// If no heartbeat is received for a timeout period, it cleans up the service.
func (s *Service) monitorHeartbeat() {
	const timeout = 30 // seconds
	// Do not monitor the 'remote' client as it has no heartbeat mechanism
	if s.ClientID == "remote" {
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if time.Now().Unix()-s.aliveTimer > timeout {
			fmt.Println("Client timeout, removing service for", s.ClientID)
			s.CloseResults()
			RemoveService(s.ClientID)
			return // End the monitoring goroutine
		}
	}
}

// createSessionDirs sets up dedicated temporary and acquisition directories for the session.
func (s *Service) createSessionDirs() error {
	// Create temp path
	tempPath := filepath.Join(global.App.Config.TempPath, s.ClientID)
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return err
	}

	// Create acquisition path
	acqPath := filepath.Join(global.App.Config.AcqPath, s.ClientID)
	return os.MkdirAll(acqPath, 0755)
}

// GetChannel returns a channel for the given frameID, creating it if it doesn't exist.
func (s *Service) GetChannel(frameID int) *Channel {
	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()

	if channel, exists := s.ChannelMap[frameID]; exists {
		return channel
	}

	newChan := newChannel()
	s.ChannelMap[frameID] = newChan
	return newChan
}

// AddDisplay tracks a new display number (e.g., for XPRA) used by the session.
func (s *Service) AddDisplay(displayNo int) {
	s.displayNos = append(s.displayNos, displayNo)
}

// ResetDisplay stops all external display processes associated with the session.
func (s *Service) ResetDisplay() {
	for _, displayNo := range s.displayNos {
		args := []string{"stop", ":" + strconv.Itoa(displayNo)}
		cmd := exec.Command("xpra", args...)
		cmd.Env = os.Environ()
		if err := cmd.Start(); err != nil {
			slog.Error("Failed to stop XPRA display", "display", displayNo, "error", err)
		}
	}
	s.displayNos = make([]int, 0)
}

// CloseResults handles the graceful shutdown of result processing goroutines.
func (s *Service) CloseResults() {
	defer runtime.GC()
	defer debug.FreeOSMemory()
	s.ResetDisplay()

	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()

	if s.Closer != nil {
		s.Closer <- true
	}
	/*
		if s.Queue != nil {
			close(s.Queue)
		}

			// Clean up embedded databases if they were used
			embedded.ReleaseAuxDB(s.ClientID)
			embedded.ReleaseVideoDB(s.ClientID)
			embedded.ReleaseSARDB(s.ClientID)
	*/
	slog.Info("Result processing closed successfully for client", "clientID", s.ClientID)
}

// --- Thread-Safe State Accessors ---

func (s *State) IsProcessingActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isProcessingActive
}

func (s *State) SetProcessingActive(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isProcessingActive = active
}

func (s *State) SetNumberOfChains(noOfChains int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.noOfChains = noOfChains
	s.acqChainStatus = make(map[int]bool)
}

func (s *State) SetAcqChainStatus(chainID int, status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acqChainStatus[chainID] = status
}

func (s *State) SetAcqFrameCount(chainID int, frames int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acqFrameCounts[chainID] = frames
}

func (s *State) SetAcqStartTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acqStartTime = t
}

func (s *State) GetAcqChainStatus() map[int]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to prevent race conditions on the map itself
	statusCopy := make(map[int]bool)
	for k, v := range s.acqChainStatus {
		statusCopy[k] = v
	}
	return statusCopy
}

func (s *State) GetAcqFrameCount() map[int]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Return a copy to prevent race conditions on the map itself
	countCopy := make(map[int]int64)
	for k, v := range s.acqFrameCounts {
		countCopy[k] = v
	}
	return countCopy
}
