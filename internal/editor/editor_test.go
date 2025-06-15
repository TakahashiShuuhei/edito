package editor

import (
	"os"
	"testing"
	"path/filepath"
)

func TestNew(t *testing.T) {
	e := New()
	if e == nil {
		t.Fatal("New() returned nil")
	}
	
	if e.lines == nil {
		t.Error("lines slice not initialized")
	}
	
	if e.keyMap == nil {
		t.Error("keyMap not initialized")
	}
	
	if e.pluginManager == nil {
		t.Error("pluginManager not initialized")
	}
	
	if e.packageManager == nil {
		t.Error("packageManager not initialized")
	}
}

func TestLoadFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "test_edito.txt")
	content := "line1\nline2\nline3"
	
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile)
	
	e := New()
	if err := e.LoadFile(tmpFile); err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}
	
	expectedLines := []string{"line1", "line2", "line3"}
	if len(e.lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(e.lines))
	}
	
	for i, expected := range expectedLines {
		if i >= len(e.lines) || e.lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, e.lines[i])
		}
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	e := New()
	if err := e.LoadFile("nonexistent.txt"); err != nil {
		t.Errorf("LoadFile should create empty file for non-existent file, got error: %v", err)
	}
	
	if len(e.lines) != 1 || e.lines[0] != "" {
		t.Error("LoadFile should create one empty line for non-existent file")
	}
}

func TestInsertChar(t *testing.T) {
	e := New()
	e.lines = []string{"hello"}
	e.cursorX = 5
	e.cursorY = 0
	
	e.insertChar(' ')
	e.insertChar('w')
	e.insertChar('o')
	e.insertChar('r')
	e.insertChar('l')
	e.insertChar('d')
	
	expected := "hello world"
	if e.lines[0] != expected {
		t.Errorf("Expected %q, got %q", expected, e.lines[0])
	}
	
	if e.cursorX != 11 {
		t.Errorf("Expected cursor at position 11, got %d", e.cursorX)
	}
}

func TestMoveCursor(t *testing.T) {
	e := New()
	e.lines = []string{"hello", "world"}
	e.cursorX = 0
	e.cursorY = 0
	
	e.moveCursor(2, 0)
	if e.cursorX != 2 || e.cursorY != 0 {
		t.Errorf("Expected cursor at (2,0), got (%d,%d)", e.cursorX, e.cursorY)
	}
	
	e.moveCursor(0, 1)
	if e.cursorX != 2 || e.cursorY != 1 {
		t.Errorf("Expected cursor at (2,1), got (%d,%d)", e.cursorX, e.cursorY)
	}
	
	e.moveCursor(-10, -10)
	if e.cursorX != 0 || e.cursorY != 0 {
		t.Errorf("Cursor should be bounded at (0,0), got (%d,%d)", e.cursorX, e.cursorY)
	}
}