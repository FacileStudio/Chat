package ui

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/Chat/internal/event"
	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model) string {
	title := titleStyle.Render("FACILE CHAT")
	jidStr := headerStyle.Render(m.cfg.JID)
	statusDot := statusDotOffline.String() + " Offline"
	switch m.status.state {
	case event.StateConnected:
		statusDot = statusDotOnline.String() + " Connected"
	case event.StateConnecting:
		statusDot = statusDotAway.String() + " Connecting..."
	case event.StateError:
		statusDot = statusDotOffline.String() + " Error"
	}
	gap := m.view.width - lipgloss.Width(title) - lipgloss.Width(jidStr) - lipgloss.Width(statusDot) - 4
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, " ", jidStr, spacer, statusDot)
}

func renderSidebar(m Model, height int) string {
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Contacts\n\n"))
	if len(m.roster) == 0 {
		sb.WriteString(contactNormal.Render("No contacts yet\nUse /add <jid>"))
	}
	for i, item := range m.roster {
		dot := renderStatusDot(item)
		label := fmt.Sprintf("%s %s", dot, item.Name)
		if i == m.selectedIdx {
			sb.WriteString(contactSelected.Render(label) + "\n")
		} else {
			sb.WriteString(contactNormal.Render(label) + "\n")
		}
	}
	style := sidebarStyle.Width(26).Height(height)
	if m.view.activePane == PaneRoster {
		style = sidebarActiveStyle.Width(26).Height(height)
	}
	return style.Render(sb.String())
}

func renderStatusDot(item event.RosterItem) string {
	if !item.Online {
		return statusDotOffline.String()
	}
	if item.Show == "away" || item.Show == "xa" || item.Show == "dnd" {
		return statusDotAway.String()
	}
	return statusDotOnline.String()
}

func renderChat(m Model, width, height int) string {
	contactHeader := "Conversation: " + m.activeJID
	if m.activeJID == "" {
		contactHeader = "No active conversation"
	}
	header := lipgloss.NewStyle().Bold(true).Foreground(colorActive).Render(contactHeader)
	body := m.viewport.View()
	inputBox := m.input.View()
	content := lipgloss.JoinVertical(lipgloss.Left, header, "\n", body, "\n", inputBox)
	style := chatPaneStyle.Width(width).Height(height)
	if m.view.activePane == PaneChat {
		style = chatPaneActiveStyle.Width(width).Height(height)
	}
	return style.Render(content)
}

func renderFooter(m Model) string {
	hints := "[Tab] Switch Pane  •  [Enter] Send  •  [/help] Commands  •  [Esc/Ctrl+C] Quit"
	status := m.status.text
	if m.err != nil {
		status = "Error: " + m.err.Error()
	}
	gap := m.view.width - lipgloss.Width(hints) - lipgloss.Width(status) - 4
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)
	return footerStyle.Render(hints + spacer + status)
}
