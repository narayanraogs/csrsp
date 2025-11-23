package microwave

import (
	"csrspServer/db"
	"csrspServer/processor/parameters"
	"csrspServer/session"
	"log/slog"
	"runtime/debug"
	"strings"
	"sync"
)

// GetMicrowaveProcessor is the main entry point for the microwave processing pipeline.
func GetMicrowaveProcessor(s *session.Service) (func(*sync.WaitGroup), bool) {
	sarProcessingIDs, err := db.GetAllSARProcessingIDs()
	if err != nil {
		slog.Error("SAR Processing cannot read details", "error", err.Error())
		return nil, false
	}
	if len(sarProcessingIDs) == 0 {
		slog.Warn("No SAR Processing IDs found in database.")
		return nil, false
	}

	var processors = make([]func(*sync.WaitGroup), 0)
	for _, processID := range sarProcessingIDs {
		processor, ok := getProcessorForID(int(processID), s)
		if !ok {
			slog.Error("Cannot get microwave processor for ID", "processID", processID)
			continue
		}
		processors = append(processors, processor)
	}

	// The main processing function that runs all individual processors.
	var process = func(group *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in MicrowaveProcessing", "panic", r, "stack", string(debug.Stack()))
			}
			group.Done()
		}()
		var wg sync.WaitGroup
		wg.Add(len(processors))
		for _, p := range processors {
			go p(&wg)
		}
		wg.Wait()
		slog.Info("Finished all microwave processing.")
	}

	return process, true
}

func getProcessorForID(processID int, s *session.Service) (func(group *sync.WaitGroup), bool) {
	sarProcessingStruct, err := db.GetSARProcessingDetails(processID)
	if err != nil {
		slog.Warn("Failed to get SARProcessing details for ID", "processID", processID, "error", err.Error())
		return nil, false
	}

	// --- Get Parameter Providers ---
	polarizationProvider, err := parameters.GetProcessedParameterValueProvider(sarProcessingStruct.Polarizationpid)
	if err != nil {
		slog.Error("Unable to get provider for polarization", "error", err.Error())
		return nil, false
	}
	timingStateProvider, err := parameters.GetProcessedParameterValueProvider(sarProcessingStruct.Timingstatepid)
	if err != nil {
		slog.Error("Unable to get provider for timing state", "error", err.Error())
		return nil, false
	}
	baqProvider, err := parameters.GetProcessedParameterValueProvider(sarProcessingStruct.Baqpid)
	if err != nil {
		slog.Error("Unable to get provider for BAQ", "error", err.Error())
		return nil, false
	}
	chirpBWProvider, err := parameters.GetProcessedParameterValueProvider(sarProcessingStruct.Chirpbwpid)
	if err != nil {
		slog.Error("Unable to get provider for Chirp BW", "error", err.Error())
		return nil, false
	}
	sarModeProvider, err := parameters.GetProcessedParameterValueProvider(sarProcessingStruct.Sarmodepid)
	if err != nil {
		slog.Error("Unable to get provider for SAR Mode", "error", err.Error())
		return nil, false
	}

	// --- Get Status Maps for converting uint values to strings ---
	polStatusMap, _ := db.GetStatusMap(sarProcessingStruct.Polarizationpid)
	sarModeStatusMap, _ := db.GetStatusMap(sarProcessingStruct.Sarmodepid)
	timingStateStatusMap, _ := db.GetStatusMap(sarProcessingStruct.Timingstatepid)

	// --- Register as consumer on session channels ---
	var inputChannels []<-chan []byte
	var inputFrameIDs []int
	if strings.EqualFold(sarProcessingStruct.Sarmodevalue, "stretch") {
		inputFrameIDs = append(inputFrameIDs, int(sarProcessingStruct.Payloadi1frameid.Int32), int(sarProcessingStruct.Payloadi2frameid.Int32), int(sarProcessingStruct.Payloadq1frameid.Int32), int(sarProcessingStruct.Payloadq2frameid.Int32))
	} else {
		// Simplified logic for non-stretch modes
		if sarProcessingStruct.Payloadi1frameid.Int32 != -1 {
			inputFrameIDs = append(inputFrameIDs, int(sarProcessingStruct.Payloadi1frameid.Int32), int(sarProcessingStruct.Payloadq1frameid.Int32))
		} else {
			inputFrameIDs = append(inputFrameIDs, int(sarProcessingStruct.Payloadi2frameid.Int32), int(sarProcessingStruct.Payloadq2frameid.Int32))
		}
	}

	for _, frameID := range inputFrameIDs {
		ch := s.GetChannel(frameID)
		if ch != nil {
			inputChannels = append(inputChannels, ch.RegisterAsConsumer())
		}
	}

	if len(inputChannels) == 0 {
		slog.Warn("No input channels found for microwave processor", "processID", processID)
		return nil, false
	}

	// --- Create the main processing function ---
	processFunc := func(wait *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in MicrowaveProcessorForID", slog.Int("processID", processID), "panic", r, "stack", string(debug.Stack()))
			}
			wait.Done()
		}()

		// Create our new, efficient processor
		processor := NewProcessor(s)
		defer processor.Shutdown()

		var prfCounter uint64 = 0
		var prevTimingState string = ""

		for {
			// Read one frame from each channel to create a bundle
			prfDataArray := make([][]byte, len(inputChannels))
			ok := true
			for i, ch := range inputChannels {
				var frameData []byte
				frameData, ok = <-ch
				if !ok {
					return // A channel was closed, so we are done.
				}
				prfDataArray[i] = frameData
			}

			// --- Extract Parameters ---
			params := SARProcessingParams{}
			timingStateUint := timingStateProvider(prfDataArray[0]).(uint64)
			params.TimingState = strings.ToUpper(strings.TrimSpace(timingStateStatusMap.Status[timingStateUint]))
			params.BaqValue = uint64(baqProvider(prfDataArray[0]).(float64))
			polUint := polarizationProvider(prfDataArray[0]).(uint64)
			params.Polarization = polStatusMap.Status[polUint]
			sarModeUint := sarModeProvider(prfDataArray[0]).(uint64)
			params.SarModeValue = sarModeStatusMap.Status[sarModeUint]
			params.SamplingFrequency = chirpBWProvider(prfDataArray[0]).(float64)
			params.SarModeValueFromDB = sarProcessingStruct.Sarmodevalue
			params.TimingStatesToExclude = strings.Join(strings.Split(sarProcessingStruct.Timingstatetoexclude, ","), ",")

			// --- Manage PRF Counter ---
			if prevTimingState != params.TimingState {
				if prfCounter > 0 {
					// In the future, we can use this to log timing counts if needed.
				}
				prevTimingState = params.TimingState
				prfCounter = 0
			}
			params.PrfIndexCounter = prfCounter
			prfCounter++

			// --- Check if processing is needed ---
			if !strings.EqualFold(params.SarModeValue, params.SarModeValueFromDB) {
				continue
			}
			if strings.Contains(params.TimingStatesToExclude, params.TimingState) {
				continue
			}

			// --- Dispatch to the processor ---
			processor.ProcessPRFBundle(params, prfDataArray)
		}
	}

	return processFunc, true
}
