package models

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martinlehoux/kagamigo/kcore"
	"github.com/martinlehoux/kagapass/internal/database"
	"github.com/martinlehoux/kagapass/internal/keyring"
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

func UnlockDatabase(dbMgr *database.Manager, keyringMgr *keyring.Manager, database types.Database, password string) tea.Cmd {
	return func() tea.Msg {
		if password != "" {
			entries, err := unlockDatabaseWithPassword(dbMgr, database, password)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			if keyringMgr != nil {
				keyringMgr.Store(database.Path, password)
				log.Println("Successfully stored password in keyring:", database.Name)
			}
			return DatabaseUnlocked{Database: database, Entries: entries}
		}
		if keyringMgr != nil {
			password, err := keyringMgr.Retrieve(database.Path)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			entries, err := unlockDatabaseWithPassword(dbMgr, database, password)
			if err != nil {
				return DatabaseUnlockFailed{Database: database, Error: err}
			}
			return DatabaseUnlocked{Database: database, Entries: entries}
		}
		return DatabaseUnlockFailed{Database: database, Error: errors.New("unknown")}
	}
}

func unlockDatabaseWithPassword(dbMgr *database.Manager, database types.Database, password string) ([]types.Entry, error) {
	err := dbMgr.Open(database.Path, password)
	if err != nil {
		return nil, kcore.Wrap(err, "failed to open database")
	}
	entries, err := dbMgr.GetEntries()
	if err != nil {
		return nil, kcore.Wrap(err, "failed to get entries")
	}
	return entries, nil
}
