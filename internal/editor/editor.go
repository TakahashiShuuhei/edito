package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/TakahashiShuuhei/edito/internal/api"
	"github.com/TakahashiShuuhei/edito/internal/buffer"
	"github.com/TakahashiShuuhei/edito/internal/command"
	"github.com/TakahashiShuuhei/edito/internal/config"
	"github.com/TakahashiShuuhei/edito/internal/keybinding"
	"github.com/TakahashiShuuhei/edito/internal/minibuffer"
	"github.com/TakahashiShuuhei/edito/internal/package_manager"
	"github.com/TakahashiShuuhei/edito/internal/plugin"
)

type Editor struct {
	width          int
	height         int
	quit           bool
	keyMap         *keybinding.KeyMap
	pluginManager  *plugin.Manager
	packageManager *package_manager.Manager
	bufferManager  *buffer.Manager
	commandRegistry *command.Registry
	minibuffer     *minibuffer.Minibuffer
	config         *config.Config
	configSettings map[string]any
	configPlugins  []string
	configKeyBindings map[string]string
	autoInstaller  *plugin.AutoInstaller
	configPluginSpecs []plugin.PluginSpec
	statusMessage  string
	messageTimeout int
}

func New() *Editor {
	e := &Editor{
		configSettings: make(map[string]any),
		configPlugins: make([]string, 0),
		configKeyBindings: make(map[string]string),
		configPluginSpecs: make([]plugin.PluginSpec, 0),
	}
	
	var err error
	e.config, err = config.New()
	if err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}
	
	err = e.loadGoConfig()
	if err != nil {
		fmt.Printf("Failed to load Go config file: %v\n", err)
		os.Exit(1)
	}
	
	e.bufferManager = buffer.NewManager()
	e.commandRegistry = command.NewRegistry()
	e.minibuffer = minibuffer.New()
	
	e.setupCommands()
	e.setupKeyBindings()
	e.setupPluginSystem()
	e.setupAutoInstaller()
	e.setupAPI()
	e.checkAndInstallPlugins()
	
	return e
}

func (e *Editor) setupPluginSystem() {
	pluginDir := e.config.PluginDir()
	
	e.pluginManager = plugin.NewManager()
	e.packageManager = package_manager.NewManager(
		"https://packages.edito.dev",
		pluginDir,
	)
	
	api := &plugin.API{
		RegisterCommand:   e.registerCommand,
		RegisterKeyBinding: e.registerKeyBinding,
		GetCurrentLine:    e.getCurrentLine,
		SetCurrentLine:    e.setCurrentLine,
		GetCursorPosition: e.getCursorPosition,
		SetCursorPosition: e.setCursorPosition,
		InsertText:        e.insertText,
		DeleteText:        e.deleteText,
		ShowMessage:       e.showMessage,
	}
	
	e.pluginManager.SetAPI(api)
	e.loadInstalledPlugins()
}

func (e *Editor) registerCommand(name string, handler func(args []string) error) {
}

func (e *Editor) registerKeyBinding(key string, handler func()) {
}

func (e *Editor) getCurrentLine() string {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil && buf.CursorY < len(buf.Lines) {
		return buf.Lines[buf.CursorY]
	}
	return ""
}

func (e *Editor) setCurrentLine(line string) {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil && buf.CursorY < len(buf.Lines) {
		buf.Lines[buf.CursorY] = line
	}
}

func (e *Editor) getCursorPosition() (int, int) {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil {
		return buf.CursorX, buf.CursorY
	}
	return 0, 0
}

func (e *Editor) setCursorPosition(x, y int) {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil {
		buf.CursorX = x
		buf.CursorY = y
		e.adjustOffset()
	}
}

func (e *Editor) insertText(text string) {
	for _, ch := range text {
		e.insertChar(ch)
	}
}

func (e *Editor) deleteText(start, end int) {
}

func (e *Editor) showMessage(message string) {
	e.statusMessage = message
	e.messageTimeout = 100 // Show message for ~5 seconds (assuming 20fps)
}

func (e *Editor) loadGoConfig() error {
	// Check if config needs to be rebuilt
	if e.shouldRebuildConfig() {
		if err := e.rebuildConfig(); err != nil {
			fmt.Printf("Warning: failed to rebuild config: %v\n", err)
		}
	}
	
	// Try to load and parse Go source first
	api := config.EditorAPI{
		BindKey: e.bindKeyFromConfig,
		LoadPlugin: e.loadPluginFromConfig,
		SetOption: e.setOptionFromConfig,
		RegisterHook: e.registerHookFromConfig,
		InstallPlugin: e.installPluginFromConfig,
	}
	
	err := config.LoadGoConfig(e.config.GoConfigFile(), api)
	if err != nil {
		// If Go config fails, try compiled config as fallback
		compiledErr := config.LoadCompiledConfig(e.config.CompiledConfigFile())
		if compiledErr != nil {
			return fmt.Errorf("failed to load Go config: %v, and compiled config: %v", err, compiledErr)
		}
	}
	
	return nil
}

func (e *Editor) shouldRebuildConfig() bool {
	configFile := e.config.GoConfigFile()
	compiledFile := e.config.CompiledConfigFile()
	
	// Check if config.go exists
	configStat, err := os.Stat(configFile)
	if err != nil {
		return false // No config.go file
	}
	
	// Check if config.so exists
	compiledStat, err := os.Stat(compiledFile)
	if err != nil {
		return true // config.so doesn't exist, need to build
	}
	
	// Check if config.go is newer than config.so
	return configStat.ModTime().After(compiledStat.ModTime())
}

func (e *Editor) rebuildConfig() error {
	fmt.Printf("Config file updated, rebuilding...\n")
	
	// Use internal edito-config functionality
	return e.compileConfig(e.config.GoConfigFile(), e.config.CompiledConfigFile())
}

func (e *Editor) compileConfig(configFile, outputFile string) error {
	// This duplicates some logic from edito-config but allows us to rebuild
	// without requiring the external binary
	
	tempDir, err := os.MkdirTemp("", "edito-config-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Initialize go module in temp directory
	cmd := exec.Command("go", "mod", "init", "temp-config")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %v", err)
	}
	
	// Read config file content
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	
	// Process config content to change package declaration
	processedConfig := e.processConfigContent(string(configContent))
	
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
		return fmt.Errorf("failed to write main.go: %v", err)
	}
	
	// Run go mod tidy to resolve dependencies
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %v", err)
	}
	
	// Build as plugin
	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", outputFile, ".")
	cmd.Dir = tempDir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (e *Editor) processConfigContent(content string) string {
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

func (e *Editor) bindKeyFromConfig(key, command string) {
	e.configKeyBindings[key] = command
	e.bindConfigKey(key, command)
}

func (e *Editor) loadPluginFromConfig(name string) {
	e.configPlugins = append(e.configPlugins, name)
}

func (e *Editor) setOptionFromConfig(key string, value any) {
	e.configSettings[key] = value
}

func (e *Editor) registerHookFromConfig(event string, handler func()) {
}

func (e *Editor) setupAutoInstaller() {
	pluginDir := e.config.PluginDir()
	cacheDir := e.config.CacheDir
	e.autoInstaller = plugin.NewAutoInstaller(pluginDir, cacheDir)
}

func (e *Editor) installPluginFromConfig(name, repository, version string) {
	spec := plugin.PluginSpec{
		Name:       name,
		Repository: repository,
		Version:    version,
	}
	e.configPluginSpecs = append(e.configPluginSpecs, spec)
}

func (e *Editor) checkAndInstallPlugins() {
	if e.autoInstaller == nil || len(e.configPluginSpecs) == 0 {
		return
	}
	
	if err := e.autoInstaller.CheckAndInstallPlugins(e.configPluginSpecs); err != nil {
		fmt.Printf("Warning: error during plugin installation: %v\n", err)
	}
}

func (e *Editor) setupAPI() {
	apiInstance := &api.EditorAPI{}
	
	// Set up the API functions
	// This allows plugins and config to call back into the editor
	
	api.Initialize(apiInstance)
}

func (e *Editor) loadInstalledPlugins() {
	installed, err := e.packageManager.ListInstalled()
	if err != nil {
		return
	}
	
	pluginDir := e.config.PluginDir()
	
	for _, name := range installed {
		pluginPath := filepath.Join(pluginDir, name+".so")
		e.pluginManager.LoadPlugin(pluginPath)
	}
	
	for _, pluginName := range e.configPlugins {
		pluginPath := filepath.Join(pluginDir, pluginName+".so")
		e.pluginManager.LoadPlugin(pluginPath)
	}
}

func (e *Editor) setupKeyBindings() {
	e.keyMap = keybinding.NewKeyMap()
	
	e.keyMap.BindKey(termbox.KeyCtrlQ, func() { e.quit = true })
	e.keyMap.BindKey(termbox.KeyCtrlS, func() { e.saveCurrentBuffer() })
	e.keyMap.BindKey(termbox.KeyArrowUp, func() { e.moveCursor(0, -1) })
	e.keyMap.BindKey(termbox.KeyArrowDown, func() { e.moveCursor(0, 1) })
	e.keyMap.BindKey(termbox.KeyArrowLeft, func() { e.moveCursor(-1, 0) })
	e.keyMap.BindKey(termbox.KeyArrowRight, func() { e.moveCursor(1, 0) })
	e.keyMap.BindKey(termbox.KeyCtrlA, func() { e.moveToLineBeginning() })
	e.keyMap.BindKey(termbox.KeyCtrlE, func() { e.moveToLineEnd() })
	e.keyMap.BindKey(termbox.KeyCtrlP, func() { e.moveCursor(0, -1) })
	e.keyMap.BindKey(termbox.KeyCtrlN, func() { e.moveCursor(0, 1) })
	e.keyMap.BindKey(termbox.KeyCtrlF, func() { e.moveCursor(1, 0) })
	e.keyMap.BindKey(termbox.KeyCtrlB, func() { e.moveCursor(-1, 0) })
	e.keyMap.BindKey(termbox.KeyEnter, func() { e.insertNewline() })
	e.keyMap.BindKey(termbox.KeyBackspace, func() { e.deleteChar() })
	e.keyMap.BindKey(termbox.KeyBackspace2, func() { e.deleteChar() })
	e.keyMap.BindKey(termbox.KeyCtrlX, func() { e.handleCtrlX() })
	
	// M-x (Alt+x) for command mode
	e.keyMap.Bind(0, 'x', termbox.ModAlt, func() { e.activateCommandMode() })
	
	// Additional bindings for command mode (useful in WSL/terminal environments where M-x might not work)
	// F1 as alternative to M-x
	e.keyMap.BindKey(termbox.KeyF1, func() { e.activateCommandMode() })
	// Ctrl+Space as alternative to M-x
	e.keyMap.BindKey(termbox.KeyCtrlSpace, func() { e.activateCommandMode() })
	
	for key, cmd := range e.configKeyBindings {
		e.bindConfigKey(key, cmd)
	}
}

func (e *Editor) LoadFile(filename string) error {
	_, err := e.bufferManager.NewBuffer(filename)
	return err
}

func (e *Editor) Run() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	e.width, e.height = termbox.Size()
	
	e.draw()
	
	for !e.quit {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			e.handleKey(ev)
		} else if ev.Type == termbox.EventResize {
			e.width, e.height = termbox.Size()
		}
		e.draw()
	}
	
	return nil
}

func (e *Editor) handleKey(ev termbox.Event) {
	if e.minibuffer.IsActive() {
		if e.minibuffer.HandleKey(ev) {
			return
		}
	}
	
	if !e.keyMap.Handle(ev) {
		if ev.Ch != 0 {
			e.insertChar(ev.Ch)
		}
	}
}

func (e *Editor) moveCursor(dx, dy int) {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return
	}
	buf.MoveCursor(dx, dy)
}

func (e *Editor) adjustOffset() {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return
	}
	
	if buf.CursorY < buf.OffsetY {
		buf.OffsetY = buf.CursorY
	}
	if buf.CursorY >= buf.OffsetY+e.height-2 {
		buf.OffsetY = buf.CursorY - e.height + 3
	}
	
	if buf.CursorX < buf.OffsetX {
		buf.OffsetX = buf.CursorX
	}
	if buf.CursorX >= buf.OffsetX+e.width {
		buf.OffsetX = buf.CursorX - e.width + 1
	}
}

func (e *Editor) insertChar(ch rune) {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return
	}
	buf.InsertChar(ch)
	e.adjustOffset()
}

func (e *Editor) insertNewline() {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return
	}
	buf.InsertNewline()
	e.adjustOffset()
}

func (e *Editor) deleteChar() {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return
	}
	buf.DeleteChar()
	e.adjustOffset()
}

func (e *Editor) saveCurrentBuffer() error {
	buf := e.bufferManager.GetCurrentBuffer()
	if buf == nil {
		return fmt.Errorf("no current buffer")
	}
	return buf.SaveFile()
}

func (e *Editor) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	buf := e.bufferManager.GetCurrentBuffer()
	if buf != nil {
		e.drawBuffer(buf)
	}
	
	e.drawStatusLine()
	e.minibuffer.Draw(e.width, e.height-1)
	
	// Handle message timeout
	if e.messageTimeout > 0 {
		e.messageTimeout--
	}
	if e.messageTimeout <= 0 {
		e.statusMessage = ""
	}
	
	termbox.Flush()
}

func (e *Editor) drawBuffer(buf *buffer.Buffer) {
	for y := 0; y < e.height-2; y++ {
		lineIndex := y + buf.OffsetY
		if lineIndex >= len(buf.Lines) {
			break
		}
		
		line := buf.Lines[lineIndex]
		for x := 0; x < e.width && x+buf.OffsetX < len(line); x++ {
			ch := line[x+buf.OffsetX]
			termbox.SetCell(x, y, rune(ch), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	
	if !e.minibuffer.IsActive() {
		cursorScreenX := buf.CursorX - buf.OffsetX
		cursorScreenY := buf.CursorY - buf.OffsetY
		if cursorScreenX >= 0 && cursorScreenX < e.width && cursorScreenY >= 0 && cursorScreenY < e.height-2 {
			termbox.SetCursor(cursorScreenX, cursorScreenY)
		}
	}
}

func (e *Editor) drawStatusLine() {
	var statusLine string
	
	// Show message if available, otherwise show buffer info
	if e.statusMessage != "" {
		statusLine = e.statusMessage
	} else {
		buf := e.bufferManager.GetCurrentBuffer()
		if buf != nil {
			modified := ""
			if buf.Modified {
				modified = "*"
			}
			statusLine = fmt.Sprintf("%s%s - Line %d, Col %d", buf.Name, modified, buf.CursorY+1, buf.CursorX+1)
		} else {
			statusLine = "No buffer"
		}
	}
	
	for i, ch := range statusLine {
		if i >= e.width {
			break
		}
		termbox.SetCell(i, e.height-2, ch, termbox.ColorBlack, termbox.ColorWhite)
	}
	
	for i := len(statusLine); i < e.width; i++ {
		termbox.SetCell(i, e.height-2, ' ', termbox.ColorBlack, termbox.ColorWhite)
	}
}