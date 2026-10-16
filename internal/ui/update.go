package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update processes incoming Bubble Tea events and key strokes.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.login.active {
		return handleLoginUpdate(m, msg)
	}
	return handleActiveUpdate(m, msg)
}

func handleActiveUpdate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cmd := handleKeyPress(&m, msg); cmd != nil {
			return m, cmd
		}
	case tea.MouseMsg:
		handleMouseClick(&m, msg)
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
	if m.view.activePane == PaneChat {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		return m, inputCmd
	}
	return m, nil
}

func handleKeyPress(m *Model, msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC:
		return tea.Quit
	case tea.KeyEsc:
		if m.view.activePane == PaneChat && len(m.roster) > 0 {
			m.view.activePane = PaneRoster
			m.input.Blur()
			return nil
		}
		return tea.Quit
	case tea.KeyTab:
		togglePane(m)
		return nil
	case tea.KeyEnter:
		if m.view.activePane == PaneRoster {
			m.view.activePane = PaneChat
			m.input.Focus()
			return nil
		}
		return handleEnter(m)
	case tea.KeyUp, tea.KeyDown:
		handleNavigation(m, msg.String())
	}
	if m.view.activePane == PaneRoster && (msg.String() == "j" || msg.String() == "k") {
		handleNavigation(m, msg.String())
	}
	return nil
}

func togglePane(m *Model) {
	if m.view.activePane == PaneRoster {
		m.view.activePane = PaneChat
		m.input.Focus()
	} else {
		m.view.activePane = PaneRoster
		m.input.Blur()
	}
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

func handleMouseClick(m *Model, msg tea.MouseMsg) {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return
	}
	if msg.X <= 26 {
		row := msg.Y - 2
		if row >= 0 && row < len(m.roster) {
			m.selectedIdx = row
			m.activeJID = m.roster[row].JID
			m.view.activePane = PaneChat
			m.input.Focus()
			refreshChatContent(m)
		}
		return
	}
	m.view.activePane = PaneChat
	m.input.Focus()
}

func handleNavigation(m *Model, key string) {
	if m.view.activePane == PaneRoster {
		handleRosterNav(m, key)
		return
	}
	if key == "up" {
		m.viewport.ScrollUp(1)
	} else if key == "down" {
		m.viewport.ScrollDown(1)
	}
}

func handleRosterNav(m *Model, key string) {
	if (key == "up" || key == "k") && m.selectedIdx > 0 {
		m.selectedIdx--
	}
	if (key == "down" || key == "j") && m.selectedIdx < len(m.roster)-1 {
		m.selectedIdx++
	}
	if len(m.roster) > 0 && m.selectedIdx < len(m.roster) {
		m.activeJID = m.roster[m.selectedIdx].JID
		refreshChatContent(m)
	}
}
