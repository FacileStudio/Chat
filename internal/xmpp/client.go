package xmpp

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"sync"
	"time"

	"github.com/FacileStudio/Chat/internal/config"
	"github.com/FacileStudio/Chat/internal/event"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/disco"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/roster"
	"mellium.im/xmpp/stanza"
)

// Client represents the active XMPP client session.
type Client struct {
	cfg     *config.Config
	addr    jid.JID
	events  chan<- any
	session *xmpp.Session
	cancel  context.CancelFunc
	mu      sync.Mutex
	roster  []event.RosterItem
	omemo   *OMEMOIdentity
}

// New creates a new uninitialized XMPP client instance.
func New(cfg *config.Config, events chan<- any) (*Client, error) {
	addr, err := jid.Parse(cfg.JID)
	if err != nil {
		return nil, fmt.Errorf("invalid jid %q: %w", cfg.JID, err)
	}
	devID := GetOrCreateDeviceID()
	om, _ := InitOMEMO(devID)
	return &Client{
		cfg:    cfg,
		addr:   addr,
		events: events,
		roster: make([]event.RosterItem, 0),
		omemo:  om,
	}, nil
}

// Connect establishes the XMPP connection and begins processing stanzas.
func (c *Client) Connect(ctx context.Context) error {
	c.events <- event.StatusEvent{State: event.StateConnecting, Message: "Connecting to server..."}
	conn, err := DialConn(ctx, c.cfg, c.addr)
	if err != nil {
		c.events <- event.StatusEvent{State: event.StateError, Message: err.Error(), Err: err}
		return err
	}
	tlsCfg := CreateTLSConfig(c.cfg, c.addr.Domain().String())
	mechanisms := BuildMechanisms(c.addr.Localpart() != "")
	saslFeature := xmpp.SASL("", c.cfg.Password, mechanisms...)
	negotiator := buildNegotiator(tlsCfg, saslFeature)
	s, err := xmpp.NewSession(ctx, c.addr.Domain(), c.addr, conn, 0, negotiator)
	if err != nil {
		c.events <- event.StatusEvent{State: event.StateError, Message: err.Error(), Err: err}
		return err
	}
	c.session = s
	c.events <- event.StatusEvent{State: event.StateConnected, Message: "Connected"}
	go c.serve()
	_ = SendPresence(ctx, s, "", "")
	_ = FetchRoster(ctx, s, c)
	_ = c.AnnounceOMEMODevice(ctx)
	return nil
}

func buildNegotiator(tlsCfg *tls.Config, sasl xmpp.StreamFeature) xmpp.Negotiator {
	return xmpp.NewNegotiator(func(*xmpp.Session, *xmpp.StreamConfig) xmpp.StreamConfig {
		return xmpp.StreamConfig{
			Features: []xmpp.StreamFeature{
				disco.StreamFeature(),
				xmpp.StartTLS(tlsCfg),
				sasl,
				roster.Versioning(),
				xmpp.BindResource(),
			},
		}
	})
}

func (c *Client) serve() {
	handler := NewXMPPHandler(c)
	err := c.session.Serve(handler)
	c.mu.Lock()
	c.session = nil
	c.mu.Unlock()
	if err != nil {
		c.events <- event.StatusEvent{State: event.StateError, Message: err.Error(), Err: err}
		return
	}
	c.events <- event.StatusEvent{State: event.StateDisconnected, Message: "Disconnected"}
}

// Send encodes and transmits a chat message to the recipient JID.
func (c *Client) Send(toJID, body string) error {
	c.mu.Lock()
	s := c.session
	c.mu.Unlock()
	if s == nil {
		return fmt.Errorf("not connected")
	}
	to, err := jid.Parse(toJID)
	if err != nil {
		return fmt.Errorf("invalid recipient jid: %w", err)
	}
	msg := stanza.Message{To: to, Type: stanza.ChatMessage}
	payload := xmlstream.Wrap(
		xmlstream.Token(xml.CharData(body)),
		xml.StartElement{Name: xml.Name{Local: "body"}},
	)
	if err := s.Send(context.Background(), msg.Wrap(payload)); err != nil {
		return err
	}
	c.events <- event.MessageEvent{
		From: c.addr.Bare().String(),
		To:   to.Bare().String(),
		Body: body,
		Time: time.Now(),
		Sent: true,
	}
	return nil
}

// SetPresence sends a presence update with show and status values.
func (c *Client) SetPresence(show, status string) error {
	c.mu.Lock()
	s := c.session
	c.mu.Unlock()
	if s == nil {
		return fmt.Errorf("not connected")
	}
	return SendPresence(context.Background(), s, show, status)
}

// AddContact adds a new contact to the roster and requests subscription.
func (c *Client) AddContact(targetJID, name string) error {
	c.mu.Lock()
	s := c.session
	c.mu.Unlock()
	if s == nil {
		return fmt.Errorf("not connected")
	}
	target, err := jid.Parse(targetJID)
	if err != nil {
		return err
	}
	item := roster.Item{JID: target.Bare(), Name: name}
	if err := roster.Set(context.Background(), s, item); err != nil {
		return err
	}
	sub := stanza.Presence{To: target.Bare(), Type: stanza.SubscribePresence}
	return s.Send(context.Background(), sub.Wrap(nil))
}

// Close closes the active session and releases resources.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
	if c.session != nil {
		_ = SendPresence(context.Background(), c.session, "offline", "")
		return c.session.Close()
	}
	return nil
}
