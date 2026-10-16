// Package event defines event types exchanged between XMPP and UI.
package event

import "time"

// ConnectionState represents the state of the XMPP connection.
type ConnectionState int

// Connection state constants.
const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateError
)

// StatusEvent signals a connection state transition.
type StatusEvent struct {
	State   ConnectionState
	Message string
	Err     error
}

// MessageEvent represents an incoming or outgoing chat message.
type MessageEvent struct {
	From      string
	To        string
	Body      string
	Time      time.Time
	Sent      bool
	Encrypted bool
}

// RosterItem represents a contact entry in the roster.
type RosterItem struct {
	JID          string
	Name         string
	Subscription string
	Show         string
	Status       string
	Online       bool
}

// RosterEvent delivers the updated roster list.
type RosterEvent struct {
	Items []RosterItem
}

// PresenceEvent delivers an individual presence status update.
type PresenceEvent struct {
	JID    string
	Show   string
	Status string
	Online bool
}
