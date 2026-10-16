package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update processes incoming Bubble Tea events and key strokes.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd := handleKeyPress(&m, msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		handleResize(&m, msg)
	case xmppEventMsg:
		cmd := handleXMPPEvent(&m, msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case connectDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status.text = msg.err.Error()
		}
	}
	if m.view.activePane == PaneChat {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		if inputCmd != nil {
			cmds = append(cmds, inputCmd)
		}
	}
	return m, tea.Batch(cmds...)
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
		return nil
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
	} else if key == "down" && m.selectedIdx < len(m.roster)-1 {
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
	sidebarWidth := 28
	chatWidth := msg.Width - sidebarWidth - 4
	if chatWidth < 20 {
		chatWidth = 20
	}
	chatHeight := msg.Height - 7
	if chatHeight < 5 {
		chatHeight = 5
	}
	m.viewport.Width = chatWidth
	m.viewport.Height = chatHeight
	m.input.Width = chatWidth - 4
	refreshChatContent(m)
}
