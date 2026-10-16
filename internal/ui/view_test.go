package ui

import (
	"strings"
	"testing"

	"github.com/FacileStudio/Chat/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func TestViewLoginModal(t *testing.T) {
	cfg := config.Default()
	events := make(chan any, 10)
	m := New(cfg, nil, events)
	m.view.width = 80
	m.view.height = 24

	output := m.View()
	if !strings.Contains(output, "FACILE CHAT SETUP") {
		t.Fatalf("expected login setup modal, got: %s", output)
	}
}

func TestViewMainChat(t *testing.T) {
	cfg := &config.Config{
		JID:      "alice@example.com",
		Password: "secret",
		Port:     5222,
	}
	events := make(chan any, 10)
	m := New(cfg, nil, events)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model := updated.(Model)

	output := model.View()
	if !strings.Contains(output, "alice@example.com") {
		t.Fatalf("expected username in header, got: %s", output)
	}
	if !strings.Contains(output, "Contacts") {
		t.Fatalf("expected Contacts in sidebar, got: %s", output)
	}
}
