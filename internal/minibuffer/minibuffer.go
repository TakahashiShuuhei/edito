package minibuffer

import (
	"strings"
	
	"github.com/nsf/termbox-go"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeCommand
	ModeSearch
	ModeInput
)

type Completion struct {
	Text        string
	Description string
}

type Minibuffer struct {
	mode         Mode
	prompt       string
	input        string
	cursorPos    int
	completions  []Completion
	selectedComp int
	active       bool
	handler      func(input string) error
}

func New() *Minibuffer {
	return &Minibuffer{
		mode:        ModeNormal,
		completions: make([]Completion, 0),
	}
}

func (mb *Minibuffer) Activate(mode Mode, prompt string, handler func(string) error) {
	mb.mode = mode
	mb.prompt = prompt
	mb.input = ""
	mb.cursorPos = 0
	mb.completions = make([]Completion, 0)
	mb.selectedComp = 0
	mb.active = true
	mb.handler = handler
}

func (mb *Minibuffer) Deactivate() {
	mb.active = false
	mb.input = ""
	mb.cursorPos = 0
	mb.completions = make([]Completion, 0)
	mb.selectedComp = 0
}

func (mb *Minibuffer) IsActive() bool {
	return mb.active
}

func (mb *Minibuffer) GetInput() string {
	return mb.input
}

func (mb *Minibuffer) SetCompletions(completions []Completion) {
	mb.completions = completions
	mb.selectedComp = 0
}

func (mb *Minibuffer) HandleKey(ev termbox.Event) bool {
	if !mb.active {
		return false
	}
	
	if len(mb.completions) > 0 {
		mb.FilterCompletions(mb.input)
	}
	
	switch ev.Key {
	case termbox.KeyEsc:
		mb.Deactivate()
		return true
		
	case termbox.KeyEnter:
		if len(mb.completions) > 0 && mb.selectedComp < len(mb.completions) {
			mb.input = mb.completions[mb.selectedComp].Text
		}
		if mb.handler != nil {
			mb.handler(mb.input)
		}
		mb.Deactivate()
		return true
		
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if mb.cursorPos > 0 {
			mb.input = mb.input[:mb.cursorPos-1] + mb.input[mb.cursorPos:]
			mb.cursorPos--
		}
		return true
		
	case termbox.KeyArrowLeft:
		if mb.cursorPos > 0 {
			mb.cursorPos--
		}
		return true
		
	case termbox.KeyArrowRight:
		if mb.cursorPos < len(mb.input) {
			mb.cursorPos++
		}
		return true
		
	case termbox.KeyArrowUp, termbox.KeyCtrlP:
		if len(mb.completions) > 0 {
			mb.selectedComp = (mb.selectedComp - 1 + len(mb.completions)) % len(mb.completions)
		}
		return true
		
	case termbox.KeyArrowDown, termbox.KeyCtrlN:
		if len(mb.completions) > 0 {
			mb.selectedComp = (mb.selectedComp + 1) % len(mb.completions)
		}
		return true
		
	case termbox.KeyTab:
		if len(mb.completions) > 0 {
			mb.input = mb.completions[mb.selectedComp].Text
			mb.cursorPos = len(mb.input)
		}
		return true
		
	default:
		if ev.Ch != 0 {
			mb.input = mb.input[:mb.cursorPos] + string(ev.Ch) + mb.input[mb.cursorPos:]
			mb.cursorPos++
			return true
		}
	}
	
	return false
}

func (mb *Minibuffer) Draw(width, y int) {
	if !mb.active {
		return
	}
	
	promptText := mb.prompt + mb.input
	
	for i := 0; i < width; i++ {
		ch := ' '
		if i < len(promptText) {
			ch = rune(promptText[i])
		}
		termbox.SetCell(i, y, ch, termbox.ColorWhite, termbox.ColorBlue)
	}
	
	cursorX := len(mb.prompt) + mb.cursorPos
	if cursorX < width {
		termbox.SetCursor(cursorX, y)
	}
	
	if len(mb.completions) > 0 {
		mb.drawCompletions(width, y-len(mb.completions)-1)
	}
}

func (mb *Minibuffer) drawCompletions(width, startY int) {
	if startY < 0 {
		return
	}
	
	maxCompletions := 10
	if len(mb.completions) < maxCompletions {
		maxCompletions = len(mb.completions)
	}
	
	for i := 0; i < maxCompletions; i++ {
		y := startY + i
		completion := mb.completions[i]
		
		bg := termbox.ColorDefault
		fg := termbox.ColorDefault
		
		if i == mb.selectedComp {
			bg = termbox.ColorWhite
			fg = termbox.ColorBlack
		}
		
		text := completion.Text
		if completion.Description != "" {
			text += " - " + completion.Description
		}
		
		if len(text) > width {
			text = text[:width-3] + "..."
		}
		
		for j := 0; j < width; j++ {
			ch := ' '
			if j < len(text) {
				ch = rune(text[j])
			}
			termbox.SetCell(j, y, ch, fg, bg)
		}
	}
}

func (mb *Minibuffer) FilterCompletions(query string) {
	if query == "" {
		return
	}
	
	filtered := make([]Completion, 0)
	query = strings.ToLower(query)
	
	for _, comp := range mb.completions {
		if strings.Contains(strings.ToLower(comp.Text), query) ||
		   strings.Contains(strings.ToLower(comp.Description), query) {
			filtered = append(filtered, comp)
		}
	}
	
	mb.completions = filtered
	if mb.selectedComp >= len(mb.completions) {
		mb.selectedComp = 0
	}
}