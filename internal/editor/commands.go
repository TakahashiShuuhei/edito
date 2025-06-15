package editor

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/nsf/termbox-go"
	"github.com/TakahashiShuuhei/edito/internal/minibuffer"
)

func (e *Editor) setupCommands() {
	e.commandRegistry.Register("save-buffer", "Save current buffer", func(args []string) error {
		return e.saveCurrentBuffer()
	})
	
	e.commandRegistry.Register("kill-buffer", "Close current buffer", func(args []string) error {
		buf := e.bufferManager.GetCurrentBuffer()
		if buf == nil {
			return fmt.Errorf("no current buffer")
		}
		return e.closeBuffer(buf.ID)
	})
	
	e.commandRegistry.Register("switch-to-buffer", "Switch to another buffer", func(args []string) error {
		if len(args) == 0 {
			return e.showBufferList()
		}
		return e.switchToBufferByName(args[0])
	})
	
	e.commandRegistry.Register("find-file", "Open a file", func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("filename required")
		}
		_, err := e.bufferManager.NewBuffer(args[0])
		return err
	})
	
	e.commandRegistry.Register("list-buffers", "List all open buffers", func(args []string) error {
		return e.showBufferList()
	})
	
	e.commandRegistry.RegisterInteractive("goto-line", "Go to line number", func(promptFunc func(prompt string) (string, error)) error {
		input, err := promptFunc("Go to line: ")
		if err != nil {
			return err
		}
		lineNum, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			return fmt.Errorf("invalid line number: %s", input)
		}
		return e.gotoLine(lineNum)
	})
	
	e.commandRegistry.Register("quit", "Quit editor", func(args []string) error {
		e.quit = true
		return nil
	})
	
	e.commandRegistry.Register("help", "Show help and available commands", func(args []string) error {
		return e.showHelp()
	})
	
	e.commandRegistry.Register("list-commands", "List all available commands", func(args []string) error {
		return e.showCommandList()
	})
}

func (e *Editor) activateCommandMode() {
	commands := e.commandRegistry.ListCommands()
	completions := make([]minibuffer.Completion, len(commands))
	
	for i, cmd := range commands {
		completions[i] = minibuffer.Completion{
			Text:        cmd.Name,
			Description: cmd.Description,
		}
	}
	
	e.minibuffer.SetCompletions(completions)
	e.minibuffer.Activate(minibuffer.ModeCommand, "M-x (or F1/C-Space) ", func(input string) error {
		parts := strings.Fields(input)
		if len(parts) == 0 {
			return nil
		}
		
		commandName := parts[0]
		args := parts[1:]
		
		// Check if command exists
		cmd := e.commandRegistry.GetCommand(commandName)
		if cmd == nil {
			e.showMessage(fmt.Sprintf("Command not found: %s. Type 'help' or 'list-commands' to see available commands.", commandName))
			return nil
		}
		
		// Execute interactive command
		if cmd.Interactive != nil {
			err := e.commandRegistry.ExecuteInteractive(commandName, e.promptUser)
			if err != nil {
				e.showMessage(fmt.Sprintf("Command failed: %v", err))
			}
			return nil
		}
		
		// Execute regular command
		err := e.commandRegistry.Execute(commandName, args)
		if err != nil {
			e.showMessage(fmt.Sprintf("Command failed: %v", err))
		}
		return nil
	})
}

func (e *Editor) promptUser(prompt string) (string, error) {
	// Simplified implementation - in practice this would need proper async handling
	// For now, we'll integrate this with the main event loop differently
	
	// This is a placeholder - real implementation would use a proper input mode
	// and integrate with the editor's event handling system
	return "", fmt.Errorf("interactive input not yet fully implemented")
}

func (e *Editor) handleCtrlX() {
}

func (e *Editor) bindConfigKey(key, cmd string) {
	// Convert config key notation to termbox events
	// For now, we'll support basic bindings
	switch key {
	case "C-x C-s":
		e.keyMap.BindKey(termbox.KeyCtrlS, func() {
			e.commandRegistry.Execute(cmd, []string{})
		})
	case "C-x C-c":
		e.keyMap.BindKey(termbox.KeyCtrlQ, func() {
			e.commandRegistry.Execute(cmd, []string{})
		})
	}
}

func (e *Editor) moveToLineBeginning() {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil {
		buf.CursorX = 0
	}
}

func (e *Editor) moveToLineEnd() {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil && buf.CursorY < len(buf.Lines) {
		buf.CursorX = len(buf.Lines[buf.CursorY])
	}
}

func (e *Editor) closeBuffer(id string) error {
	if !e.bufferManager.CloseBuffer(id) {
		return fmt.Errorf("buffer not found: %s", id)
	}
	return nil
}

func (e *Editor) switchToBufferByName(name string) error {
	buffers := e.bufferManager.ListBuffers()
	for _, buf := range buffers {
		if buf.Name == name {
			e.bufferManager.SetCurrentBuffer(buf.ID)
			return nil
		}
	}
	return fmt.Errorf("buffer not found: %s", name)
}

func (e *Editor) showBufferList() error {
	return nil
}

func (e *Editor) gotoLine(lineNum int) error {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return fmt.Errorf("no current buffer")
	}
	
	if lineNum < 1 {
		lineNum = 1
	}
	if lineNum > len(buf.Lines) {
		lineNum = len(buf.Lines)
	}
	
	buf.CursorY = lineNum - 1
	buf.CursorX = 0
	e.adjustOffset()
	
	return nil
}

func (e *Editor) showHelp() error {
	// Create a temporary buffer for help content
	helpContent := []string{
		"Edito - Emacs-like CLI Editor",
		"",
		"Key Bindings:",
		"  Ctrl+Q     - Quit editor",
		"  Ctrl+S     - Save current buffer", 
		"  Ctrl+A     - Move to line beginning",
		"  Ctrl+E     - Move to line end",
		"  Ctrl+P     - Previous line (or Up arrow)",
		"  Ctrl+N     - Next line (or Down arrow)",
		"  Ctrl+F     - Forward character (or Right arrow)",
		"  Ctrl+B     - Backward character (or Left arrow)",
		"",
		"  M-x        - Command palette (or F1, Ctrl+Space)",
		"",
		"Commands (via M-x):",
		"  help           - Show this help",
		"  list-commands  - List all available commands",
		"  goto-line      - Go to specific line number",
		"  save-buffer    - Save current buffer",
		"  find-file      - Open a file",
		"  list-buffers   - List all open buffers",
		"  quit           - Quit editor",
		"",
		"Type any command name in M-x to execute it.",
		"Press any key to close this help.",
	}
	
	return e.showHelpBuffer("*Help*", helpContent)
}

func (e *Editor) showCommandList() error {
	commands := e.commandRegistry.ListCommands()
	var helpContent []string
	
	helpContent = append(helpContent, "Available Commands:")
	helpContent = append(helpContent, "")
	for _, cmd := range commands {
		helpContent = append(helpContent, fmt.Sprintf("  %-15s - %s", cmd.Name, cmd.Description))
	}
	helpContent = append(helpContent, "")
	helpContent = append(helpContent, "Press any key to close this list.")
	
	return e.showHelpBuffer("*Commands*", helpContent)
}

func (e *Editor) showHelpBuffer(name string, content []string) error {
	// For now, show as a simple message
	// In future, this could create a temporary read-only buffer
	message := strings.Join(content, " | ")
	if len(message) > 100 {
		message = message[:100] + "... (Type 'help' or 'list-commands' for full info)"
	}
	e.showMessage(message)
	return nil
}