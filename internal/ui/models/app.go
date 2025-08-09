package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/martinlehoux/kagapass/internal/clipboard"
	"github.com/martinlehoux/kagapass/internal/config"
	"github.com/martinlehoux/kagapass/internal/database"
	"github.com/martinlehoux/kagapass/internal/keyring"
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
	currentDB *types.Database
	entries   []types.Entry

	// Service managers
	dbManager        *database.Manager
	keyringManager   *keyring.Manager
	clipboardManager *clipboard.Manager

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
	dbManager := database.New()
	keyringMgr, err := keyring.New()
	if err != nil {
		// Keyring not available, continue without it
		// TODO: Add warning or fallback behavior
		keyringMgr = nil
	}
	clipboardMgr := clipboard.New()

	app := &AppModel{
		screen:           FileSelectionScreen,
		config:           cfg,
		configMgr:        configMgr,
		databases:        databases,
		dbManager:        dbManager,
		keyringManager:   keyringMgr,
		clipboardManager: clipboardMgr,
	}

	// Initialize screen models
	unlockDatabase := func(database types.Database, password string) tea.Cmd {
		return UnlockDatabase(
			app.dbManager, app.keyringManager, database, password,
		)
	}
	app.fileSelector = NewFileSelectModel(databases, unlockDatabase)
	app.passwordModel = NewPasswordModel(unlockDatabase)
	app.searchModel = NewSearchModel()
	app.searchModel.SetClipboardManager(clipboardMgr)
	app.detailsModel = NewDetailsModel()
	app.detailsModel.SetClipboardManager(clipboardMgr)

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
				selectedDB := &m.databases.Databases[i]
				return UnlockDatabase(m.dbManager, m.keyringManager, *selectedDB, "")
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
	case SwitchScreenMsg:
		return m.switchScreen(msg)
	case UpdateDatabaseListMsg:
		m.databases = msg.DatabaseList
		// Save to config
		if m.configMgr != nil {
			m.configMgr.SaveDatabaseList(m.databases)
		}
		return m, nil
	case DatabaseUnlocked:
		return m, func() tea.Msg {
			return SwitchScreenMsg{
				Screen:   MainSearchScreen,
				Database: &msg.Database,
				Entries:  msg.Entries,
			}
		}
	case DatabaseUnlockFailed:
		if m.screen == PasswordInputScreen {
			var cmd tea.Cmd
			m.passwordModel, cmd = m.passwordModel.Update(msg)
			return m, cmd
		} else {
			return m, func() tea.Msg {
				return SwitchScreenMsg{
					Screen:   PasswordInputScreen,
					Database: &msg.Database,
					Entries:  nil,
				}
			}
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

// switchScreen handles screen switching messages
func (m *AppModel) switchScreen(msg SwitchScreenMsg) (*AppModel, tea.Cmd) {
	m.screen = msg.Screen

	switch msg.Screen {
	case PasswordInputScreen:
		if msg.Database != nil && m.passwordModel != nil {
			m.passwordModel.SetDatabase(msg.Database)
		}
	case MainSearchScreen:
		if msg.Database != nil {
			m.currentDB = msg.Database
			// Set entries if provided (from successful unlock)
			if msg.Entries != nil {
				m.entries = msg.Entries
			} else {
				m.entries = []types.Entry{}
			}
			if m.searchModel != nil {
				m.searchModel.SetEntries(m.entries)
				m.searchModel.SetDatabaseName(msg.Database.Name)
			}
			// Update LastUsed and save to config
			m.databases.LastUsed = msg.Database.Path
			if m.configMgr != nil {
				m.configMgr.SaveDatabaseList(m.databases)
			}
		}
	case EntryDetailsScreen:
		if msg.Entry != nil && m.detailsModel != nil {
			m.detailsModel.SetEntry(*msg.Entry)
		}
	}

	return m, nil
}

// TODO: Remove
// SwitchScreenMsg is sent to switch between screens
type SwitchScreenMsg struct {
	Screen   Screen
	Database *types.Database
	Entry    *types.Entry
	Entries  []types.Entry
}
