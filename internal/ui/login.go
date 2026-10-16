package ui

import (
	"strings"

	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/xmpp"
	tea "github.com/charmbracelet/bubbletea"
)

func handleLoginUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleLoginKey(&m, msg)
	case tea.WindowSizeMsg:
		handleResize(&m, msg)
	}
	return m, nil
}

func handleLoginKey(m *Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return *m, tea.Quit
	case tea.KeyTab:
		m.login.focus = 1 - m.login.focus
		if m.login.focus == 0 {
			m.input.Focus()
			m.login.pass.Blur()
		} else {
			m.login.pass.Focus()
			m.input.Blur()
		}
		return *m, nil
	case tea.KeyEnter:
		return handleLoginSubmit(m)
	}
	if m.login.focus == 0 {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return *m, cmd
	}
	var cmd tea.Cmd
	m.login.pass, cmd = m.login.pass.Update(msg)
	return *m, cmd
}

func handleLoginSubmit(m *Model) (tea.Model, tea.Cmd) {
	jidVal := strings.TrimSpace(m.input.Value())
	passVal := strings.TrimSpace(m.login.pass.Value())
	if m.login.focus == 0 && jidVal != "" && passVal == "" {
		m.login.focus = 1
		m.login.pass.Focus()
		m.input.Blur()
		return *m, nil
	}
	if jidVal == "" || passVal == "" {
		m.status.text = "Please enter both JID and password"
		return *m, nil
	}
	m.cfg.JID = jidVal
	m.cfg.Password = passVal
	_ = config.Save(m.cfg)
	client, err := xmpp.New(m.cfg, m.events)
	if err != nil {
		m.status.err = err
		return *m, nil
	}
	m.client = client
	m.login.active = false
	m.input.SetValue("")
	m.input.Placeholder = "Type a message... (/help for commands)"
	m.input.Focus()
	return *m, connectCmd(m.client)
}

func handleConnectDone(m *Model, msg connectDoneMsg) {
	if msg.err != nil {
		m.status.err = msg.err
		m.status.text = msg.err.Error()
	}
}
