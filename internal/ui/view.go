package ui

import (
	tea "github.com/charmbracelet/bubbletea"
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
	chatWidth := m.view.width - 26
	if chatWidth < 10 {
		chatWidth = 10
	}
	chat := renderChat(m, chatWidth, contentHeight)
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, chat)
	footer := renderFooter(m)
	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer)
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
