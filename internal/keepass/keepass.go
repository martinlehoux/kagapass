package keepass

import (
	"io/fs"
	"log"
	"time"

	"github.com/martinlehoux/kagapass/internal/types"
	"github.com/tobischo/gokeepasslib/v3"
)

type Loader struct {
	fs fs.FS
}

func NewLoader(fs fs.FS) *Loader {
	return &Loader{
		fs: fs,
	}
}

func (m *Loader) Load(path string, password []byte) (*KeePass, error) {
	file, err := m.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing keepass file: %v", err)
		}
	}()

	database := gokeepasslib.NewDatabase()
	database.Credentials = gokeepasslib.NewPasswordCredentials(string(password))
	err = gokeepasslib.NewDecoder(file).Decode(database)
	if err != nil {
		return nil, err
	}
	err = database.UnlockProtectedEntries()
	return &KeePass{
		database: database,
	}, err
}

type KeePass struct {
	database *gokeepasslib.Database
}

func (k *KeePass) Entries() ([]types.Entry, error) {
	var entries []types.Entry

	// Start from the root group
	if k.database.Content != nil && k.database.Content.Root != nil && len(k.database.Content.Root.Groups) > 0 {
		entries = append(entries, collectEntriesFromGroup(&k.database.Content.Root.Groups[0], "")...)
	}

	return entries, nil
}

func (k *KeePass) Close() error {
	return k.database.LockProtectedEntries()
}

func collectEntriesFromGroup(group *gokeepasslib.Group, groupPath string) []types.Entry {
	var entries []types.Entry

	if group == nil {
		return entries
	}

	// Process entries in current group
	for _, entry := range group.Entries {
		if entry.Values == nil {
			continue
		}

		entryData := types.Entry{
			Raw:      entry,
			Group:    groupPath,
			Title:    "",
			Username: "",
			Password: "",
			URL:      "",
			Notes:    "",
			Created:  time.Time{},
			Modified: time.Time{},
		}

		// Extract common fields
		for _, value := range entry.Values {
			switch value.Key {
			case "Title":
				entryData.Title = value.Value.Content
			case "UserName":
				entryData.Username = value.Value.Content
			case "Password":
				entryData.Password = value.Value.Content
			case "URL":
				entryData.URL = value.Value.Content
			case "Notes":
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

		entries = append(entries, entryData)
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
		entries = append(entries, collectEntriesFromGroup(&subGroup, subGroupPath)...)
	}

	return entries
}
