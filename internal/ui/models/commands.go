package models

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagamigo/kcore"
	"github.com/martinlehoux/kagapass/internal/database"
	"github.com/martinlehoux/kagapass/internal/secretstore"
	"github.com/martinlehoux/kagapass/internal/types"
)

type DatabaseUnlocked struct {
	Database types.Database
	Entries  []types.Entry
}

type DatabaseUnlockFailed struct {
	Database types.Database
	Error    error
}

type UnlockDatabase struct {
	databaseManager *database.Manager
	secretStore     secretstore.SecretStore
}

func (u *UnlockDatabase) Handle(database types.Database, password string) tea.Cmd {
	return func() tea.Msg {
		if password != "" {
			entries, err := u.unlockDatabaseWithPassword(database, password)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			if u.secretStore != nil {
				err := u.secretStore.Store(database.Path, []byte(password))
				if err != nil {
					log.Printf("failed to store password in keyring: %v", err)
				} else {
					log.Println("Successfully stored password in keyring:", database.Name)
				}
			}
			return DatabaseUnlocked{Database: database, Entries: entries}
		}
		if u.secretStore != nil {
			password, err := u.secretStore.Get(database.Path)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			entries, err := u.unlockDatabaseWithPassword(database, string(password))
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			return DatabaseUnlocked{Database: database, Entries: entries}
		}
		return DatabaseUnlockFailed{Database: database, Error: errors.New("unknown")}
	}
}

func (u *UnlockDatabase) unlockDatabaseWithPassword(database types.Database, password string) ([]types.Entry, error) {
	err := u.databaseManager.Open(database.Path, password)
	if err != nil {
		return nil, kcore.Wrap(err, "failed to open database")
	}
	entries, err := u.databaseManager.GetEntries()
	if err != nil {
		return nil, kcore.Wrap(err, "failed to get entries")
	}
	return entries, nil
}

type SwitchScreen struct{}

func (s *SwitchScreen) Handle() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
