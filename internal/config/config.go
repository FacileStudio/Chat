// Package config provides configuration management for the chat CLI.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	JID         string `yaml:"jid"`
	Password    string `yaml:"password"`
	Server      string `yaml:"server,omitempty"`
	Port        int    `yaml:"port,omitempty"`
	InsecureTLS bool   `yaml:"insecure_tls,omitempty"`
}

// Default returns a configuration populated with defaults.
func Default() *Config {
	return &Config{
		Port: 5222,
	}
}

// Dir returns the directory where the configuration is stored.
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("read user home: %w", err)
	}
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "chat"), nil
	}
	return filepath.Join(home, ".config", "chat"), nil
}

// Path returns the path to the configuration file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yml"), nil
}

// LoadFile reads configuration from the specified file path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	return cfg, nil
}

// Load reads configuration from the default path or environment.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	cfg, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	if envJID := os.Getenv("CHAT_JID"); envJID != "" {
		cfg.JID = envJID
	}
	if envPass := os.Getenv("CHAT_PASSWORD"); envPass != "" {
		cfg.Password = envPass
	}
	if envServer := os.Getenv("CHAT_SERVER"); envServer != "" {
		cfg.Server = envServer
	}
	return cfg, nil
}

// SaveFile writes configuration to the specified path.
func SaveFile(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

// Save writes configuration to the default path.
func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	return SaveFile(path, cfg)
}
