package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ConfigDir string
	DataDir   string
	CacheDir  string
}

func New() (*Config, error) {
	cfg := &Config{}
	
	if err := cfg.initXDGDirs(); err != nil {
		return nil, err
	}
	
	if err := cfg.ensureDirs(); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func (c *Config) initXDGDirs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	
	c.ConfigDir = os.Getenv("XDG_CONFIG_HOME")
	if c.ConfigDir == "" {
		c.ConfigDir = filepath.Join(homeDir, ".config")
	}
	c.ConfigDir = filepath.Join(c.ConfigDir, "edito")
	
	c.DataDir = os.Getenv("XDG_DATA_HOME")
	if c.DataDir == "" {
		c.DataDir = filepath.Join(homeDir, ".local", "share")
	}
	c.DataDir = filepath.Join(c.DataDir, "edito")
	
	c.CacheDir = os.Getenv("XDG_CACHE_HOME")
	if c.CacheDir == "" {
		c.CacheDir = filepath.Join(homeDir, ".cache")
	}
	c.CacheDir = filepath.Join(c.CacheDir, "edito")
	
	return nil
}

func (c *Config) ensureDirs() error {
	dirs := []string{c.ConfigDir, c.DataDir, c.CacheDir}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	
	return nil
}

func (c *Config) GoConfigFile() string {
	return filepath.Join(c.ConfigDir, "config.go")
}

func (c *Config) CompiledConfigFile() string {
	return filepath.Join(c.ConfigDir, "config.so")
}

func (c *Config) PluginDir() string {
	return filepath.Join(c.DataDir, "plugins")
}

func (c *Config) CacheFile(name string) string {
	return filepath.Join(c.CacheDir, name)
}