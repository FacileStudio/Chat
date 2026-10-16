// Package config provides configuration management for the chat CLI.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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

// Path returns the canonical path to the configuration file.
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("read user home: %w", err)
	}
	candidates := []string{
		filepath.Join(home, ".chat.yml"),
		filepath.Join(home, ".chat.yaml"),
		filepath.Join(home, ".config", "chat", "config.yml"),
		filepath.Join(home, ".config", "chat", "config.yaml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "chat", "config.yml"), nil
	}
	return filepath.Join(home, ".config", "chat", "config.yml"), nil
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

// Load reads configuration from available config paths and environment variables.
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	cfg, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	applyEnvOverrides(cfg)
	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if u := os.Getenv("CHAT_USER"); u != "" {
		cfg.JID = u
	} else if j := os.Getenv("CHAT_JID"); j != "" {
		cfg.JID = j
	}
	if p := os.Getenv("CHAT_PASSWORD"); p != "" {
		cfg.Password = p
	} else if p := os.Getenv("CHAT_PASS"); p != "" {
		cfg.Password = p
	}
	if s := os.Getenv("CHAT_SERVER"); s != "" {
		cfg.Server = s
	}
	if p := os.Getenv("CHAT_PORT"); p != "" {
		if port, err := strconv.Atoi(p); err == nil {
			cfg.Port = port
		}
	}
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
