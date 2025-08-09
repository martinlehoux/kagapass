package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/tobischo/gokeepasslib/v3"
)

// Manager handles KeePass database operations
type Manager struct {
	db       *gokeepasslib.Database
	filePath string
}

// New creates a new database manager
func New() *Manager {
	return &Manager{}
}

// Open loads and decrypts a KeePass database
func (m *Manager) Open(filePath, masterPassword string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	defer file.Close()

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(masterPassword)

	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		return fmt.Errorf("failed to decrypt database: %w", err)
	}

	// Unlock protected entries
	db.UnlockProtectedEntries()

	m.db = db
	m.filePath = filePath
	return nil
}

// GetEntries returns all entries from the database
func (m *Manager) GetEntries() ([]types.Entry, error) {
	if m.db == nil {
		return nil, fmt.Errorf("no database loaded")
	}

	var entries []types.Entry

	// Start from the root group
	if m.db.Content != nil && m.db.Content.Root != nil && len(m.db.Content.Root.Groups) > 0 {
		m.collectEntriesFromGroup(&m.db.Content.Root.Groups[0], "", &entries)
	}

	return entries, nil
}

// collectEntriesFromGroup recursively collects entries from a group and its subgroups
func (m *Manager) collectEntriesFromGroup(group *gokeepasslib.Group, groupPath string, entries *[]types.Entry) {
	if group == nil {
		return
	}

	// Process entries in current group
	for _, entry := range group.Entries {
		if entry.Values == nil {
			continue
		}

		entryData := types.Entry{
			Raw:   entry,
			Group: groupPath,
		}

		// Extract common fields
		for _, value := range entry.Values {
			if value.Key == "Title" {
				entryData.Title = value.Value.Content
			} else if value.Key == "UserName" {
				entryData.Username = value.Value.Content
			} else if value.Key == "Password" {
				entryData.Password = value.Value.Content
			} else if value.Key == "URL" {
				entryData.URL = value.Value.Content
			} else if value.Key == "Notes" {
				entryData.Notes = value.Value.Content
			}
		}

		// Extract timestamps
		if entry.Times.CreationTime != nil {
			entryData.Created = entry.Times.CreationTime.Time
		}
		if entry.Times.LastModificationTime != nil {
			entryData.Modified = entry.Times.LastModificationTime.Time
		}

		*entries = append(*entries, entryData)
	}

	// Recursively process subgroups
	for _, subGroup := range group.Groups {
		subGroupPath := groupPath
		if subGroup.Name != "" {
			if subGroupPath != "" {
				subGroupPath += "/"
			}
			subGroupPath += subGroup.Name
		}
		m.collectEntriesFromGroup(&subGroup, subGroupPath, entries)
	}
}

// Close clears sensitive data from memory
func (m *Manager) Close() {
	if m.db != nil {
		// Lock protected entries to clear passwords from memory
		m.db.LockProtectedEntries()
		m.db = nil
	}
	m.filePath = ""
}

// IsOpen returns true if a database is currently loaded
func (m *Manager) IsOpen() bool {
	return m.db != nil
}

// GetFilePath returns the path of the currently loaded database
func (m *Manager) GetFilePath() string {
	return m.filePath
}

// GetFileName returns the filename of the currently loaded database
func (m *Manager) GetFileName() string {
	if m.filePath == "" {
		return ""
	}
	return filepath.Base(m.filePath)
}
