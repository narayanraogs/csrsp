package global

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration holds all fixed parameters for the application, loaded from a JSON file.
type Configuration struct {
	MainParameters        MainParams        `json:"mainParameters"`
	AcquisitionParameters AcquisitionParams `json:"acquisitionParameters"`
	DasLockStatus         DasLockStatus     `json:"dasLockStatus"`
	DasPorts              DasPorts          `json:"dasPorts"`
	LBTParameters         LBTParams         `json:"lbtParameters"`
	DRParameters          HDRParams         `json:"globalHDRParameters"`
	NetworkInfo           NetworkInfo       `json:"networkInfo"`
}

// MainParams holds core application settings like paths and database credentials.
type MainParams struct {
	SatName                string   `json:"satName"`
	RootPath               string   `json:"rootPath"`
	AcqPath                string   `json:"acqPath"`
	TempPath               string   `json:"tempPath"`
	ArchPath               string   `json:"archPath"`
	DBServerIP             string   `json:"dbServerIP"`
	OutputLevel            string   `json:"outputLevel"`
	DBName                 string   `json:"dbName"`
	DBUser                 string   `json:"dbUser"`
	DBPassword             string   `json:"dbPassword"`
	SystemName             string   `json:"systemName"`
	LogFileDirectory       string   `json:"logFileDirectory"`
	AssetPath              string   `json:"assetPath"`
	WebPath                string   `json:"webPath"`
	DevOpsPath             string   `json:"devOpsPath"`
	ProcessingSequencePath string   `json:"processingSequencePath"`
	WhiteListedIPs         []string `json:"whiteListedIPs"`
	ResultNamesPath        string   `json:"resultNamesPath"`
}

// NetworkInfo contains information about other connected systems.
type NetworkInfo struct {
	PCCList []PCCInfo `json:"pccList"`
}

// PCCInfo holds connection details for a single PCC.
type PCCInfo struct {
	SystemName string `json:"systemName"`
	IPAddress  string `json:"ipAddress"`
	ArchPath   string `json:"archPath"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
}

// AcquisitionParams holds settings related to data acquisition.
type AcquisitionParams struct {
	NumberOfFramesInBlock int64 `json:"numberOfFramesInBlock"`
}

// DasLockStatus holds parameters for checking DAS lock status.
type DasLockStatus struct {
	ComponentCode int `json:"componentCode"`
	LockOffset    int `json:"lockOffset"`
	LockValue     int `json:"lockValue"`
	AlarmOffset   int `json:"alarmOffset"`
	AlarmValue    int `json:"alarmValue"`
}

// DasPorts specifies the network ports for various DAS services.
type DasPorts struct {
	Monitor int         `json:"monitor"`
	Control int         `json:"control"`
	Acq     DasAcqPorts `json:"acq"`
}

// DasAcqPorts specifies network ports for acquisition components.
type DasAcqPorts struct {
	HDR  int `json:"hdr"`
	TTCP int `json:"ttcp"`
	PDPS int `json:"pdps"`
	CDP  int `json:"cdp"`
}

// LBTParams holds parameters for LBT (Link Budget Tool) processing.
type LBTParams struct {
	StartWord    int    `json:"startWord"`
	CRCwords     int    `json:"crcWords"`
	LbtLocalPath string `json:"lbtLocalPath"`
}

// HDRParams holds parameters for processing the global header.
type HDRParams struct {
	CompCode            string `json:"compCode"`
	TotalLengthInWords  int    `json:"totalLengthInWords"`
	ConfigLengthInWords int    `json:"configLengthInWords"`
	BytesPerWord        int    `json:"bytesPerWord"`
	HeaderLengthInWords int    `json:"headerLengthInWords"`
	ExcludedWords       []int  `json:"excludedWords"`
	RetryLimit          int    `json:"retryLimit"`
	RetryGapInMillis    int    `json:"retryGapInMillis"`
	ComponentDelay      int    `json:"componentDelay"`
}

// DeveloperOptions holds configuration settings that are primarily for development and debugging.
type DeveloperOptions struct {
	AutomaticArchival  bool   `json:"automaticArchival"`
	LogLevel           string `json:"logLevel"`
	ParallelProcessing bool   `json:"parallelProcessing"`
	EncryptionMode     string `json:"encryptionMode"`
	EndProcessID       int32  `json:"endProcessID"`
	MaxThreads         int32  `json:"maxThreads"`
	BufferLength       int32  `json:"bufferLength"`
}

func (do *DeveloperOptions) Save(path string) error {
	file, err := json.MarshalIndent(do, "", "  ") // Using 2-space indent for readability
	if err != nil {
		return fmt.Errorf("failed to marshal developer options to JSON: %w", err)
	}

	if err := os.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", path, err)
	}

	return nil
}

// LoadDeveloperOptions reads the Developer options from the specified path.
// If the file doesn't exist, it saves defaults to that path and loads them.
func LoadDeveloperOptions(path string) (*DeveloperOptions, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Define defaults
		defaults := &DeveloperOptions{
			AutomaticArchival:  false,
			LogLevel:           "INFO",
			ParallelProcessing: false,
			EncryptionMode:     "NONE",
			EndProcessID:       -1,
			MaxThreads:         100000,
			BufferLength:       100000,
		}

		// Save defaults to the file
		if err := defaults.Save(path); err != nil {
			return nil, fmt.Errorf("failed to save default developer options: %w", err)
		}

		return defaults, nil
	}

	// Read existing file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read developer options file: %w", err)
	}

	var options DeveloperOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("failed to unmarshal developer options: %w", err)
	}

	return &options, nil
}
