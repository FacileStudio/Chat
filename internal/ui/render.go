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
	gap := m.view.width - lipgloss.Width(title) - lipgloss.Width(jidStr) - lipgloss.Width(statusDot) - 2
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)
	line := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", jidStr, spacer, statusDot)
	return lipgloss.NewStyle().Width(m.view.width).MaxHeight(1).Render(line)
}

func renderSidebar(m Model, height int) string {
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Contacts\n"))
	if len(m.roster) == 0 {
		sb.WriteString(contactNormal.Render("No contacts\n/add <jid>"))
	}
	for i, item := range m.roster {
		dot := renderStatusDot(item)
		name := item.Name
		if len(name) > 14 {
			name = name[:13] + "…"
		}
		label := fmt.Sprintf("%s %s", dot, name)
		if i == m.selectedIdx {
			sb.WriteString(contactSelected.Render(label) + "\n")
		} else {
			sb.WriteString(contactNormal.Render(label) + "\n")
		}
	}
	style := sidebarStyle.Width(20).Height(height)
	if m.view.activePane == PaneRoster {
		style = sidebarActiveStyle.Width(20).Height(height)
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
	content := lipgloss.JoinVertical(lipgloss.Left, header, body, inputBox)
	style := chatPaneStyle.Width(width).Height(height)
	if m.view.activePane == PaneChat {
		style = chatPaneActiveStyle.Width(width).Height(height)
	}
	return style.Render(content)
}

func renderFooter(m Model) string {
	hints := "[Tab] Switch Pane • [Enter] Send • [/help] Commands • [Esc] Quit"
	status := m.status.text
	if m.status.err != nil {
		status = "Error: " + m.status.err.Error()
	}
	gap := m.view.width - lipgloss.Width(hints) - lipgloss.Width(status) - 2
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)
	line := hints + spacer + status
	return lipgloss.NewStyle().Width(m.view.width).MaxHeight(1).Render(footerStyle.Render(line))
}

func renderLogin(m Model) string {
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(colorActive).
		Padding(1, 2).
		Width(46)
	title := titleStyle.Render("FACILE CHAT SETUP")
	desc := lipgloss.NewStyle().Foreground(colorMuted).Render("Enter your XMPP credentials or set CHAT_USER / CHAT_PASSWORD.")
	jidLabel := lipgloss.NewStyle().Bold(true).Render("JID (username@domain):")
	passLabel := lipgloss.NewStyle().Bold(true).Render("Password:")
	hint := lipgloss.NewStyle().Foreground(colorMuted).Render("[Tab] Switch Field  •  [Enter] Connect & Save")
	status := ""
	if m.status.text != "" {
		status = lipgloss.NewStyle().Foreground(colorYellow).Render(m.status.text)
	}
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", desc, "", jidLabel, m.input.View(), "", passLabel, m.login.pass.View(), "", hint, status)
	card := boxStyle.Render(content)
	return lipgloss.Place(m.view.width, m.view.height, lipgloss.Center, lipgloss.Center, card)
}
