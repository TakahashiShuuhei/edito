package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PluginSpec defines a plugin to be automatically installed
type PluginSpec struct {
	Name       string // プラグイン名 (例: "file-tree")
	Repository string // GitリポジトリURL (例: "github.com/TakahashiShuuhei/edito-file-tree")
	Version    string // バージョンタグ (例: "v0.1.0", "latest", "main")
}

// AutoInstaller handles automatic plugin installation
type AutoInstaller struct {
	pluginDir string
	cacheDir  string
}

// NewAutoInstaller creates a new auto installer
func NewAutoInstaller(pluginDir, cacheDir string) *AutoInstaller {
	return &AutoInstaller{
		pluginDir: pluginDir,
		cacheDir:  cacheDir,
	}
}

// InstallPlugin automatically downloads, builds, and installs a plugin
func (ai *AutoInstaller) InstallPlugin(spec PluginSpec) error {
	soPath := filepath.Join(ai.pluginDir, spec.Name+".so")
	
	// すでに.soファイルが存在する場合はスキップ
	if _, err := os.Stat(soPath); err == nil {
		fmt.Printf("Plugin %s already installed, skipping\n", spec.Name)
		return nil
	}
	
	fmt.Printf("Installing plugin %s from %s@%s\n", spec.Name, spec.Repository, spec.Version)
	
	// 一時ディレクトリでプラグインをビルド
	tempDir, err := ai.createTempBuildDir(spec)
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// プラグインをダウンロード
	if err := ai.downloadPlugin(tempDir, spec); err != nil {
		return fmt.Errorf("failed to download plugin: %v", err)
	}
	
	// プラグインをビルド
	if err := ai.buildPlugin(tempDir, spec); err != nil {
		return fmt.Errorf("failed to build plugin: %v", err)
	}
	
	// プラグインを配置
	if err := ai.installBuiltPlugin(tempDir, spec); err != nil {
		return fmt.Errorf("failed to install plugin: %v", err)
	}
	
	fmt.Printf("Plugin %s installed successfully\n", spec.Name)
	return nil
}

func (ai *AutoInstaller) createTempBuildDir(spec PluginSpec) (string, error) {
	tempDir, err := os.MkdirTemp("", "edito-plugin-"+spec.Name+"-*")
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func (ai *AutoInstaller) downloadPlugin(tempDir string, spec PluginSpec) error {
	// go mod init
	cmd := exec.Command("go", "mod", "init", "temp-plugin-build")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod init failed: %v", err)
	}
	
	// go get でプラグインをダウンロード
	version := spec.Version
	if version == "latest" {
		version = ""
	} else if version != "" && !strings.HasPrefix(version, "@") {
		version = "@" + version
	}
	
	repoURL := spec.Repository + version
	cmd = exec.Command("go", "get", repoURL)
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get %s failed: %v", repoURL, err)
	}
	
	return nil
}

func (ai *AutoInstaller) buildPlugin(tempDir string, spec PluginSpec) error {
	// main.goファイルを作成（プラグインをimportしてre-export）
	mainGoContent := fmt.Sprintf(`package main

import (
	_ "%s"
)

// プラグインのre-export
// この方法でプラグインパッケージをロード可能にする
`, spec.Repository)
	
	mainGoPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}
	
	// go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v", err)
	}
	
	// プラグインとしてビルド
	outputPath := filepath.Join(tempDir, spec.Name+".so")
	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, ".")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (ai *AutoInstaller) installBuiltPlugin(tempDir string, spec PluginSpec) error {
	sourcePath := filepath.Join(tempDir, spec.Name+".so")
	targetPath := filepath.Join(ai.pluginDir, spec.Name+".so")
	
	// プラグインディレクトリを作成
	if err := os.MkdirAll(ai.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin dir: %v", err)
	}
	
	// .soファイルをコピー
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()
	
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %v", err)
	}
	defer targetFile.Close()
	
	if _, err := sourceFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek source file: %v", err)
	}
	
	if _, err := targetFile.ReadFrom(sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}
	
	return nil
}

// CheckAndInstallPlugins checks config and installs missing plugins
func (ai *AutoInstaller) CheckAndInstallPlugins(specs []PluginSpec) error {
	for _, spec := range specs {
		if err := ai.InstallPlugin(spec); err != nil {
			fmt.Printf("Warning: failed to install plugin %s: %v\n", spec.Name, err)
			// 一つのプラグインが失敗しても他を続ける
			continue
		}
	}
	return nil
}