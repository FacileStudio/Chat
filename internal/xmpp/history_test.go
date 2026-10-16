package xmpp

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
)

type dummyReadEncoder struct {
	*xml.Decoder
}

func (d dummyReadEncoder) Encode(any) error {
	return nil
}

func (d dummyReadEncoder) EncodeElement(any, xml.StartElement) error {
	return nil
}

func (d dummyReadEncoder) EncodeToken(xml.Token) error {
	return nil
}

func (d dummyReadEncoder) Flush() error {
	return nil
}

func TestNewHistoryHandler(t *testing.T) {
	events := make(chan any, 10)
	c := &Client{
		addr:   jid.MustParse("alice@example.com/res"),
		events: events,
	}
	handler := NewHistoryHandler(c)
	rawXML := `<result xmlns="urn:xmpp:mam:2" id="1" queryid="q1"><forwarded xmlns="urn:xmpp:forward:0"><delay xmlns="urn:xmpp:delay" stamp="2026-10-15T12:00:00Z"/><message xmlns="jabber:client" from="bob@example.com/res" to="alice@example.com/res"><body>Hello from yesterday</body></message></forwarded></result>`
	dec := xml.NewDecoder(strings.NewReader(rawXML))
	re := dummyReadEncoder{Decoder: dec}

	if err := handler(stanza.Message{}, re); err != nil {
		t.Fatalf("handler failed: %v", err)
	}
	select {
	case ev := <-events:
		assertHistoryMsg(t, ev)
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for event")
	}
}

func assertHistoryMsg(t *testing.T, ev any) {
	t.Helper()
	msg, ok := ev.(event.MessageEvent)
	if !ok || msg.From != "bob@example.com" || msg.Body != "Hello from yesterday" || msg.Sent {
		t.Fatalf("unexpected event: %+v", ev)
	}
	expectedTime, _ := time.Parse(time.RFC3339, "2026-10-15T12:00:00Z")
	if !msg.Time.Equal(expectedTime) {
		t.Fatalf("expected time %v, got %v", expectedTime, msg.Time)
	}
}
