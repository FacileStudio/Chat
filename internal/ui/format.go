package ui

import (
	"strings"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
)

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
		if strings.TrimSpace(msg.Body) == "" {
			return msgEncrypted.Render("🔒 [OMEMO Encrypted for other device]")
		}
		return msgEncrypted.Render("🔒 ") + msg.Body
	}
	return msg.Body
}
