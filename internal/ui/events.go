package ui

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/Chat/internal/event"
	tea "github.com/charmbracelet/bubbletea"
)

func handleXMPPEvent(m *Model, ev any) tea.Cmd {
	switch e := ev.(type) {
	case event.StatusEvent:
		handleStatusEvent(m, e)
	case event.MessageEvent:
		handleMessageEvent(m, e)
	case event.RosterEvent:
		handleRosterEvent(m, e)
	case event.PresenceEvent:
		handlePresenceEvent(m, e)
	}
	return listenEvents(m.events)
}

func handleStatusEvent(m *Model, ev event.StatusEvent) {
	m.status.state = ev.State
	m.status.text = ev.Message
	if ev.Err != nil {
		m.err = ev.Err
	}
}

func handleMessageEvent(m *Model, ev event.MessageEvent) {
	target := ev.From
	if ev.Sent {
		target = ev.To
	}
	m.messages[target] = append(m.messages[target], ev)
	if m.activeJID == "" {
		m.activeJID = target
	}
	if m.activeJID == target {
		refreshChatContent(m)
	}
}

func handleRosterEvent(m *Model, ev event.RosterEvent) {
	m.roster = ev.Items
	if m.activeJID == "" && len(m.roster) > 0 {
		m.activeJID = m.roster[0].JID
		refreshChatContent(m)
	}
}

func handlePresenceEvent(m *Model, ev event.PresenceEvent) {
	for i := range m.roster {
		if m.roster[i].JID == ev.JID {
			m.roster[i].Online = ev.Online
			m.roster[i].Show = ev.Show
			m.roster[i].Status = ev.Status
			break
		}
	}
}

func refreshChatContent(m *Model) {
	if m.activeJID == "" {
		m.viewport.SetContent("Select a contact to view conversation.")
		return
	}
	msgs := m.messages[m.activeJID]
	if len(msgs) == 0 {
		m.viewport.SetContent(fmt.Sprintf("No messages with %s yet. Send one below!", m.activeJID))
		return
	}
	var sb strings.Builder
	for _, msg := range msgs {
		t := msgTime.Render(msg.Time.Format("15:04"))
		if msg.Sent {
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", t, msgSelf.Render("You"), msg.Body))
		} else {
			sender := msgPeer.Render(msg.From)
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", t, sender, msg.Body))
		}
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}
