package plugin

import (
	"errors"
	"fmt"
	"plugin"
	"sync"
)

type Plugin interface {
	Name() string
	Version() string
	Init(api *API) error
	Execute(command string, args []string) error
}

type API struct {
	RegisterCommand func(name string, handler func(args []string) error)
	RegisterKeyBinding func(key string, handler func())
	GetCurrentLine func() string
	SetCurrentLine func(line string)
	GetCursorPosition func() (int, int)
	SetCursorPosition func(x, y int)
	InsertText func(text string)
	DeleteText func(start, end int)
	ShowMessage func(message string)
}

type Manager struct {
	plugins map[string]Plugin
	loaded  map[string]*plugin.Plugin
	api     *API
	mutex   sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
		loaded:  make(map[string]*plugin.Plugin),
		api:     &API{},
	}
}

func (m *Manager) SetAPI(api *API) {
	m.api = api
}

func (m *Manager) LoadPlugin(path string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %v", path, err)
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'Plugin' symbol: %v", path, err)
	}

	pluginInstance, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	name := pluginInstance.Name()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already loaded", name)
	}

	if err := pluginInstance.Init(m.api); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
	}

	m.plugins[name] = pluginInstance
	m.loaded[name] = p

	return nil
}

func (m *Manager) UnloadPlugin(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	delete(m.plugins, name)
	delete(m.loaded, name)

	return nil
}

func (m *Manager) GetPlugin(name string) (Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if plugin, exists := m.plugins[name]; exists {
		return plugin, nil
	}

	return nil, errors.New("plugin not found")
}

func (m *Manager) ListPlugins() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}

	return names
}

func (m *Manager) ExecuteCommand(name string, command string, args []string) error {
	plugin, err := m.GetPlugin(name)
	if err != nil {
		return err
	}

	return plugin.Execute(command, args)
}