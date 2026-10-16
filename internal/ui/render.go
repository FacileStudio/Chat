package ui

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/Chat/internal/event"
	"github.com/charmbracelet/lipgloss"
)

func renderHeader(m Model) string {
	userTag := m.cfg.JID
	if userTag == "" {
		userTag = "chat"
	}
	userBox := titleStyle.Render(userTag)
	statusDot := statusDotOffline.String() + " Offline"
	switch m.status.state {
	case event.StateConnected:
		statusDot = statusDotOnline.String() + " Connected"
	case event.StateConnecting:
		statusDot = statusDotAway.String() + " Connecting..."
	case event.StateError:
		statusDot = statusDotOffline.String() + " Error"
	}
	gap := m.view.width - lipgloss.Width(userBox) - lipgloss.Width(statusDot) - 2
	if gap < 1 {
		gap = 1
	}
	spacer := strings.Repeat(" ", gap)
	line := lipgloss.JoinHorizontal(lipgloss.Center, userBox, spacer, statusDot)
	return lipgloss.NewStyle().Width(m.view.width).MaxHeight(1).Render(line)
}

func renderSidebar(m Model, height int) string {
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("Contacts\n"))
	if len(m.roster) == 0 {
		sb.WriteString(contactNormal.Render("No contacts\n/add <jid>"))
	}
	for i, item := range m.roster {
		label := formatContactRow(item, m.cfg.JID, i == m.selectedIdx)
		if i == m.selectedIdx {
			sb.WriteString(contactSelected.Render(label) + "\n")
		} else {
			sb.WriteString(contactNormal.Render(label) + "\n")
		}
	}
	style := sidebarStyle.Width(22).Height(height)
	if m.view.activePane == PaneRoster {
		style = sidebarActiveStyle.Width(22).Height(height)
	}
	return style.Render(sb.String())
}

func formatContactRow(item event.RosterItem, selfJID string, isSelected bool) string {
	name := item.Name
	if item.JID == selfJID || (name == "" && item.JID == selfJID) {
		name = "You (Notes)"
	}
	if len(name) > 16 {
		name = name[:15] + "…"
	}
	if isSelected {
		return fmt.Sprintf("%s %s", statusDotOnline.String(), name)
	}
	return fmt.Sprintf("  %s", name)
}

func renderChat(m Model, width, height int) string {
	contactHeader := "Conversation: " + m.activeJID
	if m.activeJID == m.cfg.JID {
		contactHeader = "Conversation: You (Notes)"
	} else if m.activeJID == "" {
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
	hints := "[Tab] Switch Pane • [Enter] Select/Send • [/help] Commands • [Esc] Quit"
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
