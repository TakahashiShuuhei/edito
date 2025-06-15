package main

import (
	"fmt"
	"os"

	"github.com/TakahashiShuuhei/edito/internal/editor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: edito <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	ed := editor.New()
	
	if err := ed.LoadFile(filename); err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}

	if err := ed.Run(); err != nil {
		fmt.Printf("Error running editor: %v\n", err)
		os.Exit(1)
	}
}