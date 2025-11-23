package microwave

import (
	"csrspServer/session"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// PRFJob represents a single PRF that needs to be processed.
type PRFJob struct {
	Polarization      string
	TimingState       string
	PRFIndexCounter   uint64
	ISamples          []float64
	QSamples          []float64
	SamplingFrequency float64
}

// Processor handles the concurrent processing of microwave data.
type Processor struct {
	jobChan chan PRFJob
	wg      sync.WaitGroup
	s       *session.Service
}

// NewProcessor creates a new microwave processor and starts its worker pool.
func NewProcessor(s *session.Service) *Processor {
	numWorkers := runtime.NumCPU()
	p := &Processor{
		jobChan: make(chan PRFJob, numWorkers*2),
		s:       s,
	}

	p.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go p.worker()
	}

	return p
}

// worker is a goroutine that receives and processes PRF jobs.
func (p *Processor) worker() {
	defer p.wg.Done()
	for job := range p.jobChan {
		// Define a unique key for this PRF instance
		key := fmt.Sprintf("%s;%s;%d", job.Polarization, job.TimingState, job.PRFIndexCounter)

		// Calculate I/Q statistics
		iStats := iqStatisticsCalculator(job.ISamples)
		qStats := iqStatisticsCalculator(job.QSamples)

		var chirp session.ChirpParams
		// Perform chirp analysis if not in imaging mode
		if strings.ToLower(job.TimingState) != "imaging" {
			chirp = generateSARChirpParameters(job.ISamples, job.QSamples, job.SamplingFrequency)
		}

		// Initialize and update the status in the session store
		p.s.Store.Microwave.InitPRFStatus(key)
		p.s.Store.Microwave.UpdatePRFStatus(key, iStats, qStats, chirp)
	}
}

// ProcessPRFBundle is the main entry point for processing a new set of PRF data.
func (p *Processor) ProcessPRFBundle(params SARProcessingParams, prfDataArray [][]byte) {
	// 1. Extract I/Q Samples
	videoDataLen := uint64(len(prfDataArray[0])) - prfHeaderSize - prfFooterSize
	getBitsProviders := make([]func(uint64, uint64) uint64, len(prfDataArray))
	for i := range prfDataArray {
		temp := prfDataArray[i][prfHeaderSize:(uint64(len(prfDataArray[i])) - prfFooterSize)]
		getBitsProviders[i] = getDataProviderForNonContinuousData(temp)
	}
	iSamples, qSamples := iqSampleProvider(params.SarModeValue, params.BaqValue, videoDataLen, getBitsProviders)

	// 2. Dispatch job to a worker
	p.jobChan <- PRFJob{
		Polarization:      params.Polarization,
		TimingState:       params.TimingState,
		PRFIndexCounter:   params.PrfIndexCounter,
		ISamples:          iSamples,
		QSamples:          qSamples,
		SamplingFrequency: params.SamplingFrequency,
	}
}

// Shutdown gracefully stops the processor's worker pool.
func (p *Processor) Shutdown() {
	close(p.jobChan)
	p.wg.Wait()
}

// SARProcessingParams is a direct copy from the old ProcessIDProcessor.go
type SARProcessingParams struct {
	PrfIndexCounter       uint64
	Polarization          string
	TimingState           string
	BaqValue              uint64
	SamplingFrequency     float64
	SarModeValue          string
	SarModeValueFromDB    string
	TimingStatesToExclude string
}
