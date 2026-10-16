package ui

import (
	"sort"
	"strings"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
	"github.com/FacileStudio/Chat/internal/xmpp"
	tea "github.com/charmbracelet/bubbletea"
)

func checkFetchHistory(m *Model, targetJID string) tea.Cmd {
	if targetJID == "" || m.client == nil {
		return nil
	}
	if m.view.historyFetched == nil {
		m.view.historyFetched = make(map[string]bool)
	}
	if m.view.historyFetched[targetJID] {
		return nil
	}
	m.view.historyFetched[targetJID] = true
	return fetchHistoryCmd(m.client, targetJID)
}

func fetchHistoryCmd(client *xmpp.Client, targetJID string) tea.Cmd {
	return func() tea.Msg {
		_ = client.FetchHistory(targetJID)
		return nil
	}
}

func insertMessage(msgs []event.MessageEvent, msg event.MessageEvent) []event.MessageEvent {
	for _, existing := range msgs {
		if existing.Body == msg.Body && existing.Sent == msg.Sent && existing.From == msg.From && existing.To == msg.To {
			diff := existing.Time.Sub(msg.Time)
			if diff < 0 {
				diff = -diff
			}
			if diff < 2*time.Second {
				return msgs
			}
		}
	}
	msgs = append(msgs, msg)
	sort.SliceStable(msgs, func(i, j int) bool {
		return msgs[i].Time.Before(msgs[j].Time)
	})
	return msgs
}

func formatMsgTime(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	return t.Format("02 Jan 15:04")
}

func isOMEMOFallback(body string) bool {
	if strings.Contains(body, "OMEMO encrypted message") {
		return true
	}
	return strings.Contains(body, "conversations.im/omemo")
}

func formatMsgBody(msg event.MessageEvent) string {
	if strings.Contains(msg.Body, "[OMEMO Encrypted for other device]") {
		return msgEncrypted.Render("🔒 [OMEMO Encrypted for other device]")
	}
	if isOMEMOFallback(msg.Body) {
		return msgEncrypted.Render("🔒 [OMEMO Encrypted for other device]")
	}
	if strings.HasPrefix(msg.Body, "🔒 [OMEMO") {
		return msgEncrypted.Render(msg.Body)
	}
	if msg.Encrypted {
		return msgEncrypted.Render("🔒 ") + msg.Body
	}
	return msg.Body
}
