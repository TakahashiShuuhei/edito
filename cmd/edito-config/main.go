// Package main provides a tool to compile user configuration files
package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const Version = "0.1.1"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: edito-config <config.go>")
		fmt.Printf("Version: %s\n", Version)
		os.Exit(1)
	}
	
	// Handle version flag
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("edito-config version %s\n", Version)
		os.Exit(0)
	}

	configFile := os.Args[1]
	
	// Get the directory containing the config file
	configDir := filepath.Dir(configFile)
	outputFile := filepath.Join(configDir, "config.so")
	
	// Create a temporary main.go that imports the config
	tempDir, err := os.MkdirTemp("", "edito-config-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)
	
	mainGoContent := fmt.Sprintf(`package main

import (
	_ "%s"
	"github.com/TakahashiShuuhei/edito/pkg/edito"
)

// Export for plugin loading
var ConfigInit = func() {
	// Configuration is loaded via init() functions
}
`, getModulePath(configFile))
	
	mainGoFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainGoFile, []byte(mainGoContent), 0644); err != nil {
		fmt.Printf("Failed to write main.go: %v\n", err)
		os.Exit(1)
	}
	
	// Build as plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputFile, mainGoFile)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Build failed: %v\n%s\n", err, string(output))
		os.Exit(1)
	}
	
	fmt.Printf("Configuration compiled to: %s\n", outputFile)
}

func getModulePath(configFile string) string {
	// Try to determine the module path
	wd, _ := os.Getwd()
	
	// Look for go.mod
	dir := filepath.Dir(configFile)
	for {
		goMod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			// Found go.mod, read module name
			content, err := os.ReadFile(goMod)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "module ") {
						module := strings.TrimSpace(strings.TrimPrefix(line, "module"))
						
						// Calculate relative path from module root to config file
						relPath, _ := filepath.Rel(dir, filepath.Dir(configFile))
						if relPath == "." {
							return module
						}
						return module + "/" + strings.ReplaceAll(relPath, string(filepath.Separator), "/")
					}
				}
			}
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	// Fallback to GOPATH style
	gopath := build.Default.GOPATH
	if gopath != "" {
		srcDir := filepath.Join(gopath, "src")
		if relPath, err := filepath.Rel(srcDir, wd); err == nil {
			return strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		}
	}
	
	return "config"
}