// Package acquisition handles the real-time data acquisition from hardware.
package buffer

import (
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

const (
	// bufferSize determines the channel buffer size for input and output.
	bufferSize = 1000
	// writeBatchSizeFrames is the number of frames to batch together for a single file write.
	writeBatchSizeFrames = 100
)

// FileBuffer provides a temporary file-based buffer to decouple a fast producer
// from a potentially slower consumer. It writes all incoming data to a file and
// reads it back out, using the file system as a large buffer.
type FileBuffer struct {
	filePath    string
	frameLength int

	input  chan []byte
	output chan []byte

	wg        sync.WaitGroup
	closeOnce sync.Once
}

// NewFileBuffer creates and initializes a new FileBuffer.
// filePath is the path to the temporary file to use.
// frameLength is the exact size in bytes of each data frame.
func NewFileBuffer(filePath string, frameLength int) *FileBuffer {
	return &FileBuffer{
		filePath:    filePath,
		frameLength: frameLength,
		input:       make(chan []byte, bufferSize),
		output:      make(chan []byte, bufferSize),
	}
}

// Start launches the internal writer and reader goroutines.
// It returns two channels for the caller to use:
// 1. A write-only channel to push data into the buffer.
// 2. A read-only channel to get data out of the buffer.
func (fb *FileBuffer) Start() (chan<- []byte, <-chan []byte) {
	fb.wg.Add(2)
	go fb.runWriter()
	go fb.runReader()

	return fb.input, fb.output
}

// Close gracefully shuts down the FileBuffer, ensuring all data is flushed.
// It closes the input channel and waits for the writer and reader goroutines to complete.
func (fb *FileBuffer) Close() {
	fb.closeOnce.Do(func() {
		close(fb.input)
		fb.wg.Wait()
	})
}

// runWriter is the internal goroutine that reads from the input channel
// and writes the data to the backing file.
func (fb *FileBuffer) runWriter() {
	defer fb.wg.Done()

	file, err := os.OpenFile(fb.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("FileBuffer writer failed to open file", "path", fb.filePath, "error", err)
		return
	}
	defer file.Close()

	writeBatchSize := fb.frameLength * writeBatchSizeFrames
	writeBuffer := make([]byte, 0, writeBatchSize)

	for frame := range fb.input {
		writeBuffer = append(writeBuffer, frame...)
		if len(writeBuffer) >= writeBatchSize {
			if _, err := file.Write(writeBuffer); err != nil {
				slog.Error("FileBuffer writer failed to write", "path", fb.filePath, "error", err)
			}
			// Reset buffer
			writeBuffer = writeBuffer[:0]
		}
	}

	// After the input channel is closed, write any remaining data.
	if len(writeBuffer) > 0 {
		if _, err := file.Write(writeBuffer); err != nil {
			slog.Error("FileBuffer writer failed to write final batch", "path", fb.filePath, "error", err)
		}
	}
}

// runReader is the internal goroutine that reads from the backing file
// and sends the data to the output channel.
func (fb *FileBuffer) runReader() {
	defer fb.wg.Done()
	defer close(fb.output)

	// The reader must wait for the writer to create the file.
	// We poll briefly, as this should happen almost instantly.
	var file *os.File
	var err error
	for i := 0; i < 10; i++ {
		file, err = os.OpenFile(fb.filePath, os.O_RDONLY, 0)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		slog.Error("FileBuffer reader failed to open file", "path", fb.filePath, "error", err)
		return
	}
	defer file.Close()

	// Create a buffer for exactly one frame.
	frameBuffer := make([]byte, fb.frameLength)

	for {
		// io.ReadFull guarantees that the frameBuffer is filled completely.
		_, err := io.ReadFull(file, frameBuffer)

		if err != nil {
			// If we hit the end of the file, it's a clean exit.
			// ErrUnexpectedEOF is also expected if the file size is not a multiple of frameLength.
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			// For any other error, log it and stop.
			slog.Error("FileBuffer reader error", "path", fb.filePath, "error", err)
			break
		}

		// We must send a copy to the output channel, as frameBuffer will be reused.
		frameCopy := make([]byte, fb.frameLength)
		copy(frameCopy, frameBuffer)

		fb.output <- frameCopy
	}
}
