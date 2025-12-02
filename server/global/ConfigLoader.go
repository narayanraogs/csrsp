package global

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfig reads the configuration from the specified path and returns a Config struct.
func LoadConfig(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Configuration
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Map Configuration to Config
	cfg := Config{
		// Paths
		RootPath:               config.MainParameters.RootPath,
		AcqPath:                config.MainParameters.AcqPath,
		TempPath:               config.MainParameters.TempPath,
		ArchPath:               config.MainParameters.ArchPath,
		LogFileDirectory:       config.MainParameters.LogFileDirectory,
		AssetPath:              config.MainParameters.AssetPath,
		WebPath:                config.MainParameters.WebPath,
		DevOpsPath:             config.MainParameters.DevOpsPath,
		ProcessingSequencePath: config.MainParameters.ProcessingSequencePath,
		ResultNamesPath:        config.MainParameters.ResultNamesPath,

		// Database Credentials
		DBServerIP: config.MainParameters.DBServerIP,
		DBName:     config.MainParameters.DBName,
		DBUser:     config.MainParameters.DBUser,
		DBPassword: config.MainParameters.DBPassword,

		// Core Parameters
		SystemName:            config.MainParameters.SystemName,
		SatName:               config.MainParameters.SatName,
		NumberOfFramesInBlock: config.AcquisitionParameters.NumberOfFramesInBlock,

		// Network & Ports
		WhiteListedIPs: config.MainParameters.WhiteListedIPs,
		PCCList:        config.NetworkInfo.PCCList,
		DasPorts:       config.DasPorts,

		// Hardware & Processing Parameters
		DasLockStatus: config.DasLockStatus,
	}

	return cfg, nil
}
