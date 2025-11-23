package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// AlignmentConfig holds the configuration for an Alignment stage.
type AlignmentConfig struct {
	ValueProviders []func(payload []byte) (string, bool)
	SyncValues     []string
	MaxFrames      int // The cutoff you suggested.
}

// candidate represents a frame that is a candidate for a sync point.
type candidate struct {
	frameIndex int
	fixedKey   string         // Key from fixed-value params (e.g., "10;")
	anyValues  map[int]string // Values of 'ANY' params (e.g., {5: "ABC"})
}

// aligner holds the state for the alignment stage.
type aligner struct {
	config          AlignmentConfig
	inputs          []<-chan pipeline.Frame
	outputs         []chan pipeline.Frame
	errChan         chan<- error
	ctx             context.Context
	frameBuffers    [][][]byte
	synced          bool
	syncIndices     []int
	framesProcessed int
	maxFrames       int
}

// NewAlignmentStage creates a new many-to-many stage that synchronizes multiple input streams.
func NewAlignmentStage(config AlignmentConfig, errChan chan<- error) pipeline.StageManyToMany {
	return func(ctx context.Context, inputs []<-chan pipeline.Frame) []<-chan pipeline.Frame {
		numStreams := len(inputs)
		outputs := make([]chan pipeline.Frame, numStreams)
		readOnlyOutputs := make([]<-chan pipeline.Frame, numStreams)
		for i := 0; i < numStreams; i++ {
			outputs[i] = make(chan pipeline.Frame)
			readOnlyOutputs[i] = outputs[i]
		}

		a := &aligner{
			config:       config,
			inputs:       inputs,
			outputs:      outputs,
			errChan:      errChan,
			ctx:          ctx,
			frameBuffers: make([][][]byte, numStreams),
			syncIndices:  make([]int, numStreams),
			maxFrames:    config.MaxFrames,
		}

		go a.run()

		return readOnlyOutputs
	}
}

func (a *aligner) run() {
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("panic in Alignment stage: %v", r)
		}
		for _, c := range a.outputs {
			close(c)
		}
	}()

	a.findSyncPoint()

	if a.synced {
		a.streamSyncedOutput()
	} else {
		// Failure case: cutoff reached. Stream all buffered frames.
		a.errChan <- fmt.Errorf("alignment: sync point not found within %d frames", a.maxFrames)
		a.streamAllOutput()
	}
}

// getSyncKeyAndAnyValues separates fixed-value keys from 'ANY' values.
func (a *aligner) getSyncKeyAndAnyValues(payload []byte) (string, map[int]string, bool) {
	var keyBuilder strings.Builder
	anyValues := make(map[int]string)

	for i, provider := range a.config.ValueProviders {
		val, ok := provider(payload)
		if !ok {
			return "", nil, false
		}

		if a.config.SyncValues[i] != "ANY" {
			if val != a.config.SyncValues[i] {
				return "", nil, false
			}
			keyBuilder.WriteString(val)
			keyBuilder.WriteString(";")
		} else {
			anyValues[i] = val
		}
	}
	return keyBuilder.String(), anyValues, true
}

func (a *aligner) findSyncPoint() {
	candidates := make([][]candidate, len(a.inputs))

	for {
		if a.framesProcessed >= a.maxFrames {
			return // Cutoff reached
		}

		// Read one frame from each input to create a "tick group"
		tickGroup := make([]pipeline.Frame, len(a.inputs))
		for i, inputChan := range a.inputs {
			select {
			case frame, ok := <-inputChan:
				if !ok {
					return // An input channel closed, cannot find sync.
				}
				tickGroup[i] = frame
				a.frameBuffers[i] = append(a.frameBuffers[i], frame.Payload)
			case <-a.ctx.Done():
				return
			}
		}

		// Generate candidates for the current tick
		for i, frame := range tickGroup {
			fixedKey, anyValues, ok := a.getSyncKeyAndAnyValues(frame.Payload)
			if ok {
				candidates[i] = append(candidates[i], candidate{
					frameIndex: a.framesProcessed,
					fixedKey:   fixedKey,
					anyValues:  anyValues,
				})
			}
		}

		// Search for a sync point using the latest candidates from the first stream
		if len(candidates[0]) > 0 {
			primaryCandidate := candidates[0][len(candidates[0])-1]
			potentialSyncIndices := make([]int, len(a.inputs))
			potentialSyncIndices[0] = primaryCandidate.frameIndex
			matchCount := 1

			for i := 1; i < len(a.inputs); i++ {
				foundInStream := false
				for _, otherCandidate := range candidates[i] {
					if otherCandidate.fixedKey == primaryCandidate.fixedKey && reflect.DeepEqual(otherCandidate.anyValues, primaryCandidate.anyValues) {
						potentialSyncIndices[i] = otherCandidate.frameIndex
						foundInStream = true
						break
					}
				}
				if foundInStream {
					matchCount++
				}
			}

			if matchCount == len(a.inputs) {
				a.synced = true
				a.syncIndices = potentialSyncIndices
				return // Sync point found!
			}
		}

		a.framesProcessed++
	}
}

func (a *aligner) streamSyncedOutput() {
	var wg sync.WaitGroup
	wg.Add(len(a.outputs))

	for i, outputChan := range a.outputs {
		go func(channelIndex int, out chan<- pipeline.Frame) {
			defer wg.Done()
			// Stream buffered frames from the sync point onwards
			for frameIdx := a.syncIndices[channelIndex]; frameIdx < len(a.frameBuffers[channelIndex]); frameIdx++ {
				out <- pipeline.Frame{Payload: a.frameBuffers[channelIndex][frameIdx]}
			}
			// Continue streaming from the live channel
			for frame := range a.inputs[channelIndex] {
				out <- frame
			}
		}(i, outputChan)
	}
	wg.Wait()
}

func (a *aligner) streamAllOutput() {
	var wg sync.WaitGroup
	wg.Add(len(a.outputs))

	for i, outputChan := range a.outputs {
		go func(channelIndex int, out chan<- pipeline.Frame) {
			defer wg.Done()
			// Stream all buffered frames
			for _, payload := range a.frameBuffers[channelIndex] {
				out <- pipeline.Frame{Payload: payload}
			}
			// Continue streaming from the live channel
			for frame := range a.inputs[channelIndex] {
				out <- frame
			}
		}(i, outputChan)
	}
	wg.Wait()
}
