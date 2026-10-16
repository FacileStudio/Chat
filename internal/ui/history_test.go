package ui

import (
	"testing"
	"time"

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
