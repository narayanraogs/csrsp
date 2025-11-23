// Package optical handles the processing of video/image data frames.
package optical

import (
	"csrspServer/db"
	"csrspServer/db/sqlc"
	"csrspServer/session"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
)

// GetProcessor is the main entry point for the optical processing block.
// It sets up the worker pool and returns a function to start the pipeline.
func GetProcessor(s *session.Service) (func(*sync.WaitGroup), bool) {
	// 1. Get all video frame types and their corresponding frame IDs.
	videoFrames, err := db.GetAllOpticalDetails()
	if err != nil {
		slog.Error("Failed to get Video FrameType map from DB", "error", err)
		return nil, false
	}

	if len(videoFrames) == 0 {
		slog.Info("No video frames configured for processing.")
		return nil, false
	}

	// 2. Create and start the shared optical processor with its worker pool.
	processor := NewProcessor(s)

	var tasks []func()

	// 3. For each video channel, create a task to process its frames.
	for frameType, frameInfo := range videoFrames {
		for _, frameID := range frameInfo {
			// Use local variables for the closure.
			ft, fid := frameType, frameID.Frameid
			info := frameID

			task := func() {
				processChannel(s, processor, ft, int(fid), info)
			}
			tasks = append(tasks, task)
		}
	}

	// 4. Return a single function that runs all channel processors and handles shutdown.
	runFunc := func(wg *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in optical.GetProcessor runFunc", "panic", r, "stack", string(debug.Stack()))
			}
			// Ensure the processor and its workers are shut down.
			processor.Shutdown()
			wg.Done()
		}()

		var childWg sync.WaitGroup
		childWg.Add(len(tasks))

		for _, task := range tasks {
			go func(t func()) {
				defer childWg.Done()
				t()
			}(task)
		}

		childWg.Wait()
	}

	return runFunc, true
}

// processChannel reads frames from a single channel and sends them to the processor.
func processChannel(s *session.Service, p *Processor, frameType string, frameID int, info sqlc.Videodataprocessing) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in optical.processChannel", "frameType", frameType, "frameID", frameID, "panic", r, "stack", string(debug.Stack()))
		}
	}()

	inputChan := s.GetChannel(frameID)

	for payload := range inputChan.RegisterAsConsumer() {
		frame := &FrameData{
			FrameType:       frameType,
			FrameIdentifier: fmt.Sprintf("%d", frameID), // Or get a more descriptive identifier
			Payload:         payload,
			Width:           int(info.Noofpixelsperline),
			Height:          int(info.Nooflinesperframe),
			BitsPerPixel:    int(info.Noofbitsperpixelcheckout),
		}
		p.ProcessFrame(frame)
	}
}
