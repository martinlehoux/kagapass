package keyring

import (
	"fmt"

	"github.com/99designs/keyring"
)

const (
	serviceName = "kagapass"
	keyPrefix   = "keepass_"
)

// Manager handles secure storage of master passwords
type Manager struct {
	ring keyring.Keyring
}

// New creates a new keyring manager
func New() (*Manager, error) {
	config := keyring.Config{
		ServiceName:          serviceName,
		AllowedBackends:      []keyring.BackendType{keyring.SecretServiceBackend},
		KeychainTrustApplication: true,
	}

	ring, err := keyring.Open(config)
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &Manager{ring: ring}, nil
}

// Store saves a master password for a database file
func (m *Manager) Store(databasePath, masterPassword string) error {
	key := keyPrefix + databasePath
	
	item := keyring.Item{
		Key:  key,
		Data: []byte(masterPassword),
		Label: fmt.Sprintf("KagaPass - %s", databasePath),
		Description: "KeePass master password for KagaPass TUI",
	}

	return m.ring.Set(item)
}

// Retrieve gets the master password for a database file
func (m *Manager) Retrieve(databasePath string) (string, error) {
	key := keyPrefix + databasePath
	
	item, err := m.ring.Get(key)
	if err != nil {
		return "", fmt.Errorf("password not found in keyring: %w", err)
	}

	return string(item.Data), nil
}

// Remove deletes the stored master password for a database file
func (m *Manager) Remove(databasePath string) error {
	key := keyPrefix + databasePath
	return m.ring.Remove(key)
}

// Exists checks if a master password is stored for a database file
func (m *Manager) Exists(databasePath string) bool {
	key := keyPrefix + databasePath
	_, err := m.ring.Get(key)
	return err == nil
}

// ListStored returns all database paths that have stored passwords
func (m *Manager) ListStored() ([]string, error) {
	keys, err := m.ring.Keys()
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, key := range keys {
		if len(key) > len(keyPrefix) && key[:len(keyPrefix)] == keyPrefix {
			path := key[len(keyPrefix):]
			paths = append(paths, path)
		}
	}

	return paths, nil
}