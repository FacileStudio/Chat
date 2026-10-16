package xmpp

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/FacileStudio/Chat/internal/event"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/roster"
	"mellium.im/xmpp/stanza"
)

// FetchRoster retrieves the user's contact roster from the server.
func FetchRoster(ctx context.Context, s *xmpp.Session, c *Client) error {
	iter := roster.Fetch(ctx, s)
	defer iter.Close()
	for iter.Next() {
		item := iter.Item()
		ApplyRosterItem(c, item)
	}
	return iter.Err()
}

// ApplyRosterItem updates or adds a contact item in the client roster.
func ApplyRosterItem(c *Client, item roster.Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	bare := item.JID.Bare().String()
	name := item.Name
	if name == "" {
		name = item.JID.Localpart()
	}
	if name == "" {
		name = bare
	}
	findOrCreateRosterItem(c, bare, name, string(item.Subscription))
	items := make([]event.RosterItem, len(c.roster))
	copy(items, c.roster)
	c.events <- event.RosterEvent{Items: items}
}

func findOrCreateRosterItem(c *Client, bare, name, sub string) {
	for i := range c.roster {
		if c.roster[i].JID == bare {
			c.roster[i].Name = name
			c.roster[i].Subscription = sub
			return
		}
	}
	c.roster = append(c.roster, event.RosterItem{
		JID:          bare,
		Name:         name,
		Subscription: sub,
		Online:       false,
	})
}

// ApplyPresence updates contact online status from a presence stanza.
func ApplyPresence(c *Client, from jid.JID, show, status string, online bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	bare := from.Bare().String()
	for i := range c.roster {
		if c.roster[i].JID == bare {
			c.roster[i].Online = online
			c.roster[i].Show = show
			c.roster[i].Status = status
			break
		}
	}
	items := make([]event.RosterItem, len(c.roster))
	copy(items, c.roster)
	c.events <- event.RosterEvent{Items: items}
	c.events <- event.PresenceEvent{
		JID:    bare,
		Show:   show,
		Status: status,
		Online: online,
	}
}

// SendPresence sends an initial or updated presence stanza.
func SendPresence(ctx context.Context, s *xmpp.Session, show, status string) error {
	p := stanza.Presence{Type: stanza.AvailablePresence}
	if strings.EqualFold(show, "offline") {
		p.Type = stanza.UnavailablePresence
		return s.Send(ctx, p.Wrap(nil))
	}
	var readers []xml.TokenReader
	if show != "" {
		r := xmlstream.Wrap(
			xmlstream.Token(xml.CharData(show)),
			xml.StartElement{Name: xml.Name{Local: "show"}},
		)
		readers = append(readers, r)
	}
	if status != "" {
		r := xmlstream.Wrap(
			xmlstream.Token(xml.CharData(status)),
			xml.StartElement{Name: xml.Name{Local: "status"}},
		)
		readers = append(readers, r)
	}
	if len(readers) == 0 {
		return s.Send(ctx, p.Wrap(nil))
	}
	return s.Send(ctx, p.Wrap(xmlstream.MultiReader(readers...)))
}
