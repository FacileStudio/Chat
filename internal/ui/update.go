package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update processes incoming Bubble Tea events and key strokes.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.login.active {
			return handleLoginKey(&m, msg)
		}
		cmd := handleKeyPress(&m, msg)
		if cmd != nil {
			return m, cmd
		}
	case tea.WindowSizeMsg:
		handleResize(&m, msg)
	case xmppEventMsg:
		return m, handleXMPPEvent(&m, msg)
	case connectDoneMsg:
		if msg.err != nil {
			m.status.err = msg.err
			m.status.text = msg.err.Error()
		}
	}
	if m.view.activePane == PaneChat && !m.login.active {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		return m, inputCmd
	}
	return m, nil
}

func handleKeyPress(m *Model, msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return tea.Quit
	case tea.KeyTab:
		if m.view.activePane == PaneRoster {
			m.view.activePane = PaneChat
			m.input.Focus()
		} else {
			m.view.activePane = PaneRoster
			m.input.Blur()
		}
		return nil
	case tea.KeyEnter:
		return handleEnter(m)
	case tea.KeyUp, tea.KeyDown:
		handleNavigation(m, msg.String())
	}
	return nil
}

func handleEnter(m *Model) tea.Cmd {
	val := strings.TrimSpace(m.input.Value())
	if val == "" {
		return nil
	}
	m.input.SetValue("")
	if strings.HasPrefix(val, "/") {
		return handleCommand(m, val)
	}
	if m.activeJID == "" {
		appendSystemMessage(m, "No contact selected. Choose a contact or use /msg <jid> <text>.")
		return nil
	}
	if m.client != nil {
		if err := m.client.Send(m.activeJID, val); err != nil {
			appendSystemMessage(m, "Failed to send: "+err.Error())
		}
	}
	refreshChatContent(m)
	return nil
}

func handleNavigation(m *Model, key string) {
	if m.view.activePane == PaneRoster {
		handleRosterNav(m, key)
		return
	}
	handleChatNav(m, key)
}

func handleRosterNav(m *Model, key string) {
	if key == "up" && m.selectedIdx > 0 {
		m.selectedIdx--
	}
	if key == "down" && m.selectedIdx < len(m.roster)-1 {
		m.selectedIdx++
	}
	if len(m.roster) > 0 && m.selectedIdx < len(m.roster) {
		m.activeJID = m.roster[m.selectedIdx].JID
		refreshChatContent(m)
	}
}

func handleChatNav(m *Model, key string) {
	if key == "up" {
		m.viewport.ScrollUp(1)
	} else if key == "down" {
		m.viewport.ScrollDown(1)
	}
}

func handleResize(m *Model, msg tea.WindowSizeMsg) {
	m.view.width = msg.Width
	m.view.height = msg.Height
	sidebarW := 24
	chatW := msg.Width - sidebarW
	innerChatW := chatW - 4
	if innerChatW < 10 {
		innerChatW = 10
	}
	innerChatH := msg.Height - 6
	if innerChatH < 3 {
		innerChatH = 3
	}
	m.viewport.Width = innerChatW
	m.viewport.Height = innerChatH
	m.input.Width = innerChatW - 2
	m.login.pass.Width = 36
	refreshChatContent(m)
}
