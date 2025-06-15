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
	
	e.commandRegistry.Register("goto-line", "Go to line number", func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("line number required")
		}
		lineNum, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid line number: %s", args[0])
		}
		return e.gotoLine(lineNum)
	})
	
	e.commandRegistry.Register("quit", "Quit editor", func(args []string) error {
		e.quit = true
		return nil
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
		
		return e.commandRegistry.Execute(commandName, args)
	})
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