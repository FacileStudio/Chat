// Package ui provides the terminal user interface using Charmbracelet Bubble Tea.
package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorPurple = lipgloss.Color("#7D56F4")
	colorGreen  = lipgloss.Color("#04B575")
	colorYellow = lipgloss.Color("#FFB800")
	colorRed    = lipgloss.Color("#FF4757")
	colorMuted  = lipgloss.Color("#6E6E77")
	colorBorder = lipgloss.Color("#3F3F46")
	colorActive = lipgloss.Color("#A855F7")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorPurple).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D4D4D8"))

	sidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	sidebarActiveStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorActive).
				Padding(0, 1)

	chatPaneStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	chatPaneActiveStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorActive).
				Padding(0, 1)

	contactSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	contactNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A1A1AA"))

	msgSelf = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPurple)

	msgPeer = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorGreen)

	msgTime = lipgloss.NewStyle().
		Foreground(colorMuted)

	msgSys = lipgloss.NewStyle().
		Foreground(colorYellow).
		Italic(true)

	msgEncrypted = lipgloss.NewStyle().
			Foreground(colorYellow).
			Italic(true)

	statusDotOnline = lipgloss.NewStyle().
			Foreground(colorGreen).
			SetString("●")

	statusDotAway = lipgloss.NewStyle().
			Foreground(colorYellow).
			SetString("◌")

	statusDotOffline = lipgloss.NewStyle().
				Foreground(colorMuted).
				SetString("○")

	footerStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)
)
