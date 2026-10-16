package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
	tea "github.com/charmbracelet/bubbletea"
)

func handleCommand(m *Model, input string) tea.Cmd {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	switch parts[0] {
	case "/help":
		return cmdHelp(m)
	case "/msg":
		return cmdMsg(m, parts)
	case "/add":
		return cmdAdd(m, parts)
	case "/status":
		return cmdStatus(m, parts)
	case "/clear":
		if m.activeJID != "" {
			m.messages[m.activeJID] = nil
			refreshChatContent(m)
		}
		return nil
	case "/quit", "/q":
		return tea.Quit
	default:
		appendSystemMessage(m, fmt.Sprintf("Unknown command: %s. Type /help for available commands.", parts[0]))
		return nil
	}
}

func cmdHelp(m *Model) tea.Cmd {
	helpText := "Available commands:\n" +
		"  /msg <jid> <message> - Send direct message\n" +
		"  /add <jid> [name]    - Add contact to roster\n" +
		"  /status <show> [msg] - Set presence (chat, away, dnd, offline)\n" +
		"  /clear               - Clear current chat history\n" +
		"  /quit                - Exit application\n" +
		"Shortcuts:\n" +
		"  Tab: Switch between contacts and chat input\n" +
		"  Up/Down: Select contact or scroll messages\n" +
		"  Enter: Send message"
	appendSystemMessage(m, helpText)
	return nil
}

func cmdMsg(m *Model, parts []string) tea.Cmd {
	if len(parts) < 3 {
		appendSystemMessage(m, "Usage: /msg <jid> <message>")
		return nil
	}
	target := parts[1]
	body := strings.Join(parts[2:], " ")
	if m.client == nil {
		appendSystemMessage(m, "Not connected to XMPP server.")
		return nil
	}
	m.activeJID = target
	if err := m.client.Send(target, body); err != nil {
		appendSystemMessage(m, fmt.Sprintf("Error sending to %s: %v", target, err))
	}
	refreshChatContent(m)
	return nil
}

func cmdAdd(m *Model, parts []string) tea.Cmd {
	if len(parts) < 2 {
		appendSystemMessage(m, "Usage: /add <jid> [name]")
		return nil
	}
	target := parts[1]
	name := target
	if len(parts) >= 3 {
		name = strings.Join(parts[2:], " ")
	}
	if m.client == nil {
		appendSystemMessage(m, "Not connected to XMPP server.")
		return nil
	}
	if err := m.client.AddContact(target, name); err != nil {
		appendSystemMessage(m, fmt.Sprintf("Error adding contact: %v", err))
		return nil
	}
	appendSystemMessage(m, fmt.Sprintf("Subscribed and added %s to roster.", target))
	return nil
}

func cmdStatus(m *Model, parts []string) tea.Cmd {
	if len(parts) < 2 {
		appendSystemMessage(m, "Usage: /status <chat|away|dnd|offline> [message]")
		return nil
	}
	show := parts[1]
	statusMsg := ""
	if len(parts) >= 3 {
		statusMsg = strings.Join(parts[2:], " ")
	}
	if m.client == nil {
		appendSystemMessage(m, "Not connected to XMPP server.")
		return nil
	}
	if err := m.client.SetPresence(show, statusMsg); err != nil {
		appendSystemMessage(m, fmt.Sprintf("Error setting status: %v", err))
		return nil
	}
	appendSystemMessage(m, fmt.Sprintf("Presence updated to %s.", show))
	return nil
}

func appendSystemMessage(m *Model, text string) {
	key := m.activeJID
	if key == "" {
		key = "system"
	}
	m.messages[key] = append(m.messages[key], event.MessageEvent{
		From: "system",
		To:   m.cfg.JID,
		Body: msgSys.Render(text),
		Time: time.Now(),
		Sent: false,
	})
	refreshChatContent(m)
}
