package types

import (
	"time"

	"github.com/tobischo/gokeepasslib/v3"
)

// Database represents a KeePass database configuration
type Database struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	LastAccessed time.Time `json:"last_accessed"`
}

// DatabaseList holds the list of configured databases
type DatabaseList struct {
	Databases []Database `json:"databases"`
	LastUsed  string     `json:"last_used"`
}

// Entry represents a KeePass entry with additional display information
type Entry struct {
	Title    string
	Username string
	Password string
	URL      string
	Notes    string
	Group    string
	Modified time.Time
	Created  time.Time
	Raw      gokeepasslib.Entry
}

// Config holds application configuration
type Config struct {
	ClipboardClearSeconds int    `json:"clipboard_clear_seconds"`
	SearchDebounceMs      int    `json:"search_debounce_ms"`
	MaxSearchResults      int    `json:"max_search_results"`
	SessionTimeoutHours   int    `json:"session_timeout_hours"`
	DefaultDatabasePath   string `json:"default_database_path"`
}

// DefaultConfig returns the default application configuration
func DefaultConfig() Config {
	return Config{
		ClipboardClearSeconds: 30,
		SearchDebounceMs:      100,
		MaxSearchResults:      50,
		SessionTimeoutHours:   0, // 0 means until logout
		DefaultDatabasePath:   "",
	}
}