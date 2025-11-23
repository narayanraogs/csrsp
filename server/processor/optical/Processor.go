// Package optical handles the processing of video/image data frames.
package optical

import (
	"csrspServer/db"
	"csrspServer/db/sqlc"
	"csrspServer/session"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"sync"
)

const (
	numWorkers = 10
)

// PixelProcessor holds the state and logic for processing a single pixel over time.
type PixelProcessor struct {
	// Configuration
	pixelNo  int
	limits   sqlc.Videolimit
	isFaulty bool

	// State
	noOfSamples int
	mean        float64
	sdSqr       float64
	sd          float64
	max         float64
	min         float64
}

// NewPixelProcessor creates a new processor for a single pixel.
func NewPixelProcessor(pixelNo int, limits sqlc.Videolimit, isFaulty bool) *PixelProcessor {
	return &PixelProcessor{
		pixelNo:  pixelNo,
		limits:   limits,
		isFaulty: isFaulty,
	}
}

// Update processes a new value for the pixel and returns the updated statistics.
func (p *PixelProcessor) Update(value int) (mean, sd, diff float64) {
	fValue := float64(value)
	if p.noOfSamples == 0 {
		p.mean = fValue
		p.sdSqr = 0
		p.sd = 0
		p.noOfSamples = 1
		p.max = fValue
		p.min = fValue
	} else {
		prevMean := p.mean
		p.mean = p.mean + ((fValue - p.mean) / float64(p.noOfSamples))
		p.sdSqr = p.sdSqr + (fValue-prevMean)*(fValue-p.mean)
		if p.noOfSamples > 1 {
			p.sd = math.Sqrt(p.sdSqr / (float64(p.noOfSamples) - 1))
		}

		if fValue > p.max {
			p.max = fValue
		}
		if fValue < p.min {
			p.min = fValue
		}

		p.noOfSamples++
	}

	// Check for errors and update status flags on the processor itself.
	if p.isFaulty {
		// Logic to handle faulty status if needed
	}
	if p.mean > p.limits.Maxmean {
		// Logic to handle mean error if needed
	}
	if p.sd > p.limits.Maxsd {
		// Logic to handle sd error if needed
	}

	return p.mean, p.sd, p.max - p.min
}

// getPixelValues extracts pixel values from a line of data. This is an optimized
// version that pre-allocates the slice to avoid repeated appends.
func getPixelValues(lineData []byte, bitsPerPixel int) []int {
	bytesPerPixel := bitsPerPixel / 8
	if bytesPerPixel == 0 || len(lineData)%bytesPerPixel != 0 {
		return nil // Invalid data
	}
	numPixels := len(lineData) / bytesPerPixel
	pixels := make([]int, numPixels)

	switch bytesPerPixel {
	case 1:
		for i := 0; i < numPixels; i++ {
			pixels[i] = int(lineData[i])
		}
	case 2:
		for i := 0; i < numPixels; i++ {
			pixels[i] = int(lineData[i*2])<<8 | int(lineData[i*2+1])
		}
	case 3:
		for i := 0; i < numPixels; i++ {
			pixels[i] = int(lineData[i*3])<<16 | int(lineData[i*3+1])<<8 | int(lineData[i*3+2])
		}
	case 4:
		for i := 0; i < numPixels; i++ {
			pixels[i] = int(lineData[i*4])<<24 | int(lineData[i*4+1])<<16 | int(lineData[i*4+2])<<8 | int(lineData[i*4+3])
		}
	}
	return pixels
}

// FrameData holds the raw data and metadata for a single video frame.
type FrameData struct {
	FrameType       string
	FrameIdentifier string
	Payload         []byte
	// Add other metadata like width, height, bitsPerPixel etc.
	Width        int
	Height       int
	BitsPerPixel int
}

// RegionJob defines a unit of work for a worker goroutine, representing a
// horizontal strip of the image to be processed.
type RegionJob struct {
	Frame *FrameData
	// The start and end line numbers for this region.
	StartLine int
	EndLine   int
}

type Processor struct {
	jobChan    chan RegionJob
	numWorkers int
	wg         sync.WaitGroup

	// state holds the persistent PixelProcessor for every pixel.
	state      map[string]*PixelProcessor
	stateMutex sync.RWMutex
}

// NewProcessor creates a new optical processor and starts its worker pool.
func NewProcessor(s *session.Service) *Processor {
	p := &Processor{
		jobChan:    make(chan RegionJob, numWorkers*2), // Buffered channel
		numWorkers: numWorkers,
		state:      make(map[string]*PixelProcessor),
	}

	// Start the worker pool.
	p.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go p.worker(i+1, s)
	}

	slog.Info("Optical processor started", "numWorkers", numWorkers)
	return p
}

// worker is a long-lived goroutine that processes regions of frames.

func (p *Processor) worker(id int, s *session.Service) {

	defer p.wg.Done()

	for job := range p.jobChan {

		// The actual processing logic is now integrated here.

		p.processRegion(&job, s)

	}

}

// processRegion handles the processing of a specific region of a frame.
func (p *Processor) processRegion(job *RegionJob, s *session.Service) {
	frame := job.Frame
	bytesPerPixel := frame.BitsPerPixel / 8
	bytesPerLine := frame.Width * bytesPerPixel

	// Initialize the processor's state on the first job of a frame.
	p.stateMutex.Lock()
	if len(p.state) == 0 {
		fid, _ := db.GetFrameIDForFrameIdentifier(frame.FrameType, frame.FrameIdentifier)
		// This is the first frame, so we create the persistent PixelProcessors.
		limits, _ := db.GetVideoLimitsForFrameID(int(fid))
		faultyPixelList, _ := db.GetFaultyPixels(int(fid))
		faultyPixels := make(map[int]bool)
		faultyPixelString := strings.Split(faultyPixelList, ",")
		for _, px := range faultyPixelString {
			p, _ := strconv.Atoi(strings.TrimSpace(px))
			faultyPixels[p] = true
		}

		numPixels := frame.Width * frame.Height
		for i := 0; i < numPixels; i++ {
			key := fmt.Sprintf("%s;%s;%d", frame.FrameType, frame.FrameIdentifier, i)
			_, isFaulty := faultyPixels[i]
			p.state[key] = NewPixelProcessor(i, limits, isFaulty)
		}
		// Also initialize the session store for the client.
		s.Store.Video.InitStatuses(frame.FrameType, frame.FrameIdentifier, numPixels)
	}
	p.stateMutex.Unlock()

	// Process the assigned region sequentially.
	for y := job.StartLine; y < job.EndLine; y++ {
		lineStart := y * bytesPerLine
		lineEnd := lineStart + bytesPerLine
		if lineEnd > len(frame.Payload) {
			continue // Avoid panic on malformed data
		}
		lineData := frame.Payload[lineStart:lineEnd]

		pixels := getPixelValues(lineData, frame.BitsPerPixel)
		for x, pixelValue := range pixels {
			pixelNo := y*frame.Width + x
			key := fmt.Sprintf("%s;%s;%d", frame.FrameType, frame.FrameIdentifier, pixelNo)

			// Get the persistent processor for this pixel.
			p.stateMutex.RLock()
			pixelProc, ok := p.state[key]
			p.stateMutex.RUnlock()

			if ok {
				// Update its state and get the new stats.
				mean, sd, diff := pixelProc.Update(pixelValue)
				// Write the new stats to the session store for the client.
				s.Store.Video.UpdatePixelStats(key, mean, sd, diff)
			}
		}
	}
}

// ProcessFrame slices a frame into regions and distributes them to the workers.
func (p *Processor) ProcessFrame(frame *FrameData) {
	// Simple division of the frame into horizontal strips.
	linesPerWorker := frame.Height / p.numWorkers
	if linesPerWorker == 0 {
		linesPerWorker = 1
	}

	for i := 0; i < p.numWorkers; i++ {
		startLine := i * linesPerWorker
		endLine := startLine + linesPerWorker

		// Ensure the last worker gets all remaining lines.
		if i == p.numWorkers-1 {
			endLine = frame.Height
		}

		if startLine >= frame.Height {
			continue
		}

		p.jobChan <- RegionJob{
			Frame:     frame,
			StartLine: startLine,
			EndLine:   endLine,
		}
	}
}

// Shutdown gracefully stops the processor and its workers.
func (p *Processor) Shutdown() {
	close(p.jobChan)
	p.wg.Wait()
	slog.Info("Optical processor shut down gracefully.")
}
