// Package xmpp provides XMPP client and communication handling.
package xmpp

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/FacileStudio/Chat/internal/event"
	"mellium.im/xmlstream"
	"mellium.im/xmpp/delay"
	"mellium.im/xmpp/history"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
)

type mamForwardedMsg struct {
	To        jid.JID                     `xml:"to,attr"`
	From      jid.JID                     `xml:"from,attr"`
	Body      string                      `xml:"body"`
	Encrypted *struct{ XMLName xml.Name } `xml:"encrypted"`
	EME       *struct{ XMLName xml.Name } `xml:"encryption"`
}

type mamForwarded struct {
	Delay delay.Delay     `xml:"urn:xmpp:delay delay"`
	Msg   mamForwardedMsg `xml:"jabber:client message"`
}

type mamResult struct {
	XMLName xml.Name     `xml:"urn:xmpp:mam:2 result"`
	QueryID string       `xml:"queryid,attr"`
	ID      string       `xml:"id,attr"`
	Forward mamForwarded `xml:"urn:xmpp:forward:0 forwarded"`
}

// NewHistoryHandler creates a message handler for incoming MAM query results.
func NewHistoryHandler(c *Client) mux.MessageHandlerFunc {
	return func(_ stanza.Message, r xmlstream.TokenReadEncoder) error {
		var res mamResult
		d := xml.NewTokenDecoder(r)
		if err := d.Decode(&res); err != nil {
			return err
		}
		isEnc := IsOMEMOMessage(res.Forward.Msg.Body, res.Forward.Msg.Encrypted != nil, res.Forward.Msg.EME != nil)
		if res.Forward.Msg.Body == "" && !isEnc {
			return nil
		}
		msgTime := res.Forward.Delay.Time
		if msgTime.IsZero() {
			msgTime = time.Now()
		}
		fromStr := res.Forward.Msg.From.Bare().String()
		toStr := res.Forward.Msg.To.Bare().String()
		c.events <- event.MessageEvent{
			From:      fromStr,
			To:        toStr,
			Body:      res.Forward.Msg.Body,
			Time:      msgTime,
			Sent:      fromStr == c.addr.Bare().String(),
			Encrypted: isEnc,
		}
		return nil
	}
}

// FetchHistory queries past conversation messages for a specific contact JID.
func (c *Client) FetchHistory(targetJID string) error {
	c.mu.Lock()
	s := c.session
	c.mu.Unlock()
	if s == nil {
		return nil
	}
	target, err := jid.Parse(targetJID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := history.Query{
		With:  target.Bare(),
		Limit: 50,
		End:   time.Now(),
	}
	_, err = history.Fetch(ctx, q, c.addr.Bare(), s)
	return err
}
