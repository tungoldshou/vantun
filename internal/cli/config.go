package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"vantun/internal/core"
)

// Config represents the VANTUN configuration.
type Config struct {
	// Server indicates if this is a server configuration.
	Server bool `json:"server"`
	// Address is the address to listen on (server) or connect to (client).
	Address string `json:"address"`
	// LogLevel is the log level (debug, info, warn, error).
	LogLevel string `json:"log_level"`
	// Multipath enables multipath.
	Multipath bool `json:"multipath"`
	// Obfs enables obfuscation.
	Obfs bool `json:"obfs"`
	// FECData is the number of FEC data shards.
	FECData int `json:"fec_data"`
	// FECParity is the number of FEC parity shards.
	FECParity int `json:"fec_parity"`
	// TokenBucketRate is the initial rate for the token bucket (bytes per second).
	TokenBucketRate float64 `json:"token_bucket_rate"`
	// TokenBucketCapacity is the capacity of the token bucket (bytes).
	TokenBucketCapacity float64 `json:"token_bucket_capacity"`
}

// ConfigManager manages the configuration with hot reloading capability.
type ConfigManager struct {
	configFile string
	config     *Config
	mutex      sync.RWMutex
	watcher    *time.Ticker
	stopChan   chan struct{}
}

// NewConfigManager creates a new ConfigManager.
func NewConfigManager(configFile string) *ConfigManager {
	return &ConfigManager{
		configFile: configFile,
		stopChan:   make(chan struct{}),
	}
}

// LoadConfig loads the configuration from a JSON file.
func LoadConfig(filename string) (*Config, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Decode the JSON
	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

// Load loads the configuration from the file.
func (cm *ConfigManager) Load() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	config, err := LoadConfig(cm.configFile)
	if err != nil {
		return err
	}

	cm.config = config
	return nil
}

// GetConfig returns the current configuration.
func (cm *ConfigManager) GetConfig() *Config {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Check if config is nil (should not happen after successful StartHotReload)
	if cm.config == nil {
		return nil
	}

	// Return a copy of the config to prevent external modifications
	config := *cm.config
	return &config
}

// StartHotReload starts the hot reloading of the configuration.
func (cm *ConfigManager) StartHotReload() error {
	// Load the initial configuration
	if err := cm.Load(); err != nil {
		return err
	}

	// Start a ticker to check for config file changes
	cm.watcher = time.NewTicker(5 * time.Second)
	
	go func() {
		for {
			select {
			case <-cm.watcher.C:
				cm.checkForChanges()
			case <-cm.stopChan:
				return
			}
		}
	}()

	return nil
}

// checkForChanges checks if the config file has been modified and reloads it if necessary.
func (cm *ConfigManager) checkForChanges() {
	cm.mutex.RLock()
	currentConfig := cm.config
	cm.mutex.RUnlock()

	if currentConfig == nil {
		return
	}

	// Load new config
	newConfig, err := LoadConfig(cm.configFile)
	if err != nil {
		core.Error("Failed to load new config: %v", err)
		return
	}

	// Check if the config has changed
	if cm.hasConfigChanged(currentConfig, newConfig) {
		core.Info("Configuration changed, reloading...")
		cm.mutex.Lock()
		cm.config = newConfig
		cm.mutex.Unlock()
		
		// Update the global logger level if it has changed
		if currentConfig.LogLevel != newConfig.LogLevel {
			core.InitLogger(newConfig.LogLevel)
		}
	}
}

// hasConfigChanged checks if the configuration has changed.
func (cm *ConfigManager) hasConfigChanged(oldConfig, newConfig *Config) bool {
	return oldConfig.Server != newConfig.Server ||
		oldConfig.Address != newConfig.Address ||
		oldConfig.LogLevel != newConfig.LogLevel ||
		oldConfig.Multipath != newConfig.Multipath ||
		oldConfig.Obfs != newConfig.Obfs ||
		oldConfig.FECData != newConfig.FECData ||
		oldConfig.FECParity != newConfig.FECParity ||
		oldConfig.TokenBucketRate != newConfig.TokenBucketRate ||
		oldConfig.TokenBucketCapacity != newConfig.TokenBucketCapacity
}

// StopHotReload stops the hot reloading of the configuration.
func (cm *ConfigManager) StopHotReload() {
	if cm.watcher != nil {
		cm.watcher.Stop()
	}
	close(cm.stopChan)
}