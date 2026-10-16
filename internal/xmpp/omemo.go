// Package xmpp provides XMPP client and communication handling.
package xmpp

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/pubsub"
)

// IsOMEMOMessage checks whether a message is an OMEMO encrypted message or fallback.
func IsOMEMOMessage(body string, hasEncrypted, hasEME bool) bool {
	if hasEncrypted || hasEME {
		return true
	}
	if strings.Contains(body, "OMEMO encrypted message") {
		return true
	}
	return strings.Contains(body, "conversations.im/omemo")
}

// ClearOMEMODevices publishes empty OMEMO device lists to PEP to reset devices.
func (c *Client) ClearOMEMODevices() error {
	c.mu.Lock()
	s := c.session
	c.mu.Unlock()
	if s == nil {
		return fmt.Errorf("not connected")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := clearPEPNode(ctx, s, "eu.siacs.conversations.axolotl.devicelist", "list", "eu.siacs.conversations.axolotl")
	err2 := clearPEPNode(ctx, s, "urn:xmpp:omemo:2:devicelist", "devices", "urn:xmpp:omemo:2")
	if err1 != nil && err2 != nil {
		return err1
	}
	return nil
}

func clearPEPNode(ctx context.Context, s *xmpp.Session, node, elem, space string) error {
	payload := xmlstream.Wrap(
		nil,
		xml.StartElement{Name: xml.Name{Space: space, Local: elem}},
	)
	_, err := pubsub.Publish(ctx, s, node, "current", payload)
	return err
}
