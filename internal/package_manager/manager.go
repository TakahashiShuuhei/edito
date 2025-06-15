package package_manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Package struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Author      string `json:"author"`
}

type Registry struct {
	URL      string
	packages map[string]Package
}

type Manager struct {
	registry     *Registry
	installedDir string
}

func NewManager(registryURL, installedDir string) *Manager {
	return &Manager{
		registry: &Registry{
			URL:      registryURL,
			packages: make(map[string]Package),
		},
		installedDir: installedDir,
	}
}

func (m *Manager) UpdateRegistry() error {
	resp, err := http.Get(m.registry.URL + "/packages.json")
	if err != nil {
		return fmt.Errorf("failed to fetch registry: %v", err)
	}
	defer resp.Body.Close()

	var packages []Package
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return fmt.Errorf("failed to decode registry: %v", err)
	}

	for _, pkg := range packages {
		m.registry.packages[pkg.Name] = pkg
	}

	return nil
}

func (m *Manager) SearchPackage(query string) []Package {
	var results []Package
	for _, pkg := range m.registry.packages {
		if containsIgnoreCase(pkg.Name, query) || containsIgnoreCase(pkg.Description, query) {
			results = append(results, pkg)
		}
	}
	return results
}

func (m *Manager) InstallPackage(name string) error {
	pkg, exists := m.registry.packages[name]
	if !exists {
		return fmt.Errorf("package %s not found in registry", name)
	}

	if err := os.MkdirAll(m.installedDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %v", err)
	}

	resp, err := http.Get(pkg.URL)
	if err != nil {
		return fmt.Errorf("failed to download package: %v", err)
	}
	defer resp.Body.Close()

	filename := filepath.Join(m.installedDir, name+".so")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create package file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save package: %v", err)
	}

	return nil
}

func (m *Manager) UninstallPackage(name string) error {
	filename := filepath.Join(m.installedDir, name+".so")
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to remove package: %v", err)
	}
	return nil
}

func (m *Manager) ListInstalled() ([]string, error) {
	files, err := os.ReadDir(m.installedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read install directory: %v", err)
	}

	var packages []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			name := file.Name()[:len(file.Name())-3]
			packages = append(packages, name)
		}
	}

	return packages, nil
}

func containsIgnoreCase(s, substr string) bool {
	sLower := string([]byte(s))
	substrLower := string([]byte(substr))
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}