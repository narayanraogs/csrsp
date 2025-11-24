package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// This is an example of an external processor compatible with the GenericProcessor stage.
// It acts as a "pass-through" or "echo" processor, but you can modify the `process` function
// to perform actual work (e.g., decryption, image processing).

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <read_socket_path> <write_socket_path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Note: <read_socket_path> is where this program WRITES output.\n")
		fmt.Fprintf(os.Stderr, "      <write_socket_path> is where this program READS input.\n")
		os.Exit(1)
	}

	// The naming in GenericProcessor is from the Go server's perspective.
	// Arg 1: The socket the Server listens on for READING. We must CONNECT and WRITE to it.
	// Arg 2: The socket the Server listens on for WRITING. We must CONNECT and READ from it.
	outputSocketPath := os.Args[1]
	inputSocketPath := os.Args[2]

	fmt.Printf("ExternalProcessor: Starting...\n")
	fmt.Printf("ExternalProcessor: Connecting to input (Server Write): %s\n", inputSocketPath)
	fmt.Printf("ExternalProcessor: Connecting to output (Server Read): %s\n", outputSocketPath)

	// Retry logic is helpful because the server might take a millisecond to start listening
	inputConn := connectWithRetry(inputSocketPath)
	defer inputConn.Close()

	outputConn := connectWithRetry(outputSocketPath)
	defer outputConn.Close()

	fmt.Println("ExternalProcessor: Connected. Starting processing loop.")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("ExternalProcessor: Shutting down...")
		inputConn.Close()
		outputConn.Close()
		os.Exit(0)
	}()

	// Buffer for reading headers
	headerBuf := make([]byte, 16)

	for {
		// 1. Read Header (Channel Index + Payload Length)
		_, err := io.ReadFull(inputConn, headerBuf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("ExternalProcessor: Input closed, exiting.")
				return
			}
			fmt.Fprintf(os.Stderr, "ExternalProcessor: Error reading header: %v\n", err)
			return
		}

		channelIndex := binary.BigEndian.Uint64(headerBuf[:8])
		payloadLen := binary.BigEndian.Uint64(headerBuf[8:])

		// 2. Read Payload
		payload := make([]byte, payloadLen)
		_, err = io.ReadFull(inputConn, payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ExternalProcessor: Error reading payload: %v\n", err)
			return
		}

		// 3. Process Data
		// This is where your custom logic goes (Decryption, etc.)
		processedPayload := process(payload)

		// 4. Write Header back
		// We keep the same channel index to route it back to the correct pipeline output
		// Update length if the payload size changed
		binary.BigEndian.PutUint64(headerBuf[:8], channelIndex)
		binary.BigEndian.PutUint64(headerBuf[8:], uint64(len(processedPayload)))

		_, err = outputConn.Write(headerBuf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ExternalProcessor: Error writing header: %v\n", err)
			return
		}

		// 5. Write Payload back
		_, err = outputConn.Write(processedPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ExternalProcessor: Error writing payload: %v\n", err)
			return
		}
	}
}

func connectWithRetry(path string) net.Conn {
	for i := 0; i < 10; i++ {
		conn, err := net.Dial("unix", path)
		if err == nil {
			return conn
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Fprintf(os.Stderr, "ExternalProcessor: Failed to connect to %s after retries\n", path)
	os.Exit(1)
	return nil
}

// process implements your custom logic.
// Currently, it just returns the data as-is (Echo).
func process(input []byte) []byte {
	// Example: Invert bytes (just to show some work)
	/*
		output := make([]byte, len(input))
		for i, b := range input {
			output[i] = ^b
		}
		return output
	*/

	// For now, just pass through
	return input
}
