// Example plugin for syntax highlighting
package main

import (
	"fmt"
	"strings"
	
	"github.com/TakahashiShuuhei/edito/internal/plugin"
)

type SyntaxHighlightingPlugin struct {
	api *plugin.API
}

func (p *SyntaxHighlightingPlugin) Name() string {
	return "syntax-highlighting"
}

func (p *SyntaxHighlightingPlugin) Version() string {
	return "1.0.0"
}

func (p *SyntaxHighlightingPlugin) Init(api *plugin.API) error {
	p.api = api
	
	api.RegisterCommand("highlight-syntax", p.highlightSyntax)
	api.RegisterCommand("toggle-highlighting", p.toggleHighlighting)
	
	return nil
}

func (p *SyntaxHighlightingPlugin) Execute(command string, args []string) error {
	switch command {
	case "highlight-syntax":
		return p.highlightSyntax(args)
	case "toggle-highlighting":
		return p.toggleHighlighting(args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (p *SyntaxHighlightingPlugin) highlightSyntax(args []string) error {
	currentLine := p.api.GetCurrentLine()
	
	if strings.Contains(currentLine, "func") ||
	   strings.Contains(currentLine, "package") ||
	   strings.Contains(currentLine, "import") {
		p.api.ShowMessage("Go keyword detected!")
	}
	
	return nil
}

func (p *SyntaxHighlightingPlugin) toggleHighlighting(args []string) error {
	p.api.ShowMessage("Syntax highlighting toggled")
	return nil
}

var Plugin SyntaxHighlightingPlugin