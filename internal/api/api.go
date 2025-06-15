// Package api provides the public API for edito plugins and configuration
package api

import (
	"github.com/nsf/termbox-go"
)

// Editor represents the main editor instance
var Editor *EditorAPI

// EditorAPI provides the public API for plugins and configuration
type EditorAPI struct {
	bindKey       func(key, command string)
	loadPlugin    func(name string)
	setOption     func(key string, value any)
	registerHook  func(event string, handler func())
	getCurrentBuffer func() Buffer
	showMessage   func(message string)
	executeCommand func(command string, args []string) error
}

// Buffer represents a text buffer
type Buffer interface {
	GetLines() []string
	SetLines(lines []string)
	GetCursorPosition() (x, y int)
	SetCursorPosition(x, y int)
	GetFilename() string
	IsModified() bool
	Save() error
}

// Initialize sets up the global Editor API
func Initialize(api *EditorAPI) {
	Editor = api
}

// BindKey binds a key combination to a command
func (e *EditorAPI) BindKey(key, command string) {
	if e.bindKey != nil {
		e.bindKey(key, command)
	}
}

// LoadPlugin loads a plugin by name
func (e *EditorAPI) LoadPlugin(name string) {
	if e.loadPlugin != nil {
		e.loadPlugin(name)
	}
}

// SetOption sets an editor option
func (e *EditorAPI) SetOption(key string, value any) {
	if e.setOption != nil {
		e.setOption(key, value)
	}
}

// RegisterHook registers an event hook
func (e *EditorAPI) RegisterHook(event string, handler func()) {
	if e.registerHook != nil {
		e.registerHook(event, handler)
	}
}

// GetCurrentBuffer returns the current active buffer
func (e *EditorAPI) GetCurrentBuffer() Buffer {
	if e.getCurrentBuffer != nil {
		return e.getCurrentBuffer()
	}
	return nil
}

// ShowMessage displays a message to the user
func (e *EditorAPI) ShowMessage(message string) {
	if e.showMessage != nil {
		e.showMessage(message)
	}
}

// ExecuteCommand executes an editor command
func (e *EditorAPI) ExecuteCommand(command string, args []string) error {
	if e.executeCommand != nil {
		return e.executeCommand(command, args)
	}
	return nil
}

// KeyBinding represents a key binding
type KeyBinding struct {
	Key     termbox.Key
	Char    rune
	Command string
}

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Version() string
	Init(api *EditorAPI) error
	Cleanup() error
}