package ui

import tea "github.com/charmbracelet/bubbletea"

func cmdOMEMO(m *Model, parts []string) tea.Cmd {
	if len(parts) >= 2 && (parts[1] == "clear" || parts[1] == "reset" || parts[1] == "purge") {
		return cmdOMEMOClear(m)
	}
	infoText := "OMEMO (XEP-0384) E2EE Status:\n" +
		"  • Encrypted messages are safely detected and displayed with 🔒.\n" +
		"  • To receive unencrypted messages from Conversations, disable OMEMO via the lock icon.\n" +
		"  • Use '/omemo clear' to reset stale OMEMO device lists on your server."
	appendSystemMessage(m, infoText)
	return nil
}

func cmdOMEMOClear(m *Model) tea.Cmd {
	if m.client == nil {
		appendSystemMessage(m, "Not connected to XMPP server.")
		return nil
	}
	appendSystemMessage(m, "Clearing OMEMO device lists on server...")
	return func() tea.Msg {
		if err := m.client.ClearOMEMODevices(); err != nil {
			return omemoDoneMsg{err: err}
		}
		return omemoDoneMsg{}
	}
}

type omemoDoneMsg struct {
	err error
}

func handleOMEMODone(m *Model, msg omemoDoneMsg) {
	if msg.err != nil {
		appendSystemMessage(m, "Error clearing OMEMO devices: "+msg.err.Error())
		return
	}
	appendSystemMessage(m, "OMEMO device lists cleared on server. Other clients will now send plain text.")
}
