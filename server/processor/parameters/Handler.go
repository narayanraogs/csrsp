// Package parameters handles the validation of header parameters.
package parameters

import (
	"csrsp/server/db"
	"csrsp/server/db/sqlc"
	"csrsp/server/session"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"sync"
)

// GetProcessor is the main entry point for the header validation processing block.
// It constructs and returns a function that runs the entire validation pipeline.
func GetProcessor(s *session.Service) (func(*sync.WaitGroup), bool) {
	// 1. PREPARE: Get a map of FrameType -> list of FrameIDs to process.
	frameTypes, err := db.GetAllAuxFrameTypes()
	if err != nil {
		slog.Error("Failed to get Aux FrameTypes from DB", "error", err)
		return nil, false
	}

	var tasks []func()

	// 2. CREATE TASKS: For each frame channel, create a task to process it.
	for _, frameType := range frameTypes {
		frameIDs, err := db.GetFrameTypeFrameIdentifiers(frameType)
		if err != nil {
			slog.Error("Failed to get Aux Frames for Type from DB", "FrameType", frameType, "error", err)
			return nil, false
		}
		for _, frameID := range frameIDs {
			// Use local variables for the closure to ensure they are captured correctly.
			ft, fid := frameType, frameID

			task := func() {
				processFrameChannel(s, ft, int(fid.Frameid))
			}
			tasks = append(tasks, task)
		}
	}

	// 3. RETURN: Return a single function that launches all tasks in parallel.
	runFunc := func(wg *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in parameters.GetProcessor runFunc", "panic", r, "stack", string(debug.Stack()))
			}
			wg.Done()
		}()

		var childWg sync.WaitGroup
		childWg.Add(len(tasks))

		for _, task := range tasks {
			go func(t func()) {
				defer childWg.Done()
				t() // Execute the task (which calls processFrameChannel)
			}(task)
		}

		childWg.Wait()
		slog.Info("All parameter validation tasks have completed.")
	}
	return runFunc, true
}

// processFrameChannel handles all validation for a single frame channel.
func processFrameChannel(s *session.Service, frameType string, frameID int) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in processFrameChannel", "frameType", frameType, "frameID", frameID, "panic", r, "stack", string(debug.Stack()))
		}
	}()

	// 1. Get all parameters for this frameType and create their validators.
	params, err := db.GetFrameTypeParamterDetails(frameType)
	if err != nil || len(params) == 0 {
		return // No parameters to validate for this frame type.
	}

	var validators []Validator
	for _, p := range params {
		v, err := createValidator(p)
		if err == nil && v != nil {
			validators = append(validators, v)
		} else if err != nil {
			slog.Error("Failed to create validator", "param", p.Parametername, "error", err)
		}
	}

	if len(validators) == 0 {
		return // No valid validators could be created.
	}

	// 2. Get the single input channel for this frame.
	inputChan := s.GetChannel(frameID)

	// 3. Create an intermediate channel for each validator (fan-out).
	intermediateChans := make([]chan []byte, len(validators))
	for i := range validators {
		intermediateChans[i] = make(chan []byte, 100) // Use a small buffer
	}

	// 4. Launch one goroutine per validator.
	var wg sync.WaitGroup
	wg.Add(len(validators))
	for i, v := range validators {
		go func(validator Validator, in <-chan []byte) {
			defer wg.Done()
			key := fmt.Sprintf("%s;%d;%s", frameType, frameID, validator.ParamName())
			validator.Process(in, s.Store, key)
		}(v, intermediateChans[i])
	}

	// 5. Main loop: Broadcast incoming frames to all validator channels.
	for frame := range inputChan.RegisterAsConsumer() {
		for _, ch := range intermediateChans {
			ch <- frame
		}
	}

	// 6. Close intermediate channels and wait for validators to finish.
	for _, ch := range intermediateChans {
		close(ch)
	}
	wg.Wait()
}

// createValidator is a factory function that builds the correct validator.
func createValidator(param sqlc.Parameterdetail) (Validator, error) {
	pType := string(param.Parametertype)
	switch strings.ToLower(pType) {
	case "analog":
		return NewAnalogValidator(param.Paramid, param.Parametername)
	case "radix":
		return NewRadixValidator(param.Paramid, param.Parametername)
	case "status":
		return NewStatusValidator(param.Paramid, param.Parametername)
	case "fsc":
		return NewFSCValidator(param.Paramid, param.Parametername)
	case "increment":
		return NewIncrementValidator(param.Paramid, param.Parametername)
	case "crc":
		return NewCRCValidator(param.Paramid, param.Parametername)
	default:
		return nil, fmt.Errorf("unknown parameter type: %s", param.Parametertype)
	}
}
