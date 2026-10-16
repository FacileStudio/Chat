package xmpp

import (
	"encoding/xml"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/disco"
	"mellium.im/xmpp/history"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/roster"
	"mellium.im/xmpp/stanza"
)

type rawKey struct {
	RID uint32 `xml:"rid,attr"`
}

type rawKeysJID struct {
	Keys []rawKey `xml:"key"`
}

type rawEncrypted struct {
	Header struct {
		SID     uint32       `xml:"sid,attr"`
		Keys    []rawKey     `xml:"key"`
		KeysJID []rawKeysJID `xml:"keys"`
	} `xml:"header"`
}

type rawMessage struct {
	stanza.Message
	Body      string                      `xml:"body"`
	Encrypted *rawEncrypted               `xml:"encrypted"`
	EME       *struct{ XMLName xml.Name } `xml:"encryption"`
}

// NewXMPPHandler builds the multiplexer for incoming XMPP stanzas.
func NewXMPPHandler(c *Client) xmpp.Handler {
	msgHandler := newMessageHandler(c)
	histHandler := history.Handle(history.NewHandler(NewHistoryHandler(c)))
	return mux.New(
		c.session.In().XMLNS,
		disco.Handle(),
		roster.Handle(roster.Handler{
			Push: func(_ string, item roster.Item) error {
				ApplyRosterItem(c, item)
				return nil
			},
		}),
		mux.Presence("", xml.Name{}, newPresenceHandler(c)),
		mux.Message(stanza.NormalMessage, xml.Name{Local: "body"}, msgHandler),
		mux.Message(stanza.ChatMessage, xml.Name{Local: "body"}, msgHandler),
		mux.Message(stanza.GroupChatMessage, xml.Name{Local: "body"}, msgHandler),
		mux.Message(stanza.NormalMessage, xml.Name{Local: "encrypted"}, msgHandler),
		mux.Message(stanza.ChatMessage, xml.Name{Local: "encrypted"}, msgHandler),
		histHandler,
	)
}

func newMessageHandler(c *Client) mux.MessageHandlerFunc {
	return func(_ stanza.Message, r xmlstream.TokenReadEncoder) error {
		var msg rawMessage
		d := xml.NewTokenDecoder(r)
		if err := d.Decode(&msg); err != nil {
			return err
		}
		isEnc := IsOMEMOMessage(msg.Body, msg.Encrypted != nil, msg.EME != nil)
		if msg.Body == "" && !isEnc {
			return nil
		}
		c.events <- event.MessageEvent{
			From:      msg.From.Bare().String(),
			To:        msg.To.Bare().String(),
			Body:      resolveMessageBody(c, msg, isEnc),
			Time:      time.Now(),
			Sent:      false,
			Encrypted: isEnc,
		}
		return nil
	}
}

func resolveMessageBody(c *Client, msg rawMessage, isEnc bool) string {
	if !isEnc || msg.Encrypted == nil {
		return msg.Body
	}
	if c.omemo != nil && hasKeyForDevice(msg.Encrypted, c.omemo.DeviceID) {
		return "🔒 [OMEMO Encrypted Message]"
	}
	return "🔒 [OMEMO Encrypted for other device]"
}

func hasKeyForDevice(enc *rawEncrypted, devID uint32) bool {
	if enc == nil {
		return false
	}
	for _, k := range enc.Header.Keys {
		if k.RID == devID {
			return true
		}
	}
	for _, kj := range enc.Header.KeysJID {
		for _, k := range kj.Keys {
			if k.RID == devID {
				return true
			}
		}
	}
	return false
}

func newPresenceHandler(c *Client) mux.PresenceHandlerFunc {
	return func(p stanza.Presence, t xmlstream.TokenReadEncoder) error {
		if p.Type == stanza.UnavailablePresence {
			ApplyPresence(c, p.From, "offline", "", false)
			return nil
		}
		show, status := parsePresence(t)
		ApplyPresence(c, p.From, show, status, true)
		return nil
	}
}

func parsePresence(t xmlstream.TokenReadEncoder) (string, string) {
	var show, status string
	for {
		tok, err := t.Token()
		if err != nil || tok == nil {
			break
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local == "show" {
			show = readText(t)
		} else if start.Name.Local == "status" {
			status = readText(t)
		}
	}
	return show, status
}

func readText(t xmlstream.TokenReadEncoder) string {
	tok, err := t.Token()
	if err != nil || tok == nil {
		return ""
	}
	if chars, ok := tok.(xml.CharData); ok {
		_ = xmlstream.Skip(t)
		return string(chars)
	}
	_ = xmlstream.Skip(t)
	return ""
}
