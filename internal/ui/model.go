package ui

import (
	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/event"
	"github.com/FacileStudio/Chat/internal/xmpp"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// Pane represents the currently focused UI section.
type Pane int

// Active pane constants.
const (
	PaneRoster Pane = iota
	PaneChat
)

// ViewState encapsulates window dimensions and focused pane.
type ViewState struct {
	width      int
	height     int
	activePane Pane
}

// ConnStatus encapsulates connection state and display text.
type ConnStatus struct {
	state event.ConnectionState
	text  string
}

// Model holds the Bubble Tea application state.
type Model struct {
	cfg         *config.Config
	client      *xmpp.Client
	events      <-chan any
	roster      []event.RosterItem
	selectedIdx int
	activeJID   string
	messages    map[string][]event.MessageEvent
	viewport    viewport.Model
	input       textinput.Model
	view        ViewState
	status      ConnStatus
	err         error
}

type xmppEventMsg any

type connectDoneMsg struct {
	err error
}
