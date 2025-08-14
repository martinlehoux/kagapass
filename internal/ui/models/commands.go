package models

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagamigo/kcore"
	"github.com/martinlehoux/kagapass/internal/keepass"
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
	keepassLoader *keepass.Loader
	secretStore   secretstore.SecretStore
}

func (u *UnlockDatabase) Handle(database types.Database, password []byte) tea.Cmd {
	return func() tea.Msg {
		if len(password) > 0 {
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

			entries, err := u.unlockDatabaseWithPassword(database, password)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}

			return DatabaseUnlocked{Database: database, Entries: entries}
		}

		return DatabaseUnlockFailed{Database: database, Error: errors.New("unknown")}
	}
}

func (u *UnlockDatabase) unlockDatabaseWithPassword(database types.Database, password []byte) ([]types.Entry, error) {
	keepass, err := u.keepassLoader.Load(database.Path, password)
	if err != nil {
		return nil, kcore.Wrap(err, "failed to open database")
	}

	defer func() {
		err := keepass.Close()
		if err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	entries, err := keepass.Entries()
	if err != nil {
		return nil, kcore.Wrap(err, "failed to get entries")
	}

	return entries, nil
}
