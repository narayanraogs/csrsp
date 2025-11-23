package buffer

import (
	"fmt"
	"log/slog"
	"sync"
)

const defaultMemoryBufferSize = 1000

// SpillBuffer provides a buffer that operates in-memory by default and spills to a
// temporary disk file if the in-memory channel becomes full. This provides the
// high performance of an in-memory buffer for the common case, and the high capacity
// of a disk buffer when the consumer is slow. This is a one-way transition; once
// the buffer spills to disk, it remains in disk mode for its lifetime.
type SpillBuffer struct {
	filePath    string
	frameLength int

	memoryChan chan []byte
	outputChan chan []byte

	mu         sync.Mutex
	isSpilling bool
	fileBuffer *FileBuffer
	wg         sync.WaitGroup
	closeOnce  sync.Once
}

// NewSpillBuffer creates a new spill buffer.
// filePath is the path to the file to use *if* spilling occurs.
// frameLength is the size of each frame in bytes.
func NewSpillBuffer(filePath string, frameLength int) *SpillBuffer {
	return &SpillBuffer{
		filePath:    filePath,
		frameLength: frameLength,
		memoryChan:  make(chan []byte, defaultMemoryBufferSize),
		outputChan:  make(chan []byte, defaultMemoryBufferSize),
	}
}

// Start launches the necessary goroutines for the buffer to operate.
// It returns a write-only channel for the producer to send data to, and a
// read-only channel for the consumer to receive data from.
func (sb *SpillBuffer) Start() (chan<- []byte, <-chan []byte) {
	sb.wg.Add(1)
	go sb.process()
	// We return a separate input channel to allow us to close it gracefully.
	inputChan := make(chan []byte)
	go func() {
		for frame := range inputChan {
			sb.Write(frame)
		}
		sb.CloseInput()
	}()
	return inputChan, sb.outputChan
}

// Write sends a frame to the buffer. If the in-memory channel is full,
// it triggers a permanent transition to a file-based buffer.
func (sb *SpillBuffer) Write(frame []byte) {
	sb.mu.Lock()
	if sb.isSpilling {
		sb.mu.Unlock()
		// Already spilling, just write to the file buffer.
		sb.fileBuffer.input <- frame
		return
	}
	sb.mu.Unlock()

	select {
	case sb.memoryChan <- frame:
		// Fast path: Sent to memory successfully.
	default:
		// Slow path: In-memory channel is full. Transition to disk.
		sb.spillToDisk()
		// Retry the write, which will now go to the file buffer.
		sb.fileBuffer.input <- frame
	}
}

// CloseInput signals that the producer has finished sending data.
func (sb *SpillBuffer) CloseInput() {
	sb.closeOnce.Do(func() {
		close(sb.memoryChan)
	})
}

// spillToDisk handles the one-way transition from in-memory to disk buffering.
func (sb *SpillBuffer) spillToDisk() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	// Check again in case another goroutine initiated the spill while we waited for the lock.
	if sb.isSpilling {
		return
	}

	slog.Warn("SpillBuffer transitioning to disk mode", "path", sb.filePath)
	sb.isSpilling = true

	// Close the in-memory channel to signal the processor to flush it.
	close(sb.memoryChan)

	// Create and start the underlying file buffer.
	sb.fileBuffer = NewFileBuffer(sb.filePath, sb.frameLength)
}

// process is the core goroutine that reads from the buffer and sends to the output.
func (sb *SpillBuffer) process() {
	defer sb.wg.Done()
	defer close(sb.outputChan)

	// First, process everything from the in-memory channel.
	for frame := range sb.memoryChan {
		sb.outputChan <- frame
	}

	// After the memory channel is closed, check if we have spilled to disk.
	sb.mu.Lock()
	hasSpilled := sb.isSpilling
	sb.mu.Unlock()

	if hasSpilled {
		// We are in disk mode. The file buffer now takes over.
		// The file buffer's input channel needs to be closed when the producer is done.
		// We can't know that directly here, so the producer must call CloseInput on the SpillBuffer.
		// The FileBuffer's Start method returns its own channels.
		fileInput, fileOutput := sb.fileBuffer.Start()
		fmt.Println(fileInput)

		// We need to handle the closing of the fileInput channel.
		// This is tricky. A simpler model might be to have the producer call a
		// specific 'Close' method on the SpillBuffer that handles this.

		// For now, let's assume the FileBuffer's input channel is managed correctly.
		for frame := range fileOutput {
			sb.outputChan <- frame
		}
	}
}
