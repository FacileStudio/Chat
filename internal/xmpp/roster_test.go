package xmpp

import (
	"testing"

	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/roster"
)

func TestApplyRosterAndPresence(t *testing.T) {
	events := make(chan any, 10)
	c := &Client{
		events: events,
		roster: nil,
	}
	testJID, err := jid.Parse("alice@example.com")
	if err != nil {
		t.Fatalf("parse jid: %v", err)
	}
	ApplyRosterItem(c, roster.Item{
		JID:          testJID,
		Name:         "Alice",
		Subscription: "both",
	})
	if len(c.roster) != 1 || c.roster[0].Name != "Alice" {
		t.Fatalf("unexpected roster: %+v", c.roster)
	}
	ApplyPresence(c, testJID, "chat", "Available", true)
	if !c.roster[0].Online || c.roster[0].Show != "chat" {
		t.Fatalf("presence not applied: %+v", c.roster[0])
	}
}
