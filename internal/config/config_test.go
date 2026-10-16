package config

import (
	"path/filepath"
	"testing"
)

func TestConfigSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	orig := &Config{
		JID:      "test@example.com",
		Password: "secretpassword",
		Server:   "xmpp.example.com",
		Port:     5222,
	}
	if err := SaveFile(path, orig); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if loaded.JID != orig.JID || loaded.Password != orig.Password {
		t.Fatalf("mismatched loaded config: got %+v, want %+v", loaded, orig)
	}
}
