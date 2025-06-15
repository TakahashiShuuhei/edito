package config

import (
	"fmt"
	"os"
	"plugin"
)

// LoadCompiledConfig loads a compiled configuration file (.so)
func LoadCompiledConfig(filepath string) error {
	// Check if compiled config exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// No compiled config found, that's OK
		return nil
	}
	
	// Load the plugin
	p, err := plugin.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to load compiled config: %v", err)
	}
	
	// Look for the ConfigInit symbol
	initSymbol, err := p.Lookup("ConfigInit")
	if err != nil {
		// No ConfigInit found, configuration is loaded via init() functions
		// which have already been executed when the plugin was loaded
		return nil
	}
	
	// Call ConfigInit if it exists
	if initFunc, ok := initSymbol.(func()); ok {
		initFunc()
	}
	
	return nil
}