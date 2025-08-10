# KagaPass - KeePass TUI

A Terminal User Interface for KeePass databases focused on efficient password retrieval and workflow optimization.

## Overview

KagaPass is a Go-based TUI application designed for quick access to KeePass databases without the overhead of a GUI. It prioritizes search speed, clipboard integration, and session persistence to streamline password management workflows.

## Core Features

### Multi-File Management
- **Persistent File List**: Maintains a permanent list of KeePass database files
- **Session Continuity**: Remembers and reopens the last used database in new sessions
- **Quick Switching**: Easy navigation between different databases via file selection prompt

### Global Fuzzy Search
- **Real-time Search**: Results update as you type
- **Case Insensitive**: Search without worrying about capitalization
- **Fuzzy Matching**: Find entries even with partial or inexact queries
- **Title-based**: Searches entry titles across all groups and folders

### Entry Display & Navigation
- **Path Context**: Shows entries as "Title (Group/Subgroup)" format
- **Keyboard Navigation**: Vim-like movement through search results
- **Quick Access**: Single-key shortcuts for common operations

### Secure Clipboard Integration
- **Auto-copy**: Quick username and password copying to clipboard
- **Auto-clear**: Clipboard automatically cleared after 30 seconds
- **Security**: No password echoing or logging

### Session Persistence
- **Linux Keyring Integration**: Uses system keyring to store master passwords
- **Session-based**: Database remains accessible throughout Linux session
- **Secure Storage**: Master passwords stored using OS-level security

## Keyboard Shortcuts

### File Selection Screen
- `↑/↓` or `j/k`: Navigate file list
- `Enter`: Open selected database
- `a`: Add new database file to list
- `d`: Remove selected file from list
- `Esc`: Quit application

### Main Search Interface
- `Type`: Real-time fuzzy search
- `↑/↓` or `j/k`: Navigate search results
- `Ctrl+B`: Copy username to clipboard
- `Ctrl+C`: Copy password to clipboard
- `Enter`: View entry details
- `Esc`: Return to file selection
- `Ctrl+Q`: Quit application
- `Ctrl+L`: Clear search

### Entry Details View
- `Ctrl+B`: Copy username to clipboard
- `Ctrl+C`: Copy password to clipboard
- `Esc`: Return to search
- `↑/↓` or `j/k`: Scroll through long notes

## Technical Strategy

### Architecture Overview
The application follows a clean architecture pattern with separate layers for UI, business logic, and data access. This ensures maintainability and testability while keeping the codebase simple and focused.

### Core Libraries

#### TUI Framework
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Elm-inspired TUI framework
  - Event-driven architecture perfect for real-time search
  - Clean state management for multiple screens
  - Excellent keyboard handling and composable components

#### KeePass Integration
- **[gokeepasslib](https://github.com/tobischo/gokeepasslib)**: Pure Go KeePass library
  - Native KDBX v3/v4 format support
  - No external dependencies on KeePass binaries
  - Memory-safe handling of encrypted databases

#### Fuzzy Search
- **[fzf](https://github.com/junegunn/fzf) algorithm** or **[fuzzy](https://github.com/sahilm/fuzzy)**: Fast fuzzy matching
  - Optimized for real-time search as-you-type
  - Relevance scoring for better result ordering
  - Case-insensitive matching with highlight support

#### System Integration
- **[keyring](https://github.com/99designs/keyring)**: Cross-platform keyring access
  - Linux Secret Service integration
  - Secure master password storage
  - Session-based credential caching

#### Clipboard Management
- **[clipboard](https://github.com/atotto/clipboard)**: Cross-platform clipboard access
  - Simple read/write operations
  - Background cleanup goroutines
  - Memory-safe password handling

### Security
- Master passwords never written to disk
- No logging of passwords or search terms

### Error Handling Strategy
- Graceful degradation for missing keyring support
- User-friendly error messages for database issues
- Automatic retry logic for transient failures
- Detailed logging for debugging (non-sensitive data only)

### Testing Approach
- Unit tests for core business logic
- Integration tests with mock KeePass databases
- UI component testing with Bubble Tea test utilities
- Security testing for memory leaks and data exposure

## Technical Requirements

### Dependencies
- Go 1.24+
- Linux keyring support (libsecret)
- KeePass database format support (KDBX)

### File Storage
- **Database List**: `~/.config/kagapass/databases.json`
- **Configuration**: `~/.config/kagapass/config.json`
- **Session Cache**: Linux keyring service

### Security Features
- Master passwords stored in OS keyring only
- No password logging or caching to disk
- Automatic clipboard clearing
- Session timeout configurable (default: until logout)

## Configuration

### Default Settings
```json
{
  "clipboard_clear_seconds": 30,
  "search_debounce_ms": 100,
  "max_search_results": 50,
  "session_timeout_hours": 0,
  "default_database_path": ""
}
```

### Database Configuration
```json
{
  "databases": [
    {
      "name": "personal.kdbx",
      "path": "/home/user/passwords/personal.kdbx",
      "last_accessed": "2024-12-20T14:30:22Z"
    }
  ],
  "last_used": "/home/user/passwords/personal.kdbx"
}
```

## Future Considerations
- Password generation
- Entry creation/editing
- TOTP support
- Custom field display
- Multi-database search
- Backup verification
- Tea Model reference vs value
- fmt.Errorf
- password input
- open previous database
- Performance
  - Database unlock: < 2 seconds
  - Search response: < 100ms for databases with 1000+ entries
  - Memory usage: < 50MB per open database

## Resources

- https://github.com/akazukin5151/kpxhs: TUI is good, but search is only at the folder level
