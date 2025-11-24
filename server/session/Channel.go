// Package session manages the state and data flow for individual client sessions.
package session

import (
	"log/slog"
	"os"
	"sync"
)

const channelLimit = 100000

// Channel provides a one-to-many broadcast mechanism for data frames.
// It allows a single provider to send data to multiple consumers.
type Channel struct {
	input           chan []byte
	output          []chan []byte
	fd              *os.File
	filePath        string
	providerPresent bool
	write           bool
	frameNos        uint64
	mutex           sync.Mutex
}

// newChannel creates and initializes a new Channel.
func newChannel() *Channel {
	return &Channel{
		output: make([]chan []byte, 0),
	}
}

// RegisterAsProvider registers a goroutine as the single source of data for this Channel.
// It returns the input channel to which the provider should send data.
func (c *Channel) RegisterAsProvider() chan []byte {
	if c.providerPresent {
		slog.Warn("Provider already present for channel, possibly due to a configuration issue. Returning a dummy channel.", "path", c.filePath)
		return make(chan []byte, channelLimit) // Return a dummy channel to prevent blocking
	}

	c.providerPresent = true
	c.input = make(chan []byte, channelLimit)

	// Launch the processor that broadcasts data from the input to all outputs.
	go c.process()

	return c.input
}

// RegisterAsConsumer registers a goroutine as a receiver of data from this Channel.
// It returns the output channel from which the consumer can read data.
func (c *Channel) RegisterAsConsumer() chan []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tempChan := make(chan []byte, channelLimit)
	c.output = append(c.output, tempChan)
	return tempChan
}

// process is the core loop that reads from the single input channel
// and broadcasts the data to all registered output channels.
func (c *Channel) process() {
	defer c.closeAllOutputs()

	// Buffer for file writing
	var writeBuffer []byte
	const writeBatchSize = 100
	framesToBatch := 0

	if c.fd != nil {
		defer c.fd.Close()
	}

	for data := range c.input {
		c.frameNos++

		if c.write {
			framesToBatch++
			writeBuffer = append(writeBuffer, data...)
			if framesToBatch >= writeBatchSize {
				_, err := c.fd.Write(writeBuffer)
				if err != nil {
					slog.Error("Cannot write to file from channel", "path", c.filePath, "error", err)
				}
				writeBuffer = nil // Clear the buffer
				framesToBatch = 0
			}
		}

		// Broadcast to all consumers
		for _, outputChan := range c.output {
			// Create a copy to ensure each consumer gets its own instance of the data.
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)
			outputChan <- dataCopy
		}
	}

	// Write any remaining data in the buffer before closing.
	if len(writeBuffer) > 0 && c.write {
		_, err := c.fd.Write(writeBuffer)
		if err != nil {
			slog.Error("Cannot write remaining buffer to file", "path", c.filePath, "error", err)
		}
	}
}

// closeAllOutputs closes all registered consumer channels.
func (c *Channel) closeAllOutputs() {
	for _, outputChan := range c.output {
		close(outputChan)
	}
}

// SetFileForWriting configures the channel to write all data it receives to a file.
func (c *Channel) SetFileForWriting(path string) {
	var err error
	c.filePath = path
	c.fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		slog.Error("Unable to open file for channel writing", "path", path, "error", err)
		c.fd = nil
		return
	}
	c.write = true
}

// FilePath returns the path of the file this channel is writing to.
func (c *Channel) FilePath() string {
	return c.filePath
}
