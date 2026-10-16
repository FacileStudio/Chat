package ui

import (
	"context"

	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/event"
	"github.com/FacileStudio/Chat/internal/xmpp"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// New initializes the Bubble Tea model with configuration and client bindings.
func New(cfg *config.Config, client *xmpp.Client, events <-chan any) Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message... (/help for commands)"
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 60

	vp := viewport.New(60, 20)
	vp.SetContent("Welcome to Facile Chat. Select a contact or type /help.")

	return Model{
		cfg:         cfg,
		client:      client,
		events:      events,
		roster:      make([]event.RosterItem, 0),
		messages:    make(map[string][]event.MessageEvent),
		viewport:    vp,
		input:       ti,
		view:        ViewState{activePane: PaneChat},
		status:      ConnStatus{state: event.StateDisconnected, text: "Ready"},
		selectedIdx: 0,
	}
}

// Init starts background event listening and initiates connection.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		connectCmd(m.client),
		listenEvents(m.events),
	)
}

func connectCmd(client *xmpp.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return connectDoneMsg{err: nil}
		}
		err := client.Connect(context.Background())
		return connectDoneMsg{err: err}
	}
}

func listenEvents(events <-chan any) tea.Cmd {
	return func() tea.Msg {
		if events == nil {
			return nil
		}
		ev, ok := <-events
		if !ok {
			return nil
		}
		return xmppEventMsg(ev)
	}
}
