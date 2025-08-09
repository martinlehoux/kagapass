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

## User Interface Specification

### File Selection Screen
```
┌─ KagaPass - Select Database ─────────────────────────────────┐
│                                                              │
│  Select KeePass Database:                                    │
│                                                              │
│  ▶ personal.kdbx        (/home/user/passwords/personal.kdbx) │
│    work.kdbx            (/home/user/passwords/work.kdbx)     │
│    family.kdbx          (/home/user/passwords/family.kdbx)   │
│                                                              │
│  [Enter] Open  [Esc] Quit  [a] Add new file                 │
└──────────────────────────────────────────────────────────────┘
```

### Main Search Interface
```
┌─ KagaPass - personal.kdbx ───────────────────────────────────┐
│                                                              │
│  Search: github_                                             │
│  ────────────────────────────────────────────────────────── │
│                                                              │
│  ▶ GitHub Personal      (Personal/Development)              │
│    GitHub Work          (Work/Development)                   │
│    GitHub API Token     (Personal/Development/Tokens)       │
│                                                              │
│                                                              │
│  [b] Copy User  [c] Copy Pass  [Enter] Details  [Esc] Files │
└──────────────────────────────────────────────────────────────┘
```

### Entry Details View
```
┌─ Entry Details ──────────────────────────────────────────────┐
│                                                              │
│  Title:    GitHub Personal                                   │
│  Username: kagamino                                          │
│  Password: ************                                      │
│  URL:      https://github.com                               │
│  Group:    Personal/Development                              │
│                                                              │
│  Notes:                                                      │
│  Recovery codes stored in separate entry.                   │
│  Two-factor authentication enabled.                         │
│  Last backup: 2024-12-15                                    │
│                                                              │
│  Modified: 2024-12-20 14:30:22                              │
│  Created:  2024-01-15 09:15:45                              │
│                                                              │
│  [b] Copy User  [c] Copy Pass  [Esc] Back                   │
└──────────────────────────────────────────────────────────────┘
```

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
- `b`: Copy username to clipboard
- `c`: Copy password to clipboard
- `Enter`: View entry details
- `Esc`: Return to file selection
- `Ctrl+C`: Quit application
- `Ctrl+L`: Clear search

### Entry Details View
- `b`: Copy username to clipboard
- `c`: Copy password to clipboard
- `Esc`: Return to search
- `↑/↓` or `j/k`: Scroll through long notes

## Technical Requirements

### Dependencies
- Go 1.21+
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
- Memory zeroing for sensitive data
- Session timeout configurable (default: until logout)

### Performance Targets
- Database unlock: < 2 seconds
- Search response: < 100ms for databases with 1000+ entries
- Memory usage: < 50MB per open database

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

## Resources

- https://github.com/akazukin5151/kpxhs: TUI is good, but search is only at the folder level
