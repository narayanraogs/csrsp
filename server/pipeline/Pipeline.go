// Package pipeline provides a generic framework for building multi-stage, concurrent data processing pipelines.
package pipeline

import (
	"context"
	"log/slog"
	"sync"
)

// Frame is the unit of data that flows through the pipeline.
type Frame struct {
	FrameType       string // e.g., "AUX_HK_1"
	FrameIdentifier string // The specific FrameID, e.g., "101"
	ID              int
	Payload         []byte
}

// Stage functions define the different shapes of a processing stage.
type StageOneToOne func(context.Context, <-chan Frame) <-chan Frame
type StageOneToMany func(context.Context, <-chan Frame) []<-chan Frame
type StageManyToOne func(context.Context, []<-chan Frame) <-chan Frame
type StageManyToMany func(context.Context, []<-chan Frame) []<-chan Frame
type Source func(context.Context) <-chan Frame // Add this line
type Sink func(context.Context, <-chan Frame)

// Factory functions are used to create instances of stages.
type StageFactoryOneToOne func(errChan chan<- error) StageOneToOne
type StageFactoryOneToMany func(errChan chan<- error) StageOneToMany
type StageFactoryManyToOne func(errChan chan<- error) StageManyToOne
type StageFactoryManyToMany func(errChan chan<- error) StageManyToMany
type SourceFactory func(errChan chan<- error) Source // Add this line
type SinkFactory func(errChan chan<- error) Sink

// Pipeline manages the execution of a series of connected stages.
type Pipeline struct {
	runFuncs      []func(context.Context)
	errChan       chan error
	TelemetryChan chan TelemetryEvent // Add this line
	wg            sync.WaitGroup
	runMutex      sync.Mutex
	isRunning     bool
}

// New creates a new, empty Pipeline.
func New() *Pipeline {
	return &Pipeline{
		errChan:       make(chan error, 100),
		TelemetryChan: make(chan TelemetryEvent, 100), // Add this line
	}
}

// ErrorChannel returns the channel from which pipeline errors can be read.
func (p *Pipeline) ErrorChannel() <-chan error {
	return p.errChan
}

// AddStage adds a one-to-one processing stage to the pipeline.
func (p *Pipeline) AddStage(factory StageFactoryOneToOne, input <-chan Frame) <-chan Frame {
	stage := factory(p.errChan)
	output := make(chan Frame)

	p.runFuncs = append(p.runFuncs, func(ctx context.Context) {
		stageOutput := stage(ctx, input)
		for frame := range stageOutput {
			output <- frame
		}
		close(output)
	})

	return output
}

// AddSink adds a terminal stage (a sink) to the pipeline.
func (p *Pipeline) AddSink(factory SinkFactory, input <-chan Frame) {
	sink := factory(p.errChan)

	p.runFuncs = append(p.runFuncs, func(ctx context.Context) {
		// A sink has no output to forward, so we just run it.
		sink(ctx, input)
	})
}

// AddSource adds a source stage to the pipeline, which generates the initial data.
func (p *Pipeline) AddSource(factory SourceFactory) <-chan Frame {
	source := factory(p.errChan)
	output := make(chan Frame)

	p.runFuncs = append(p.runFuncs, func(ctx context.Context) {
		// A source has no input, it only produces an output channel.
		sourceOutput := source(ctx)
		for frame := range sourceOutput {
			output <- frame
		}
		close(output)
	})

	return output
}

// AddManyToOne adds a many-to-one processing stage to the pipeline.
func (p *Pipeline) AddManyToOne(factory StageFactoryManyToOne, inputs []<-chan Frame) <-chan Frame {
	stage := factory(p.errChan)
	output := make(chan Frame)

	p.runFuncs = append(p.runFuncs, func(ctx context.Context) {
		stageOutput := stage(ctx, inputs)
		for frame := range stageOutput {
			output <- frame
		}
		close(output)
	})

	return output
}

// Run starts the execution of all stages in the pipeline.
func (p *Pipeline) Run(ctx context.Context) {
	p.runMutex.Lock()
	if p.isRunning {
		p.runMutex.Unlock()
		slog.Warn("Pipeline.Run() called more than once.")
		return
	}
	p.isRunning = true
	p.runMutex.Unlock()

	p.wg.Add(len(p.runFuncs) + 1) // Add one for the supervisor

	go p.supervisor(ctx)

	for _, runFunc := range p.runFuncs {
		go func(rf func(context.Context)) {
			defer p.wg.Done()
			rf(ctx)
		}(runFunc)
	}
}

// Wait blocks until all stages in the pipeline have completed.
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

// supervisor runs in the background, listening for errors and context cancellation.
func (p *Pipeline) supervisor(ctx context.Context) {
	defer p.wg.Done()
	slog.Info("Pipeline supervisor started.")

	for {
		select {
		case err, ok := <-p.errChan:
			if !ok {
				slog.Info("Error channel closed, pipeline supervisor shutting down.")
				return
			}
			// Centralized error logging.
			slog.Error("Pipeline Error", "error", err)
		case <-ctx.Done():
			slog.Info("Context cancelled, pipeline supervisor shutting down.")
			// Drain the error channel before exiting
			for err := range p.errChan {
				slog.Error("Pipeline Error (draining on shutdown)", "error", err)
			}
			return
		}
	}
}
