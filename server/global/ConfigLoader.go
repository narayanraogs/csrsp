package global

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	// 1. Validate RootPath is absolute
	rootPath := config.MainParameters.RootPath
	if !filepath.IsAbs(rootPath) {
		return Config{}, fmt.Errorf("RootPath must be absolute: %s", rootPath)
	}

	// Helper to resolve paths
	resolvePath := func(target string) string {
		if filepath.IsAbs(target) {
			return target
		}
		return filepath.Join(rootPath, target)
	}

	// Map Configuration to Config
	cfg := Config{
		// Paths
		RootPath:               rootPath,
		AcqPath:                resolvePath(config.MainParameters.AcqPath),
		TempPath:               resolvePath(config.MainParameters.TempPath),
		ArchPath:               resolvePath(config.MainParameters.ArchPath),
		LogFileDirectory:       resolvePath(config.MainParameters.LogFileDirectory),
		AssetPath:              resolvePath(config.MainParameters.AssetPath),
		WebPath:                resolvePath(config.MainParameters.WebPath),
		DevOpsPath:             resolvePath(config.MainParameters.DevOpsPath),
		ProcessingSequencePath: resolvePath(config.MainParameters.ProcessingSequencePath),
		ResultNamesPath:        resolvePath(config.MainParameters.ResultNamesPath),

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
