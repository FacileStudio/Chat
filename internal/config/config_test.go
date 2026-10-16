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

func TestApplyEnvOverrides(t *testing.T) {
	t.Setenv("CHAT_USER", "envuser@domain.com")
	t.Setenv("CHAT_PASSWORD", "envsecret")
	t.Setenv("CHAT_SERVER", "xmpp.custom.com")
	t.Setenv("CHAT_PORT", "5223")

	cfg := Default()
	applyEnvOverrides(cfg)

	if cfg.JID != "envuser@domain.com" {
		t.Fatalf("expected JID envuser@domain.com, got %s", cfg.JID)
	}
	if cfg.Password != "envsecret" {
		t.Fatalf("expected password envsecret, got %s", cfg.Password)
	}
	if cfg.Server != "xmpp.custom.com" || cfg.Port != 5223 {
		t.Fatalf("expected server xmpp.custom.com:5223, got %s:%d", cfg.Server, cfg.Port)
	}
}
