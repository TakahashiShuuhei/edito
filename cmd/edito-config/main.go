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

const Version = "0.2.0"

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
	
	// Create a temporary build directory with proper module setup
	tempDir, err := os.MkdirTemp("", "edito-config-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)
	
	// Initialize go module in temp directory
	if err := initGoModule(tempDir); err != nil {
		fmt.Printf("Failed to initialize Go module: %v\n", err)
		os.Exit(1)
	}
	
	// Read config file content
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}
	
	// Process config content to change package declaration
	processedConfig := processConfigContent(string(configContent))
	
	// Create main.go that includes the config content
	mainGoContent := fmt.Sprintf(`package main

%s

// Export for plugin loading
var ConfigInit = func() {
	// Configuration is loaded via init() functions
}
`, processedConfig)
	
	mainGoFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainGoFile, []byte(mainGoContent), 0644); err != nil {
		fmt.Printf("Failed to write main.go: %v\n", err)
		os.Exit(1)
	}
	
	// Run go mod tidy to resolve dependencies
	if err := runGoModTidy(tempDir); err != nil {
		fmt.Printf("Failed to run go mod tidy: %v\n", err)
		os.Exit(1)
	}
	
	// Build as plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputFile, ".")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Build failed: %v\n%s\n", err, string(output))
		os.Exit(1)
	}
	
	fmt.Printf("Configuration compiled to: %s\n", outputFile)
}

func initGoModule(tempDir string) error {
	// Initialize go module
	cmd := exec.Command("go", "mod", "init", "temp-config")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func runGoModTidy(tempDir string) error {
	// Run go mod tidy to resolve dependencies
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func processConfigContent(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip package declaration since we'll use package main
		if strings.HasPrefix(trimmed, "package ") {
			continue
		}
		result = append(result, line)
	}
	
	return strings.Join(result, "\n")
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