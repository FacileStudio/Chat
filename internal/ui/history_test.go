package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/event"
)

func TestFormatContactRow(t *testing.T) {
	selfItem := event.RosterItem{JID: "alice@example.com", Name: "alice"}
	rowSelfSelected := formatContactRow(selfItem, "alice@example.com", true)
	if rowSelfSelected != "● You (Notes)" {
		t.Fatalf("expected '● You (Notes)', got %q", rowSelfSelected)
	}

	peerItem := event.RosterItem{JID: "bob@example.com", Name: "Bob"}
	rowPeerUnselected := formatContactRow(peerItem, "alice@example.com", false)
	if rowPeerUnselected != "  Bob" {
		t.Fatalf("expected '  Bob', got %q", rowPeerUnselected)
	}

	rowPeerSelected := formatContactRow(peerItem, "alice@example.com", true)
	if rowPeerSelected != "● Bob" {
		t.Fatalf("expected '● Bob', got %q", rowPeerSelected)
	}
}

func TestInsertMessageDeduplicationAndSorting(t *testing.T) {
	t1 := time.Date(2026, 10, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 10, 16, 10, 0, 0, 0, time.UTC)

	m1 := event.MessageEvent{From: "bob@example.com", To: "alice@example.com", Body: "Hello", Time: t2, Sent: false}
	m2 := event.MessageEvent{From: "bob@example.com", To: "alice@example.com", Body: "Earlier", Time: t1, Sent: false}
	m3 := event.MessageEvent{From: "bob@example.com", To: "alice@example.com", Body: "Hello", Time: t2.Add(500 * time.Millisecond), Sent: false}

	var msgs []event.MessageEvent
	msgs = insertMessage(msgs, m1)
	msgs = insertMessage(msgs, m2)
	msgs = insertMessage(msgs, m3)

	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages after dedup, got %d", len(msgs))
	}
	if msgs[0].Body != "Earlier" || msgs[1].Body != "Hello" {
		t.Fatalf("expected sorted messages, got %+v", msgs)
	}
}

func TestFormatMsgBodyOMEMO(t *testing.T) {
	fb1 := "I sent you an OMEMO encrypted message but your client doesn’t seem to support that."
	m1 := event.MessageEvent{Body: fb1}
	res1 := formatMsgBody(m1)
	if !strings.Contains(res1, "🔒 [OMEMO Encrypted Message]") {
		t.Fatalf("expected encrypted badge for fb1, got: %s", res1)
	}

	fb2 := "I sent you an OMEMO encrypted message but your client doesn’t seem to support that. Find more information on https://conversations.im/omemo"
	m2 := event.MessageEvent{Body: fb2}
	res2 := formatMsgBody(m2)
	if !strings.Contains(res2, "🔒 [OMEMO Encrypted Message]") {
		t.Fatalf("expected encrypted badge for fb2, got: %s", res2)
	}

	m3 := event.MessageEvent{Body: "Hello world"}
	res3 := formatMsgBody(m3)
	if res3 != "Hello world" {
		t.Fatalf("expected plain body, got: %s", res3)
	}
}

func TestSidebarRenderOutput(t *testing.T) {
	cfg := &config.Config{JID: "alice@example.com"}
	m := New(cfg, nil, nil)
	m.roster = []event.RosterItem{
		{JID: "alice@example.com", Name: "You (Notes)"},
		{JID: "bob@example.com", Name: "Bob"},
	}
	m.selectedIdx = 0
	out := renderSidebar(m, 6)
	if !strings.Contains(out, "│ ● You (Notes)") {
		t.Fatalf("expected left-aligned selected contact, got:\n%s", out)
	}
	if !strings.Contains(out, "│   Bob") {
		t.Fatalf("expected left-aligned unselected contact, got:\n%s", out)
	}
}
