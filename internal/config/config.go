package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/martinlehoux/kagapass/internal/types"
)

// Manager handles configuration loading and saving
type Manager struct {
	configDir    string
	configPath   string
	databasePath string
}

// New creates a new configuration manager
func New() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "kagapass")
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return nil, err
	}

	return &Manager{
		configDir:    configDir,
		configPath:   filepath.Join(configDir, "config.json"),
		databasePath: filepath.Join(configDir, "databases.json"),
	}, nil
}

// LoadConfig loads the application configuration
func (m *Manager) LoadConfig() (types.Config, error) {
	config := types.DefaultConfig()

	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Create default config file
		return config, m.SaveConfig(config)
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves the application configuration
func (m *Manager) SaveConfig(config types.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0o600)
}

// LoadDatabaseList loads the list of configured databases
func (m *Manager) LoadDatabaseList() (types.DatabaseList, error) {
	dbList := types.DatabaseList{
		Databases: []types.Database{},
		LastUsed:  "",
	}

	if _, err := os.Stat(m.databasePath); os.IsNotExist(err) {
		// Create empty database list
		return dbList, m.SaveDatabaseList(dbList)
	}

	data, err := os.ReadFile(m.databasePath)
	if err != nil {
		return dbList, err
	}

	if err := json.Unmarshal(data, &dbList); err != nil {
		return dbList, err
	}

	return dbList, nil
}

// SaveDatabaseList saves the list of configured databases
func (m *Manager) SaveDatabaseList(dbList types.DatabaseList) error {
	data, err := json.MarshalIndent(dbList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.databasePath, data, 0o600)
}
