// Package edito provides the public API for edito configuration and plugins
// This package should be imported by user configuration files and plugins
package edito

// Re-export the API for user convenience
import "github.com/TakahashiShuuhei/edito/internal/api"

// Global editor instance - this will be set by the main edito binary
var editor = api.Editor

// BindKey binds a key combination to a command
// Usage: edito.BindKey("C-x C-s", "save-buffer")
func BindKey(key, command string) {
	if editor != nil {
		editor.BindKey(key, command)
	}
}

// LoadPlugin loads a plugin by name
// Usage: edito.LoadPlugin("syntax-highlighting")
func LoadPlugin(name string) {
	if editor != nil {
		editor.LoadPlugin(name)
	}
}

// SetOption sets an editor option
// Usage: edito.SetOption("tab-width", 4)
func SetOption(key string, value any) {
	if editor != nil {
		editor.SetOption(key, value)
	}
}

// RegisterHook registers an event hook
// Usage: edito.RegisterHook("file-opened", func() { ... })
func RegisterHook(event string, handler func()) {
	if editor != nil {
		editor.RegisterHook(event, handler)
	}
}

// GetCurrentBuffer returns the current active buffer
func GetCurrentBuffer() api.Buffer {
	if editor != nil {
		return editor.GetCurrentBuffer()
	}
	return nil
}

// ShowMessage displays a message to the user
func ShowMessage(message string) {
	if editor != nil {
		editor.ShowMessage(message)
	}
}

// ExecuteCommand executes an editor command
func ExecuteCommand(command string, args []string) error {
	if editor != nil {
		return editor.ExecuteCommand(command, args)
	}
	return nil
}

// InstallPlugin installs a plugin from a git repository
// Usage: edito.InstallPlugin("file-tree", "github.com/TakahashiShuuhei/edito-file-tree", "v0.1.0")
func InstallPlugin(name, repository, version string) {
	if editor != nil {
		editor.InstallPlugin(name, repository, version)
	}
}