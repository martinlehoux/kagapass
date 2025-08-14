package models

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagamigo/kcore"
	"github.com/martinlehoux/kagapass/internal/clipboard"
	"github.com/martinlehoux/kagapass/internal/config"
	"github.com/martinlehoux/kagapass/internal/keepass"
	"github.com/martinlehoux/kagapass/internal/secretstore"
	"github.com/martinlehoux/kagapass/internal/types"
)

// Screen represents the current screen being displayed
type Screen int

const (
	FileSelectionScreen Screen = iota
	PasswordInputScreen
	MainSearchScreen
	EntryDetailsScreen
)

// AppModel is the main application model
type AppModel struct {
	screen    Screen
	config    types.Config
	configMgr *config.Manager
	databases types.DatabaseList

	// Service managers
	keepassLoader *keepass.Loader
	secretStore   secretstore.SecretStore
	clipboard     *clipboard.Clipboard

	// Commands
	unlockDatabase *UnlockDatabase

	// Screen-specific models
	fileSelector  *FileSelectModel
	passwordModel *PasswordModel
	searchModel   *SearchModel
	detailsModel  *DetailsModel
}

// NewAppModel creates a new application model
func NewAppModel() (*AppModel, error) {
	configMgr, err := config.New()
	if err != nil {
		return nil, err
	}

	cfg, err := configMgr.LoadConfig()
	if err != nil {
		return nil, err
	}

	databases, err := configMgr.LoadDatabaseList()
	if err != nil {
		return nil, err
	}

	// Initialize service managers
	keepassLoader := keepass.NewLoader(os.DirFS("/home/kagamino"))
	secretStore, err := secretstore.NewKeyring()
	kcore.Expect(err, "error initializing keyring")
	clipboard := clipboard.New()
	unlockDatabase := &UnlockDatabase{
		keepassLoader: keepassLoader,
		secretStore:   &secretStore,
	}
	app := &AppModel{
		screen:         FileSelectionScreen,
		config:         cfg,
		configMgr:      configMgr,
		databases:      databases,
		keepassLoader:  keepassLoader,
		secretStore:    &secretStore,
		clipboard:      clipboard,
		unlockDatabase: unlockDatabase,
		fileSelector:   NewFileSelectModel(databases, unlockDatabase),
		passwordModel:  nil,
		searchModel:    nil,
		detailsModel:   nil,
	}

	return app, nil
}

// Init implements tea.Model
func (m *AppModel) Init() tea.Cmd {
	// Check if we should automatically try to unlock the last used database
	if m.databases.LastUsed != "" {
		// Find the database that matches LastUsed path
		for i, db := range m.databases.Databases {
			if db.Path == m.databases.LastUsed {
				// Found the last used database, try to unlock it automatically
				return m.unlockDatabase.Handle(m.databases.Databases[i], []byte{})
			}
		}
	}
	return nil
}

// Update implements tea.Model
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		case "esc":
			return m.handleEscape()
		}
	case UpdateDatabaseListMsg:
		m.databases = msg.DatabaseList
		// Save to config
		if m.configMgr != nil {
			err := m.configMgr.SaveDatabaseList(m.databases)
			if err != nil {
				log.Printf("failed to save database list: %v", err)
			}
		}
		return m, nil
	case DatabaseUnlocked:
		m.switchMainSearchScreen(msg.Database, msg.Entries)
		return m, nil
	case DatabaseUnlockFailed:
		if m.screen == PasswordInputScreen {
			var cmd tea.Cmd
			m.passwordModel, cmd = m.passwordModel.Update(msg)
			return m, cmd
		} else {
			m.switchPasswordInputScreen(msg.Database)
		}
	}

	// Delegate to current screen
	switch m.screen {
	case FileSelectionScreen:
		var cmd tea.Cmd
		m.fileSelector, cmd = m.fileSelector.Update(msg)
		return m, cmd
	case PasswordInputScreen:
		var cmd tea.Cmd
		m.passwordModel, cmd = m.passwordModel.Update(msg)
		return m, cmd
	case MainSearchScreen:
		var cmd tea.Cmd
		m.searchModel, cmd = m.searchModel.Update(msg)
		return m, cmd
	case EntryDetailsScreen:
		var cmd tea.Cmd
		m.detailsModel, cmd = m.detailsModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model
func (m *AppModel) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(m.chooseView())
}

// View implements tea.Model
func (m *AppModel) chooseView() string {
	switch m.screen {
	case FileSelectionScreen:
		return m.fileSelector.View()
	case PasswordInputScreen:
		return m.passwordModel.View()
	case MainSearchScreen:
		return m.searchModel.View()
	case EntryDetailsScreen:
		return m.detailsModel.View()
	}
	return "Loading..."
}

// handleEscape handles the escape key based on current screen
func (m *AppModel) handleEscape() (*AppModel, tea.Cmd) {
	switch m.screen {
	case FileSelectionScreen:
		return m, tea.Quit
	case PasswordInputScreen:
		m.screen = FileSelectionScreen
		return m, nil
	case MainSearchScreen:
		m.screen = FileSelectionScreen
		return m, nil
	case EntryDetailsScreen:
		m.screen = MainSearchScreen
		return m, nil
	}
	return m, nil
}

func (m *AppModel) switchPasswordInputScreen(database types.Database) {
	m.passwordModel = NewPasswordModel(m.unlockDatabase, m.switchFileSelectionScreen, database)
	m.screen = PasswordInputScreen
}

func (m *AppModel) switchFileSelectionScreen() {
	m.fileSelector = NewFileSelectModel(m.databases, m.unlockDatabase)
	m.screen = FileSelectionScreen
}

func (m *AppModel) switchMainSearchScreen(database types.Database, entries []types.Entry) (*AppModel, tea.Cmd) {
	m.searchModel = NewSearchModel(m.clipboard, entries, m.switchEntryDetailsScreen, database.Name)
	m.screen = MainSearchScreen
	m.databases.LastUsed = database.Path

	return m, func() tea.Msg {
		err := m.configMgr.SaveDatabaseList(m.databases)
		if err != nil {
			log.Printf("failed to save database list: %v", err)
		}
		return UpdateDatabaseListMsg{DatabaseList: m.databases}
	}
}

func (m *AppModel) switchEntryDetailsScreen(entry types.Entry) {
	m.detailsModel = NewDetailsModel(m.clipboard, entry)
	m.screen = EntryDetailsScreen
}
