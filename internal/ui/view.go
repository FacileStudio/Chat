package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// View renders the complete terminal interface layout.
func (m Model) View() string {
	if m.view.width == 0 || m.view.height == 0 {
		return "Starting Facile Chat..."
	}
	if m.login.active {
		return renderLogin(m)
	}
	header := renderHeader(m)
	contentHeight := m.view.height - 4
	if contentHeight < 2 {
		contentHeight = 2
	}
	sidebar := renderSidebar(m, contentHeight)
	chatWidth := m.view.width - 28
	if chatWidth < 10 {
		chatWidth = 10
	}
	chat := renderChat(m, chatWidth, contentHeight)
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chat)
	footer := renderFooter(m)
	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer)
}
