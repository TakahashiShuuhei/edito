package buffer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Buffer struct {
	ID        string
	Name      string
	Filename  string
	Lines     []string
	CursorX   int
	CursorY   int
	OffsetX   int
	OffsetY   int
	Modified  bool
	ReadOnly  bool
}

type Manager struct {
	buffers       map[string]*Buffer
	currentBuffer string
	nextID        int
}

func NewManager() *Manager {
	return &Manager{
		buffers: make(map[string]*Buffer),
		nextID:  1,
	}
}

func (m *Manager) NewBuffer(filename string) (*Buffer, error) {
	id := fmt.Sprintf("buffer-%d", m.nextID)
	m.nextID++
	
	buffer := &Buffer{
		ID:       id,
		Name:     filepath.Base(filename),
		Filename: filename,
		Lines:    []string{""},
	}
	
	if filename != "" {
		if err := buffer.LoadFile(filename); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}
	
	m.buffers[id] = buffer
	m.currentBuffer = id
	
	return buffer, nil
}

func (m *Manager) GetBuffer(id string) *Buffer {
	return m.buffers[id]
}

func (m *Manager) GetCurrentBuffer() *Buffer {
	if m.currentBuffer == "" {
		return nil
	}
	return m.buffers[m.currentBuffer]
}

func (m *Manager) SetCurrentBuffer(id string) bool {
	if _, exists := m.buffers[id]; exists {
		m.currentBuffer = id
		return true
	}
	return false
}

func (m *Manager) ListBuffers() []*Buffer {
	buffers := make([]*Buffer, 0, len(m.buffers))
	for _, buffer := range m.buffers {
		buffers = append(buffers, buffer)
	}
	return buffers
}

func (m *Manager) CloseBuffer(id string) bool {
	if _, exists := m.buffers[id]; !exists {
		return false
	}
	
	delete(m.buffers, id)
	
	if m.currentBuffer == id {
		m.currentBuffer = ""
		for bufferID := range m.buffers {
			m.currentBuffer = bufferID
			break
		}
	}
	
	return true
}

func (b *Buffer) LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	b.Lines = make([]string, 0)
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		b.Lines = append(b.Lines, scanner.Text())
	}
	
	if len(b.Lines) == 0 {
		b.Lines = []string{""}
	}
	
	b.Filename = filename
	b.Name = filepath.Base(filename)
	b.Modified = false
	
	return scanner.Err()
}

func (b *Buffer) SaveFile() error {
	if b.Filename == "" {
		return fmt.Errorf("no filename specified")
	}
	
	file, err := os.Create(b.Filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	for i, line := range b.Lines {
		if i > 0 {
			file.WriteString("\n")
		}
		file.WriteString(line)
	}
	
	b.Modified = false
	return nil
}

func (b *Buffer) InsertChar(ch rune) {
	if b.ReadOnly {
		return
	}
	
	if b.CursorY >= len(b.Lines) {
		return
	}
	
	line := b.Lines[b.CursorY]
	if b.CursorX > len(line) {
		b.CursorX = len(line)
	}
	
	newLine := line[:b.CursorX] + string(ch) + line[b.CursorX:]
	b.Lines[b.CursorY] = newLine
	b.CursorX++
	b.Modified = true
}

func (b *Buffer) DeleteChar() {
	if b.ReadOnly {
		return
	}
	
	if b.CursorX == 0 && b.CursorY == 0 {
		return
	}
	
	if b.CursorX == 0 {
		line := b.Lines[b.CursorY]
		b.CursorY--
		b.CursorX = len(b.Lines[b.CursorY])
		b.Lines[b.CursorY] += line
		
		newLines := make([]string, len(b.Lines)-1)
		copy(newLines[:b.CursorY+1], b.Lines[:b.CursorY+1])
		copy(newLines[b.CursorY+1:], b.Lines[b.CursorY+2:])
		b.Lines = newLines
	} else {
		line := b.Lines[b.CursorY]
		b.Lines[b.CursorY] = line[:b.CursorX-1] + line[b.CursorX:]
		b.CursorX--
	}
	
	b.Modified = true
}

func (b *Buffer) InsertNewline() {
	if b.ReadOnly {
		return
	}
	
	if b.CursorY >= len(b.Lines) {
		return
	}
	
	line := b.Lines[b.CursorY]
	newLines := make([]string, len(b.Lines)+1)
	
	copy(newLines[:b.CursorY], b.Lines[:b.CursorY])
	newLines[b.CursorY] = line[:b.CursorX]
	newLines[b.CursorY+1] = line[b.CursorX:]
	copy(newLines[b.CursorY+2:], b.Lines[b.CursorY+1:])
	
	b.Lines = newLines
	b.CursorY++
	b.CursorX = 0
	b.Modified = true
}

func (b *Buffer) MoveCursor(dx, dy int) {
	b.CursorX += dx
	b.CursorY += dy
	
	if b.CursorY < 0 {
		b.CursorY = 0
	}
	if b.CursorY >= len(b.Lines) {
		b.CursorY = len(b.Lines) - 1
	}
	
	if b.CursorX < 0 {
		b.CursorX = 0
	}
	if b.CursorY < len(b.Lines) && b.CursorX > len(b.Lines[b.CursorY]) {
		b.CursorX = len(b.Lines[b.CursorY])
	}
}