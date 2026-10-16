package xmpp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meszmate/xmpp-go/crypto/omemo"
	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/pubsub"
)

// OMEMOIdentity stores the local cryptographic identity and key manager.
type OMEMOIdentity struct {
	DeviceID    uint32
	Fingerprint string
	Manager     *omemo.Manager
	Store       *omemo.MemoryStore
	Bundle      *omemo.Bundle
}

// GetOrCreateDeviceID retrieves or creates a persistent device ID.
func GetOrCreateDeviceID() uint32 {
	home, err := os.UserHomeDir()
	if err != nil {
		return fallbackRandomID()
	}
	dir := filepath.Join(home, ".config", "chat")
	_ = os.MkdirAll(dir, 0o700)
	idPath := filepath.Join(dir, "device_id")
	data, err := os.ReadFile(idPath)
	if err == nil {
		if id, parseErr := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 32); parseErr == nil && id > 0 {
			return uint32(id)
		}
	}
	id := fallbackRandomID()
	_ = os.WriteFile(idPath, []byte(strconv.FormatUint(uint64(id), 10)), 0o600)
	return id
}

func fallbackRandomID() uint32 {
	n, err := rand.Int(rand.Reader, big.NewInt(899999999))
	if err != nil {
		return 133742
	}
	return uint32(n.Int64() + 100000000)
}

// InitOMEMO generates local keys and initializes the OMEMO manager.
func InitOMEMO(deviceID uint32) (*OMEMOIdentity, error) {
	store := omemo.NewMemoryStore(deviceID)
	mgr := omemo.NewManager(store)
	bundle, err := mgr.GenerateBundle(50)
	if err != nil {
		return nil, fmt.Errorf("generate omemo bundle: %w", err)
	}
	rawFingerprint := hex.EncodeToString(bundle.IdentityKey)
	var formatted strings.Builder
	for i := 0; i < len(rawFingerprint); i += 8 {
		if i > 0 {
			formatted.WriteString(" ")
		}
		end := i + 8
		if end > len(rawFingerprint) {
			end = len(rawFingerprint)
		}
		formatted.WriteString(rawFingerprint[i:end])
	}
	return &OMEMOIdentity{
		DeviceID:    deviceID,
		Fingerprint: strings.ToUpper(formatted.String()),
		Manager:     mgr,
		Store:       store,
		Bundle:      bundle,
	}, nil
}

// OMEMOInfo formats the local OMEMO identity status.
func (c *Client) OMEMOInfo() string {
	if c.omemo == nil {
		return "OMEMO: Not initialized"
	}
	return fmt.Sprintf("OMEMO E2EE (XEP-0384):\n  • Device ID: %d\n  • Fingerprint: %s\n  • Pre-keys Generated: %d\n  • Engine: Pure Go (crypto/ecdh, X25519)",
		c.omemo.DeviceID, c.omemo.Fingerprint, len(c.omemo.Bundle.PreKeys))
}

type deviceNodeInfo struct {
	node  string
	root  string
	space string
}

// AnnounceOMEMODevice publishes our device ID and bundle to server PEP.
func (c *Client) AnnounceOMEMODevice(ctx context.Context) error {
	c.mu.Lock()
	s := c.session
	om := c.omemo
	c.mu.Unlock()
	if s == nil || om == nil {
		return fmt.Errorf("not connected or omemo not ready")
	}
	modern := deviceNodeInfo{node: "urn:xmpp:omemo:2:devices", root: "devices", space: "urn:xmpp:omemo:2"}
	legacy := deviceNodeInfo{node: "eu.siacs.conversations.axolotl.devicelist", root: "list", space: "eu.siacs.conversations.axolotl"}
	_ = publishDevice(ctx, s, modern, om.DeviceID)
	_ = publishDevice(ctx, s, legacy, om.DeviceID)
	return publishBundle(ctx, s, om)
}

func publishDevice(ctx context.Context, s *xmpp.Session, info deviceNodeInfo, id uint32) error {
	item := xmlstream.Wrap(
		xmlstream.Wrap(
			nil,
			xml.StartElement{
				Name: xml.Name{Space: info.space, Local: "device"},
				Attr: []xml.Attr{{Name: xml.Name{Local: "id"}, Value: strconv.FormatUint(uint64(id), 10)}},
			},
		),
		xml.StartElement{Name: xml.Name{Space: info.space, Local: info.root}},
	)
	_, err := pubsub.Publish(ctx, s, info.node, "current", item)
	return err
}

func publishBundle(ctx context.Context, s *xmpp.Session, om *OMEMOIdentity) error {
	spkB64 := base64.StdEncoding.EncodeToString(om.Bundle.SignedPreKey)
	sigB64 := base64.StdEncoding.EncodeToString(om.Bundle.SignedPreKeySignature)
	ikB64 := base64.StdEncoding.EncodeToString(om.Bundle.IdentityKey)
	bundleItem := xmlstream.Wrap(
		xmlstream.MultiReader(
			xmlstream.Wrap(xmlstream.Token(xml.CharData(spkB64)), xml.StartElement{
				Name: xml.Name{Space: "urn:xmpp:omemo:2", Local: "spk"},
				Attr: []xml.Attr{{Name: xml.Name{Local: "id"}, Value: strconv.FormatUint(uint64(om.Bundle.SignedPreKeyID), 10)}},
			}),
			xmlstream.Wrap(xmlstream.Token(xml.CharData(sigB64)), xml.StartElement{Name: xml.Name{Space: "urn:xmpp:omemo:2", Local: "spks"}}),
			xmlstream.Wrap(xmlstream.Token(xml.CharData(ikB64)), xml.StartElement{Name: xml.Name{Space: "urn:xmpp:omemo:2", Local: "ik"}}),
		),
		xml.StartElement{Name: xml.Name{Space: "urn:xmpp:omemo:2", Local: "bundle"}},
	)
	itemID := strconv.FormatUint(uint64(om.DeviceID), 10)
	_, err := pubsub.Publish(ctx, s, "urn:xmpp:omemo:2:bundles", itemID, bundleItem)
	return err
}
