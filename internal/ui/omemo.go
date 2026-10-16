package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func cmdOMEMO(m *Model, parts []string) tea.Cmd {
	if len(parts) >= 2 {
		switch parts[1] {
		case "clear", "reset", "purge":
			return cmdOMEMOClear(m)
		case "register", "publish", "announce":
			return cmdOMEMORegister(m)
		}
	}
	info := "OMEMO: Client not initialized."
	if m.client != nil {
		info = m.client.OMEMOInfo()
	}
	infoText := info + "\n\n" +
		"OMEMO (XEP-0384) E2EE Guidance:\n" +
		"  • Messages encrypted for your phone display as '🔒 [OMEMO Encrypted for other device]'.\n" +
		"  • To chat in plain text, disable OMEMO in your phone client or use '/omemo clear'.\n" +
		"  • Commands:\n" +
		"    - /omemo status   : Show local OMEMO identity & fingerprint\n" +
		"    - /omemo register : Re-publish device ID and key bundle to server PEP\n" +
		"    - /omemo clear    : Reset PEP device lists on server"
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

func cmdOMEMORegister(m *Model) tea.Cmd {
	if m.client == nil {
		appendSystemMessage(m, "Not connected to XMPP server.")
		return nil
	}
	appendSystemMessage(m, "Publishing OMEMO device ID and bundle to server PEP...")
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.client.AnnounceOMEMODevice(ctx); err != nil {
			return omemoRegisterDoneMsg{err: err}
		}
		return omemoRegisterDoneMsg{}
	}
}

type omemoDoneMsg struct {
	err error
}

type omemoRegisterDoneMsg struct {
	err error
}

func handleOMEMODone(m *Model, msg omemoDoneMsg) {
	if msg.err != nil {
		appendSystemMessage(m, "Error clearing OMEMO devices: "+msg.err.Error())
		return
	}
	appendSystemMessage(m, "OMEMO device lists cleared on server. Other clients will now send plain text.")
}

func handleOMEMORegisterDone(m *Model, msg omemoRegisterDoneMsg) {
	if msg.err != nil {
		appendSystemMessage(m, "Error publishing OMEMO bundle: "+msg.err.Error())
		return
	}
	appendSystemMessage(m, "OMEMO device ID and bundle published to server PEP successfully.")
}
