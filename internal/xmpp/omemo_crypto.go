package xmpp

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/meszmate/xmpp-go/crypto/omemo"
)

// DecryptOMEMO attempts to decrypt an incoming OMEMO message.
func (c *Client) DecryptOMEMO(senderJID string, enc *rawEncrypted) ([]byte, error) {
	if c.omemo == nil || enc == nil {
		return nil, fmt.Errorf("omemo not ready")
	}
	keyData, isPreKey, found := findDeviceKey(enc, c.omemo.DeviceID)
	if !found {
		return nil, fmt.Errorf("no key for device %d", c.omemo.DeviceID)
	}
	encMsg, keyBytes, err := buildEncryptedMessage(enc, c.omemo.DeviceID, keyData, isPreKey)
	if err != nil {
		return nil, err
	}
	addr := omemo.Address{JID: senderJID, DeviceID: enc.Header.SID}
	plaintext, err := c.omemo.Manager.Decrypt(addr, encMsg)
	if err == nil && len(plaintext) > 0 {
		return plaintext, nil
	}
	if isPreKey && len(keyBytes) >= 32 {
		return c.omemo.Manager.DecryptPreKeyMessage(addr, nil, keyBytes[:32], nil, c.omemo.Bundle.SignedPreKeyID, encMsg)
	}
	return nil, err
}

func buildEncryptedMessage(enc *rawEncrypted, devID uint32, keyData string, isPreKey bool) (*omemo.EncryptedMessage, []byte, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(keyData))
	if err != nil {
		return nil, nil, fmt.Errorf("decode key data: %w", err)
	}
	payloadBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(enc.Payload))
	if err != nil {
		return nil, nil, fmt.Errorf("decode payload: %w", err)
	}
	var ivBytes []byte
	if enc.Header.IV != "" {
		ivBytes, _ = base64.StdEncoding.DecodeString(strings.TrimSpace(enc.Header.IV))
	}
	return &omemo.EncryptedMessage{
		SenderDeviceID: enc.Header.SID,
		Keys:           []omemo.MessageKey{{DeviceID: devID, Data: keyBytes, IsPreKey: isPreKey}},
		IV:             ivBytes,
		Payload:        payloadBytes,
	}, keyBytes, nil
}

func findDeviceKey(enc *rawEncrypted, devID uint32) (string, bool, bool) {
	for _, k := range enc.Header.Keys {
		if k.RID == devID {
			return k.Data, k.PreKey || k.Kex, true
		}
	}
	for _, kj := range enc.Header.KeysJID {
		for _, k := range kj.Keys {
			if k.RID == devID {
				return k.Data, k.PreKey || k.Kex, true
			}
		}
	}
	return "", false, false
}
