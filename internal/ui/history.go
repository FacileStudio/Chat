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
	if strings.TrimSpace(msg.Body) == "" && !msg.Encrypted {
		return msgs
	}
	if shouldSkipOrUpdate(msgs, msg) {
		return msgs
	}
	msgs = append(msgs, msg)
	sort.SliceStable(msgs, func(i, j int) bool {
		return msgs[i].Time.Before(msgs[j].Time)
	})
	return msgs
}

func shouldSkipOrUpdate(msgs []event.MessageEvent, msg event.MessageEvent) bool {
	for i, existing := range msgs {
		if isSameMessageWindow(existing, msg) {
			return mergeDuplicate(msgs, i, msg)
		}
	}
	return false
}

func isSameMessageWindow(a, b event.MessageEvent) bool {
	if a.From != b.From || a.To != b.To {
		return false
	}
	diff := a.Time.Sub(b.Time)
	if diff < 0 {
		diff = -diff
	}
	return diff < 5*time.Second
}

func mergeDuplicate(msgs []event.MessageEvent, idx int, msg event.MessageEvent) bool {
	existing := msgs[idx]
	if existing.Body == msg.Body {
		return true
	}
	if isPlaceholder(existing.Body) && !isPlaceholder(msg.Body) {
		msgs[idx] = msg
		return true
	}
	return isPlaceholder(msg.Body)
}

func isPlaceholder(body string) bool {
	clean := strings.TrimSpace(body)
	if clean == "" || clean == "🔒" || clean == "🔒 [OMEMO Encrypted Message]" {
		return true
	}
	return strings.Contains(body, "[OMEMO Encrypted") || isOMEMOFallback(body)
}
